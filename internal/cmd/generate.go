package cmd

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"github.com/nerve-stack/nerve-cli/internal/codegen/golang"
	"github.com/nerve-stack/nerve-cli/internal/config"
	"github.com/nerve-stack/nerve-cli/internal/schema"
	"github.com/nerve-stack/nerve-cli/pkg/cases"
	"github.com/spf13/cobra"
)

func codegenCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "generate",
		Aliases: []string{"gen"},
		Short:   "Generate code from spec",
		Args:    cobra.NoArgs,
		RunE: func(cmd *cobra.Command, _ []string) error {
			cfg, cfgPath, err := config.GetConfig()
			if err != nil {
				return fmt.Errorf("failed to get config: %w", err)
			}

			baseDir := filepath.Dir(cfgPath)

			spec, err := schema.Parse(cfg.Schema)
			if err != nil {
				return fmt.Errorf("failed to parse schema: %w", err)
			}

			for _, output := range cfg.Outputs {
				switch output.Target {
				case config.TargetClient:
					return errors.New("NOT IMPLEMENTED YET")

				case config.TargetServer:
					switch output.Language {
					case config.LanguageGo:
						model, err := golang.ParseSpec(spec)
						if err != nil {
							return err
						}

						model.NerveVersion = Version
						model.Package = cases.ToGoPkgName(filepath.Base(output.Out))

						outSource, err := golang.RenderModelToBuffer(model)
						if err != nil {
							return err
						}

						for filename, source := range outSource {
							path := filepath.Join(baseDir, output.Out, filename)

							if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
								return fmt.Errorf("failed to mkdirall %s: %w", filepath.Base(path), err)
							}

							if err := os.WriteFile(path, []byte(source), 0o644); err != nil {
								return fmt.Errorf("failed to write file %s: %w", filename, err)
							}
						}

						return nil
					case config.LanguageTS:
						return errors.New("NOT IMPLEMENTED YET")
					default:
						panic(fmt.Sprintf("unknown language: %s", output.Language))
					}
				default:
					panic(fmt.Sprintf("unknown target: %s", output.Target))
				}
			}
			return nil
		},
	}

	return cmd
}
