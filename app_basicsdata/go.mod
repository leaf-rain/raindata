module github.com/leaf-rain/raindata/app_basicsdata

go 1.22.0

replace github.com/leaf-rain/raindata/common => ../common

replace github.com/leaf-rain/raindata/app_report => ../app_report

require (
	github.com/fsnotify/fsnotify v1.7.0
	github.com/google/wire v0.6.0
	github.com/leaf-rain/raindata/app_report v0.0.0-00010101000000-000000000000
	github.com/leaf-rain/raindata/common v0.0.0-00010101000000-000000000000
	github.com/mitchellh/mapstructure v1.5.0
	github.com/spf13/viper v1.18.2
	go.etcd.io/etcd/client/v3 v3.5.13
	go.uber.org/zap v1.27.0
	google.golang.org/grpc v1.63.2
	google.golang.org/protobuf v1.34.1
)

require (
	github.com/ClickHouse/ch-go v0.61.5 // indirect
	github.com/ClickHouse/clickhouse-go/v2 v2.25.0 // indirect
	github.com/andybalholm/brotli v1.1.0 // indirect
	github.com/coreos/go-semver v0.3.0 // indirect
	github.com/coreos/go-systemd/v22 v22.3.2 // indirect
	github.com/go-faster/city v1.0.1 // indirect
	github.com/go-faster/errors v0.7.1 // indirect
	github.com/gogo/protobuf v1.3.2 // indirect
	github.com/golang/protobuf v1.5.4 // indirect
	github.com/google/uuid v1.6.0 // indirect
	github.com/hashicorp/hcl v1.0.0 // indirect
	github.com/jmoiron/sqlx v1.2.0 // indirect
	github.com/klauspost/compress v1.17.8 // indirect
	github.com/leaf-rain/fastjson v0.0.0-20240605084444-e45ca8b5db75 // indirect
	github.com/magiconair/properties v1.8.7 // indirect
	github.com/natefinch/lumberjack v2.0.0+incompatible // indirect
	github.com/paulmach/orb v0.11.1 // indirect
	github.com/pelletier/go-toml/v2 v2.1.0 // indirect
	github.com/pierrec/lz4/v4 v4.1.21 // indirect
	github.com/pkg/errors v0.9.1 // indirect
	github.com/sagikazarmark/locafero v0.4.0 // indirect
	github.com/sagikazarmark/slog-shim v0.1.0 // indirect
	github.com/segmentio/asm v1.2.0 // indirect
	github.com/shopspring/decimal v1.4.0 // indirect
	github.com/sourcegraph/conc v0.3.0 // indirect
	github.com/spf13/afero v1.11.0 // indirect
	github.com/spf13/cast v1.6.0 // indirect
	github.com/spf13/pflag v1.0.5 // indirect
	github.com/subosito/gotenv v1.6.0 // indirect
	go.etcd.io/etcd/api/v3 v3.5.13 // indirect
	go.etcd.io/etcd/client/pkg/v3 v3.5.13 // indirect
	go.opentelemetry.io/otel v1.27.0 // indirect
	go.opentelemetry.io/otel/trace v1.27.0 // indirect
	go.uber.org/multierr v1.11.0 // indirect
	golang.org/x/exp v0.0.0-20230905200255-921286631fa9 // indirect
	golang.org/x/net v0.22.0 // indirect
	golang.org/x/sys v0.20.0 // indirect
	golang.org/x/text v0.14.0 // indirect
	google.golang.org/genproto v0.0.0-20240227224415-6ceb2ff114de // indirect
	google.golang.org/genproto/googleapis/api v0.0.0-20240227224415-6ceb2ff114de // indirect
	google.golang.org/genproto/googleapis/rpc v0.0.0-20240227224415-6ceb2ff114de // indirect
	gopkg.in/ini.v1 v1.67.0 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)
