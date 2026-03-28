package cache

import (
	"siteLetterJob/internal/context"
	"sync"
	"time"

	"github.com/RussellLuo/timingwheel"
)

// 本地缓存方案
func init() {
	gTimingWheel.Start()
}

type (
	rotatePeriod  = time.Duration
	ValueFunction = func() (interface{}, error)
)

const (
	RotatePeriod10 = 10 * time.Second
	RotatePeriod30 = 30 * time.Second
	RotatePeriod60 = 60 * time.Second
)

var (
	gCache       = make(map[string]*cache)
	gTimingWheel = timingwheel.NewTimingWheel(time.Millisecond, 20)
	mux          = sync.RWMutex{}
)

type RotateScheduler struct{ Interval time.Duration }

func (s *RotateScheduler) Next(prev time.Time) time.Time {
	return prev.Add(s.Interval)
}

func getCache(c *context.Context, k string, s rotatePeriod, f ValueFunction) (*cache, error) {
	mux.RLock()
	ca, ok := gCache[k]
	if ok {
		mux.RUnlock()
		return ca, nil
	}

	mux.RUnlock()

	mux.Lock()

	data, err := f()
	if err != nil {
		mux.Unlock()
		c.Errorf("localCache |k=%s |err=%v", k, err)
		return nil, err
	}

	ca = &cache{
		k:      k,
		data:   data,
		rotate: s,
		f:      f,
		l:      new(sync.RWMutex),
	}

	gCache[k] = ca

	gCache[k].tw = gTimingWheel.ScheduleFunc(&RotateScheduler{Interval: s}, func() {
		//定期更新
		_, _ = gCache[k].load(c)
	})

	mux.Unlock()

	return ca, nil
}

type cache struct {
	k      string
	data   interface{}
	rotate rotatePeriod
	f      ValueFunction
	l      *sync.RWMutex
	tw     *timingwheel.Timer
}

func GetOrSet(c *context.Context, k string, s rotatePeriod, f ValueFunction) (interface{}, error) {
	ca, err := getCache(c, k, s, f)
	if err != nil {
		return nil, err
	}
	ca.l.RLock()
	data := ca.data
	if data == nil {
		ca.l.RUnlock()
		return ca.load(c)
	}
	ca.l.RUnlock()
	return data, nil
}

func (ca *cache) load(c *context.Context) (interface{}, error) {
	data, err := ca.f()
	if err != nil {
		c.Errorf("localCache |err=%v", err)
		return ca.data, err
	}
	ca.l.Lock()
	ca.data = data
	ca.l.Unlock()
	return data, nil
}

func Close() {
	mux.Lock()
	for _, c := range gCache {
		c.tw.Stop()
	}
	mux.Unlock()
	gTimingWheel.Stop()
}
