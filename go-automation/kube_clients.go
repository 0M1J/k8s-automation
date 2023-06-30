package goautomation

import (
	"context"
	"encoding/base64"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/container/v1"
	"google.golang.org/api/option"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

func GetGKEClinetset(cluster *container.Cluster, ts oauth2.TokenSource) (*kubernetes.Clientset, error) {
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

func GcpKubernetesSetup() (*kubernetes.Clientset, error){
    gcpServiceAccountKeyFile := "gcp_service_account_key.json"
    gkeLocation := "us-central1-c"
    gkeClusterName := "kube-cluster"
    // gkeNamespace := "default"

    data, err := ioutil.ReadFile(gcpServiceAccountKeyFile)
    if err != nil {
        log.Fatalf("Failed to read GCP service account key file: %s", err)
        return nil, err
    }
    ctx := context.Background()

    creds, err := google.CredentialsFromJSON(ctx, data, container.CloudPlatformScope)
    if err != nil {
        log.Fatalf("Failed to load GCP service account credintials: %s", err)
        return nil, err
    }

    gkeService, err := container.NewService(ctx, option.WithHTTPClient(oauth2.NewClient(ctx, creds.TokenSource)))
    if err != nil {
        log.Fatalf("Failed to initialize kubernetes engine service: %s", err)
        return nil, err
    }

    name := fmt.Sprintf("projects/%s/locations/%s/clusters/%s", creds.ProjectID, gkeLocation, gkeClusterName)
    cluster, err := container.NewProjectsLocationsClustersService(gkeService).Get(name).Do()
    if err != nil {
        log.Fatalf("Failed to load GKE cluster %q: %s", name, err)
        return nil, err
    }

    clientset, err := GetGKEClinetset(cluster, creds.TokenSource)
    if err != nil {
        log.Fatalf("Failed to initialize kubernetes clientset: %s", err)
        return nil, err
    }

    return clientset, nil
}

func LocalKubernetesSetup() (*kubernetes.Clientset, error) {
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

    return clientset, nil
}


