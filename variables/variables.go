package variables

//var ServiceAccountName = "wso2is-serviceaccount"

// var ServiceName = "wso2is-service"
var RoleName = "endpoints-reader-role" // wso2is-role
var SecretName = "wso2is-secret"

// var StatefulSetName = "wso2is-statefulset"
var RoleBindingName = "endpoints-reader-role-wso2-binding"

//var ConfigMapName = "wso2is-configmap"

// var ConfigMapName = "identity-server-conf"
// var IngressName = "wso2is-ingress"
var PersistenVolumeName = "wso2is-persistent-volume"
var UserstorePVCName string = "identity-server-shared-userstores-volume-claim"

const ContainerPortHttps int32 = 9443
const ContainerPortHttp int32 = 9763
const ServicePortHttp int32 = 9763
const ServicePortHttps int32 = 9443

const ContainerImage string = "rukshanjs/wso2is:v6.1.0"

//const DeploymentName string = "wso2is"

const Wso2IsNamespace string = "wso2-iam-system"
