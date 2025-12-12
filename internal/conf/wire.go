package conf

import "github.com/google/wire"

var ProviderSet = wire.NewSet(
	NewFileSource,
	NewConsulSource,
	NewBootstrap,
)
