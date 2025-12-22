package app

import (
	"github.com/Everestown/Outfit_backend/internal/core/module"
	"sync"
)

type ModuleRegistry struct {
	modules map[string]module.Module
	mu      sync.RWMutex
}

func NewModuleRegistry() *ModuleRegistry {
	return &ModuleRegistry{
		modules: make(map[string]module.Module),
	}
}

func (r *ModuleRegistry) RegisterModule(m module.Module) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if err := m.Init(); err != nil {
		return err
	}

	r.modules[m.GetName()] = m
	return nil
}

func (r *ModuleRegistry) GetModule(name string) (module.Module, bool) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	m, exists := r.modules[name]
	return m, exists
}

func (r *ModuleRegistry) GetAllModules() []module.Module {
	r.mu.RLock()
	defer r.mu.RUnlock()

	modules := make([]module.Module, 0, len(r.modules))
	for _, m := range r.modules {
		modules = append(modules, m)
	}
	return modules
}

func (r *ModuleRegistry) CloseAll() {
	r.mu.Lock()
	defer r.mu.Unlock()

	for _, m := range r.modules {
		err := m.Close()
		if err != nil {
			return
		}
	}
	r.modules = make(map[string]module.Module)
}
