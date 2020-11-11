package model

import "database/sql"

type ProductGroup struct {
	Id        int32
	Name      sql.NullString `db:"name"`
	CreatedAt string         `db:"created_at"`
	UpdatedAt string         `db:"updated_at"`
}
