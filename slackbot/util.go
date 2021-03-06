package slackbot

import (
	"errors"
	"io/ioutil"
	"os"
	"strings"

	"github.com/lestrrat/go-pdebug"
)

func readEnvConfigFile(name, ename string, dst *string) error {
	f := os.Getenv(ename)
	if f == "" {
		return errors.New("environment " + ename + " not found")
	}

	if pdebug.Enabled {
		pdebug.Printf("Using %s from file specified in environment variable %s (%s)", name, ename, f)
	}

	v, err := ioutil.ReadFile(f)
	if err != nil {
		return err
	}

	if pdebug.Enabled {
		pdebug.Printf("Read %d bytes from %s", len(v), f)
	}
	*dst = strings.TrimSpace(string(v))
	return nil
}
