package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"time"
)

type EcommerceClient interface {
	GetApiKey(username, password, tokenUrl string) (string, error)
	GetItems(baseUrl, apiKey string, page, limit int, publishedStatus bool, filters map[string]string) ([]byte, error)
	GetItemByID(baseUrl, apiKey, itemId string) ([]byte, error)
	GetCustomers(baseUrl, apiKey string) ([]byte, error)
	GetCustomerByID(baseUrl, apiKey, id string) ([]byte, error)
	GetOrders(baseUrl, apiKey string) ([]byte, error)
	GetAllItems(baseUrl, apiKey string) ([]byte, error)
	GetStores(baseUrl, apiKey string) ([]byte, error)
	CreateCustomer(baseUrl, apiKey string, customerData []byte) ([]byte, error)
	CreateBillingAddress(baseUrl, apiKey string, customerID int, addressData []byte) ([]byte, error)
	CreateShippingAddress(baseUrl, apiKey string, customerID int, addressData []byte) ([]byte, error)
	CreateShoppingCartItem(baseUrl, apiKey string, cartItemData []byte) ([]byte, error)
	CreateOrder(baseUrl, apiKey string, orderData []byte) ([]byte, error)
	CountEcommerceItems(baseUrl, apiKey string, filters map[string]string) (int64, error)
	UpdateOrderItemPrice(baseUrl, apiKey string, orderID, itemID int, orderItemData []byte) error
	UpdateOrder(baseUrl, apiKey string, orderID int, orderData []byte) error
}

type ecommerceClient struct {
	httpClient *http.Client
}

func NewEcommerceClient() EcommerceClient {
	return &ecommerceClient{
		httpClient: &http.Client{
			Timeout: time.Second * 30,
		},
	}
}

func (c *ecommerceClient) GetApiKey(username, password, tokenUrl string) (string, error) {
	payload := map[string]interface{}{
		"guest":       true,
		"username":    username,
		"password":    password,
		"remember_me": true,
	}

	body, err := json.Marshal(payload)
	if err != nil {
		return "", fmt.Errorf("failed to marshal payload: %w", err)
	}

	req, err := http.NewRequest("POST", tokenUrl, bytes.NewBuffer(body))
	if err != nil {
		return "", fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("accept", "text/plain")
	req.Header.Set("Content-Type", "application/json-patch+json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("failed to get API key, status code: %d", resp.StatusCode)
	}

	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to read response body: %w", err)
	}

	var result map[string]interface{}
	if err := json.Unmarshal(respBody, &result); err != nil {
		return "", fmt.Errorf("failed to unmarshal response: %w", err)
	}

	apiKey, ok := result["access_token"].(string)
	if !ok {
		return "", fmt.Errorf("API key not found in response")
	}

	return apiKey, nil
}

func (c *ecommerceClient) GetItems(baseUrl, apiKey string, page, limit int, publishedStatus bool, filters map[string]string) ([]byte, error) {
	url := fmt.Sprintf("%s/api/products?Page=%d&Limit=%d&PublishedStatus=%t&Name=%s", baseUrl, page, limit, publishedStatus, url.PathEscape(filters["name"]))

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", apiKey))

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to get items, status code: %d", resp.StatusCode)
	}

	return ioutil.ReadAll(resp.Body)
}

func (c *ecommerceClient) CountEcommerceItems(baseUrl, apiKey string, filters map[string]string) (int64, error) {

	url := fmt.Sprintf("%s/api/products/count?PublishedStatus=true&Name=%s", baseUrl, url.PathEscape(filters["name"]))

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return 0, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", apiKey))

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return 0, fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return 0, fmt.Errorf("failed to get items, status code: %d", resp.StatusCode)
	}

	var result struct {
		Count int64 `json:"count"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return 0, fmt.Errorf("failed to decode JSON: %w", err)
	}

	fmt.Println("Total items:", result.Count)
	return result.Count, nil
}

func (c *ecommerceClient) GetItemByID(baseUrl, apiKey, itemId string) ([]byte, error) {
	url := fmt.Sprintf("%s/api/products/%s", baseUrl, itemId)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", apiKey))

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to get item, status code: %d", resp.StatusCode)
	}

	return ioutil.ReadAll(resp.Body)
}

func (c *ecommerceClient) GetCustomers(baseUrl, apiKey string) ([]byte, error) {
	req, err := http.NewRequest("GET", fmt.Sprintf("%s/api/customers", baseUrl), nil)
	if err != nil {
		return nil, fmt.Errorf("error creating request: %w", err)
	}

	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", apiKey))
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error making request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	return ioutil.ReadAll(resp.Body)
}

func (c *ecommerceClient) GetCustomerByID(baseUrl, apiKey, id string) ([]byte, error) {
	req, err := http.NewRequest("GET", fmt.Sprintf("%s/customers/%s", baseUrl, id), nil)
	if err != nil {
		return nil, fmt.Errorf("error creating request: %w", err)
	}

	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", apiKey))
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error making request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNotFound {
		return nil, fmt.Errorf("customer not found")
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	return ioutil.ReadAll(resp.Body)
}

func (c *ecommerceClient) GetOrders(baseUrl, apiKey string) ([]byte, error) {
	req, err := http.NewRequest("GET", fmt.Sprintf("%s/api/orders", baseUrl), nil)
	if err != nil {
		return nil, fmt.Errorf("error creating request: %w", err)
	}

	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", apiKey))
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error making request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	return ioutil.ReadAll(resp.Body)
}

func (c *ecommerceClient) GetAllItems(baseUrl, apiKey string) ([]byte, error) {
	var allProducts []map[string]interface{}
	page := 1
	limit := 100

	for {
		url := fmt.Sprintf("%s/api/products?Page=%d&Limit=%d", baseUrl, page, limit)

		req, err := http.NewRequest("GET", url, nil)
		if err != nil {
			return nil, fmt.Errorf("failed to create request: %w", err)
		}
		req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", apiKey))

		resp, err := c.httpClient.Do(req)
		if err != nil {
			return nil, fmt.Errorf("failed to send request: %w", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			return nil, fmt.Errorf("failed to get items, status code: %d", resp.StatusCode)
		}

		bodyBytes, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return nil, fmt.Errorf("failed to read response body: %w", err)
		}

		var response struct {
			Products []map[string]interface{} `json:"products"`
		}
		if err := json.Unmarshal(bodyBytes, &response); err != nil {
			return nil, fmt.Errorf("failed to unmarshal response: %w", err)
		}

		if len(response.Products) == 0 {
			break
		}

		allProducts = append(allProducts, response.Products...)

		if len(response.Products) < limit {
			break
		}

		page++
	}

	finalResponse := map[string]interface{}{
		"products": allProducts,
	}

	return json.Marshal(finalResponse)
}

func (c *ecommerceClient) CreateCustomer(baseUrl, apiKey string, customerData []byte) ([]byte, error) {
	url := fmt.Sprintf("%s/api/customers", baseUrl)
	fmt.Printf("[HTTP] POST %s\n", url)

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(customerData))
	if err != nil {
		fmt.Printf("[HTTP] ERROR: Failed to create request: %v\n", err)
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", apiKey))
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		fmt.Printf("[HTTP] ERROR: Failed to send request: %v\n", err)
		return nil, fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	fmt.Printf("[HTTP] Response Status: %d\n", resp.StatusCode)

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		fmt.Printf("[HTTP] ERROR: Unexpected status code: %d\n", resp.StatusCode)
		return nil, fmt.Errorf("failed to create customer, status code: %d", resp.StatusCode)
	}

	return ioutil.ReadAll(resp.Body)
}

func (c *ecommerceClient) CreateBillingAddress(baseUrl, apiKey string, customerID int, addressData []byte) ([]byte, error) {
	url := fmt.Sprintf("%s/api/customers/%d/billingaddress", baseUrl, customerID)
	fmt.Printf("[HTTP] POST %s\n", url)

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(addressData))
	if err != nil {
		fmt.Printf("[HTTP] ERROR: Failed to create request: %v\n", err)
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", apiKey))
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		fmt.Printf("[HTTP] ERROR: Failed to send request: %v\n", err)
		return nil, fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	fmt.Printf("[HTTP] Response Status: %d\n", resp.StatusCode)

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		fmt.Printf("[HTTP] ERROR: Unexpected status code: %d\n", resp.StatusCode)
		return nil, fmt.Errorf("failed to create billing address, status code: %d", resp.StatusCode)
	}

	return ioutil.ReadAll(resp.Body)
}

func (c *ecommerceClient) CreateShippingAddress(baseUrl, apiKey string, customerID int, addressData []byte) ([]byte, error) {
	url := fmt.Sprintf("%s/api/customers/%d/shippingaddress", baseUrl, customerID)
	fmt.Printf("[HTTP] POST %s\n", url)

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(addressData))
	if err != nil {
		fmt.Printf("[HTTP] ERROR: Failed to create request: %v\n", err)
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", apiKey))
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		fmt.Printf("[HTTP] ERROR: Failed to send request: %v\n", err)
		return nil, fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	fmt.Printf("[HTTP] Response Status: %d\n", resp.StatusCode)

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		fmt.Printf("[HTTP] ERROR: Unexpected status code: %d\n", resp.StatusCode)
		return nil, fmt.Errorf("failed to create shipping address, status code: %d", resp.StatusCode)
	}

	return ioutil.ReadAll(resp.Body)
}

func (c *ecommerceClient) CreateShoppingCartItem(baseUrl, apiKey string, cartItemData []byte) ([]byte, error) {
	url := fmt.Sprintf("%s/api/shopping_cart_items", baseUrl)
	fmt.Printf("[HTTP] POST %s\n", url)

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(cartItemData))
	if err != nil {
		fmt.Printf("[HTTP] ERROR: Failed to create request: %v\n", err)
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", apiKey))
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		fmt.Printf("[HTTP] ERROR: Failed to send request: %v\n", err)
		return nil, fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	fmt.Printf("[HTTP] Response Status: %d\n", resp.StatusCode)

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		fmt.Printf("[HTTP] ERROR: Unexpected status code: %d\n", resp.StatusCode)
		return nil, fmt.Errorf("failed to create shopping cart item, status code: %d", resp.StatusCode)
	}

	return ioutil.ReadAll(resp.Body)
}

func (c *ecommerceClient) CreateOrder(baseUrl, apiKey string, orderData []byte) ([]byte, error) {
	url := fmt.Sprintf("%s/api/orders", baseUrl)
	fmt.Printf("[HTTP] POST %s\n", url)

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(orderData))
	if err != nil {
		fmt.Printf("[HTTP] ERROR: Failed to create request: %v\n", err)
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", apiKey))
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		fmt.Printf("[HTTP] ERROR: Failed to send request: %v\n", err)
		return nil, fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	fmt.Printf("[HTTP] Response Status: %d\n", resp.StatusCode)

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		fmt.Printf("[HTTP] ERROR: Unexpected status code: %d\n", resp.StatusCode)
		return nil, fmt.Errorf("failed to create order, status code: %d", resp.StatusCode)
	}

	return ioutil.ReadAll(resp.Body)
}

func (c *ecommerceClient) UpdateOrderItemPrice(baseUrl, apiKey string, orderID, itemID int, orderItemData []byte) error {
	url := fmt.Sprintf("%s/api/orders/%d/items/%d", baseUrl, orderID, itemID)
	fmt.Printf("[HTTP] PUT %s\n", url)
	fmt.Printf("[HTTP] Payload: %s\n", string(orderItemData))

	req, err := http.NewRequest("PUT", url, bytes.NewBuffer(orderItemData))
	if err != nil {
		fmt.Printf("[HTTP] ERROR: Failed to create request: %v\n", err)
		return fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", apiKey))
	req.Header.Set("Content-Type", "application/json-patch+json")
	req.Header.Set("accept", "text/plain")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		fmt.Printf("[HTTP] ERROR: Failed to send request: %v\n", err)
		return fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	fmt.Printf("[HTTP] Response Status: %d\n", resp.StatusCode)

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusNoContent {
		bodyBytes, _ := ioutil.ReadAll(resp.Body)
		fmt.Printf("[HTTP] ERROR: Unexpected status code: %d, Body: %s\n", resp.StatusCode, string(bodyBytes))
		return fmt.Errorf("failed to update order item price, status code: %d", resp.StatusCode)
	}

	fmt.Printf("[HTTP] SUCCESS: Order item price updated\n")
	return nil
}

func (c *ecommerceClient) GetStores(baseUrl, apiKey string) ([]byte, error) {
	url := fmt.Sprintf("%s/api/stores", baseUrl)
	fmt.Printf("[HTTP] GET %s\n", url)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		fmt.Printf("[HTTP] ERROR: Failed to create request: %v\n", err)
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", apiKey))
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		fmt.Printf("[HTTP] ERROR: Failed to send request: %v\n", err)
		return nil, fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	fmt.Printf("[HTTP] Response Status: %d\n", resp.StatusCode)

	if resp.StatusCode != http.StatusOK {
		bodyBytes, _ := ioutil.ReadAll(resp.Body)
		fmt.Printf("[HTTP] ERROR: Unexpected status code: %d, Body: %s\n", resp.StatusCode, string(bodyBytes))
		return nil, fmt.Errorf("failed to get stores, status code: %d", resp.StatusCode)
	}

	return ioutil.ReadAll(resp.Body)
}

func (c *ecommerceClient) UpdateOrder(baseUrl, apiKey string, orderID int, orderData []byte) error {
	url := fmt.Sprintf("%s/api/orders/%d", baseUrl, orderID)
	fmt.Printf("[HTTP] PUT %s\n", url)
	fmt.Printf("[HTTP] Payload: %s\n", string(orderData))

	req, err := http.NewRequest("PUT", url, bytes.NewBuffer(orderData))
	if err != nil {
		fmt.Printf("[HTTP] ERROR: Failed to create request: %v\n", err)
		return fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", apiKey))
	req.Header.Set("Content-Type", "application/json-patch+json")
	req.Header.Set("accept", "text/plain")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		fmt.Printf("[HTTP] ERROR: Failed to send request: %v\n", err)
		return fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	fmt.Printf("[HTTP] Response Status: %d\n", resp.StatusCode)

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusNoContent {
		bodyBytes, _ := ioutil.ReadAll(resp.Body)
		fmt.Printf("[HTTP] ERROR: Unexpected status code: %d, Body: %s\n", resp.StatusCode, string(bodyBytes))
		return fmt.Errorf("failed to update order, status code: %d", resp.StatusCode)
	}

	fmt.Printf("[HTTP] SUCCESS: Order updated\n")
	return nil
}
