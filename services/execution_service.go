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

// ExecuteRequest holds code execution details.
type ExecuteRequest struct {
    Code         string `json:"code"`
    Language     string `json:"language"`
    QuestionCode int    `json:"questionCode"`
}

// ExecutionService provides methods to execute user code within Kubernetes.
type ExecutionService struct {
    KubeClient *kubernetes.Clientset
}

// NewExecutionService initializes a new ExecutionService.
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

// Execute runs user code in a Kubernetes pod and returns the output or any errors.
func (service *ExecutionService) Execute(req ExecuteRequest) (string, error) {
    codeDir := filepath.Join("running_user_code", req.Language)
    codeFile := filepath.Join(codeDir, "user_code."+getFileExtension(req.Language))
    if err := os.WriteFile(codeFile, []byte(req.Code), 0644); err != nil {
        return "", errors.New("failed to write code to file")
    }

    if os.Getenv("KUBERNETES_SERVICE_HOST") == "" {
        return "", errors.New("Kubernetes environment not detected; skipping pod creation")
    }

    podName := fmt.Sprintf("user-code-%d", time.Now().UnixNano())
    configMap := &v1.ConfigMap{
        ObjectMeta: metav1.ObjectMeta{
            Name: podName,
        },
        Data: map[string]string{
            "user_code": req.Code,
        },
    }

    ctx := context.Background()
    if _, err := service.KubeClient.CoreV1().ConfigMaps("default").Create(ctx, configMap, metav1.CreateOptions{}); err != nil {
        return "", fmt.Errorf("error creating ConfigMap: %v", err)
    }

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
                    Command: []string{
                        "/bin/sh",
                        "-c",
                        fmt.Sprintf("cat /config/user_code | %s", getRunCommand(req.Language)),
                    },
                    VolumeMounts: []v1.VolumeMount{
                        {
                            Name:      "config-volume",
                            MountPath: "/config",
                            SubPath:   "user_code",
                        },
                    },
                },
            },
            Volumes: []v1.Volume{
                {
                    Name: "config-volume",
                    VolumeSource: v1.VolumeSource{
                        ConfigMap: &v1.ConfigMapVolumeSource{
                            LocalObjectReference: v1.LocalObjectReference{Name: podName},
                        },
                    },
                },
            },
            RestartPolicy: v1.RestartPolicyNever,
        },
    }

    if _, err := service.KubeClient.CoreV1().Pods("default").Create(ctx, pod, metav1.CreateOptions{}); err != nil {
        return "", fmt.Errorf("error creating pod: %v", err)
    }

    defer func() {
        _ = service.KubeClient.CoreV1().Pods("default").Delete(ctx, podName, metav1.DeleteOptions{})
        _ = service.KubeClient.CoreV1().ConfigMaps("default").Delete(ctx, podName, metav1.DeleteOptions{})
    }()

    if err := service.waitForPodCompletion(ctx, podName); err != nil {
        return "", err
    }

    logs, err := service.KubeClient.CoreV1().Pods("default").GetLogs(podName, &v1.PodLogOptions{}).Stream(ctx)
    if err != nil {
        return "", fmt.Errorf("error getting pod logs: %v", err)
    }
    defer logs.Close()

    output, err := io.ReadAll(logs)
    if err != nil {
        return "", fmt.Errorf("error reading pod logs: %v", err)
    }

    return string(output), nil
}

// getRunCommand returns the command to run code based on the language.
func getRunCommand(language string) string {
    switch language {
    case "go":
        return "go run"
    case "python":
        return "python3"
    default:
        return ""
    }
}

// getFileExtension returns the file extension for the specified language.
func getFileExtension(language string) string {
    switch language {
    case "go":
        return "go"
    case "python":
        return "py"
    default:
        return ""
    }
}

// waitForPodCompletion waits until the specified pod completes its execution.
func (service *ExecutionService) waitForPodCompletion(ctx context.Context, podName string) error {
    for {
        pod, err := service.KubeClient.CoreV1().Pods("default").Get(ctx, podName, metav1.GetOptions{})
        if err != nil {
            return fmt.Errorf("error getting pod: %v", err)
        }
        if pod.Status.Phase == v1.PodSucceeded || pod.Status.Phase == v1.PodFailed {
            return nil
        }
        time.Sleep(1 * time.Second)
    }
}
