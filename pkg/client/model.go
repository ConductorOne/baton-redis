package client

import "time"

type User struct {
	AuthMethod             string        `json:"auth_method"`
	BdbsEmailAlerts        []interface{} `json:"bdbs_email_alerts"`
	CertificateSubjectLine string        `json:"certificate_subject_line"`
	ClusterEmailAlerts     bool          `json:"cluster_email_alerts"`
	Email                  string        `json:"email"`
	EmailAlerts            bool          `json:"email_alerts"`
	Name                   string        `json:"name"`
	PasswordIssueDate      time.Time     `json:"password_issue_date"`
	Role                   string        `json:"role"`
	RoleUIDs               []int         `json:"role_uids"`
	Status                 string        `json:"status"`
	UID                    int           `json:"uid"`
}

type Setting struct {
	Enabled   bool   `json:"enabled"`
	Threshold string `json:"threshold"`
}

type Cluster struct {
	AlertSettings struct {
		ClusterCaCertAboutToExpire     Setting `json:"cluster_ca_cert_about_to_expire"`
		ClusterCertsAboutToExpire      Setting `json:"cluster_certs_about_to_expire"`
		ClusterLicenseAboutToExpire    Setting `json:"cluster_license_about_to_expire"`
		ClusterNodeOperationFailed     bool    `json:"cluster_node_operation_failed"`
		ClusterOcspQueryFailed         bool    `json:"cluster_ocsp_query_failed"`
		ClusterOcspStatusRevoked       bool    `json:"cluster_ocsp_status_revoked"`
		NodeChecksError                bool    `json:"node_checks_error"`
		NodeEphemeralStorage           Setting `json:"node_ephemeral_storage"`
		NodeFreeFlash                  Setting `json:"node_free_flash"`
		NodeInternalCertsAboutToExpire Setting `json:"node_internal_certs_about_to_expire"`
		NodePersistentStorage          Setting `json:"node_persistent_storage"`
	} `json:"alert_settings"`
	BigstoreDriver                      string        `json:"bigstore_driver"`
	BlockClusterChanges                 bool          `json:"block_cluster_changes"`
	CcsInternodeEncryption              bool          `json:"ccs_internode_encryption"`
	ClusterSSHPublicKey                 string        `json:"cluster_ssh_public_key"`
	CmPort                              int           `json:"cm_port"`
	CmServerVersion                     int           `json:"cm_server_version"`
	CmSessionTimeoutMinutes             int           `json:"cm_session_timeout_minutes"`
	CnmHttpMaxThreadsPerWorker          int           `json:"cnm_http_max_threads_per_worker"`
	CnmHttpPort                         int           `json:"cnm_http_port"`
	CnmHttpWorkers                      int           `json:"cnm_http_workers"`
	CnmHTTPSPort                        int           `json:"cnm_https_port"`
	ControlCipherSuites                 string        `json:"control_cipher_suites"`
	ControlCipherSuitesTLS13            string        `json:"control_cipher_suites_tls_1_3"`
	CrdbCoordinatorIgnoreRequests       bool          `json:"crdb_coordinator_ignore_requests"`
	CrdbCoordinatorPort                 int           `json:"crdb_coordinator_port"`
	CrdtSupportedFeaturesetVersion      int           `json:"crdt_supported_featureset_version"`
	CrdtSupportedProtocolVersions       []string      `json:"crdt_supported_protocol_versions"`
	CreatedTime                         time.Time     `json:"created_time"`
	DataCipherList                      string        `json:"data_cipher_list"`
	DataCipherSuitesTLS13               []interface{} `json:"data_cipher_suites_tls_1_3"`
	DebuginfoPath                       string        `json:"debuginfo_path"`
	EmailAlerts                         bool          `json:"email_alerts"`
	EncryptPkeys                        bool          `json:"encrypt_pkeys"`
	EnvoyAdminPort                      int           `json:"envoy_admin_port"`
	EnvoyMaxDownstreamConnections       int           `json:"envoy_max_downstream_connections"`
	EnvoyMgmtServerPort                 int           `json:"envoy_mgmt_server_port"`
	GossipEnvoyAdminPort                int           `json:"gossip_envoy_admin_port"`
	HandleMetricsRedirects              bool          `json:"handle_metrics_redirects"`
	HandleRedirects                     bool          `json:"handle_redirects"`
	HttpSupport                         bool          `json:"http_support"`
	MetricsSystem                       int           `json:"metrics_system"`
	MinControlTLSVersion                string        `json:"min_control_TLS_version"`
	MinDataTLSVersion                   string        `json:"min_data_TLS_version"`
	MinSentinelTLSVersion               string        `json:"min_sentinel_TLS_version"`
	ModuleUploadMaxSizeMb               int           `json:"module_upload_max_size_mb"`
	MtlsAuthorizedSubjects              []interface{} `json:"mtls_authorized_subjects"`
	MtlsCertificateAuthentication       bool          `json:"mtls_certificate_authentication"`
	MtlsClientCertSubjectValidationType string        `json:"mtls_client_cert_subject_validation_type"`
	Name                                string        `json:"name"`
	OptionsMethodForbidden              bool          `json:"options_method_forbidden"`
	PasswordComplexity                  bool          `json:"password_complexity"`
	PasswordExpirationDuration          int           `json:"password_expiration_duration"`
	PasswordMinLength                   int           `json:"password_min_length"`
	ProxyCertificate                    string        `json:"proxy_certificate"`
	RackAware                           bool          `json:"rack_aware"`
	ReservedPorts                       []interface{} `json:"reserved_ports"`
	S3CertificateVerification           bool          `json:"s3_certificate_verification"`
	SentinelCipherSuites                []interface{} `json:"sentinel_cipher_suites"`
	SentinelCipherSuitesTLS13           string        `json:"sentinel_cipher_suites_tls_1_3"`
	SentinelTLSMode                     string        `json:"sentinel_tls_mode"`
	SlaveHa                             bool          `json:"slave_ha"`
	SlaveHaBdbCooldownPeriod            int           `json:"slave_ha_bdb_cooldown_period"`
	SlaveHaCooldownPeriod               int           `json:"slave_ha_cooldown_period"`
	SlaveHaGracePeriod                  int           `json:"slave_ha_grace_period"`
	SlowlogInSanitizedSupport           bool          `json:"slowlog_in_sanitized_support"`
	SMTPTLSMode                         string        `json:"smtp_tls_mode"`
	SMTPUseTLS                          bool          `json:"smtp_use_tls"`
	SyncerCertificate                   string        `json:"syncer_certificate"`
	SystemReservedPorts                 []int         `json:"system_reserved_ports"`
	UpgradeMode                         bool          `json:"upgrade_mode"`
	UseExternalIpv6                     bool          `json:"use_external_ipv6"`
	UseIpv6                             bool          `json:"use_ipv6"`
	WaitCommand                         bool          `json:"wait_command"`
}

type Role struct {
	Management string `json:"management"`
	Name       string `json:"name"`
	UID        int    `json:"uid"`
}
