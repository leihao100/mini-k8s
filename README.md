# 🚀 Mini-Kubernetes (mini-k8s)

A lightweight Kubernetes-inspired container orchestration system built from scratch using Go, Docker, and Nginx.  
This project implements core Kubernetes concepts such as Pods, Services, Deployments, and Horizontal Pod Autoscaling (HPA).

---

## 📌 Overview

Mini-k8s is designed to simulate the architecture and workflow of Kubernetes while remaining simple and easy to understand.  
It provides a minimal yet functional orchestration platform for managing containerized applications.

---

## 🏗️ Key Features

### 📦 Container Orchestration
- Pod abstraction for container grouping
- Deployment management for scaling and updates
- Service abstraction for load balancing

### 🌐 Networking
- Intra-pod communication
- Inter-pod communication
- Cross-node communication
- DNS-based service discovery

### ⚖️ Load Balancing
- Nginx-based traffic distribution
- Service-level routing

### 📈 Auto Scaling
- Horizontal Pod Autoscaler (HPA)
- Dynamic scaling based on workload

### 🔄 Fault Tolerance
- System state persistence
- Recovery after node or system restart

---

## 🧰 Tech Stack

- **Language:** Go  
- **Container Runtime:** Docker  
- **Networking / Load Balancer:** Nginx  
- **System Concepts:** Kubernetes architecture  

---

## 📂 Project Structure

```
mini-k8s/
│
├── cmd/                # Entry points for components
├── pkg/                # Core logic (scheduler, controller, etc.)
├── api/                # API definitions
├── configs/            # Configuration files
├── scripts/            # Deployment scripts
└── docs/               # Documentation
```

---

## ⚙️ Getting Started

### 1️⃣ Prerequisites
- Go installed (>= 1.18)
- Docker installed and running
- Nginx installed

---

### 2️⃣ Clone Repository
```bash
git clone <your-repo-url>
cd mini-k8s
```

---

### 3️⃣ Build Project
```bash
go build ./...
```

---

### 4️⃣ Run System
```bash
./mini-k8s
```

---

## 🔧 Core Components

- **API Server** – Handles REST requests and cluster state  
- **Scheduler** – Assigns pods to nodes  
- **Controller Manager** – Maintains desired state  
- **Kubelet (Simulated)** – Manages containers on nodes  

---

## 📊 Highlights

- Designed a simplified Kubernetes control plane  
- Implemented container lifecycle management  
- Built custom scheduler and controller logic  
- Achieved reliable service networking and discovery  
- Ensured system recovery after crashes  

---

## 📌 Future Improvements

- Add Web UI dashboard  
- Implement RBAC authorization  
- Improve scheduling algorithms  
- Support rolling updates  
- Integrate monitoring (Prometheus/Grafana)  


