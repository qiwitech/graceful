package integration

import (
	"context"
	"testing"
	"time"

	"github.com/golang/protobuf/proto"
	api "github.com/qiwitech/graceful/systest/api_v1"
)

func BenchmarkAPI(t *testing.B) {

	apiClient, stop, err := FullServiceChain(&PlutoMockProcessing{}, &PlutoMockStorage{})
	defer stop(context.Background())
	if err != nil {
		t.Fatal(err)
	}

	// гоняем http запросы:
	t.ResetTimer()
	for i := 0; i < t.N; i++ {
		callAllbench(t, apiClient)
	}
	t.StopTimer() // чтоб не учитывать время остановки всего
}

func callAllbench(t *testing.B, client api.APIInterface) {

	var (
		ctx  context.Context
		resp proto.Message
		err  error
	)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	resp, err = client.Transfer(ctx, &api.TransferRequest{
		Sender: 0x1000,
	})
	cancel()
	if err != nil {
		t.Fatal(err)
	} else if resp == nil {
		t.Fatal("Empty response")
	}

	ctx, cancel = context.WithTimeout(context.Background(), time.Second)
	resp, err = client.UpdateSettings(ctx, &api.UpdateSettingsRequest{})
	cancel()
	if err != nil {
		t.Fatal(err)
	} else if resp == nil {
		t.Fatal("Empty response")
	}

	ctx, cancel = context.WithTimeout(context.Background(), time.Second)
	resp, err = client.GetPrevHash(ctx, &api.PrevHashRequest{})
	cancel()
	if err != nil {
		t.Fatal(err)
	} else if resp == nil {
		t.Fatal("Empty response")
	}

	ctx, cancel = context.WithTimeout(context.Background(), time.Second)
	resp, err = client.GetHistory(ctx, &api.HistoryRequest{})
	cancel()
	if err != nil {
		t.Fatal(err)
	} else if resp == nil {
		t.Fatal("Empty response")
	}

	ctx, cancel = context.WithTimeout(context.Background(), time.Second)
	resp, err = client.GetStats(ctx, &api.StatsRequest{})
	cancel()
	if err != nil {
		t.Fatal(err)
	} else if resp == nil {
		t.Fatal("Empty response")
	}

	ctx, cancel = context.WithTimeout(context.Background(), time.Second)
	resp, err = client.GetAccounts(ctx, &api.AccountsRequest{})
	cancel()
	if err != nil {
		t.Fatal(err)
	} else if resp == nil {
		t.Fatal("Empty response")
	}

	ctx, cancel = context.WithTimeout(context.Background(), time.Second)
	resp, err = client.GetAccountSettings(ctx, &api.AccountSettingsRequest{})
	cancel()
	if err != nil {
		t.Fatal(err)
	} else if resp == nil {
		t.Fatal("Empty response")
	}

}
