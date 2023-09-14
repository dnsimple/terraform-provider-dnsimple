package consts

const (
	BaseURLSandbox        = "https://api.sandbox.dnsimple.com"
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
