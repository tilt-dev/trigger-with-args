package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/pkg/errors"
	clientset "github.com/tilt-dev/tilt-ci-status/pkg/clientset/versioned"
	"github.com/tilt-dev/tilt-ci-status/pkg/config"
	"github.com/tilt-dev/tilt/pkg/apis/core/v1alpha1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func run() error {
	flag.Parse()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	fmt.Println("hello world")
	fmt.Println("here are all cmds:")
	cmds, err := allCmdsByName(ctx)
	if err != nil {
		return err
	}

	for n, c := range cmds {
		fmt.Println(n)
		fmt.Printf("\t%s\n", strings.Join(c.Spec.Args, " "))
		fmt.Println("---")
	}
	return nil
}

func allCmdsByName(ctx context.Context) (map[string]v1alpha1.Cmd, error) {
	// TODO - how do we handle multiple tilt instances?
	cfg, err := config.NewConfig()
	if err != nil {
		return nil, errors.Wrap(err, "getting tilt api config")
	}

	cli := clientset.NewForConfigOrDie(cfg)

	cmds, err := cli.TiltV1alpha1().Cmds().List(ctx, v1.ListOptions{})
	if err != nil {
		return nil, errors.Wrap(err, "error watching tilt sessions")
	}

	ret := make(map[string]v1alpha1.Cmd, len(cmds.Items))
	for _, cmd := range cmds.Items {
		ret[cmd.Name] = cmd
	}
	return ret, nil
}

func main() {
	err := run()
	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "%v", err)
		os.Exit(1)
	}
}
