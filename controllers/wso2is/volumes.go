package wso2is

import (
	"context"

	"github.com/go-logr/logr"
	wso2v1beta1 "github.com/wso2/k8s-wso2is-operator/api/v1beta1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/types"
	ctrl "sigs.k8s.io/controller-runtime"
)

func reconcileVolume(r *Wso2IsReconciler, instance wso2v1beta1.Wso2Is, log logr.Logger, err error, ctx context.Context) (ctrl.Result, error) {
	// Check for persistent volume claim
	pvcFound := &corev1.PersistentVolumeClaim{}
	err = r.Get(ctx, types.NamespacedName{Name: instance.Name + "-userstore-pvc", Namespace: instance.Namespace}, pvcFound)
	if err != nil && errors.IsNotFound(err) {
		log.Info("Unable to detect PVC claim in your cluster. You may configure your own")
	} else if err != nil {
		log.Error(err, "Failed to get PersistentVolumeClaim")
		return ctrl.Result{Requeue: true}, err

	} else {
		log.Info("Found PVC")
	}

	return ctrl.Result{}, err
}

func MakeVolumes(instance wso2v1beta1.Wso2Is) []corev1.Volume {
	volumes := []corev1.Volume{
		{
			Name: instance.Name + "-pv",
			VolumeSource: corev1.VolumeSource{
				PersistentVolumeClaim: &corev1.PersistentVolumeClaimVolumeSource{
					ClaimName: instance.Name + "-userstore-pvc",
				},
			},
		},
		{
			Name: instance.Name + "-config",
			VolumeSource: corev1.VolumeSource{
				ConfigMap: &corev1.ConfigMapVolumeSource{
					LocalObjectReference: corev1.LocalObjectReference{
						//Name: variables.ConfigMapName,
						Name: getConfigMapName(instance),
					},
				},
			},
		},
		{
			Name: instance.Name + "-secret",
			VolumeSource: corev1.VolumeSource{
				Secret: &corev1.SecretVolumeSource{
					//TODO:
					SecretName: instance.Name + "-secret",
				},
			},
		},
	}

	return volumes
}
