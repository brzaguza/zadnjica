package nocache

import (
	"time"
)

type DB struct{}

func New() (DB, error) { return DB{}, nil }

func (db DB) Close() {}

func (db DB) Set(k string, v interface{}, ttl ...time.Duration) error { return nil }

func (db DB) Get(k string, o interface{}) error { return nil }

func (db DB) GetTTL(k string) (time.Duration, error) { return 0, nil }
