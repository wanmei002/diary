## 我们从Finish入手
func Finish(... func() error) error {
 .......
}

这个里面运行了一个 MapReduceVoid(func ...); 这个方法有四个参数，都是 func, 我们分别对这些方法映射下:
```go
// 把函数放入管道里
a := func(source chan<- interface{}) {
        for _, fn := range fns {// 这个 fns 就是传入的要执行的方法
            source <- fn
        }
    }
```

```go
// 这里执行管道里的方法
b := func(item interface{}, writer Writer, cancel func(error)) {
        fn := item.(func() error)
        if err := fn(); err != nil {
            cancel(err)
        }
    }
```

```go
c := func(pipe <-chan interface{}, cancel func(error)) {
        drain(pipe)
    }
```

```go
d := func(opts *mapReduceOptions) {
        if workers < minWorkers {// workers 是传入方法的数量
            opts.workers = minWorkers
        } else {
            opts.workers = workers
        }
    }
```

把上面的四个方法传入了 MapReduceVoid 这个方法里；让我们看看这个方法怎么运行的
```go
func MapReduceVoid(generate GenerateFunc, mapper MapperFunc, reducer VoidReducerFunc, opts ...Option) error {
	_, err := MapReduce(generate, mapper, func(input <-chan interface{}, writer Writer, cancel func(error)) {
		reducer(input, cancel)
		// We need to write a placeholder to let MapReduce to continue on reducer done,
		// otherwise, all goroutines are waiting. The placeholder will be discarded by MapReduce.
		writer.Write(lang.Placeholder)
	}, opts...)
	return err
}
```
这个方法引用了  MapReduce 这个方法
```go
e := func(input <-chan interface{}, writer Writer, cancel func(error)) {
        c(input, cancel)
        // We need to write a placeholder to let MapReduce to continue on reducer done,
        // otherwise, all goroutines are waiting. The placeholder will be discarded by MapReduce.
        writer.Write(lang.Placeholder)
    }

MapReduceVoid(a, b, c, d func() error) error {
    MapReduce(a, b, e, d)
}
```

我们再看看 MapReduce方法是怎么运行的
```go
func MapReduce(a, b, e, d) (interface{}, error) {
	source := buildSource(a)
	return MapReduceWithSource(source, b, e, d)
}
```

上面执行的 buildSource 方法如下:
```go
func buildSource(a) chan interface{} {
	source := make(chan interface{})
	threading.GoSafe(func() {
		defer close(source)
		a(source)
	})

	return source
}
```
让我们看看 `MapReduceWithSource` 这个方法的逻辑
```go
func MapReduceWithSource(execFuncChan, b, e, d) (interface{}, error) {
    // 设置 workers 属性为传进来方法的个数
	options := buildOptions(d)
	output := make(chan interface{})
	// 创建 chan，长度为要执行的函数的个数
	collector := make(chan interface{}, options.workers)
    // &DoneChan{done: make(chan lang.PlaceholderType)}
	done := syncx.NewDoneChan()
    // guardedWriter{ channel: output, done:DoneChan.done}
	writer := newGuardedWriter(output, done.Done())
	var closeOnce sync.Once
	var retErr errorx.AtomicError
	finish := func() {// 关闭开启的chan
		closeOnce.Do(func() {
			done.Close()
			close(output)
		})
	}
	// 执行一次 传进来的方法，这个方法会把管道里的方法都取出来但是不执行
	cancel := once(func(err error) {
		if err != nil {
			retErr.Set(err)
		} else {
			retErr.Set(ErrCancelWithNil)
		}
        // 开始从管道里取数据
		drain(execFuncChan)
        // 开始关闭管道
		finish()
	})

	go func() {
		defer func() {
			drain(collector)

			if r := recover(); r != nil {
				cancel(fmt.Errorf("%v", r))
			} else {
				finish()
			}
		}()
        // 所有要执行的函数执行完了会 关闭 collector,此时关闭剩下的管道
		e(collector, writer, cancel)
	}()
    //e := func(input <-chan interface{}, writer Writer, cancel func(error)) {
    //        c(input, cancel)
    //        writer.Write(lang.Placeholder)
    //    }
    // c := func(pipe <-chan interface{}, cancel func(error)) {
    //     drain(pipe)
    // }
    
    MapReduceVoid(a, b, c, d func() error) error {
        MapReduce(a, b, e, d)
    }
    // 这个函数的作用就是 如果 done 管道里有数据就直接返回，如果没有就开启一个协程执行传入要执行的方法
	go executeMappers(func(item interface{}, w Writer) {
		b(item, w, cancel)
	}, execFuncChan, collector, done.Done(), options.workers)

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

```go
func executeMappers(mapper MapFunc, input <-chan interface{}, collector chan<- interface{},
	done <-chan lang.PlaceholderType, workers int) {
	var wg sync.WaitGroup
	defer func() {
		wg.Wait()
		close(collector)
	}()

	pool := make(chan lang.PlaceholderType, workers)
	writer := newGuardedWriter(collector, done)
	for {
		select {
        // done 管道里有值，就直接返回
		case <-done:
			return
		case pool <- lang.Placeholder:
            // 管道已经关闭了，取不出来值 直接返回
			item, ok := <-input
			if !ok {
				<-pool
				return
			}
            // 开启协程执行方法
			wg.Add(1)
			// better to safely run caller defined method
			threading.GoSafe(func() {
				defer func() {
					wg.Done()
					<-pool
				}()
                // 里面执行传入的b方法
				mapper(item, writer)
                // func(item interface{}, w Writer) {
                // 		b(item, w, cancel)
                // 	}
                //b := func(item interface{}, writer Writer, cancel func(error)) {
            //        fn := item.(func() error)
            //        if err := fn(); err != nil {
            //            cancel(err)
            //        }
            //    }
			})
		}
	}
}
```

### 主要逻辑
#### 第一步 执行 buildSource 运行传入的 a 方法
在这里创建了一个无缓存的管道，每次往里面存入一个要执行的方法


