# hemera

[hemera](https://en.wikipedia.org/wiki/Hemera) is a zero-dependency [statsd](https://github.com/statsd/statsd) implementation written in Go. The project takes advantage of modular packages and interfaces to make connecting to any backend simple. By default, the hemera binary connects to a [Graphite](http://graphiteapp.org/) server over TCP. 

## Documentation

* Installing hemera
* Basic usage
* Supported metrics
* Creating/using new backend interfaces
* Future plans

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
hemera -p 8000 -t 5
```


## Supported metrics

hemera supports the four metrics that statsd supports. They are of the following form where bucket represents the name metric to update. 

` <bucket>:<value>|<metric-type>|@<sampling-rate>` 

### Counters

Counters represent metrics that can only increase. During each purge cycle the counter is reset to 0. For example, the following packet will increment the 'button' counter by 2. 

```sh
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

