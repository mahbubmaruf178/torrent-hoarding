package main

func GetTelegramByName(cfg *Config, name string) (*TelegramConfig, bool) {
	for _, tg := range cfg.Telegram {
		if tg.Name == name {
			return &tg, true
		}
	}
	return nil, false
}

type Config struct {
	Telegram          []TelegramConfig `json:"telegram"`
	DeleteAfterUpload bool             `json:"deleteaferupload"`
	DownloadPath      string           `json:"download_path"`
}

type TelegramConfig struct {
	Name         string `json:"name"`
	APIID        int    `json:"api_id"`
	APIHash      string `json:"api_hash"`
	SessionName  string `json:"session_name"`
	DownloadPath string `json:"download_path"`
	ChannelID    string `json:"channelid"`
}
