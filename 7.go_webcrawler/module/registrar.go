package module

import (
	"demo/7.go_webcrawler/errors"
	"fmt"
	"sync"
)

type Registrar interface {
	Register(module Module) (bool, error)                 //注册
	Unregister(mid MID) (bool, error)                     //取消注册
	Get(moduleType Type) (Module, error)                  //根据类型获取一个组件实例
	GetAllByType(moduleType Type) (map[MID]Module, error) //根据类型获取对应组件所有的实例
	GetAll() map[MID]Module                               //获取所有的组件实例
	Clear()                                               //清除所有组件的注册纪录
}

func NewRegistrar() Registrar {
	return &myRegistrar{
		moduleTypeMap: map[Type]map[MID]Module{},
	}
}

type myRegistrar struct {
	moduleTypeMap map[Type]map[MID]Module //组件类型与对应组件实例的映射
	rwlock        sync.RWMutex            //组件注册读写锁
}

func (registrar *myRegistrar) Register(module Module) (bool, error) {
	if module == nil {
		return false, errors.NewIllegalParameterError("nil module instatnce")
	}
	mid := module.ID()
	parts, err := SplitMID(mid)
	if err != nil {
		return false, err
	}
	moduleType := legalLetterTypeMap[parts[0]]
	if !CheckType(moduleType, module) {
		return false, errors.NewIllegalParameterError(fmt.Sprintf("incorrect module type: %s", moduleType))
	}
	registrar.rwlock.Lock()
	defer registrar.rwlock.Unlock()

	modules := registrar.moduleTypeMap[moduleType]
	if modules == nil {
		modules = map[MID]Module{}
	}
	if _, ok := modules[mid]; ok {
		return false, nil
	}
	modules[mid] = module
	registrar.moduleTypeMap[moduleType] = modules
	return true, nil
}

func (registrar *myRegistrar) Unregister(mid MID) (bool, error) {
	parts, err := SplitMID(mid)
	if err != nil {
		return false, err
	}
	moduleType := legalLetterTypeMap[parts[0]]
	var deleted bool
	registrar.rwlock.Lock()
	defer registrar.rwlock.Unlock()
	if modules, ok := registrar.moduleTypeMap[moduleType]; ok {
		if _, ok := modules[mid]; ok {
			delete(modules, mid)
			deleted = true
		}
	}
	return deleted, nil
}

func (registrar *myRegistrar) Get(moduleType Type) (Module, error) {
	modules, err := registrar.GetAllByType(moduleType)
	if err != nil {
		return nil, err
	}
	minScore := uint64(0)
	var selectdModule Module
	for _, module := range modules {
		SetScore(module)
		if err != nil {
			return nil, err
		}
		score := module.Score()
		if minScore == 0 || score < minScore {
			selectdModule = module
			minScore = score
		}
	}
	return selectdModule, nil
}

func (registrar *myRegistrar) GetAllByType(moduleType Type) (map[MID]Module, error) {
	if !LegalType(moduleType) {
		return nil, errors.NewIllegalParameterError(fmt.Sprintf("illegal module type: %s", moduleType))
	}
	registrar.rwlock.Lock()
	defer registrar.rwlock.Unlock()

	modules := registrar.moduleTypeMap[moduleType]
	if len(modules) == 0 {
		return nil, ErrNotFoundModuleInstance
	}
	result := map[MID]Module{}
	for mid, module := range modules {
		result[mid] = module
	}
	return result, nil
}

func (registrar *myRegistrar) GetAll() map[MID]Module {
	result := map[MID]Module{}

	registrar.rwlock.Lock()
	defer registrar.rwlock.Unlock()

	for _, modules := range registrar.moduleTypeMap {
		for mid, module := range modules {
			result[mid] = module
		}
	}
	return result
}

func (registrar *myRegistrar) Clear() {
	registrar.rwlock.Lock()
	defer registrar.rwlock.Unlock()

	registrar.moduleTypeMap = map[Type]map[MID]Module{}
}
