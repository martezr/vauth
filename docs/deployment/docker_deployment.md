# Docker Deployment

The vAuth platform can be deployed with Docker using the following command:

```
docker run --name vauth \
-e UI_PORT=9000 -e DATA_DIR=/app \
-e VSPHERE_SERVER=grtvcenter01.grt.local \
-e VSPHERE_USERNAME=vauth@vsphere.local \
-e VSPHERE_PASSWORD="Password123#" \
-e VAULT_ADDRESS="https://grtmanage01.grt.local:8200" \ -e VAULT_TOKEN="s.ewdkchV1oqIwTxxI8G3INWVG" \
-e VAULT_APPROLE_MOUNT=approle public.ecr.aws/i4r5n0t9/vauth:1.0
```