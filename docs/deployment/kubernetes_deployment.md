# Kubernetes Deployment

The vAuth platform can be deployed to a Kubernetes cluster using the following manifest:

```yaml
apiVersion: v1
kind: Namespace
metadata:
  name: vauth
---
apiVersion: v1
kind: ConfigMap
metadata:
  name: vauth-config
  namespace: vauth
data:
  VSPHERE_SERVER: "grtvcenter01.grt.local"
  DATA_DIR: "/vauthdata"
  VSPHERE_TLS_SKIP_VERIFY: "true"
  VSPHERE_USERNAME: "vauth@vsphere.local"
  VSPHERE_PASSWORD: "Password123#"
  VSPHERE_DATACENTERS: ["DC1","DC2"]
  VAULT_ADDRESS: "https://10.0.0.202:8200"
  VAULT_TOKEN: "s.r5A9FBMiQyRzXcEh7Ab7ZE4K"
  VAULT_APPROLE_MOUNT: "approle"
  VAULT_WRAP_RESPONSE: "true"
  VAULT_TLS_SKIP_VERIFY: "true"
---
apiVersion: v1
kind: Pod
metadata:
  name: vauth
  namespace: vauth
  labels:
    name: vauth
spec:
  containers:
  - name: vauth
    image: public.ecr.aws/i4r5n0t9/vauth:1.0
    imagePullPolicy: Always
    envFrom:
      - configMapRef:
          name: vauth-config
    volumeMounts:
    - mountPath: /vauthdata
      name: cache-volume
  volumes:
  - name: cache-volume
    emptyDir: {}
```