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
	"context"
	"fmt"
	appsv1 "k8s.io/api/apps/v1"
	rbacv1 "k8s.io/api/rbac/v1"
	"reflect"

	wso2v1beta1 "github.com/wso2/k8s-wso2is-operator/api/v1beta1"
	corev1 "k8s.io/api/core/v1"
	networkingv1 "k8s.io/api/networking/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/types"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

// +kubebuilder:rbac:groups=iam.wso2.com,resources=wso2is,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=iam.wso2.com,resources=wso2is/status,verbs=get;update;patch

func (r *Wso2IsReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {

	// Get logger
	logger := r.Log.WithValues(deploymentName, req.NamespacedName)

	// Fetch the WSO2IS instance
	instance := wso2v1beta1.Wso2Is{}
	err := r.Get(ctx, req.NamespacedName, &instance)

	// Check if WSO2 custom resource is present
	err = r.Get(ctx, req.NamespacedName, &instance)
	if err != nil {
		if errors.IsNotFound(err) {
			// Request object not found, could have been deleted after reconcile request.
			// Owned objects are automatically garbage collected. For additional cleanup logic use finalizers.
			// Return and don't requeue
			logger.Info("WSO2IS resource not found. Ignoring since object must be deleted")
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
	// Check if the tomlConfig field has changed
	if confMap.Data[configFileName] != instance.Spec.TomlConfig {
		logger.Info("Updating the ConfigMap due to tomlConfig change")
		remountVolume(r, instance, logger, ctx)
		if err != nil {
			logger.Error(err, "Failed to update ConfigMap")
			return ctrl.Result{}, err
		}
		logger.Info("ConfigMap updated successfully")
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

	//---------------------------StatefulSet pattern-------------------------------------------
	// Check if the StatefulSet already exists, if not create a new one
	statefulSet := &appsv1.StatefulSet{}
	err = r.Get(ctx, types.NamespacedName{Name: instance.Name, Namespace: instance.Namespace}, statefulSet)
	if err != nil && errors.IsNotFound(err) {
		return reconcileStatefulSet(r, instance, logger, err, ctx)
	} else if err != nil {
		logger.Error(err, "Failed to get StatefulSet")
		return ctrl.Result{}, err
	}

	wso2Is := &wso2v1beta1.Wso2Is{}
	err = r.Get(ctx, types.NamespacedName{Name: instance.Name, Namespace: instance.Namespace}, wso2Is)
	if err != nil {
		logger.Error(err, "Failed to get Wso2Is resource")
		return ctrl.Result{}, err
	}

	if statefulSet.Spec.Template.ObjectMeta.Annotations["configmapHash"] != wso2Is.Spec.Template.Annotations["configmapHash"] {
		statefulSet.Spec.Template.ObjectMeta.Annotations["configmapHash"] = wso2Is.Spec.Template.Annotations["configmapHash"]
		err = r.Update(ctx, statefulSet)
		if err != nil {
			logger.Error(err, "Failed to update Wso2Is resource")
			return ctrl.Result{}, err
		}
	}

	// Update replica status
	instance.Status.Replicas = fmt.Sprint(*statefulSet.Spec.Replicas)

	// Ensure the StatefulSet size is the same as the spec
	size := instance.Spec.Size

	if *statefulSet.Spec.Replicas != size {
		statefulSet.Spec.Replicas = &size
		err = r.Update(ctx, statefulSet)
		if err != nil {
			logger.Error(err, "Failed to update StatefulSet", "StatefulSet.Namespace", statefulSet.Namespace, "StatefulSet.Name", statefulSet.Name)
			return ctrl.Result{}, err
		}
		// Spec updated - return and requeue
		return ctrl.Result{Requeue: true}, nil
	}

	//---------------------------StatefulSet pattern-------------------------------------------

	//---------------------------Deployment pattern-------------------------------------------
	//// Check if the deployment already exists, if not create a new one
	//depFound := &appsv1.Deployment{}
	//err = r.Get(ctx, types.NamespacedName{Name: instance.Name, Namespace: instance.Namespace}, depFound)
	//if err != nil && errors.IsNotFound(err) {
	//	return reconcileDeployment(r, instance, logger, err, ctx)
	//} else if err != nil {
	//	logger.Error(err, "Failed to get Deployment")
	//	return ctrl.Result{}, err
	//}
	//
	//wso2Is := &wso2v1beta1.Wso2Is{}
	//err = r.Get(ctx, types.NamespacedName{Namespace: "wso2-iam-system", Name: "identity-server-test"}, wso2Is)
	//if err != nil {
	//	logger.Error(err, "Failed to get Wso2Is resource")
	//	return ctrl.Result{}, err
	//}
	//
	//if depFound.Spec.Template.ObjectMeta.Annotations["configmapHash"] != wso2Is.Spec.Template.Annotations["configmapHash"] {
	//	depFound.Spec.Template.ObjectMeta.Annotations["configmapHash"] = wso2Is.Spec.Template.Annotations["configmapHash"]
	//	err = r.Update(ctx, depFound)
	//	if err != nil {
	//		logger.Error(err, "Failed to update Wso2Is resource")
	//		return ctrl.Result{}, err
	//	}
	//}
	//
	//// Update replica status
	//instance.Status.Replicas = fmt.Sprint(depFound.Spec.Replicas)
	//
	//// Ensure the deployment size is the same as the spec
	//size := instance.Spec.Size
	//
	//if depFound.Spec.Replicas == nil || *depFound.Spec.Replicas != size {
	//	depFound.Spec.Replicas = &size
	//	err = r.Update(ctx, depFound)
	//	if err != nil {
	//		logger.Error(err, "Failed to update Deployment", "Deployment.Namespace", depFound.Namespace, "Deployment.Name", depFound.Name)
	//		return ctrl.Result{}, err
	//	}
	//	// Spec updated - return and requeue
	//	return ctrl.Result{Requeue: true}, nil
	//}
	//---------------------------Deployment pattern-------------------------------------------

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

	// Add new role if not present
	roleFound := &rbacv1.Role{}
	err = r.Get(ctx, types.NamespacedName{Name: roleName, Namespace: instance.Namespace}, roleFound)
	if err != nil && errors.IsNotFound(err) {
		return reconcileRole(r, instance, logger, err, ctx)
	} else if err != nil {
		logger.Error(err, "Failed to get Role")
		return ctrl.Result{}, err
	}

	// Update role details in status
	//instance.Status.ServiceName = serviceFound.Name

	// Add new role binding if not present
	roleBindingFound := &rbacv1.RoleBinding{}
	err = r.Get(ctx, types.NamespacedName{Name: roleBindingName, Namespace: instance.Namespace}, roleBindingFound)
	if err != nil && errors.IsNotFound(err) {
		return reconcileRoleBinding(r, instance, logger, err, ctx)
	} else if err != nil {
		logger.Error(err, "Failed to get RoleBinding")
		return ctrl.Result{}, err
	}

	// Update role binding details in status
	//instance.Status.ServiceName = serviceFound.Name

	return ctrl.Result{}, nil

}
