package appschema

import (
	"sync"
)

type RequestStore struct {
	Mu       sync.Mutex
	Requests map[string]map[string]int 
}