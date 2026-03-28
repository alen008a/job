package grpc

import (
	"siteLetterJob/config"
	"siteLetterJob/internal/glog"

	"context"
	"time"

	"github.com/processout/grpc-go-pool"
	"google.golang.org/grpc"
)

var g *transport

type transport struct {
	grpcAddress string
	factory     func() (*grpc.ClientConn, error)
	initCap     int
	cap         int
	maxLife     time.Duration
	pool        *grpcpool.Pool
}

func InitTransport() error {
	var err error
	opts := []grpc.DialOption{grpc.WithInsecure()}

	g = &transport{
		grpcAddress: config.GetConfig().VenueGRPCAddr,
		initCap:     500,
		cap:         1000,
		maxLife:     time.Minute * 5,
	}

	g.factory = func() (*grpc.ClientConn, error) {
		return grpc.Dial(g.grpcAddress, opts...)
	}

	// 创建一个初始化为 1000， 容量为 8000 的连接池
	g.pool, err = grpcpool.New(g.factory, g.initCap, g.cap, g.maxLife)
	if err != nil {
		glog.Errorf("GRPC |VenueGRPCAddr=%s |err=%v", g.grpcAddress, err)
		return err
	}

	return nil
}

func GetConn() (*grpcpool.ClientConn, error) {
	return g.pool.Get(context.Background())
}
