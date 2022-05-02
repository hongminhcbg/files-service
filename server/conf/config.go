package conf

type Config struct {
	RedisDns      string
	GcpBucketName string
	MySqlUrl      string
}

func Load() *Config {
	return &Config{
		RedisDns:      "redis://localhost:6379",
		GcpBucketName: "local-test-file",
		MySqlUrl:      "root:12345678@tcp(localhost:3306)/file_service?charset=utf8mb4&parseTime=True&loc=Local",
	}
}
