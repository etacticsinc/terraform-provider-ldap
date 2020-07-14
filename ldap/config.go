package ldap

type Config struct {
	Server       string
	BindDN       string
	BindPassword string
}

func (c *Config) Client() (interface{}, error) {
	client := &Client{
		Server:       c.Server,
		BindDN:       c.BindDN,
		BindPassword: c.BindPassword,
	}
	return client, nil
}
