package main

type Config struct {
	Server       string
	BindDN       string
	BindPassword string
}

func (c *Config) Client() (interface{}, error) {
	client := &LdapClient{
		Server:       c.Server,
		BindDN:       c.BindDN,
		BindPassword: c.BindPassword,
	}
	return client, nil
}
