# Trunks [![Build Status](https://travis-ci.org/straightdave/trunks.svg?branch=master)](https://travis-ci.org/straightdave/trunks)

Trunks, like every son, is derived from the father Vegeta with some enhanced skills:
1. dump HTTP reponses
2. gRPC support

>My purpose is building a universal gRPC tool. But it's Golang and it's gRPC, we need the client stub and we need to do more than just one command. So far we can build custom program with Trunks as a lib. Sooner I would add some glue code to generate such client stub code and run it automatically.

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

### gRPC perf test (as a lib)

use 'burn' to burn the gRPC services. Capable to any gRPC service.

Sample file:
```go
package main

import (
  "fmt"
  "log"
  "time"

  trunks "github.com/straightdave/trunks/lib"
)

func main() {
  log.Println("hello")
  tgt := &trunks.GTargeter{
    Target:     ":8087",
    IsEtcd:     false,
    MethodName: "/myapppb.MyApp/Hello",
    Request:    &HelloRequest{Name: "dave"},   // from xxx.pb.go
    Response:   &HelloResponse{},
  }

  b, err := tgt.GenBurner()
  if err != nil {
    fmt.Println(err)
    return
  }
  defer b.Conn.Close()

  var metrics trunks.Metrics

  startT := time.Now()
  for res := range b.Burn(tgt, uint64(5), 10*time.Second) {
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

For details please read the code (currently smell but promised to be better)

## TODO

### connection pool
to support client side load balance testing against a cluster of service instances with service discovery.

### glue code for gRPC perf
to make it a portable one-statement command. Quick and easy.
