package flusher

import (
	"errors"
	"github.com/golang/mock/gomock"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/ozonva/ova-purchase-api/internal/mocks"
	"github.com/ozonva/ova-purchase-api/internal/purchase"
)

var _ = Describe("Flusher", func() {
	var (
		repoMock  *mocks.MockRepo
		purchases []purchase.Purchase
	)
	BeforeEach(func() {
		ctrl := gomock.NewController(GinkgoT())
		repoMock = mocks.NewMockRepo(ctrl)
		defer ctrl.Finish()

		purchases = []purchase.Purchase{
			purchase.New(), purchase.New(), purchase.New(),
		}
	})

	Describe("Flush purchases", func() {
		Context("Batch size is 2", func() {
			It("should return nil", func() {
				flusherInstance := NewFlusher(2, repoMock)

				gomock.InOrder(
					repoMock.EXPECT().AddPurchases(gomock.Any(), gomock.Eq(purchases[0:2])).Return(nil),
					repoMock.EXPECT().AddPurchases(gomock.Any(), gomock.Eq(purchases[2:3])).Return(nil),
				)

				Expect(flusherInstance.Flush(purchases) == nil).To(BeTrue())
			})
		})

		Context("Batch size is 3", func() {
			It("should return nil", func() {
				flusherInstance := NewFlusher(3, repoMock)

				gomock.InOrder(
					repoMock.EXPECT().AddPurchases(gomock.Any(), gomock.Eq(purchases[0:3])).Return(nil),
				)

				Expect(flusherInstance.Flush(purchases) == nil).To(BeTrue())
			})
		})

		Context("Batch size is 10", func() {
			It("should return nil", func() {
				flusherInstance := NewFlusher(10, repoMock)

				gomock.InOrder(
					repoMock.EXPECT().AddPurchases(gomock.Any(), gomock.Eq(purchases[0:3])).Return(nil),
				)

				Expect(flusherInstance.Flush(purchases) == nil).To(BeTrue())
			})
		})

		Context("Batch size is 0", func() {
			It("should return all purchases", func() {
				flusherInstance := NewFlusher(0, repoMock)

				Expect(flusherInstance.Flush(purchases)).To(Equal(purchases))
			})
		})

		Context("All add calls return error", func() {
			It("should return all purchases", func() {
				flusherInstance := NewFlusher(2, repoMock)

				gomock.InOrder(
					repoMock.EXPECT().AddPurchases(gomock.Any(), gomock.Eq(purchases[0:2])).Return(errors.New("Opps, error happens")),
					repoMock.EXPECT().AddPurchases(gomock.Any(), gomock.Eq(purchases[2:3])).Return(errors.New("Opps, error happens")),
				)

				Expect(flusherInstance.Flush(purchases)).To(Equal(purchases))
			})
		})

		Context("One of call return error", func() {
			It("should return not saved purchases", func() {
				flusherInstance := NewFlusher(2, repoMock)

				gomock.InOrder(
					repoMock.EXPECT().AddPurchases(gomock.Any(), gomock.Eq(purchases[0:2])).Return(errors.New("Opps, error happens")),
					repoMock.EXPECT().AddPurchases(gomock.Any(), gomock.Eq(purchases[2:3])).Return(nil),
				)

				Expect(flusherInstance.Flush(purchases)).To(Equal(purchases[0:2]))
			})
		})
	})
})
