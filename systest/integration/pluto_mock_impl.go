package integration

import (
	"context"

	"github.com/qiwitech/graceful/systest/pluto"
)

// Pluto mock implementation

type PlutoMockProcessing struct {
	pluto.ProcessingInterface
}

type PlutoMockStorage struct {
	pluto.StorageInterface
}

func (s *PlutoMockProcessing) Transfer(context.Context, *pluto.TransferRequest) (*pluto.TransferResponse, error) {
	resp := &pluto.TransferResponse{
		Result: &pluto.TransferOK{},
		Error:  nil,
	}
	return resp, nil
}

func (s *PlutoMockProcessing) UpdateSettings(context.Context, *pluto.UpdateSettingsRequest) (*pluto.UpdateSettingsResponse, error) {
	resp := &pluto.UpdateSettingsResponse{
		Result: &pluto.UpdateSettingsOK{},
		Error:  nil,
	}
	return resp, nil
}

func (s *PlutoMockStorage) GetPrevHash(context.Context, *pluto.PrevHashRequest) (*pluto.PrevHashResponse, error) {
	return &pluto.PrevHashResponse{}, nil
}

func (s *PlutoMockStorage) GetHistory(context.Context, *pluto.HistoryRequest) (*pluto.HistoryResponse, error) {
	return &pluto.HistoryResponse{}, nil
}

func (s *PlutoMockStorage) GetStats(context.Context, *pluto.StatsRequest) (*pluto.StatsResponse, error) {
	return &pluto.StatsResponse{}, nil
}

func (s *PlutoMockStorage) GetAccounts(context.Context, *pluto.AccountsRequest) (*pluto.AccountsResponse, error) {
	return &pluto.AccountsResponse{}, nil
}

func (s *PlutoMockStorage) GetAccountSettings(context.Context, *pluto.AccountSettingsRequest) (*pluto.AccountSettingsResponse, error) {
	return &pluto.AccountSettingsResponse{}, nil
}
