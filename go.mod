module github.com/jrasell/sherpa

go 1.12

require (
	github.com/armon/go-metrics v0.0.0-20190430140413-ec5e00d3c878
	github.com/gofrs/uuid v3.2.0+incompatible
	github.com/google/btree v1.0.0 // indirect
	github.com/gorilla/mux v1.7.3
	github.com/hashicorp/consul/api v1.4.0
	github.com/hashicorp/go-cleanhttp v0.5.1
	github.com/hashicorp/go-rootcerts v1.0.2
	github.com/hashicorp/golang-lru v0.5.3 // indirect
	github.com/hashicorp/nomad/api v0.0.0-20190508234936-7ba2378a159e
	github.com/influxdata/influxdb v1.7.2
	github.com/influxdata/platform v0.0.0-20181130221400-c51d92363ec2 // indirect
	github.com/liamg/tml v0.2.0
	github.com/mattn/go-isatty v0.0.12
	github.com/oklog/run v1.0.0
	github.com/panjf2000/ants/v2 v2.1.1
	github.com/pkg/errors v0.8.1
	github.com/prometheus/client_golang v1.3.0
	github.com/rs/zerolog v1.14.3
	github.com/ryanuber/columnize v2.1.0+incompatible
	github.com/sean-/sysexits v0.0.0-20171026162210-598690305aaa
	github.com/spf13/cobra v0.0.3
	github.com/spf13/viper v1.3.2
	github.com/stretchr/testify v1.4.0
)

replace github.com/hashicorp/consul v1.4.0 => github.com/hashicorp/consul v1.7.0
