package crud

import (
	"errors"
	"time"
)

type service struct {
	repo1 *repository
}

func NewService(repo *repository) *service {
	return &service{repo1: repo}
}

func (s *service) select1() ([]dual, error) {
	return s.repo1.select1()
}

func (s *service) selectAll() ([]Category, error) {
	return s.repo1.selectAll()
}

func (s *service) selectAllProducts(name string) ([]Product, error) {
	return s.repo1.selectAllProducts(name)
}

func (s *service) insert(data *Category) error {
	return s.repo1.insert(data)
}

func (s *service) selectByID(id int) (*Category, error) {
	return s.repo1.selectByID(id)
}

func (s *service) update(category *Category) error {
	return s.repo1.update(category)
}

func (s *service) delete(id int) error {
	return s.repo1.delete(id)
}

func (s *service) insertShoppingCart(co *CheckoutItem1) error {
	return s.repo1.insertShoppingCart(co)
}

func (s *service) checkout(items []CheckoutItem1) (*Transaction, error) {
	return s.repo1.createTransaction(items)
}

type ReportService interface {
	GetTodayReport() (*TodayReportResponse, error)
}

func (s *service) GetTodayReport(startDate, endDate time.Time) (*TodayReportResponse, error) {
	
	if startDate.After(endDate) {
		return nil, errors.New("Start_date cannot be after end_date")
	}

	if endDate.Sub(startDate) > 365*24*time.Hour {
		return nil, errors.New("Date range cannot exceed 1 year")
	}
	totalRevenue, totalTransaksi, bestProduct, err := s.repo1.GetTodaySummary(startDate, endDate)
	if err != nil {
		return nil, err
	}

	return &TodayReportResponse{
		TotalRevenue:   totalRevenue,
		TotalTransaksi: totalTransaksi,
		ProdukTerlaris: bestProduct,
	}, nil
}
