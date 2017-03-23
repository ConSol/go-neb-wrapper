package neb

/*
#include "naemon/naemon.h"
*/
import "C"
import (
	"sync"
	"time"
	"unsafe"
)

//Callback defines an function, which will be called by the core
//Return your result in the channel
type Callback func(int, unsafe.Pointer) int
type callbackMapping map[int][]Callback

var usedCallbackMapping = callbackMapping{}
var callbackMutex = sync.Mutex{}

//CallbackTimeout is the duration each callback has to return.
var CallbackTimeout = time.Duration(10) * time.Millisecond

//Generic_Callback this is a mapping function for C. Don't use it.
//export Generic_Callback
func Generic_Callback(callbackType int, data unsafe.Pointer) int {
	returnCode := Ok
	var callbacks []Callback
	var contains bool

	//Test if this call is increasing
	callbackMutex.Lock()
	callbacks, contains = usedCallbackMapping[callbackType]
	callbackMutex.Unlock()
	callbackAmount := len(callbacks)
	if !contains || callbackAmount == 0 {
		return returnCode
	}

	resultChannels := make([]chan int, callbackAmount)
	for i := range resultChannels {
		resultChannels[i] = make(chan int, 1)
	}
	resultList := make([]int, callbackAmount)
	//start all handlers for this callback
	for i, c := range callbacks {
		go func(result chan int, call Callback) {
			result <- call(callbackType, data)
		}(resultChannels[i], c)
	}

	//Wait for all callbacks to signal that they are done and collect the returncodes
	var result int
	for i, c := range resultChannels {
		select {
		case result = <-c:
			resultList[i] = result
		case <-time.After(CallbackTimeout):
			return Error
		}
	}

	//Test the returncodes if any is not OK return the error otherwise OK
	for _, r := range resultList {
		if r != returnCode {
			return r
		}
	}
	return returnCode
}

//AddCallback can be uses to register a function for a certain event
func AddCallback(callbackType int, callback Callback) {
	callbackMutex.Lock()
	usedCallbackMapping[callbackType] = append(usedCallbackMapping[callbackType], callback)
	callbackMutex.Unlock()
}

func initCallbacks() {
	callbackMutex.Lock()
	for callbackType := range usedCallbackMapping {
		RegisterGenericCallback(int64(callbackType))
	}
	callbackMutex.Unlock()
}

func deinitCallbacks() {
	callbackMutex.Lock()
	for callbackType := range usedCallbackMapping {
		RegisterGenericCallback(int64(callbackType))
	}
	callbackMutex.Unlock()
}
