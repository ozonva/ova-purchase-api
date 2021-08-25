package repo

import "github.com/ozonva/ova-purchase-api/internal/purchase"

// Repo - интерфейс хранилища для сущности
type Repo interface {
	AddPurchases(purchases []purchase.Purchase) error
	ListPurchases(limit, offset uint64) ([]purchase.Purchase, error)
	DescribePurchase(purchaseId uint64) (*purchase.Purchase, error)
}
