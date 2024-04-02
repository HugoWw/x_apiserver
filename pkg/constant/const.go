package constant

const (
	BaseResource    = "base_resource"
	BaseResourceAPI = "/v1/apis"
)

// define request header info
const (
	X_AUTH_TOKEN    string = "X-Auth-Token"
	X_NM_PAGE       string = "X-Nm-Page"
	X_TRN_ID        string = "X-Transaction-Id"
	X_AS_STANDALONE string = "X-As-Standalone"
	X_SUSE_TOKEN    string = "X-R-Sess"
	X_LX_PROXY      string = "X-LX-Proxy"
)

const (
	ComplianceTemplatePCI   = "PCI"
	ComplianceTemplateGDPR  = "GDPR"
	ComplianceTemplateHIPAA = "HIPAA"
	ComplianceTemplateNIST  = "NIST" // NIST SP 800-190
)
