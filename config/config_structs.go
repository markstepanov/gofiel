package config

type Config struct {
	Port          int    `yaml:"port"`
	BaseDir       string `yaml:"base-dir"`
	AdminUsername string `yaml:"admin-username"`
	AdminPassword string `yaml:"admin-password"`
}
