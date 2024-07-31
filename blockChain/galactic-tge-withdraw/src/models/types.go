package models

// Account 表示数据库中的账户表
type Account struct {
	ID      int    `gorm:"primaryKey"`
	Account string `json:"account"`
	Star    int    `json:"star"`
}

// TableName 设置Account的表名
func (Account) TableName() string {
	return "account" // 或者 "Accounts"
}

type WithDraw struct {
	ID        int     `gorm:"primaryKey"`
	Account   string  `json:"account"`
	Stars     int     `json:"stars"`
	Chain     string  `json:"chain"`
	Height    int64   `json:"height"`
	CoinKind  string  `json:"coin_kind"`
	CoinCount float64 `json:"coin_count"`
	CoinPrice float64 `json:"coin_price"`
	Value     float64 `json:"value"`
}

func (WithDraw) TableName() string {
	return "with_draw"
}

type WithCount struct {
	ID         int     `gorm:"primaryKey"`
	Account    string  `json:"account"`
	Stars      int     `json:"stars"`
	CountValue float64 `json:"count_value"`
}

func (WithCount) TableName() string {
	return "with_count"
}
