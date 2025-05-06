The package is a set of `<-chan T` utitilty functions, your `rochannel` std lib module.

## Background

Golang's channel is very versatile construct that allows to build variety of parallel processing pipelines with relative ease, while keeping things modular, type-safe and expressive.

Here is the sequence of intuitions that underline the designs found in this repo
- a size-1 readonly-channel can represent an awaitable
- a strongly-typed readonly channel can be used to interpreted a typed async iterable/enumerable, expressed with channels
- a collection of functions that use goroutines to create and serve new channels can be used to build parallel and highly composable sdks (like, Akka Streams or PLINQ)
also, if you go meta enough  
- a N-sized channel of 1-sized channels can represent a not initialised yet group of workers and more  

Here, we expliot the concept of typed readonly-channel as main subject of our API, thus making it possible to build very expressive dataflow sdk around it.

# The library

The `<-chan streaming lib` is a small collection of utilities built around golang's channels that implements **dataflow** paradigm, intended as convenient tool for building modular workflows and data streaming applications.

The functions featured here accept read-only typed channels as inputs and return readonly-typed channels as outputs, allowing for easy composition and chaining of operations.

Here you'll find:
- `Map(fn, maxWorkers)` & `MapUnordered(fn, maxWorkers)`
- `Partition(maxPartitions, partitioner)` & `Merge(sources)`
- `Batch(maxLength, maxInterval)` & `BatchWeighted(sizeFn, maxSize, maxCount, maxInterval)`
- `WithContext(context)` to make the channel close on cancellation
- `WhenDone(callback)` to invoke a callback on cancelling
- `Throttle(interval)`, `Jitter(interval)` for rate-limiting & extra randomness
- `Scan(fn, zero)`, `Fold(fn, zero)`, `WithSlidingWindowCount(count)`, `WithSlidingWindowTimed(interval)` for stateful processing
- `Mapped(fn)` and `Apply(fn)` for simple transformations and logging
- `FromSlice(slice)` and `ToSlice(source)` for converting channels and slices and more.

It relies heavily on Golang's generics for type safety, so this is not back-portable to golang pre-1.18.

Some of the functions implement patterns seen at https://go.dev/blog/pipelines, but this time taking advantage of generics to build a versatile toolkit.

Check tests to see intended usage.

Below you can read a fun summary of the core functions.

## The Instruments of Control

### **`Batch(maxLength, maxInterval)`: The Hoarding of Power**
A stream is held, gathered in force, then released at scale—or at time’s command, should patience wear thin.  
Used for **batch database inserts, rate-limited API calls, or bundling events before network transmission**.

### **`MapUnordered(fn, maxWorkers)`: Unleashing the Horde**
Fan out, consume, transform, and return—**ideal for compute-heavy workloads where order is irrelevant**.  
Useful for **CPU/GPU-bound tasks, parallel data crunching, and unordered batch processing** where throughput is the priority.

### **`Map(fn, maxWorkers)`: Sequential Precision**
Each element is processed in turn, **ensuring outputs match the order of inputs**.  
Suited for **web request processing, database writes, filesystem operations**, and any scenario where **causality and sequence integrity** must be maintained.

### **`Partition(maxPartitions, partitioner)`: The Dividing Blade**
The flow is split as dictated by a higher will—**the function decides, the system obeys**.  
Designed for **sharding workloads, distributing traffic, and parallelizing processing across consumer groups**.

### **`Merge(sources)`: The Great Convergence**
Many become one. When all sources are exhausted, the system closes itself—**no watchers, no counters, no waste**.  
Particularly useful for **aggregating multiple event sources, log streams, or external APIs into a unified pipeline**.

## There's go-streams and others, why another one?
- The `chanstreaming` lib addresses roughly same class of data/control streaming scenarios, but chooses to use the `<-chan T` (read-only channel) primitive as the central type of the module's API surface. Decouple, extend, test & rearrange the workflows in type-safe way by using pre-existing builtins.
- For production use, the real difference would be the style of the execution. For explicit control on both producer and consumer ends, one could consider to use `go-streams` first, or inline the timed-concurrency-critical pieces in their own coroutine or combine the approaches.
- There are no generic methods in golang, so the `chanstreaming` lib does not try to implement the wrapper that attempts to sidestep it with the use of `reflect` and `any`. We simply expose higher order functions in the API instead for easier composition.
