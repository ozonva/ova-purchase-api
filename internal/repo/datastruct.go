package repo

import (
	"database/sql"
	"github.com/ozonva/ova-purchase-api/internal/purchase"
	"github.com/shopspring/decimal"
	"time"
)

type PurchaseRow struct {
	Id        uint64
	ItemId    sql.NullInt64 `db:"item_id"`
	UserId    uint64        `db:"user_id"`
	Total     decimal.Decimal
	Name      sql.NullString
	Price     sql.NullFloat64
	Quantity  sql.NullInt32
	CreatedAt time.Time `db:"created_at"`
	UpdatedAt time.Time `db:"updated_at"`
	Status    purchase.Status
}
