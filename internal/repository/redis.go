package repository

import (
	"encoding/json"
	"time"

	"github.com/Khovanskiy5/yopass/internal/secret/domain"
	"github.com/go-redis/redis/v7"
)

type Redis struct {
	client *redis.Client
}

func NewRedis(url string) (domain.Repository, error) {
	options, err := redis.ParseURL(url)
	if err != nil {
		return nil, err
	}
	client := redis.NewClient(options)
	return &Redis{client}, nil
}

func (r *Redis) Get(key string) (domain.Secret, error) {
	var s domain.Secret
	val, err := r.client.Get(key).Result()
	if err != nil {
		if err == redis.Nil {
			return s, domain.ErrNotFound
		}
		return s, err
	}

	if err := json.Unmarshal([]byte(val), &s); err != nil {
		return s, err
	}

	if s.OneTime {
		_, _ = r.Delete(key)
	}
	return s, nil
}

func (r *Redis) Put(key string, secret domain.Secret) error {
	data, err := secret.ToJSON()
	if err != nil {
		return err
	}
	return r.client.Set(
		key,
		data,
		time.Duration(secret.Expiration)*time.Second,
	).Err()
}

func (r *Redis) Delete(key string) (bool, error) {
	res, err := r.client.Del(key).Result()
	if err != nil {
		if err == redis.Nil {
			return false, nil
		}
		return false, err
	}
	return res > 0, nil
}

func (r *Redis) Status(key string) (bool, error) {
	val, err := r.client.Get(key).Result()
	if err != nil {
		if err == redis.Nil {
			return false, domain.ErrNotFound
		}
		return false, err
	}
	var s domain.Secret
	if err := json.Unmarshal([]byte(val), &s); err != nil {
		return false, err
	}
	return s.OneTime, nil
}
