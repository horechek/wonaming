package consul

import (
	"fmt"
	consul "github.com/hashicorp/consul/api"
	"google.golang.org/grpc/resolver"
)

// Builder is the implementaion of grpc.naming.Resolver
type Builder struct {
	addr string //service name
}

// NewBuilder return Builder with service name
func NewBuilder(addr string) resolver.Builder {
	return &Builder{addr: addr}
}

func (cr *Builder) Build(target resolver.Target, cc resolver.ClientConn, opts resolver.BuildOption) (resolver.Resolver, error) {
	var err error

	// generate consul client, return if error
	conf := &consul.Config{
		Scheme:  "http",
		Address: cr.addr,
	}
	client, err := consul.NewClient(conf)
	if err != nil {
		return nil, fmt.Errorf("wonaming: creat consul error: %v", err)
	}

	r := &Resolver{
		cc:     cc,
		consul: client,
	}

	go r.watch(target)

	return r, nil
}

func (cr *Builder) Scheme() string {
	return ""
}

type Resolver struct {
	cc     resolver.ClientConn
	consul *consul.Client
}

// Resolve to resolve the service from consul, target is the dial address of consul
func (cr *Resolver) ResolveNow(option resolver.ResolveNowOption) {
	fmt.Println(option)
}

func (cr *Resolver) Close() {
	fmt.Println("close")
}

func (r *Resolver) watch(target resolver.Target) {
	fmt.Println(target)
}
