package spec

import (
	"sync"

	"github.com/ranna-go/ranna/pkg/models"
)

type SafeSpecMap struct {
	m *sync.Map
}

func NewSafeSpecMap(m models.SpecMap) (ssm *SafeSpecMap) {
	ssm = &SafeSpecMap{&sync.Map{}}

	for k, v := range m {
		ssm.m.Store(k, v)
	}

	return
}

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

func (ssm *SafeSpecMap) Get(key string) (models.Spec, bool) {
	return ssm.get(key, false)
}

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
