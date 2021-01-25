package memory

import (
	"encoding/json"
	"errors"
	"sync"
	"time"
)

func NewStore() *Store {

	return &Store{
		db: new(sync.Map),
	}
}

type Store struct {
	db *sync.Map
}

type Data struct {
	value      []byte
	Time       time.Time
	expiration int64
}

func (s Store) Set(key string, val interface{}, expiration int) error {
	bytes, err := json.Marshal(val)
	if err != nil {
		return err
	}
	s.db.Store(key, Data{value: bytes, Time: time.Now(), expiration: int64(expiration)})
	return nil
}

func (s Store) Get(key string, rel interface{}) error {
	value, ok := s.db.Load(key)
	if ok {
		data := value.(Data)
		bytes := data.value
		if data.Time.Unix()+data.expiration >= time.Now().Unix() {
			err := json.Unmarshal(bytes, rel)
			if err != nil {
				return err
			}
			return nil
		}
		s.Del(key)
	}
	return errors.New("Not Find Record")
}

func (s Store) Del(key string) error {
	s.db.Delete(key)
	return nil
}

func (s Store) Close() error {
	return nil
}
