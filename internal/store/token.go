package store

type Token struct {
	Token string
}

func (s *store) SaveToken(save *Token) (*Token, error) {
	return s.driver.SaveToken(save)
}
