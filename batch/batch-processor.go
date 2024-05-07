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
	default_interval = 100
)

type CallbackFunc func(ctx context.Context, task any) (any, error)

type InterruptFunc func() bool

type processor struct {
	ctx           context.Context
	interval      int
	taskList      []any
	cocurrentCnt  int
	processNum    int
	callback      CallbackFunc
	interruptFunc InterruptFunc
}

// init processor
func NewBatchProcessor(ctx context.Context, taskList any, callback CallbackFunc) *processor {
	return &processor{
		ctx:          ctx,
		interval:     default_interval,
		taskList:     transformInterfaceToAnySlice(taskList),
		cocurrentCnt: default_cocurrent_count,
		processNum:   default_process_num,
		callback:     callback,
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
	if len(p.taskList) == 0 {
		return nil, errors.New("Empty input task list")
	}

	goPool := pool.New().WithErrors().WithContext(p.ctx).WithMaxGoroutines(p.cocurrentCnt)

	SplitTaskList := SplitByLength(p.taskList, p.processNum)

	for index, SliceTaskList := range SplitTaskList {
		for _, task := range SliceTaskList {
			if p.isInterrupted() {
				return nil, errors.New("process interrupted, process index=" + string(index)
			}
			tmpTask := task
			goPool.Go(func(ctx context.Context) error {
				return p.callback(ctx, tmpTask)
			})
			time.Sleep(time.Millisecond * time.Duration(p.interval))
		}
	}

	return nil, goPool.Wait()
}
