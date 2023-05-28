package config

type (
	Config struct {
		MySQL   `yaml:"mysql"`
	}
	MySQL struct {
		Host     string `yaml:"host" env-default:"localhost"`
		Port     string `yaml:"port" env-default:"3306"`
		User     string `yaml:"user" env-required:"true"`
		Password string `yaml:"password" env-required:"true"`
	}
)
