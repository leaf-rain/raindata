package service

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/leaf-rain/raindata/app_bi/internal/conf"
	"github.com/leaf-rain/raindata/app_bi/internal/data/dto"
	"github.com/leaf-rain/raindata/app_bi/internal/data/entity"
	"github.com/leaf-rain/raindata/app_bi/third_party/rhttp"
	"github.com/leaf-rain/raindata/app_bi/third_party/utils"
	"go.uber.org/zap"
	"gorm.io/gorm"
	"mime/multipart"
	"net/url"
	"strconv"
	"strings"
	"time"
)

type SysExportTemplateService struct {
	data *entity.Data
	log  *zap.Logger
	conf *conf.Bootstrap
}

var SysExportTemplateServiceApp = new(SysExportTemplateService)

// CreateSysExportTemplate 创建导出模板记录
func (svc *SysExportTemplateService) CreateSysExportTemplate(sysExportTemplate *entity.SysExportTemplate) (err error) {
	err = svc.data.SqlClient.Create(sysExportTemplate).Error
	return err
}

// DeleteSysExportTemplate 删除导出模板记录
func (svc *SysExportTemplateService) DeleteSysExportTemplate(sysExportTemplate entity.SysExportTemplate) (err error) {
	err = svc.data.SqlClient.Delete(&sysExportTemplate).Error
	return err
}

// DeleteSysExportTemplateByIds 批量删除导出模板记录
func (svc *SysExportTemplateService) DeleteSysExportTemplateByIds(ids rhttp.IdsReq) (err error) {
	err = svc.data.SqlClient.Delete(&[]entity.SysExportTemplate{}, "id in ?", ids.Ids).Error
	return err
}

// UpdateSysExportTemplate 更新导出模板记录
func (svc *SysExportTemplateService) UpdateSysExportTemplate(sysExportTemplate entity.SysExportTemplate) (err error) {
	return svc.data.SqlClient.Transaction(func(tx *gorm.DB) error {
		conditions := sysExportTemplate.Conditions
		e := tx.Delete(&[]entity.Condition{}, "template_id = ?", sysExportTemplate.TemplateID).Error
		if e != nil {
			return e
		}
		sysExportTemplate.Conditions = nil

		joins := sysExportTemplate.JoinTemplate
		e = tx.Delete(&[]entity.JoinTemplate{}, "template_id = ?", sysExportTemplate.TemplateID).Error
		if e != nil {
			return e
		}
		sysExportTemplate.JoinTemplate = nil

		e = tx.Updates(&sysExportTemplate).Error
		if e != nil {
			return e
		}
		if len(conditions) > 0 {
			for i := range conditions {
				conditions[i].ID = 0
			}
			e = tx.Create(&conditions).Error
		}
		if len(joins) > 0 {
			for i := range joins {
				joins[i].ID = 0
			}
			e = tx.Create(&joins).Error
		}
		return e
	})
}

// GetSysExportTemplate 根据id获取导出模板记录
func (svc *SysExportTemplateService) GetSysExportTemplate(id uint) (sysExportTemplate entity.SysExportTemplate, err error) {
	err = svc.data.SqlClient.Where("id = ?", id).Preload("JoinTemplate").Preload("Conditions").First(&sysExportTemplate).Error
	return
}

// GetSysExportTemplateInfoList 分页获取导出模板记录
func (svc *SysExportTemplateService) GetSysExportTemplateInfoList(info dto.SysExportTemplateSearch) (list []entity.SysExportTemplate, total int64, err error) {
	limit := info.PageSize
	offset := info.PageSize * (info.Page - 1)
	// 创建db
	db := svc.data.SqlClient.Model(&entity.SysExportTemplate{})
	var sysExportTemplates []entity.SysExportTemplate
	// 如果有条件搜索 下方会自动创建搜索语句
	if info.StartCreatedAt != nil && info.EndCreatedAt != nil {
		db = db.Where("created_at BETWEEN ? AND ?", info.StartCreatedAt, info.EndCreatedAt)
	}
	if info.Name != "" {
		db = db.Where("name LIKE ?", "%"+info.Name+"%")
	}
	if info.TableName != "" {
		db = db.Where("table_name = ?", info.TableName)
	}
	if info.TemplateID != "" {
		db = db.Where("template_id = ?", info.TemplateID)
	}
	err = db.Count(&total).Error
	if err != nil {
		return
	}

	if limit != 0 {
		db = db.Limit(limit).Offset(offset)
	}

	err = db.Find(&sysExportTemplates).Error
	return sysExportTemplates, total, err
}

// ExportExcel 导出Excel
func (svc *SysExportTemplateService) ExportExcel(templateID string, values url.Values) (file *bytes.Buffer, name string, err error) {
	var template entity.SysExportTemplate
	err = svc.data.SqlClient.Preload("Conditions").Preload("JoinTemplate").First(&template, "template_id = ?", templateID).Error
	if err != nil {
		return nil, "", err
	}
	f := excelize.NewFile()
	defer func() {
		if err := f.Close(); err != nil {
			fmt.Println(err)
		}
	}()
	// Create a new sheet.
	index, err := f.NewSheet("Sheet1")
	if err != nil {
		fmt.Println(err)
		return
	}
	var templateInfoMap = make(map[string]string)
	columns, err := utils.GetJSONKeys(template.TemplateInfo)
	if err != nil {
		return nil, "", err
	}
	err = json.Unmarshal([]byte(template.TemplateInfo), &templateInfoMap)
	if err != nil {
		return nil, "", err
	}
	var tableTitle []string
	var selectKeyFmt []string
	for _, key := range columns {
		selectKeyFmt = append(selectKeyFmt, fmt.Sprintf("`%s`", key))
		tableTitle = append(tableTitle, templateInfoMap[key])
	}

	selects := strings.Join(selectKeyFmt, ", ")
	var tableMap []map[string]interface{}
	db := svc.data.SqlClient
	if len(template.JoinTemplate) > 0 {
		for _, join := range template.JoinTemplate {
			db = db.Joins(join.JOINS + "`" + join.Table + "`" + " ON " + join.ON)
		}
	}

	db = db.Select(selects).Table(template.TableName)

	if len(template.Conditions) > 0 {
		for _, condition := range template.Conditions {
			sql := fmt.Sprintf("%s %s ?", condition.Column, condition.Operator)
			value := values.Get(condition.From)
			if value != "" {
				if condition.Operator == "LIKE" {
					value = "%" + value + "%"
				}
				db = db.Where(sql, value)
			}
		}
	}
	// 通过参数传入limit
	limit := values.Get("limit")
	if limit != "" {
		l, e := strconv.Atoi(limit)
		if e == nil {
			db = db.Limit(l)
		}
	}
	// 模板的默认limit
	if limit == "" && template.Limit != nil && *template.Limit != 0 {
		db = db.Limit(*template.Limit)
	}

	// 通过参数传入offset
	offset := values.Get("offset")
	if offset != "" {
		o, e := strconv.Atoi(offset)
		if e == nil {
			db = db.Offset(o)
		}
	}

	// 获取当前表的所有字段
	table := template.TableName
	orderColumns, err := db.Migrator().ColumnTypes(table)
	if err != nil {
		return nil, "", err
	}

	// 创建一个 map 来存储字段名
	fields := make(map[string]bool)

	for _, column := range orderColumns {
		fields[column.Name()] = true
	}

	// 通过参数传入order
	order := values.Get("order")

	if order == "" && template.Order != "" {
		// 如果没有order入参，这里会使用模板的默认排序
		order = template.Order
	}

	if order != "" {
		checkOrderArr := strings.Split(order, " ")
		orderStr := ""
		// 检查请求的排序字段是否在字段列表中
		if _, ok := fields[checkOrderArr[0]]; !ok {
			return nil, "", fmt.Errorf("order by %s is not in the fields", order)
		}
		orderStr = checkOrderArr[0]
		if len(checkOrderArr) > 1 {
			if checkOrderArr[1] != "asc" && checkOrderArr[1] != "desc" {
				return nil, "", fmt.Errorf("order by %s is not secure", order)
			}
			orderStr = orderStr + " " + checkOrderArr[1]
		}
		db = db.Order(orderStr)
	}

	err = db.Debug().Find(&tableMap).Error
	if err != nil {
		return nil, "", err
	}
	var rows [][]string
	rows = append(rows, tableTitle)
	for _, exTable := range tableMap {
		var row []string
		for _, column := range columns {
			if len(template.JoinTemplate) > 0 {
				columnAs := strings.Split(column, " as ")
				if len(columnAs) > 1 {
					column = strings.TrimSpace(strings.Split(column, " as ")[1])
				} else {
					columnArr := strings.Split(column, ".")
					if len(columnArr) > 1 {
						column = strings.Split(column, ".")[1]
					}
				}
			}
			// 需要对时间类型特殊处理
			if t, ok := exTable[column].(time.Time); ok {
				row = append(row, t.Format("2006-01-02 15:04:05"))
			} else {
				row = append(row, fmt.Sprintf("%v", exTable[column]))
			}
		}
		rows = append(rows, row)
	}
	for i, row := range rows {
		for j, colCell := range row {
			sErr := f.SetCellValue("Sheet1", fmt.Sprintf("%s%d", getColumnName(j+1), i+1), colCell)
			if sErr != nil {
				return nil, "", sErr
			}
		}
	}
	f.SetActiveSheet(index)
	file, err = f.WriteToBuffer()
	if err != nil {
		return nil, "", err
	}

	return file, template.Name, nil
}

// ExportTemplate 导出Excel模板
func (svc *SysExportTemplateService) ExportTemplate(templateID string) (file *bytes.Buffer, name string, err error) {
	var template entity.SysExportTemplate
	err = svc.data.SqlClient.First(&template, "template_id = ?", templateID).Error
	if err != nil {
		return nil, "", err
	}
	f := excelize.NewFile()
	defer func() {
		if err := f.Close(); err != nil {
			fmt.Println(err)
		}
	}()
	// Create a new sheet.
	index, err := f.NewSheet("Sheet1")
	if err != nil {
		fmt.Println(err)
		return
	}
	var templateInfoMap = make(map[string]string)

	columns, err := utils.GetJSONKeys(template.TemplateInfo)

	err = json.Unmarshal([]byte(template.TemplateInfo), &templateInfoMap)
	if err != nil {
		return nil, "", err
	}
	var tableTitle []string
	for _, key := range columns {
		tableTitle = append(tableTitle, templateInfoMap[key])
	}

	for i := range tableTitle {
		fErr := f.SetCellValue("Sheet1", fmt.Sprintf("%s%d", getColumnName(i+1), 1), tableTitle[i])
		if fErr != nil {
			return nil, "", fErr
		}
	}
	f.SetActiveSheet(index)
	file, err = f.WriteToBuffer()
	if err != nil {
		return nil, "", err
	}

	return file, template.Name, nil
}

// ImportExcel 导入Excel
func (svc *SysExportTemplateService) ImportExcel(templateID string, file *multipart.FileHeader) (err error) {
	var template entity.SysExportTemplate
	err = svc.data.SqlClient.First(&template, "template_id = ?", templateID).Error
	if err != nil {
		return err
	}

	src, err := file.Open()
	if err != nil {
		return err
	}
	defer src.Close()

	f, err := excelize.OpenReader(src)
	if err != nil {
		return err
	}

	rows, err := f.GetRows("Sheet1")
	if err != nil {
		return err
	}

	var templateInfoMap = make(map[string]string)
	err = json.Unmarshal([]byte(template.TemplateInfo), &templateInfoMap)
	if err != nil {
		return err
	}

	var titleKeyMap = make(map[string]string)
	for key, title := range templateInfoMap {
		titleKeyMap[title] = key
	}

	db := svc.data.SqlClient
	return db.Transaction(func(tx *gorm.DB) error {
		excelTitle := rows[0]
		values := rows[1:]
		items := make([]map[string]interface{}, 0, len(values))
		for _, row := range values {
			var item = make(map[string]interface{})
			for ii, value := range row {
				key := titleKeyMap[excelTitle[ii]]
				item[key] = value
			}

			needCreated := tx.Migrator().HasColumn(template.TableName, "created_at")
			needUpdated := tx.Migrator().HasColumn(template.TableName, "updated_at")

			if item["created_at"] == nil && needCreated {
				item["created_at"] = time.Now()
			}
			if item["updated_at"] == nil && needUpdated {
				item["updated_at"] = time.Now()
			}

			items = append(items, item)
		}
		cErr := tx.Table(template.TableName).CreateInBatches(&items, 1000).Error
		return cErr
	})
}

func getColumnName(n int) string {
	columnName := ""
	for n > 0 {
		n--
		columnName = string(rune('A'+n%26)) + columnName
		n /= 26
	}
	return columnName
}
