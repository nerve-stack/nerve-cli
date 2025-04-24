package config

type Config struct {
	Version string   `json:"version" yaml:"version"`
	Schema  string   `json:"schema"  yaml:"schema"`
	Outputs []Output `json:"outputs" yaml:"outputs"`
}

type target string

const (
	TargetClient target = "client"
	TargetServer target = "server"
)

type language string

const (
	LanguageTS language = "ts"
	LanguageGo language = "go"
)

type Output struct {
	Target   target   `json:"target"   yaml:"target"`
	Language language `json:"language" yaml:"language"`
	Out      string   `json:"out"      yaml:"out"`
}
