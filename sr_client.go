package main

import (
	"fmt"
	"strings"

	"github.com/linkedin/goavro/v2"
	"github.com/riferrei/srclient"
)

type Client struct {
	Session *srclient.SchemaRegistryClient
}

func NewSchemaRegistryClient(endpoint string, refreshIntervalSeconds int) *Client {
	session := srclient.CreateSchemaRegistryClient(endpoint)

	client := &Client{
		Session: session,
	}

	return client
}

func (c *Client) GetIdByVersion(schemaKey string, schemaVersion int) (int, error) {
	schema, err := c.Session.GetSchemaByVersion(schemaKey, schemaVersion)
	if err != nil {
		return 0, err
	}
	return schema.ID(), nil
}

func (c *Client) GetFullSchemaById(id int) ([]byte, error) {
	rawSchema, err := c.Session.GetSchema(id)
	if err != nil {
		return nil, err
	}

	schema, err := c.getCompleteSchema(rawSchema)
	if err != nil {
		return nil, err
	}

	schemaStr := "[" + strings.Join(schema, ",") + "]"
	codec, err := goavro.NewCodecForStandardJSON(schemaStr)
	if err != nil {
		return nil, err
	}

	return []byte(codec.CanonicalSchema()), nil
}

func (c *Client) getCompleteSchema(schema *srclient.Schema) ([]string, error) {
	var finalSchema []string

	schemaReferences, err := c.getReferences(schema)
	if err != nil {
		return nil, err
	}

	// Appends references into finalSchema
	finalSchema = append(finalSchema, schemaReferences...)

	// Appends envelope schema to entities schema
	finalSchema = append(finalSchema, schema.Schema())

	if len(finalSchema) == 0 {
		return []string{}, fmt.Errorf("[schemaRegistryHandler] empty schema for given subject name")
	}

	return finalSchema, nil
}

func (c *Client) getReferences(parent *srclient.Schema) ([]string, error) {
	if parent.References() == nil {
		return []string{}, nil
	}

	var current []string
	for _, reference := range parent.References() {
		// Fetch reference schema by Version ID and appends to new constructor
		schema, err := c.Session.GetSchemaByVersion(reference.Subject, reference.Version)
		if err != nil {
			return nil, err
		}

		schemaReferences, err := c.getReferences(schema)
		if err != nil {
			return nil, err
		}
		current = append(current, schemaReferences...)

		current = append(current, schema.Schema())
	}

	return current, nil
}
