package rng

import (
	"context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"time"
)

type SimpleClient struct {
	api               RNGClient
	MaxProcessingTime time.Duration
}

type Config struct {
	Host              string
	Port              string
	MaxProcessingTime time.Duration
}

func newClient(host, port string) (RNGClient, error) {
	addr := host + ":" + port
	conn, err := grpc.Dial(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))

	if err != nil {
		return nil, err
	}

	return NewRNGClient(conn), nil
}

func New(cfg *Config) (Client, error) {
	var err error

	client := &SimpleClient{}
	client.api, err = newClient(cfg.Host, cfg.Port)

	if err != nil {
		return nil, err
	}

	client.MaxProcessingTime = cfg.MaxProcessingTime * time.Millisecond

	return client, nil
}

func (c *SimpleClient) Rand(max uint64) (rand uint64, err error) {
	ctx, cancel := context.WithTimeout(context.Background(), c.MaxProcessingTime)
	defer cancel()

	in := &RandRequest{Max: []uint64{max}}
	resp, err := c.api.Rand(ctx, in)

	if err != nil {
		return 0, err
	}

	return resp.Result[0], nil
}

func (c *SimpleClient) RandSlice(maxSlice []uint64) (rand []uint64, err error) {
	ctx, cancel := context.WithTimeout(context.Background(), c.MaxProcessingTime)
	defer cancel()

	in := &RandRequest{Max: maxSlice}
	resp, err := c.api.Rand(ctx, in)

	if err != nil {
		return rand, err
	}

	return resp.Result, nil
}

func (c *SimpleClient) RandFloat() (float64, error) {
	ctx, cancel := context.WithTimeout(context.Background(), c.MaxProcessingTime)
	defer cancel()

	in := &RandRequestFloat{Max: uint64(1)}
	resp, err := c.api.RandFloat(ctx, in)

	if err != nil {
		return 0, err
	}

	return resp.Result[0], nil
}

func (c *SimpleClient) RandFloatSlice(count int) ([]float64, error) {
	ctx, cancel := context.WithTimeout(context.Background(), c.MaxProcessingTime)
	defer cancel()

	in := &RandRequestFloat{Max: uint64(count)}
	resp, err := c.api.RandFloat(ctx, in)

	if err != nil {
		return nil, err
	}

	return resp.Result, nil
}
