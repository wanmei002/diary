## mapreduce 并发请求，减少时间的等待
### 并发请求的没有关联 可以用 Finish()
```go
    func Finish(fns ...func() error) error {
        if len(fns) == 0 {
            return nil
        }
        // 传入的函数 都是  MapReduceVoid 处理了
        return MapReduceVoid(func(source chan<- interface{}) {
            for _, fn := range fns {
                source <- fn
            }
        }, func(item interface{}, writer Writer, cancel func(error)) {
            fn := item.(func() error)
            if err := fn(); err != nil {
                cancel(err)
            }
        }, func(pipe <-chan interface{}, cancel func(error)) {
            drain(pipe)
        }, WithWorkers(len(fns)))
    }
```

#### MapReduceVoid
```go
    func MapReduceVoid(generate GenerateFunc, mapper MapperFunc, reducer VoidReducerFunc, opts ...Option) error {
        _, err := MapReduce(generate, mapper, func(input <-chan interface{}, writer Writer, cancel func(error)) {
            reducer(input, cancel)
            drain(input)
            // We need to write a placeholder to let MapReduce to continue on reducer done,
            // otherwise, all goroutines are waiting. The placeholder will be discarded by MapReduce.
            writer.Write(lang.Placeholder)
        }, opts...)
        return err
    }
```

#### 紧接着到 MapReduce 了
```go
    func MapReduce(generate GenerateFunc, mapper MapperFunc, reducer ReducerFunc, opts ...Option) (interface{}, error) {
        source := buildSource(generate)
        return MapReduceWithSource(source, mapper, reducer, opts...)
    }
```

#### buildSource 内部是
```go
    func buildSource(generate GenerateFunc) chan interface{} {
        source := make(chan interface{})
        threading.GoSafe(func() {
            defer close(source)
            generate(source)
        })
    
        return source
    }
```
实际上 generate 方法是传入的 
```go
    func(source chan<- interface{}) {
        for _, fn := range fns {
            source <- fn
        }
    }
```

#### MapReduceWithSource
```go
// MapReduceWithSource maps all elements from source, and reduce the output elements with given reducer.
func MapReduceWithSource(source <-chan interface{}, mapper MapperFunc, reducer ReducerFunc,
	opts ...Option) (interface{}, error) {
	options := buildOptions(opts...)
	output := make(chan interface{})
                                        // options.workers 默认是 16
	collector := make(chan interface{}, options.workers)
            // 空结构体
	done := syncx.NewDoneChan()
	writer := newGuardedWriter(output, done.Done())
	var closeOnce sync.Once
	var retErr errorx.AtomicError
	finish := func() {
		closeOnce.Do(func() {
			done.Close()
			close(output)
		})
	}
	cancel := once(func(err error) {
		if err != nil {
			retErr.Set(err)
		} else {
			retErr.Set(ErrCancelWithNil)
		}

		drain(source)
		finish()
	})

	go func() {
		defer func() {
			if r := recover(); r != nil {
				cancel(fmt.Errorf("%v", r))
			} else {
				finish()
			}
		}()
		reducer(collector, writer, cancel)
		drain(collector)
	}()

	go executeMappers(func(item interface{}, w Writer) {
		mapper(item, w, cancel)
	}, source, collector, done.Done(), options.workers)

	value, ok := <-output
	if err := retErr.Load(); err != nil {
		return nil, err
	} else if ok {
		return value, nil
	} else {
		return nil, ErrReduceNoOutput
	}
}
```