package controllers

import (
	"fmt"
	wso2v1beta1 "github.com/wso2/k8s-wso2is-operator/api/v1beta1"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
	ctrl "sigs.k8s.io/controller-runtime"
)

func (r *Wso2IsReconciler) statefulSetForWso2Is(m wso2v1beta1.Wso2Is) *appsv1.StatefulSet {
	ls := labelsForWso2IS(m.Name, m.Spec.Version)
	replicas := m.Spec.Size
	runasuser := int64(802)

	statefulSet := &appsv1.StatefulSet{
		ObjectMeta: metav1.ObjectMeta{
			Name:      m.Name,
			Namespace: m.Namespace,
		},
		Spec: appsv1.StatefulSetSpec{
			Replicas:    &replicas,
			ServiceName: svcName,
			Selector: &metav1.LabelSelector{
				MatchLabels: ls,
			},
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels:      ls,
					Annotations: m.Spec.Template.Annotations,
				},
				Spec: corev1.PodSpec{
					Volumes: []corev1.Volume{
						{
							Name: pvName,
							VolumeSource: corev1.VolumeSource{
								PersistentVolumeClaim: &corev1.PersistentVolumeClaimVolumeSource{
									ClaimName: usPvClaimName,
								},
							},
						},
						{
							Name: configMapName,
							VolumeSource: corev1.VolumeSource{
								ConfigMap: &corev1.ConfigMapVolumeSource{
									LocalObjectReference: corev1.LocalObjectReference{
										Name: configMapName,
									},
								},
							},
						},
						{
							Name: secretName,
							VolumeSource: corev1.VolumeSource{
								Secret: &corev1.SecretVolumeSource{
									SecretName: secretName,
								},
							},
						},
					},
					Containers: []corev1.Container{{
						Name:  deploymentName,
						Image: containerImage,
						Ports: []corev1.ContainerPort{{
							ContainerPort: containerPortHttps,
							Protocol:      "TCP",
						}, {
							ContainerPort: containerPortHttp,
							Protocol:      "TCP",
						}},
						Env: []corev1.EnvVar{{
							Name: "NODE_IP",
							ValueFrom: &corev1.EnvVarSource{
								FieldRef: &corev1.ObjectFieldSelector{
									FieldPath: "status.podIP",
								},
							},
						},
						//{
						//	Name:  "HOST_NAME",
						//	Value: m.Spec.Configurations.Host,
						//}
						},
						/* @TODO Please uncomment for live production
						Resources: corev1.ResourceRequirements{
							Requests: corev1.ResourceList{
								corev1.ResourceCPU:    resource.MustParse("1Gi"),
								corev1.ResourceMemory: resource.MustParse("1000m"),
							},
							Limits: corev1.ResourceList{
								corev1.ResourceCPU:    resource.MustParse("2Gi"),
								corev1.ResourceMemory: resource.MustParse("2000m"),
							},
						},
						*/
						VolumeMounts: []corev1.VolumeMount{
							{
								Name:      pvName,
								MountPath: fmt.Sprintf("/home/wso2carbon/wso2is-%s/repository/deployment/server/userstores", m.Spec.Version),
							},
							{
								Name:        configMapName,
								MountPath:   "/home/wso2carbon/wso2-config-volume/repository/conf/deployment.toml",
								SubPathExpr: configFileName,
							},
							{
								Name:      secretName,
								MountPath: fmt.Sprintf("/home/wso2carbon/wso2is-%s/repository/resources/security/controller-keystores", m.Spec.Version),
								ReadOnly:  true,
							},
						},
						StartupProbe: &corev1.Probe{
							ProbeHandler: corev1.ProbeHandler{
								Exec: &corev1.ExecAction{
									Command: []string{
										"/bin/sh",
										"-c",
										"nc -z localhost 9443",
									},
								},
							},
							InitialDelaySeconds: 60,
							PeriodSeconds:       5,
							FailureThreshold:    30,
						},
						LivenessProbe: &corev1.Probe{
							ProbeHandler: corev1.ProbeHandler{
								HTTPGet: &corev1.HTTPGetAction{
									Path:   "/carbon/admin/login.jsp",
									Port:   intstr.FromInt(9443),
									Scheme: corev1.URISchemeHTTPS,
								},
							},
							PeriodSeconds: 1,
						},
						ReadinessProbe: &corev1.Probe{
							ProbeHandler: corev1.ProbeHandler{
								HTTPGet: &corev1.HTTPGetAction{
									Path:   "/api/health-check/v1.0/health",
									Port:   intstr.FromInt(9443),
									Scheme: corev1.URISchemeHTTPS,
								},
							},
							InitialDelaySeconds: 60,
							PeriodSeconds:       10,
						},
						Lifecycle: &corev1.Lifecycle{
							PreStop: &corev1.LifecycleHandler{
								Exec: &corev1.ExecAction{
									Command: []string{
										"sh",
										"-c",
										"${WSO2_SERVER_HOME}/bin/wso2server.sh stop",
									},
								},
							},
						},
						ImagePullPolicy: "IfNotPresent",
						SecurityContext: &corev1.SecurityContext{
							RunAsUser: &runasuser,
						},
					}},
					ServiceAccountName: svcAccountName,
					//HostAliases: []corev1.HostAlias{{
					//	IP:        "127.0.0.1",
					//	Hostnames: []string{m.Spec.Configurations.Host},
					//}},
				},
			},
			MinReadySeconds: 30,
		},
	}
	// Set WSO2IS instance as the owner and controller
	ctrl.SetControllerReference(&m, statefulSet, r.Scheme)
	return statefulSet
}
