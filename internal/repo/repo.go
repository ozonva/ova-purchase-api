package repo

import "github.com/ozonva/ova-purchase-api/internal/purchase"

// Repo - интерфейс хранилища для сущности
type Repo interface {
	AddEntities(entities []purchase.Purchase) error
	ListEntities(limit, offset uint64) ([]purchase.Purchase, error)
	DescribeEntity(entityId uint64) (*purchase.Purchase, error)
}
