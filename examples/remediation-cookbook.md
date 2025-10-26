# Kubernetes Remediation Cookbook

A comprehensive guide to common Kubernetes issues and their fixes, powered by kubectl-pilot.

## Table of Contents

1. [Pod Issues](#pod-issues)
2. [Deployment Issues](#deployment-issues)
3. [Service Issues](#service-issues)
4. [Storage Issues](#storage-issues)
5. [Network Issues](#network-issues)
6. [Resource Issues](#resource-issues)

## Pod Issues

### 1. CrashLoopBackOff

**Symptoms**: Pod continuously crashes and restarts

**Diagnosis**:
```bash
kubectl-pilot diagnose pod <pod-name>
```

**Common Causes**:
- Application startup failure
- Missing environment variables
- Database connection issues
- Out of memory (OOMKilled)

**Remediations**:

```bash
# Check logs for errors
kubectl logs <pod-name> --previous

# Describe pod for events
kubectl describe pod <pod-name>

# Check resource limits
kubectl-pilot run "check if pod has resource limits"

# Increase memory if OOMKilled
kubectl set resources deployment <deployment> --limits=memory=1Gi

# Add missing environment variables
kubectl set env deployment/<deployment> NEW_VAR=value
```

### 2. ImagePullBackOff

**Symptoms**: Pod cannot pull container image

**Diagnosis**:
```bash
kubectl-pilot explain "why can't my pod pull the image"
```

**Common Causes**:
- Image doesn't exist
- Private registry authentication failure
- Network connectivity issues
- Typo in image name

**Remediations**:

```bash
# Verify image exists
docker pull <image-name>

# Check image pull secrets
kubectl get secrets -n <namespace>

# Create image pull secret
kubectl create secret docker-registry regcred \
  --docker-server=<registry> \
  --docker-username=<username> \
  --docker-password=<password>

# Add secret to service account
kubectl patch serviceaccount default -p '{"imagePullSecrets": [{"name": "regcred"}]}'

# Fix image name
kubectl set image deployment/<deployment> container=<correct-image>
```

### 3. Pending Pod

**Symptoms**: Pod stuck in Pending state

**Diagnosis**:
```bash
kubectl-pilot diagnose pod <pod-name>
```

**Common Causes**:
- Insufficient cluster resources
- Node selector doesn't match any node
- PVC not bound
- Taints/tolerations mismatch

**Remediations**:

```bash
# Check node resources
kubectl top nodes

# Check PVC status
kubectl get pvc

# Remove restrictive node selectors
kubectl patch deployment <deployment> --type=json \
  -p='[{"op": "remove", "path": "/spec/template/spec/nodeSelector"}]'

# Scale cluster (cloud-specific)
# GKE: gcloud container clusters resize <cluster> --num-nodes=3
# EKS: eksctl scale nodegroup --cluster=<cluster> --nodes=3
```

### 4. Evicted Pod

**Symptoms**: Pod evicted due to resource pressure

**Diagnosis**:
```bash
kubectl-pilot run "show evicted pods"
```

**Remediations**:

```bash
# Delete evicted pods
kubectl delete pod --field-selector=status.phase==Failed

# Add resource requests/limits
kubectl set resources deployment <deployment> \
  --requests=cpu=100m,memory=128Mi \
  --limits=cpu=500m,memory=512Mi

# Increase node capacity
kubectl-pilot run "recommend scaling strategy for cluster"
```

## Deployment Issues

### 5. Deployment Not Rolling Out

**Symptoms**: New pods not being created

**Diagnosis**:
```bash
kubectl rollout status deployment/<deployment>
kubectl-pilot diagnose deployment <deployment>
```

**Remediations**:

```bash
# Check rollout history
kubectl rollout history deployment/<deployment>

# Force restart
kubectl rollout restart deployment/<deployment>

# Rollback if needed
kubectl rollout undo deployment/<deployment>

# Check for resource quotas
kubectl describe resourcequota -n <namespace>
```

### 6. Old ReplicaSets Not Scaling Down

**Symptoms**: Multiple ReplicaSets with replicas

**Remediations**:

```bash
# Set revision history limit
kubectl patch deployment <deployment> -p '{"spec":{"revisionHistoryLimit":3}}'

# Manually scale down old ReplicaSets
kubectl scale replicaset <old-rs> --replicas=0
```

## Service Issues

### 7. Service Not Accessible

**Symptoms**: Cannot access service from inside/outside cluster

**Diagnosis**:
```bash
kubectl-pilot explain "why can't I access my service"
```

**Remediations**:

```bash
# Check service endpoints
kubectl get endpoints <service>

# Verify selectors match pods
kubectl get pods -l <selector>

# Test from inside cluster
kubectl run test --image=busybox --rm -it -- wget -O- http://<service>:<port>

# Check network policies
kubectl get networkpolicies

# Fix service selector
kubectl patch service <service> -p '{"spec":{"selector":{"app":"correct-label"}}}'
```

### 8. LoadBalancer Service Pending

**Symptoms**: External IP shows <pending>

**Remediations**:

```bash
# Check cloud provider quotas
kubectl-pilot run "check cloud provider limits"

# Verify cloud controller is running
kubectl get pods -n kube-system | grep cloud-controller

# Switch to NodePort temporarily
kubectl patch service <service> -p '{"spec":{"type":"NodePort"}}'
```

## Storage Issues

### 9. PVC Pending

**Symptoms**: PersistentVolumeClaim not binding

**Diagnosis**:
```bash
kubectl describe pvc <pvc-name>
```

**Remediations**:

```bash
# Check available PVs
kubectl get pv

# Check storage class
kubectl get storageclass

# Create PV manually if needed
kubectl apply -f pv.yaml

# Check provisioner logs
kubectl logs -n kube-system <provisioner-pod>

# Use existing storage class
kubectl patch pvc <pvc> -p '{"spec":{"storageClassName":"gp2"}}'
```

### 10. Volume Mount Failures

**Symptoms**: Pod cannot mount volume

**Remediations**:

```bash
# Check PVC status
kubectl get pvc

# Verify node has volume driver
kubectl get nodes -o wide

# Check CSI driver pods
kubectl get pods -n kube-system | grep csi

# Restart kubelet on node
# ssh to node: systemctl restart kubelet
```

## Network Issues

### 11. DNS Resolution Failures

**Symptoms**: Pods cannot resolve DNS names

**Diagnosis**:
```bash
kubectl-pilot run "test DNS resolution in cluster"
```

**Remediations**:

```bash
# Check CoreDNS status
kubectl get pods -n kube-system -l k8s-app=kube-dns

# Test DNS resolution
kubectl run test --image=busybox --rm -it -- nslookup kubernetes.default

# Restart CoreDNS
kubectl rollout restart deployment/coredns -n kube-system

# Check DNS configuration
kubectl get configmap coredns -n kube-system -o yaml
```

### 12. Inter-Pod Communication Failure

**Symptoms**: Pods cannot communicate with each other

**Remediations**:

```bash
# Check network policies
kubectl get networkpolicies --all-namespaces

# Allow all traffic (dev only)
kubectl delete networkpolicy <policy>

# Check CNI plugin
kubectl get pods -n kube-system | grep -E "calico|flannel|weave"

# Verify pod network connectivity
kubectl-pilot run "test network connectivity between pods"
```

## Resource Issues

### 13. High Memory Usage

**Symptoms**: Pods using excessive memory

**Diagnosis**:
```bash
kubectl top pods
kubectl-pilot diagnose --all-namespaces
```

**Remediations**:

```bash
# Set memory limits
kubectl set resources deployment <deployment> --limits=memory=512Mi

# Add vertical pod autoscaler
kubectl apply -f vpa.yaml

# Check for memory leaks
kubectl-pilot run "analyze memory usage trends"

# Increase node size
kubectl-pilot run "recommend node scaling strategy"
```

### 14. High CPU Usage

**Symptoms**: Pods consuming too much CPU

**Remediations**:

```bash
# Set CPU limits
kubectl set resources deployment <deployment> --limits=cpu=500m

# Add horizontal pod autoscaler
kubectl autoscale deployment <deployment> --cpu-percent=70 --min=2 --max=10

# Profile application
kubectl-pilot explain "how to profile CPU usage in pod"
```

### 15. Resource Quota Exceeded

**Symptoms**: Cannot create new pods due to quota

**Remediations**:

```bash
# Check current usage
kubectl describe resourcequota -n <namespace>

# Increase quota
kubectl patch resourcequota <quota> -p '{"spec":{"hard":{"pods":"100"}}}'

# Delete unused resources
kubectl-pilot run "find unused resources in namespace"
```

## Probe Failures

### 16. Liveness Probe Failures

**Symptoms**: Pod restarting due to liveness probe

**Remediations**:

```bash
# Increase probe timeout
kubectl patch deployment <deployment> --type=json \
  -p='[{"op": "replace", "path": "/spec/template/spec/containers/0/livenessProbe/timeoutSeconds", "value": 5}]'

# Increase initial delay
kubectl patch deployment <deployment> --type=json \
  -p='[{"op": "replace", "path": "/spec/template/spec/containers/0/livenessProbe/initialDelaySeconds", "value": 60}]'

# Check probe endpoint
kubectl exec <pod> -- curl localhost:8080/health
```

### 17. Readiness Probe Failures

**Symptoms**: Pod not receiving traffic

**Remediations**:

```bash
# Check readiness endpoint
kubectl exec <pod> -- curl localhost:8080/ready

# Adjust probe settings
kubectl patch deployment <deployment> --type=json \
  -p='[{"op": "replace", "path": "/spec/template/spec/containers/0/readinessProbe/failureThreshold", "value": 5}]'

# Temporarily disable probe (dev only)
kubectl patch deployment <deployment> --type=json \
  -p='[{"op": "remove", "path": "/spec/template/spec/containers/0/readinessProbe"}]'
```

## Security Issues

### 18. RBAC Permission Denied

**Symptoms**: User/SA cannot perform action

**Remediations**:

```bash
# Check current permissions
kubectl auth can-i <verb> <resource> --as=<user>

# Create role/rolebinding
kubectl create role pod-reader --verb=get,list --resource=pods
kubectl create rolebinding user-pod-reader --role=pod-reader --user=<user>

# Use ClusterRole for cluster-wide access
kubectl create clusterrolebinding admin-binding --clusterrole=admin --user=<user>
```

### 19. Image Security Issues

**Symptoms**: Image fails security scanning

**Remediations**:

```bash
# Scan image
trivy image <image-name>

# Update base image
docker build --build-arg BASE_IMAGE=<new-base> .

# Use non-root user
# Add to Dockerfile: USER 1000
```

### 20. Secret Management

**Symptoms**: Secrets exposed or not accessible

**Remediations**:

```bash
# Create secret from file
kubectl create secret generic <secret> --from-file=key=./secret.txt

# Use external secrets operator
kubectl apply -f external-secrets-operator.yaml

# Encrypt secrets at rest
# Enable encryption in kube-apiserver configuration

# Rotate secrets
kubectl-pilot run "rotate secrets in namespace"
```

## Advanced Troubleshooting

### Using kubectl-pilot for Complex Issues

```bash
# Get AI-powered insights
kubectl-pilot diagnose --all-namespaces

# Natural language troubleshooting
kubectl-pilot run "find pods that are failing and suggest fixes"

# Explain complex situations
kubectl-pilot explain "why is my StatefulSet not scaling"

# Get remediation recommendations
kubectl-pilot diagnose pod <pod> | grep "Recommended Fixes"
```

## Prevention Best Practices

1. **Always set resource requests and limits**
2. **Implement health checks (liveness and readiness probes)**
3. **Use Pod Disruption Budgets for high availability**
4. **Enable monitoring and alerting**
5. **Regular security scanning of images**
6. **Implement network policies**
7. **Use namespaces for isolation**
8. **Regular cluster maintenance and upgrades**
9. **Backup critical data and configurations**
10. **Document custom configurations**

## Getting Help

```bash
# Analyze any issue with AI
kubectl-pilot diagnose <resource-type> <resource-name>

# Ask questions in natural language
kubectl-pilot explain "your question here"

# Get command suggestions
kubectl-pilot run "describe what you want to do"
```

---

**Note**: Always test remediations in a non-production environment first. Use `--dry-run` flags when available.
