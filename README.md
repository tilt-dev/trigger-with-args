# Description
If Tilt is running as `tilt ci`, periodically logs what ci is waiting on before completing, e.g.:
```
report_tilt_… │ ⌛ tilt ci is waiting on:
report_tilt_… │   doggos:runtime: waiting-for-pod
report_tilt_… │   doggos:update: executing
report_tilt_… │   emoji:runtime: waiting-for-pod
report_tilt_… │   emoji:update: unknown
report_tilt_… │   fe:runtime: waiting-for-pod
report_tilt_… │   fe:update: executing
```

# Usage:
1. Clone the repo to somewhere on your drive.
2. Put this near the top of your Tiltfile to ensure it starts early in your `tilt ci` run.
```
load('path/to/tilt-ci-status/Tiltfile', 'report_tilt_ci_status')
report_tilt_ci_status()
```
