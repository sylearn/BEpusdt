package task

import (
	"strings"
	"testing"
	"time"
)

func TestSolanaGetBlockPayloadSupportsVersionOne(t *testing.T) {
	payload := string(solanaGetBlockPayload(456519571))
	if !strings.Contains(payload, `"maxSupportedTransactionVersion":1`) {
		t.Fatalf("payload does not support Solana v1 transactions: %s", payload)
	}
}

func TestValidateSolanaBlockResponse(t *testing.T) {
	tests := []struct {
		name    string
		body    string
		wantErr bool
	}{
		{name: "valid empty block", body: `{"jsonrpc":"2.0","result":{"blockTime":1760380000,"transactions":[]},"id":1}`},
		{name: "rpc error", body: `{"jsonrpc":"2.0","error":{"code":-32015,"message":"unsupported transaction version"},"id":1}`, wantErr: true},
		{name: "missing result", body: `{"jsonrpc":"2.0","id":1}`, wantErr: true},
		{name: "skipped slot", body: `{"jsonrpc":"2.0","result":null,"id":1}`},
		{name: "invalid json", body: `<html>rate limited</html>`, wantErr: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := validateSolanaBlockResponse([]byte(tt.body))
			if (err != nil) != tt.wantErr {
				t.Fatalf("error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestSolanaRetryDelayIsBounded(t *testing.T) {
	tests := []struct {
		attempt int
		want    time.Duration
	}{
		{attempt: 1, want: 2 * time.Second},
		{attempt: 2, want: 5 * time.Second},
		{attempt: 3, want: 15 * time.Second},
		{attempt: 4, want: 30 * time.Second},
		{attempt: 5, want: time.Minute},
		{attempt: 100, want: time.Minute},
	}
	for _, tt := range tests {
		if got := solanaRetryDelay(tt.attempt); got != tt.want {
			t.Fatalf("attempt %d delay = %s, want %s", tt.attempt, got, tt.want)
		}
	}
}
