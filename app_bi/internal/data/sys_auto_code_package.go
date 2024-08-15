package data

type SysAutoCodePackage struct {
	GVA_MODEL
	Desc        string `json:"desc" gorm:"comment:描述"`
	Label       string `json:"label" gorm:"comment:展示名"`
	Template    string `json:"template"  gorm:"comment:模版"`
	PackageName string `json:"packageName" gorm:"comment:包名"`
}

func (s *SysAutoCodePackage) TableName() string {
	return "sys_auto_code_packages"
}
