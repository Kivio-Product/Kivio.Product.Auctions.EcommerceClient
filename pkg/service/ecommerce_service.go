package service

import (
	"context"
	"fmt"

	"github.com/Kivio-Product/Kivio.Product.Auctions.EcommerceClient/pkg/domain"
	"github.com/Kivio-Product/Kivio.Product.Auctions.EcommerceClient/pkg/repository"
)

type EcommerceService interface {
	GetItems(ctx context.Context, apiUrl, apiKey string, page, limit int) ([]domain.Item, error)
	GetItemsRaw(ctx context.Context, apiUrl, apiKey string, page, limit int, publishedStatus bool) ([]byte, error)
	GetItemByID(ctx context.Context, id, apiUrl, apiKey string) (*domain.Item, error)
	GetItemByIDRaw(ctx context.Context, id, apiUrl, apiKey string) ([]byte, error)
	GetCustomers(ctx context.Context, apiUrl, apiKey string) ([]domain.Customer, error)
	GetCustomerByID(ctx context.Context, id, apiUrl, apiKey string) (*domain.Customer, error)
	GetApiKey(ctx context.Context, username, password, tokenUrl string) (string, error)
	UpdateItemStock(ctx context.Context, apiUrl, apiKey, itemId string, newStock int) error
	GetAllItemsRaw(ctx context.Context, apiUrl, apiKey string) ([]byte, error)
	CreateEcommerceCustomer(ctx context.Context, apiUrl, apiKey string, customerData []byte) ([]byte, error)
	CreateEcommerceBillingAddress(ctx context.Context, apiUrl, apiKey string, customerID int, addressData []byte) ([]byte, error)
	CreateEcommerceShippingAddress(ctx context.Context, apiUrl, apiKey string, customerID int, addressData []byte) ([]byte, error)
	CreateEcommerceShoppingCartItem(ctx context.Context, apiUrl, apiKey string, cartItemData []byte) ([]byte, error)
	CreateEcommerceOrder(ctx context.Context, apiUrl, apiKey string, orderData []byte) ([]byte, error)
}

type ecommerceService struct {
	repo repository.EcommerceRepository
}

func NewEcommerceService(repo repository.EcommerceRepository) EcommerceService {
	return &ecommerceService{
		repo: repo,
	}
}

func (s *ecommerceService) GetItems(ctx context.Context, apiUrl, apiKey string, page, limit int) ([]domain.Item, error) {
	return s.repo.GetItems(apiUrl, apiKey, page, limit)
}

func (s *ecommerceService) GetItemsRaw(ctx context.Context, apiUrl, apiKey string, page, limit int, publishedStatus bool) ([]byte, error) {
	return s.repo.GetItemsRaw(apiUrl, apiKey, page, limit, publishedStatus)
}

func (s *ecommerceService) GetItemByID(ctx context.Context, id, apiUrl, apiKey string) (*domain.Item, error) {
	return s.repo.GetItemByID(apiUrl, apiKey, id)
}

func (s *ecommerceService) GetItemByIDRaw(ctx context.Context, id, apiUrl, apiKey string) ([]byte, error) {
	fmt.Println("Fetching item by ID:", id, "from API URL:", apiUrl)
	return s.repo.GetItemByIDRaw(apiUrl, apiKey, id)
}

func (s *ecommerceService) GetCustomers(ctx context.Context, apiUrl, apiKey string) ([]domain.Customer, error) {
	return s.repo.GetCustomers(apiUrl, apiKey)
}

func (s *ecommerceService) GetCustomerByID(ctx context.Context, id, apiUrl, apiKey string) (*domain.Customer, error) {
	return s.repo.GetCustomerByID(id, apiUrl, apiKey)
}

func (s *ecommerceService) GetApiKey(ctx context.Context, username, password, tokenUrl string) (string, error) {
	return s.repo.GetApiKey(username, password, tokenUrl)
}

func (s *ecommerceService) UpdateItemStock(ctx context.Context, apiUrl, apiKey, itemId string, newStock int) error {
	return s.repo.UpdateItemStock(apiUrl, apiKey, itemId, int64(newStock))
}

func (s *ecommerceService) GetAllItemsRaw(ctx context.Context, apiUrl, apiKey string) ([]byte, error) {
	return s.repo.GetAllItemsRaw(apiUrl, apiKey)
}

func (s *ecommerceService) CreateEcommerceCustomer(ctx context.Context, apiUrl, apiKey string, customerData []byte) ([]byte, error) {
	fmt.Printf("Creating customer in ecommerce with data: %s\n", string(customerData))

	respBody, err := s.repo.CreateCustomer(apiUrl, apiKey, customerData)
	if err != nil {
		return nil, fmt.Errorf("failed to create customer: %w", err)
	}

	return respBody, nil
}

func (s *ecommerceService) CreateEcommerceBillingAddress(ctx context.Context, apiUrl, apiKey string, customerID int, addressData []byte) ([]byte, error) {
	fmt.Printf("Creating billing address for customer %d in ecommerce with data: %s\n", customerID, string(addressData))

	respBody, err := s.repo.CreateBillingAddress(apiUrl, apiKey, customerID, addressData)
	if err != nil {
		return nil, fmt.Errorf("failed to create billing address: %w", err)
	}

	return respBody, nil
}

func (s *ecommerceService) CreateEcommerceShippingAddress(ctx context.Context, apiUrl, apiKey string, customerID int, addressData []byte) ([]byte, error) {
	fmt.Printf("Creating shipping address for customer %d in ecommerce with data: %s\n", customerID, string(addressData))

	respBody, err := s.repo.CreateShippingAddress(apiUrl, apiKey, customerID, addressData)
	if err != nil {
		return nil, fmt.Errorf("failed to create shipping address: %w", err)
	}

	return respBody, nil
}

func (s *ecommerceService) CreateEcommerceShoppingCartItem(ctx context.Context, apiUrl, apiKey string, cartItemData []byte) ([]byte, error) {
	fmt.Printf("Creating shopping cart item in ecommerce with data: %s\n", string(cartItemData))

	respBody, err := s.repo.CreateShoppingCartItem(apiUrl, apiKey, cartItemData)
	if err != nil {
		return nil, fmt.Errorf("failed to create shopping cart item: %w", err)
	}

	return respBody, nil
}

func (s *ecommerceService) CreateEcommerceOrder(ctx context.Context, apiUrl, apiKey string, orderData []byte) ([]byte, error) {
	fmt.Printf("Creating order in ecommerce with data: %s\n", string(orderData))

	respBody, err := s.repo.CreateOrder(apiUrl, apiKey, orderData)
	if err != nil {
		return nil, fmt.Errorf("failed to create order: %w", err)
	}

	return respBody, nil
}
