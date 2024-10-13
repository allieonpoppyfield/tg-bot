package cache

import (
	"context"
	"fmt"

	re "github.com/go-redis/redis/v8"
)

type Cache struct {
	Client *re.Client
}

func New(ctx context.Context) (*Cache, error) {
	client := re.NewClient(&re.Options{
		Addr:     "redis:6379", // адрес и порт Redis
		Password: "",           // пароль, если он есть
		DB:       0,            // номер базы данных
	})

	pong, err := client.Ping(ctx).Result()
	if err != nil {
		return nil, err
	}
	fmt.Println("Соединение с Redis установлено:", pong)

	return &Cache{Client: client}, nil
}
