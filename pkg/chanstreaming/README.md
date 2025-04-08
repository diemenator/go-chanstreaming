# chanstreaming

## A set of functions to command channel-driven workflows in golang

Golang's channel is a very generic construct that allows to build a variety of parallel processing pipelines with ease, while keeping things modular, type-safe and expressive.

Here is a sequence of intuitions that underline the designs found in this repo
- a size-1 readonly-channel can represent an awaitable Future
- a N-sized channel of 1-sized channels can represent a group of workers
- a strongly-typed readonly channel can be used to implement a typed async iterable/enumerable (like in C# and Python) programming models
- a collection of channels and goroutines can be used to implement a parallel processing pipeline sdk (like in Akka Streams and PLINQ)

# The chanstreaming lib

The `<-chan streaming lib` is a small collection of utilities built around golang's channels that implements **dataflow** paradigm, intended as convenient tool for building modular workflows and data streaming applications.

The functions featured here accept read-only typed channels as inputs and return typed channels as outputs, allowing for easy composition and chaining of operations.

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

## There's go-streams, why another one?
- The `go-streams` lib implements idiomatical, Java-inspired Streams API and/or similar data flow building frameworks seen in other languages. This is done traditionally in objective way, exposing the Fluent-style interface to give you a concise workflow builder.
- The `chanstreaming` lib addresses roughly same class of data/control streaming scenarios, but advocates for the re-use the `<-chan T` primitive as the main object of the module's API surface. Decouple, extend, test & rearrange the workflows in type-safe way.
- For production use, the real difference would be the style of the execution. Depending on the task and background one could choose to describe complex precisely-timed workflows with go-streams first, or inline the timed-concurrency-critical piece in its own code space.
- The two are very much compatible as they can be used together in same project.
- There are no generic methods in golang, so the `chanstreaming` lib does not try to implement them by hacking around `reflect` and `any`. We simply expose higher order functions in the API instead.

### What's missing?
- [ ] Examples:
  - [ ] Basic system & IO (FromCSV(filename), FromShell(command, args...) (to produce a stream of stdout+stderr+eof+exitCode messages)
  - [ ] Integration demos (shell, kafka, sql, rpc)
  - [ ] More tests
