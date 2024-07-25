package entity

import (
	"context"
	pb_metadata "github.com/leaf-rain/raindata/app_basicsdata/api/grpc"
	"github.com/leaf-rain/raindata/app_basicsdata/internal/infrastructure/config"
	commonConfig "github.com/leaf-rain/raindata/common/config"
	"github.com/leaf-rain/raindata/common/consts"
	"github.com/leaf-rain/raindata/common/rclickhouse"
	"go.uber.org/zap"
	"sort"
	"strconv"
	"strings"
)

type Metadata struct {
	ctx           context.Context
	fields        *pb_metadata.MetadataRequest
	dynamicConfig commonConfig.ConfigInterface
	ck            *rclickhouse.Conn
	logger        *zap.Logger
	cfg           *config.Config
}

func (repo *Repository) NewMetadata(ctx context.Context, logger *zap.Logger, fields *pb_metadata.MetadataRequest, ck *rclickhouse.Conn) *Metadata {
	metadata := &Metadata{
		ctx:           ctx,
		logger:        logger,
		fields:        fields,
		dynamicConfig: repo.dynamicConfig,
		cfg:           repo.cfg,
		ck:            ck,
	}
	return metadata
}

func sortedFields(fields []*pb_metadata.Field) {
	sort.SliceIsSorted(fields, func(i, j int) bool {
		return fields[i].Name < fields[j].Name
	})
}

func GetMetadataKey(fields []*pb_metadata.Field) string {
	sortedFields(fields)
	// 避免生成的key顺序不一样
	var builder strings.Builder
	for _, item := range fields {
		builder.WriteString("_")
		builder.WriteString(item.Name)
		builder.WriteString("_")
		builder.WriteString(item.Type)
	}
	return builder.String()
}

func (repo *Metadata) LockKey() string {
	var builder strings.Builder
	builder.WriteString(consts.ETCDKeyPre)
	builder.WriteString("/")
	builder.WriteString(repo.fields.EventName)
	return builder.String()
}

func (repo *Metadata) key(k string) string {
	var builder strings.Builder
	builder.WriteString(repo.cfg.MetadataPath)
	builder.WriteString("/")
	builder.WriteString(k)
	return builder.String()
}

func (repo *Metadata) SetFields(data *pb_metadata.MetadataRequest) {
	repo.fields = data
}

func (repo *Metadata) GetMetadata() (*pb_metadata.MetadataResponse, error) {
	resultI := repo.dynamicConfig.GetByKey(repo.fields.App, repo.fields.EventName)
	var result *pb_metadata.MetadataResponse
	if resultI != nil {
		result, _ = resultI.(*pb_metadata.MetadataResponse)
	}
	return result, nil
}

type fileInfo struct {
	Name              string `db:"name"`
	Type              string `db:"type"`
	DefaultType       string `db:"default_type"`
	DefaultExpression string `db:"default_expression"`
	Comment           string `db:"comment"`
	CodecExpression   string `db:"codec_expression"`
	TtlExpression     string `db:"ttl_expression"`
}

func (repo *Metadata) PutMetadata(ty int) (*pb_metadata.MetadataResponse, error) {
	var result = &pb_metadata.MetadataResponse{
		Metadata: repo.fields.Fields,
	}
	err := repo.dynamicConfig.Storage(commonConfig.ConfigInfo{
		App:  repo.fields.App,
		Name: repo.fields.EventName,
		Info: result,
	})
	if err != nil {
		return nil, err
	}
	var tableName = getEventTableName(repo.fields.App)
	if ty == 1 {
		err = repo.CreateEventTable(repo.fields.App)
		if err != nil {
			return nil, err
		}
	}
	var rows *rclickhouse.Rows
	rows, err = repo.ck.Query("DESCRIBE " + tableName)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var fieldsForDb = make(map[string]string)
	var name, fieldType, defaultType, defaultExpression, comment, codecExpression, ttlExpression string
	for rows.Next() {
		if err = rows.Scan(
			&name,
			&fieldType,
			&defaultType,
			&defaultExpression,
			&comment,
			&codecExpression,
			&ttlExpression,
		); err != nil {
			return nil, err
		}
		fieldsForDb[name] = fieldType
	}
	var ok bool
	for _, item := range repo.fields.Fields {
		if _, ok = fieldsForDb[item.Field]; !ok {
			err = repo.ck.Exec("ALTER TABLE " + tableName + " ADD COLUMN " + item.Field + " " + item.Type)
			if err != nil {
				repo.logger.Error("[PutMetadata] add field failed.", zap.Error(err))
			}
		}
	}
	return result, nil
}

func (repo *Metadata) GroupMetadata(fields []*pb_metadata.Field) *pb_metadata.MetadataResponse {
	if len(repo.fields.Fields) == 0 {
		repo.fields.Fields = []*pb_metadata.Field{
			{
				Name:  "event",
				Type:  "String",
				Field: "event",
			}, {
				Name:  "sort_id",
				Type:  "Int64",
				Field: "sort_id",
			}, {
				Name:  "insert_at",
				Type:  "Int64",
				Field: "insert_at",
			}, {
				Name:  "created_at",
				Type:  "DateTime",
				Field: "created_at",
			},
		}
	}
	var fieldsMap = make(map[string][]*pb_metadata.Field)
	var ok bool
	var length int
	for _, item := range repo.fields.Fields {
		if _, ok = fieldsMap[item.Type]; !ok {
			fieldsMap[item.Type] = make([]*pb_metadata.Field, 0)
		}
		fieldsMap[item.Type] = append(fieldsMap[item.Type], item)
		length += 1
	}
	for _, item := range fields {
		_, ok = fieldsMap[item.Type]
		if ok && IsExistByName(fieldsMap[item.Type], item.Name) {
			continue
		}
		length += 1
		if !ok {
			item.Field = item.Type + "_1"
			fieldsMap[item.Type] = []*pb_metadata.Field{item}
		} else {
			item.Field += item.Type + "_" + strconv.Itoa(len(fieldsMap[item.Type])+1)
			fieldsMap[item.Type] = append(fieldsMap[item.Type], item)
		}
	}
	var result = &pb_metadata.MetadataResponse{
		Metadata: make([]*pb_metadata.Field, length),
	}
	for _, items := range fieldsMap {
		for _, item := range items {
			result.Metadata[length-1] = item
			length -= 1
		}
	}
	repo.fields.Fields = result.Metadata
	return result
}

func IsExistByName(fields []*pb_metadata.Field, name string) bool {
	for i := range fields {
		if fields[i].Name == name {
			return true
		}
	}
	return false
}

func getEventTableName(appId string) string {
	return "app_" + appId
}

func (repo *Metadata) CreateEventTable(appId string) error {
	tableName := getEventTableName(appId)
	sql := `CREATE TABLE IF NOT EXISTS ` + tableName + `(event String, sort_id Int64, insert_at Int64, created_at DateTime) ENGINE = MergeTree ORDER BY (event, sort_id, created_at) PARTITION BY toYYYYMM(toMonday(created_at)) SETTINGS index_granularity = 8192;`
	err := repo.ck.Exec(sql)
	return err
}
