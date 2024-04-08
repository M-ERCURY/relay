package initcmd

import (
	"flag"
	"log"

	"github.com/M-ERCURY/core/api/tlscert"
	"github.com/M-ERCURY/core/cli"
	"github.com/M-ERCURY/core/cli/commonsub/initcmd"
	"github.com/M-ERCURY/core/cli/fsdir"
	"github.com/M-ERCURY/relay/filenames"
	"github.com/M-ERCURY/relay/sub/initcmd/embedded"
)

var Cmd = &cli.Subcmd{
	FlagSet: flag.NewFlagSet("init", flag.ExitOnError),
	Desc:    "Generate ed25519 keypair and TLS cert/key",
	Run: func(fm fsdir.T) {
		initcmd.Cmd.Run(fm)
		privkey, err := cli.LoadKey(fm, filenames.Seed)
		if err != nil {
			log.Fatal(err)
		}
		if err := tlscert.Generate(fm.Path(filenames.TLSCert), fm.Path(filenames.TLSKey), privkey); err != nil {
			log.Fatal(err)
		}
		if err := cli.UnpackEmbedded(embedded.FS, fm, false); err != nil {
			log.Fatalf("error while unpacking embedded files: %s", err)
		}
	},
}
