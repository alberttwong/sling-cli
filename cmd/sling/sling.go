package main

import (
	"context"
	"embed"
	"net/http"
	_ "net/http/pprof"
	"os"
	"os/signal"
	"runtime"
	"time"

	"github.com/flarco/sling/core"
	"github.com/flarco/sling/core/env"
	"github.com/getsentry/sentry-go"
	"github.com/rudderlabs/analytics-go"

	"github.com/flarco/dbio/database"
	"github.com/flarco/dbio/iop"
	"github.com/flarco/g"
	"github.com/integrii/flaggy"
	"github.com/spf13/cast"
)

//go:embed *
var slingFolder embed.FS
var examples = ``
var ctx, cancel = context.WithCancel(context.Background())

var cliRun = &g.CliSC{
	Name:        "run",
	Description: "execute an ad-hoc task",
	Flags: []g.Flag{
		{
			Name:        "config",
			ShortName:   "c",
			Type:        "string",
			Description: "The config string or file to use (JSON or YAML).\n",
		},
		{
			Name:        "src-conn",
			ShortName:   "",
			Type:        "string",
			Description: "The source database / API connection (name, conn string or URL).",
		},
		{
			Name:        "src-stream",
			ShortName:   "",
			Type:        "string",
			Description: "The source table (schema.table), local / cloud file path or API supported object name. Can also be the path of sql file or in-line text to use as query. Use `file://` for local paths.",
		},
		{
			Name:        "src-options",
			Type:        "string",
			Description: "in-line options to further configure source (JSON or YAML).\n",
		},
		{
			Name:        "tgt-conn",
			ShortName:   "",
			Type:        "string",
			Description: "The target database connection (name, conn string or URL).",
		},
		{
			Name:        "tgt-object",
			ShortName:   "",
			Type:        "string",
			Description: "The target table (schema.table) or local / cloud file path. Use `file://` for local paths.",
		},
		{
			Name:        "tgt-options",
			Type:        "string",
			Description: "in-line options to further configure target (JSON or YAML).\n",
		},
		{
			Name:        "stdout",
			ShortName:   "",
			Type:        "bool",
			Description: "Output the stream to standard output (STDOUT).",
		},
		{
			Name:        "mode",
			ShortName:   "m",
			Type:        "string",
			Description: "The target load mode to use: append, upsert, truncate, drop.\n                     Default is append. For upsert, must provide `update_key` and `primary_key` in a configuration file.\n                     All modes load into a new temp table on tgtConn prior to final load.",
		},
		{
			Name:        "examples",
			ShortName:   "e",
			Type:        "bool",
			Description: "Shows some examples.",
		},
	},
	ExecProcess: processRun,
}

var cliInteractive = &g.CliSC{
	Name:        "it",
	Description: "launch interactive mode",
	ExecProcess: slingPrompt,
}

var cliServer = &g.CliSC{
	Name:        "server",
	Description: "launch REST API server mode",
	ExecProcess: slingServer,
	Flags: []g.Flag{
		{
			Name:        "host",
			Type:        "string",
			Description: "The host to bind to. Default is `localhost`",
		},
		{
			Name:        "port",
			Type:        "string",
			Description: "The port to bind to. Default is `9876`",
		},
	},
}

var cliUpdate = &g.CliSC{
	Name:        "update",
	Description: "update to the latest version",
	ExecProcess: updateCLI,
}

var cliConns = &g.CliSC{
	Name:        "conns",
	Singular:    "local connection",
	Description: "manage local connections",
	SubComs: []*g.CliSC{
		// {
		// 	Name:        "add",
		// 	Description: "add new connection",
		// },
		{
			Name:        "show",
			Description: "list local connections detected",
		},
		{
			Name:        "test",
			Description: "test a local connection",
			Flags: []g.Flag{
				{
					Name:        "name",
					ShortName:   "",
					Type:        "string",
					Description: "The name of the connection to test",
				},
			},
		},
	},
	ExecProcess: processConns,
}

func init() {
	// we need a webserver to get the pprof webserver
	if os.Getenv("SLING_PPROF") == "TRUE" {
		go func() {
			g.Trace("Starting pprof webserver @ localhost:6060")
			g.LogError(http.ListenAndServe("localhost:6060", nil))
		}()
	}

	// collect examples
	examplesBytes, _ := slingFolder.ReadFile("examples.sh")
	examples = string(examplesBytes)
	cliInteractive.Make().Add()
	cliConns.Make().Add()
	cliRun.Make().Add()
	cliServer.Make().Add()
	cliUpdate.Make().Add()

	sentry.Init(sentry.ClientOptions{
		// Either set your DSN here or set the SENTRY_DSN environment variable.
		Dsn: "https://abb36e36341a4a2fa7796b6f9a0b3766@o881232.ingest.sentry.io/5835484",
		// Either set environment and release here or set the SENTRY_ENVIRONMENT
		// and SENTRY_RELEASE environment variables.
		// Environment: "",
		Release: "sling@" + core.Version,
		// Enable printing of SDK debug messages.
		// Useful when getting started or trying to figure something out.
		Debug: false,
	})
}

func getRsClient() analytics.Client {
	return analytics.New("1uXKxEPgIB9HXkTduvfwXmFak2l", "https://liveflarccszw.dataplane.rudderstack.com")
}

func Track(event string, props ...g.Map) {
	if val := os.Getenv("SLING_SEND_ANON_USAGE"); val != "" {
		if !cast.ToBool(val) {
			return
		}
	}

	rsClient := getRsClient()
	properties := analytics.NewProperties().
		Set("application", "sling-cli").
		Set("version", core.Version).
		Set("os", runtime.GOOS)

	if len(props) > 0 {
		for k, v := range props[0] {
			properties.Set(k, v)
		}
	}

	rsClient.Enqueue(analytics.Track{
		UserId:     "sling",
		Event:      event,
		Properties: properties,
	})
	rsClient.Close()
}

func main() {

	exitCode := 11
	done := make(chan struct{})
	interrupt := make(chan os.Signal, 1)
	kill := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)
	signal.Notify(kill, os.Kill)

	iop.ShowProgress = os.Getenv("SLING_SHOW_PROGRESS") != "false"
	database.UseBulkExportFlowCSV = cast.ToBool(os.Getenv("SLING_BULK_EXPORT_FLOW_CSV"))

	go func() {
		defer close(done)
		exitCode = cliInit()
	}()

	select {
	case <-done:
		os.Exit(exitCode)
	case <-kill:
		println("\nkilling process...")
		os.Exit(111)
	case <-interrupt:
		if cliRun.Sc.Used {
			println("\ninterrupting...")
			cancel()
			select {
			case <-done:
			case <-time.After(5 * time.Second):
			}
		}
		os.Exit(exitCode)
		return
	}
}

func cliInit() int {
	env.InitLogger()

	// Set your program's name and description.  These appear in help output.
	flaggy.SetName("sling")
	flaggy.SetDescription("An Extract-Load tool.")
	flaggy.DefaultParser.ShowHelpOnUnexpected = true
	flaggy.DefaultParser.AdditionalHelpPrepend = "Slings data from a data source to a data target.\nVersion " + core.Version

	flaggy.SetVersion(core.Version)
	for _, cli := range g.CliArr {
		flaggy.AttachSubcommand(cli.Sc, 1)
	}

	flaggy.Parse()

	ok, err := g.CliProcess()
	if ok {
		if err != nil {
			sentry.CaptureException(err)
		}
		g.LogFatal(err)
	} else {
		flaggy.ShowHelp("")
		Track("ShowHelp")
	}
	return 0
}
