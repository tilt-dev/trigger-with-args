# cmd-with-args
A hackweek proof-of-concept for triggering a local resource with user-provided args.

## Why?
We know that Tilt users dig having a menu of common operations (in the form of trigger-able Local Resources) at their fingertips in the Web UI. One comment we've heard is "it'd be great to be able to call these with arbitrary args"--e.g. maybe there's a `seed-db.sh` script that targets a different database depending on the arg passed.

This CLI is a proof-of-concept that with the new API, we can swap out the args of an existing Cmd with user-specified ones.

## Usage:
```
go run cmd/main.go resourcename newarg1 [newarg2...]
```
* specify `resourcename`: the name of the local resource to operate on
* provide one or more space-separated args to be passed to that resource's ServeCmd

Try running `tilt up` in this repo and modifying the args for the `change-me` resource:
```
go run cmd/main.go changeme ✨
```

### What's Happening/Limitations
Currently this CLI can only modify `serve_cmd`'s (b/c those are the only ones implemented as API objects today). It will probably be most useful operating on `update_cmd`'s (at least for the "configurable task menu" vision of the world). In the `update_cmd` case, the logic changes somewhat, as there's no long-lived Cmd to find, modify, and upsert--instead, we would probably find the resource, read its `update_cmd`, modify the args as needed, and dispatch that Cmd ourselves.

The way we decide what args of the existing Cmd to keep vs. discard is clunky. Currently we keep only the first arg, and replace the rest. If we're okay with the first run of the Cmd failing, maybe the Cmd should be initially as just the base command (+ args we always want), and the first time we run this tool against it, we append our args to the end, and store the number of added args in an annotation; if we modify that same Cmd again, we chop off only args we added via script (as noted in annotation). (This gets slightly easier if we're dispatching `update_cmd`'s ourselves, rather than modifying an existing `serve_cmd`--we don't need to track state of "how many args did we add last time?".)

Finally, this tool is kinda clunky as a CLI--if we want to invest in this workflow in future, it should be in the UI (e.g. user clicks a button to "run with args", gets a text box where they can specify args and hit "enter".)
