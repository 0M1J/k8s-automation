package goautomation

import (
	"context"
	"fmt"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
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

