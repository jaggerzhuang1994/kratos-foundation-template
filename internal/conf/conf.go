package conf

import "github.com/jaggerzhuang1994/kratos-foundation/pkg/config"

func NewBootstrap(c config.Config) (*Bootstrap, error) {
	var bc Bootstrap
	if err := c.Scan(&bc); err != nil {
		return nil, err
	}
	return &bc, nil
}
