package main

import (
	"context"
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

	ticker := time.NewTicker(5 * time.Second)

	// https://blog.boot.dev/golang/range-over-ticker-in-go-with-immediate-first-tick/
	for t := time.Now(); true; t = <-ticker.C {
		isoDate := t.Format(time.RFC3339Nano)

		const queue = "dev:queue:00"
		const ch = "dev:ch:00"
		rCmd := rdb.RPush(ctx, queue, isoDate)
		if rCmd.Err() != nil {
			log.Println("rCmd.Err()", rCmd.Err())
			continue
		}

		cmd := rdb.Publish(ctx, ch, isoDate)
		if cmd.Err() != nil {
			log.Println("cmd.Err()", cmd.Err())
		}
	}

	// ticker.Stop()
}
