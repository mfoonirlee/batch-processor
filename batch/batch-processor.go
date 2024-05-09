package processor

import (
	"context"
	"errors"
	"time"

	"github.com/sourcegraph/conc/pool"
)

const (
	// default process request number
	default_process_num int = 10
	// default process goruntine count
	default_cocurrent_count = 10
	// default process interval with Millisecond
	default_interval = 10
)

type CallerFunc func(ctx context.Context, param any) (any, error)

type InterruptFunc func() bool

type processor struct {
	ctx           context.Context
	interval      int
	paramsList    []any
	cocurrentCnt  int
	processNum    int
	callerFunc    CallerFunc
	interruptFunc InterruptFunc
}

// init processor
func NewBatchProcessor(ctx context.Context, paramsList any, callerFunc CallerFunc) *processor {
	return &processor{
		ctx:          ctx,
		interval:     default_interval,
		paramsList:   transformInterfaceToAnySlice(paramsList),
		cocurrentCnt: default_cocurrent_count,
		processNum:   default_process_num,
		callerFunc:   callerFunc,
	}
}

func (p *processor) SetInterval(interval int) *processor {
	p.interval = interval
	return p
}

func (p *processor) SetCocurrentCnt(cocurrentCnt int) *processor {
	p.cocurrentCnt = cocurrentCnt
	return p
}

func (p *processor) SetProcessNum(processNum int) *processor {
	p.processNum = processNum
	return p
}

func (p *processor) SetInterruptFunc(interruptFunc InterruptFunc) *processor {
	p.interruptFunc = interruptFunc
	return p
}

// to control pool excution
func (p *processor) isInterrupted() bool {
	if p.interruptFunc != nil {
		return p.interruptFunc()
	}
	return false
}

func (p *processor) Run() (any, error) {
	if len(p.paramsList) == 0 {
		return nil, errors.New("Empty input params list")
	}

	goPool := pool.NewWithResults[any]().WithErrors().WithContext(p.ctx).WithMaxGoroutines(p.cocurrentCnt)

	splitParamsList := SplitByLength(p.paramsList, p.processNum)

	for index, SliceParamsList := range splitParamsList {
		for _, param := range SliceParamsList {
			if p.isInterrupted() {
				return nil, errors.New("process interrupted, process index=" + string(rune(index)))
			}
			tmpParam := param
			goPool.Go(func(ctx context.Context) (any, error) {
				return p.callerFunc(ctx, tmpParam)
			})
			time.Sleep(time.Millisecond * time.Duration(p.interval))
		}
	}

	return goPool.Wait()
}
