package store

import "sync"

type InMemDB struct {
	mu    sync.Mutex
	kvmap map[string]string
}

func (db *InMemDB) Get(key string) (string, bool) {
	val, ok := db.kvmap[key]
	if ok {
		return val, true
	}
	return "", false
}

func (db *InMemDB) Set(key, value string) bool {
	db.mu.Lock()
	defer db.mu.Unlock()
	db.kvmap[key] = value
	return true
}

func (db *InMemDB) Delete(key string) bool {
	db.mu.Lock()
	defer db.mu.Unlock()
	delete(db.kvmap, key)
	return true
}

func InitInMemDb() IDB {
	db := &InMemDB{}
	db.kvmap = make(map[string]string)
	return db
}
