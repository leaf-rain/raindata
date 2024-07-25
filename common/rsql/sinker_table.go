package rsql

import (
	"bytes"
	"context"
	"database/sql"
	"fmt"
	"github.com/tidwall/gjson"
	"go.uber.org/zap"
	"go.uber.org/zap/buffer"
	"io"
	"net/http"
	"net/url"
	"os"
	"sort"
	"strconv"
	"strings"
	"sync"
	"sync/atomic"
	"time"
)

const (
	StateRunning uint32 = 0
	StateStopped uint32 = 1
)

var bufferPool = buffer.NewPool()

// Decoder decodes the contents of b into v.
// It's primarily used for decoding contents of a file into a map[string]any.
type Decoder interface {
	Decode(b []byte, v map[string]any) error
}

type SinkerTableConfig struct {
	TableName      string           `json:"table_name,omitempty"`
	Database       string           `json:"database,omitempty"`
	FlushInterval  int              `json:"flush_interval,omitempty"`
	BufferSize     int              `json:"buffer_size,omitempty"`
	BaseColumn     []ColumnWithType `json:"base_column,omitempty"`
	TableType      int              `json:"table_type"`
	PrimaryKey     string           `json:"primary_key"`
	DistributedKey string           `json:"distributed_key"`
	OrderByKey     string           `json:"order_by_key"`

	Username   string `json:"username,omitempty"`
	Password   string `json:"password,omitempty"`
	FeHost     string `json:"fe_host,omitempty"`
	FeHttpPort string `json:"fe_http_port,omitempty"`
}

// SinkerTable object maintains number of task for each partition
type SinkerTable struct {
	logger     *zap.Logger
	db         *sql.DB
	ctx        context.Context
	cancel     context.CancelFunc
	fetchCH    chan Fetch
	exitCh     chan struct{}
	processWg  sync.WaitGroup
	mux        sync.RWMutex
	commitDone *sync.Cond
	state      atomic.Uint32
	parser     Parser
	fieldMap   *sync.Map // key:columnName, value:*ColumnWithType
	// table相关
	sinkerTableConfig *SinkerTableConfig
	fd                *os.File
	needComma         bool // 需要逗号
}

// NewSinkerTable get an instance of sinker with the task list
func NewSinkerTable(ctx context.Context, db *sql.DB, logger *zap.Logger, sinkerTableConfig *SinkerTableConfig, parserTy string) (*SinkerTable, error) {
	s := &SinkerTable{
		logger:            logger,
		db:                db,
		ctx:               ctx,
		fetchCH:           make(chan Fetch, 1000),
		exitCh:            make(chan struct{}),
		sinkerTableConfig: sinkerTableConfig,
		fieldMap:          new(sync.Map),
	}
	s.state.Store(StateStopped)
	s.commitDone = sync.NewCond(&s.mux)
	s.parser, _ = NewParse(ParserConfig{
		Ty:     parserTy,
		Logger: logger,
	})
	// 存储默认字段信息
	for i := range sinkerTableConfig.BaseColumn {
		s.fieldMap.Store(sinkerTableConfig.BaseColumn[i].Name, &sinkerTableConfig.BaseColumn[i])
	}
	err := s.flushFields()
	if err != nil {
		return nil, err
	}
	err = s.wal()
	if err != nil {
		return nil, err
	}
	if s.sinkerTableConfig.FlushInterval <= 0 {
		s.sinkerTableConfig.FlushInterval = 1
	}
	s.fd, err = os.OpenFile(sinkerTableConfig.TableName+".wal", os.O_CREATE|os.O_RDWR|os.O_APPEND, 0666)
	if err != nil {
		return nil, err
	}
	var tail = getTailChar(s.fd)
	if tail == "" {
		s.fd.Write([]byte("["))
	} else if tail == "}" {
		s.fd.Write([]byte(","))
	}
	return s, nil
}

func getTailChar(fd *os.File) string {
	var tail = make([]byte, 1)
	tailCount, _ := fd.Read(tail)
	if tailCount == 0 {
		return ""
	}
	return string(tail)
}

func (c *SinkerTable) wal() error {
	// todo:wal恢复未处理的数据
	return nil
}

func (c *SinkerTable) getAllFields() []string {
	var allColumns []string
	c.fieldMap.Range(func(key, value any) bool {
		allColumns = append(allColumns, key.(string))
		return true
	})
	sort.Slice(allColumns, func(i, j int) bool {
		return allColumns[i] < allColumns[j]
	})
	return allColumns
}

type sqlFields struct {
	Field string `db:"Field"`
	Type  string `db:"Type"`
	Null  string `db:"Null"`
}

func (c *SinkerTable) flushFields() error {
	var query = getTableColumns(c.sinkerTableConfig.Database, c.sinkerTableConfig.TableName)
	var err error
	var rs *sql.Rows
	rs, err = c.db.Query(query)
	if err == nil {
		var columns = make(map[string]*ColumnWithType)
		defer rs.Close()
		var name, typ, nullable, key, extra string
		var de sql.NullString
		for rs.Next() {
			if err = rs.Scan(&name, &typ, &nullable, &key, &de, &extra); err != nil {
				return err
			}
			ct := &ColumnWithType{Name: name, Type: WhichType(typ, nullable == "YES")}
			c.fieldMap.Store(name, ct)
			columns[name] = ct
		}
		var addColumns []string
		var allColumns []string
		c.fieldMap.Range(func(key, value any) bool {
			if _, ok := columns[key.(string)]; !ok {
				addColumns = append(addColumns, addTableColumns(c.sinkerTableConfig.Database, c.sinkerTableConfig.TableName, key.(string), value.(*ColumnWithType).Type.ToString(), value.(*ColumnWithType).Type.Nullable))
			}
			allColumns = append(allColumns, key.(string))
			return true
		})
		sort.Slice(allColumns, func(i, j int) bool {
			return allColumns[i] < allColumns[j]
		})
		if len(addColumns) > 0 {
			for i := range addColumns {
				_, err = c.db.Exec(addColumns[i])
				if err != nil && !strings.Contains(err.Error(), "already exists") {
					c.logger.Error(c.sinkerTableConfig.TableName+" add table columns error", zap.Error(err), zap.String("sql", addColumns[i]))
					return err
				}
			}
		}
	} else if err != nil && strings.Contains(err.Error(), "Unknown table") {
		query = createTable(c.sinkerTableConfig.TableType, c.sinkerTableConfig.Database, c.sinkerTableConfig.TableName, c.sinkerTableConfig.PrimaryKey, c.sinkerTableConfig.DistributedKey, c.sinkerTableConfig.OrderByKey, c.fieldMap)
		_, err = c.db.Exec(query)
		if err != nil {
			return err
		}
		return c.flushFields()
	}
	return err
}

func (c *SinkerTable) Start() {
	if c.state.Load() == StateRunning {
		return
	}
	c.state.Store(StateRunning)
	go c.processFetch()
}

func (c *SinkerTable) Stop() {
	if c.state.Load() == StateStopped {
		return
	}
	c.state.Store(StateStopped)
	c.exitCh <- struct{}{}
	c.processWg.Wait()
}

func (c *SinkerTable) Restart() {
	c.Stop()
	c.Start()
}

func (c *SinkerTable) WriterMsg(msg FetchSingle) {
	c.fetchCH <- msg
}

func (c *SinkerTable) processFetch() {
	c.processWg.Add(1)
	defer c.processWg.Done()
	var data FetchArray
	var err error
	flushFn := func(data Fetch, fields []string) {
		bufLength := len(data.GetData())
		c.logger.Info(c.sinkerTableConfig.TableName+" flush msg.", zap.Int("bufLength", bufLength))
		if bufLength == 0 {
			return
		}
		c.mux.Lock()
		c.fd.Write([]byte("]"))
		c.fd.Close()
		var label = strconv.FormatInt(time.Now().UnixNano(), 10)
		var newFileName = c.sinkerTableConfig.TableName + "." + label + ".wal.pending"
		err = os.Rename(c.sinkerTableConfig.TableName+".wal", newFileName)
		if err != nil {
			c.logger.Error(c.sinkerTableConfig.TableName+" rename wal error", zap.Error(err))
			return
		}
		c.needComma = false
		c.fd, err = os.OpenFile(c.sinkerTableConfig.TableName+".wal", os.O_CREATE|os.O_RDWR|os.O_APPEND, 0666)
		if err != nil {
			c.logger.Error(c.sinkerTableConfig.TableName+" open wal error", zap.Error(err))
		}
		if getTailChar(c.fd) == "" {
			c.fd.Write([]byte("["))
		}
		c.mux.Unlock()
		c.sendStarRocks(newFileName, label)
	}

	ticker := time.NewTicker(time.Duration(c.sinkerTableConfig.FlushInterval) * time.Second)
	defer ticker.Stop()
	var newKeys map[string]string
	var count int
	for {
		select {
		case fetch := <-c.fetchCH:
			if c.state.Load() == StateStopped {
				continue
			}
			tmpData := fetch.GetData()
			count = 0
			var tmpBuffer = bufferPool.Get()
			for i := range tmpData {
				count++
				var parse Metric
				parse, err = c.parser.Parse([]byte(tmpData[i]))
				if err != nil {
					c.logger.Error(c.sinkerTableConfig.TableName+" parse error", zap.Error(err), zap.String("data", tmpData[i]))
					continue
				}
				newKeys = parse.GetNewKeys(c.fieldMap)
				if len(newKeys) > 0 {
					for k, v := range newKeys {
						c.fieldMap.Store(k, &ColumnWithType{Name: k, Type: WhichType(v, true)})
					}
					err = c.flushFields()
					if err != nil {
						c.logger.Error(c.sinkerTableConfig.TableName+" flush fields error", zap.Error(err))
						continue
					}
				}
				// todo:修改原json添加not null字段的默认值
				c.fieldMap.Range(func(key, value any) bool {
					if !value.(*ColumnWithType).Type.Nullable {
						var obj any
						switch value.(*ColumnWithType).Type.Type {
						case TINYINT,
							SMALLINT,
							INT,
							BIGINT,
							LARGEINT,
							DECIMAL,
							DOUBLE,
							FLOAT,
							BOOLEAN:
							obj = 0
						case CHAR,
							STRING,
							VARCHAR,
							BINARY:
							obj = ""
						case DATE, DATETIME:
							obj = time.Now()
						default:
						}
						parse.Set(key.(string), obj)
					}
					return true
				})
				tmpBuffer.AppendBytes([]byte(parse.Value()))
				if i != len(tmpData)-1 {
					tmpBuffer.AppendByte(44)
				}
				parse.Close()
			}
			c.mux.Lock()
			if c.needComma {
				c.fd.Write([]byte(","))
			}
			c.needComma = true
			_, err = c.fd.Write(tmpBuffer.Bytes())
			if err != nil {
				c.logger.Error(c.sinkerTableConfig.TableName+" conn.Write error", zap.Error(err))
			}
			err = c.fd.Sync()
			if err != nil {
				c.logger.Error(c.sinkerTableConfig.TableName+" conn.Sync error", zap.Error(err))
			}
			c.mux.Unlock()
			tmpBuffer.Free()
			data.Data = append(data.Data, tmpData...)
			data.Callback = append(data.Callback, fetch.GetCallback()...)
			bufLength := len(data.Data)
			if bufLength >= c.sinkerTableConfig.BufferSize {
				tmp := data.Copy()
				go flushFn(tmp, c.getAllFields())
				data = FetchArray{}
				ticker.Reset(time.Duration(c.sinkerTableConfig.FlushInterval) * time.Second)
			}
		case <-ticker.C:
			c.logger.Debug("--------------------------------------------")
			tmp := data.Copy()
			go flushFn(tmp, c.getAllFields())
			data = FetchArray{}
		case <-c.ctx.Done():
			c.logger.Info("stopped processing loop", zap.String("table", c.sinkerTableConfig.Database+"."+c.sinkerTableConfig.TableName))
			return
		}
	}
}

func (c *SinkerTable) sendStarRocks(fileName, label string) {
	urlStr := fmt.Sprintf("http://%s:%s/api/%s/%s/_stream_load", c.sinkerTableConfig.FeHost, c.sinkerTableConfig.FeHttpPort, c.sinkerTableConfig.Database, c.sinkerTableConfig.TableName)
	targetURL, _ := url.Parse(urlStr)

	// Set up the request body
	bodyContent, _ := os.ReadFile(fileName)
	body := bytes.NewBuffer(bodyContent)

	// Create the request object
	req, err := http.NewRequest("PUT", targetURL.String(), body)
	if err != nil {
		fmt.Println("Error creating request:", err)
		return
	}

	// Set the request headers
	req.Header.Set("label", label)
	req.Header.Set("format", "json")
	req.Header.Set("Expect", "100-continue")
	req.Header.Set("strip_outer_array", "true")

	// Set additional headers with column names (this part might need adjustment based on the API requirements)
	columns := c.getAllFields()
	str := strings.Join(columns, ",")
	req.Header.Set("columns", str)

	// Set up basic authentication
	req.SetBasicAuth(c.sinkerTableConfig.Username, c.sinkerTableConfig.Password)

	// Create an HTTP client
	client := &http.Client{}

	// Send the request
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Error executing request:", err)
		return
	}
	defer resp.Body.Close()
	respBody, _ := io.ReadAll(resp.Body)
	// Print response status code
	fmt.Println("Response Status Code:", resp.StatusCode)
	fmt.Println("Response:", string(respBody))
	if gjson.GetBytes(respBody, "Status").String() == "Success" {
		c.logger.Info(c.sinkerTableConfig.TableName+" send starrocks success", zap.String("fileName", fileName))
		// todo:后续这里可以根据配置删除策略来删除文件，目前觉得不需要
		// Delete the file
		os.Remove(fileName)
	} else {
		c.logger.Error(c.sinkerTableConfig.TableName+" send starrocks error", zap.String("fileName", fileName), zap.String("respBody", string(respBody)))
	}
}
