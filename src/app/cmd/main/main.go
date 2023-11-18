package main

import (
	"context"
	"fmt"
	"log"
	"time"

	env "github.com/Netflix/go-env"
	"github.com/redis/go-redis/v9"
)

type Environment struct {
	Home string `env:"HOME"`

	Jenkins struct {
		BuildId     *string `env:"BUILD_ID"`
		BuildNumber int     `env:"BUILD_NUMBER"`
		Ci          bool    `env:"CI"`
	}

	Node struct {
		ConfigCache *string `env:"npm_config_cache,NPM_CONFIG_CACHE"`
	}

	Extras env.EnvSet

	Duration     time.Duration `env:"TYPE_DURATION"`
	DefaultValue string        `env:"MISSING_VAR,default=default_value"`
	// RequiredValue string        `env:"IM_REQUIRED,required=true"`

	Redis struct {
		Addr     string `env:"REDIS_ADDR,required=true"`
		Password string `env:"REDIS_PASSWORD,required=true"`
		Db       int    `env:"REDIS_DB,required=true"`
	}
}

func main() {
	var environment Environment
	es, err := env.UnmarshalFromEnviron(&environment)
	if err != nil {
		log.Fatal(err)
	}
	// Remaining environment variables.
	environment.Extras = es

	log.Println("Hello, World!")
	ExampleClient(&environment)
}

func ExampleClient(environment *Environment) {
	var ctx = context.Background()
	rdb := redis.NewClient(&redis.Options{
		Addr:     environment.Redis.Addr,
		Password: environment.Redis.Password, // "" no password set
		DB:       environment.Redis.Db,       // 0 use default DB
	})

	const ch = "dev:ch:00"
	sub := rdb.Subscribe(ctx, ch)
	iface, err := sub.Receive(ctx)
	if err != nil {
		// handle error
		panic(err)
	}

	// Should be *Subscription, but others are possible if other actions have been
	// taken on sub since it was created.
	switch iface.(type) {
	case *redis.Subscription:
		// subscribe succeeded
		fmt.Println("Subscription")
	case *redis.Message:
		// received first message
		fmt.Println("message")
	case *redis.Pong:
		// pong received
		fmt.Println("pong")
	default:
		// handle error
		panic(err)
	}

	for rMsg := range sub.Channel() {
		fmt.Println("rMsg", rMsg)
	}
}
