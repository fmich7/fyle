package config

// NewTestingConfig returns config that is used for testing
func NewTestingConfig() *Config {
	return &Config{
		ServerPort:      ":0",
		JWTsecretKey:    "supersecret",
		UploadsLocation: "/uploads",
	}
}
