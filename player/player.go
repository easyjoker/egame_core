package player

type Player struct {
	Id       int64   `json:"id" gorm:"primaryKey;autoIncrement"  mapstructure:"id"`
	Account  string  `json:"account" gorm:"column:account" mapstructure:"account"`
	Password string  `json:"password" gorm:"column:password" mapstructure:"password"`
	Name     string  `json:"name" gorm:"column:name" mapstructure:"name"`
	Balance  float64 `json:"balance" gorm:"column:balance" mapstructure:"balance"`
}

func (p *Player) CheckEnoughBalance(amount float64) bool {
	return p.Balance >= amount
}

func (p *Player) DeductBalance(amount float64) {
	p.Balance -= amount
}

func (p *Player) AddBalance(amount float64) {
	p.Balance += amount
}
