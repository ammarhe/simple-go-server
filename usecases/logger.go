package usecases

//
//import (
//	"fmt"
//	"sync"
//	"time"
//)
//
//var (
//	loggedReq = make(map[string]bool)
//)
//
//func logCounter(mu *sync.Mutex) {
//	ticker := time.NewTicker(1 * time.Second)
//	mu.Lock()
//	defer mu.Unlock()
//	<-ticker.C
//	for i, logged := range loggedReq {
//		if !logged {
//			loggedReq[i] = true
//			fmt.Println("logged:", i)
//		}
//	}
//	loggedReq = make(map[string]bool)
//}
