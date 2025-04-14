package postgres

import (
	"context"

	"github.com/safe-homie/backend/internal/store"
)

func (p *_postgres) SaveToken(save *store.Token) (*store.Token, error) {
	stmt := `UPDATE tokens
			 SET token = $1
			 WHERE device_id = $2
			 RETURNING device_id, token
			`
	if _, err := p.db.Query(context.Background(), stmt, save.Token, "default-device-id"); err != nil {
		return nil, err
	}
	return save, nil
}
