package pnl

import "testing"

func TestPair_Calculate(t *testing.T) {
	type fields struct {
		ID           PairID
		BaseAsset    Asset
		QuoteAsset   Asset
		Trades       []Trade
		Calculations Calculations
	}
	tests := []struct {
		name    string
		fields  fields
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tp := &Pair{
				ID:           tt.fields.ID,
				BaseAsset:    tt.fields.BaseAsset,
				QuoteAsset:   tt.fields.QuoteAsset,
				Trades:       tt.fields.Trades,
				Calculations: tt.fields.Calculations,
			}
			if err := tp.Calculate(); (err != nil) != tt.wantErr {
				t.Errorf("Calculate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
