package wso2is

import (
	"context"

	"github.com/go-logr/logr"
	wso2v1beta1 "github.com/wso2/k8s-wso2is-operator/api/v1beta1"
	"github.com/wso2/k8s-wso2is-operator/variables"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/types"
	ctrl "sigs.k8s.io/controller-runtime"
)

func reconcileVolume(r *Wso2IsReconciler, instance wso2v1beta1.Wso2Is, log logr.Logger, err error, ctx context.Context) (ctrl.Result, error) {
	// Check for persistent volume claim
	pvcFound := &corev1.PersistentVolumeClaim{}
	err = r.Get(ctx, types.NamespacedName{Name: variables.UserstorePVCName, Namespace: instance.Namespace}, pvcFound)
	if err != nil && errors.IsNotFound(err) {
		log.Info("Unable to detect PVC claim in your cluster. You may configure your own")
	} else if err != nil {
		log.Error(err, "Failed to get PersistentVolumeClaim")

	} else {
		log.Info("Found PVC")
	}

	return ctrl.Result{}, err
}

//
//func remountVolume(r *Wso2IsReconciler, instance wso2v1beta1.Wso2Is, log logr.Logger, ctx context.Context) (ctrl.Result, error) {
//	// log.Info("Remounting volume")
//
//	// Get the ConfigMap
//	configMap := &corev1.ConfigMap{}
//	err := r.Get(ctx, types.NamespacedName{Name: variables.ConfigMapName, Namespace: instance.Namespace}, configMap)
//	if err != nil {
//		log.Error(err, "Failed to get ConfigMap")
//		return ctrl.Result{}, err
//	}
//
//	// Update the ConfigMap data with the new content
//	configMap.Data["deployment.toml"] = getTomlConfig(instance.Spec)
//
//	// Update the ConfigMap
//	err = r.Update(ctx, configMap)
//	if err != nil {
//		log.Error(err, "Failed to update ConfigMap")
//		return ctrl.Result{}, err
//	}
//
//	wso2Is := &wso2v1beta1.Wso2Is{}
//	err = r.Get(ctx, types.NamespacedName{Namespace: "wso2-iam-system", Name: "identity-server-test"}, wso2Is)
//	if err != nil {
//		log.Error(err, "Failed to get Wso2Is resource")
//		return ctrl.Result{}, err
//	}
//	wso2Is.Spec.Template.Annotations["configmapHash"] = calculateConfigMapHash(configMap)
//	err = r.Update(ctx, wso2Is) //This doesn't work 2nd time?
//	if err != nil {
//		log.Error(err, "Failed to update Wso2Is resource")
//		return ctrl.Result{}, err
//	}
//
//	log.Info("ConfigMap updated successfully")
//
//	// Return reconcile result
//	return ctrl.Result{}, nil
//}
