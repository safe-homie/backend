package store

type store struct {
	driver Driver
}

type Store interface {
	Migrate() error
}

func New(driver Driver) *store {
	return &store{driver: driver}
}

func (s *store) Migrate() error {
	return s.driver.Migrate()
}
