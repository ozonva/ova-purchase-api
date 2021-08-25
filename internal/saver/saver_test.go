package saver

import (
	"github.com/golang/mock/gomock"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/ozonva/ova-purchase-api/internal/mocks"
	"github.com/ozonva/ova-purchase-api/internal/purchase"
	"sync"
	"time"
)

var _ = Describe("Flusher", func() {
	var (
		flusherMock *mocks.MockFlusher
		ctrl        *gomock.Controller
		purchases   []purchase.Purchase
	)

	Describe("Saver tests", func() {
		BeforeEach(func() {
			ctrl = gomock.NewController(GinkgoT())
			flusherMock = mocks.NewMockFlusher(ctrl)
			purchases = []purchase.Purchase{
				purchase.New(), purchase.New(), purchase.New(),
			}
		})
		AfterEach(func() {
			ctrl.Finish()
		})
		Context("Successful flush", func() {
			It("Should flush one time when close saver", func() {
				saverInstance, err := NewSaver(3, flusherMock, 2*time.Hour)

				Expect(err).Should(BeNil())
				wg := sync.WaitGroup{}
				wg.Add(1)

				flusherMock.EXPECT().Flush(gomock.Any()).DoAndReturn(func(ps []purchase.Purchase) []purchase.Purchase {
					wg.Done()
					return nil
				}).Times(1)

				err = saverInstance.Save(purchases[0])
				Expect(err).Should(BeNil())
				err = saverInstance.Save(purchases[1])
				Expect(err).Should(BeNil())
				err = saverInstance.Save(purchases[2])
				Expect(err).Should(BeNil())

				saverInstance.Close()
				wg.Wait()
			})

			It("Should flush one time when timer is ready", func() {
				saverInstance, err := NewSaver(3, flusherMock, 2*time.Second)

				Expect(err).Should(BeNil())
				wg := sync.WaitGroup{}
				wg.Add(1)

				flusherMock.EXPECT().Flush(gomock.Any()).DoAndReturn(func(ps []purchase.Purchase) []purchase.Purchase {
					wg.Done()
					return nil
				}).Times(1)

				err = saverInstance.Save(purchases[0])
				Expect(err).Should(BeNil())
				err = saverInstance.Save(purchases[1])
				Expect(err).Should(BeNil())
				err = saverInstance.Save(purchases[2])
				Expect(err).Should(BeNil())
				wg.Wait()
			})

			It("Should flush two times when close and timer", func() {
				saverInstance, err := NewSaver(3, flusherMock, 2*time.Second)

				Expect(err).Should(BeNil())
				wg := sync.WaitGroup{}
				wg.Add(1)

				flusherMock.EXPECT().Flush(gomock.Any()).DoAndReturn(func(ps []purchase.Purchase) []purchase.Purchase {
					wg.Done()
					return nil
				}).Times(2)

				err = saverInstance.Save(purchases[0])
				Expect(err).Should(BeNil())
				err = saverInstance.Save(purchases[1])
				Expect(err).Should(BeNil())
				err = saverInstance.Save(purchases[2])
				Expect(err).Should(BeNil())

				wg.Wait() // Wait timer

				wg.Add(1)
				err = saverInstance.Save(purchases[0])
				Expect(err).Should(BeNil())
				saverInstance.Close()
				wg.Wait() // Wait close flush
			})
		})

		Context("Errors happens", func() {
			It("Should return error when timeout is zero", func() {
				saverInstance, err := NewSaver(3, flusherMock, 0*time.Second)

				Expect(err).Should(Equal(TimeoutNotValidError))
				Expect(saverInstance).Should(BeNil())
			})

			It("Should return error when capacity is zero", func() {
				saverInstance, err := NewSaver(0, flusherMock, 5*time.Hour)

				Expect(err).Should(Equal(CapacityNotValidError))
				Expect(saverInstance).Should(BeNil())
			})

			It("Should return error when added purchases more than capacity", func() {
				saverInstance, err := NewSaver(1, flusherMock, 5*time.Hour)

				Expect(err).Should(BeNil())

				err = saverInstance.Save(purchases[0])
				Expect(err).Should(BeNil())
				err = saverInstance.Save(purchases[1])
				Expect(err).Should(Equal(CapacityReachedError))
			})

			It("Should flush failed records second time", func() {

				saverInstance, err := NewSaver(3, flusherMock, 3*time.Second)

				Expect(err).Should(BeNil())

				wg := sync.WaitGroup{}
				wg.Add(1)

				flusherMock.EXPECT().Flush(gomock.Eq(purchases[:1])).DoAndReturn(func(ps []purchase.Purchase) []purchase.Purchase {
					wg.Done()
					return ps
				}).Times(1)

				flusherMock.EXPECT().Flush(gomock.Eq(purchases)).DoAndReturn(func(ps []purchase.Purchase) []purchase.Purchase {
					wg.Done()
					return nil
				}).Times(1)

				err = saverInstance.Save(purchases[0])
				Expect(err).Should(BeNil())
				wg.Wait()

				err = saverInstance.Save(purchases[1])
				Expect(err).Should(BeNil())
				err = saverInstance.Save(purchases[2])
				Expect(err).Should(BeNil())
				wg.Add(1)
				saverInstance.Close()
				wg.Wait()
			})
		})
	})
})
