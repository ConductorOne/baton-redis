package config

//go:generate go run ./gen

import (
	"github.com/conductorone/baton-sdk/pkg/field"
)

var (
	ClusterHostField = field.StringField(
		"cluster-host",
		field.WithDescription("The enterprise cluster host"),
		field.WithRequired(true),
	)
	ApiPortField = field.StringField(
		"api-port",
		field.WithDescription("The enterprise API port"),
		field.WithDefaultValue("9443"),
	)
	UsernameField = field.StringField(
		"username",
		field.WithDescription("The enterprise cluster admin username"),
		field.WithRequired(true),
	)
	PasswordField = field.StringField(
		"password",
		field.WithDescription("The enterprise cluster admin password"),
		field.WithRequired(true),
	)

	// ConfigurationFields defines the external configuration required for the
	// connector to run.
	ConfigurationFields = []field.SchemaField{ClusterHostField, ApiPortField, UsernameField, PasswordField}

	// FieldRelationships defines relationships between the fields.
	FieldRelationships = []field.SchemaFieldRelationship{}

	// Config is the configuration schema for the connector.
	Config = field.Configuration{
		Fields:      ConfigurationFields,
		Constraints: FieldRelationships,
	}
)
