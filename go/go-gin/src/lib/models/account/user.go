package account

type User struct {
	ID     uint64 `gorm:"primaryKey; autoIncrement" json:"id"`
	Name   string `gorm:"type:varchar(32);" json:"username"`
	RoleId uint64 `json:"role_id"`
	Role   Role   `gorm:"foreignKey:RoleId; constraint:OnDelete:CASCADE" json:"role"`
	Uuid   string `gorm:"type:varchar(36); user" json:"uuid"`
}

func (User) TableName() string {
	return "user"
}
