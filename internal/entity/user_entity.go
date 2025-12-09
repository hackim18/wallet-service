package entity

type User struct {
	ID           string `gorm:"column:id;primaryKey"`
	Name         string `gorm:"column:name"`
	Email        string `gorm:"column:email;uniqueIndex"`
	PasswordHash string `gorm:"column:password"`
	CreatedAt    int64  `gorm:"column:created_at;autoCreateTime:milli"`
	UpdatedAt    int64  `gorm:"column:updated_at;autoCreateTime:milli;autoUpdateTime:milli"`
}

func (u *User) TableName() string {
	return "users"
}
