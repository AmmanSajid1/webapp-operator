# WebApp Kubernetes Operator

A custom Kubernetes Operator built using Go and Kubebuilder that extends the Kubernetes API with a new resource type (`WebApp`) and automatically manages application deployments based on declarative specifications.

---

## 🚀 Overview

This project demonstrates how to extend Kubernetes using:

- Custom Resource Definitions (CRDs) to define new resource types  
- Custom Controllers to implement reconciliation logic  
- The Operator Pattern to automate infrastructure management  

Instead of manually creating Deployments, users can define a simple `WebApp` resource:

```yaml
apiVersion: apps.amman.dev/v1
kind: WebApp
metadata:
  name: my-webapp
spec:
  image: nginx:latest
  replicas: 2
```

The operator automatically creates and manages the corresponding Kubernetes Deployment.

---

## 🧠 How It Works

WebApp (Custom Resource) <br>
            ↓ <br>
Controller (Reconciliation Loop) <br>
            ↓ <br>
Deployment (Managed Resource) <br>
            ↓ <br>
Pods

The controller continuously ensures that the actual cluster state matches the desired state defined in the `WebApp` spec.

---

## ⚙️ Features

- Defines a custom Kubernetes resource: `WebApp`
- Automatically creates a Deployment based on:
  - Container image
  - Replica count
- Uses owner references for automatic cleanup
- Implements a basic reconciliation loop

---

## 🛠️ Tech Stack

- Go (Golang)  
- Kubebuilder  
- controller-runtime  
- Kubernetes  

---

## 📦 Getting Started

### Prerequisites

- Go v1.24+  
- Docker  
- kubectl  
- A running Kubernetes cluster (e.g. Minikube)  

---

### 1. Install CRDs

```sh
make install
```

---

### 2. Run the controller locally

```sh
make run
```

---

### 3. Apply sample resource

```sh
kubectl apply -f config/samples/apps_v1_webapp.yaml
```

---

### 4. Verify resources

```sh
kubectl get webapps
kubectl get deployments
kubectl get pods
```

You should see a Deployment created automatically.

---

## 🧪 Example

```yaml
apiVersion: apps.amman.dev/v1
kind: WebApp
metadata:
  name: example
spec:
  image: nginx:latest
  replicas: 2
```

---

## 🧹 Cleanup

```sh
kubectl delete -f config/samples/apps_v1_webapp.yaml
make uninstall
```

---

## 🎯 Learning Goals

This project was built to demonstrate:

- How Kubernetes can be extended beyond native resources  
- How controllers implement reconciliation loops  
- How to build custom operators using Go  

## 📄 License

Apache 2.0