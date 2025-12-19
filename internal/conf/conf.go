package conf

import "github.com/go-kratos/kratos/v2/config"

func NewBootstrap(c config.Config) (*Bootstrap, error) {
	var bc Bootstrap
	if err := c.Scan(&bc); err != nil {
		return nil, err
	}
	return &bc, nil
}
