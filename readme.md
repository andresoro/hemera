# hemera

[hemera](https://en.wikipedia.org/wiki/Hemera) is a zero-dependency [statsd](https://github.com/statsd/statsd) implementation written in Go. The project takes advantage of modular packages and interfaces to make connecting to any backend simple. By default, the hemera binary connects to a [Graphite](http://graphiteapp.org/) server over TCP. 

## Documentation

* [Installing hemera](#installing-Hemera)
* [Basic usage](#basic-usage)
* [Supported metrics](#supported-metrics)
* [Creating new backend interfaces](#creating-new-backend-interfaces)
* [Using new backend interfaces](#using-new-backend-interfaces)
* [To Do](#to-do)

## Installing hemera

Install the binary by running the following:

``` 
go install github.com/andresoro/hemera
```

Install both the binary and the package by running:

```
go get github.com/andresoro/hemera
```

## Basic usage

First make sure the package has been installed then simply run the following to run hemera with default configuration.

```
hemera
```

This will start the server with default configuration with the following flags but each one can be changed individually:

* `-p 8484` port where hemera server listens for metrics over UDP
* `-g 2003` port to connect to graphite instance to purge metrics over TCP
* `-s localhost` server host
* `-t 10` interval in seconds for purging metrics to the server 


```sh
# listen for metrics on port 8000 and purge to backend every 5 seconds
hemera -p 8000 -t 5
```


## Supported metrics

hemera supports the four metrics that statsd supports. They are of the following form where bucket represents the name of the metric to update. 

` <bucket>:<value>|<metric-type>|@<sampling-rate>` 

### Counters

Counters represent metrics that can only be incremented. During each purge cycle the counter is reset to 0. 

```sh
# increment the 'button' bucket by 2
button:2|c
```

### Sets 

Sets hold a unique collection of values which represent events from the client. What is purge is the cardinality of the set at the purge interval.

```sh
# set named uniques will add the value '22' only if it does not already exist
uniques:22|s
```

### Gauges

A gauge can fluctuate positively or negatively and will take on an arbitrary value assigned to it. 

```sh
# set the 'gaugor' gauge to 100
gaugor:100|g

# set the 'gaugor' gauge to 1000
gaugor:1000|g
```

### Timers

Timers generate an array of statistics that are purged to the backend:

* Min/Max value
* Count
* Average
* Median
* 95th percentile
* Standard Deviation

Timers currently only support the `ms` metric tag.

```sh
# load-time took 225ms to complete this time
load-time:225|ms
```

## Creating new backend interfaces

Backends must simply satisfy the following interface:

``` go
type Backend interface {
    Purge(c *cache.Cache) error
}
```

It is up to the user to decide how to purge the actual values out. The `cache.Cache` struct exports all the values that are currently held since last purge cycle.

Warning: Do not clear the cache from the Purge() function as it is done by the server. Would cause issues when using multiple backends. 

For example, a backend implementation where we only purge counters to the standard output would look like this. 

```go
type ConsoleBackend struct{}

// implement backend interface
func (cb *ConsoleBackend) Purge(c *cache.Cache) error {

    // only handling counters
    for name, value := range c.Counters {
        fmt.Printf("counter name: %s value: %f \n", name, value)
    }

    return nil
}
```

Take a look at the [graphite implementation](https://github.com/andresoro/hemera/blob/master/pkg/backend/graphite.go) for a more robust example.

## Using new backend interfaces

The server takes care of the cache and metric collection. If you want to use a new backend interface simply define it and add it to a new server instance. The `server.New()` function can take in a variadic amount of backends.


```go
import github.com/andresoro/hemera/pkg/backend
import github.com/andresoro/hemera/pkg/server

// import a backend from the hemera backend package
graphite := &backend.Graphite{Addr: "localhost:2003"}

// and use a backend you implemented
console := &ConsoleBackend{}

// new server with given purge interval, host/port, and the backends that we would like to purge to. 
srv, err := server.New(purgeTime, host, port, graphite, console)

srv.Run()
```


## To Do

* Add support for incrementing/decrementing gauges with '+' or '-' signs in metric value. 
* Benchmark tests