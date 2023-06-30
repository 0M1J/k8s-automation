package main

import (
	"log"

	goautomation "github.com/ozx1812/go-automation/go-automation"
)


func main() {
    // accessing kube cluster through the gcp service acc key
    clientset, err := goautomation.GcpKubernetesSetup()
    if err != nil {
        log.Fatalf("error while setting up gke kubernetes clientset: %v", err)
    }

    // accessing kube cluster through the local kube config file 
    clientset, err = goautomation.LocalKubernetesSetup()
    if err != nil {
        log.Fatalf("error while setting up local kubernetes clientset: %v", err)
    }

    goautomation.ListAllResources(clientset)
    err = goautomation.DeployResourcesFromDirectory(clientset, "kube/nginx")
    if err != nil {
        log.Fatalf("error while doing deployments", err)
    }
}
