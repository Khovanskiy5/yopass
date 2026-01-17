package config

import (
	"strings"

	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

type Config struct {
	Address             string
	Port                int
	Database            string
	AssetPath           string
	MaxLength           int
	Memcached           string
	MetricsPort         int
	Redis               string
	TLSCert             string
	TLSKey              string
	ForceOneTimeSecrets bool
	CORSAllowOrigin     string
	DisableUpload       bool
	PrefetchSecret      bool
	DisableFeatures     bool
	NoLanguageSwitcher  bool
	TrustedProxies      []string
	PrivacyNoticeURL    string
	ImprintURL          string
	AllowedExpirations  []int
}

func Load() (*Config, error) {
	pflag.String("address", "", "listen address (default 0.0.0.0)")
	pflag.Int("port", 1337, "listen port")
	pflag.String("database", "memcached", "database backend ('memcached' or 'redis')")
	pflag.String("asset-path", "public", "path to the assets folder")
	pflag.Int("max-length", 5242880, "max length of encrypted secret")
	pflag.String("memcached", "localhost:11211", "memcached address")
	pflag.Int("metrics-port", -1, "metrics server listen port")
	pflag.String("redis", "redis://localhost:6379/0", "Redis URL")
	pflag.String("tls-cert", "", "path to TLS certificate")
	pflag.String("tls-key", "", "path to TLS key")
	pflag.Bool("force-onetime-secrets", false, "reject non onetime secrets from being created")
	pflag.String("cors-allow-origin", "*", "Access-Control-Allow-Origin")
	pflag.Bool("disable-upload", false, "disable the /file upload endpoints")
	pflag.Bool("prefetch-secret", true, "Display information that the secret might be one time use")
	pflag.Bool("disable-features", false, "disable features")
	pflag.Bool("no-language-switcher", false, "disable the language switcher in the UI")
	pflag.StringSlice("trusted-proxies", []string{}, "trusted proxy IP addresses or CIDR blocks for X-Forwarded-For header validation")
	pflag.String("privacy-notice-url", "", "URL to privacy notice page")
	pflag.String("imprint-url", "", "URL to imprint/legal notice page")
	pflag.IntSlice("allowed-expirations", []int{3600, 86400, 604800}, "allowed expiration times in seconds")

	viper.SetEnvPrefix("yopass")
	viper.AutomaticEnv()
	viper.SetEnvKeyReplacer(strings.NewReplacer("-", "_"))
	if err := viper.BindPFlags(pflag.CommandLine); err != nil {
		return nil, err
	}

	pflag.Parse()

	return &Config{
		Address:             viper.GetString("address"),
		Port:                viper.GetInt("port"),
		Database:            viper.GetString("database"),
		AssetPath:           viper.GetString("asset-path"),
		MaxLength:           viper.GetInt("max-length"),
		Memcached:           viper.GetString("memcached"),
		MetricsPort:         viper.GetInt("metrics-port"),
		Redis:               viper.GetString("redis"),
		TLSCert:             viper.GetString("tls-cert"),
		TLSKey:              viper.GetString("tls-key"),
		ForceOneTimeSecrets: viper.GetBool("force-onetime-secrets"),
		CORSAllowOrigin:     viper.GetString("cors-allow-origin"),
		DisableUpload:       viper.GetBool("disable-upload"),
		PrefetchSecret:      viper.GetBool("prefetch-secret"),
		DisableFeatures:     viper.GetBool("disable-features"),
		NoLanguageSwitcher:  viper.GetBool("no-language-switcher"),
		TrustedProxies:      viper.GetStringSlice("trusted-proxies"),
		PrivacyNoticeURL:    viper.GetString("privacy-notice-url"),
		ImprintURL:          viper.GetString("imprint-url"),
		AllowedExpirations:  viper.GetIntSlice("allowed-expirations"),
	}, nil
}
