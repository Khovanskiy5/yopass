package repository

import (
	"encoding/json"

	"github.com/Khovanskiy5/yopass/internal/secret/domain"
	"github.com/bradfitz/gomemcache/memcache"
)

type Memcached struct {
	client *memcache.Client
}

func NewMemcached(server string) domain.Repository {
	return &Memcached{memcache.New(server)}
}

func (m *Memcached) Get(key string) (domain.Secret, error) {
	var s domain.Secret

	item, err := m.client.Get(key)
	if err != nil {
		if err == memcache.ErrCacheMiss {
			return s, domain.ErrNotFound
		}
		return s, err
	}

	if err := json.Unmarshal(item.Value, &s); err != nil {
		return s, err
	}

	if s.OneTime {
		_ = m.client.Delete(key)
	}

	return s, nil
}

func (m *Memcached) Put(key string, secret domain.Secret) error {
	data, err := secret.ToJSON()
	if err != nil {
		return err
	}

	return m.client.Set(&memcache.Item{
		Key:        key,
		Value:      data,
		Expiration: secret.Expiration,
	})
}

func (m *Memcached) Delete(key string) (bool, error) {
	err := m.client.Delete(key)
	if err != nil {
		if err == memcache.ErrCacheMiss {
			return false, nil
		}
		return false, err
	}
	return true, nil
}

func (m *Memcached) Status(key string) (bool, error) {
	item, err := m.client.Get(key)
	if err != nil {
		if err == memcache.ErrCacheMiss {
			return false, domain.ErrNotFound
		}
		return false, err
	}

	var s domain.Secret
	if err := json.Unmarshal(item.Value, &s); err != nil {
		return false, err
	}
	return s.OneTime, nil
}
