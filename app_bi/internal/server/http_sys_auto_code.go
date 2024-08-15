package server

import (
	"github.com/leaf-rain/raindata/app_bi/internal/conf"
	"github.com/leaf-rain/raindata/app_bi/internal/data"
	"io"
	"strings"

	"github.com/goccy/go-json"
	"github.com/leaf-rain/raindata/app_bi/third_party/rhttp"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type AutoCodeApi struct {
	data *data.Data
	log  *zap.Logger
	conf *conf.Bootstrap
}

// GetDB
// @Tags      AutoCode
// @Summary   获取当前所有数据库
// @Security  ApiKeyAuth
// @accept    application/json
// @Produce   application/json
// @Success   200  {object}  rhttp.Response{data=map[string]interface{},msg=string}  "获取当前所有数据库"
// @Router    /autoCode/getDatabase [get]
func (autoApi *AutoCodeApi) GetDB(c *gin.Context) {
	businessDB := c.Query("businessDB")
	dbs, err := autoCodeService.Database(businessDB).GetDB(businessDB)
	var dbList []map[string]interface{}
	for _, db := range svc.conf.DBList {
		var item = make(map[string]interface{})
		item["aliasName"] = db.AliasName
		item["dbName"] = db.Dbname
		item["disable"] = db.Disable
		item["dbtype"] = db.Type
		dbList = append(dbList, item)
	}
	if err != nil {
		svr.log.Error("获取失败!", zap.Error(err))
		rhttp.FailWithMessage("获取失败", c)
	} else {
		rhttp.OkWithDetailed(gin.H{"dbs": dbs, "dbList": dbList}, "获取成功", c)
	}
}

// GetTables
// @Tags      AutoCode
// @Summary   获取当前数据库所有表
// @Security  ApiKeyAuth
// @accept    application/json
// @Produce   application/json
// @Success   200  {object}  rhttp.Response{data=map[string]interface{},msg=string}  "获取当前数据库所有表"
// @Router    /autoCode/getTables [get]
func (autoApi *AutoCodeApi) GetTables(c *gin.Context) {
	dbName := c.DefaultQuery("dbName", svc.conf.Mysql.Dbname)
	businessDB := c.Query("businessDB")
	tables, err := autoCodeService.Database(businessDB).GetTables(businessDB, dbName)
	if err != nil {
		svr.log.Error("查询table失败!", zap.Error(err))
		rhttp.FailWithMessage("查询table失败", c)
	} else {
		rhttp.OkWithDetailed(gin.H{"tables": tables}, "获取成功", c)
	}
}

// GetColumn
// @Tags      AutoCode
// @Summary   获取当前表所有字段
// @Security  ApiKeyAuth
// @accept    application/json
// @Produce   application/json
// @Success   200  {object}  rhttp.Response{data=map[string]interface{},msg=string}  "获取当前表所有字段"
// @Router    /autoCode/getColumn [get]
func (autoApi *AutoCodeApi) GetColumn(c *gin.Context) {
	businessDB := c.Query("businessDB")
	dbName := c.DefaultQuery("dbName", svc.conf.Mysql.Dbname)
	tableName := c.Query("tableName")
	columns, err := autoCodeService.Database(businessDB).GetColumn(businessDB, tableName, dbName)
	if err != nil {
		svr.log.Error("获取失败!", zap.Error(err))
		rhttp.FailWithMessage("获取失败", c)
	} else {
		rhttp.OkWithDetailed(gin.H{"columns": columns}, "获取成功", c)
	}
}

func (autoApi *AutoCodeApi) LLMAuto(c *gin.Context) {
	prompt := c.Query("prompt")
	mode := c.Query("mode")
	params := make(map[string]string)
	params["prompt"] = prompt
	params["mode"] = mode
	path := strings.ReplaceAll(svc.conf.AutoCode.AiPath, "{FUNC}", "api/chat/ai")
	res, err := dto.HttpRequest(
		path,
		"POST",
		nil,
		params,
		nil,
	)
	if err != nil {
		svr.log.Error("大模型生成失败!", zap.Error(err))
		rhttp.FailWithMessage("大模型生成失败"+err.Error(), c)
		return
	}
	var resStruct rhttp.Response
	b, err := io.ReadAll(res.Body)
	defer res.Body.Close()
	if err != nil {
		svr.log.Error("大模型生成失败!", zap.Error(err))
		rhttp.FailWithMessage("大模型生成失败"+err.Error(), c)
		return
	}
	err = json.Unmarshal(b, &resStruct)
	if err != nil {
		svr.log.Error("大模型生成失败!", zap.Error(err))
		rhttp.FailWithMessage("大模型生成失败"+err.Error(), c)
		return
	}

	if resStruct.Code == 7 {
		svr.log.Error("大模型生成失败!"+resStruct.Msg, zap.Error(err))
		rhttp.FailWithMessage("大模型生成失败"+resStruct.Msg, c)
		return
	}
	rhttp.OkWithData(resStruct.Data, c)
}
