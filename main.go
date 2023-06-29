package main

import (
	"context"
	"encoding/base64"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	goautomation "github.com/ozx1812/go-automation/go-automation"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/container/v1"
	"google.golang.org/api/option"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

func getGKEClinetset(cluster *container.Cluster, ts oauth2.TokenSource) (kubernetes.Interface, error) {
    capem, err := base64.StdEncoding.DecodeString(cluster.MasterAuth.ClusterCaCertificate)
    if err != nil {
        return nil, fmt.Errorf("failed to decode cluster CA cert: %s", err)
    }

    config := &rest.Config{
        Host: cluster.Endpoint,
        TLSClientConfig: rest.TLSClientConfig{
            CAData: capem,
        },
    }
    config.Wrap(func(rt http.RoundTripper) http.RoundTripper {
        return &oauth2.Transport{
            Source: ts,
            Base: rt,
        }
    }) 

    clientset, err := kubernetes.NewForConfig(config)
    if err != nil {
        return nil, fmt.Errorf("failed to initialize clientset from config: %s", err) 
    }
    return clientset, nil
}

func gcpKubernetesSetup() {
    gcpServiceAccountKeyFile := "gcp_service_account_key.json"
    gkeLocation := "us-central1-c"
    gkeClusterName := "kube-cluster"
    gkeNamespace := "default"

    data, err := ioutil.ReadFile(gcpServiceAccountKeyFile)
    if err != nil {
        log.Fatalf("Failed to read GCP service account key file: %s", err)
    }
    ctx := context.Background()

    creds, err := google.CredentialsFromJSON(ctx, data, container.CloudPlatformScope)
    if err != nil {
        log.Fatalf("Failed to load GCP service account credintials: %s", err)
    }

    gkeService, err := container.NewService(ctx, option.WithHTTPClient(oauth2.NewClient(ctx, creds.TokenSource)))
    if err != nil {
        log.Fatalf("Failed to initialize kubernetes engine service: %s", err)
    }

    name := fmt.Sprintf("projects/%s/locations/%s/clusters/%s", creds.ProjectID, gkeLocation, gkeClusterName)
    cluster, err := container.NewProjectsLocationsClustersService(gkeService).Get(name).Do()
    if err != nil {
        log.Fatalf("Failed to load GKE cluster %q: %s", name, err)
    }

    clientset, err := getGKEClinetset(cluster, creds.TokenSource)
    if err != nil {
        log.Fatalf("Failed to initialize kubernetes clientset: %s", err)
    }

    pods, err := clientset.CoreV1().Pods(gkeNamespace).List(ctx, metav1.ListOptions{})
    if err != nil {
        log.Fatalf("Failed to list pods: %s", err)
    }
    log.Printf("There are %d pods in the namespace", len(pods.Items))

}

func localKubernetesSetup() {
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
    // testing list resource

    // list the deployments in the default namespace
    deployments, err := goautomation.ListResources(clientset, "Deployment")
    if err != nil {
        panic(err.Error())
    }
    fmt.Println("Deployments :")
    for _, deployment := range deployments{
        fmt.Println(deployment)
    }
    
    // list the services in the default namespace
    services, err := goautomation.ListResources(clientset, "Service")
    if err != nil {
        panic(err.Error())
    }
    fmt.Println("Services:")
    for _, service:= range services{
        fmt.Println(service)
    }

    // list the configmaps in the default namespace
    configmaps, err := goautomation.ListResources(clientset, "ConfigMap")
    if err != nil {
        panic(err.Error())
    }
    fmt.Println("Configmaps:")
    for _, configmap:= range configmaps{
        fmt.Println(configmap)
    }

    // create deployment from yaml file
    yamlFile := "nginx-deployment.yaml"
    err = goautomation.CreateResources(clientset, yamlFile)
    if err != nil {
        panic(err.Error())
    }


}

func main() {
    // accessing kube cluster through the gcp service acc key
    gcpKubernetesSetup()

    // accessing kube cluster through the local kube config file 
    localKubernetesSetup()
}
