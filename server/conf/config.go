package conf
type Config struct {
	RedisDns string
	GcpBucketName string
}

func Load() *Config {
	return &Config{
		RedisDns:      "redis://localhost:6379",
		GcpBucketName: "local-test-file",
	}
}

