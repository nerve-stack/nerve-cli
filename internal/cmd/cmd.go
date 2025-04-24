package cmd

import (
	"context"
	"errors"
	"fmt"
	"io"
	"os/exec"

	"github.com/nerve-stack/nerve-cli/internal/schema"
	"github.com/spf13/cobra"
)

func Do(args []string, stdin io.Reader, stdout io.Writer, stderr io.Writer) int {
	rootCmd := &cobra.Command{
		Use:   "nerve",
		Short: "Nerve CLI",
	}

	rootCmd.SetArgs(args)
	rootCmd.SetIn(stdin)
	rootCmd.SetOut(stdout)
	rootCmd.SetErr(stderr)

	rootCmd.AddCommand(versionCmd())
	rootCmd.AddCommand(codegenCmd())

	ctx := context.Background()
	if err := rootCmd.ExecuteContext(ctx); err != nil {
		var exitError *exec.ExitError
		if errors.As(err, &exitError) {
			return exitError.ExitCode()
		}

		return 1
	}

	return 0
}

const Version = "v0.0.1"

func versionCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "version",
		Short: "Print version",
		RunE: func(cmd *cobra.Command, args []string) error {
			_, err := fmt.Fprintf(cmd.OutOrStdout(), "%s\n", Version)
			return err
		},
	}
}

func codegenCmd() *cobra.Command {
	command := &cobra.Command{
		Use:   "codegen",
		Short: "Generate code for specified language",
		RunE: func(cmd *cobra.Command, args []string) error {
			language, _ := cmd.Flags().GetString("language")
			target, _ := cmd.Flags().GetString("target")
			schemaPath, _ := cmd.Flags().GetString("schema")

			spec, err := schema.Parse(schemaPath)
			if err != nil {
				return err
			}

			// Print out the parsed schema and options (just for illustration)
			fmt.Printf("Generating %s code for %s with schema: %s\n", language, target, schemaPath)
			fmt.Printf("Parsed Spec: %+v\n", spec)

			return nil
		},
	}

	// Add flags inside this function
	command.Flags().StringP("language", "l", "", "Programming language for code generation (ts, go, rust)")
	command.Flags().StringP("target", "t", "", "Target code generation (client, server)")
	command.Flags().StringP("schema", "s", "", "Path to the schema file (JSON or YAML)")

	// Mark flags as required
	command.MarkFlagRequired("language")
	command.MarkFlagRequired("target")
	command.MarkFlagRequired("schema")

	return command
}
