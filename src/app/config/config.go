package config

type Config struct {
	Port           int    // CONTROTTO_PORT
	DBPath         string // CONTROTTO_DB_PATH
	UncommonPairs  bool   // CONTROTTO_UNCOMMON_PAIRS
	LoadSampleData bool   // CONTROTTO_LOAD_SAMPLE_DATA
}
