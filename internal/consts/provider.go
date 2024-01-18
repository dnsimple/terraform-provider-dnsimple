package consts

const (
	BaseURLSandbox = "https://api.sandbox.dnsimple.com"

	// Certificate states
	CertificateStateCancelled = "cancelled"
	CertificateStateFailed    = "failed"
	CertificateStateIssued    = "issued"
	CertificateStateRefunded  = "refunded"

	// Domain states
	DomainStateRegistered = "registered"
	DomainStateHosted     = "hosted"
	DomainStateNew        = "new"
	DomainStateFailed     = "failed"
	DomainStateCancelling = "cancelling"
	DomainStateCancelled  = "cancelled"

	// Domain Registrant Change (Contact change) states
	RegistrantChangeStateNew        = "new"
	RegistrantChangeStatePending    = "pending"
	RegistrantChangeStateCompleted  = "completed"
	RegistrantChangeStateCancelling = "cancelling"
	RegistrantChangeStateCancelled  = "cancelled"
)
