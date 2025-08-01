package config

type HTTPServer struct {
	Address string `yaml:"address"`
}

type Config struct {
	Env		  string `yaml:"env"`
	StoragePath string `yaml:"storage_path"`
	HTTPServer `yaml:"http_server"`
}