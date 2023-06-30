package goautomation

import (
    "fmt"
	"k8s.io/client-go/kubernetes"
)

func ListAllResources(clientset *kubernetes.Clientset) {
    // testing list resource

    // list the deployments in the default namespace
    deployments, err := ListResources(clientset, "Deployment")
    if err != nil {
        panic(err.Error())
    }
    fmt.Println("Deployments :")
    for _, deployment := range deployments{
        fmt.Println(deployment)
    }
    
    // list the services in the default namespace
    services, err := ListResources(clientset, "Service")
    if err != nil {
        panic(err.Error())
    }
    fmt.Println("Services:")
    for _, service:= range services{
        fmt.Println(service)
    }

    // list the configmaps in the default namespace
    configmaps, err := ListResources(clientset, "ConfigMap")
    if err != nil {
        panic(err.Error())
    }
    fmt.Println("Configmaps:")
    for _, configmap:= range configmaps{
        fmt.Println(configmap)
    }
}
