package cli

import "github.com/99designs/keyring"

// setJWTToken stores token in keyring
func (c *CliClient) setJWTToken(tokenString string) error {
	ring, err := keyring.Open(keyring.Config{
		ServiceName: c.KeyRingName,
	})
	if err != nil {
		return err
	}

	return ring.Set(keyring.Item{
		Key:  "jwt_token",
		Data: []byte(tokenString),
	})
}

// getJWTToken return storerd jwt token
func (c *CliClient) getJWTToken() (string, error) {
	ring, err := keyring.Open(keyring.Config{
		ServiceName: c.KeyRingName,
	})

	if err != nil {
		return "", err
	}

	i, err := ring.Get("jwt_token")
	if err != nil {
		return "", err
	}

	return string(i.Data), nil
}
