package services

import (
    "context"
    "fmt"
    "os"
    "path/filepath"
    "os/exec"
	"time"
    "strings"
    "encoding/json"
    
    "leetcodeBackend/repositories"
    "k8s.io/client-go/kubernetes"
    "k8s.io/client-go/rest"
    corev1 "k8s.io/api/core/v1"
    metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
    "k8s.io/apimachinery/pkg/api/resource"
)

type ExecutionService struct {
    clientset *kubernetes.Clientset
    questionRepo repositories.QuestionRepositoryInterface
}

func NewExecutionService(questionRepo repositories.QuestionRepositoryInterface) (*ExecutionService, error) {
    if os.Getenv("KUBERNETES_SERVICE_HOST") == "" {
        return &ExecutionService{
            questionRepo: questionRepo,
        }, nil
    }

    config, err := rest.InClusterConfig()
    if err != nil {
        return nil, fmt.Errorf("failed to get k8s config: %v", err)
    }

    clientset, err := kubernetes.NewForConfig(config)
    if err != nil {
        return nil, fmt.Errorf("failed to create k8s clientset: %v", err)
    }

    return &ExecutionService{
        clientset: clientset,
        questionRepo: questionRepo,
    }, nil
}

func (es *ExecutionService) CreateDockerImage(code string, language string, questionId int) (string, error) {
    question, err := es.questionRepo.GetQuestionByID(fmt.Sprintf("%d", questionId))
    if err != nil {
        return "", fmt.Errorf("failed to get question: %v", err)
    }

    tmpDir, err := os.MkdirTemp("", "code-execution")
    if err != nil {
        return "", fmt.Errorf("failed to create temp dir: %v", err)
    }
    defer os.RemoveAll(tmpDir)

    var codeContent string
    if language == "python" {
        codeContent = fmt.Sprintf(`
import sys
import json

%s

if __name__ == "__main__":
    input_data = json.loads(sys.argv[1])
    args = input_data["args"]
    result = %s(*args)
    print(json.dumps(result))
`, code, question.Title)
    } else if language == "go" {
        codeContent = fmt.Sprintf(`package main

import (
    "fmt"
    "encoding/json"
    "os"
)

%s

func main() {
    input := os.Args[1]
    var data struct {
        Args []interface{} ` + "`json:\"args\"`" + `
    }
    
    if err := json.Unmarshal([]byte(input), &data); err != nil {
        fmt.Printf("Error parsing input: %%v\n", err)
        os.Exit(1)
    }

    args := make([]reflect.Value, len(data.Args))
    for i, arg := range data.Args {
        args[i] = reflect.ValueOf(arg)
    }

    result := reflect.ValueOf(%s).Call(args)

    output, _ := json.Marshal(result[0].Interface())
    fmt.Print(string(output))
}`, code, question.Title)
    }

    codeFile := "main"
    if language == "python" {
        codeFile += ".py"
    } else {
        codeFile += ".go"
    }
    
    if err := os.WriteFile(filepath.Join(tmpDir, codeFile), []byte(codeContent), 0644); err != nil {
        return "", fmt.Errorf("failed to write code file: %v", err)
    }

    var dockerfileContent string
    if language == "python" {
        dockerfileContent = `FROM python:3.9-alpine
COPY main.py /app/main.py
WORKDIR /app
ENTRYPOINT ["python", "main.py"]`
    } else {
        dockerfileContent = `FROM golang:1.21-alpine
COPY main.go /app/main.go
WORKDIR /app
RUN go mod init usercode && go build -o main
ENTRYPOINT ["./main"]`
    }

    if err := os.WriteFile(filepath.Join(tmpDir, "Dockerfile"), []byte(dockerfileContent), 0644); err != nil {
        return "", fmt.Errorf("failed to write Dockerfile: %v", err)
    }

    imageName := fmt.Sprintf("user-code-%s:%s", language, filepath.Base(tmpDir))
    buildCmd := exec.Command("docker", "build", "-t", imageName, tmpDir)
    if output, err := buildCmd.CombinedOutput(); err != nil {
        return "", fmt.Errorf("failed to build docker image: %v, output: %s", err, string(output))
    }

    return imageName, nil
}

func (es *ExecutionService) RunOnKubernetes(imageName string, args []interface{}) (string, error) {
    if es.clientset == nil {
        containerName := fmt.Sprintf("user-code-%d", time.Now().Unix())
        
        inputData := map[string]interface{}{
            "args": args,
        }
        
        jsonArgs, err := json.Marshal(inputData)
        if err != nil {
            return "", fmt.Errorf("failed to marshal arguments: %v", err)
        }

        //  Run the container
        runCmd := exec.Command("docker", "run",
            "--name", containerName,
            "--rm",
            imageName,
            string(jsonArgs),
        )

        output, err := runCmd.CombinedOutput()
        if err != nil {
            return "", fmt.Errorf("failed to run docker container: %v, output: %s", err, string(output))
        }

        return strings.TrimSpace(string(output)), nil
    }

    // Create pod specification for Kubernetes
    pod := &corev1.Pod{
        ObjectMeta: metav1.ObjectMeta{
            Name: fmt.Sprintf("user-code-%d", time.Now().Unix()),
        },
        Spec: corev1.PodSpec{
            RestartPolicy: corev1.RestartPolicyNever,
            Containers: []corev1.Container{
                {
                    Name:  "user-code",
                    Image: imageName,
                    Args:  []string{fmt.Sprintf("%v", args[0])},
                    Resources: corev1.ResourceRequirements{
                        Limits: corev1.ResourceList{
                            corev1.ResourceCPU:    resource.MustParse("100m"),
                            corev1.ResourceMemory: resource.MustParse("100Mi"),
                        },
                    },
                    SecurityContext: &corev1.SecurityContext{
                        AllowPrivilegeEscalation: &[]bool{false}[0],
                        RunAsNonRoot:            &[]bool{true}[0],
                        RunAsUser:               &[]int64{1000}[0],
                    },
                },
            },
        },
    }

    // Create pod in Kubernetes cluster
    createdPod, err := es.clientset.CoreV1().Pods("default").Create(context.Background(), pod, metav1.CreateOptions{})
    if err != nil {
        return "", fmt.Errorf("failed to create pod: %v", err)
    }

    // Ensure pod cleanup after execution
    defer func() {
        es.clientset.CoreV1().Pods("default").Delete(context.Background(), createdPod.Name, metav1.DeleteOptions{})
    }()

    // Watch pod status
    watch, err := es.clientset.CoreV1().Pods("default").Watch(context.Background(), metav1.ListOptions{
        FieldSelector: fmt.Sprintf("metadata.name=%s", createdPod.Name),
    })
    if err != nil {
        return "", fmt.Errorf("failed to watch pod: %v", err)
    }
    defer watch.Stop()

    // Get execution output from pod logs
    logs, err := es.clientset.CoreV1().Pods("default").
        GetLogs(createdPod.Name, &corev1.PodLogOptions{}).
        DoRaw(context.Background())
    if err != nil {
        return "", fmt.Errorf("failed to get logs: %v", err)
    }

    return string(logs), nil
}