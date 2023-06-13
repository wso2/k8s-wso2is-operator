package wso2is

import (
	"context"
	"fmt"

	"github.com/go-logr/logr"
	wso2v1beta1 "github.com/wso2/k8s-wso2is-operator/api/v1beta1"
	"github.com/wso2/k8s-wso2is-operator/variables"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/apimachinery/pkg/util/intstr"
	ctrl "sigs.k8s.io/controller-runtime"
)

func (r *Wso2IsReconciler) defineStatefulSet(m wso2v1beta1.Wso2Is) *appsv1.StatefulSet {
	ls := labelsForWso2IS(m.Name, m.Spec.Version)
	replicas := m.Spec.Size
	runasuser := int64(802)

	sfs := &appsv1.StatefulSet{
		ObjectMeta: metav1.ObjectMeta{
			Name:      m.Name,
			Namespace: m.Namespace,
		},
		Spec: appsv1.StatefulSetSpec{
			Replicas:    &replicas,
			ServiceName: variables.ServiceName,
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
							Name: variables.PersistenVolumeName,
							VolumeSource: corev1.VolumeSource{
								PersistentVolumeClaim: &corev1.PersistentVolumeClaimVolumeSource{
									ClaimName: variables.UserstorePVCName,
								},
							},
						},
						{
							Name: variables.ConfigMapName,
							VolumeSource: corev1.VolumeSource{
								ConfigMap: &corev1.ConfigMapVolumeSource{
									LocalObjectReference: corev1.LocalObjectReference{
										Name: variables.ConfigMapName,
									},
								},
							},
						},
						{
							Name: variables.SecretName,
							VolumeSource: corev1.VolumeSource{
								Secret: &corev1.SecretVolumeSource{
									SecretName: variables.SecretName,
								},
							},
						},
					},
					Containers: []corev1.Container{{
						Name:  variables.DeploymentName,
						Image: variables.ContainerImage,
						Ports: []corev1.ContainerPort{{
							ContainerPort: variables.ContainerPortHttps,
							Protocol:      "TCP",
						}, {
							ContainerPort: variables.ContainerPortHttp,
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
								Name:      variables.PersistenVolumeName,
								MountPath: fmt.Sprintf("/home/wso2carbon/wso2is-%s/repository/deployment/server/userstores", m.Spec.Version),
							},
							{
								Name:        variables.ConfigMapName,
								MountPath:   "/home/wso2carbon/wso2-config-volume/repository/conf/deployment.toml",
								SubPathExpr: "deployment.toml",
							},
							{
								Name:      variables.SecretName,
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
							PeriodSeconds: 10,
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
					ServiceAccountName: variables.ServiceAccountName,
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
	ctrl.SetControllerReference(&m, sfs, r.Scheme)
	return sfs
}

func reconcileStatefulSet(r *Wso2IsReconciler, instance wso2v1beta1.Wso2Is, log logr.Logger, err error, ctx context.Context) (ctrl.Result, error) {
	sfsDefinition := r.defineStatefulSet(instance)
	sfs := &appsv1.StatefulSet{}
	err = r.Get(ctx, types.NamespacedName{Name: instance.Name, Namespace: instance.Namespace}, sfs)
	if err != nil {
		if errors.IsNotFound(err) {
			log.Info("StatefulSet resource " + instance.Name + " not found. Creating or re-creating statefulset")
			err = r.Create(ctx, sfsDefinition)
			if err != nil {
				log.Error(err, "Failed to create new StatefulSet", "StatefulSet.Namespace", sfsDefinition.Namespace, "StatefulSet.Name", sfsDefinition.Name)
				return ctrl.Result{}, err
			}
		} else {
			log.Info("Failed to get statefulset resource " + instance.Name + ". Re-running reconcile.")
			return ctrl.Result{}, err
		}
	} else {
		log.Info("Found StatefulSet")
		// Note: For simplication purposes StatefulSets are not updated - see deployment section
		// wso2Is := &wso2v1beta1.Wso2Is{}
		// err = r.Get(ctx, types.NamespacedName{Name: instance.Name, Namespace: instance.Namespace}, wso2Is)
		// if err != nil {
		// 	log.Error(err, "Failed to get Wso2Is resource")
		// 	return ctrl.Result{}, err
		// }

		//if sfs.Spec.Template.ObjectMeta.Annotations["configmapHash"] != instance.Spec.Template.Annotations["configmapHash"] {
		//	sfs.Spec.Template.ObjectMeta.Annotations["configmapHash"] = instance.Spec.Template.Annotations["configmapHash"]
		//	err = r.Update(ctx, sfs)
		//	if err != nil {
		//		log.Error(err, "Failed to update Wso2Is resource")
		//		return ctrl.Result{}, err
		//	}
		//}

		// Update replica status
		instance.Status.Replicas = fmt.Sprint(*sfs.Spec.Replicas)

		// Ensure the StatefulSet size is the same as the spec
		size := instance.Spec.Size

		if *sfs.Spec.Replicas != size {
			sfs.Spec.Replicas = &size
			err = r.Update(ctx, sfs)
			if err != nil {
				log.Error(err, "Failed to update StatefulSet", "StatefulSet.Namespace", sfs.Namespace, "StatefulSet.Name", sfs.Name)
				return ctrl.Result{}, err
			}
			// Spec updated - return and requeue
			return ctrl.Result{Requeue: true}, nil
		}
	}
	return ctrl.Result{}, nil
}
