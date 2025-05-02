package console

import (
	"io"
)

type Config struct {
	DataDir         string       // Data directory to store the console history at
	DocRoot         string       // Filesystem path from where to load JavaScript files from
	Prompt          string       // Input prompt prefix string (defaults to DefaultPrompt)
	Prompter        UserPrompter // Input prompter to allow interactive user feedback (defaults to TerminalPrompter)
	Printer         io.Writer    // Output writer to serialize any display strings to (defaults to os.Stdout)
	Preload         []string     // Absolute paths to JavaScript files to preload
	DefaultCategory string       // the default cmd category
}

func NewConsoleConfig() Config {
	return Config{
		DataDir:         "./data",
		DocRoot:         "./doc",
		DefaultCategory: "ll",
	}
}
