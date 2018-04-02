package main

import (
	"flag"
	"time"

	trunks "github.com/straightdave/trunks/lib"
)

func burnCmd() command {
	fs := flag.NewFlagSet("trunks burn", flag.ExitOnError)
	opts := &burnOpts{}

	fs.StringVar(&opts.target, "target", "stdin", "Target host:port/default stdin")
	fs.BoolVar(&opts.isTargetEtcd, "etcd", false, "If the target is an etcd service")
	fs.StringVar(&opts.serviceName, "service", "", "Service name")
	fs.StringVar(&opts.methodName, "method", "", "Method name")
	fs.StringVar(&opts.pbgof, "pbgo", "", "pb.go file (compiled pb)")
	fs.DurationVar(&opts.duration, "duration", 0, "Duration of the test [0 = forever]")
	fs.DurationVar(&opts.timeout, "timeout", trunks.DefaultTimeout, "Requests timeout")
	fs.Uint64Var(&opts.rate, "rate", 50, "Requests per second")

	return command{fs, func(args []string) error {
		fs.Parse(args)
		return burn(opts)
	}}
}

type burnOpts struct {
	target       string // service host
	isTargetEtcd bool   // if set, will treat target as etcd instance to query service hosts

	serviceName string
	methodName  string
	pbgof       string // pb.go file (compiled from gRPC pb definition)

	duration time.Duration
	timeout  time.Duration
	rate     uint64
}

func burn(opts *burnOpts) (err error) {
	if opts.rate == 0 {
		return errZeroRate
	}

	// step 1: create stub client

	// step 2: invoke stub client

	return

}
