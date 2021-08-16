package purchase

import (
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

func TestPurchase_NewPurchase(t *testing.T) {
	first := New()

	require.NotEqual(t, Purchase{}, first)
	require.NotNil(t, first.Items)
	require.NotEqual(t, time.Time{}, first.CreatedAt)

	require.Equal(t, decimal.Zero, first.Total)
	require.Equal(t, Created, first.Status)
}

func TestPurchase_String(t *testing.T) {
	first := New()

	ok, err := first.Add(Item{
		Id:       1,
		Price:    decimal.NewFromInt(1990),
		Quantity: 5,
		Name:     "Iphone",
	})

	require.NoError(t, err)
	require.Equal(t, true, ok)

	str := first.String()

	require.Contains(t, str, "Iphone")
	require.Contains(t, str, "9950")
}

func TestPurchase_AddPurchaseSuccess(t *testing.T) {
	first := New()
	item := Item{
		Id:       1,
		Name:     "Iphone",
		Price:    decimal.NewFromInt(999),
		Quantity: 2,
	}
	ok, err := first.Add(item)

	assert.Equal(t, true, ok)
	assert.NoError(t, err)

	assert.Equal(t, 1, len(first.Items))
	assert.Equal(t, decimal.NewFromInt(1998), first.Total)
}

func TestPurchase_AddPurchaseFailure(t *testing.T) {
	first := New()
	first.Status = Success
	item := Item{
		Id:       1,
		Name:     "Iphone",
		Price:    decimal.NewFromInt(999),
		Quantity: 2,
	}
	ok, err := first.Add(item)

	assert.Equal(t, false, ok)
	assert.Error(t, err)

	assert.Equal(t, 0, len(first.Items))
	assert.Equal(t, decimal.Zero, first.Total)
}

func TestPurchase_AddPurchaseTwoTimes(t *testing.T) {
	first := New()

	item := Item{
		Id:       1,
		Name:     "Iphone",
		Price:    decimal.NewFromInt(100),
		Quantity: 2,
	}
	ok1, err1 := first.Add(item)

	assert.Equal(t, true, ok1)
	assert.NoError(t, err1)

	ok2, err2 := first.Add(item)

	assert.Equal(t, true, ok2)
	assert.NoError(t, err2)

	assert.Equal(t, 1, len(first.Items))
	assert.Equal(t, decimal.NewFromInt(400), first.Total)
	assert.Equal(t, uint(4), first.Items[0].Quantity)
}

func TestPurchase_RemoveOneQuantity(t *testing.T) {
	first := New()
	item := Item{
		Id:       1,
		Name:     "Iphone",
		Price:    decimal.NewFromInt(999),
		Quantity: 1,
	}

	addOk, addErr := first.Add(item)

	assert.Equal(t, true, addOk)
	assert.NoError(t, addErr)

	removeOk, removeErr := first.Remove(1)

	assert.Equal(t, true, removeOk)
	assert.NoError(t, removeErr)

	assert.Equal(t, 0, len(first.Items))
	assert.Equal(t, decimal.Zero, first.Total)
}

func TestPurchase_RemoveFromManyQuantity(t *testing.T) {
	first := New()
	item := Item{
		Id:       1,
		Name:     "Iphone",
		Price:    decimal.NewFromInt(999),
		Quantity: 3,
	}

	addOk, addErr := first.Add(item)

	assert.Equal(t, true, addOk)
	assert.NoError(t, addErr)

	removeOk, removeErr := first.Remove(1)

	assert.Equal(t, true, removeOk)
	assert.NoError(t, removeErr)

	assert.Equal(t, 1, len(first.Items))
	assert.Equal(t, decimal.NewFromInt(1998), first.Total)
	assert.Equal(t, uint(2), first.Items[0].Quantity)
}
