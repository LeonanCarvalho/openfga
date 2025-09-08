package storage

import (
	"context"
	"testing"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/require"
)

func TestRedisCache(t *testing.T) {
       t.Run("corrupted_data", func(t *testing.T) {
	       prefix := "test:" + t.Name() + ":"
	       cache := NewRedisCache[string]("localhost:6379", "", prefix)
	       defer cache.Stop()

	       key := "corrupted"
	       client := redis.NewClient(&redis.Options{Addr: "localhost:6379"})
	       client.Set(context.Background(), prefix+key, []byte{0x01, 0x02, 0x03}, 5*time.Second)

	       got := cache.Get(key)
	       require.Equal(t, "", got)
       })

       t.Run("type_mismatch", func(t *testing.T) {
	       prefix := "test:" + t.Name() + ":"
	       cache := NewRedisCache[string]("localhost:6379", "", prefix)
	       defer cache.Stop()

	       key := "typemismatch"
	       cacheInt := NewRedisCache[int]("localhost:6379", "", prefix)
	       cacheInt.Set(key, 42, 5*time.Second)

	       got := cache.Get(key)
	       require.Equal(t, "", got)
       })

       t.Run("tuple_iterator_cache_entry", func(t *testing.T) {
	       prefix := "test:" + t.Name() + ":"
	       cache := NewRedisCache[*TupleIteratorCacheEntry]("localhost:6379", "", prefix)
	       defer cache.Stop()

	       key := "tuple_entry"
	       entry := &TupleIteratorCacheEntry{
		       Tuples: []*TupleRecord{{ObjectID: "obj1", Relation: "rel1"}},
		       LastModified: time.Now(),
	       }
	       cache.Set(key, entry, 5*time.Second)

	       got := cache.Get(key)
	       require.NotNil(t, got)
	       require.Equal(t, entry.Tuples[0].ObjectID, got.Tuples[0].ObjectID)
	       require.Equal(t, entry.Tuples[0].Relation, got.Tuples[0].Relation)
       })

       t.Run("set_get_delete", func(t *testing.T) {
	       prefix := "test:" + t.Name() + ":"
	       cache := NewRedisCache[string]("localhost:6379", "", prefix)
	       defer cache.Stop()

	       key := "mykey"
	       value := "myvalue"
	       cache.Set(key, value, 5*time.Second)

	       got := cache.Get(key)
	       require.Equal(t, value, got)

	       cache.Delete(key)
	       got = cache.Get(key)
	       require.Equal(t, "", got)
       })

       t.Run("ttl", func(t *testing.T) {
	       prefix := "test:" + t.Name() + ":"
	       cache := NewRedisCache[string]("localhost:6379", "", prefix)
	       defer cache.Stop()

	       key := "ttlkey"
	       value := "ttlvalue"
	       cache.Set(key, value, 1*time.Second)

	       time.Sleep(2 * time.Second)
	       got := cache.Get(key)
	       require.Equal(t, "", got)
       })

       t.Run("set_and_get", func(t *testing.T) {
	       prefix := "test:" + t.Name() + ":"
	       cache := NewRedisCache[string]("localhost:6379", "", prefix)
	       defer cache.Stop()
	       key := "key"
	       value := "value"
	       cache.Set(key, value, 1*time.Second)
	       result := cache.Get(key)
	       require.Equal(t, value, result)
       })

       t.Run("set_and_get_more_than_one_year", func(t *testing.T) {
	       prefix := "test:" + t.Name() + ":"
	       cache := NewRedisCache[string]("localhost:6379", "", prefix)
	       defer cache.Stop()
	       key := "key"
	       value := "value"
	       cache.Set(key, value, time.Duration(1<<32-1))
	       result := cache.Get(key)
	       require.Equal(t, value, result)
       })

       t.Run("negative_ttl_ignored", func(t *testing.T) {
	       prefix := "test:" + t.Name() + ":"
	       cache := NewRedisCache[string]("localhost:6379", "", prefix)
	       defer cache.Stop()
	       key := "key"
	       value := "value"
	       cache.Set(key, value, -2)
	       result := cache.Get(key)
	       require.Equal(t, "", result, "Valor não deve ser gravado se o TTL for negativo")
       })

       t.Run("stop_multiple_times", func(t *testing.T) {
	       prefix := "test:" + t.Name() + ":"
	       cache := NewRedisCache[string]("localhost:6379", "", prefix)
	       cache.Stop()
	       cache.Stop()
       })

       t.Run("stop_concurrently", func(t *testing.T) {
	       prefix := "test:" + t.Name() + ":"
	       cache := NewRedisCache[string]("localhost:6379", "", prefix)
	       done := make(chan struct{}, 2)
	       go func() {
		       cache.Stop()
		       done <- struct{}{}
	       }()
	       go func() {
		       cache.Stop()
		       done <- struct{}{}
	       }()
	       <-done
	       <-done
       })
}
