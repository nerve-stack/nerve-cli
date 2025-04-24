package cmd

import (
	"errors"
	"fmt"
	"slices"
	"strings"

	"github.com/nerve-stack/nerve-cli/internal/codegen/golang"
	"github.com/nerve-stack/nerve-cli/internal/schema"
	"github.com/spf13/cobra"
)

type target string

const (
	targetClient target = "client"
	targetServer target = "server"
)

type lang string

const (
	langTS lang = "ts"
	langGo lang = "go"
)

var configExts = []string{
	string(langTS),
	string(langGo),
}

func (e *lang) String() string {
	return string(*e)
}

func (e *lang) Set(v string) error {
	if slices.Contains(configExts, v) {
		*e = lang(v)

		return nil
	}

	return fmt.Errorf("must be one of: %s", strings.Join(configExts, ", "))
}

func (e *lang) Type() string {
	return "string"
}

func codegenCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "codegen",
		Short: "Generate code for specified language and target",
	}

	// Add subcommands for different targets
	cmd.AddCommand(newTargetCmd(targetServer))
	cmd.AddCommand(newTargetCmd(targetClient))

	return cmd
}

func newTargetCmd(t target) *cobra.Command {
	var language lang

	subCmd := &cobra.Command{
		Use:   fmt.Sprintf("%s <schema-path>", t),
		Short: fmt.Sprintf("Generate %s code", t),
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			schemaPath := args[0]

			spec, err := schema.Parse(schemaPath)
			if err != nil {
				return err
			}

			fmt.Printf("Generating %s code for %s with schema: %s\n", language, t, schemaPath)

			switch t {
			case targetClient:
				return errors.New("NOT IMPLEMENTED YET")

			case targetServer:
				switch language {
				case langGo:
					return golang.GenServer(cmd.OutOrStdout(), spec)
				case langTS:
					return errors.New("NOT IMPLEMENTED YET")
				default:
					panic(fmt.Sprintf("unknown language: %s", language))
				}
			default:
				panic(fmt.Sprintf("unknown target: %s", t))
			}
		},
	}

	subCmd.Flags().VarP(&language, "language", "l", "Programming language for code generation (ts, go)")
	subCmd.MarkFlagRequired("language")

	return subCmd
}
