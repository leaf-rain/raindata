module github.com/leaf-rain/raindata/app_report

go 1.22.0

replace github.com/leaf-rain/raindata/common => ../common

replace github.com/leaf-rain/raindata/app_basicsdata => ../app_basicsdata

require (
	github.com/fsnotify/fsnotify v1.7.0
	github.com/google/wire v0.6.0
	github.com/leaf-rain/fastjson v0.0.0-20240605084444-e45ca8b5db75
	github.com/leaf-rain/raindata/app_basicsdata v0.0.0-00010101000000-000000000000
	github.com/leaf-rain/raindata/common v0.0.0-00010101000000-000000000000
	github.com/natefinch/lumberjack v2.0.0+incompatible
	github.com/spf13/viper v1.18.2
	go.uber.org/zap v1.27.0
	google.golang.org/grpc v1.63.2
	google.golang.org/protobuf v1.34.1
)

require (
	github.com/BurntSushi/toml v1.3.2 // indirect
	github.com/hashicorp/hcl v1.0.0 // indirect
	github.com/magiconair/properties v1.8.7 // indirect
	github.com/mitchellh/mapstructure v1.5.0 // indirect
	github.com/pelletier/go-toml/v2 v2.1.0 // indirect
	github.com/sagikazarmark/locafero v0.4.0 // indirect
	github.com/sagikazarmark/slog-shim v0.1.0 // indirect
	github.com/sourcegraph/conc v0.3.0 // indirect
	github.com/spf13/afero v1.11.0 // indirect
	github.com/spf13/cast v1.6.0 // indirect
	github.com/spf13/pflag v1.0.5 // indirect
	github.com/subosito/gotenv v1.6.0 // indirect
	go.uber.org/multierr v1.11.0 // indirect
	golang.org/x/exp v0.0.0-20230905200255-921286631fa9 // indirect
	golang.org/x/net v0.22.0 // indirect
	golang.org/x/sys v0.20.0 // indirect
	golang.org/x/text v0.14.0 // indirect
	google.golang.org/genproto/googleapis/rpc v0.0.0-20240227224415-6ceb2ff114de // indirect
	gopkg.in/ini.v1 v1.67.0 // indirect
	gopkg.in/natefinch/lumberjack.v2 v2.2.1 // indirect
	gopkg.in/yaml.v2 v2.4.0 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)
