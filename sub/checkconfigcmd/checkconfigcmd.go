package checkconfigcmd

import (
	"flag"
	"fmt"
	"os"

	"github.com/M-ERCURY/core/cli"
	"github.com/M-ERCURY/core/cli/fsdir"
	"github.com/M-ERCURY/relay/filenames"
	"github.com/M-ERCURY/relay/relaycfg"
)

var Cmd = &cli.Subcmd{
	FlagSet: flag.NewFlagSet("check-config", flag.ExitOnError),
	Desc:    "Validate mercury-relay config file",
	Run: func(fm fsdir.T) {
		c := relaycfg.Defaults()
		if err := fm.Get(&c, filenames.Config); err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		if err := c.Validate(); err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		fmt.Println("OK")
	},
}
