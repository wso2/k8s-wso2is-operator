package controllers

const containerImage string = "rukshanjs/wso2is:v6.1.0"

const containerPortHttps int32 = 9443

const containerPortHttp int32 = 9763

const servicePortHttp int32 = 9763

const servicePortHttps int32 = 9443

const svcAccountName string = "wso2is-svcaccount"

const svcName string = "wso2is-service"

const configMapName string = "identity-server-conf"

const ingName string = "wso2is-ingress"

const deploymentName string = "wso2is"

const configFileName string = "deployment.toml"

const secretName string = "wso2is-secret"

const roleName string = "endpoints-reader-role"
const roleBindingName string = "endpoints-reader-role-wso2-binding"

const usPvClaimName string = "identity-server-shared-userstores-volume-claim"

const pvName string = "user-store-persistent-storage"
