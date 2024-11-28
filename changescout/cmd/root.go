package cmd

import (
	"github.com/gelleson/changescout/changescout/cmd/commands"
	"github.com/gelleson/changescout/changescout/internal/pkg/clis"
	"github.com/gelleson/changescout/changescout/internal/platform/logger"
	"github.com/gelleson/changescout/changescout/pkg/flags"
	"github.com/urfave/cli/v2"
	"log"
	"os"
)

var root = cli.App{
	Name:  "changescout",
	Usage: "A service for monitoring websites",

	Flags: flags.Build(
		clis.FlagsLogLevel,
	),
	Commands: []*cli.Command{
		commands.StartServer,
	},
	Action: func(c *cli.Context) error {
		logger.SetLevel(c.String("log-level"))
		return nil
	},
}

func Execute() {
	err := root.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
