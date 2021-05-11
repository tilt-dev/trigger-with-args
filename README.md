# Description
If Tilt is running as `tilt ci`, periodically logs what ci is waiting on before completing.

# Usage:
1. Clone the repo to somewhere on your drive.
2. Put this near the top of your Tiltfile to ensure it starts early in your `tilt ci` run.
```
load('path/to/tilt-ci-status/Tiltfile', 'report_tilt_ci_status')
report_tilt_ci_status()
```
