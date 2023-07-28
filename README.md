# WSO2 Identity Server - Kubernetes Operator

## Introduction

This Kubernetes (k8s) operator allows you to create a clustered deployment of the WSO2 Identity Server (WSO2 IS) with very less friction. All you have to do is have a k8s cluster running this operator and then you can do the deployment using just a .yaml file.

The other (more complex) alternative to do such a deployment of the WSO2 IS is using helm. If you want to deploy the WSO2 IS via helm operator, please refer to the
link [https://github.com/wso2/kubernetes-is](https://github.com/wso2/kubernetes-is)

If you want to easily do the deployment, read along!

### Key benefits of using this operator

- Auto healing
- Horizontal Pod Autoscaler (HPA) enabled
- Ability to make a test clusters
- Ability to provision multiple ISs on same cluster
- Custom Keystore addition
- Ability to mount custom deployment TOML files
- Seameless updates
- No need to manually edit a deployment.toml file
- A controlled, clustered deployment of the WSO2 IS can be easily deployed

## Usage

### Pre-requisites

1. [A working k8s cluster](https://kubernetes.io/docs/setup/production-environment/tools/kubeadm/create-cluster-kubeadm/)
2. kubectl CLI
3. Databases configured. (optional, these configurations are recommended outside the k8s cluster) Required databases configured, and it should be accessible by all pods
   1. Please read the documentation at: [https://is.docs.wso2.com/en/latest/setup/working-with-databases/](https://is.docs.wso2.com/en/latest/setup/working-with-databases/)
   2. The following two databases are required for a standard WSO2 IS deployment
      1. WSO2_IDENTITY_DB
      2. WSO2_SHARED_DB
4. A [persistence volume](https://kubernetes.io/docs/concepts/storage/persistent-volumes/) has be configured with ReadWriteMany permission
   1. The easiest way is to use the NFS server using the NFS provisioner
   2. For AWS users, you can refer to Elastic File System (EFS) docs and learn about the configurations: [https://docs.aws.amazon.com/eks/latest/userguide/efs-csi.html](https://docs.aws.amazon.com/eks/latest/userguide/efs-csi.html)
   3. Microsoft Azure users can use [AzureFile](https://docs.microsoft.com/en-us/azure/aks/azure-files-dynamic-pv) as the persistent storage
   4. Google Cloud users may use [GCEPersistentDisk](https://cloud.google.com/kubernetes-engine/docs/concepts/persistent-volumes)
5. Metrics server running (required for HPA)
6. Certificate setup
7. Also you will need to have an [Ingress](https://kubernetes.io/docs/concepts/services-networking/ingress/) ready to route requests from your endpoint to service, your ingress can vary from cloud provider to provider.

```
Note : Untill the operator is released as an artifact, you would have to clone the repo to execute the commands given below.
```

### Example guide on setting up pre-requisites

The following is one way to setup the pre-requisites mentioned above.

#### 1. Setting up an AKS cluster

https://learn.microsoft.com/en-us/azure/aks/learn/quick-kubernetes-deploy-portal?tabs=azure-cli

#### 2. Install kubectl

https://learn.microsoft.com/en-us/cli/azure/aks?view=azure-cli-latest#az-aks-install-cli

#### 3. External Database Setup

If you want to quickly view the operator functioning, you can deploy a MySQL database within the cluster itself, execute the following.

But this MySQL deployment assumes the `wso2-iam-system` namespace by default. Therefore, it's required to create it. If not, you can execute the yaml file to do that and install the CRDs that would otherwise have to be installed in a later step.

```
kubectl apply -f artifacts/k8s-wso2is-operator.yaml
```

```
kubectl apply -k config/samples/mysql/overlay/dev
```

If you use the above command to create a MySQL instance, it'll automatically execute the relevant database scripts required for the WSO2 IS deployment.

Otherwise you have to connect an external MySQL database (Eg. a managed MySQL server) and properly configure it when apply the Wso2Is resource at a later step.

Please follow the instructions given in the documentation to setup the external MySQL databases

- https://is.docs.wso2.com/en/latest/setup/changing-to-mysql/
- https://is.docs.wso2.com/en/latest/setup/changing-datasource-bpsds/
- https://is.docs.wso2.com/en/latest/setup/changing-datasource-consent-management/

##### Databases to be created

- WSO2_IDENTITY_DB
- WSO2_SHARED_DB
- WSO2_CONSENT_DB (Optional)
- WSO2_BPS_DB (Optional)

#### 4. Volume setup

PVCs and PVs are required for the functioning of the WSO2 IS deployment.

1. Install the NFS provisioner

```
helm repo add wso2 https://helm.wso2.com

helm repo update

helm install nfs-server-provisioner wso2/nfs-server-provisioner -n wso2-iam-system
```

2. Then you can create a PVC using the default `nfs` storageClassName

#### 5. Metrics server setup

For the proper working of the HPA, it needs to access the utilization metrics of the pods. For this, the cluster should have a running metrics server.

To install one within your cluster, follow the instructions.

```
kubectl apply -f https://github.com/kubernetes-sigs/metrics-server/releases/latest/download/components.yaml
```

#### 6. Certificate setup

To allow access to the WSO2 IS instance via HTTPS through the ingress (will be created in next step), we need to setup a TLS certificate for the deployment.

The whole process of managing this TLS cert will be automatically handled by the cert-manager that we are going to install in our cluster next.

For this, follow the instructions.

```
helm repo add jetstack https://charts.jetstack.io

helm repo update

helm install \
  cert-manager jetstack/cert-manager \
  --namespace cert-manager \
  --create-namespace \
  --version v1.12.0 \
  --set installCRDs=true
```

#### 7. Ingress setup

The setting up of an ingress can vary from cloud provider to provider. The following commands have been tested with an AKS (Azure) cluster.

1. Install the nginx ingress controller in the cluster

```
NAMESPACE=ingress-basic

helm repo add ingress-nginx https://kubernetes.github.io/ingress-nginx
helm repo update

helm install ingress-nginx ingress-nginx/ingress-nginx \
  --create-namespace \
  --namespace $NAMESPACE \
  --set controller.service.annotations."service\.beta\.kubernetes\.io/azure-load-balancer-health-probe-request-path"=/healthz
```

Refer https://learn.microsoft.com/en-us/azure/aks/ingress-basic?tabs=azure-cli for further info

2. After this you can create your ingress resource within the cluster. If everything is correctly setup, this ingress resource should be assigned a new IP and it should use the certificate generated by the cert-manager's issuer.

### Steps to deploy the WSO2 IS within the cluster

#### 1. Install the operator

```
kubectl apply -f artifacts/k8s-wso2is-operator.yaml
```

Applying this will,

- Create the namespaces required
- Create the `ClusterRole` required for the operator
- Create the `ServiceAccount` required for the operator
- Create the `ClusterRoleBinding` required for connecting the `ServiceAccount` and `ClusterRole`
- Create the `CustomResourceDefinition` required for the WSO2 IS deployment
- Create the `CustomResourceDefinition` required for the Userstores
- Create a `Deployment` using the operator image

#### 2. Create the TLS cert

- Since we have already installed the cert-manager controller above, now what we have to do is apply the k8s `Issuer` resource which will take care of provisioning and managing the TLS cert for our deployment.
- Here it's very important to set a valid email address to the `spec.acme.email` field in `/config/certmanager/certificate.yaml`

#### 3. Apply the Wso2Is .yaml

The following command uses kustomize to apply all the .yaml files with appropriate sample values to simplify the deployment.

Feel free to use a separate .yaml files if required. Otherwise you can customize the values within the `overlay` folder within each sample within the `samples` folder.

```
kubectl apply -k config/samples/04-azure-basic/overlay/dev
```

Applying this will,

- Create the `Ingress` resource
- Create the `PersistentVolumeClaim` resource
- Create the `Wso2Is` custom resource (using the ealier created CRD)

4. Run the scenario tests by following instrutions within `/testbin/` directory
5. Visit the URL (hostname) given within the ingress and you should be able to access your WSO2 IS deployment with TLS enable, over HTTPS.

## System Architecture

### General architecture

![enter image description here](https://user-images.githubusercontent.com/3047253/105663226-b9149900-5ef7-11eb-825b-0413649a99ed.jpg)

### Sample deployment scenario

![](https://raw.githubusercontent.com/RukshanJS/k8s-wso2is-operator/is-6.1.0-clustering/image1.png)

## Development

### Pre-requisites

- A working k8s cluster (Minikube is preferred for development)
- Access to a Kubernetes v1.26.3+ cluster
- kubectl CLI
- [Homebrew](https://brew.sh/) installed
- mkcert
- maven
- OpenJDK 11
- [GoLang](https://golang.org/) installed and correctly configured, including system paths
- Git command line installed and configured
- User logged with admin permission.
  See [how to grant yourself cluster-admin privileges or be logged in as admin](https://cloud.google.com/kubernetes-engine/docs/how-to/role-based-access-control#iam-rolebinding-bootstrap)

### Steps

This is a small guidance on how to develop this operator further.

1. Fork the repo

2. Install `Operator SDK` in your development environment.

```
brew install operator-sdk
```

3. Clone the repository by running the following command

```
git clone https://github.com/wso2/k8s-wso2is-operator.git
```

4. Navigate to the project directory

```
cd k8s-wso2is-operator
```

5. Run the following command to install dependancies

```
make install
```

6. Apply the CRDs by running the following command

```
kubectl apply -f config/crd/bases/iam.wso2.com_wso2is.yaml
kubectl apply -f config/crd/bases/iam.wso2.com_userstores.yaml
```

7. Customize the files and configurations inside the `/config/samples/` folder. There are three types of configuration presets.

- `/config/samples/01-inline-configs`
- `/config/samples/02-ref-configs`
- `/config/samples/03-custom-configs`

Inside each of these folders, there are three presets for `dev`, `stage` and `prod`. You can customize these values to easily have different sample deployments for various environments.

8. Once you do the config changes apply the config by running

```
kubectl apply -f config/samples/wso2_v1_wso2is.yaml
```

9. Finally run the following command to run the operator in your cluster

```
make run
```

10. Put a PR to the original repo

## Authors

- [Rukshan J. Senanayaka](https://github.com/RukshanJS)
- [Suresh Peiris](https://github.com/tsuresh)
- [Deependra Ariyadewa](https://github.com/gnudeep)
- [Maheshika Goonetilleke](https://github.com/maheshika)

## References

- https://sdk.operatorframework.io/docs/building-operators/golang/tutorial/#prerequisites
- [Part 01: Deploying WSO2 Identity Server 5.11.0 on Kubernetes with all new K8s Operator](https://tsmpeiris.medium.com/part-01-deploying-wso2-identity-server-5-11-0-on-kubernetes-with-all-new-k8s-operator-e6d9e76d7e7)
- [Part 02: Deploying WSO2 Identity Server 5.11.0 on Kubernetes with all new K8s Operator](https://medium.com/@tsmpeiris/part-02-deploying-wso2-identity-server-5-11-0-on-kubernetes-with-all-new-k8s-operator-5d751c1f4ba0)
- https://github.com/wso2/kubernetes-is
- https://kubernetes.io/docs/setup/production-environment/tools/kubeadm/create-cluster-kubeadm/

## Disclaimer

This project is still under active development and should not be considered ready for production. PRs and Issues are welcome and highly appreciated.
