// Package transform applies a pipeline of named transformations to KV pairs.
package transform

import (
	"strings"

	"github.com/hashicorp/consul/api"
)

// Stage is a single transformation step.
type Stage struct {
	Name string
	Fn   func([]*api.KVPair) []*api.KVPair
}

// Pipeline holds an ordered list of stages.
type Pipeline struct {
	stages []Stage
}

// New returns an empty Pipeline.
func New() *Pipeline {
	return &Pipeline{}
}

// Add appends a stage to the pipeline.
func (p *Pipeline) Add(name string, fn func([]*api.KVPair) []*api.KVPair) *Pipeline {
	p.stages = append(p.stages, Stage{Name: name, Fn: fn})
	return p
}

// Run executes all stages in order and returns the final slice.
func (p *Pipeline) Run(pairs []*api.KVPair) []*api.KVPair {
	out := pairs
	for _, s := range p.stages {
		out = s.Fn(out)
	}
	return out
}

// Stages returns the names of registered stages.
func (p *Pipeline) Stages() []string {
	names := make([]string, len(p.stages))
	for i, s := range p.stages {
		names[i] = s.Name
	}
	return names
}

// StripPrefix returns a stage that removes a leading key prefix.
func StripPrefix(prefix string) func([]*api.KVPair) []*api.KVPair {
	prefix = strings.TrimRight(prefix, "/") + "/"
	return func(pairs []*api.KVPair) []*api.KVPair {
		out := make([]*api.KVPair, 0, len(pairs))
		for _, p := range pairs {
			if p == nil {
				continue
			}
			cp := *p
			cp.Key = strings.TrimPrefix(cp.Key, prefix)
			out = append(out, &cp)
		}
		return out
	}
}

// LowerKeys returns a stage that lower-cases all keys.
func LowerKeys() func([]*api.KVPair) []*api.KVPair {
	return func(pairs []*api.KVPair) []*api.KVPair {
		out := make([]*api.KVPair, 0, len(pairs))
		for _, p := range pairs {
			if p == nil {
				continue
			}
			cp := *p
			cp.Key = strings.ToLower(cp.Key)
			out = append(out, &cp)
		}
		return out
	}
}
