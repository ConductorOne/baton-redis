package config

//go:generate go run ./gen

import (
	"github.com/conductorone/baton-sdk/pkg/field"
)

var Config = field.NewConfiguration([]field.SchemaField{
	field.StringField(
		"cluster-host",
		field.WithDescription("The enterprise cluster host"),
		field.WithRequired(true),
	),
	field.StringField(
		"api-port",
		field.WithDescription("The enterprise API port"),
		field.WithDefaultValue("9443"),
	),
	field.StringField(
		"username",
		field.WithDescription("The enterprise cluster admin username"),
		field.WithRequired(true),
	),
	field.StringField(
		"password",
		field.WithDescription("The enterprise cluster admin password"),
		field.WithRequired(true),
	),
})

func ValidateConfig(c *Redis) error {
	return nil
}
