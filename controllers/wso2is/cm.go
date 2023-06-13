package wso2is

import (
	"context"
	"github.com/go-logr/logr"
	wso2v1beta1 "github.com/wso2/k8s-wso2is-operator/api/v1beta1"
	"github.com/wso2/k8s-wso2is-operator/variables"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	ctrl "sigs.k8s.io/controller-runtime"
)

// Find the index of the volume with the given name
func getVolumeIndex(volumes []corev1.Volume, volumeName string) int {
	for i, v := range volumes {
		if v.Name == volumeName {
			return i
		}
	}
	return -1
}

func (r *Wso2IsReconciler) defineConfigMap(m wso2v1beta1.Wso2Is) *corev1.ConfigMap {
	configMap := &corev1.ConfigMap{
		ObjectMeta: metav1.ObjectMeta{
			Name:      variables.ConfigMapName,
			Namespace: m.Namespace,
		},
		Data: map[string]string{
			//variables.ConfigMapName: getTomlConfig(m.Spec),
			"deployment.toml": getTomlConfig(m.Spec),
		},
	}
	ctrl.SetControllerReference(&m, configMap, r.Scheme)
	return configMap
}

func reconcileConfigMap(r *Wso2IsReconciler, instance wso2v1beta1.Wso2Is, log logr.Logger, err error, ctx context.Context) (ctrl.Result, error) {

	configMap := &corev1.ConfigMap{}

	if instance.Spec.TomlConfigFile != "" {
		// If configs are specified using the ConfigMap reference
		log.Info("ConfigMap ref found in CRD yaml")

		err = r.Get(ctx, types.NamespacedName{Name: instance.Spec.TomlConfigFile, Namespace: instance.Namespace}, configMap)
		if err != nil {
			log.Error(err, "Failed to update ConfigMap")
			return ctrl.Result{}, err
		}
		sfs := &appsv1.StatefulSet{}
		err = r.Get(ctx, types.NamespacedName{Name: instance.Name, Namespace: instance.Namespace}, sfs)
		if err != nil {
			log.Info("Couldn't obtain StatefulSet. It'll be found in next reconcile loop if this is the first run.")
		} else {
			// If toml ConfigMap contents have changed,
			volumeIndex := getVolumeIndex(sfs.Spec.Template.Spec.Volumes, variables.ConfigMapName)

			log.Info("The volumeIndex current value is ", "volumeIndex", volumeIndex)

			if volumeIndex != -1 {
				// Volume exists, update its configuration
				// Change the ConfigMap which acts as the source for the Volume with name `variables.ConfigMap` to the one provided by user.
				if sfs.Spec.Template.Spec.Volumes[volumeIndex].VolumeSource.ConfigMap.LocalObjectReference.Name != instance.Spec.TomlConfigFile {
					sfs.Spec.Template.Spec.Volumes[volumeIndex].VolumeSource.ConfigMap.LocalObjectReference.Name = instance.Spec.TomlConfigFile
				}
				// If the name is not changed, but the content is changed,
				currentHash, err := calculateConfigMapHash(configMap)
				log.Info("The annotation current value is ", "configmapHash", sfs.Spec.Template.Annotations["configmapHash"])
				log.Info("The currentHash current value is ", "currentHash", currentHash)
				if err != nil {
					log.Error(err, "Failed to calculate ConfigMap hash")
				} else if sfs.Spec.Template.Annotations["configmapHash"] != currentHash {
					log.Info("A change is observed in the content of ConfigMap")
					// Perform actions when the ConfigMap content has changed
					// Update the stored hash with the new calculated hash
					sfs.Spec.Template.SetAnnotations(map[string]string{
						"configmapHash": currentHash,
					})
					// Update the object
					if err := r.Update(ctx, sfs); err != nil {
						return ctrl.Result{}, err
					}
				} else {
					log.Info("No change happened to the ConfigMap")
				}
			} else {
				// Volume does not exist, add it
				volume := corev1.Volume{
					Name: variables.ConfigMapName,
					VolumeSource: corev1.VolumeSource{
						ConfigMap: &corev1.ConfigMapVolumeSource{
							LocalObjectReference: corev1.LocalObjectReference{
								Name: variables.ConfigMapName,
							},
						},
					},
				}
				sfs.Spec.Template.Spec.Volumes = append(sfs.Spec.Template.Spec.Volumes, volume)
			}
		}
	} else {
		configMapDefinition := r.defineConfigMap(instance)
		err = r.Get(ctx, types.NamespacedName{Name: variables.ConfigMapName, Namespace: instance.Namespace}, configMap)
		if err != nil {
			if errors.IsNotFound(err) {
				log.Info("ConfigMap resource " + variables.ConfigMapName + " not found. Creating or re-creating configMap")
				err = r.Create(ctx, configMapDefinition)
				if err != nil {
					log.Error(err, "Failed to create new ConfigMap", "ConfigMap.Namespace", configMapDefinition.Namespace, "ConfigMap.Name", configMapDefinition.Name)
					return ctrl.Result{}, err
				}
			} else {
				log.Info("Failed to get configMap resource " + variables.ConfigMapName + ". Re-running reconcile.")
				return ctrl.Result{}, err
			}
		} else {
			// Note: For simplication purposes ConfigMaps are not updated - see deployment section
		}

	}

	return ctrl.Result{}, nil
}
