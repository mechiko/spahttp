package config

var TomlConfig = []byte(`
# This is a TOML document.
hostname = "127.0.0.1"
hostport = "8888"
browser = ""

[layouts]
timelayout = "2006-01-02T15:04:05-0700"
timelayoutclear = "2006.01.02 15:04:05"
timelayoutday = "2006.01.02"
timelayoututc = "2006-01-02T15:04:05"

[db]
driver = 'sqlite'
file = 'spahttp.db'

`)

type Configuration struct {
	Hostname string `mapstructure:"hostname"`
	HostPort string `mapstructure:"hostport"`
	Browser  string `mapstructure:"browser"`

	Application AppConfiguration      `mapstructure:"application"`
	Layouts     LayoutConfiguration   `mapstructure:"layouts"`
	Db          DatabaseConfiguration `mapstructure:"db"`
}

type LayoutConfiguration struct {
	TimeLayout      string `mapstructure:"timelayout"`
	TimeLayoutClear string `mapstructure:"timelayoutclear"`
	TimeLayoutDay   string `mapstructure:"timelayoutday"`
	TimeLayoutUTC   string `mapstructure:"timelayoututc"`
}

type DatabaseConfiguration struct {
	Name       string `mapstructure:"name"`
	Connection string `mapstructure:"connection"`
	Driver     string `mapstructure:"driver"`
	DbName     string `mapstructure:"dbname"`
	File       string `mapstructure:"file"`
	User       string `mapstructure:"user"`
	Pass       string `mapstructure:"pass"`
	Host       string `mapstructure:"host"`
	Port       string `mapstructure:"port"`
}

type AppConfiguration struct {
	License string `mapstructure:"license"`
}
