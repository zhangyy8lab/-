package account

type Role struct {
	ID       uint64 `gorm:"foreignKey; autoIncrement" json:"id"`
	Name     string `gorm:"type:varchar(32);" json:"name"`
	RoleKind int    `json:"role_kind" describe:"类别id"` //
	Title    string `gorm:"type:varchar(32);" json:"title" describe:"中文名称"`
	AsName   string `gorm:"type:varchar(32);" json:"-" describe:"中文名称"`
	RoleDesc string `gorm:"type:varchar(300);" json:"desc" describe:"中文描述信息"`
	EnDesc   string `gorm:"type:varchar(300);" json:"en_desc" describe:"英文描述信息"`
}

func (Role) TableName() string {
	return "role"
}

// RoleKind
/*
roleKind	asName
	1		user
	2		user
	3		ops
	4		None
	5		ops
	6		test
	7		dev
*/
