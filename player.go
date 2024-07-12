package egamecore

type Player struct {
	Id      uint64  `json:"id"`
	Name    string  `json:"name"`
	Balance float64 `json:"balance"`
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
