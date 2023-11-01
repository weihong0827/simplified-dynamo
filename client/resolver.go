package main

import (
	"google.golang.org/grpc/resolver"
)

const (
	exampleScheme = "example"
)

type exampleResolverBuilder struct {
	addrStore map[string][]string
}

func NewExampleResolverBuilder(addrStore map[string][]string) *exampleResolverBuilder {
	return &exampleResolverBuilder{addrStore: addrStore}
}

func (e *exampleResolverBuilder) Build(target resolver.Target, cc resolver.ClientConn, opts resolver.BuildOptions) (resolver.Resolver, error) {
	r := &exampleResolver{
		target:     target,
		cc:         cc,
		addrsStore: e.addrStore,
	}
	r.start()
	return r, nil
}

func (e *exampleResolverBuilder) Scheme() string { return exampleScheme }

type exampleResolver struct {
	target     resolver.Target
	cc         resolver.ClientConn
	addrsStore map[string][]string
}

func (r *exampleResolver) start() {
	addrStrs := r.addrsStore[r.target.Endpoint()]
	addrs := make([]resolver.Address, len(addrStrs))
	for i, s := range addrStrs {
		addrs[i] = resolver.Address{Addr: s}
	}
	r.cc.UpdateState(resolver.State{Addresses: addrs})
}

func (*exampleResolver) ResolveNow(o resolver.ResolveNowOptions) {}
func (*exampleResolver) Close()                                  {}

func init() {
	resolver.Register(&exampleResolverBuilder{})
}
