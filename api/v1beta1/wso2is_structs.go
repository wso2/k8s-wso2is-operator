package v1beta1

type Configurations struct {
	// +kubebuilder:default:="NodePort"
	ServiceType string `json:"serviceType,omitempty"`

	Server Server `json:"server,omitempty" toml:"server,omitempty"`

	// +kubebuilder:default:={ "username" : "admin", "password": "admin", "createAdminAccount": true }
	SuperAdmin SuperAdmin `json:"superAdmin,omitempty" toml:"super_admin,omitempty"`

	// +kubebuilder:default:={ "type" : "database_unique_id"}
	UserStore UserStore `json:"userStore,omitempty" toml:"user_store,omitempty"`

	// +kubebuilder:default:={   "userDb": {     "type": "h2",     "url": "jdbc:h2:./repository/database/UM_DB;DB_CLOSE_ON_EXIT=false;LOCK_TIMEOUT=60000",     "username": "wso2carbon",     "password": "wso2carbon"   },   "identityDb": {     "type": "h2",     "url": "jdbc:h2:./repository/database/WSO2IS_IDENTITY_DB;DB_CLOSE_ON_EXIT=false;LOCK_TIMEOUT=60000",     "username": "wso2carbon",     "password": "wso2carbon"   },   "sharedDb": {     "type": "h2",     "url": "jdbc:h2:./repository/database/WSO2IS_SHARED_DB;DB_CLOSE_ON_EXIT=false;LOCK_TIMEOUT=60000",     "username": "wso2carbon",     "password": "wso2carbon"   } }
	Database Database `json:"database,omitempty" toml:"database,omitempty"`

	// +kubebuilder:default:={ "https" : { "properties" : { "proxyPort" : 443 } } }
	Transport Transport `json:"transport,omitempty" toml:"transport,omitempty"`

	// +kubebuilder:default:={ "primary":{"name":"wso2carbon.jks","password":"wso2carbon"} }
	KeyStore KeyStore `json:"keystore,omitempty" toml:"keystore,omitempty"`

	TrustStore TrustStore `json:"truststore,omitempty" toml:"truststore,omitempty"`

	Identity Identity `json:"identity,omitempty" toml:"identity,omitempty"`

	AccountRecovery AccountRecovery `json:"accountRecovery,omitempty" toml:"account_recovery,omitempty"`

	Monitoring Monitoring `json:"monitoring,omitempty" toml:"monitoring,omitempty"`

	// +kubebuilder:default:={"shutdownhook": {"enabled": false},"logging": {"type": "log4j2"}}
	Hazelcast Hazelcast `json:"hazelcast,omitempty" toml:"hazelcast,omitempty"`

	Authentication Authentication `json:"authentication,omitempty" toml:"authentication,omitempty"`

	Recaptcha Recaptcha `json:"recaptcha,omitempty" toml:"recaptcha,omitempty"`

	OutputAdapter OutputAdapter `json:"outputAdapter,omitempty" toml:"output_adapter,omitempty"`

	// +kubebuilder:default:={"membershipScheme":"kubernetes","properties":{"membershipSchemeClassName":"org.wso2.carbon.membership.scheme.kubernetes.KubernetesMembershipScheme","kubernetesNamespace":"wso2-iam-system","kubernetesServices":"wso2is-service","kubernetesMasterSkipSslVerification":true,"useDns":false} }
	Clustering Clustering `json:"clustering,omitempty" toml:"clustering,omitempty"`

	TenantMgt TenantMgt `json:"tenantMgt,omitempty" toml:"tenant_mgt,omitempty"`

	// +kubebuilder:default:={ "enableTenantQualifiedUrls":true }
	TenantContext TenantContext `json:"tenantContext,omitempty" toml:"tenant_context,omitempty"`

	AdminService AdminService `json:"adminService,omitempty" toml:"admin_service,omitempty"`

	SystemRoles SystemRoles `json:"systemRoles,omitempty" toml:"system_roles,omitempty"`

	CarbonHealthCheck CarbonHealthCheck `json:"carbonHealthCheck,omitempty" toml:"carbon_health_check,omitempty"`

	SystemApplications SystemApplications `json:"systemApplications,omitempty" toml:"system_applications,omitempty"`

	Catalina Catalina `json:"catalina,omitempty" toml:"catalina,omitempty"`
}

type Server struct {
	// +kubebuilder:default:= "wso2is.com"
	HostName string `json:"hostname,omitempty" toml:"hostname,omitempty"`

	// +kubebuilder:default:="$env{NODE_IP}"
	NodeIp string `json:"nodeIp,omitempty" toml:"node_ip,omitempty"`

	// +kubebuilder:default:="https://$ref{server.hostname}"
	BasePath string `json:"basePath,omitempty" toml:"base_path,omitempty"`
}

type SuperAdmin struct {
	Username string `json:"username,omitempty" toml:"username,omitempty"`

	Password string `json:"password,omitempty" toml:"password,omitempty"`

	AdminRole string `json:"adminRole,omitempty" toml:"admin_role,omitempty"`

	CreateAdminAccount bool `json:"createAdminAccount,omitempty" toml:"create_admin_account,omitempty"`
}

type UserStore struct {
	Type string `json:"type,omitempty" toml:"type,omitempty"`

	UserNameJavaRegex string `json:"usernameJavaRegex,omitempty" toml:"username_java_regex,omitempty"`

	ConnectionUrl string `json:"connectionUrl,omitempty" toml:"connection_url,omitempty"`

	ConnectionName string `json:"connectionName,omitempty" toml:"connection_name,omitempty"`

	ConnectionPassword string `json:"connectionPassword,omitempty" toml:"connection_password,omitempty"`

	Properties UserStoreProperties `json:"properties,omitempty" toml:"properties,omitempty"`

	BaseDn string `json:"baseDn,omitempty" toml:"base_dn,omitempty"`

	UsernameAttribute string `json:"usernameAttribute,omitempty" toml:"user_name_attribute,omitempty"`
}

type Database struct {
	UserDb UserDb `json:"userDb,omitempty" toml:"user,omitempty"`

	IdentityDb IdentityDb `json:"identityDb,omitempty" toml:"identity_db,omitempty"`

	SharedDb SharedDb `json:"sharedDb,omitempty" toml:"shared_db,omitempty"`
}

type Transport struct {
	Https Https `json:"https,omitempty" toml:"https,omitempty"`
}

type KeyStore struct {
	Primary Primary `json:"primary,omitempty" toml:"primary,omitempty"`

	Internal Internal `json:"internal,omitempty" toml:"internal,omitempty"`

	Tls Tls `json:"tls,omitempty" toml:"tls,omitempty"`
}

type TrustStore struct {
	FileName string `json:"fileName,omitempty" toml:"file_name,omitempty"`

	Password string `json:"password,omitempty" toml:"password,omitempty"`

	Type string `json:"type,omitempty" toml:"type,omitempty"`
}

type Identity struct {
	AuthFramework AuthFramework `json:"authFramework,omitempty" toml:"auth_framework,omitempty"`
}

type AccountRecovery struct {
	Endpoint AccoutnRecoveryEndpoint `json:"endpoint,omitempty" toml:"endpoint,omitempty"`
}

type Monitoring struct {
	Jmx Jmx `json:"jmx,omitempty" toml:"jmx,omitempty"`
}

type Hazelcast struct {
	HazelcastShutdownHook HazelcastShutdownHook `json:"shutdownhook,omitempty" toml:"shutdownhook,omitempty"`

	Logging Logging `json:"logging,omitempty" toml:"logging,omitempty"`
}

type Authentication struct {
	Consent Consent `json:"consent,omitempty" toml:"consent,omitempty"`

	Authenticator Authenticator `json:"authenticator,omitempty" toml:"authenticator,omitempty"`

	Endpoint AuthEndpoint `json:"endpoint,omitempty" toml:"endpoint,omitempty"`

	Adaptive Adaptive `json:"adaptive,omitempty" toml:"adaptive,omitempty"`

	CustomAuthenticator string `json:"customAuthenticator,omitempty" toml:"custom_authenticator,omitempty"`

	JitProvisioning JitProvisioning `json:"jitProvisioning,omitempty" toml:"jit_provisioning,omitempty"`
}

type Recaptcha struct {
	// +kubebuilder:default:=false
	Enabled bool `json:"enabled,omitempty" toml:"enabled,omitempty"`

	// +kubebuilder:default:="https://www.google.com/recaptcha/api.js"
	ApiUrl string `json:"apiUrl,omitempty" toml:"api_url,omitempty"`

	// +kubebuilder:default:="https://www.google.com/recaptcha/api/siteverify"
	VerifyUrl string `json:"verifyUrl,omitempty" toml:"verify_url,omitempty"`

	RequestWrapUrls string `json:"requestWrapUrls,omitempty" toml:"request_wrap_urls,omitempty"`

	SiteKey string `json:"siteKey,omitempty" toml:"site_key,omitempty"`

	SecretKey string `json:"secretKey,omitempty" toml:"secret_key,omitempty"`
}

type OutputAdapter struct {
	Email Email `json:"email,omitempty" toml:"email,omitempty"`
}

type Clustering struct {
	MembershipScheme string `json:"membershipScheme,omitempty" toml:"membership_scheme,omitempty"`

	Properties ClusteringProperties `json:"properties,omitempty" toml:"properties,omitempty"`
}

type TenantMgt struct {
	// +kubebuilder:default:=false
	EnableEmailDomain bool `json:"enableEmailDomain,omitempty" toml:"enable_email_domain,omitempty"`
}

type TenantContext struct {
	EnableTenantQualifiedUrls bool `json:"enableTenantQualifiedUrls,omitempty" toml:"enable_tenant_qualified_urls,omitempty"`

	EnableTenantedSessions bool `json:"enableTenantedSessions,omitempty" toml:"enable_tenanted_sessions,omitempty"`

	Rewrite Rewrite `json:"rewrite,omitempty" toml:"rewrite,omitempty"`
}

type AdminService struct {
	Wsdl Wsdl `json:"wsdl,omitempty" toml:"wsdl,omitempty"`
}

type SystemRoles struct {
	Enable bool `json:"enable,omitempty" toml:"enable,omitempty"`

	ReadOnlyRoles bool `json:"readOnlyRoles,omitempty" toml:"read_only_roles,omitempty"`
}

type CarbonHealthCheck struct {
	Enable bool `json:"enable,omitempty" toml:"enable,omitempty"`

	HealthChecker HealthChecker `json:"healthChecker,omitempty" toml:"health_checker,omitempty"`
}

type SystemApplications struct {
	ReadOnlyRoles bool `json:"readOnlyRoles,omitempty" toml:"read_only_roles,omitempty"`

	FidpRoleBasedAuthzEnabledApps string `json:"fidpRoleBasedAuthzEnabledApps,omitempty" toml:"fidp_role_based_authz_enabled_apps,omitempty"`
}

type Catalina struct {
	Valves Valves `json:"valves,omitempty" toml:"valves,omitempty"`
}

type UserStoreProperties struct {
	CaseInsensitiveUsername bool `json:"caseInsensitiveUsername,omitempty" toml:"CaseInsensitiveUsername,omitempty"`
}

type UserDb struct {
	Type string `json:"type,omitempty" toml:"type,omitempty"`

	Url string `json:"url,omitempty" toml:"url,omitempty"`

	Hostname string `json:"hostname,omitempty" toml:"hostname,omitempty"`

	Username string `json:"username,omitempty" toml:"username,omitempty"`

	Password string `json:"password,omitempty" toml:"password,omitempty"`

	Driver string `json:"driver,omitempty" toml:"driver,omitempty"`

	PoolOptions UserPoolOptions `json:"poolOptions,omitempty" toml:"pool_options,omitempty"`
}

type IdentityDb struct {
	Type string `json:"type,omitempty" toml:"type,omitempty"`

	Url string `json:"url,omitempty" toml:"url,omitempty"`

	Hostname string `json:"hostname,omitempty" toml:"hostname,omitempty"`

	Username string `json:"username,omitempty" toml:"username,omitempty"`

	Password string `json:"password,omitempty" toml:"password,omitempty"`

	Driver string `json:"driver,omitempty" toml:"driver,omitempty"`

	PoolOptions IdentityDbPoolOptions `json:"poolOptions,omitempty" toml:"pool_options,omitempty"`
}

type SharedDb struct {
	Type string `json:"type,omitempty" toml:"type,omitempty"`

	Url string `json:"url,omitempty" toml:"url,omitempty"`

	Hostname string `json:"hostname,omitempty" toml:"hostname,omitempty"`

	Username string `json:"username,omitempty" toml:"username,omitempty"`

	Password string `json:"password,omitempty" toml:"password,omitempty"`

	Driver string `json:"driver,omitempty" toml:"driver,omitempty"`

	PoolOptions SharedDbPoolOptions `json:"poolOptions,omitempty" toml:"pool_options,omitempty"`
}

type Https struct {
	Properties HttpsProperties `json:"properties,omitempty" toml:"properties,omitempty"`
}

type Primary struct {
	Name string `json:"name,omitempty" toml:"name,omitempty"`

	Password string `json:"password,omitempty" toml:"password,omitempty"`
}

type Internal struct {
	FileName string `json:"fileName,omitempty" toml:"file_name,omitempty"`

	Type string `json:"type,omitempty" toml:"type,omitempty"`

	Password string `json:"password,omitempty" toml:"password,omitempty"`

	Alias string `json:"alias,omitempty" toml:"alias,omitempty"`

	KeyPassword string `json:"keyPassword,omitempty" toml:"key_password,omitempty"`
}

type Tls struct {
	FileName string `json:"fileName,omitempty" toml:"file_name,omitempty"`

	Type string `json:"type,omitempty" toml:"type,omitempty"`

	Password string `json:"password,omitempty" toml:"password,omitempty"`

	Alias string `json:"alias,omitempty" toml:"alias,omitempty"`

	KeyPassword string `json:"keyPassword,omitempty" toml:"key_password,omitempty"`
}

type AuthFramework struct {
	Endpoint AuthFrameworkEndpoint `json:"endpoint,omitempty" toml:"endpoint,omitempty"`
}

type AccoutnRecoveryEndpoint struct {
	Auth EndpointAuth `json:"auth,omitempty" toml:"auth,omitempty"`
}

type Jmx struct {
	// +kubebuilder:default:=true
	RmiServerStart bool `json:"rmiServerStart,omitempty" toml:"rmi_server_start,omitempty"`
}

type HazelcastShutdownHook struct {
	Enabled bool `json:"enabled,omitempty" toml:"enabled,omitempty"`
}

type Logging struct {
	Type string `json:"type,omitempty" toml:"type,omitempty"`
}

type Consent struct {
	DataSource string `json:"dataSource,omitempty" toml:"data_source,omitempty"`
}

type Authenticator struct {
	Basic Basic `json:"basic,omitempty" toml:"basic,omitempty"`

	Totp Totp `json:"totp,omitempty" toml:"totp,omitempty"`

	EmailOtp EmailOtp `json:"emailOtp,omitempty" toml:"email_otp,omitempty"`

	SmsOtp SmsOtp `json:"smsOtp,omitempty" toml:"sms_otp,omitempty"`

	MagicLink MagicLink `json:"magiclink,omitempty" toml:"magiclink,omitempty"`

	Fido Fido `json:"fido,omitempty" toml:"fido,omitempty"`
}

type AuthEndpoint struct {
	EnableCustomClaimMappings bool `json:"enableCustomClaimMappings,omitempty" toml:"enableCustomClaimMappings,omitempty"`

	EnableMergingCustomClaimMappingsWithDefault bool `json:"enableMergingCustomClaimMappingsWithDefault,omitempty" toml:"enableMergingCustomClaimMappingsWithDefault,omitempty"`
}

type Adaptive struct {
	AllowLoops bool `json:"allowLoops,omitempty" toml:"allow_loops,omitempty"`

	ExecutionSupervisor ExecutionSupervisor `json:"executionSupervisor,omitempty" toml:"execution_supervisor,omitempty"`
}

type JitProvisioning struct {
	EnableEnhancedFeature bool `json:"enableEnhancedFeature,omitempty" toml:"enable_enhanced_feature,omitempty"`
}

type Email struct {
	FromAddress string `json:"fromAddress,omitempty" toml:"from_address,omitempty"`

	Username string `json:"username,omitempty" toml:"username,omitempty"`

	Password string `json:"password,omitempty" toml:"password,omitempty"`

	Hostname string `json:"hostname,omitempty" toml:"hostname,omitempty"`

	Port string `json:"port,omitempty" toml:"port,omitempty"`

	EnableStartTls string `json:"enableStartTls,omitempty" toml:"enable_start_tls,omitempty"`

	EnableAuthentication string `json:"enableAuthentication,omitempty" toml:"enable_authentication,omitempty"`
}

type ClusteringProperties struct {
	MembershipSchemeClassName string `json:"membershipSchemeClassName,omitempty" toml:"membershipSchemeClassName,omitempty"`

	KubernetesNamespace string `json:"kubernetesNamespace,omitempty" toml:"KUBERNETES_NAMESPACE,omitempty"`

	KubernetesServices string `json:"kubernetesServices,omitempty" toml:"KUBERNETES_SERVICES,omitempty"`

	KubernetesMasterSkipSslVerification bool `json:"kubernetesMasterSkipSslVerification,omitempty" toml:"KUBERNETES_MASTER_SKIP_SSL_VERIFICATION,omitempty"`

	UseDns bool `json:"useDns,omitempty" toml:"USE_DNS,omitempty"`
}

type Rewrite struct {
	CustomWebapps bool `json:"customWebapps,omitempty" toml:"custom_webapps,omitempty"`
}

type Wsdl struct {
	// +kubebuilder:default:=false
	Enable bool `json:"enable,omitempty" toml:"enable,omitempty"`
}

type HealthChecker struct {
	DatasourceHealthChecker DatasourceHealthChecker `json:"dataSourceHealthChecker,omitempty" toml:"data_source_health_checker,omitempty"`
}

type Valves struct {
	Valve Valve `json:"valve,omitempty" toml:"valve,omitempty"`
}

type UserPoolOptions struct {
	ValidationQuery string `json:"validationQuery,omitempty" toml:"validationQuery,omitempty"`

	MaxActive string `json:"maxActive,omitempty" toml:"maxActive,omitempty"`
}

type IdentityDbPoolOptions struct {
	ValidationQuery string `json:"validationQuery,omitempty" toml:"validationQuery,omitempty"`

	MaxActive string `json:"maxActive,omitempty" toml:"maxActive,omitempty"`

	MaxWait string `json:"maxWait,omitempty" toml:"maxWait,omitempty"`

	MinIdle string `json:"minIdle,omitempty" toml:"minIdle,omitempty"`

	TestOnBorrow string `json:"testOnBorrow,omitempty" toml:"testOnBorrow,omitempty"`

	ValidationInterval string `json:"validationInterval,omitempty" toml:"validationInterval,omitempty"`

	DefaultAutoCommit string `json:"defaultAutoCommit,omitempty" toml:"defaultAutoCommit,omitempty"`

	CommitOnReturn string `json:"commitOnReturn,omitempty" toml:"commitOnReturn,omitempty"`
}

type SharedDbPoolOptions struct {
	// +kubebuilder:default:="SELECT 1"
	ValidationQuery string `json:"validationQuery,omitempty" toml:"validationQuery,omitempty"`

	MaxActive string `json:"maxActive,omitempty" toml:"maxActive,omitempty"`
}

type HttpsProperties struct {
	ProxyPort int32 `json:"proxyPort,omitempty" toml:"proxyPort,omitempty"`
}

type AuthFrameworkEndpoint struct {
	AppPassword string `json:"appPassword,omitempty" toml:"app_password,omitempty"`
}

type EndpointAuth struct {
	Hash string `json:"hash,omitempty" toml:"hash,omitempty"`
}

type Basic struct {
	Parameters BasicParameters `json:"parameters,omitempty" toml:"parameters,omitempty"`
}

type Totp struct {
	// +kubebuilder:default:={"Issuer":"WSO2","totpAuthenticationEndpointEnableTotpPage":"/totpauthenticationendpoint/enableTOTP.jsp","totpAuthenticationEndpointErrorPage":"/totpauthenticationendpoint/totpError.jsp","totpAuthenticationEndpointUrl":"/totpauthenticationendpoint/totp.jsp","useCommonIssuer":true,"authenticationMandatory":true,"encodingMethod":"Base32","enrolUserInAuthenticationFlow":true,"secondaryUserstore":"primary","timeStepSize":"30","usecase":"local","windowSize":"3"}
	Parameters TotpParameters `json:"parameters,omitempty" toml:"parameters,omitempty"`

	// +kubebuilder:default:=false
	Enable bool `json:"enable,omitempty" toml:"enable,omitempty"`
}

type EmailOtp struct {
	// +kubebuilder:default:="EmailOTP"
	Name string `json:"name,omitempty" toml:"name,omitempty"`

	// +kubebuilder:default:=false
	Enable bool `json:"enable,omitempty" toml:"enable,omitempty"`

	// +kubebuilder:default:={"captureAndUpdateEmailAddress":true,"emailOtpAuthenticationEndpointUrl":"/emailotpauthenticationendpoint/emailotp.jsp","emailOtpMandatory":false,"emailAddressRequestPage":"/emailotpauthenticationendpoint/emailAddress.jsp","emailOtpAuthenticationEndpointErrorPage":"/emailotpauthenticationendpoint/emailotpError.jsp","emailOtpEnableByUserClaim":true,"federatedEmailAttributeKey":"email","secondaryUserstore":"primary","sendOtpToFederatedEmailAttribute":false,"showEmailAddressInUi":true,"useEventHandlerBasedEmailSender":true,"usecase":"local"}
	Parameters EmailOtpParameters `json:"parameters,omitempty" toml:"parameters,omitempty"`
}

type SmsOtp struct {
	// +kubebuilder:default:="SmsOTP"
	Name string `json:"name,omitempty" toml:"name,omitempty"`

	// +kubebuilder:default:={"backupCode":true,"captureAndUpdateMobileNumber":true,"mobileNumberRegPage":"/smsotpauthenticationendpoint/mobile.jsp","resendEnable":true,"retryEnable":true,"smsOtpAuthenticationEndpointErrorPage":"/smsotpauthenticationendpoint/smsotpError.jsp","smsOtpAuthenticationEndpointUrl":"/smsotpauthenticationendpoint/smsotp.jsp","smsOtpEnableByUserClaim":true,"smsOtpMandatory":false,"sendOtpDirectlyToMobile":false,"redirectToMultiOptionPageOnFailure":false}
	Parameters SmsParameters `json:"parameters,omitempty" toml:"parameters,omitempty"`

	// +kubebuilder:default:=false
	Enable bool `json:"enable,omitempty" toml:"enable,omitempty"`
}

type MagicLink struct {
	Name string `json:"name,omitempty" toml:"name,omitempty"`

	// +kubebuilder:default:=false
	Enable bool `json:"enable,omitempty" toml:"enable,omitempty"`

	Parameters MagicLinkParameters `json:"parameters,omitempty" toml:"parameters,omitempty"`
}

type Fido struct {
	Parameters FidoParameters `json:"parameters,omitempty" toml:"parameters,omitempty"`
}

type ExecutionSupervisor struct {
	Timeout int32 `json:"timeout,omitempty" toml:"timeout,omitempty"`
}

type DatasourceHealthChecker struct {
	Properties DataSourceHealthCheckerProperties `json:"properties,omitempty" toml:"properties,omitempty"`
}

type Valve struct {
	Properties ValveProperties `json:"properties,omitempty" toml:"properties,omitempty"`
}

type BasicParameters struct {
	ShowAuthFailureReason bool `json:"showAuthFailureReason,omitempty" toml:"showAuthFailureReason,omitempty"`

	ShowAuthFailureReasonOnLoginPage bool `json:"showAuthFailureReasonOnLoginPage,omitempty" toml:"showAuthFailureReasonOnLoginPage,omitempty"`
}

type TotpParameters struct {
	ShowAuthFailureReason bool `json:"showAuthFailureReason,omitempty" toml:"showAuthFailureReason,omitempty"`

	ShowAuthFailureReasonOnLoginPage bool `json:"showAuthFailureReasonOnLoginPage,omitempty" toml:"showAuthFailureReasonOnLoginPage,omitempty"`

	EncodingMethod string `json:"encodingMethod,omitempty" toml:"encodingMethod,omitempty"`

	TimeStepSize string `json:"timeStepSize,omitempty" toml:"timeStepSize,omitempty"`

	WindowSize string `json:"windowSize,omitempty" toml:"windowSize,omitempty"`

	AuthenticationMandatory bool `json:"authenticationMandatory,omitempty" toml:"authenticationMandatory,omitempty"`

	EnrolUserInAuthenticationFlow bool `json:"enrolUserInAuthenticationFlow,omitempty" toml:"enrolUserInAuthenticationFlow,omitempty"`

	UseCase string `json:"usecase,omitempty" toml:"usecase,omitempty"`

	SecondaryUserstore string `json:"secondaryUserstore,omitempty" toml:"secondaryUserstore,omitempty"`

	TotpAuthenticationEndpointUrl string `json:"totpAuthenticationEndpointUrl,omitempty" toml:"TOTPAuthenticationEndpointURL,omitempty"`

	TotpAuthenticationEndpointErrorPage string `json:"totpAuthenticationEndpointErrorPage,omitempty" toml:"TOTPAuthenticationEndpointErrorPage,omitempty"`

	TotpAuthenticationEndpointEnableTotpPage string `json:"totpAuthenticationEndpointEnableTotpPage,omitempty" toml:"TOTPAuthenticationEndpointEnableTOTPPage,omitempty"`

	Issuer string `json:"Issuer,omitempty" toml:"Issuer,omitempty"`

	UseCommonIssuer bool `json:"useCommonIssuer,omitempty" toml:"UseCommonIssuer,omitempty"`
}

type EmailOtpParameters struct {
	ShowAuthFailureReason bool `json:"showAuthFailureReason,omitempty" toml:"showAuthFailureReason,omitempty"`

	ShowAuthFailureReasonOnLoginPage bool `json:"showAuthFailureReasonOnLoginPage,omitempty" toml:"showAuthFailureReasonOnLoginPage,omitempty"`

	EmailOtpAuthenticationEndpointUrl string `json:"emailOtpAuthenticationEndpointUrl,omitempty" toml:"EMAILOTPAuthenticationEndpointURL,omitempty"`

	EmailOtpAuthenticationEndpointErrorPage string `json:"emailOtpAuthenticationEndpointErrorPage,omitempty" toml:"EmailOTPAuthenticationEndpointErrorPage,omitempty"`

	EmailAddressRequestPage string `json:"emailAddressRequestPage,omitempty" toml:"EmailAddressRequestPage,omitempty"`

	UseCase string `json:"usecase,omitempty" toml:"usecase,omitempty"`

	SecondaryUserstore string `json:"secondaryUserstore,omitempty" toml:"secondaryUserstore,omitempty"`

	EmailOtpMandatory bool `json:"emailOtpMandatory,omitempty" toml:"EMAILOTPMandatory,omitempty"`

	SendOtpToFederatedEmailAttribute bool `json:"sendOtpToFederatedEmailAttribute,omitempty" toml:"sendOTPToFederatedEmailAttribute,omitempty"`

	FederatedEmailAttributeKey string `json:"federatedEmailAttributeKey,omitempty" toml:"federatedEmailAttributeKey,omitempty"`

	EmailOtpEnableByUserClaim bool `json:"emailOtpEnableByUserClaim,omitempty" toml:"EmailOTPEnableByUserClaim,omitempty"`

	CaptureAndUpdateEmailAddress bool `json:"captureAndUpdateEmailAddress,omitempty" toml:"CaptureAndUpdateEmailAddress,omitempty"`

	ShowEmailAddressInUi bool `json:"showEmailAddressInUi,omitempty" toml:"showEmailAddressInUI,omitempty"`

	UseEventHandlerBasedEmailSender bool `json:"useEventHandlerBasedEmailSender,omitempty" toml:"useEventHandlerBasedEmailSender,omitempty"`
}

type SmsParameters struct {
	SmsOtpAuthenticationEndpointUrl string `json:"smsOtpAuthenticationEndpointUrl,omitempty" toml:"SMSOTPAuthenticationEndpointURL,omitempty"`

	SmsOtpAuthenticationEndpointErrorPage string `json:"smsOtpAuthenticationEndpointErrorPage,omitempty" toml:"SMSOTPAuthenticationEndpointErrorPage,omitempty"`

	MobileNumberRegPage string `json:"mobileNumberRegPage,omitempty" toml:"MobileNumberRegPage,omitempty"`

	RetryEnable bool `json:"retryEnable,omitempty" toml:"RetryEnable,omitempty"`

	ResendEnable bool `json:"resendEnable,omitempty" toml:"ResendEnable,omitempty"`

	BackupCode bool `json:"backupCode,omitempty" toml:"BackupCode,omitempty"`

	SmsOtpEnabledByUserClaim bool `json:"smsOtpEnableByUserClaim,omitempty" toml:"SMSOTPEnableByUserClaim,omitempty"`

	SmsOtpMandatory bool `json:"smsOtpMandatory,omitempty" toml:"SMSOTPMandatory,omitempty"`

	CaptureAndUpdateMobileNumber bool `json:"captureAndUpdateMobileNumber,omitempty" toml:"CaptureAndUpdateMobileNumber,omitempty"`

	SendOtpDirectlyToMobile bool `json:"sendOtpDirectlyToMobile,omitempty" toml:"SendOTPDirectlyToMobile,omitempty"`

	RedirectToMultiOptionPageOnFailure bool `json:"redirectToMultiOptionPageOnFailure,omitempty" toml:"redirectToMultiOptionPageOnFailure,omitempty"`
}

type MagicLinkParameters struct {
	ExpiryTime string `json:"expiryTime,omitempty" toml:"ExpiryTime,omitempty"`

	Tags string `json:"tags,omitempty" toml:"Tags,omitempty"`

	BlockedUserStoreDomains string `json:"blockedUserStoreDomains,omitempty" toml:"BlockedUserStoreDomains,omitempty"`
}

type FidoParameters struct {
	Tags string `json:"tags,omitempty" toml:"Tags,omitempty"`
}

type DataSourceHealthCheckerProperties struct {
	Monitored Monitored `json:"monitored,omitempty" toml:"monitored,omitempty"`
}

type ValveProperties struct {
	ClassName string `json:"className,omitempty" toml:"className,omitempty"`
}

type Monitored struct {
	Datasources string `json:"datasources,omitempty" toml:"datasources,omitempty"`
}
