package store

type (
	Token struct {
		Token string
	}
	FindToken struct {
		DeviceID string
	}
)

func (s *store) SaveToken(save *Token) (*Token, error) {
	return s.driver.SaveToken(save)
}

func (s *store) GetToken(find *FindToken) (*Token, error) {
	return s.driver.GetToken(find)
}
