/*


Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package controllers

import (
	"bytes"
	"context"
	"fmt"
	"github.com/BurntSushi/toml"
	"github.com/go-logr/logr"
	wso2v1beta1 "github.com/wso2/k8s-wso2is-operator/api/v1beta1"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/api/extensions/v1beta1"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/apimachinery/pkg/util/intstr"
	"log"
	"reflect"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

// Wso2IsReconciler reconciles a Wso2Is object
type Wso2IsReconciler struct {
	client.Client
	Log    logr.Logger
	Scheme *runtime.Scheme
}

// +kubebuilder:rbac:groups=wso2.wso2.com,resources=wso2is,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=wso2.wso2.com,resources=wso2is/status,verbs=get;update;patch

func (r *Wso2IsReconciler) Reconcile(req ctrl.Request) (ctrl.Result, error) {

	// Get context
	ctx := context.Background()

	// Get logger
	log := r.Log.WithValues(deployment_name, req.NamespacedName)

	// Fetch the WSO2IS instance
	instance := wso2v1beta1.Wso2Is{}

	// Check if WSO2 custom resource is present
	err := r.Get(ctx, req.NamespacedName, &instance)
	if err != nil {
		if errors.IsNotFound(err) {
			// Request object not found, could have been deleted after reconcile request.
			// Owned objects are automatically garbage collected. For additional cleanup logic use finalizers.
			// Return and don't requeue
			log.Info("WSO2IS resource not found. Ignoring since object must be deleted")
			return ctrl.Result{}, nil
		}
		// Error reading the object - requeue the request.
		log.Error(err, "Failed to get WSO2IS Instance")
		return ctrl.Result{}, err
	}

	// Add new service account if not present
	saFound := &corev1.ServiceAccount{}
	err = r.Get(ctx, types.NamespacedName{Name: svc_account_name, Namespace: instance.Namespace}, saFound)
	if err != nil && errors.IsNotFound(err) {
		// Define a new deployment
		svc := r.addServiceAccount(instance)
		log.Info("Creating a new ServiceAccount", "ServiceAccount.Namespace", svc.Namespace, "ServiceAccount.Name", svc.Name)
		err = r.Create(ctx, svc)
		if err != nil {
			log.Error(err, "Failed to create new ServiceAccount", "ServiceAccount.Namespace", svc.Namespace, "ServiceAccount.Name", svc.Name)
			return ctrl.Result{}, err
		} else {
			log.Info("Successfully created new ServiceAccount", "ServiceAccount.Namespace", svc.Namespace, "ServiceAccount.Name", svc.Name)
		}
		// ServiceAccount created successfully - return and requeue
		return ctrl.Result{Requeue: true}, nil
	} else if err != nil {
		log.Error(err, "Failed to get ServiceAccount")
		return ctrl.Result{}, err
	}

	// Add new persistant volume claim
	pvcFound := &corev1.PersistentVolumeClaim{}
	err = r.Get(ctx, types.NamespacedName{Name: us_pv_claim_name, Namespace: instance.Namespace}, pvcFound)
	if err != nil && errors.IsNotFound(err) {
		// Define a new deployment
		pvc := r.addPersistentVolumeClaim(instance)
		log.Info("Creating a new PersistentVolumeClaim", "PersistentVolumeClaim.Namespace", pvc.Namespace, "PersistentVolumeClaim.Name", pvc.Name)
		err = r.Create(ctx, pvc)
		if err != nil {
			log.Error(err, "Failed to create new PersistentVolumeClaim", "PersistentVolumeClaim.Namespace", pvc.Namespace, "PersistentVolumeClaim.Name", pvc.Name)
			return ctrl.Result{}, err
		} else {
			log.Info("Successfully created new PersistentVolumeClaim", "PersistentVolumeClaim.Namespace", pvc.Namespace, "PersistentVolumeClaim.Name", pvc.Name)
		}
		// PersistentVolumeClaim created successfully - return and requeue
		return ctrl.Result{Requeue: true}, nil
	} else if err != nil {
		log.Error(err, "Failed to get ServiceAccount")
		return ctrl.Result{}, err
	}

	// Add new config map if not present
	confMap := &corev1.ConfigMap{}
	err = r.Get(ctx, types.NamespacedName{Name: config_map_name, Namespace: instance.Namespace}, confMap)
	if err != nil && errors.IsNotFound(err) {
		// Define a new deployment
		cfgMap := r.addConfigMap(instance, log)
		log.Info("Creating a new ConfigMap", "ConfigMap.Namespace", cfgMap.Namespace, "ConfigMap.Name", cfgMap.Name)
		err = r.Create(ctx, cfgMap)
		if err != nil {
			log.Error(err, "Failed to create new ConfigMap", "ConfigMap.Namespace", cfgMap.Namespace, "ConfigMap.Name", cfgMap.Name)
			return ctrl.Result{}, err
		} else {
			log.Info("Successfully created new ConfigMap", "ConfigMap.Namespace", cfgMap.Namespace, "ConfigMap.Name", cfgMap.Name)
		}
		// ServiceAccount created successfully - return and requeue
		return ctrl.Result{Requeue: true}, nil
	} else if err != nil {
		log.Error(err, "Failed to get ConfigMap")
		return ctrl.Result{}, err
	}

	// Add new service if not present
	serviceFound := &corev1.Service{}
	err = r.Get(ctx, types.NamespacedName{Name: svc_name, Namespace: instance.Namespace}, serviceFound)
	if err != nil && errors.IsNotFound(err) {
		// Define a new deployment
		svc := r.addNewService(instance)
		log.Info("Creating a new Service", "Service.Namespace", svc.Namespace, "Service.Name", svc.Name)
		err = r.Create(ctx, svc)
		if err != nil {
			log.Error(err, "Failed to create new Service", "Service.Namespace", svc.Namespace, "Service.Name", svc.Name)
			return ctrl.Result{}, err
		} else {
			log.Info("Successfully created new Service", "Service.Namespace", svc.Namespace, "Service.Name", svc.Name)
		}
		// ServiceAccount created successfully - return and requeue
		return ctrl.Result{Requeue: true}, nil
	} else if err != nil {
		log.Error(err, "Failed to get Service")
		return ctrl.Result{}, err
	}

	// Update service name in status
	instance.Status.ServiceName = serviceFound.Name

	// Add Ingress if not present
	ingressFound := v1beta1.Ingress{}
	err = r.Get(ctx, types.NamespacedName{Name: ing_name, Namespace: instance.Namespace}, &ingressFound)
	if err != nil && errors.IsNotFound(err) {
		// Define a new Ingress
		svc := r.addNewIngress(instance)
		log.Info("Creating new Ingress", "Ingress.Namespace", svc.Namespace, "Ingress.Name", svc.Name)
		err = r.Create(ctx, svc)
		if err != nil {
			log.Error(err, "Failed to create new Ingress", "Ingress.Namespace", svc.Namespace, "Ingress.Name", svc.Name)
			return ctrl.Result{}, err
		} else {
			log.Info("Successfully created new Ingress", "Ingress.Namespace", svc.Namespace, "Ingress.Name", svc.Name)
		}
		// Ingress created successfully - return and requeue
		return ctrl.Result{Requeue: true}, nil
	} else if err != nil {
		log.Error(err, "Failed to get Ingress")
		return ctrl.Result{}, err
	}

	// Update ingress details in status
	instance.Status.IngressName = ingressFound.Name
	if len(ingressFound.Status.LoadBalancer.Ingress) > 0 {
		instance.Status.IngressHostname = ingressFound.Status.LoadBalancer.Ingress[0].Hostname
		instance.Status.IngressIP = ingressFound.Status.LoadBalancer.Ingress[0].IP
	}

	// Check if the deployment already exists, if not create a new one
	found := &appsv1.Deployment{}
	err = r.Get(ctx, types.NamespacedName{Name: instance.Name, Namespace: instance.Namespace}, found)
	if err != nil && errors.IsNotFound(err) {
		// Define a new deployment
		dep := r.deploymentForWso2Is(instance)
		log.Info("Creating a new Deployment", "Deployment.Namespace", dep.Namespace, "Deployment.Name", dep.Name)
		err = r.Create(ctx, dep)
		if err != nil {
			log.Error(err, "Failed to create new Deployment", "Deployment.Namespace", dep.Namespace, "Deployment.Name", dep.Name)
			return ctrl.Result{}, err
		} else {
			log.Info("Successfully added new Deployment", "Deployment.Namespace", dep.Namespace, "Deployment.Name", dep.Name)
		}
		// Deployment created successfully - return and requeue
		return ctrl.Result{Requeue: true}, nil
	} else if err != nil {
		log.Error(err, "Failed to get Deployment")
		return ctrl.Result{}, err
	}

	// Ensure the deployment size is the same as the spec
	size := instance.Spec.Size
	foundReplicas := found.Spec.Replicas

	// Update replica status
	instance.Status.Replicas = fmt.Sprint(*foundReplicas)

	if *found.Spec.Replicas != size {
		found.Spec.Replicas = &size
		err = r.Update(ctx, found)
		if err != nil {
			log.Error(err, "Failed to update Deployment", "Deployment.Namespace", found.Namespace, "Deployment.Name", found.Name)
			return ctrl.Result{}, err
		}
		// Spec updated - return and requeue
		return ctrl.Result{Requeue: true}, nil
	}

	// Update the IS status with the pod names
	// List the pods for this IS's deployment
	podList := &corev1.PodList{}
	listOpts := []client.ListOption{
		client.InNamespace(instance.Namespace),
		client.MatchingLabels(labelsForWso2IS(instance.Name, instance.Spec.Version)),
	}
	if err = r.List(ctx, podList, listOpts...); err != nil {
		log.Error(err, "Failed to list pods", "WSO2IS.Namespace", instance.Namespace, "WSO2IS.Name", instance.Name)
		return ctrl.Result{}, err
	}
	podNames := getPodNames(podList.Items)

	// Update status.Nodes if needed
	if !reflect.DeepEqual(podNames, instance.Status.Nodes) {
		instance.Status.Nodes = podNames
		err := r.Status().Update(ctx, &instance)
		if err != nil {
			log.Error(err, "Failed to update WSO2IS status")
			return ctrl.Result{}, err
		}
	}

	return ctrl.Result{}, nil
}

// labelsForWso2IS returns the labels for selecting the resources
// belonging to the given WSO2IS CR name.
func labelsForWso2IS(depname string, version string) map[string]string {
	return map[string]string{
		"deployment": depname,
		"app":        depname,
		"monitoring": "jmx",
		"pod":        depname,
		"version":    version,
	}
}

// getPodNames returns the pod names of the array of pods passed in
func getPodNames(pods []corev1.Pod) []string {
	var podNames []string
	for _, pod := range pods {
		podNames = append(podNames, pod.Name)
	}
	return podNames
}

// addServiceAccount adds a new ServiceAccount
func (r *Wso2IsReconciler) addServiceAccount(m wso2v1beta1.Wso2Is) *corev1.ServiceAccount {
	svc := &corev1.ServiceAccount{
		ObjectMeta: metav1.ObjectMeta{
			Name:      svc_account_name,
			Namespace: m.Namespace,
		},
	}
	ctrl.SetControllerReference(&m, svc, r.Scheme)
	return svc
}

// add adds a new PersistentVolume
func (r *Wso2IsReconciler) addPersistentVolumeClaim(m wso2v1beta1.Wso2Is) *corev1.PersistentVolumeClaim {
	pvc := &corev1.PersistentVolumeClaim{
		ObjectMeta: metav1.ObjectMeta{
			Name:      us_pv_claim_name,
			Namespace: m.Namespace,
		},
		Spec: corev1.PersistentVolumeClaimSpec{
			AccessModes: []corev1.PersistentVolumeAccessMode{
				"ReadWriteMany",
			},
			Resources: corev1.ResourceRequirements{
				Requests: corev1.ResourceList{
					corev1.ResourceStorage: resource.MustParse("1Gi"),
				},
			},
		},
	}
	ctrl.SetControllerReference(&m, pvc, r.Scheme)
	return pvc
}

// addConfigMap adds a new ConfigMap
func (r *Wso2IsReconciler) addConfigMap(m wso2v1beta1.Wso2Is, logger logr.Logger) *corev1.ConfigMap {
	configMap := &corev1.ConfigMap{
		ObjectMeta: metav1.ObjectMeta{
			Name:      config_map_name,
			Namespace: m.Namespace,
		},
		Data: map[string]string{
			config_file_name: getTomlConfig(m.Spec, logger),
		},
	}
	ctrl.SetControllerReference(&m, configMap, r.Scheme)
	return configMap
}

func getTomlConfig(spec wso2v1beta1.Wso2IsSpec, logger logr.Logger) string {
	if len(spec.TomlConfig) == 0 {
		buf := new(bytes.Buffer)
		if err := toml.NewEncoder(buf).Encode(spec.Configurations); err != nil {
			log.Println(err)
		}
		logger.Info(buf.String())
		return buf.String()
	} else {
		return spec.TomlConfig
	}
}

// addNewIngress adds a new Ingress Controller
func (r *Wso2IsReconciler) addNewIngress(m wso2v1beta1.Wso2Is) *v1beta1.Ingress {
	ingress := &v1beta1.Ingress{
		ObjectMeta: metav1.ObjectMeta{
			Name:      ing_name,
			Namespace: m.Namespace,
			Annotations: map[string]string{
				"kubernetes.io/ingress.class":                     "nginx",
				"nginx.ingress.kubernetes.io/backend-protocol":    "HTTPS",
				"nginx.ingress.kubernetes.io/affinity":            "cookie",
				"nginx.ingress.kubernetes.io/session-cookie-name": "route",
				"nginx.ingress.kubernetes.io/session-cookie-hash": "sha1",
			},
		},
		Spec: v1beta1.IngressSpec{
			TLS: []v1beta1.IngressTLS{
				{
					Hosts: []string{m.Spec.Configurations.Host},
				},
			},
			Rules: []v1beta1.IngressRule{
				{
					Host: m.Spec.Configurations.Host,
					IngressRuleValue: v1beta1.IngressRuleValue{
						HTTP: &v1beta1.HTTPIngressRuleValue{
							Paths: []v1beta1.HTTPIngressPath{{
								Path: "/",
								Backend: v1beta1.IngressBackend{
									ServiceName: svc_name,
									ServicePort: intstr.IntOrString{
										IntVal: servicePort,
									},
								},
							}},
						},
					},
				},
			},
		},
	}
	ctrl.SetControllerReference(&m, ingress, r.Scheme)
	return ingress
}

// addNewService adds a new Service
func (r *Wso2IsReconciler) addNewService(m wso2v1beta1.Wso2Is) *corev1.Service {
	svc := &corev1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name:      svc_name,
			Namespace: m.Namespace,
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
			Type:     corev1.ServiceTypeLoadBalancer,
		},
	}
	ctrl.SetControllerReference(&m, svc, r.Scheme)
	return svc
}

// New deployment for WSO2IS
func (r *Wso2IsReconciler) deploymentForWso2Is(m wso2v1beta1.Wso2Is) *appsv1.Deployment {
	ls := labelsForWso2IS(m.Name, m.Spec.Version)
	replicas := m.Spec.Size
	runasuser := int64(802)

	dep := &appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name:      m.Name,
			Namespace: m.Namespace,
		},
		Spec: appsv1.DeploymentSpec{
			Replicas: &replicas,
			Selector: &metav1.LabelSelector{
				MatchLabels: ls,
			},
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: ls,
				},
				Spec: corev1.PodSpec{
					Volumes: []corev1.Volume{
						{
							Name: pvc_name,
							VolumeSource: corev1.VolumeSource{
								PersistentVolumeClaim: &corev1.PersistentVolumeClaimVolumeSource{
									ClaimName: us_pv_claim_name,
								},
							},
						},
						{
							Name: config_map_name,
							VolumeSource: corev1.VolumeSource{
								ConfigMap: &corev1.ConfigMapVolumeSource{
									LocalObjectReference: corev1.LocalObjectReference{
										Name: config_map_name,
									},
								},
							},
						},
					},
					Containers: []corev1.Container{{
						Name:  deployment_name,
						Image: "sureshmichael/wso2-is-5.11.0:rc1",
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
						}, {
							Name:  "HOST_NAME",
							Value: m.Spec.Configurations.Host,
						}},
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
								Name:      pvc_name,
								MountPath: "/home/wso2carbon/wso2is-5.11.0/repository/deployment/server/userstores",
							},
							{
								Name:        config_map_name,
								MountPath:   "/home/wso2carbon/wso2-config-volume/repository/conf/deployment.toml",
								SubPathExpr: config_file_name,
							},
						},
						LivenessProbe: &corev1.Probe{
							Handler: corev1.Handler{
								Exec: &corev1.ExecAction{
									Command: []string{"/bin/sh", "-c", "nc -z localhost " + string(containerPortHttps)},
								},
							},
							InitialDelaySeconds: 250,
							PeriodSeconds:       10,
						},
						ReadinessProbe: &corev1.Probe{
							Handler: corev1.Handler{
								Exec: &corev1.ExecAction{
									Command: []string{"/bin/sh", "-c", "nc -z localhost " + string(containerPortHttps)},
								},
							},
							InitialDelaySeconds: 250,
							PeriodSeconds:       10,
						},
						Lifecycle: &corev1.Lifecycle{
							PreStop: &corev1.Handler{
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
					ServiceAccountName: svc_account_name,
					HostAliases: []corev1.HostAlias{{
						IP:        "127.0.0.1",
						Hostnames: []string{m.Spec.Configurations.Host},
					}},
				},
			},
			Strategy: appsv1.DeploymentStrategy{
				Type: "RollingUpdate",
				RollingUpdate: &appsv1.RollingUpdateDeployment{
					MaxUnavailable: &intstr.IntOrString{
						Type:   intstr.Int,
						IntVal: 0,
					},
					MaxSurge: &intstr.IntOrString{
						Type:   intstr.Int,
						IntVal: 1,
					},
				},
			},
			MinReadySeconds: 30,
		},
	}
	// Set WSO2IS instance as the owner and controller
	ctrl.SetControllerReference(&m, dep, r.Scheme)
	return dep
}

func (r *Wso2IsReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&wso2v1beta1.Wso2Is{}).
		Complete(r)
}
