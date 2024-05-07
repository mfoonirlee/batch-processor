# batch-processor

this is a batch processor for golang, it provide a simple way to process a batch of tasks.

## How to use

### Create

```golang
// initialize a batch processor
processor := NewBatchProcessor(ctx, taskList, callback)
```

### Config processor
```golang
// set batch request num 10
processor.SetProcessNum(10)

// set process interval by 100ms
processor.SetInterval(100)

// set process concurrent count 10
processor.SetConcurrent(10)
```

### Run processor
```golang
// run processor
processor.Run()
```

### Interrupt processor
```golang
// set interrupt func，when process once, it will call interruptFunc for check if need to interrupt
processor.SetInterruptFunc(interruptFunc)
```


