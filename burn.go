package main

import (
	"flag"

	trunks "github.com/straightdave/trunks/lib"
)

func burnCmd() command {
	fs := flag.NewFlagSet("trunks attack", flag.ExitOnError)
	opts := &burnOpts{}

	fs.StringVar(&opts.target, "target", "stdin", "Target host:port/default stdin")
	fs.StringVar(&opts.serviceName, "service", "", "Service name")
	fs.StringVar(&opts.methodName, "method", "", "Method name")
	fs.StringVar(&opts.pbgof, "pbgo", "", "pb.go file (compiled pb)")
	fs.DurationVar(&opts.duration, "duration", 0, "Duration of the test [0 = forever]")
	fs.DurationVar(&opts.timeout, "timeout", trunks.DefaultTimeout, "Requests timeout")
	fs.Uint64Var(&opts.rate, "rate", 50, "Requests per second")
	fs.Uint64Var(&opts.workers, "workers", trunks.DefaultWorkers, "Initial number of workers")

	return command{fs, func(args []string) error {
		fs.Parse(args)
		return burn(opts)
	}}
}

type burnOpts struct {
	target      string // single for now
	serviceName string
	methodName  string
	pbgof       string
	duration    time.Duration
	timeout     time.Duration
	rate        uint64
	workers     uint64
}

func burn(opts *burnOpts) (err error) {

}
