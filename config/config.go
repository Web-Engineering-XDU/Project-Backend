package config

type (
	Config struct {
		MySQL `yaml:"mysql"`
	}

	MySQL struct {
		Host string `yaml:"host" env-defualt:"localhost"`
		Port string `yaml:"port" env-defualt:"3306"`
		User string `yaml:"user" env-required="true"`
		Password string `yaml:"password" env-required="true"`
	}
)
