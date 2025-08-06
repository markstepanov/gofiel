package config

var GlobalServerConfig Config

type Config struct {
	Port          string `yaml:"port"`
	Basedir       string `yaml:"base-dir"`
	Adminusername string `yaml:"admin-username"`
	Adminpassword string `yaml:"admin-password"`
}
