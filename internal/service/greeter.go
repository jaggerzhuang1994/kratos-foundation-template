package service

import (
	"context"
	"fmt"
	"time"

	v1 "github.com/jaggerzhuang1994/kratos-foundation-template/api/helloworld/v1"
	"github.com/jaggerzhuang1994/kratos-foundation-template/internal/biz/biz1"
	"github.com/jaggerzhuang1994/kratos-foundation-template/internal/biz/biz2"
	"github.com/jaggerzhuang1994/kratos-foundation/pkg/app_info"
	"github.com/jaggerzhuang1994/kratos-foundation/pkg/component/log"
	"github.com/jaggerzhuang1994/kratos-foundation/pkg/component/metric"
	"github.com/jaggerzhuang1994/kratos-foundation/pkg/component/tracing"
)

// GreeterService is a greeter service.
type GreeterService struct {
	v1.UnimplementedGreeterServer
	tracing *tracing.Tracing
	metrics *metric.Metrics
	log     *log.Log

	biz1 *biz1.Biz
	biz2 *biz2.Biz
}

// NewGreeterService new a greeter service.
func NewGreeterService(
	metrics *metric.Metrics,
	tracing *tracing.Tracing,
	log *log.Log,
	biz1 *biz1.Biz,
	biz2 *biz2.Biz,
) *GreeterService {
	return &GreeterService{tracing: tracing, metrics: metrics, log: log, biz1: biz1, biz2: biz2}
}

// SayHello implements helloworld.GreeterServer.
func (s *GreeterService) SayHello(ctx context.Context, in *v1.HelloRequest) (*v1.HelloReply, error) {
	s.log.WithContext(ctx).Info("Counter", s.metrics.AddCounter(ctx, "hello1", 1))
	s.log.WithContext(ctx).Info("Gauge", s.metrics.RecordGauge(ctx, "hello2", time.Now().Unix()))
	s.log.WithContext(ctx).Info("Histogram", s.metrics.RecordHistogram(ctx, "hello3", 6))

	s.tracing.SimpleTrace(ctx, "input", func(ctx context.Context) {
		s.log.WithContext(ctx).Info(app_info.FromContext(ctx))
	})
	s.log.WithContext(ctx).NewHelper().Info("hello2")
	//err := utils.SleepWithContext(ctx, time.Second)
	//if err != nil {
	//	return nil, err
	//}
	err := s.biz2.Set(ctx, "1")
	if err != nil {
		return nil, err
	}
	val, _ := s.biz2.Get(ctx)
	s.log.WithContext(ctx).Info("get", val)

	fmt.Println("biz1.get", s.biz1.Get())

	return &v1.HelloReply{Message: string("Hello " + val)}, nil
}
