package core

type Service interface {
	Lifecycle
	Start() chan error
	Stop() chan error
}

func NewServiceManager() *ServiceManager {
	s := &ServiceManager{}
	s.services = make(map[interface{}]Service, 0)
	return s
}

type ServiceManager struct {
	services map[interface{}]Service
}

func (sm *ServiceManager) Services() map[interface{}]Service {
	return sm.services
}

func (sm *ServiceManager) Get(key interface{}) Service {
	s := sm.services[key]
	return s
}

func (sm *ServiceManager) Set(key interface{}, s Service) {
	sm.services[key] = s
}
