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
	// +kubebuilder:default:={ "type" : "database_unique_id"}
	UserStore UserStore `json:"userStore,omitempty" toml:"user_store"`
	// +kubebuilder:default:={ "identityDb":{"password":"wso2carbon","type":"h2","url":"jdbc:h2:./repository/database/WSO2IDENTITY_DB;DB_CLOSE_ON_EXIT=FALSE;LOCK_TIMEOUT=60000","username":"wso2carbon"},"sharedDb":{"password":"wso2carbon","type":"h2","url":"jdbc:h2:./repository/database/WSO2SHARED_DB;DB_CLOSE_ON_EXIT=FALSE;LOCK_TIMEOUT=60000","username":"wso2carbon"} }
	Database Database `json:"database,omitempty" toml:"database"`
	// +kubebuilder:default:={ "https" : { "properties" : { "proxyPort" : 443 } } }
	Transport Transport `json:"transport,omitempty" toml:"transport,omitempty"`
	// +kubebuilder:default:= { "primary":{"name":"wso2carbon.jks","password":"wso2carbon"} }
	Keystore Keystore `json:"keystore,omitempty" toml:"keystore,omitempty"`
	// +kubebuilder:default:={ "jmx" : { "rmi_server_start" : true } }
	Monitoring Monitoring `json:"monitoring,omitempty" toml:"monitoring,omitempty"`
	// +kubebuilder:default:={ "hazelcastShutdownhookEnabled" : false, "hazelcastLoggingType" : "log4j"  }
	Hazelcast      Hazelcast          `json:"hazelcast,omitempty" toml:"hazelcast,omitempty"`
	Authentication StepAuthentication `json:"authentication,omitempty" toml:"authentication,omitempty"`
	Recaptcha      Recaptcha          `json:"recaptcha,omitempty" toml:"recaptcha,omitempty"`
	OutputAdapter  OutputAdapter      `toml:"output_adapter,omitempty" json:"output_adapter,omitempty"`
	Clustering     Clustering         `json:"clustering,omitempty" toml:"clustering,omitempty"`
	TenantMgt      TenantMgt          `json:"tenant_mgt,omitempty" toml:"tenant_mgt,omitempty"`
	// +kubebuilder:default:={ "enableTenantQualifiedUrls":true }
	TenantCtx    TenantCtx    `json:"tenant_context,omitempty" toml:"tenant_context,omitempty"`
	AdminService AdminService `json:"admin_service,omitempty" toml:"admin_service,omitempty"`
}

/* Primary server configs */
type Server struct {
	// +kubebuilder:default:="$env{HOST_NAME}"
	Hostname string `json:"hostname,omitempty" toml:"hostname"`
	// +kubebuilder:default:="$env{NODE_IP}"
	NodeIP string `json:"nodeIp,omitempty" toml:"node_ip"`
}

/* Create admin accounts */
type SuperAdmin struct {
	// +kubebuilder:default:="admin"
	Username string `json:"username,omitempty" toml:"username"`
	// +kubebuilder:default:="admin"
	Password string `json:"password,omitempty" toml:"password"`
	// +kubebuilder:default:=true
	CreateAdminAccount bool `json:"createAdminAccount,omitempty" toml:"create_admin_account"`
}

/* UserStore configs */
type UserStore struct {
	Type               string `json:"type" toml:"type"`
	ConnectionURL      string `json:"connection_url,omitempty" toml:"connection_url,omitempty"`
	ConnectionName     string `json:"connection_name,omitempty" toml:"connection_name,omitempty"`
	ConnectionPassword string `json:"connection_password,omitempty" toml:"connection_password,omitempty"`
	BaseDN             string `json:"base_dn,omitempty" toml:"base_dn,omitempty"`
	UsernameAttrib     string `json:"user_name_attribute,omitempty" toml:"user_name_attribute,omitempty"`
}

/* Hazelcast clustering configs */
type Hazelcast struct {
	// +kubebuilder:default:=false
	ShutdownHookEnabled bool `json:"hazelcastShutdownhookEnabled,omitempty" toml:"hazelcast.shutdownhook.enabled"`
	// +kubebuilder:default:="log4j"
	LoggingType string `json:"hazelcastLoggingType,omitempty" toml:"hazelcast.logging.type"`
}

/* MySQL pool options */
type PoolOptions struct {
	// +kubebuilder:default:="SELECT 1"
	ValidationQuery string `json:"validationQuery,omitempty" toml:"validationQuery"`
}

/* User database */
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

/* Identity database */
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

/* Shared database */
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

/* BPS database */
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

/* Database connections and configs */
type Database struct {
	IdentityDb IdentityDb `json:"identityDb" toml:"identity_db"`
	SharedDb   SharedDb   `json:"sharedDb" toml:"shared_db"`
}

/* Proxyport properties */
type Properties struct {
	ProxyPort int `json:"proxyPort" toml:"proxyPort"`
}

/* Transport protocol HTTP */
type HTTPS struct {
	Properties Properties `json:"properties" toml:"properties"`
}

/* Transport protocols */
type Transport struct {
	HTTPS HTTPS `json:"https" toml:"https"`
}

/* Single datasource */
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

/* Consent DB data source */
type Consent struct {
	DataSource string `json:"data_source" toml:"data_source"`
}

/* Consent DB configurations */
type Authentication struct {
	Consent Consent `json:"consent" toml:"consent"`
}

/* Primary Keystore configurations */
type Primary struct {
	Name     string `json:"name" toml:"name"`
	Password string `json:"password" toml:"password"`
}

/* Keystore configurations */
type Keystore struct {
	Primary Primary `json:"primary" toml:"primary"`
}

/* Clustering configurations */
type Clustering struct {
	// +kubebuilder:default:="kubernetes"
	MembershipScheme string `json:"membership_scheme,omitempty" toml:"membership_scheme,omitempty"`
	// +kubebuilder:default:="wso2.is.domain"
	Domain     string               `json:"domain,omitempty" toml:"domain,omitempty"`
	Properties ClusteringProperties `json:"properties" toml:"properties,omitempty"`
}

/* Clustering Properties*/
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

/* Jmx monitoring configurations */
type Jmx struct {
	// +kubebuilder:default:=true
	RmiServerStart bool `toml:"rmi_server_start" json:"rmi_server_start,omitempty"`
}

/* Monitoring configurations */
type Monitoring struct {
	Jmx Jmx `toml:"jmx" json:"jmx,omitempty"`
}

/* TOTP parameters */
type TotpParameters struct {
	// +kubebuilder:default:="Base32"
	EncodingMethod string `toml:"encodingMethod,omitempty" json:"encodingMethod,omitempty"`
	// +kubebuilder:default:="30"
	TimeStepSize string `toml:"timeStepSize,omitempty" json:"timeStepSize,omitempty"`
	// +kubebuilder:default:="3"
	WindowSize string `toml:"windowSize,omitempty" json:"windowSize,omitempty"`
	// +kubebuilder:default:=true
	AuthenticationMandatory bool `toml:"authenticationMandatory,omitempty" json:"authenticationMandatory,omitempty"`
	// +kubebuilder:default:=true
	EnrolUserInAuthenticationFlow bool `toml:"enrolUserInAuthenticationFlow,omitempty" json:"enrolUserInAuthenticationFlow,omitempty"`
	// +kubebuilder:default:="local"
	Usecase string `toml:"usecase,omitempty" json:"usecase,omitempty"`
	// +kubebuilder:default:="primary"
	SecondaryUserstore string `toml:"secondaryUserstore,omitempty" json:"secondaryUserstore,omitempty"`
	// +kubebuilder:default:="/totpauthenticationendpoint/totp.jsp"
	TOTPAuthenticationEndpointURL string `toml:"TOTPAuthenticationEndpointURL,omitempty" json:"TOTPAuthenticationEndpointURL,omitempty"`
	// +kubebuilder:default:="/totpauthenticationendpoint/totpError.jsp"
	TOTPAuthenticationEndpointErrorPage string `toml:"TOTPAuthenticationEndpointErrorPage,omitempty" json:"TOTPAuthenticationEndpointErrorPage,omitempty"`
	// +kubebuilder:default:="/totpauthenticationendpoint/enableTOTP.jsp"
	TOTPAuthenticationEndpointEnableTOTPPage string `toml:"TOTPAuthenticationEndpointEnableTOTPPage,omitempty" json:"TOTPAuthenticationEndpointEnableTOTPPage,omitempty"`
	// +kubebuilder:default:="WSO2"
	Issuer string `toml:"Issuer,omitempty" json:"Issuer,omitempty"`
	// +kubebuilder:default:=true
	UseCommonIssuer bool `toml:"UseCommonIssuer,omitempty" json:"UseCommonIssuer,omitempty"`
}

/* TOTP configurations */
type Totp struct {
	// +kubebuilder:default:={"Issuer":"WSO2","TOTPAuthenticationEndpointEnableTOTPPage":"/totpauthenticationendpoint/enableTOTP.jsp","TOTPAuthenticationEndpointErrorPage":"/totpauthenticationendpoint/totpError.jsp","TOTPAuthenticationEndpointURL":"/totpauthenticationendpoint/totp.jsp","UseCommonIssuer":true,"authenticationMandatory":true,"encodingMethod":"Base32","enrolUserInAuthenticationFlow":true,"secondaryUserstore":"primary","timeStepSize":"30","usecase":"local","windowSize":"3"}
	Parameters TotpParameters `toml:"parameters,omitempty" json:"parameters,omitempty"`
	// +kubebuilder:default:=false
	Enable bool `toml:"enable,omitempty" json:"enable,omitempty"`
}

/* Email authentication parameters */
type EmailOtpParameters struct {
	// +kubebuilder:default:="/emailotpauthenticationendpoint/emailotp.jsp"
	EMAILOTPAuthenticationEndpointURL string `toml:"EMAILOTPAuthenticationEndpointURL,omitempty" json:"EMAILOTPAuthenticationEndpointURL,omitempty"`
	// +kubebuilder:default:="/emailotpauthenticationendpoint/emailotpError.jsp"
	EmailOTPAuthenticationEndpointErrorPage string `toml:"EmailOTPAuthenticationEndpointErrorPage,omitempty" json:"EmailOTPAuthenticationEndpointErrorPage,omitempty"`
	// +kubebuilder:default:="/emailotpauthenticationendpoint/emailAddress.jsp"
	EmailAddressRequestPage string `toml:"EmailAddressRequestPage,omitempty" json:"EmailAddressRequestPage,omitempty"`
	// +kubebuilder:default:="local"
	Usecase string `toml:"usecase,omitempty" json:"usecase,omitempty"`
	// +kubebuilder:default:="primary"
	SecondaryUserstore string `toml:"secondaryUserstore,omitempty" json:"secondaryUserstore,omitempty"`
	// +kubebuilder:default:=false
	EMAILOTPMandatory bool `toml:"EMAILOTPMandatory,omitempty" json:"EMAILOTPMandatory,omitempty"`
	// +kubebuilder:default:=false
	SendOTPToFederatedEmailAttribute bool `toml:"sendOTPToFederatedEmailAttribute,omitempty" json:"sendOTPToFederatedEmailAttribute,omitempty"`
	// +kubebuilder:default:="email"
	FederatedEmailAttributeKey string `toml:"federatedEmailAttributeKey,omitempty" json:"federatedEmailAttributeKey,omitempty"`
	// +kubebuilder:default:=true
	EmailOTPEnableByUserClaim bool `toml:"EmailOTPEnableByUserClaim,omitempty" json:"EmailOTPEnableByUserClaim,omitempty"`
	// +kubebuilder:default:=true
	CaptureAndUpdateEmailAddress bool `toml:"CaptureAndUpdateEmailAddress,omitempty" json:"CaptureAndUpdateEmailAddress,omitempty"`
	// +kubebuilder:default:=true
	ShowEmailAddressInUI bool `toml:"showEmailAddressInUI,omitempty" json:"showEmailAddressInUI,omitempty"`
	// +kubebuilder:default:=true
	UseEventHandlerBasedEmailSender bool `toml:"useEventHandlerBasedEmailSender,omitempty" json:"useEventHandlerBasedEmailSender,omitempty"`
}

/* Enable email authentication */
type EmailOtp struct {
	// +kubebuilder:default:="EmailOTP"
	Name string `toml:"name,omitempty" json:"name,omitempty"`
	// +kubebuilder:default:=false
	Enable bool `toml:"enable,omitempty" json:"enable,omitempty"`
	// +kubebuilder:default:={"CaptureAndUpdateEmailAddress":true,"EMAILOTPAuthenticationEndpointURL":"/emailotpauthenticationendpoint/emailotp.jsp","EMAILOTPMandatory":false,"EmailAddressRequestPage":"/emailotpauthenticationendpoint/emailAddress.jsp","EmailOTPAuthenticationEndpointErrorPage":"/emailotpauthenticationendpoint/emailotpError.jsp","EmailOTPEnableByUserClaim":true,"federatedEmailAttributeKey":"email","secondaryUserstore":"primary","sendOTPToFederatedEmailAttribute":false,"showEmailAddressInUI":true,"useEventHandlerBasedEmailSender":true,"usecase":"local"}
	Parameters EmailOtpParameters `toml:"parameters,omitempty" json:"parameters,omitempty"`
}

/* SMS OTP parameters */
type SmsOtpParameters struct {
	// +kubebuilder:default:="/smsotpauthenticationendpoint/smsotp.jsp"
	SMSOTPAuthenticationEndpointURL string `toml:"SMSOTPAuthenticationEndpointURL,omitempty" json:"SMSOTPAuthenticationEndpointURL,omitempty"`
	// +kubebuilder:default:="/smsotpauthenticationendpoint/smsotpError.jsp"
	SMSOTPAuthenticationEndpointErrorPage string `toml:"SMSOTPAuthenticationEndpointErrorPage,omitempty" json:"SMSOTPAuthenticationEndpointErrorPage,omitempty"`
	// +kubebuilder:default:="/smsotpauthenticationendpoint/mobile.jsp"
	MobileNumberRegPage string `toml:"MobileNumberRegPage,omitempty" json:"MobileNumberRegPage,omitempty"`
	// +kubebuilder:default:=true
	RetryEnable bool `toml:"RetryEnable,omitempty" json:"RetryEnable,omitempty"`
	// +kubebuilder:default:=true
	ResendEnable bool `toml:"ResendEnable,omitempty" json:"ResendEnable,omitempty"`
	// +kubebuilder:default:=true
	BackupCode bool `toml:"BackupCode,omitempty" json:"BackupCode,omitempty"`
	// +kubebuilder:default:=true
	SMSOTPEnableByUserClaim bool `toml:"SMSOTPEnableByUserClaim,omitempty" json:"SMSOTPEnableByUserClaim,omitempty"`
	// +kubebuilder:default:=false
	SMSOTPMandatory bool `toml:"SMSOTPMandatory,omitempty" json:"SMSOTPMandatory,omitempty"`
	// +kubebuilder:default:=true
	CaptureAndUpdateMobileNumber bool `toml:"CaptureAndUpdateMobileNumber,omitempty" json:"CaptureAndUpdateMobileNumber,omitempty"`
	// +kubebuilder:default:=false
	SendOTPDirectlyToMobile bool `toml:"SendOTPDirectlyToMobile,omitempty" json:"SendOTPDirectlyToMobile,omitempty"`
	// +kubebuilder:default:=false
	RedirectToMultiOptionPageOnFailure bool `toml:"redirectToMultiOptionPageOnFailure,omitempty" json:"redirectToMultiOptionPageOnFailure,omitempty"`
}

/* SMS OTP configurations */
type SmsOtp struct {
	// +kubebuilder:default:="SmsOTP"
	Name string `toml:"name,omitempty" json:"name,omitempty"`
	// +kubebuilder:default:={"BackupCode":true,"CaptureAndUpdateMobileNumber":true,"MobileNumberRegPage":"/smsotpauthenticationendpoint/mobile.jsp","ResendEnable":true,"RetryEnable":true,"SMSOTPAuthenticationEndpointErrorPage":"/smsotpauthenticationendpoint/smsotpError.jsp","SMSOTPAuthenticationEndpointURL":"/smsotpauthenticationendpoint/smsotp.jsp","SMSOTPEnableByUserClaim":true,"SMSOTPMandatory":false,"SendOTPDirectlyToMobile":false,"redirectToMultiOptionPageOnFailure":false}
	Parameters SmsOtpParameters `toml:"parameters,omitempty" json:"parameters,omitempty"`
	// +kubebuilder:default:=false
	Enable bool `toml:"enable,omitempty" json:"enable,omitempty"`
}

/* Multi factor authenticators */
type Authenticator struct {
	Totp     Totp     `toml:"totp,omitempty" json:"totp,omitempty"`
	EmailOtp EmailOtp `toml:"email_otp,omitempty" json:"email_otp,omitempty"`
	SmsOtp   SmsOtp   `toml:"sms_otp,omitempty" json:"sms_otp,omitempty"`
}

/* Second step authentication configurations */
type StepAuthentication struct {
	Authenticator Authenticator `toml:"authenticator,omitempty" json:"authenticator,omitempty"`
}

/* Recaptcha configurations */
type Recaptcha struct {
	// +kubebuilder:default:=false
	Enabled bool `toml:"enabled,omitempty" json:"enabled,omitempty"`
	// +kubebuilder:default:="https://www.google.com/recaptcha/api.js"
	APIURL string `toml:"api_url,omitempty" json:"api_url,omitempty"`
	// +kubebuilder:default:="https://www.google.com/recaptcha/api/siteverify"
	VerifyURL string `toml:"verify_url,omitempty" json:"verify_url,omitempty"`
	// +kubebuilder:validation:Required
	SiteKey string `toml:"site_key,omitempty" json:"site_key,omitempty"`
	// +kubebuilder:validation:Required
	SecretKey string `toml:"secret_key,omitempty" json:"secret_key,omitempty"`
}

/* SMTP email configurations */
type Email struct {
	// +kubebuilder:validation:Required
	FromAddress string `toml:"from_address,omitempty" json:"from_address,omitempty"`
	// +kubebuilder:validation:Required
	Username string `toml:"username,omitempty" json:"username,omitempty"`
	// +kubebuilder:validation:Required
	Password string `toml:"password,omitempty" json:"password,omitempty"`
	// +kubebuilder:validation:Required
	Hostname string `toml:"hostname,omitempty" json:"hostname,omitempty"`
	// +kubebuilder:default:="587"
	Port string `toml:"port,omitempty" json:"port,omitempty"`
	// +kubebuilder:default:=true
	EnableStartTLS bool `toml:"enable_start_tls,omitempty" json:"enable_start_tls,omitempty"`
	// +kubebuilder:default:=true
	EnableAuthentication bool `toml:"enable_authentication,omitempty" json:"enable_authentication,omitempty"`
}
type OutputAdapter struct {
	Email Email `toml:"email,omitempty" json:"email,omitempty"`
}
type TenantMgt struct {
	// +kubebuilder:default:=false
	EnableEmailDomain bool `toml:"enable_email_domain,omitempty" json:"enable_email_domain,omitempty"`
}
type TenantCtx struct {
	EnableQualifiedUrls bool `toml:"enable_tenant_qualified_urls,omitempty" json:"enableTenantQualifiedUrls,omitempty"`
}
type Wsdl struct {
	// +kubebuilder:default:=false
	Enable bool `toml:"enable,omitempty" json:"enable,omitempty"`
}
type AdminService struct {
	Wsdl Wsdl `toml:"wsdl,omitempty" json:"wsdl,omitempty"`
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
