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
	networkingv1 "k8s.io/api/networking/v1"
	"k8s.io/apimachinery/pkg/api/errors"
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

// +kubebuilder:rbac:groups=iam.wso2.com,resources=wso2is,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=iam.wso2.com,resources=wso2is/status,verbs=get;update;patch

func (r *Wso2IsReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {

	// Get logger
	logger := r.Log.WithValues(deploymentName, req.NamespacedName)

	// Fetch the WSO2IS instance
	instance := wso2v1beta1.Wso2Is{}

	// Check if WSO2 custom resource is present
	err := r.Get(ctx, req.NamespacedName, &instance)
	if err != nil {
		if errors.IsNotFound(err) {
			// Request object not depFound, could have been deleted after reconcile request.
			// Owned objects are automatically garbage collected. For additional cleanup logic use finalizers.
			// Return and don't requeue
			logger.Info("WSO2IS resource not depFound. Ignoring since object must be deleted")
			return ctrl.Result{}, nil
		}
		// Error reading the object - requeue the request.
		logger.Error(err, "Failed to get WSO2IS Instance")
		return ctrl.Result{}, err
	}

	// Add new service account if not present
	saFound := &corev1.ServiceAccount{}
	err = r.Get(ctx, types.NamespacedName{Name: svcAccountName, Namespace: instance.Namespace}, saFound)
	if err != nil && errors.IsNotFound(err) {
		return reconcileSva(r, instance, logger, err, ctx)
	} else if err != nil {
		logger.Error(err, "Failed to get ServiceAccount")
		return ctrl.Result{}, err
	}

	// Check for persistent volume claim
	pvcFound := &corev1.PersistentVolumeClaim{}
	err = r.Get(ctx, types.NamespacedName{Name: usPvClaimName, Namespace: instance.Namespace}, pvcFound)
	if err != nil && errors.IsNotFound(err) {
		logger.Info("Unable to detect PVC claim in your cluster. You may configure your own")
		return ctrl.Result{}, nil
	} else if err != nil {
		logger.Error(err, "Failed to get PersistentVolumeClaim")
		return ctrl.Result{}, err
	}

	// Add new config map if not present
	confMap := &corev1.ConfigMap{}
	err = r.Get(ctx, types.NamespacedName{Name: configMapName, Namespace: instance.Namespace}, confMap)
	if err != nil && errors.IsNotFound(err) {
		return reconcileCfg(r, instance, logger, err, ctx)
	} else if err != nil {
		logger.Error(err, "Failed to get ConfigMap")
		return ctrl.Result{}, err
	}

	// Add new secret if not present
	secretFound := &corev1.Secret{}
	err = r.Get(ctx, types.NamespacedName{Name: secretName, Namespace: instance.Namespace}, secretFound)
	if err != nil && errors.IsNotFound(err) {
		return reconcileSecret(r, instance, logger, err, ctx)
	} else if err != nil {
		logger.Error(err, "Failed to get Secret")
		return ctrl.Result{}, err
	}

	// Add new service if not present
	serviceFound := &corev1.Service{}
	err = r.Get(ctx, types.NamespacedName{Name: svcName, Namespace: instance.Namespace}, serviceFound)
	if err != nil && errors.IsNotFound(err) {
		return reconcileSvc(r, instance, logger, err, ctx)
	} else if err != nil {
		logger.Error(err, "Failed to get Service")
		return ctrl.Result{}, err
	}

	// Update service details in status
	instance.Status.ServiceName = serviceFound.Name

	// Check for ingress
	ingressFound := networkingv1.Ingress{}
	logger.Info("ingress name: " + ingName + "namespace: " + instance.Namespace)
	err = r.Get(ctx, types.NamespacedName{Name: ingName, Namespace: instance.Namespace}, &ingressFound)
	if err != nil && errors.IsNotFound(err) {
		logger.Info("Unable to detect Ingress in your cluster. You may configure your own")
		return ctrl.Result{}, nil
	} else if err != nil {
		logger.Error(err, "Failed to get Ingress")
		return ctrl.Result{}, err
	}

	// Update ingress details in status
	if len(ingressFound.Status.LoadBalancer.Ingress) > 0 {
		instance.Status.IngressHostname = ingressFound.Status.LoadBalancer.Ingress[0].Hostname
	}

	// Check if the deployment already exists, if not create a new one
	depFound := &appsv1.Deployment{}
	err = r.Get(ctx, types.NamespacedName{Name: instance.Name, Namespace: instance.Namespace}, depFound)
	if err != nil && errors.IsNotFound(err) {
		return reconcileDeployment(r, instance, logger, err, ctx)
	} else if err != nil {
		logger.Error(err, "Failed to get Deployment")
		return ctrl.Result{}, err
	}

	// Update replica status
	instance.Status.Replicas = fmt.Sprint(depFound.Spec.Replicas)

	// Ensure the deployment size is the same as the spec
	size := instance.Spec.Size

	if *depFound.Spec.Replicas != size {
		depFound.Spec.Replicas = &size
		err = r.Update(ctx, depFound)
		if err != nil {
			logger.Error(err, "Failed to update Deployment", "Deployment.Namespace", depFound.Namespace, "Deployment.Name", depFound.Name)
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
		logger.Error(err, "Failed to list pods", "WSO2IS.Namespace", instance.Namespace, "WSO2IS.Name", instance.Name)
		return ctrl.Result{}, err
	}
	podNames := getPodNames(podList.Items)

	// Update status.Nodes if needed
	if !reflect.DeepEqual(podNames, instance.Status.Nodes) {
		instance.Status.Nodes = podNames
		err := r.Status().Update(ctx, &instance)
		if err != nil {
			logger.Error(err, "Failed to update WSO2IS status")
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

func reconcileDeployment(r *Wso2IsReconciler, instance wso2v1beta1.Wso2Is, log logr.Logger, err error, ctx context.Context) (ctrl.Result, error) {
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
}

func reconcileSvc(r *Wso2IsReconciler, instance wso2v1beta1.Wso2Is, log logr.Logger, err error, ctx context.Context) (ctrl.Result, error) {
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
}

func reconcileCfg(r *Wso2IsReconciler, instance wso2v1beta1.Wso2Is, log logr.Logger, err error, ctx context.Context) (ctrl.Result, error) {
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
}

func reconcileSecret(r *Wso2IsReconciler, instance wso2v1beta1.Wso2Is, log logr.Logger, err error, ctx context.Context) (ctrl.Result, error) {
	// Define a new secret
	secret := r.addSecret(instance, log)
	log.Info("Creating a new Secret", "Secret.Namespace", secret.Namespace, "Secret.Name", secret.Name)
	err = r.Create(ctx, secret)
	if err != nil {
		log.Error(err, "Failed to create new Secret", "Secret.Namespace", secret.Namespace, "Secret.Name", secret.Name)
		return ctrl.Result{}, err
	} else {
		log.Info("Successfully created new Secret", "Secret.Namespace", secret.Namespace, "Secret.Name", secret.Name)
	}
	// Secret created successfully - return and requeue
	return ctrl.Result{Requeue: true}, nil
}

func reconcileSva(r *Wso2IsReconciler, instance wso2v1beta1.Wso2Is, log logr.Logger, err error, ctx context.Context) (ctrl.Result, error) {
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
}

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

// addNewService adds a new Service
func (r *Wso2IsReconciler) addNewService(m wso2v1beta1.Wso2Is) *corev1.Service {

	// Make Service type configurable
	serviceType := corev1.ServiceTypeNodePort
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
							Name: pvcName,
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
								Name:      pvcName,
								MountPath: "/home/wso2carbon/wso2is-5.11.0/repository/deployment/server/userstores",
							},
							{
								Name:        configMapName,
								MountPath:   "/home/wso2carbon/wso2-config-volume/repository/conf/deployment.toml",
								SubPathExpr: configFileName,
							},
							{
								Name:      secretName,
								MountPath: "/home/wso2carbon/wso2is-5.11.0/repository/resources/security/controller-keystores",
								ReadOnly:  true,
							},
						},
						LivenessProbe: &corev1.Probe{
							ProbeHandler: corev1.ProbeHandler{
								Exec: &corev1.ExecAction{
									Command: []string{"/bin/sh", "-c", "nc -z localhost " + fmt.Sprint(containerPortHttps)},
								},
							},
							InitialDelaySeconds: 250,
							PeriodSeconds:       10,
						},
						ReadinessProbe: &corev1.Probe{
							ProbeHandler: corev1.ProbeHandler{
								Exec: &corev1.ExecAction{
									Command: []string{"/bin/sh", "-c", "nc -z localhost " + fmt.Sprint(containerPortHttps)},
								},
							},
							InitialDelaySeconds: 250,
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
