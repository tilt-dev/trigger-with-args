package main

import (
	"context"
	"fmt"
	"os"

	"github.com/tilt-dev/tilt-ci-status/pkg/editor"
)

func parseArgs() (resource string, args []string, err error) {
	if len(os.Args) < 3 {
		return "", nil,
			fmt.Errorf("not enough args (need resource name and at minimum one arg to pass to its Cmd). Got args: %v", os.Args[1:])
	}
	return os.Args[1], os.Args[2:], nil
}
func run() error {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	resource, newArgs, err := parseArgs()
	if err != nil {
		return err
	}

	ce, err := editor.NewCmdEditor()
	if err != nil {
		return err
	}

	return ce.CallCmdForResourceWithArgs(ctx, resource, newArgs)
}

func main() {
	err := run()
	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "%v\n", err)
		os.Exit(1)
	}
}
