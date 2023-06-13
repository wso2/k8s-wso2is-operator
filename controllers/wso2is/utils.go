package wso2is

import (
	"bytes"
	"crypto/md5"
	"encoding/hex"
	"k8s.io/apimachinery/pkg/util/json"
	"log"

	"github.com/BurntSushi/toml"
	wso2v1beta1 "github.com/wso2/k8s-wso2is-operator/api/v1beta1"
	corev1 "k8s.io/api/core/v1"
)

//	func calculateConfigMapHash(configMap *corev1.ConfigMap) string {
//		data := ""
//		for _, value := range configMap.Data {
//			data += value
//		}
//
//		hash := md5.Sum([]byte(data))
//		return hex.EncodeToString(hash[:])
//	}
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
