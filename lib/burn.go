package lib

import (
	"fmt"
	"log"
	"runtime"
	"sync"
	"time"

	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/health/grpc_health_v1"
)

const (
	DefaultGRPCWorkers = runtime.NumCPU()
)

type GTargeter struct {
	Target     string
	IsEtcd     bool
	MethodName string
	Args       []*interface{}
}

type Burner struct {
	conn    *grpc.ClientConn
	workers uint64
	stopch  chan struct{}
}

func (t *GTargeter) GenBurner() (burner *Burner, err error) {
	if t.IsEtcd {
		return nil, fmt.Errorf("Etcd is not supported yet")
	}

	// dialing
	c, err := grpc.Dial(t.Target, grpc.WithInsecure())
	if err != nil {
		return nil, err
	}

	// healthy check
	grpcCheck := grpc_health_v1.NewHealthClient(c)
	checkReq := &grpc_health_v1.HealthCheckRequest{
		Service: "", // leave empty to check all services
	}

	_, checkErr := grpcCheck.Check(ctx, checkReq)
	if checkErr != nil {
		return nil, fmt.Errorf("Not Healthy")
	}

	return &Burner{
		conn:    c,
		workers: DefaultGRPCWorkers,
	}, nil
}

func (b *Burner) Burn(rate uint64, du time.Duration) <-chan *Result {

	var workers sync.WaitGroup
	results := make(chan *Result)
	ticks := make(chan time.Time)
	for i := uint64(0); i < b.workers; i++ {
		workers.Add(1)
		go b.burn(&workers, ticks, results)
	}

	go func() {
		defer close(results)
		defer workers.Wait()
		defer close(ticks)
		interval := 1e9 / rate
		hits := rate * uint64(du.Seconds())
		began, done := time.Now(), uint64(0)
		for {
			now, next := time.Now(), began.Add(time.Duration(done*interval))
			time.Sleep(next.Sub(now))
			select {
			case ticks <- max(next, now):
				if done++; done == hits {
					return
				}
			case <-b.stopch:
				return
			default: // all workers are blocked. start one more and try again
				workers.Add(1)
				go b.burn(tr, &workers, ticks, results)
			}
		}
	}()

	return results
}

func (b *Burner) Close() {
	b.conn.Close()
}

func (b *Burner) Stop() {
	select {
	case <-b.stopch:
		return
	default:
		close(b.stopch)
	}
}

func (b *Burner) burn(workers *sync.WaitGroup, ticks <-chan time.Time, results chan<- *Result) {
	defer workers.Done()
	for tm := range ticks {
		results <- b.hit(tm)
	}
}
