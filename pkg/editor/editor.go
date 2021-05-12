package editor

import (
	"context"
	"fmt"
	"strings"

	"github.com/pkg/errors"
	clientset "github.com/tilt-dev/tilt-ci-status/pkg/clientset/versioned"
	"github.com/tilt-dev/tilt-ci-status/pkg/config"
	"github.com/tilt-dev/tilt/pkg/apis/core/v1alpha1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type CmdEditor struct {
	Cli *clientset.Clientset
}

func NewCmdEditor() (CmdEditor, error) {
	// TODO - how do we handle multiple tilt instances?
	cfg, err := config.NewConfig()
	if err != nil {
		return CmdEditor{}, errors.Wrap(err, "getting tilt api config")
	}

	return CmdEditor{
		Cli: clientset.NewForConfigOrDie(cfg),
	}, nil
}

func (ce *CmdEditor) CallCmdForResourceWithArgs(ctx context.Context, resource string, newArgs []string) error {
	oldCmd, err := ce.CmdForResource(ctx, resource)
	if err != nil {
		return errors.Wrap(err, "getting cmd for resource")
	}
	fmt.Printf("For resource %s, found command:\n", resource)
	fmt.Printf("\t%s\n", strings.Join(oldCmd.Spec.Args, " "))
	fmt.Printf("will call with args: %s\n", strings.Join(newArgs, " "))

	newCmd, err := ce.NewCmdWithArgs(oldCmd, newArgs)
	if err != nil {
		return errors.Wrap(err, "making new command with args")
	}

	return ce.UpsertCmd(ctx, newCmd)

}
func (ce *CmdEditor) CmdForResource(ctx context.Context, resource string) (v1alpha1.Cmd, error) {
	cmds, err := ce.AllCmdsByResource(ctx)
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
func (ce *CmdEditor) AllCmdsByResource(ctx context.Context) (map[string]v1alpha1.Cmd, error) {
	cmds, err := ce.Cli.TiltV1alpha1().Cmds().List(ctx, v1.ListOptions{})
	if err != nil {
		return nil, errors.Wrap(err, "error watching tilt sessions")
	}

	ret := make(map[string]v1alpha1.Cmd, len(cmds.Items))
	for _, cmd := range cmds.Items {
		ret[cmd.ObjectMeta.Annotations[v1alpha1.AnnotationManifest]] = cmd
	}
	return ret, nil
}

func (ce *CmdEditor) NewCmdWithArgs(cmd v1alpha1.Cmd, args []string) (*v1alpha1.Cmd, error) {
	curArgs := cmd.Spec.Args
	if len(curArgs) == 0 {
		return nil, fmt.Errorf("cannot run this command with new args: need at least one existing arg (found none)")
	}
	if curArgs[0] == "/bin/sh" || curArgs[0] == "/bin/bash" || curArgs[0] == "sh" {
		return nil, fmt.Errorf("cannot run this command with new args: must be in argv format, found a bash command (in your Tiltfile, pass the command as `List[str]`, not `str`)")
	}

	// Keep the first arg and chop off+replace all the rest. There's probably a
	// smarter way to do this--maybe indicating via annotation how many args to chop off?
	// If we're okay with the first run of the Cmd failing, maybe the Cmd should be
	// specified as just the base command+args, and the first time we run this script
	// against it, we append our args to the end, and store the number of added
	// args in an annotation; if we modify that same Cmd again, we chop off only
	// args we added via script (as noted in annotation).
	newArgs := append(curArgs[:1], args...)

	cp := cmd.DeepCopy()
	cp.Spec.Args = newArgs
	return cp, nil
}

func (ce *CmdEditor) UpsertCmd(ctx context.Context, cmd *v1alpha1.Cmd) error {
	cmdCli := ce.Cli.TiltV1alpha1().Cmds()

	fmt.Printf("🤖 deleting previous version of command %q\n", cmd.Name)
	err := cmdCli.Delete(ctx, cmd.Name, v1.DeleteOptions{})
	if err != nil {
		return errors.Wrapf(err, "deleting existing cmd %q", cmd.Name)
	}

	fmt.Printf("🤖 creating new version with args: %v\n", cmd.Spec.Args)
	_, err = cmdCli.Create(ctx, cmd, v1.CreateOptions{})
	if err != nil {
		return errors.Wrapf(err, "recreating cmd %q", cmd.Name)
	}

	return nil
}
