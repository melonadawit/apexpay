package secrets

import "os"

// Vault is an abstraction over secret storage. In dev it reads from the environment; in
// production it can be backed by a real secret manager (e.g. AWS/GCP KMS, HashiCorp Vault)
// by implementing the same interface and wiring it in main.
type Vault interface {
	// Get returns the secret value for a key.
	Get(key string) (string, error)
	// GetOrEnv returns the secret, falling back to the given env var if the vault has none.
	GetOrEnv(key, envName string) (string, error)
}

// EnvVault reads secrets from the process environment. Suitable for dev/local; swap in a
// KMS-backed Vault for production without changing callers.
type EnvVault struct{}

func NewEnv() *EnvVault { return &EnvVault{} }

func (EnvVault) Get(key string) (string, error) {
	v := os.Getenv("APEXPAY_SECRET_" + key)
	if v == "" {
		return "", os.ErrNotExist
	}
	return v, nil
}

func (EnvVault) GetOrEnv(key, envName string) (string, error) {
	v, err := EnvVault{}.Get(key)
	if err == nil {
		return v, nil
	}
	return os.Getenv(envName), nil
}
