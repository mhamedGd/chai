package ecs

import (
	"reflect"
	"sync"
)

var componentIdMutex sync.Mutex
var registeredComponents = make(map[reflect.Type]componentId, maxComponentId)
var inavlidComponentId componentId = 0
var componentRegistryCounter componentId = 1

func name(t any) componentId {
	componentIdMutex.Lock()
	defer componentIdMutex.Unlock()

	typeof := reflect.TypeOf(t)
	compId, ok := registeredComponents[typeof]

	if !ok {
		compId = componentRegistryCounter
		registeredComponents[typeof] = compId
		componentRegistryCounter++
	}

	return compId
}
