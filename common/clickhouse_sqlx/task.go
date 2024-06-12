package clickhouse_sqlx

import "strings"

// TaskConfig parameters
type TaskConfig struct {
	Name string

	Topic         string
	ConsumerGroup string

	// Earliest set to true to consume the message from oldest position
	Earliest bool
	Parser   string
	// the csv cloum title if Parser is csv
	CsvFormat []string
	Delimiter string

	TableName       string
	SeriesTableName string

	// AutoSchema will auto fetch the schema from clickhouse
	AutoSchema     bool
	ExcludeColumns []string
	Dims           []struct {
		Name       string
		Type       string
		SourceName string
	} `json:"dims"`
	// DynamicSchema will add columns present in message to clickhouse. Requires AutoSchema be true.
	DynamicSchema struct {
		Enable      bool
		NotNullable bool
		MaxDims     int // the upper limit of dynamic columns number, <=0 means math.MaxInt16. protecting dirty data attack
		// A column is added for new key K if all following conditions are true:
		// - K isn't in ExcludeColumns
		// - number of existing columns doesn't reach MaxDims-1
		// - WhiteList is empty, or K matchs WhiteList
		// - BlackList is empty, or K doesn't match BlackList
		WhiteList string // the regexp of white list
		BlackList string // the regexp of black list
	}
	// additional fields to be appended to each input message, should be a valid json string
	Fields string `json:"fields,omitempty"`
	// PrometheusSchema expects each message is a Prometheus metric(timestamp, value, metric name and a list of labels).
	PrometheusSchema bool
	// fields match PromLabelsBlackList are not considered as labels. Requires PrometheusSchema be true.
	PromLabelsBlackList string // the regexp of black list

	// ShardingKey is the column name to which sharding against
	ShardingKey string `json:"shardingKey,omitempty"`
	// ShardingStripe take effect if the sharding key is numerical
	ShardingStripe uint64 `json:"shardingStripe,omitempty"`

	FlushInterval int     `json:"flushInterval,omitempty"`
	BufferSize    int     `json:"bufferSize,omitempty"`
	TimeZone      string  `json:"timeZone"`
	TimeUnit      float64 `json:"timeUnit"`
}

type GroupConfig struct {
	Name          string
	Topics        []string
	Earliest      bool
	FlushInterval int
	BufferSize    int
	Configs       map[string]*TaskConfig
}

type Assignment struct {
	Version   int
	UpdatedAt int64               // timestamp when created
	UpdatedBy string              // leader instance
	Map       map[string][]string // map instance to a list of task_name
}

func (cfg *Assignment) IsAssigned(instance, task string) (assigned bool) {
	if taskNames, ok := cfg.Map[instance]; ok {
		for _, taskName := range taskNames {
			if taskName == task {
				assigned = true
				return
			}
		}
	}
	return
}

func readConfig(config string) map[string]string {
	configMap := make(map[string]string)
	config = strings.TrimSuffix(config, ";")
	fields := strings.Split(config, " ")
	for _, field := range fields {
		if strings.Contains(field, "=") {
			key := strings.Split(field, "=")[0]
			value := strings.Split(field, "=")[1]
			value = strings.Trim(value, "\"")
			configMap[key] = value
		}
	}
	return configMap
}
