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
