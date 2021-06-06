package spec

import (
	"sync"

	"github.com/ranna-go/ranna/pkg/models"
)

// SafeSpecMap wraps a sync.Map containing
// the actual loaded spec map so that it can
// be accessed and updated across multiple
// go routines.
type SafeSpecMap struct {
	m *sync.Map
}

// NewSafeSpecMap initializes a new SafeSpecMap
// wrapping the passed models.SpecMap m.
func NewSafeSpecMap(m models.SpecMap) (ssm *SafeSpecMap) {
	ssm = &SafeSpecMap{&sync.Map{}}
	ssm.storeMap(m)
	return
}

// get tries to retrieve a spec by the given key.
// If the retrieved spec is an alias (utilizes the 'use'
// spec property), get is performed with the value
// of 'use' as key and isAlias as true to prevent alias
// cycles.
func (ssm *SafeSpecMap) get(key string, isAlias bool) (s models.Spec, ok bool) {
	_sp, _ := ssm.m.Load(key)
	sp, ok := _sp.(*models.Spec)
	if !ok {
		return
	}

	if sp.Use != "" {
		if isAlias {
			ok = false
			return
		}
		return ssm.get(sp.Use, true)
	}

	s = *sp
	return
}

// Get tries to retrieve a Spec from the internal
// spec map by the given key. If no spec was found,
// a nil map and false is returned.
//
// This also resolved aliases (see 'use' spec property).
func (ssm *SafeSpecMap) Get(key string) (models.Spec, bool) {
	return ssm.get(key, false)
}

// GetSnapshot initializes a new SpecMap from the current
// state of the internal spec map and returns it.
func (ssm *SafeSpecMap) GetSnapshot() (m models.SpecMap) {
	m = make(models.SpecMap)
	ssm.m.Range(func(_key, _value interface{}) bool {
		key, kOk := _key.(string)
		value, vOk := _value.(*models.Spec)
		if kOk && vOk {
			m[key] = value
		}
		return true
	})
	return m
}

// Update purges the internal spec map and sets the
// values from the provided new spec map m.
func (ssm *SafeSpecMap) Update(m models.SpecMap) {
	ssm.m.Range(func(key, value interface{}) bool {
		ssm.m.Delete(key)
		return true
	})
	ssm.storeMap(m)
}

// storeMap iterates through all key-value paris
// of the given SpecMap and sets them to the internal
// spec map.
func (ssm *SafeSpecMap) storeMap(m models.SpecMap) {
	for k, v := range m {
		ssm.m.Store(k, v)
	}
}
