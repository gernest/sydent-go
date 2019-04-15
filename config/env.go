package config

import (
	"os"
	"reflect"
	"strings"
)

func LoadFromEnv() *Matrix {
	return &Matrix{
		Mode: os.Getenv("MX_MODE"),
		Server: Server{
			Name:           env("MX_SERVER_NAME"),
			Port:           env("MX_SERVER_PORT"),
			ClientHTTPBase: env("MX_CLIENT_HTTP_BASE"),
			Crypto: Crypto{
				Algorithm:  env("MX_SERVER_CRYPTO_ALG"),
				Version:    env("MX_SERVER_CRYPTO_VER"),
				SingingKey: env("MX_SERVER_CRYPTO_SIGN_KEY"),
				VerifyKey:  env("MX_SERVER_CRYPTO_VERIFY_KEY"),
			},
		},
		DB: DB{
			Driver: env("MX_DB_DRIVER"),
			Conn:   env("MX_DB_CONN"),
		},
		Email: Email{
			Invite: Invite{
				From:     env("MX_EMAIL_INVITE_FROM"),
				Template: env("MX_EMAIL_INVITE_TEMPLATE"),
			},
			Verification: Verification{
				From:         env("MX_EMAIL_VERIFY_FROM"),
				Template:     env("MX_EMAIL_VERIFY_TEMPLATE"),
				ResponsePage: env("MX_EMAIL_VERIFY_RESPONSE_TEMPLATE"),
			},
			Providers: providers(os.Environ()),
		},
	}
}

func env(name string) string {
	return os.Getenv(name)
}

type eVar struct {
	k []string
	v string
}

func providers(environ []string) []Provider {
	var values []eVar
	for _, v := range environ {
		p := strings.Split(v, "=")
		if len(p) == 2 {
			if strings.HasPrefix(p[0], "MX_EMAIL_PROVIDER_") {
				values = append(values, eVar{
					k: strings.Split(strings.TrimPrefix(p[0], "MX_EMAIL_PROVIDER_"), "_"),
					v: p[1],
				})
			}
		}

	}
	var emailProviders []string
	for _, v := range values {
		if len(v.k) == 1 {
			emailProviders = append(emailProviders, v.k[0])
		}
	}
	var p []Provider
	for _, v := range emailProviders {
		pr := Provider{
			Name:     strings.ToLower(v),
			Settings: make(map[string]interface{}),
		}
		for _, x := range values {
			if len(x.k) > 0 && x.k[0] == v {
				switch len(x.k) {
				case 1:
					pr.State = x.v
				default:
					pr.Settings[strings.ToLower(strings.Join(x.k[1:], "_"))] = x.v
				}
			}

		}
		p = append(p, pr)
	}
	return p
}

// Empty returns true if c is empty.
func Empty(c *Matrix) bool {
	if c == nil {
		return true
	}
	a := *c
	return reflect.DeepEqual(a, Matrix{})
}
