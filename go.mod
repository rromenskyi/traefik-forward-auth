module github.com/thomseddon/traefik-forward-auth

go 1.27

require (
	github.com/coreos/go-oidc v2.5.0+incompatible
	github.com/go-jose/go-jose/v4 v4.1.4
	github.com/gorilla/mux v1.8.1
	github.com/sirupsen/logrus v1.10.2
	github.com/stretchr/testify v1.12.1
	github.com/thomseddon/go-flags v1.4.1-0.20190507184247-a3629c504486
	github.com/vulcand/predicate v1.3.0
	golang.org/x/oauth2 v0.37.0
	golang.org/x/text v0.42.0
)

require (
	github.com/google/go-cmp v0.7.0 // indirect
	github.com/gravitational/trace v1.5.4 // indirect
	github.com/pquerna/cachecontrol v0.2.0 // indirect
	go.yaml.in/yaml/v3 v3.0.5 // indirect
	golang.org/x/crypto v0.57.0 // indirect
	golang.org/x/sys v0.48.0 // indirect
	gopkg.in/go-jose/go-jose.v2 v2.6.3 // indirect
)

// Containous forks
replace (
	github.com/abbot/go-http-auth => github.com/containous/go-http-auth v0.4.1-0.20200324110947-a37a7636d23e
	github.com/go-check/check => github.com/containous/check v0.0.0-20170915194414-ca0bf163426a
	github.com/gorilla/mux => github.com/containous/mux v0.0.0-20250523120546-41b6ec3aed59
	github.com/mailgun/minheap => github.com/containous/minheap v0.0.0-20190809180810-6e71eb837595
)
