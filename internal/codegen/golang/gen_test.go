package golang_test

import (
	"bytes"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/nerve-stack/nerve-cli/internal/codegen/golang"
	"github.com/nerve-stack/nerve-cli/internal/schema"
)

func strPtr(s string) *string {
	return &s
}

func TestMethodsAndEventsDecalrations(t *testing.T) {
	tests := []struct {
		name     string
		spec     *schema.Spec
		expected string
	}{
		{
			name: "Single method and event",
			spec: &schema.Spec{
				Version: "1.0.0",
				Info: schema.Info{
					Name:    "ExampleService",
					Version: strPtr("v1"),
				},
				Schemas: schema.Schemas{
					Types: map[string]schema.SchemaOrRefWithDesc{
						"User": {
							SchemaOrRef: schema.SchemaOrRef{
								Type: strPtr("object"),
								Properties: map[string]schema.SchemaOrRef{
									"id":   {Type: strPtr("string")},
									"name": {Type: strPtr("string")},
								},
								Required: []string{"id"},
							},
							Description: strPtr("Represents a user."),
						},
					},
					Entities: map[string]schema.SchemaOrRefWithDesc{},
					Errors: map[string]schema.ErrorSchema{
						"NotFound": {
							Code:        404,
							Description: strPtr("User not found"),
						},
					},
				},
				Methods: map[string]schema.Method{
					"getUser": {
						Description: strPtr("Fetch user by ID"),
						Params: &schema.SchemaOrRef{
							Type: strPtr("object"),
							Properties: map[string]schema.SchemaOrRef{
								"userId": {Type: strPtr("string")},
							},
							Required: []string{"userId"},
						},
						Result: &schema.SchemaOrRef{
							Ref: strPtr("#/schemas/types/User"),
						},
					},
				},
				Events: map[string]schema.SchemaOrRefWithDesc{
					"userJoined": {
						SchemaOrRef: schema.SchemaOrRef{
							Ref: strPtr("#/schemas/types/User"),
						},
						Description: strPtr("Triggered when a user joins."),
					},
				},
			},
			expected: `
package gen

const (
	MethodGetUser = "getUser"
	EventUserJoined = "userJoined"
)
`[1:],
		},
		{
			name: "Multiple methods and events",
			spec: &schema.Spec{
				Version: "1.0.1",
				Info: schema.Info{
					Name:    "AnotherService",
					Version: strPtr("v2"),
				},
				Schemas: schema.Schemas{
					Types:    map[string]schema.SchemaOrRefWithDesc{},
					Entities: map[string]schema.SchemaOrRefWithDesc{},
					Errors:   map[string]schema.ErrorSchema{},
				},
				Methods: map[string]schema.Method{
					"createUser": {Description: strPtr("Create a new user")},
					"deleteUser": {Description: strPtr("Delete a user")},
				},
				Events: map[string]schema.SchemaOrRefWithDesc{
					"userCreated": {Description: strPtr("Emitted when a new user is created")},
					"userDeleted": {Description: strPtr("Emitted when a user is deleted")},
				},
			},
			expected: `
package gen

const (
	MethodCreateUser = "createUser"
	MethodDeleteUser = "deleteUser"
	EventUserCreated = "userCreated"
	EventUserDeleted = "userDeleted"
)
`[1:],
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var buf bytes.Buffer

			if err := golang.GenServer(&buf, tt.spec); err != nil {
				t.Fatalf("GenServer failed: %v", err)
			}

			output := buf.String()

			if diff := cmp.Diff(tt.expected, output); diff != "" {
				t.Errorf("Generated output mismatch (-want +got):\n%s", diff)
			}
		})
	}
}
