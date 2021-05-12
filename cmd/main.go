package main

import (
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/pkg/errors"
	clientset "github.com/tilt-dev/tilt-ci-status/pkg/clientset/versioned"
	"github.com/tilt-dev/tilt-ci-status/pkg/config"
	"github.com/tilt-dev/tilt/pkg/apis/core/v1alpha1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func parseArgs() (resource string, args []string, err error) {
	if len(os.Args) < 3 {
		return "", nil,
			fmt.Errorf("not enough args (need resource name and at minimum one arg to pass to its Cmd). Got args: %v", os.Args[1:])
	}
	return os.Args[1], os.Args[2:], nil
}
func run() error {
	// flag.Parse()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	resource, newArgs, err := parseArgs()
	if err != nil {
		return err
	}

	cmd, err := cmdForResource(ctx, resource)
	if err != nil {
		return err
	}

	fmt.Printf("For resource %s, found command:\n", resource)
	fmt.Printf("\t%s\n", strings.Join(cmd.Spec.Args, " "))
	fmt.Printf("will call with args: %s\n", strings.Join(newArgs, " "))
	return nil
}

func cmdForResource(ctx context.Context, resource string) (v1alpha1.Cmd, error) {
	cmds, err := allCmdsByResource(ctx)
	if err != nil {
		return v1alpha1.Cmd{}, errors.Wrap(err, "getting all Cmds")
	}
	cmd, ok := cmds[resource]
	if !ok {
		var allResources []string
		for resource := range cmds {
			allResources = append(allResources, resource)
		}
		return v1alpha1.Cmd{}, fmt.Errorf("no Cmd found for resource %q (found Cmds for the following resources: %s)",
			resource, strings.Join(allResources, ", "))
	}
	return cmd, nil
}
func allCmdsByResource(ctx context.Context) (map[string]v1alpha1.Cmd, error) {
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
		ret[cmd.ObjectMeta.Annotations[v1alpha1.AnnotationManifest]] = cmd
	}
	return ret, nil
}

func main() {
	err := run()
	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "%v\n", err)
		os.Exit(1)
	}
}
