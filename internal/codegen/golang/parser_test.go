package golang

import (
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/nerve-stack/nerve-cli/internal/schema"
	"github.com/nerve-stack/nerve-cli/pkg/ptrto"
)

type testCase struct {
	name     string
	spec     *schema.Spec
	expected *Model
}

func runTestCases(t *testing.T, testCases []testCase) {
	t.Helper()

	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ParseSpec(tt.spec)
			if err != nil {
				t.Fatalf("unexpected error: %s", err)
			}

			if diff := cmp.Diff(got, tt.expected); diff != "" {
				t.Errorf("mismatch (-got +want):\n%s", diff)
			}
		})
	}
}

func TestMethods(t *testing.T) {
	tests := []testCase{
		{
			name: "No params and result",
			spec: &schema.Spec{
				Methods: map[string]schema.Method{
					"sendMessage": {},
				},
			},
			expected: &Model{
				Methods: []GoMethod{
					{
						Name:            "sendMessage",
						CapitalizedName: "SendMessage",
					},
				},
				Imports: map[string]struct{}{},
				Structs: map[string]GoStruct{},
				Errors:  map[string]GoError{},
			},
		},
		{
			name: "String param",
			spec: &schema.Spec{
				Methods: map[string]schema.Method{
					"sendMessage": {
						Params: &schema.SchemaOrRef{
							Type: ptrto.PtrTo("string"),
						},
					},
				},
			},
			expected: &Model{
				Methods: []GoMethod{
					{
						Name:            "sendMessage",
						CapitalizedName: "SendMessage",
						ParamsType:      ptrto.PtrTo("string"),
					},
				},
				Imports: map[string]struct{}{},
				Structs: map[string]GoStruct{},
				Errors:  map[string]GoError{},
			},
		},
		{
			name: "String result",
			spec: &schema.Spec{
				Methods: map[string]schema.Method{
					"sendMessage": {
						Result: &schema.SchemaOrRef{
							Type: ptrto.PtrTo("string"),
						},
					},
				},
			},
			expected: &Model{
				Methods: []GoMethod{
					{
						Name:            "sendMessage",
						CapitalizedName: "SendMessage",
						ResultType:      ptrto.PtrTo("string"),
					},
				},
				Imports: map[string]struct{}{},
				Structs: map[string]GoStruct{},
				Errors:  map[string]GoError{},
			},
		},
		{
			name: "Errors",
			spec: &schema.Spec{
				Errors: map[string]schema.ErrorSchema{
					"UserNotFound": {
						Code:        -32009,
						Description: ptrto.PtrTo("user does not exist"),
					},
					"ChatNotFound": {
						Code: -32010,
						Data: &schema.SchemaOrRef{
							Type: ptrto.PtrTo("string"),
						},
					},
					"MessageTooLong": {
						Code: -32011,
						Data: &schema.SchemaOrRef{
							Type: ptrto.PtrTo("object"),
							Properties: map[string]schema.SchemaOrRef{
								"user_id": {
									Type: ptrto.PtrTo("string"),
								},
								"content": {
									Type: ptrto.PtrTo("string"),
								},
							},
							Required: []string{"content"},
						},
					},
				},
				Methods: map[string]schema.Method{
					"sendMessage": {
						Result: &schema.SchemaOrRef{
							Type: ptrto.PtrTo("string"),
						},
						Errors: []string{
							"UserNotFound",
							"ChatNotFound",
							"MessageTooLong",
						},
					},
				},
			},
			expected: &Model{
				Methods: []GoMethod{
					{
						Name:            "sendMessage",
						CapitalizedName: "SendMessage",
						ResultType:      ptrto.PtrTo("string"),
					},
				},
				Imports: map[string]struct{}{},
				Structs: map[string]GoStruct{
					"MessageTooLongData": {
						Fields: []GoField{
							{
								Name:     "UserId",
								Type:     "string",
								Tag:      "user_id",
								Optional: false,
							},
							{
								Name:     "Content",
								Type:     "string",
								Tag:      "content",
								Optional: false,
							},
						},
					},
				},
				Errors: map[string]GoError{
					"UserNotFoundError": {
						Code: -32009,
					},
					"ChatNotFoundError": {
						Code:     -32010,
						DataType: ptrto.PtrTo("string"),
					},
					"MessageTooLongError": {
						Code:     -32011,
						DataType: ptrto.PtrTo("MessageTooLongData"),
					},
				},
			},
		},
		{
			name: "Methods",
			spec: &schema.Spec{
				Errors: map[string]schema.ErrorSchema{
					"UserNotFound": {
						Code:        -32000,
						Description: ptrto.PtrTo("user does not exist"),
						Data: &schema.SchemaOrRef{
							Type:   ptrto.PtrTo("string"),
							Format: ptrto.PtrTo("uuid"),
						},
					},
				},
				Methods: map[string]schema.Method{
					"sendMessage": {
						Params: &schema.SchemaOrRef{
							Type: ptrto.PtrTo("object"),
							Properties: map[string]schema.SchemaOrRef{
								"user_id": {
									Type:   ptrto.PtrTo("string"),
									Format: ptrto.PtrTo("uuid"),
								},
								"content": {
									Type: ptrto.PtrTo("string"),
								},
							},
							Required: []string{"user_id"},
						},
						Result: &schema.SchemaOrRef{
							Type: ptrto.PtrTo("array"),
							Items: &schema.SchemaOrRef{
								Type: ptrto.PtrTo("string"),
							},
						},
						Errors: []string{
							"UserNotFound",
						},
					},
				},
			},
			expected: &Model{
				Imports: map[string]struct{}{
					UUIDImport: {},
				},
				Structs: map[string]GoStruct{
					"SendMessageParams": {
						Fields: []GoField{
							{
								Name:     "UserId",
								Type:     "uuid.UUID",
								Tag:      "user_id",
								Optional: false,
							},
							{
								Name:     "Content",
								Type:     "string",
								Tag:      "content",
								Optional: true,
							},
						},
					},
				},
				Methods: []GoMethod{
					{
						Name:            "sendMessage",
						CapitalizedName: "SendMessage",
						ParamsType:      ptrto.PtrTo("SendMessageParams"),
						ResultType:      ptrto.PtrTo("[]string"),
					},
				},
				Errors: map[string]GoError{
					"UserNotFoundError": {
						Code:     -32000,
						DataType: ptrto.PtrTo("uuid.UUID"),
					},
				},
			},
		},
	}

	runTestCases(t, tests)
}
