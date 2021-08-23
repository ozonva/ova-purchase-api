package flusher

import (
	"github.com/ozonva/ova-purchase-api/internal/purchase"
	"github.com/ozonva/ova-purchase-api/internal/repo"
	"github.com/ozonva/ova-purchase-api/internal/utils"
)

type Flusher interface {
	Flush(entities []purchase.Purchase) []purchase.Purchase
}

func NewFlusher(chunkSize uint, repository repo.Repo) *flusher {
	return &flusher{
		chunkSize: chunkSize,
		repo:      repository,
	}
}

type flusher struct {
	chunkSize uint
	repo      repo.Repo
}

func (s *flusher) Flush(entities []purchase.Purchase) []purchase.Purchase {
	batch, err := utils.SplitToBulks(entities, s.chunkSize)
	if err != nil {
		return entities
	}
	result := make([]purchase.Purchase, 0)
	for _, items := range batch {
		if err := s.repo.AddEntities(items); err != nil {
			result = append(result, items...)
		}
	}
	if len(result) == 0 {
		return nil
	}
	return result
}
