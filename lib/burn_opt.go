package trunks

import (
	"context"

	"google.golang.org/grpc/metadata"
)

type BurnOpt func(*Burner)

// Deprecating; use WithNumWorker(uint64)
func NumWorker(num uint64) BurnOpt {
	return func(b *Burner) { b.numWorker = num }
}

// Deprecating; use WithLooping(bool)
func WithLoop() BurnOpt {
	return func(b *Burner) { b.loop = true }
}

func WithNumWorker(num uint64) BurnOpt {
	return func(b *Burner) { b.numWorker = num }
}

func WithNumConnPerHost(num uint64) BurnOpt {
	return func(b *Burner) { b.numConnPerHost = num }
}

func WithLooping(yesno bool) BurnOpt {
	return func(b *Burner) { b.loop = yesno }
}

func WithDumpFile(fileName string) BurnOpt {
	return func(b *Burner) {
		b.dumpFile = fileName
		if fileName != "" {
			b.dump = true
		}
	}
}

func WithMetadata(md metadata.MD) BurnOpt {
	return func(b *Burner) {
		b.ctx = metadata.NewOutgoingContext(context.Background(), md)
	}
}

func WithMaxRecvSize(s int) BurnOpt {
	return func(b *Burner) {
		b.maxRecv = s
	}
}

func WithMaxSendSize(s int) BurnOpt {
	return func(b *Burner) {
		b.maxSend = s
	}
}
