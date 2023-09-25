package config

import (
	"os"
)

// TODO: Make Markets into config, so I not only handle the environment variables from
// here, but I will also handle the available market instances here.

const (
	TOKEN_TIINGO   string = "CONTROTTO_TIINGO_TOKEN"
	TOKEN_AVANTAGE string = "CONTROTTO_AVANTAGE_TOKEN"
	PORT           string = "CONTROTTO_PORT"
	UNCOMMON_PAIRS string = "CONTROTTO_UNCOMMON_PAIRS"
)

type GetResp struct {
	TiingoAPIToken   string
	AVantageAPIToken string
	Port             string
	UncommonPairs    bool
}

// type IsSetReq struct{}
type IsSetResp struct {
	TiingoAPIToken   bool
	AVantageAPIToken bool
	Port             bool
	UncommonPairs    bool
}

type Service struct{}

func NewConfig() *Service {
	if _, present := os.LookupEnv(PORT); !present {
		os.Setenv(PORT, "8000")
	}
	return &Service{}
}

func (c *Service) Get() *GetResp {
	return &GetResp{
		TiingoAPIToken:   os.Getenv(TOKEN_AVANTAGE),
		AVantageAPIToken: os.Getenv(TOKEN_TIINGO),
		Port:             os.Getenv(PORT),
		UncommonPairs:    os.Getenv(UNCOMMON_PAIRS) != "true",
	}
}
func (c *Service) IsSet() *IsSetResp {
	return &IsSetResp{
		AVantageAPIToken: len(os.Getenv(TOKEN_AVANTAGE)) > 0,
		TiingoAPIToken:   len(os.Getenv(TOKEN_TIINGO)) > 0,
		Port:             len(os.Getenv(PORT)) > 0,
		UncommonPairs:    os.Getenv(UNCOMMON_PAIRS) == "true",
	}
}
