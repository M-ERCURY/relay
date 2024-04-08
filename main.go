package main

import (
	"fmt"
	"os"

	"github.com/M-ERCURY/core/cli"
	"github.com/M-ERCURY/core/cli/commonsub/commonlib"
	"github.com/M-ERCURY/core/cli/commonsub/migratecmd"
	"github.com/M-ERCURY/core/cli/commonsub/reloadcmd"
	"github.com/M-ERCURY/core/cli/commonsub/restartcmd"
	"github.com/M-ERCURY/core/cli/commonsub/rollbackcmd"
	"github.com/M-ERCURY/core/cli/commonsub/statuscmd"
	"github.com/M-ERCURY/core/cli/commonsub/stopcmd"
	"github.com/M-ERCURY/core/cli/commonsub/superviseupgradecmd"
	"github.com/M-ERCURY/core/cli/commonsub/upgradecmd"
	"github.com/M-ERCURY/core/cli/commonsub/versioncmd"
	"github.com/M-ERCURY/core/cli/upgrade"

	"github.com/M-ERCURY/core/mrnet"
	"github.com/M-ERCURY/relay/sub/balancecmd"
	"github.com/M-ERCURY/relay/sub/checkconfigcmd"
	"github.com/M-ERCURY/relay/sub/initcmd"
	"github.com/M-ERCURY/relay/sub/startcmd"
	"github.com/M-ERCURY/relay/sub/withdrawcmd"
	"github.com/M-ERCURY/relay/version"
)

const binname = "mercury-relay"

func main() {
	cli.CLI{
		Subcmds: []*cli.Subcmd{
			initcmd.Cmd,
			startcmd.Cmd(),
			stopcmd.Cmd(binname),
			restartcmd.Cmd(binname, startcmd.Cmd().Run, stopcmd.Cmd(binname).Run),
			reloadcmd.Cmd(binname),
			statuscmd.Cmd(binname),
			upgradecmd.Cmd(
				binname,
				upgrade.ExecutorSupervised,
				version.VERSION,
				version.LatestChannelVersion,
			),
			superviseupgradecmd.Cmd(commonlib.Context{
				BinName:    binname,
				NewVersion: version.VERSION,
			}),
			migratecmd.Cmd(binname, version.MIGRATIONS, version.VERSION),
			rollbackcmd.Cmd(commonlib.Context{
				BinName:  binname,
				PostHook: version.PostRollbackHook,
			}),
			checkconfigcmd.Cmd,
			balancecmd.Cmd(),
			withdrawcmd.Cmd(),
			versioncmd.Cmd(fmt.Sprintf(
				"%s, protocol version %s",
				version.VERSION_STRING,
				mrnet.PROTO_VERSION.String(),
			)),
		},
	}.Parse(os.Args).Run(cli.Home())
}
