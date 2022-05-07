package config

import "github.com/Todorov99/mailsender/pkg/global"

var serverConfig *serverCfg

type serverCfg struct {
	props *ApplicationProperties
	tls   *TLS
}

func init() {
	serverCfg := &serverCfg{}

	applicationProperties, err := LoadApplicationProperties(global.ApplicationPropertyFile)
	if err != nil {
		panic(err)
	}

	serverCfg.props = applicationProperties
	serverCfg.tls = &applicationProperties.Security.TLS

	serverConfig = serverCfg
}

func ServerProps() *ApplicationProperties {
	return serverConfig.props
}

func GetTLSCfg() *TLS {
	return serverConfig.tls
}
