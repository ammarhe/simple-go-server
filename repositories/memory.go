package repositories

import (
	"log"
	"os"
	"sync"
	"time"
)

type MemoryStore struct {
	Mutex        *sync.Mutex
	LoggedReqIds map[int]bool
	logFile      *os.File
}

func NewMemoryStore() *MemoryStore {
	return &MemoryStore{
		LoggedReqIds: make(map[int]bool),
		Mutex:        &sync.Mutex{},
	}
}

func (s *MemoryStore) LogRequest(reqID int) {
	s.Mutex.Lock()
	defer s.Mutex.Unlock()
	s.LoggedReqIds[reqID] = true
}
func (s *MemoryStore) IsLogged(reqID int) bool {
	s.Mutex.Lock()
	defer s.Mutex.Unlock()
	_, ok := s.LoggedReqIds[reqID]
	return ok
}
func (s *MemoryStore) Reset() {
	s.Mutex.Lock()
	defer s.Mutex.Unlock()
	s.LoggedReqIds = make(map[int]bool)
}

func (s *MemoryStore) OpenLogFile() error {
	f, err := os.OpenFile("logs", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	s.logFile = f
	log.SetOutput(s.logFile)
	return err
}
func (s *MemoryStore) CloseLogFile() error {
	err := s.logFile.Close()
	return err
}

func (s *MemoryStore) LogCounter() {
	ticker := time.NewTicker(1 * time.Minute)
	defer ticker.Stop()
	defer func() {
		err := s.CloseLogFile()
		if err != nil {
			log.Fatalf("failed to close log file: %v", err)
		}
	}()
	for {
		select {
		case <-ticker.C:
			s.Mutex.Lock()
			log.Printf("Unique request count: %d\n", len(s.LoggedReqIds))
			s.Mutex.Unlock()
			s.Reset()
		}
	}

}
