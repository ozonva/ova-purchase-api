package purchase

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/shopspring/decimal"
	"log"
	"time"
)

type Status uint

const (
	Created Status = iota
	Pending
	Success
	Failure
)

type Item struct {
	Id       uint64
	Name     string
	Price    decimal.Decimal
	Quantity uint
}

type Purchase struct {
	Id     uint64
	UserID uint64

	Total decimal.Decimal
	Items []Item

	CreatedAt time.Time
	UpdatedAt time.Time
	Status    Status
}

func New() Purchase {
	items := make([]Item, 0)
	return Purchase{
		Items:     items,
		Total:     decimal.Zero,
		Status:    Created,
		CreatedAt: time.Now(),
	}
}

func (s *Purchase) isChangeAllowed() (bool, error) {
	if s.Status == Created {
		return true, nil
	}
	return false, errors.New("purchase changes do not allow in this state")
}

func (s *Purchase) calculateTotal() {
	total := decimal.Zero
	for _, v := range s.Items {
		total = total.Add(v.Price.Mul(decimal.NewFromInt(int64(v.Quantity))))
	}
	s.Total = total
}

func (s *Purchase) String() string {
	bytes, err := json.Marshal(s)
	if err != nil {
		log.Fatal("can't marshal purchase")
	}
	return fmt.Sprintf("Purchase(%s)", string(bytes))
}

func (s *Purchase) Add(item Item) (bool, error) {
	if ok, err := s.isChangeAllowed(); !ok {
		return false, err
	}
	if item.Quantity == 0 {
		return false, errors.New("you can't add item with zero quantity")
	}
	found := false
	for k, v := range s.Items {
		if v.Id == item.Id {
			s.Items[k].Quantity = v.Quantity + item.Quantity
			found = true
			break
		}
	}
	if !found {
		s.Items = append(s.Items, item)
	}
	s.UpdatedAt = time.Now()
	s.calculateTotal()
	return true, nil
}

func (s *Purchase) Remove(itemId uint64) (bool, error) {
	if ok, err := s.isChangeAllowed(); !ok {
		return false, err
	}
	index := -1
	for k, v := range s.Items {
		if v.Id == itemId {
			if v.Quantity > 1 {
				s.Items[k].Quantity = v.Quantity - 1
				s.UpdatedAt = time.Now()
			} else {
				index = k
				break
			}
		}
	}
	if index != -1 {
		items := make([]Item, 0, len(s.Items)-1)
		items = append(items, s.Items[:index]...)
		items = append(items, s.Items[index+1:]...)
		s.Items = items
		s.UpdatedAt = time.Now()
	}
	s.calculateTotal()
	return true, nil
}

func (s *Purchase) Pending() (bool, error) {
	if s.Status == Pending {
		return true, nil
	}
	if s.Status == Created || s.Status == Failure {
		s.Status = Pending
		s.UpdatedAt = time.Now()
		return true, nil
	}
	return false, errors.New("you can't change status to pending from Success status")
}

func (s *Purchase) Success() (bool, error) {
	if s.Status == Success {
		return true, nil
	}
	if s.Status == Pending {
		s.Status = Success
		s.UpdatedAt = time.Now()
		return true, nil
	}
	return false, errors.New("you can't change status to success from not Pending status")
}

func (s *Purchase) Failure() (bool, error) {
	if s.Status == Failure {
		return true, nil
	}
	if s.Status == Pending {
		s.Status = Failure
		s.UpdatedAt = time.Now()
		return true, nil
	}
	return false, errors.New("you can't change status to failure from not Pending status")
}
