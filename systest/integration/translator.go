package integration

import (
	"context"

	api "github.com/qiwitech/graceful/systest/api_v1"
	"github.com/qiwitech/graceful/systest/pluto"
)

type APItoPlutoTranslator struct {
	api.APIInterface
	processing pluto.ProcessingInterface
	storage    pluto.StorageInterface
}

// translate api-calls to pluto-calls
func NewAPItoPlutoTranslator(p pluto.ProcessingInterface, s pluto.StorageInterface) api.APIInterface {
	return &APItoPlutoTranslator{
		processing: p,
		storage:    s,
	}
}

func (t *APItoPlutoTranslator) Transfer(ctx context.Context, in *api.TransferRequest) (*api.TransferResponse, error) {
	req := &pluto.TransferRequest{}
	//TODO (ag): fill "req" fields from "in" HERE
	resp, err := t.processing.Transfer(ctx, req)
	if err != nil {
		return nil, err
	}
	out := &api.TransferResponse{
		State:  &api.ActualState{},
		Result: &api.TransferOK{},
		Error:  nil,
	}
	//TODO (ag): fill "out" fields from "resp" HERE
	_ = resp
	return out, nil
}

func (t *APItoPlutoTranslator) UpdateSettings(ctx context.Context, in *api.UpdateSettingsRequest) (*api.UpdateSettingsResponse, error) {
	req := &pluto.UpdateSettingsRequest{}
	//TODO (ag): fill "req" fields from "in" HERE
	resp, err := t.processing.UpdateSettings(ctx, req)
	if err != nil {
		return nil, err
	}
	out := &api.UpdateSettingsResponse{
		State:  &api.ActualState{},
		Result: &api.UpdateSettingsOK{},
		Error:  nil,
	}
	//TODO (ag): fill "out" fields from "resp" HERE
	_ = resp
	return out, nil
}

func (t *APItoPlutoTranslator) GetPrevHash(ctx context.Context, in *api.PrevHashRequest) (*api.PrevHashResponse, error) {
	req := &pluto.PrevHashRequest{}
	//TODO (ag): fill "req" fields from "in" HERE
	resp, err := t.storage.GetPrevHash(ctx, req)
	if err != nil {
		return nil, err
	}
	out := &api.PrevHashResponse{}
	//TODO (ag): fill "out" fields from "resp" HERE
	_ = resp
	return out, nil
}

func (t *APItoPlutoTranslator) GetHistory(ctx context.Context, in *api.HistoryRequest) (*api.HistoryResponse, error) {
	req := &pluto.HistoryRequest{}
	//TODO (ag): fill "req" fields from "in" HERE
	resp, err := t.storage.GetHistory(ctx, req)
	if err != nil {
		return nil, err
	}
	out := &api.HistoryResponse{}
	//TODO (ag): fill "out" fields from "resp" HERE
	_ = resp
	return out, nil
}

func (t *APItoPlutoTranslator) GetStats(ctx context.Context, in *api.StatsRequest) (*api.StatsResponse, error) {
	req := &pluto.StatsRequest{}
	//TODO (ag): fill "req" fields from "in" HERE
	resp, err := t.storage.GetStats(ctx, req)
	if err != nil {
		return nil, err
	}
	out := &api.StatsResponse{}
	//TODO (ag): fill "out" fields from "resp" HERE
	_ = resp
	return out, nil
}

func (t *APItoPlutoTranslator) GetAccounts(ctx context.Context, in *api.AccountsRequest) (*api.AccountsResponse, error) {
	req := &pluto.AccountsRequest{}
	//TODO (ag): fill "req" fields from "in" HERE
	resp, err := t.storage.GetAccounts(ctx, req)
	if err != nil {
		return nil, err
	}
	out := &api.AccountsResponse{}
	//TODO (ag): fill "out" fields from "resp" HERE
	_ = resp
	return out, nil
}

func (t *APItoPlutoTranslator) GetAccountSettings(ctx context.Context, in *api.AccountSettingsRequest) (*api.AccountSettingsResponse, error) {
	req := &pluto.AccountSettingsRequest{}
	//TODO (ag): fill "req" fields from "in" HERE
	resp, err := t.storage.GetAccountSettings(ctx, req)
	if err != nil {
		return nil, err
	}
	out := &api.AccountSettingsResponse{}
	//TODO (ag): fill "out" fields from "resp" HERE
	_ = resp
	return out, nil
}
