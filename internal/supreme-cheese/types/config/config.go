package config

type Config struct {
	MySQLDsn            string              `yaml:"mysql_dsn"`
	ServerPort          string              `yaml:"server_port"`
	SecretKeyJWT        string              `yaml:"secret_key_jwt"`
	Email               *ConfigForSendEmail `yaml:"server_email"`
	Telegram            *Telegram           `yaml:"telegram"`
	PathForBill         string              `yaml:"path_for_bill"`
	PathForSert         string              `yaml:"path_for_sert"`
	UseSSL              bool                `yaml:"use_ssl"`
	CertPathSSL         string              `yaml:"cert_path_ssl"`
	KeyPathSSL          string              `yaml:"key_path_ssl"`
	GifterySecret       string              `yaml:"giftery_secret"`
	ProverkaChekaSecret string              `yaml:"proverka_cheka_secret"`
	YoomoneySecret      string              `yaml:"yoomoney_secret_bearer"`
	SMSLogin            string              `yaml:"sms_login"`
	SMSPassword         string              `yaml:"sms_pass"`
}

type ConfigForSendEmail struct {
	EmailHost  string `yaml:"host"`
	EmailPort  string `yaml:"port"`
	EmailLogin string `yaml:"login"`
	EmailPass  string `yaml:"pass"`
}

type Telegram struct {
	TelegramToken string `yaml:"telegram_token"`
	ChatID        string `yaml:"chat_id"`
	SashaID       string `yaml:"sasha_id"`
}
