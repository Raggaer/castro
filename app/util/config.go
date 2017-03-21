package util

import (
	"github.com/BurntSushi/toml"
	"time"
)

// CookieConfig struct used for the cookies configuration options
type CookieConfig struct {
	Name     string
	MaxAge   int
	HashKey  string
	BlockKey string
}

// RateLimiterConfig struct used for the rate limiting configration options
type RateLimiterConfig struct {
	Number int64
	Time   time.Duration
}

// CacheConfig struct used for the cache configuration options
type CacheConfig struct {
	Default time.Duration
	Purge   time.Duration
}

// SSLConfig struct used for the ssl configuration options
type SSLConfig struct {
	Enabled bool
	Auto    bool
	Cert    string
	Key     string
}

// MailConfig struct used for the mail configuration options
type MailConfig struct {
	Enabled  bool
	Server   string
	Port     int
	Username string
	Password string
}

// PaygolConfig struct used for the paygol configuration options
type PaygolConfig struct {
	Enabled  bool
	Service  int
	Currency string
	Language string
}

// PayPalConfig struct used for the paypal configuration options
type PayPalConfig struct {
	Enabled   bool
	PublicKey string
	SecretKey string
	Currency  string
	SandBox   bool
}

// Configuration struct used for the main Castro config file TOML file
type Configuration struct {
	Mode      string
	Port      int
	URL       string
	Datapack  string
	Mail      MailConfig
	Captcha   CaptchaConfig
	SSL       SSLConfig
	PayPal    PayPalConfig
	PayGol    PaygolConfig
	Cookies   CookieConfig
	Cache     CacheConfig
	RateLimit RateLimiterConfig
	Custom    map[string]interface{}
}

// Config holds the main configuration file
var Config *Configuration

func init() {
	Config = &Configuration{}
}

// LoadConfig loads the configuration file to the given interface pointer
func LoadConfig(data string, dest interface{}) error {
	// Decode the given file to the given interface
	if _, err := toml.DecodeFile(data, dest); err != nil {
		return err
	}

	return nil
}

// IsDev checks if castro is running on development mode
func (c Configuration) IsDev() bool {
	return c.Mode == "dev"
}

// IsLog checks if castro is running on log mode
func (c Configuration) IsLog() bool {
	return c.Mode == "log"
}
