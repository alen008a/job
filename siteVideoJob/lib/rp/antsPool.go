package rp

import (
	"fmt"
	"github.com/panjf2000/ants/v2"
	"log"
	"runtime"
	"siteVideoJob/internal/glog"
	"strings"
)

var (
	pool     *ants.Pool
	maxStack = 20
)

func InitGlobal(size int) {
	var err error
	pool, err = ants.NewPool(size)
	if err != nil {
		log.Fatalf("ants协程池创建失败，%v", err)
	}
}

func HandlePanic(msg interface{}) {
	sb := strings.Builder{}
	sb.WriteString(fmt.Sprintf("%v", msg))
	sb.WriteString("\n")
	i := 2
	for {
		pc, file, line, ok := runtime.Caller(i)
		if !ok || i > maxStack {
			break
		}
		sb.WriteString(fmt.Sprintf("[stack: %d] %s:%d %s\n", i-1, file, line, runtime.FuncForPC(pc).Name()))
		i++
	}
	glog.Error(sb.String())
}

func Go(f func()) {
	err := pool.Submit(f)
	if err != nil {
		glog.Errorf("ants协程池添加任务失败，%v", err)
	}
}

func ReleaseGlobal() {
	pool.Release()
}

func InitPool(size int) *ants.Pool {
	p, err := ants.NewPool(size, ants.WithPanicHandler(HandlePanic))
	if err != nil {
		log.Fatalf("ants协程池创建失败，%v", err)
	}
	return p
}

func GoWithPool(p *ants.Pool, f func()) {
	err := p.Submit(f)
	if err != nil {
		glog.Errorf("ants协程池添加任务失败，%v", err)
	}
}
