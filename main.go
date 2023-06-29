package main

import (
	"fmt"
    "context"

    metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
)

func main() {
    fmt.Print("Go Automation for Kubernetes deployments...")
    // load kubeconfig from the default config file 
    kubeConfigPath := clientcmd.RecommendedHomeFile

    // build config from config path
    config, err := clientcmd.BuildConfigFromFlags("", kubeConfigPath)
    if err != nil {
        panic(err.Error())
    }

    // create a kubernetes clientset
    clientset, err := kubernetes.NewForConfig(config)
    if err != nil {
        panic(err.Error())
    }
    
    ctx := context.TODO()

    // list the deployments in the default namespace
    deployments, err := clientset.AppsV1().Deployments("default").List(ctx, metav1.ListOptions{})
    if err != nil {
        panic(err.Error())
    }

    fmt.Println("Deployments :")
    for _, deployment := range deployments.Items {
        fmt.Println(deployment.Name)
    }
    
}
