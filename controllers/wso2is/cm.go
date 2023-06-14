package wso2is

import (
	"bytes"
	"context"
	"crypto/md5"
	"encoding/hex"
	"github.com/BurntSushi/toml"
	"github.com/go-logr/logr"
	wso2v1beta1 "github.com/wso2/k8s-wso2is-operator/api/v1beta1"
	"github.com/wso2/k8s-wso2is-operator/variables"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/apimachinery/pkg/util/json"
	"log"
	ctrl "sigs.k8s.io/controller-runtime"
)

func calculateConfigMapHash(configMap *corev1.ConfigMap) (string, error) {
	// Serialize the data in the ConfigMap
	dataBytes, err := json.Marshal(configMap.Data)
	if err != nil {
		return "", err
	}

	// Calculate the MD5 hash of the serialized data
	hasher := md5.New()
	_, err = hasher.Write(dataBytes)
	if err != nil {
		return "", err
	}
	hashBytes := hasher.Sum(nil)

	// Convert the hash bytes to a hexadecimal string
	hashString := hex.EncodeToString(hashBytes)

	return hashString, nil
}

func getTomlConfig(spec wso2v1beta1.Wso2IsSpec) string {
	if len(spec.TomlConfig) == 0 {
		buf := new(bytes.Buffer)
		if err := toml.NewEncoder(buf).Encode(spec.Configurations); err != nil {
			log.Println(err)
		}
		//logger.Info(buf.String())
		return buf.String()
	} else {
		return spec.TomlConfig
	}
}

func remountConfigMap(r *Wso2IsReconciler, ctx context.Context, log logr.Logger, instance wso2v1beta1.Wso2Is, configMap *corev1.ConfigMap, configMapRefName string) (ctrl.Result, error) {
	sfs := &appsv1.StatefulSet{}
	err := r.Get(ctx, types.NamespacedName{Name: instance.Name, Namespace: instance.Namespace}, sfs)
	if err != nil {
		log.Error(err, "Error getting the StatefulSet")
		return ctrl.Result{}, err
	}
	volumeIndex := getVolumeIndex(sfs.Spec.Template.Spec.Volumes, variables.ConfigMapName)
	if volumeIndex != -1 {
		// Change the ConfigMap which acts as the source for the Volume with name `variables.ConfigMap` to the default ConfigMap by getting inline configs.
		if sfs.Spec.Template.Spec.Volumes[volumeIndex].VolumeSource.ConfigMap.LocalObjectReference.Name != configMapRefName {
			sfs.Spec.Template.Spec.Volumes[volumeIndex].VolumeSource.ConfigMap.LocalObjectReference.Name = configMapRefName
		}
		currentHash, err := calculateConfigMapHash(configMap)
		if err != nil {
			log.Error(err, "Failed to calculate ConfigMap hash")
		} else if sfs.Spec.Template.Annotations["configmapHash"] != currentHash {
			log.Info("A change is observed in the content of ConfigMap")
			sfs.Spec.Template.SetAnnotations(map[string]string{
				"configmapHash": currentHash,
			})
			if err := r.Update(ctx, sfs); err != nil {
				return ctrl.Result{}, err
			}
		}
	}

	return ctrl.Result{}, err
}

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
			"deployment.toml": getTomlConfig(m.Spec),
		},
	}
	ctrl.SetControllerReference(&m, configMap, r.Scheme)
	return configMap
}

func reconcileConfigMap(r *Wso2IsReconciler, instance wso2v1beta1.Wso2Is, log logr.Logger, err error, ctx context.Context) (ctrl.Result, error) {
	if instance.Spec.TomlConfigFile != "" {
		// If configs are specified using the ConfigMap reference
		log.Info("ConfigMap ref found in CRD yaml")

		configMap := &corev1.ConfigMap{}
		err = r.Get(ctx, types.NamespacedName{Name: instance.Spec.TomlConfigFile, Namespace: instance.Namespace}, configMap)
		if err != nil {
			log.Error(err, "Failed to get ConfigMap with the given name in CRD yaml")
			return ctrl.Result{}, err
		}
		sfs := &appsv1.StatefulSet{}
		err = r.Get(ctx, types.NamespacedName{Name: instance.Name, Namespace: instance.Namespace}, sfs)
		if err != nil {
			log.Info("Couldn't obtain StatefulSet. It'll be found in next reconcile loop if this is the first run.")
		} else {
			remountConfigMap(r, ctx, log, instance, configMap, instance.Spec.TomlConfigFile)
		}
	} else {
		// Using inline configs because ConfigMap ref not found.
		log.Info("ConfigMap ref is not found in CRD yaml. Checking for inline toml configs")

		configMapDefinition := r.defineConfigMap(instance)
		remountConfigMap(r, ctx, log, instance, configMapDefinition, variables.ConfigMapName)
	}

	return ctrl.Result{}, nil
}
