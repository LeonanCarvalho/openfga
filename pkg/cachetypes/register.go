package cachetypes

import (
	"encoding/gob"
	"sync"

	"github.com/openfga/openfga/internal/graph"
	"github.com/openfga/openfga/pkg/storage"
)

var once sync.Once

func RegisterCacheTypes() {
	once.Do(func() {
		gob.Register(&storage.ChangelogCacheEntry{})
		gob.Register(&storage.TupleIteratorCacheEntry{})
		gob.Register(&storage.InvalidEntityCacheEntry{})
		gob.Register(&graph.CheckResponseCacheEntry{})
		// Adicione outros tipos aqui
		// gob.Register(&OutroTipo{})
	})
}
