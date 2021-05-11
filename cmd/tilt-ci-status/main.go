package main

import (
  "context"
  "flag"
  "fmt"
  "github.com/pkg/errors"
  clientset "github.com/tilt-dev/tilt-ci-status/pkg/clientset/versioned"
  "github.com/tilt-dev/tilt-ci-status/pkg/config"
  "github.com/tilt-dev/tilt/pkg/apis/core/v1alpha1"
  v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
  watch "k8s.io/apimachinery/pkg/watch"
  "os"
  "sync"
  "time"
)

var ownResourceName = flag.String("resourcename", "", "name of the resource running this program")
var logInterval = flag.Duration("loginterval", 10 * time.Second, "how often to log tilt-ci status")

type WaitingResource struct {
  Name string
  Reason string
}

func waitingResourcesFromSession(session *v1alpha1.Session) []WaitingResource {
  var ret []WaitingResource
  if session.Status.Done {
    return nil
  }
  for _, t := range session.Status.Targets {
    if t.State.Waiting != nil {
      ret = append(ret, WaitingResource{t.Name, t.State.Waiting.WaitReason})
    } else if t.Type == v1alpha1.TargetTypeJob && t.State.Active != nil {
      ret = append(ret, WaitingResource{t.Name, "executing"})
    } else if t.Type == v1alpha1.TargetTypeServer && t.State.Active != nil && !t.State.Active.Ready {
      ret = append(ret, WaitingResource{t.Name, "waiting for server to be ready"})
    }
  }

  return ret
}

type waitingResourcesWatcher struct {
  wrs []WaitingResource
  mu sync.Mutex
}

func (wrw *waitingResourcesWatcher) subscribe(w watch.Interface, cancel context.CancelFunc) {
  for e := range w.ResultChan() {
    session := e.Object.(*v1alpha1.Session)
    wrs := waitingResourcesFromSession(session)
    // filter out the resource that's running this code
    for i, wr := range wrs {
      if wr.Name == *ownResourceName {
        wrs = append(wrs[:i], wrs[i+1:]...)
      }
    }
    if len(wrs) == 0 {
      cancel()
    }
    wrw.mu.Lock()
    wrw.wrs = wrs
    wrw.mu.Unlock()
  }
}

func (wrw *waitingResourcesWatcher) WaitingResources() []WaitingResource {
  wrw.mu.Lock()
  defer wrw.mu.Unlock()
  return append([]WaitingResource{}, wrw.wrs...)
}

func run() error {
  flag.Parse()

  ctx, cancel := context.WithCancel(context.Background())
  defer cancel()

  w, err := watchSessions(ctx)
  if err != nil {
    return err
  }
  defer w.Stop()

  wrw := waitingResourcesWatcher{}

  go wrw.subscribe(w, cancel)

  for {
    select {
    case <-ctx.Done():
      fmt.Println("🚀 tilt ci all resources finished and healthy!")
      return nil
    case <-time.After(*logInterval):
      fmt.Println("⌛ tilt ci is waiting on:")
      for _, wr := range wrw.WaitingResources() {
        fmt.Printf("  %s: %s\n", wr.Name, wr.Reason)
      }
    }
  }
}

func watchSessions(ctx context.Context) (watch.Interface, error) {
  // TODO - how do we handle multiple tilt instances?
  cfg, err := config.NewConfig()
  if err != nil {
    return nil, errors.Wrap(err, "getting tilt api config")
  }

  cli := clientset.NewForConfigOrDie(cfg)

  w, err := cli.TiltV1alpha1().Sessions().Watch(ctx, v1.ListOptions{})
  if err != nil {
    return nil, errors.Wrap(err, "error watching tilt sessions")
  }
  return w, nil
}

func main() {
  err := run()
  if err != nil {
    _, _ = fmt.Fprintf(os.Stderr, "%v", err)
    os.Exit(1)
  }
}
