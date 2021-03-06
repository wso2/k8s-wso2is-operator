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

package v1beta1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

// Wso2IsSpec defines the desired state of Wso2Is
type Wso2IsSpec struct {
	Size int32 `json:"replicas"`
	// +kubebuilder:default:="5.11.0"
	Version        string          `json:"version,omitempty"`
	Configurations Configurations  `json:"configurations"`
	TomlConfig     string          `json:"tomlConfig,omitempty"`
	KeystoreMounts []KeystoreMount `json:"keystoreMounts,omitempty"`
}

type KeystoreMount struct {
	Name string `json:"name"`
	Data string `json:"data"`
}

type Configurations struct {
	Host string `json:"host"`
	// +kubebuilder:default:="NodePort"
	ServiceType string `json:"serviceType,omitempty"`
	// +kubebuilder:default:={ "hostname" : "$env{HOST_NAME}", "nodeIp": "$env{NODE_IP}" }
	Server Server `json:"server,omitempty" toml:"server"`
	// +kubebuilder:default:={ "username" : "admin", "password": "admin", "createAdminAccount": true }
	SuperAdmin SuperAdmin `json:"superAdmin,omitempty" toml:"super_admin"`
	// +kubebuilder:default:={ "type" : "read_write_ldap_unique_id", "connection_url": "ldap://localhost:${Ports.EmbeddedLDAP.LDAPServerPort}", "connection_name" : "uid=admin,ou=system", "connection_password" : "admin", "base_dn" : "dc=wso2,dc=org" }
	UserStore UserStore `json:"userStore,omitempty" toml:"user_store"`
	// +kubebuilder:default:={ "identityDb":{"password":"wso2carbon","type":"h2","url":"jdbc:h2:./repository/database/WSO2IDENTITY_DB;DB_CLOSE_ON_EXIT=FALSE;LOCK_TIMEOUT=60000","username":"wso2carbon"},"sharedDb":{"password":"wso2carbon","type":"h2","url":"jdbc:h2:./repository/database/WSO2SHARED_DB;DB_CLOSE_ON_EXIT=FALSE;LOCK_TIMEOUT=60000","username":"wso2carbon"} }
	Database Database `json:"database,omitempty" toml:"database"`
	// +kubebuilder:default:={ "https" : { "properties" : { "proxyPort" : 443 } } }
	Transport Transport `json:"transport,omitempty" toml:"transport,omitempty"`
	// +kubebuilder:default:= { "primary":{"name":"wso2carbon.jks","password":"wso2carbon"} }
	Keystore   Keystore   `json:"keystore,omitempty" toml:"keystore,omitempty"`
	Clustering Clustering `json:"clustering,omitempty" toml:"clustering,omitempty"`
	// +kubebuilder:default:={ "jmx" : { "rmi_server_start" : true } }
	Monitoring Monitoring `json:"monitoring,omitempty" toml:"monitoring,omitempty"`
	// +kubebuilder:default:={ "hazelcastShutdownhookEnabled" : false, "hazelcastLoggingType" : "log4j"  }
	Hazelcast Hazelcast `json:"hazelcast,omitempty" toml:"hazelcast,omitempty"`
}
type Server struct {
	// +kubebuilder:default:="$env{HOST_NAME}"
	Hostname string `json:"hostname,omitempty" toml:"hostname"`
	// +kubebuilder:default:="$env{NODE_IP}"
	NodeIP string `json:"nodeIp,omitempty" toml:"node_ip"`
}
type SuperAdmin struct {
	// +kubebuilder:default:="admin"
	Username string `json:"username,omitempty" toml:"username"`
	// +kubebuilder:default:="admin"
	Password string `json:"password,omitempty" toml:"password"`
	// +kubebuilder:default:=true
	CreateAdminAccount bool `json:"createAdminAccount,omitempty" toml:"create_admin_account"`
}
type UserStore struct {
	Type               string `json:"type" toml:"type"`
	ConnectionURL      string `json:"connection_url,omitempty" toml:"connection_url,omitempty"`
	ConnectionName     string `json:"connection_name,omitempty" toml:"connection_name,omitempty"`
	ConnectionPassword string `json:"connection_password,omitempty" toml:"connection_password,omitempty"`
	BaseDN             string `json:"base_dn,omitempty" toml:"base_dn,omitempty"`
}
type Hazelcast struct {
	// +kubebuilder:default:=false
	ShutdownHookEnabled bool `json:"hazelcastShutdownhookEnabled,omitempty" toml:"hazelcast.shutdownhook.enabled"`
	// +kubebuilder:default:="log4j"
	LoggingType string `json:"hazelcastLoggingType,omitempty" toml:"hazelcast.logging.type"`
}
type PoolOptions struct {
	// +kubebuilder:default:="SELECT 1"
	ValidationQuery string `json:"validationQuery,omitempty" toml:"validationQuery"`
}
type User struct {
	Type     string `json:"type,omitempty" toml:"type,omitempty"`
	URL      string `json:"url,omitempty" toml:"url,omitempty"`
	Hostname string `json:"hostname,omitempty" toml:"hostname,omitempty"`
	Username string `json:"username" toml:"username"`
	Password string `json:"password" toml:"password"`
	Driver   string `json:"driver,omitempty" toml:"driver,omitempty"`
	// +kubebuilder:default:={ "validationQuery" : "SELECT 1" }
	PoolOptions PoolOptions `json:"pool_options,omitempty" toml:"pool_options,omitempty"`
}
type IdentityDb struct {
	Type     string `json:"type,omitempty" toml:"type,omitempty"`
	URL      string `json:"url,omitempty" toml:"url,omitempty"`
	Hostname string `json:"hostname,omitempty" toml:"hostname,omitempty"`
	Username string `json:"username" toml:"username"`
	Password string `json:"password" toml:"password"`
	Driver   string `json:"driver,omitempty" toml:"driver,omitempty"`
	// +kubebuilder:default:={ "validationQuery" : "SELECT 1" }
	PoolOptions PoolOptions `json:"pool_options,omitempty" toml:"pool_options,omitempty"`
}
type SharedDb struct {
	Type     string `json:"type,omitempty" toml:"type,omitempty"`
	URL      string `json:"url,omitempty" toml:"url,omitempty"`
	Hostname string `json:"hostname,omitempty" toml:"hostname,omitempty"`
	Username string `json:"username" toml:"username"`
	Password string `json:"password" toml:"password"`
	Driver   string `json:"driver,omitempty" toml:"driver,omitempty"`
	// +kubebuilder:default:={ "validationQuery" : "SELECT 1" }
	PoolOptions PoolOptions `json:"pool_options,omitempty" toml:"pool_options,omitempty"`
}
type BpsDatabase struct {
	Type     string `json:"type,omitempty" toml:"type,omitempty"`
	URL      string `json:"url,omitempty" toml:"url,omitempty"`
	Hostname string `json:"hostname,omitempty" toml:"hostname,omitempty"`
	Username string `json:"username" toml:"username"`
	Password string `json:"password" toml:"password"`
	Driver   string `json:"driver,omitempty" toml:"driver,omitempty"`
	// +kubebuilder:default:={ "validationQuery" : "SELECT 1" }
	PoolOptions PoolOptions `json:"pool_options,omitempty" toml:"pool_options,omitempty"`
}
type Database struct {
	IdentityDb IdentityDb `json:"identityDb" toml:"identity_db"`
	SharedDb   SharedDb   `json:"sharedDb" toml:"shared_db"`
}
type Properties struct {
	ProxyPort int `json:"proxyPort" toml:"proxyPort"`
}
type HTTPS struct {
	Properties Properties `json:"properties" toml:"properties"`
}
type Transport struct {
	HTTPS HTTPS `json:"https" toml:"https"`
}
type Datasource struct {
	ID                         string `json:"id" toml:"id"`
	URL                        string `json:"url" toml:"url"`
	Username                   string `json:"username" toml:"username"`
	Password                   string `json:"password" toml:"password"`
	Driver                     string `json:"driver" toml:"driver"`
	PoolOptionsValidationQuery string `json:"pool_options_validationQuery" toml:"pool_options.validationQuery"`
	PoolOptionsMaxActive       int    `json:"pool_options_maxActive" toml:"pool_options.maxActive"`
	PoolOptionsMaxWait         int    `json:"pool_options_maxWait" toml:"pool_options.maxWait"`
	PoolOptionsTestOnBorrow    bool   `json:"pool_options_testOnBorrow" toml:"pool_options.testOnBorrow"`
	PoolOptionsJmxEnabled      bool   `json:"pool_options_jmxEnabled" toml:"pool_options.jmxEnabled"`
}
type Consent struct {
	DataSource string `json:"data_source" toml:"data_source"`
}
type Authentication struct {
	Consent Consent `json:"consent" toml:"consent"`
}
type Primary struct {
	Name     string `json:"name" toml:"name"`
	Password string `json:"password" toml:"password"`
}
type Keystore struct {
	Primary Primary `json:"primary" toml:"primary"`
}
type Clustering struct {
	// +kubebuilder:default:="kubernetes"
	MembershipScheme string `json:"membership_scheme,omitempty" toml:"membership_scheme,omitempty"`
	// +kubebuilder:default:="wso2.is.domain"
	Domain     string               `json:"domain,omitempty" toml:"domain,omitempty"`
	Properties ClusteringProperties `json:"properties" toml:"properties,omitempty"`
}
type ClusteringProperties struct {
	// +kubebuilder:default:="org.wso2.carbon.membership.scheme.kubernetes.KubernetesMembershipScheme"
	PropertiesMembershipSchemeClassName string `json:"membershipSchemeClassName,omitempty" toml:"membershipSchemeClassName,omitempty"`
	// +kubebuilder:default:="default"
	PropertiesKUBERNETESNAMESPACE string `json:"KUBERNETES_NAMESPACE,omitempty" toml:"KUBERNETES_NAMESPACE,omitempty"`
	// +kubebuilder:default:="wso2is-service"
	PropertiesKUBERNETESSERVICES string `json:"KUBERNETES_SERVICES,omitempty" toml:"KUBERNETES_SERVICES,omitempty"`
	// +kubebuilder:default:=true
	PropertiesKUBERNETESMASTERSKIPSSLVERIFICATION bool `json:"KUBERNETES_MASTER_SKIP_SSL_VERIFICATION,omitempty" toml:"KUBERNETES_MASTER_SKIP_SSL_VERIFICATION,omitempty"`
	// +kubebuilder:default:=false
	PropertiesUSEDNS                bool   `json:"USE_DNS,omitempty" toml:"USE_DNS,omitempty"`
	PropertiesKUBERNETES_API_SERVER string `json:"KUBERNETES_API_SERVER" toml:"KUBERNETES_API_SERVER,omitempty"`
}
type Jmx struct {
	// +kubebuilder:default:=true
	RmiServerStart bool `toml:"rmi_server_start" json:"rmi_server_start,omitempty"`
}
type Monitoring struct {
	Jmx Jmx `toml:"jmx" json:"jmx,omitempty"`
}

// Wso2IsStatus defines the observed state of Wso2Is
type Wso2IsStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "make" to regenerate code after modifying this file
	Nodes           []string `json:"nodes" toml:"nodes"`
	ServiceName     string   `json:"serviceName"`
	IngressHostname string   `json:"ingressHostname"`
	Replicas        string   `json:"replicas"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
// +k8s:openapi-gen=true
// +kubebuilder:subresource:status
// +kubebuilder:resource:scope=Namespaced
// +kubebuilder:printcolumn:name="Service Name",type="string",JSONPath=`.status.serviceName`
// +kubebuilder:printcolumn:name="Ingress Hostname",type="string",JSONPath=`.status.ingressHostname`
// +kubebuilder:printcolumn:name="Desired Replicas",type="string",JSONPath=`.spec.replicas`
// +kubebuilder:printcolumn:name="Current Replicas",type="string",JSONPath=`.status.replicas`
// +kubebuilder:printcolumn:name="Host Name",type="string",JSONPath=`.spec.configurations.host`

// Wso2Is is the Schema for the wso2is API
type Wso2Is struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   Wso2IsSpec   `json:"spec,omitempty"`
	Status Wso2IsStatus `json:"status,omitempty"`
}

// +kubebuilder:object:root=true

// Wso2IsList contains a list of Wso2Is
type Wso2IsList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []Wso2Is `json:"items"`
}

func init() {
	SchemeBuilder.Register(&Wso2Is{}, &Wso2IsList{})
}
