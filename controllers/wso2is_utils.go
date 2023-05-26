package controllers

import (
	"bytes"
	"crypto/md5"
	"encoding/hex"
	"github.com/BurntSushi/toml"
	"github.com/go-logr/logr"
	wso2v1beta1 "github.com/wso2/k8s-wso2is-operator/api/v1beta1"
	corev1 "k8s.io/api/core/v1"
	"log"
	ctrl "sigs.k8s.io/controller-runtime"
)

func (r *Wso2IsReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&wso2v1beta1.Wso2Is{}).
		Complete(r)
}

func calculateConfigMapHash(configMap *corev1.ConfigMap) string {
	data := ""
	for _, value := range configMap.Data {
		data += value
	}

	hash := md5.Sum([]byte(data))
	return hex.EncodeToString(hash[:])
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
