package global

import (
	"os"

	"github.com/urfave/cli/v2"
	"github.com/vincent78/butil/timewheel"
)

var CliContext *cli.Context

var Timewheel *timewheel.TimeWheel

var GinRepBodyBlackPath = map[string]struct{}{}

const (
	RequestTokenKey = "REQ_TOKEN"

	LogFileHttpName  = "http"  // web logger name
	LogFileNetName   = "net"   // net logger name
	LogFileDBName    = "db"    // mysql logger name
	LogFileMongoName = "mongo" // mongodb logger name
	LogFileWSSName   = "wss"   // websocket server logger name
	LogFileWSCName   = "wsc"   // websocket client logger name

)

var (
	DefautTmpPath = "/tmp"
)

func init() {
	if p, exist := os.LookupEnv("GC_PATH_TEMP"); exist {
		DefautTmpPath = p
	}
}

const (
	AppHelpTemplate = `{{.Name}} {{if .Flags}}[common options] {{end}}command{{if .Flags}} [command options]{{end}} [arguments...]

VERSION:
   {{.Version}}

COMMANDS:
   {{range .Commands}}{{.Name}}{{with .Aliases}}, {{.}}{{end}}{{ "\t" }}{{.Usage}}
   {{end}}{{if .Flags}}

GLOBAL OPTIONS:
   {{range .Flags}}{{.}}
   {{end}}{{end}}
`
	CommandHelpTemplate = `{{.cmd.Name}}{{if .cmd.Subcommands}} command{{end}}{{if .cmd.Flags}} [command options]{{end}} [arguments...]
{{if .cmd.Description}}{{.cmd.Description}}
{{end}}{{if .cmd.Subcommands}}
SUBCOMMANDS:
	{{range .cmd.Subcommands}}{{.Name}}{{with .Aliases}}, {{.}}{{end}}{{ "\t" }}{{.Usage}}
	{{end}}{{end}}{{if .categorizedFlags}}
{{range $idx, $categorized := .categorizedFlags}}{{$categorized.Name}} OPTIONS:
{{range $categorized.Flags}}{{"\t"}}{{.}}
{{end}}
{{end}}{{end}}`

	OriginCommandHelpTemplate = `{{.Name}}{{if .Subcommands}} command{{end}}{{if .Flags}} [command options]{{end}} [arguments...]
{{if .Description}}{{.Description}}
{{end}}{{if .Subcommands}}
SUBCOMMANDS:
	{{range .Subcommands}}{{.Name}}{{with .Aliases}}, {{.}}{{end}}{{ "\t" }}{{.Usage}}
	{{end}}{{end}}{{if .Flags}}
OPTIONS:
{{range $.Flags}}{{"\t"}}{{.}}
{{end}}
{{end}}`
)

const (
	ErrorCodeBase         = 40000
	ErrorCodeHttpBase     = 41000
	ErrorCodeWSServerBase = 42000
	ErrorCodeWSClientBase = 43000
	ErrorCodeWSCommonBase = 44000
	ErrorCodeConsoleBase  = 44000
)
