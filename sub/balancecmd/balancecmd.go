package balancecmd

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"

	"github.com/M-ERCURY/core/api/auth"
	"github.com/M-ERCURY/core/api/client"
	"github.com/M-ERCURY/core/api/signer"
	"github.com/M-ERCURY/core/api/texturl"
	"github.com/M-ERCURY/core/cli"
	"github.com/M-ERCURY/core/cli/fsdir"
	"github.com/M-ERCURY/relay/filenames"
	"github.com/M-ERCURY/relay/relaycfg"
	"github.com/M-ERCURY/relay/relaylib"
)

func Cmd() *cli.Subcmd {
	c := relaycfg.Defaults()

	fs := flag.NewFlagSet("balance", flag.ExitOnError)
	contract := fs.String("contract", "", "Service contract URL")

	run := func(fm fsdir.T) {
		err := fm.Get(&c, filenames.Config)

		if err != nil {
			log.Fatal(err)
		}

		scs := []*texturl.URL{}

		type Results struct {
			Contract   string           `json:"contract"`
			Balance    *json.RawMessage `json:"balance"`
			Tokens     *json.RawMessage `json:"tokens"`
			Withdrawal *json.RawMessage `json:"withdrawal"`
		}

		switch {
		case *contract == "":
			// all contracts
			log.Println("no contract specified, listing balance for all contracts...")

			for u, _ := range c.Contracts {
				scs = append(scs, &u)
			}
		default:
			scs = []*texturl.URL{texturl.URLMustParse(*contract)}
		}

		privkey, err := cli.LoadKey(fm, filenames.Seed)

		if err != nil {
			log.Fatal(err)
		}

		var (
			s  = signer.New(privkey)
			cl = client.New(s, auth.Relay)
		)

		for _, u := range scs {
			sc := u.String()

			b, err := relaylib.Get(cl, sc+"/payout/balance")

			if err != nil {
				log.Fatalf("error while getting balance: %s", err)
			}

			ts, err := relaylib.Get(cl, sc+"/payout/tokens")

			if err != nil {
				log.Fatalf("error while getting accumulated tokens: %s", err)
			}

			w, err := relaylib.Get(cl, sc+"/payout/withdrawals")

			if err != nil {
				log.Fatalf("error while getting withdrawal: %s", err)
			}

			r := Results{
				Contract:   sc,
				Balance:    b,
				Tokens:     ts,
				Withdrawal: w,
			}

			data, err := json.MarshalIndent(r, "", "    ")

			if err != nil {
				log.Fatal(err)
			}

			fmt.Println(string(data))
		}
	}

	r := &cli.Subcmd{
		FlagSet: fs,
		Desc:    "Show balance, pending sharetokens and last withdrawal",
		Run:     run,
	}

	return r
}
