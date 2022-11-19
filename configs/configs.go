package configs

import "time"

type Config struct {
	Debug         bool   `mapstructure:"DEBUG"`
	UseBrowser    string `mapstructure:"USE_BROWSER"`
	WebdriverPort uint16 `mapstructure:"WEBDRIVER_PORT"`

	FirefoxBinaryPath string   `mapstructure:"FIREFOX_BINARY_PATH"`
	FirefoxArgs       []string `mapstructure:"FIREFOX_ARGS"`
	GeckodriverPath   string   `mapstructure:"GECKODRIVER_PATH"`

	ChromeBinaryPath string   `mapstructure:"CHROME_BINARY_PATH"`
	ChromeArgs       []string `mapstructure:"CHROME_ARGS"`
	ChromedriverPath string   `mapstructure:"CHROMEDRIVER_PATH"`

	ServerListenAddr   string        `mapstructure:"SERVER_LISTEN_ADDR"`
	ServerReadTimeout  time.Duration `mapstructure:"SERVER_READ_TIMEOUT"`
	ServerWriteTimeout time.Duration `mapstructure:"SERVER_WRITE_TIMEOUT"`
}

var Configuration = new(Config)
