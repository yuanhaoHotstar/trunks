# Trunks [![Build Status](https://travis-ci.org/straightdave/trunks.svg?branch=master)](https://travis-ci.org/straightdave/trunks)

Trunks, like every son, is derived from the father Vegeta with some enhanced skills:
1. dump HTTP reponses
2. gRPC support

![Trunks](http://images2.wikia.nocookie.net/__cb20100725123520/dragonballfanon/images/5/52/Future_Trunks_SSJ2.jpg)

## Usage manual

for original usage of Vegeta, please refer to [vegeta' readme](https://github.com/tsenart/vegeta/blob/master/README.md)

## More functionalities

### dump http attck response to file
```console
(add one more option '-respf' to 'attack')
-respf string
      Dump responses to file
```

### gRPC perf test (now as a lib)

Using the _burn_ attack to burn the gRPC services. Capable with any gRPC service.

Given multiple hosts ('IP:port' of service instances), Trunks will use simple round-robin as client-side load balance mechanism.

>Burning duration is shorter. No need to do the real __service discovery__ or __watch / real-time live connection's update__. Say, if testing service instances registered to Etcd, we don't need to watch the status of all instances and adjust the connection pool. During the short time of perf testing, we assume the instance hosts are _not-changing_.

Example:

>This example is using "google.golang.org/grpc/examples/route_guide" as the target server.

```go
package main

import (
    "fmt"
    "time"

    trunks "github.com/straightdave/trunks/lib"
    // for convenience, change the client_stub.pb.go into package main
)

func main() {
    tgt := &trunks.Gtarget{
        MethodName: "/routeguide.RouteGuide/GetFeature",
        Request:    &Point{Latitude: 10000, Longitude: 10000},
        Response:   &Feature{},
    }

    burner, err := trunks.NewBurner(
        []string{":10000"},   // a pool with round-robin picker
        trunks.NumWorker(20), // if not specified, 10 is default
    )
    if err != nil {
        fmt.Println(err)
        return
    }
    defer burner.Close()

    var metrics trunks.Metrics
    startT := time.Now()
    for res := range burner.Burn(tgt, uint64(5), 10*time.Second) {
        metrics.Add(res)
    }
    dur := time.Since(startT)
    metrics.Close()

    fmt.Printf("dur: %v\n", dur.Seconds())
    fmt.Printf("earliest: %v\n", metrics.Earliest.Sub(startT).Nanoseconds())
    fmt.Printf("latest: %v\n", metrics.Latest.Sub(startT).Nanoseconds())
    fmt.Printf("end: %v\n", metrics.End.Sub(startT).Nanoseconds())
    fmt.Printf("reqs: %d\n", metrics.Requests)
    fmt.Printf("50th: %s\n", metrics.Latencies.P50)
    fmt.Printf("95th: %s\n", metrics.Latencies.P95)
    fmt.Printf("99th: %s\n", metrics.Latencies.P99)
    fmt.Printf("mean: %s\n", metrics.Latencies.Mean)
    fmt.Printf("max: %s\n", metrics.Latencies.Max)
}
```

For this code snippet, it would result in:
```console
dur: 9.802099215
earliest: 68670
latest: 9800068577
end: 9802058490
reqs: 50
50th: 5.974748ms
95th: 6.084433ms
99th: 6.10946ms
mean: 5.143272ms
max: 6.19225ms
```

_NOTE_
* Currently Trunks is ignoring response unmarshalling. We don't need to provide response data structures to call gRPC endpoints.

## Arion as gRPC
Arion makes it easy to use Trunks. Please check https://github.com/straightdave/arion


