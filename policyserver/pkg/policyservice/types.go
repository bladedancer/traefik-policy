package policyservice

// PolicyServiceConfig for the PolicyService
type PolicyServiceConfig struct {
	Host string
	Port uint32

	CertFile string
	KeyFile  string

	HeaderPolicy bool
	BodyPolicy   bool
}
