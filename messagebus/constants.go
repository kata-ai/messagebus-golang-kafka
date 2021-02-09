package messagebus

type SecurityProtocol string
type SASLMechanism string

const VERSION = "1.0.0"
const (
	SASL_PLAINTEXT SecurityProtocol = "SASL_PLAINTEXT"
)

const (
	SCRAM_SHA_512 SASLMechanism = "SCRAM-SHA-512"
)
