// The release version is defined here.
package version

import (
	"fmt"
	"log"

	"github.com/M-ERCURY/core/api/auth"
	"github.com/M-ERCURY/core/api/client"
	"github.com/M-ERCURY/core/api/consume"
	"github.com/M-ERCURY/core/api/relayentry"
	"github.com/M-ERCURY/core/api/signer"
	"github.com/M-ERCURY/core/api/texturl"
	"github.com/M-ERCURY/core/cli"
	"github.com/M-ERCURY/core/cli/fsdir"
	"github.com/M-ERCURY/core/cli/upgrade"
	"github.com/M-ERCURY/relay/filenames"
	"github.com/M-ERCURY/relay/relaycfg"
	"github.com/blang/semver"
)

// old name compat
var GITREV string = "1.0.0"

// VERSION_STRING is the current version string, set by the linker via go build
// -X flag.
var VERSION_STRING = GITREV

// VERSION is the semver version struct of VERSION_STRING.
var VERSION = semver.MustParse(VERSION_STRING)

// Post-rollback hook for rollbackcmd.
func PostRollbackHook(f fsdir.T) (err error) {
	// get old binary back up asap, try 3 times
	log.Printf("starting old mercury-relay...")
	for i := 0; i < 3; i++ {
		if err = cli.RunChild(f.Path("mercury-relay"), "start"); err == nil {
			// ok
			return nil
		} else {
			err = fmt.Errorf("FAILED to start old mercury-relay, try %d: %s", i, err)
		}
	}
	// hard fail
	return fmt.Errorf("failed to bring old binary up -- there is no mercury-relay running! %s", err)
}

// MIGRATIONS is the slice of versioned migrations.
var MIGRATIONS = []*upgrade.Migration{}

// LatestChannelVersion is a special function for mercury-relay which will
// obtain the latest version supported by the currently configured update
// channel from the directory.
func LatestChannelVersion(f fsdir.T) (semver.Version, error) {
	c := relaycfg.Defaults()
	// err := fm.Get(&c, filenames.Config)
	if err := f.Get(&c, filenames.Config); err != nil {
		return semver.Version{}, err
	}
	if err := c.Validate(); err != nil {
		return semver.Version{}, err
	}
	privkey, err := cli.LoadKey(f, filenames.Seed)
	if err != nil {
		return semver.Version{}, err
	}
	cl := client.New(signer.New(privkey), auth.Relay)
	// NOTE: this depends on there being only 1 contract
	var (
		scurl texturl.URL
		sccfg *relayentry.T
	)
	for k, v := range c.Contracts {
		scurl, sccfg = k, v
		break
	}
	dinfo, err := consume.DirectoryInfo(cl, &scurl)
	if err != nil {
		return semver.Version{}, err
	}
	v, ok := dinfo.Channels[sccfg.Channel]
	if !ok {
		return v, fmt.Errorf("no version for channel '%s' is provided by directory", sccfg.Channel)
	}
	return v, nil
}
