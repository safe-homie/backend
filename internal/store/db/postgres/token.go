package postgres

import (
	"context"

	"github.com/safe-homie/backend/internal/domain"
	"github.com/safe-homie/backend/internal/store"
)

func (p *_postgres) SaveToken(save *store.Token) (*store.Token, error) {
	stmt := `UPDATE tokens
			 SET token = $1
			 WHERE device_id = $2
			 RETURNING device_id, token
			`
	if _, err := p.db.Query(context.Background(), stmt, save.Token, domain.DefaultDeviceID); err != nil {
		return nil, err
	}
	return save, nil
}

func (p *_postgres) GetToken(find *store.FindToken) (*store.Token, error) {
	stmt := `SELECT token FROM tokens`
	var token store.Token
	if err := p.db.QueryRow(context.Background(), stmt).Scan(
		&token.Token,
	); err != nil {
		return nil, err
	}
	return &token, nil
}
