module github.com/tilt-dev/tilt-ci-status

go 1.16

replace (
	github.com/pkg/browser v0.0.0-00010101000000-000000000000 => github.com/pkg/browser v0.0.0-20210115035449-ce105d075bb4

	k8s.io/apimachinery => github.com/tilt-dev/apimachinery v0.20.2-tilt-20210505
)

require (
	github.com/pkg/errors v0.9.1
	github.com/tilt-dev/tilt v0.20.2
	github.com/tilt-dev/wmclient v0.0.0-20201109174454-1839d0355fbc
	k8s.io/apimachinery v0.21.0
	k8s.io/client-go v0.21.0
)
