package config

//go:generate go run ./gen

import (
	"github.com/conductorone/baton-sdk/pkg/field"
)

var (
	clusterHostField = field.StringField(
		"cluster-host",
		field.WithDescription("The enterprise cluster host"),
		field.WithDisplayName("Cluster Host"),
		field.WithPlaceholder("Your Redis Enterprise cluster host"),
		field.WithRequired(true),
	)
	apiPortField = field.StringField(
		"api-port",
		field.WithDescription("The enterprise API port"),
		field.WithDisplayName("API Port"),
		field.WithPlaceholder("9443"),
		field.WithDefaultValue("9443"),
	)
	usernameField = field.StringField(
		"username",
		field.WithDescription("The enterprise cluster admin username"),
		field.WithDisplayName("Username"),
		field.WithPlaceholder("Your Redis Enterprise admin username"),
		field.WithRequired(true),
	)
	passwordField = field.StringField(
		"password",
		field.WithDescription("The enterprise cluster admin password"),
		field.WithDisplayName("Password"),
		field.WithPlaceholder("Your Redis Enterprise admin password"),
		field.WithRequired(true),
		field.WithIsSecret(true),
	)
)

var Config = field.NewConfiguration(
	[]field.SchemaField{
		clusterHostField,
		apiPortField,
		usernameField,
		passwordField,
	},
	field.WithConnectorDisplayName("Redis"),
	field.WithIconUrl("/static/app-icons/redis.svg"),
	field.WithHelpUrl("/docs/baton/redis"),
)

func ValidateConfig(c *Redis) error {
	return nil
}
