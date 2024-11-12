// services/execution_service.go
package services

import (
    "context"
    "errors"
    "fmt"
    "io"
    "os"
    "path/filepath"
    "time"

    v1 "k8s.io/api/core/v1"
    metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
    "k8s.io/client-go/kubernetes"
    "k8s.io/client-go/rest"
)

type ExecuteRequest struct {
    Code         string `json:"code"`
    Language     string `json:"language"`
    QuestionCode int    `json:"questionCode"`
}

type ExecutionService struct {
    KubeClient *kubernetes.Clientset
}

func NewExecutionService() (*ExecutionService, error) {
    config, err := rest.InClusterConfig()
    if err != nil {
        return nil, err
    }
    kubeClient, err := kubernetes.NewForConfig(config)
    if err != nil {
        return nil, err
    }

    return &ExecutionService{
        KubeClient: kubeClient,
    }, nil
}

func (service *ExecutionService) Execute(req ExecuteRequest) (string, error) {
    println("3")
    codeDir := filepath.Join("running_user_code", req.Language)
    codeFile := ""
    if req.Language == "go" {
        codeFile = filepath.Join(codeDir, "user_code.go")
    } else if req.Language == "python" {
        codeFile = filepath.Join(codeDir, "user_code.py")
    }

    println("4")
    println(codeFile)
    if err := os.WriteFile(codeFile, []byte(req.Code), 0644); err != nil {
        return "", errors.New("failed to write code to file")
    }

    println("5")
    podName := fmt.Sprintf("user-code-%d", time.Now().UnixNano())
    pod := &v1.Pod{
        ObjectMeta: metav1.ObjectMeta{
            Name: podName,
        },
        Spec: v1.PodSpec{
            Containers: []v1.Container{
                {
                    Name:            "user-code",
                    Image:           fmt.Sprintf("user_code_%s", req.Language),
                    ImagePullPolicy: v1.PullIfNotPresent,
                    Command:         []string{"/bin/sh", "-c", "cat /code/user_code.go | go run"},
                    VolumeMounts: []v1.VolumeMount{
                        {
                            Name:      "code",
                            MountPath: "/code",
                        },
                    },
                },
            },
            Volumes: []v1.Volume{
                {
                    Name: "code",
                    VolumeSource: v1.VolumeSource{
                        HostPath: &v1.HostPathVolumeSource{
                            Path: codeDir,
                        },
                    },
                },
            },
            RestartPolicy: v1.RestartPolicyNever,
        },
    }

    ctx := context.Background()

    // Create the pod in Kubernetes
    _, err := service.KubeClient.CoreV1().Pods("default").Create(ctx, pod, metav1.CreateOptions{})
    if err != nil {
        return "", fmt.Errorf("error creating pod: %v", err)
    }

    // Wait for the pod to finish
    for {
        waitPod, err := service.KubeClient.CoreV1().Pods("default").Get(ctx, podName, metav1.GetOptions{})
        if err != nil {
            return "", fmt.Errorf("error getting pod: %v", err)
        }

        if waitPod.Status.Phase == v1.PodSucceeded || waitPod.Status.Phase == v1.PodFailed {
            break
        }
        time.Sleep(1 * time.Second)
    }

    // Get the logs from the pod
    logs, err := service.KubeClient.CoreV1().Pods("default").GetLogs(podName, &v1.PodLogOptions{}).Stream(ctx)
    if err != nil {
        return "", fmt.Errorf("error getting pod logs: %v", err)
    }
    defer logs.Close()

    output, err := io.ReadAll(logs)
    if err != nil {
        return "", fmt.Errorf("error reading pod logs: %v", err)
    }

    // Cleanup: Delete the pod
    err = service.KubeClient.CoreV1().Pods("default").Delete(ctx, podName, metav1.DeleteOptions{})
    if err != nil {
        fmt.Printf("warning: failed to delete pod %s: %v\n", podName, err)
    }

    return string(output), nil
}
