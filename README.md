# go-automation

## glossary

### Kubernetes: 
- Kubernetes is an open-source container orchestration platform used for automating the deployment, scaling, and management of containerized applications. 
- It provides features such as load balancing, auto-scaling, self-healing, and rolling deployments to ensure efficient and reliable application management.

### Kubeconfig: 
- Kubeconfig is a configuration file used by Kubernetes client tools to specify the connection details and authentication credentials for accessing a Kubernetes cluster. 
- It contains information such as the cluster's API server endpoint, authentication methods (e.g., username and password, client certificates, or bearer tokens), and contextual information like the default namespace and cluster context. 
- The kubeconfig file allows users to manage and switch between different Kubernetes clusters, providing flexibility and control over cluster access and configuration. 
- It is typically located at the ~/.kube/config file path by default, but can be customized using the KUBECONFIG environment variable.

### Kubernetes Clientset: 
- A Kubernetes clientset is a client library in Go that provides a programmatic interface to interact with the Kubernetes API server. 
- It offers methods and functions to perform operations on Kubernetes resources, allowing developers to create, update, delete, and retrieve information about objects within a Kubernetes cluster. 
- The clientset abstracts the underlying communication details, making it easier to interact with the Kubernetes API programmatically.

## [Basic Kubernetes](./docs/kubernetes.md)
