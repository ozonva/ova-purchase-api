package repo

import (
	"context"
	"errors"
	"github.com/jmoiron/sqlx"
	"github.com/ozonva/ova-purchase-api/internal/purchase"
	"github.com/shopspring/decimal"
)

var PurchaseNotFoundError = errors.New("Purchase not found")

// Repo - интерфейс хранилища для сущности
type Repo interface {
	AddPurchase(ctx context.Context, purchases purchase.Purchase) (uint64, error)
	UpdatePurchase(ctx context.Context, purchases purchase.Purchase) error
	AddPurchases(ctx context.Context, purchases []purchase.Purchase) ([]uint64, error)
	ListPurchases(ctx context.Context, limit, offset uint) ([]purchase.Purchase, error)
	DescribePurchase(ctx context.Context, purchaseId uint64) (*purchase.Purchase, error)
	RemovePurchase(ctx context.Context, purchaseId uint64) error
}

type repo struct {
	db *sqlx.DB
}

func NewRepo(db *sqlx.DB) Repo {
	return &repo{
		db: db,
	}
}

func insertPurchase(ctx context.Context, tx *sqlx.Tx, purchase purchase.Purchase) (uint64, error) {
	row := tx.QueryRowContext(ctx, `insert into purchases (user_id, total, updated_at, status) values ($1, $2, $3, $4) RETURNING id`, purchase.UserID, purchase.Total, purchase.UpdatedAt, purchase.Status.String())
	purchaseId := 0
	err := row.Scan(&purchaseId)
	if err != nil {
		return 0, err
	}
	statement, err := tx.Prepare("insert into purchase_items (purchase_id, name, price, quantity) values ($1, $2, $3, $4)")
	if err != nil {
		return 0, err
	}
	for _, item := range purchase.Items {
		_, err = statement.ExecContext(ctx, purchaseId, item.Name, item.Price, item.Quantity)
		if err != nil {
			return 0, err
		}
	}
	return uint64(purchaseId), nil
}

func (s *repo) AddPurchases(ctx context.Context, purchases []purchase.Purchase) ([]uint64, error) {
	tx, err := s.db.BeginTxx(ctx, nil)
	if err != nil {
		return nil, err
	}

	defer tx.Rollback()

	ids := make([]uint64, 0, len(purchases))

	for _, p := range purchases {
		id, err := insertPurchase(ctx, tx, p)
		if err != nil {
			return nil, err
		}
		ids = append(ids, id)
	}
	err = tx.Commit()
	return ids, err
}

func (s *repo) AddPurchase(ctx context.Context, purchase purchase.Purchase) (uint64, error) {
	tx, err := s.db.BeginTxx(ctx, nil)
	if err != nil {
		return 0, err
	}
	defer tx.Rollback()

	id, err := insertPurchase(ctx, tx, purchase)
	if err != nil {
		return 0, err
	}
	err = tx.Commit()
	return id, err
}

func (s *repo) ListPurchases(ctx context.Context, limit, offset uint) ([]purchase.Purchase, error) {
	rows, err := s.db.QueryxContext(ctx, `
			with ps as (select * from purchases offset $1 limit $2)
			select ps.*, i.id as item_id, i.name, i.price, i.quantity from ps left join purchase_items i on i.purchase_id = ps.id
		`, offset, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	selected := make([]PurchaseRow, 0)
	for rows.Next() {
		entity := PurchaseRow{}
		err = rows.StructScan(&entity)
		if err != nil {
			return nil, err
		}
		selected = append(selected, entity)
	}
	purchases := make([]purchase.Purchase, 0)
	exist := make(map[uint64]purchase.Purchase)
	for _, item := range selected {
		if p, ok := exist[item.Id]; !ok {
			added := purchase.Purchase{
				Id:        item.Id,
				UserID:    item.UserId,
				Total:     item.Total,
				CreatedAt: item.CreatedAt,
				UpdatedAt: item.UpdatedAt,
				Items:     make([]purchase.Item, 0),
				Status:    item.Status,
			}
			exist[item.Id] = added
			if item.ItemId.Valid {
				added.Items = append(added.Items, purchase.Item{
					Id:       uint64(item.ItemId.Int64),
					Name:     item.Name.String,
					Price:    decimal.NewFromFloat(item.Price.Float64),
					Quantity: uint(item.Quantity.Int32),
				})
			}
			purchases = append(purchases, added)
		} else {
			p.Items = append(p.Items, purchase.Item{
				Id:       uint64(item.ItemId.Int64),
				Name:     item.Name.String,
				Price:    decimal.NewFromFloat(item.Price.Float64),
				Quantity: uint(item.Quantity.Int32),
			})
		}
	}
	return purchases, nil
}

func (s *repo) DescribePurchase(ctx context.Context, purchaseId uint64) (*purchase.Purchase, error) {
	rows, err := s.db.QueryxContext(ctx, `
			select p.*, i.id as item_id, i.name, i.price, i.quantity 
					from purchases p 
					    	left join purchase_items i on i.purchase_id = p.id
						where p.id = $1
	`, purchaseId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var p *purchase.Purchase
	for rows.Next() {
		row := PurchaseRow{}
		if err = rows.StructScan(&row); err != nil {
			return nil, err
		}
		if p == nil {
			p = &purchase.Purchase{
				Id:        row.Id,
				Total:     row.Total,
				CreatedAt: row.CreatedAt,
				UpdatedAt: row.UpdatedAt,
				Status:    row.Status,
				Items:     make([]purchase.Item, 0),
			}
		}
		if row.ItemId.Valid {
			p.Items = append(p.Items, purchase.Item{
				Id:       uint64(row.ItemId.Int64),
				Name:     row.Name.String,
				Price:    decimal.NewFromFloat(row.Price.Float64),
				Quantity: uint(row.Quantity.Int32),
			})
		}
	}
	if p == nil {
		return nil, PurchaseNotFoundError
	}
	return p, nil
}

func (s *repo) RemovePurchase(ctx context.Context, purchaseId uint64) error {
	tx, err := s.db.BeginTxx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	_, err = tx.ExecContext(ctx, `delete from purchase_items where purchase_id = $1`, purchaseId)
	if err != nil {
		return err
	}
	result, err := tx.ExecContext(ctx, `delete from purchases where id = $1`, purchaseId)
	if err != nil {
		return err
	}
	if count, _ := result.RowsAffected(); count == 0 {
		return PurchaseNotFoundError
	}
	err = tx.Commit()
	return err
}

func (s *repo) UpdatePurchase(ctx context.Context, purchase purchase.Purchase) error {
	tx, err := s.db.BeginTxx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	_, err = tx.ExecContext(ctx, `delete from purchase_items where purchase_id = $1`, purchase.Id)
	if err != nil {
		return err
	}
	statement, err := tx.Prepare("insert into purchase_items (purchase_id, name, price, quantity) values ($1, $2, $3, $4)")
	if err != nil {
		return err
	}
	for _, item := range purchase.Items {
		_, err = statement.ExecContext(ctx, purchase.Id, item.Name, item.Price, item.Quantity)
		if err != nil {
			return err
		}
	}
	_, err = tx.ExecContext(ctx, `update purchases set total=$2, updated_at=now() where purchase_id = $1`, purchase.Id, purchase.Total)
	if err != nil {
		return err
	}
	err = tx.Commit()
	return err
}
