package integration

import (
	"context"
	"testing"
	"time"

	api "github.com/qiwitech/graceful/systest/api_v1"
)

func TestHttp(t *testing.T) {

	client, stop, err := FullServiceChain(&PlutoMockProcessing{}, &PlutoMockStorage{})
	defer stop(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	// гоняем http запросы:
	callAlltest(t, client)
}

func callAlltest(t *testing.T, client api.APIInterface) {

	t.Run("Client call Transfer", func(t *testing.T) {
		ctx, cancel := context.WithTimeout(context.Background(), time.Second)
		defer cancel()
		resp, err := client.Transfer(ctx, &api.TransferRequest{
			Sender: 0x1000,
		})
		if err != nil {
			t.Error(err)
			return
		}
		if resp == nil {
			t.Error("Empty response")
			return
		}
		if resp.GetError() != nil {
			t.Error("Error response")
		}
		if resp.GetResult() == nil {
			t.Logf("TransferResponse: %+v", *resp)
			t.Error("Empty response")
		}
	})

	t.Run("Client call UpdateSettings", func(t *testing.T) {
		ctx, cancel := context.WithTimeout(context.Background(), time.Second)
		defer cancel()
		resp, err := client.UpdateSettings(ctx, &api.UpdateSettingsRequest{})
		if err != nil {
			t.Error(err)
			return
		}
		if resp == nil {
			t.Error("Empty response")
			return
		}
		if resp.GetError() != nil {
			t.Error("Error response")
		}
		if resp.GetResult() == nil {
			t.Logf("UpdateSettingsResponse: %+v", *resp)
			t.Error("Empty response")
		}
	})

	t.Run("Client call GetPrevHash", func(t *testing.T) {
		ctx, cancel := context.WithTimeout(context.Background(), time.Second)
		defer cancel()
		resp, err := client.GetPrevHash(ctx, &api.PrevHashRequest{})
		if err != nil {
			t.Error(err)
			return
		}
		if resp == nil {
			t.Error("Empty response")
			return
		}
	})

	t.Run("Client call GetHistory", func(t *testing.T) {
		ctx, cancel := context.WithTimeout(context.Background(), time.Second)
		defer cancel()
		resp, err := client.GetHistory(ctx, &api.HistoryRequest{})
		if err != nil {
			t.Error(err)
			return
		}
		if resp == nil {
			t.Error("Empty response")
			return
		}
	})

	t.Run("Client call GetStats", func(t *testing.T) {
		ctx, cancel := context.WithTimeout(context.Background(), time.Second)
		defer cancel()
		resp, err := client.GetStats(ctx, &api.StatsRequest{})
		if err != nil {
			t.Error(err)
			return
		}
		if resp == nil {
			t.Error("Empty response")
			return
		}
	})

	t.Run("Client call GetAccounts", func(t *testing.T) {
		ctx, cancel := context.WithTimeout(context.Background(), time.Second)
		defer cancel()
		resp, err := client.GetAccounts(ctx, &api.AccountsRequest{})
		if err != nil {
			t.Error(err)
			return
		}
		if resp == nil {
			t.Error("Empty response")
			return
		}
	})

	t.Run("Client call GetAccountSettings", func(t *testing.T) {
		ctx, cancel := context.WithTimeout(context.Background(), time.Second)
		defer cancel()
		resp, err := client.GetAccountSettings(ctx, &api.AccountSettingsRequest{})
		if err != nil {
			t.Error(err)
			return
		}
		if resp == nil {
			t.Error("Empty response")
			return
		}
	})
}
