package cli

import "github.com/99designs/keyring"

// setKeyringValue stores data in keyring with given key.
func (c *CliClient) setKeyringValue(key string, data []byte) error {
	ring, err := keyring.Open(keyring.Config{
		ServiceName: c.KeyRingName,
	})
	if err != nil {
		return err
	}

	return ring.Set(keyring.Item{
		Key:  key,
		Data: data,
	})
}

// getKeyringValue retrieves key's data from keyring.
func (c *CliClient) getKeyringValue(key string) ([]byte, error) {
	ring, err := keyring.Open(keyring.Config{
		ServiceName: c.KeyRingName,
	})

	if err != nil {
		return nil, err
	}

	i, err := ring.Get(key)
	if err != nil {
		return nil, err
	}

	return i.Data, nil
}
