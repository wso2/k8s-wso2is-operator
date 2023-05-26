package controllers

import (
	"github.com/go-logr/logr"
	wso2v1beta1 "github.com/wso2/k8s-wso2is-operator/api/v1beta1"
	corev1 "k8s.io/api/core/v1"
	rbacv1 "k8s.io/api/rbac/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
	ctrl "sigs.k8s.io/controller-runtime"
)

// addServiceAccount adds a new ServiceAccount
func (r *Wso2IsReconciler) addServiceAccount(m wso2v1beta1.Wso2Is) *corev1.ServiceAccount {
	svc := &corev1.ServiceAccount{
		ObjectMeta: metav1.ObjectMeta{
			Name:      svcAccountName,
			Namespace: m.Namespace,
		},
	}
	ctrl.SetControllerReference(&m, svc, r.Scheme)
	return svc
}

// addNewService adds a new Service
func (r *Wso2IsReconciler) addNewService(m wso2v1beta1.Wso2Is) *corev1.Service {
	ls := labelsForWso2IS(m.Name, m.Spec.Version)

	// Make Service type configurable
	//serviceType := corev1.ServiceTypeNodePort
	serviceType := corev1.ServiceTypeClusterIP
	if m.Spec.Configurations.ServiceType == "NodePort" {
		serviceType = corev1.ServiceTypeNodePort
	} else if m.Spec.Configurations.ServiceType == "ClusterIP" {
		serviceType = corev1.ServiceTypeClusterIP
	} else if m.Spec.Configurations.ServiceType == "LoadBalancer" {
		serviceType = corev1.ServiceTypeLoadBalancer
	}

	svc := &corev1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name:      svcName,
			Namespace: m.Namespace,
			Labels:    ls,
		},
		Spec: corev1.ServiceSpec{
			Ports: []corev1.ServicePort{{
				Name:     "servlet-http",
				Protocol: "TCP",
				Port:     servicePortHttp,
				TargetPort: intstr.IntOrString{
					IntVal: servicePortHttp,
				},
			}, {
				Name:     "servlet-https",
				Protocol: "TCP",
				Port:     servicePortHttps,
				TargetPort: intstr.IntOrString{
					IntVal: servicePortHttps,
				},
			}},
			Selector: labelsForWso2IS(m.Name, m.Spec.Version),
			Type:     serviceType,
		},
	}
	ctrl.SetControllerReference(&m, svc, r.Scheme)
	return svc
}

// addConfigMap adds a new ConfigMap
func (r *Wso2IsReconciler) addConfigMap(m wso2v1beta1.Wso2Is, logger logr.Logger) *corev1.ConfigMap {
	configMap := &corev1.ConfigMap{
		ObjectMeta: metav1.ObjectMeta{
			Name:      configMapName,
			Namespace: m.Namespace,
		},
		Data: map[string]string{
			configFileName: getTomlConfig(m.Spec, logger),
		},
	}
	ctrl.SetControllerReference(&m, configMap, r.Scheme)
	return configMap
}

func (r *Wso2IsReconciler) addSecret(m wso2v1beta1.Wso2Is, logger logr.Logger) *corev1.Secret {
	//append all secrets to here
	secretsMap := map[string][]byte{}
	//mount keystore secrets to Kubernetes secrets
	for _, element := range m.Spec.KeystoreMounts {
		secretsMap[element.Name] = []byte(element.Data)
	}
	secret := &corev1.Secret{
		TypeMeta: metav1.TypeMeta{},
		ObjectMeta: metav1.ObjectMeta{
			Name:      secretName,
			Namespace: m.Namespace,
		},
		Data: secretsMap,
	}
	ctrl.SetControllerReference(&m, secret, r.Scheme)
	return secret
}

func (r *Wso2IsReconciler) addRole(m wso2v1beta1.Wso2Is) *rbacv1.Role {
	role := &rbacv1.Role{
		ObjectMeta: metav1.ObjectMeta{
			Name:      roleName,
			Namespace: m.Namespace,
		},
		Rules: []rbacv1.PolicyRule{
			{
				APIGroups: []string{""},
				Verbs:     []string{"get", "list"},
				Resources: []string{"endpoints"},
			},
		},
	}
	ctrl.SetControllerReference(&m, role, r.Scheme)
	return role
}

func (r *Wso2IsReconciler) addRoleBinding(m wso2v1beta1.Wso2Is) *rbacv1.RoleBinding {
	roleBinding := &rbacv1.RoleBinding{
		ObjectMeta: metav1.ObjectMeta{
			Name:      roleBindingName,
			Namespace: m.Namespace,
		},
		Subjects: []rbacv1.Subject{
			{
				Kind:      "ServiceAccount",
				Name:      roleBindingName,
				Namespace: m.Namespace,
			},
		},
		RoleRef: rbacv1.RoleRef{
			Kind:     "Role",
			Name:     roleName, // Name of the Role created by addRole function
			APIGroup: "rbac.authorization.k8s.io",
		},
	}
	ctrl.SetControllerReference(&m, roleBinding, r.Scheme)
	return roleBinding
}
