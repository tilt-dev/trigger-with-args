def report_tilt_ci_status(log_interval="10s"):
  if config.tilt_subcommand == 'ci':
    local_resource('report_tilt_ci_status', serve_cmd='cd "%s" && go run ./cmd/tilt-ci-status --resourcename report_tilt_ci_status --loginterval "%s"' % (os.path.dirname(__file__), log_interval))
