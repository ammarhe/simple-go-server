package services

import (
	"context"
	"log"
	"os"
	"sync"
	"time"
	"verve-task/infrastructure"
)

type MemoryStore struct {
	Mutex        *sync.Mutex
	LoggedReqIds map[int]bool
	logFile      *os.File
	redisClient  *infrastructure.RedisClient
	ctx          context.Context
}

func NewMemoryStore(redisClient *infrastructure.RedisClient, ctx context.Context) *MemoryStore {
	redisClient.Client.Expire(ctx, "uniqueIDs", time.Minute) // reset every minute

	return &MemoryStore{
		LoggedReqIds: make(map[int]bool),
		Mutex:        &sync.Mutex{},
		redisClient:  redisClient,
		ctx:          ctx,
	}
}

func (s *MemoryStore) LogRequest(reqID int) error {
	s.Mutex.Lock()
	defer s.Mutex.Unlock()
	s.LoggedReqIds[reqID] = true
	_, err := s.redisClient.Client.SAdd(s.ctx, "uniqueIDs", reqID).Result() // id is added only if it is not exists, no need to check if it is existing or not
	if err != nil {
		return err
	}
	return nil
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
			//pipe := s.redisClient.Pipeline()
			//keys := s.redisClient.Get(s.ctx, "uniqueIDs")
			//for _, key := range keys.Val() {
			//	pipe.Del(key)
			//}
			//_, exec := pipe.Exec(s.ctx)
			//if exec != nil {
			//	return
			//}
			s.Mutex.Unlock()
			s.Reset()
		}
	}

}
