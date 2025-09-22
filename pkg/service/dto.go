package service

import "time"

type IntegrationResponse struct {
	IntegrationID string                      `json:"integrationId"`
	PosID         string                      `json:"posId"`
	Name          string                      `json:"name"`
	Type          string                      `json:"type"`
	Status        string                      `json:"status"`
	LastSync      time.Time                   `json:"lastSync"`
	CreatedAt     time.Time                   `json:"createdAt"`
	Configs       []IntegrationConfigResponse `json:"configs"`
}

type IntegrationConfigResponse struct {
	IntegrationConfigID string `json:"integrationConfigId"`
	Key                 string `json:"key"`
	Value               string `json:"value"`
}
