package ecommerce

import (
	"github.com/Kivio-Product/Kivio.Product.Auctions.EcommerceClient/pkg/domain"
	"github.com/Kivio-Product/Kivio.Product.Auctions.EcommerceClient/pkg/repository"
	"github.com/Kivio-Product/Kivio.Product.Auctions.EcommerceClient/pkg/service"
)

type Item = domain.Item
type Customer = domain.Customer

type EcommerceService = service.EcommerceService
type EcommerceCredentialsService = service.EcommerceCredentialsService
type EcommerceCredentials = service.EcommerceCredentials
type IntegrationService = service.IntegrationService
type IntegrationResponse = service.IntegrationResponse
type IntegrationConfigResponse = service.IntegrationConfigResponse

func NewEcommerceService() EcommerceService {
	repo := repository.NewEcommerceRepository()
	return service.NewEcommerceService(repo)
}

func NewEcommerceCredentialsService(integrationService IntegrationService) EcommerceCredentialsService {
	ecommerceService := NewEcommerceService()
	return service.NewEcommerceCredentialsService(integrationService, ecommerceService)
}
