package consumer

type SecurityProtocol string
type SASLMechanism string

const (
	SASL_PLAINTEXT SecurityProtocol = "SASL_PLAINTEXT"
)

const (
	SCRAM_SHA_512 SASLMechanism = "SCRAM-SHA-512"
)
