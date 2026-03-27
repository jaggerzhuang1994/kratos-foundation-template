package job

import (
	"github.com/google/wire"
	"github.com/jaggerzhuang1994/kratos-foundation-template/internal/job/cronjob"
	"github.com/jaggerzhuang1994/kratos-foundation-template/internal/job/globaljob"
	"github.com/jaggerzhuang1994/kratos-foundation-template/internal/job/kafka_consumer"
)

var ProviderSet = wire.NewSet(
	globaljob.NewJob,
	cronjob.NewJob,
	kafka_consumer.NewConsumer,
)
