package config

type ApplicationProperties struct {
	SMTPServerCfg SMTPServerCfg `yaml:"SMTPServerCfg,omitempty"`
}

type SMTPServerCfg struct {
	Host              string `yaml:"host,omitempty"`
	Port              int    `yaml:"port,omitempty"`
	PasswordSecret    string `yaml:"passwordSecret,omitempty"`
	KeepAlive         bool   `yaml:"keepAlive,omitempty"`
	ConnectionTimeout string `yaml:"connectionTimeout,omitempty"`
	SendTimeout       string `yaml:"sendTimeout,omitempty"`
}
