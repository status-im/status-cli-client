module github.com/status-im/status-cli-client

go 1.13

require (
	github.com/allegro/bigcache v1.2.1 // indirect
	github.com/ethereum/go-ethereum v1.9.5
	github.com/fjl/memsize v0.0.0-20190710130421-bcb5799ab5e5 // indirect
	github.com/google/uuid v1.1.1
	github.com/karalabe/hid v1.0.0 // indirect
	github.com/pkg/errors v0.8.1
	github.com/rjeczalik/notify v0.9.2 // indirect
	github.com/status-im/migrate/v4 v4.6.2-status.2
	github.com/status-im/status-go v0.34.0-beta.3
	github.com/status-im/status-protocol-go v0.4.2
	gopkg.in/olebedev/go-duktape.v3 v3.0.0-20190709231704-1e4459ed25ff // indirect
)

replace github.com/ethereum/go-ethereum v1.9.5 => github.com/status-im/go-ethereum v1.9.5-status.4

replace github.com/NaySoftware/go-fcm => github.com/status-im/go-fcm v1.0.0-status
