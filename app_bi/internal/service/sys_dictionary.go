package service

import (
	"errors"
	"github.com/leaf-rain/raindata/app_bi/internal/conf"
	"github.com/leaf-rain/raindata/app_bi/internal/data/entity"
	"go.uber.org/zap"

	"gorm.io/gorm"
)

//@function: CreateSysDictionary
//@description: 创建字典数据
//@param: sysDictionary data.SysDictionary
//@return: err error

type DictionaryService struct {
	data *entity.Data
	log  *zap.Logger
	conf *conf.Bootstrap
}

var DictionaryServiceApp = new(DictionaryService)

func (svc *DictionaryService) CreateSysDictionary(sysDictionary entity.SysDictionary) (err error) {
	if (!errors.Is(svc.data.SqlClient.First(&entity.SysDictionary{}, "type = ?", sysDictionary.Type).Error, gorm.ErrRecordNotFound)) {
		return errors.New("存在相同的type，不允许创建")
	}
	err = svc.data.SqlClient.Create(&sysDictionary).Error
	return err
}

//@function: DeleteSysDictionary
//@description: 删除字典数据
//@param: sysDictionary data.SysDictionary
//@return: err error

func (svc *DictionaryService) DeleteSysDictionary(sysDictionary entity.SysDictionary) (err error) {
	err = svc.data.SqlClient.Where("id = ?", sysDictionary.ID).Preload("SysDictionaryDetails").First(&sysDictionary).Error
	if err != nil && errors.Is(err, gorm.ErrRecordNotFound) {
		return errors.New("请不要搞事")
	}
	if err != nil {
		return err
	}
	err = svc.data.SqlClient.Delete(&sysDictionary).Error
	if err != nil {
		return err
	}

	if sysDictionary.SysDictionaryDetails != nil {
		return svc.data.SqlClient.Where("sys_dictionary_id=?", sysDictionary.ID).Delete(sysDictionary.SysDictionaryDetails).Error
	}
	return
}

//@function: UpdateSysDictionary
//@description: 更新字典数据
//@param: sysDictionary *data.SysDictionary
//@return: err error

func (svc *DictionaryService) UpdateSysDictionary(sysDictionary *entity.SysDictionary) (err error) {
	var dict entity.SysDictionary
	sysDictionaryMap := map[string]interface{}{
		"Name":   sysDictionary.Name,
		"Type":   sysDictionary.Type,
		"Status": sysDictionary.Status,
		"Desc":   sysDictionary.Desc,
	}
	err = svc.data.SqlClient.Where("id = ?", sysDictionary.ID).First(&dict).Error
	if err != nil {
		svc.log.Debug(err.Error())
		return errors.New("查询字典数据失败")
	}
	if dict.Type != sysDictionary.Type {
		if !errors.Is(svc.data.SqlClient.First(&entity.SysDictionary{}, "type = ?", sysDictionary.Type).Error, gorm.ErrRecordNotFound) {
			return errors.New("存在相同的type，不允许创建")
		}
	}
	err = svc.data.SqlClient.Model(&dict).Updates(sysDictionaryMap).Error
	return err
}

//@function: GetSysDictionary
//@description: 根据id或者type获取字典单条数据
//@param: Type string, Id uint
//@return: err error, sysDictionary data.SysDictionary

func (svc *DictionaryService) GetSysDictionary(Type string, Id uint, status *bool) (sysDictionary entity.SysDictionary, err error) {
	var flag = false
	if status == nil {
		flag = true
	} else {
		flag = *status
	}
	err = svc.data.SqlClient.Where("(type = ? OR id = ?) and status = ?", Type, Id, flag).Preload("SysDictionaryDetails", func(db *gorm.DB) *gorm.DB {
		return db.Where("status = ?", true).Order("sort")
	}).First(&sysDictionary).Error
	return
}

//@function: GetSysDictionaryInfoList
//@description: 分页获取字典列表
//@param: info dto.SysDictionarySearch
//@return: err error, list interface{}, total int64

func (svc *DictionaryService) GetSysDictionaryInfoList() (list interface{}, err error) {
	var sysDictionarys []entity.SysDictionary
	err = svc.data.SqlClient.Find(&sysDictionarys).Error
	return sysDictionarys, err
}
