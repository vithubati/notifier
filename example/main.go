package main

import (
	"context"
	"github.com/vithubati/go-notifier/config"
	"github.com/vithubati/go-notifier/service"
	"log"
	"net"
	"net/http"
	"sync"
	"time"
)

func newConfig() *config.Config {
	return &config.Config{
		Notifier: config.Notifier{
			Webhook:          true,
			ConnString:       "root:password@/notifier?parseTime=true",
			DeliveryInterval: 5 * time.Second,
			Migrations:       true,
		},
		Trace: true,
	}
}

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		s, err := Notifier(newConfig())
		if err != nil {
			log.Fatalf("Notifier() error = %v", err)
			return
		}
		if err := s.KickOff(ctx); err != nil {
			log.Fatalf("Notifier() error = %v", err)
			return
		}
		return
	}()
	wg.Wait()
}

func Notifier(cfg *config.Config) (service.Service, error) {
	var netTransport = &http.Transport{
		DialContext: (&net.Dialer{
			Timeout: 5 * time.Second,
		}).DialContext,
		TLSHandshakeTimeout: 5 * time.Second,
	}
	c := &http.Client{
		Timeout:   time.Second * 10,
		Transport: netTransport,
	}

	s, err := service.New(service.Opts{
		DeliveryInterval: cfg.Notifier.DeliveryInterval,
		ConnString:       cfg.Notifier.ConnString,
		Client:           c,
		Migrations:       cfg.Notifier.Migrations,
		WebhookEnabled:   cfg.Notifier.Webhook,
	})
	if err != nil {
		return nil, err
	}
	return s, nil
}