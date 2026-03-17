package metrics

import "hash/fnv"

func shardFor(name string, shards int) int {
	h := fnv.New32a()
	_, _ = h.Write([]byte(name))
	return int(h.Sum32() % uint32(shards))
}

