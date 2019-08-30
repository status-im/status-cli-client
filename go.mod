module github.com/status-im/status-cli-client

go 1.12

require (
	github.com/allegro/bigcache v1.2.1 // indirect
	github.com/ethereum/go-ethereum v1.8.27
	github.com/fjl/memsize v0.0.0-20190710130421-bcb5799ab5e5 // indirect
	github.com/golang-migrate/migrate/v4 v4.6.1 // indirect
	github.com/google/uuid v1.1.1
	github.com/karalabe/hid v1.0.0 // indirect
	github.com/pkg/errors v0.8.1
	github.com/rjeczalik/notify v0.9.2 // indirect
	github.com/status-im/migrate/v4 v4.3.1-status.0.20190822050738-a9d340ec8fb7 // indirect
	github.com/status-im/status-go v0.30.1-beta.2.0.20190828210454-4761179cc0673a0bb8fe0d33ca6a035ca62c28c7
	github.com/status-im/status-protocol-go v0.2.0
	github.com/vacp2p/mvds v0.0.21-0.20190824144946-3233b2308076 // indirect
	gopkg.in/olebedev/go-duktape.v3 v3.0.0-20190709231704-1e4459ed25ff // indirect
)

replace github.com/ethereum/go-ethereum v1.8.27 => github.com/status-im/go-ethereum v1.8.27-status.5
