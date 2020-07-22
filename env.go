package db

import "taylz.io/env"

var _env = []string{"USER", "PASSWORD", "HOST", "PORT", "NAME"}

// ENV returns env.Service containing relevant settings
func ENV() env.Service {
	env := env.Service{}
	for _, v := range _env {
		env[v] = ""
	}
	return env
}

// ParseDSN returns DSN from env type
func ParseDSN(env env.Service) string {
	return DSN(env[_env[0]], env[_env[1]], env[_env[2]], env[_env[3]], env[_env[4]])
}
