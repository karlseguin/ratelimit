# Rate Limiter

Generic approach to tracking if an action goes over a certain treshold, measured in # / seconds. The library is generic, but the obvious example would be to limit the number of requests a user can make.

The library functions in two mode.

## Embedded Mode

In this mode, the tracking code integrates with your existing code. This mode is ideal if you already maintain long-lived objects. For example, if you were building a TCP chat server, you'd use the embedded mode and associate a `*Tracker` with your existing `*User`.

```go
    // your code embeds a *ratelimit.Tracker
    // create an instance via ratelimit.NewTracker()
    type User struct {
      tracker *ratelimit.Tracker
    }


    // whenever an action is taken that you want to limit:
    if user.tracker.Track(5) < 0 {
      // we've seen more than 5 requests per second
    }
```
If you want to track different events independently, create multiple trackers. A tracker is composed of an `int32` and an `int64` and uses the `sync/atomic` package for concurrency control

## Standalone Mode

In standalone mode, a rate limiter is backed by an LRU cache and tracks usage based on an arbitrary string key (such as an IP address).

Configure a new rate limiter instance:

```go
limiter := ratelimit.New(ratelimit.Configure().
              MaxAllowance(5).
              MaxItems(5000)
```

and use it to track requests:

```go
if limiter.Track("SOME_KEY") < 0 {

}
```

When configuring the limiter:

* `MaxAllowance(int)` - maximum number of requests per second allowed (default: 5)
* `MaxItems(int)` - maximum number of values to track (default: 5000)

It takes roughly 1.2MB to track 10 000 items. The limiter is thread-safe.

## Return Value

The return value is the number of allowed actions remaining. A value less than 0 means we're over the allocation. This can be used to take different actions based on how much over we are:

```go
  if allowed := tracker.Track(5); allowed < -5 {
    //way over, disconnect!
  } else if allowed < 0 {
    //a little over, warn
  }
```
