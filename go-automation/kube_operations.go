package goautomation

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"

	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/kubernetes"
)

func ListDeployments(clientset *kubernetes.Clientset) ([]string, error) {
    deployments, err := clientset.AppsV1().Deployments("default").List(context.TODO(), metav1.ListOptions{})
    if err != nil {
        return nil, err
    }
    var deploymentNames []string
    for _, deployment := range deployments.Items {
        deploymentNames = append(deploymentNames, deployment.Name)
    }
    return deploymentNames, nil
}

func ListConfigmaps(clientset *kubernetes.Clientset) ([]string, error) {
    configmaps, err := clientset.CoreV1().ConfigMaps("default").List(context.TODO(), metav1.ListOptions{})
    if err != nil {
        return nil, err
    }
    var configmapNames []string
    for _, configmap := range configmaps.Items {
        configmapNames = append(configmapNames, configmap.Name)
    }
    return configmapNames, nil
}

func ListServices(clientset *kubernetes.Clientset) ([]string, error) {
    services, err := clientset.CoreV1().Services("default").List(context.TODO(), metav1.ListOptions{})
    if err != nil {
        return nil, err
    }
    var serviceNames []string
    for _, service := range services.Items {
        serviceNames = append(serviceNames, service.Name)
    }
    return serviceNames, nil
}

func ListResources(clientset *kubernetes.Clientset, resourceKind string) ([]string, error) {
    switch resourceKind {
    case "Deployment":
        deployments, err := ListDeployments(clientset)
        if err != nil {
            return nil, err
        }
        return deployments, nil
    case "ConfigMap":
        configmaps, err := ListConfigmaps(clientset)
        if err != nil {
            return nil, err
        }
        return configmaps, nil
    case "Service":
        services, err := ListServices(clientset)
        if err != nil {
            return nil, err
        }
        return services, nil
    default:
        return nil, fmt.Errorf("unsupported resourceKind: %s", resourceKind)
    }
}

func CreateDeployment(clientset *kubernetes.Clientset, deployment *appsv1.Deployment) (error) {
    _, err := clientset.AppsV1().Deployments("default").Create(context.TODO(), deployment, metav1.CreateOptions{})
    if err != nil {
        return err
    }
    return nil
}

func CreateConfigmap(clientset *kubernetes.Clientset, configmap *corev1.ConfigMap) (error) {
    _, err := clientset.CoreV1().ConfigMaps("default").Create(context.TODO(), configmap, metav1.CreateOptions{})
    if err != nil {
        return err
    }
    return nil
}

func CreateService(clientset *kubernetes.Clientset, service *corev1.Service) (error) {
    _, err := clientset.CoreV1().Services("default").Create(context.TODO(), service, metav1.CreateOptions{})
    if err != nil {
        return err
    }
    return nil
}

func CreateResources(clientset *kubernetes.Clientset, yamlPath string) (error) {
    var yamlObjects []*unstructured.Unstructured
    yamlObjects, err := ParseYaml(yamlPath)
    if err != nil {
        return err
    }
    for _, obj := range yamlObjects {
        resourceKind := obj.GetKind()
        switch resourceKind {
        case "Deployment":
            deployment := &appsv1.Deployment{}
            err := runtime.DefaultUnstructuredConverter.FromUnstructured(obj.Object, deployment)
            if err != nil {
                return err
            }
            err = CreateDeployment(clientset, deployment)
            if err != nil {
                return err
            }
        case "ConfigMap":
            configmap := &corev1.ConfigMap{}
            err := runtime.DefaultUnstructuredConverter.FromUnstructured(obj.Object, configmap)
            if err != nil {
                return err
            }
            err = CreateConfigmap(clientset, configmap)
            if err != nil {
                return err
            }
        case "Service":
            service := &corev1.Service{}
            err := runtime.DefaultUnstructuredConverter.FromUnstructured(obj.Object, service)
            if err != nil {
                return err
            }
            err = CreateService(clientset, service)
            if err != nil {
                return err
            }
        default:
            return fmt.Errorf("unsupported resourceKind: %s", resourceKind)
        }
    }
    return nil
}

func getResourceKind(path string) string {
    dir := filepath.Dir(path)

    parts := strings.Split(dir, string(os.PathSeparator))
    kind := parts[len(parts)-1]
    return kind
}

func getOrderIndex(kind string) int {
    deploymentOrder := map[string]int{
        "configmap": 0,
        "deployment": 1,
        "service": 2,
    }
    return deploymentOrder[kind]
}

func DeployResourcesFromDirectory(clientset *kubernetes.Clientset, kubeDir string) error {
    deploymentOrder := []string{"configmap", "deployment", "service"}
    resourceFiles := make(map[string][]string)
    err := filepath.Walk(kubeDir, func(path string, info os.FileInfo, err error) error {
        if err != nil {
            return err
        }

        if !info.IsDir() {
            resourceKind := getResourceKind(path)
            resourceFiles[resourceKind] = append(resourceFiles[resourceKind], path)
        }
        return nil
    })
    if err != nil {
        return err
    }

    sort.SliceStable(deploymentOrder, func(i, j int) bool {
        return getOrderIndex(deploymentOrder[i]) < getOrderIndex(deploymentOrder[j])
    })

    for _, kind := range deploymentOrder{
        files := resourceFiles[kind]
        for _, file := range files {
            fmt.Printf("file : %s", file)
            err = CreateResources(clientset, file)
            if err != nil {
                return err
            }
        }
    }
    return nil
}
