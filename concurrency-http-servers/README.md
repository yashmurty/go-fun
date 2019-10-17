# Concurrency in Go HTTP Servers 

## Part 1/3 : 1_concurrent-maps 

```sh
cd 1_concurrent-maps
```

### HTTP Server

**Start the server**:
```sh
go run main.go
```


The server implements the following endpoints:

```sh
$ curl "localhost:8000/set?name=x&val=0"
ok
$ curl "localhost:8000/get?name=x"
x: 0
$ curl "localhost:8000/inc?name=x"
ok
$ curl "localhost:8000/get?name=x"
x: 1
```
The requests manipulate a shared `CounterStore` which is essentially a simple `map`.

**Simulate many concurrent connections with `ApacheBench`**:

```sh
$ ab -n 20000 -c 200 "127.0.0.1:8000/inc?name=i"

Benchmarking 127.0.0.1 (be patient)
Completed 2000 requests
Completed 4000 requests

Test aborted after 10 failures

apr_socket_connect(): Connection reset by peer (104)
Total of 4622 requests completed
```
We can see that the tests fail.

**Server logs from our Go HTTP Server**:

```sh
2019/10/16 17:57:11 inc &{GET /inc?name=i HTTP/1.0 1 0 map[Accept:[*/*] User-Agent:[ApacheBench/2.3]] {} <nil> 0 [] true 127.0.0.1:8000 map[] map[] <nil> map[] 127.0.0.1:45452 /inc?name=i <nil> <nil> <nil> 0xc000142640}

fatal error: concurrent map writes
2019/10/16 17:57:11 inc &{GET /inc?name=i HTTP/1.0 1 0 map[Accept:[*/*] User-Agent:[ApacheBench/2.3]] {} <nil> 0 [] true 127.0.0.1:8000 map[] map[] <nil> map[] 127.0.0.1:45446 /inc?name=i <nil> <nil> <nil> 0xc000418680}

goroutine 2319 [running]:
runtime.throw(0x6ecf75, 0x15)
	/usr/local/go/src/runtime/panic.go:617 +0x72 fp=0xc00050db50 sp=0xc00050db20 pc=0x42cf12
runtime.mapassign_faststr(0x68f5a0, 0xc000094ba0, 0xc000270eee, 0x1, 0xc00009ce58)
	/usr/local/go/src/runtime/map_faststr.go:211 +0x42a fp=0xc00050dbb8 sp=0xc00050db50 pc=0x413cda
main.CounterStore.inc(0xc000094ba0, 0x74a300, 0xc0004948c0, 0xc0002e8100)

```

The error that caused the failure is:    
`fatal error: concurrent map writes`

The request handlers can run concurrently but they all manipulate a shared `CounterStore`.    
For example, the `inc` handler is being called concurrently for multiple requests and attempts to mutate the `map` in the `CounterStore`.    
This leads to a race condition since in Go, map operations are not atomic.
https://golang.org/doc/faq#atomic_maps

## PART 2/3 : 2_mutex-maps

To fix the race condition, we will add a `mutex`.

We add the following 2 changes:

* We embed a `sync.Mutex` in `CounterStore`, and each handler starts by locking the mutex (and deferring an unlock).
* We change the receiver `inc` is defined on to a pointer `*CounterStore`. In fact, the previous version of the code was wrong in this respect - methods that modify data should always be defined with pointer receivers. We got lucky that the data was shared at all with value receivers because maps are reference types. Pointer receivers are particularly critical when mutexes are involved.

## PART 3/3 : 3_channel-commands

Instead of mutexes, we could use channels to synchronize access to shared data.
We start by defining a "counter manager" which is a background goroutine with access to a closure that stores the actual data.

Instead of accessing the map of counters directly, handlers will send `Commands` on a channel and will receive replies on a reply channel they provide.

The shared object for the handlers will now be a Server.
```sh
type Server struct {
  cmds chan<- Command
}
```

Each handler deals with the manager synchronously; the `Command` send is blocking, and so is the read from the reply channel. But note - not a mutex in sight! Mutual exclusion is accomplished by a single goroutine having access to the actual data.

While it certainly looks like an interesting technique, for our particular use case this approach seems like an overkill. In fact, overuse of channels is one of the common gotchas for Go beginners. 

These tests are taken from: https://eli.thegreenplace.net/2019/on-concurrency-in-go-http-servers/