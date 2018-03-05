package domain

type Config struct {
	Hosts     map[string]Host `yaml:"hosts"`
	Broadcast string          `yaml:"broadcast"`
}

type Host struct {
	MacAddress string `yaml:"mac"`
	Username   string `yaml:"username"`
	Password   string `yaml:"password"`
}
