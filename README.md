# Kivio Ecommerce Client

Este módulo contiene toda la lógica para interactuar con servicios de ecommerce externos. Ha sido extraído del repositorio principal para aislar las dependencias de terceros y mantener una arquitectura limpia.

## Estructura

```
ecommerce-client/
├── pkg/
│   ├── domain/           # Entidades de dominio (Item, Customer)
│   ├── client/           # Cliente HTTP para APIs externas
│   ├── repository/       # Capa de persistencia/adaptadores
│   └── service/          # Lógica de negocio y servicios
├── ecommerce.go          # API pública del módulo
├── go.mod               # Definición del módulo
└── README.md            # Este archivo
```

## Uso

### Servicio básico de ecommerce

```go
import "github.com/Kivio-Product/Kivio.Product.Auctions.EcommerceClient"

// Crear instancia del servicio
ecommerceService := ecommerce.NewEcommerceService()

// Usar el servicio
items, err := ecommerceService.GetItems(ctx, apiUrl, apiKey, page, limit)
```

### Servicio de credenciales

```go
// Necesitas implementar IntegrationService en tu aplicación
type MyIntegrationService struct {
    // Tu implementación
}

func (s *MyIntegrationService) GetIntegrationsByPosID(ctx context.Context, posID string) ([]*service.IntegrationResponse, error) {
    // Tu implementación
}

// Crear el servicio de credenciales
integrationService := &MyIntegrationService{}
credentialsService := ecommerce.NewEcommerceCredentialsService(integrationService)

// Obtener credenciales
credentials, err := credentialsService.GetCredentials(ctx, posID)
```

## Interfaces Principales

- `EcommerceService`: Servicio principal para operaciones de ecommerce
- `EcommerceCredentialsService`: Manejo de credenciales y autenticación
- `IntegrationService`: Interfaz que debe implementar el consumidor para acceso a integraciones

## Dependencias

Este módulo depende únicamente del módulo padre para tipos de dominio compartidos. Todas las dependencias HTTP y de terceros están encapsuladas dentro del módulo.