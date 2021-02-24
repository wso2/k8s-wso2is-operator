
## WSO2 Identity Server - K8S Operator

The following CRD operator can be used to deploy WSO2 IS on your Kubernates Cluster. If you want to deploy the Identity
Server via Helm operator please refer to the given
link:  [https://github.com/wso2/kubernetes-is](https://github.com/wso2/kubernetes-is)

#### Key benefits
- Auto healing
- Ability to make a test clusters
- Ability to provision multiple ISs on same cluster 
- Custom Keystore addition 
- Ability to mount custom deployment TOML files
- Seameless updates


## Prerequisites (Development)[](https://sdk.operatorframework.io/docs/building-operators/golang/tutorial/#prerequisites)

- Access to a Kubernetes v1.11.3+ cluster (v1.16.0+ if using  `apiextensions.k8s.io/v1`  CRDs).
- User logged with admin permission.
  See  [how to grant yourself cluster-admin privileges or be logged in as admin](https://cloud.google.com/kubernetes-engine/docs/how-to/role-based-access-control#iam-rolebinding-bootstrap)
- [Homebrew](https://brew.sh/) installed
- Git command line installed and configured
- [GoLang](https://golang.org/) installed and correctly configured, including system paths

## System Architecture
![enter image description here](https://user-images.githubusercontent.com/3047253/105663226-b9149900-5ef7-11eb-825b-0413649a99ed.jpg)

## External Database Setup

Please follow the instructions given in the documentation to setup the external MySQL databases

- https://is.docs.wso2.com/en/5.11.0/setup/changing-to-mysql/
- https://is.docs.wso2.com/en/5.11.0/setup/changing-datasource-bpsds/
- https://is.docs.wso2.com/en/5.11.0/setup/changing-datasource-consent-management/

### Databases to be created

- WSO2_IDENTITY_DB
- WSO2_SHARED_DB
- WSO2_CONSENT_DB (Optional)
- WSO2_BPS_DB (Optional)

## Development Environment Setup

Please follow the following instructions to install Operator-SDK in your development environment.

    brew install operator-sdk

Clone the repository by running the following command

    git clone https://github.com/wso2/k8s-wso2is-operator.git

Navigate to the project directory

    cd k8s-wso2is-operator

Run the following command to install dependancies

    make install


Apply the CRDs by running the following command

    kubectl apply -f config/crd/bases/iam.wso2.com_wso2is.yaml
    kubectl apply -f config/crd/bases/iam.wso2.com_userstores.yaml

Feel free to change any configurations at **config/samples/wso2_v1_wso2is.yaml**
Once you do the config changes apply the config by running

    kubectl apply -f config/samples/wso2_v1_wso2is.yaml

Finally run the following command to run the operator in your cluster

    make run

## Installation (Stand alone)

It is possible to deploy a stand alone version of the IS Operator in your cluster as well. You many follow the given steps in order to setup correctly.

**Prerequisites** 
- A database configured and it should be accessible by all pods
- A persistence storage has be configured with ReadWriteMany permission
-- For AWS users, you can refer to Elastic File System (EFS) docs and learn about the configuration https://docs.aws.amazon.com/eks/latest/userguide/efs-csi.html
-- Microsoft Azure users can use AzureFile as the persistent storage
-- Google Cloud users may use GCEPersistentDisk
- Also you will need to have a Ingress controller ready, something that matches to your endpoint


Run the given command  within your cluster

    kubectl apply -f https://raw.githubusercontent.com/wso2/k8s-wso2is-operator/main/artifacts/operator.yaml
    
Finally you may apply your own configurations by refering to the formats given in samples
https://github.com/wso2/k8s-wso2is-operator/tree/main/config/samples
