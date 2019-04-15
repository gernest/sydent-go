package config

import (
	"reflect"
	"testing"
)

func TestProviders(t *testing.T) {

	p := providers([]string{
		"MX_EMAIL_PROVIDER_SMTP=enabled",
		"MX_EMAIL_PROVIDER_SMTP_USERNAME=root",
		"MX_EMAIL_PROVIDER_SMTP_PASSWORD=pass",
	})
	expect := []Provider{
		{
			Name:  "smtp",
			State: "enabled",
			Settings: map[string]interface{}{
				"username": "root",
				"password": "pass",
			},
		},
	}
	if !reflect.DeepEqual(p, expect) {
		t.Errorf("expected %#v got %#v", expect, p)
	}
}

func TestEmptry(t *testing.T) {
	x := Empty(LoadFromEnv())
	if !x {
		t.Error("expected true")
	}
}
