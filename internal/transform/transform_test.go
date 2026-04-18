package transform_test

import (
	"testing"

	"github.com/hashicorp/consul/api"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"consul-drift-check/internal/transform"
)

func makePairs(kvs ...string) []*api.KVPair {
	pairs := make([]*api.KVPair, 0, len(kvs)/2)
	for i := 0; i+1 < len(kvs); i += 2 {
		pairs = append(pairs, &api.KVPair{Key: kvs[i], Value: []byte(kvs[i+1])})
	}
	return pairs
}

func pairKeys(pairs []*api.KVPair) []string {
	keys := make([]string, len(pairs))
	for i, p := range pairs {
		keys[i] = p.Key
	}
	return keys
}

func TestPipeline_Empty_ReturnsUnchanged(t *testing.T) {
	input := makePairs("a/b", "1", "a/c", "2")
	out := transform.New().Run(input)
	assert.Equal(t, pairKeys(input), pairKeys(out))
}

func TestPipeline_StripPrefix(t *testing.T) {
	input := makePairs("prod/db/host", "x", "prod/db/port", "y")
	p := transform.New().Add("strip", transform.StripPrefix("prod/db"))
	out := p.Run(input)
	assert.Equal(t, []string{"host", "port"}, pairKeys(out))
}

func TestPipeline_LowerKeys(t *testing.T) {
	input := makePairs("Config/DB", "v", "Config/HOST", "w")
	p := transform.New().Add("lower", transform.LowerKeys())
	out := p.Run(input)
	assert.Equal(t, []string{"config/db", "config/host"}, pairKeys(out))
}

func TestPipeline_Stages_ReturnsNames(t *testing.T) {
	p := transform.New().
		Add("strip", transform.StripPrefix("x")).
		Add("lower", transform.LowerKeys())
	assert.Equal(t, []string{"strip", "lower"}, p.Stages())
}

func TestPipeline_NilPairs_Skipped(t *testing.T) {
	input := []*api.KVPair{nil, {Key: "a/b", Value: []byte("1")}}
	p := transform.New().Add("strip", transform.StripPrefix("a"))
	out := p.Run(input)
	require.Len(t, out, 1)
	assert.Equal(t, "b", out[0].Key)
}
