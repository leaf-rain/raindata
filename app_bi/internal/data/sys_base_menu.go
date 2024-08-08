package data

type SysBaseMenu struct {
	GVA_MODEL
	MenuLevel     uint                                       `json:"-"`
	ParentId      uint                                       `json:"parentId" gorm:"comment:父菜单ID"`     // 父菜单ID
	Path          string                                     `json:"path" gorm:"comment:路由path"`        // 路由path
	Name          string                                     `json:"name" gorm:"comment:路由name"`        // 路由name
	Hidden        bool                                       `json:"hidden" gorm:"comment:是否在列表隐藏"`     // 是否在列表隐藏
	Component     string                                     `json:"component" gorm:"comment:对应前端文件路径"` // 对应前端文件路径
	Sort          int                                        `json:"sort" gorm:"comment:排序标记"`          // 排序标记
	Meta          `json:"meta" gorm:"embedded;comment:附加属性"` // 附加属性
	SysAuthoritys []SysAuthority                             `json:"authoritys" gorm:"many2many:sys_authority_menus;"`
	Children      []SysBaseMenu                              `json:"children" gorm:"-"`
	Parameters    []SysBaseMenuParameter                     `json:"parameters"`
	MenuBtn       []SysBaseMenuBtn                           `json:"menuBtn"`
}

type SysBaseMenuBtn struct {
	GVA_MODEL
	Name          string `json:"name" gorm:"comment:按钮关键key"`
	Desc          string `json:"desc" gorm:"按钮备注"`
	SysBaseMenuID uint   `json:"sysBaseMenuID" gorm:"comment:菜单ID"`
}

type Meta struct {
	ActiveName  string `json:"activeName" gorm:"comment:高亮菜单"`
	KeepAlive   bool   `json:"keepAlive" gorm:"comment:是否缓存"`           // 是否缓存
	DefaultMenu bool   `json:"defaultMenu" gorm:"comment:是否是基础路由（开发中）"` // 是否是基础路由（开发中）
	Title       string `json:"title" gorm:"comment:菜单名"`                // 菜单名
	Icon        string `json:"icon" gorm:"comment:菜单图标"`                // 菜单图标
	CloseTab    bool   `json:"closeTab" gorm:"comment:自动关闭tab"`         // 自动关闭tab
}

type SysBaseMenuParameter struct {
	SysBaseMenuID uint
	Type          string `json:"type" gorm:"comment:地址栏携带参数为params还是query"` // 地址栏携带参数为params还是query
	Key           string `json:"key" gorm:"comment:地址栏携带参数的key"`            // 地址栏携带参数的key
	Value         string `json:"value" gorm:"comment:地址栏携带参数的值"`            // 地址栏携带参数的值
}

func (SysBaseMenu) TableName() string {
	return "sys_base_menus"
}
