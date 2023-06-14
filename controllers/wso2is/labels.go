package wso2is

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
