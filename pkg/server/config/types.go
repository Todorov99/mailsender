package config

type ApplicationProperties struct {
	SMTPServerCfg SMTPServerCfg `yaml:"SMTPServerCfg,omitempty"`
	Security      Security      `yaml:"security,omitempty"`
	VaultType     string        `yaml:"vaultType,omitempty"`
}

type Security struct {
	TLS TLS `yaml:"tls,omitempty"`
}

type TLS struct {
	CertFile   string `yaml:"certFile,omitempty"`
	PrivateKey string `yaml:"privateKey,omitempty"`
	RootCACert string `yaml:"rootCACert,omitempty"`
	RootCAKey  string `yaml:"rootCAKey,omitempty"`
}

type SMTPServerCfg struct {
	Host              string `yaml:"host,omitempty"`
	Port              int    `yaml:"port,omitempty"`
	PasswordSecret    string `yaml:"passwordSecret,omitempty"`
	KeepAlive         bool   `yaml:"keepAlive,omitempty"`
	ConnectionTimeout string `yaml:"connectionTimeout,omitempty"`
	SendTimeout       string `yaml:"sendTimeout,omitempty"`
}
