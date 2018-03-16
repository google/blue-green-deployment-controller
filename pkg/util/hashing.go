package util

import (
	"hash/fnv"

	"k8s.io/api/core/v1"
	hashutil "k8s.io/kubernetes/pkg/util/hash"
)

// ComputeHash returns a hash value calculated from pod template to avoid hash collision.
func ComputeHash(template *v1.PodSpec) uint32 {
	podTemplateSpecHasher := fnv.New32a()
	hashutil.DeepHashObject(podTemplateSpecHasher, *template)
	return podTemplateSpecHasher.Sum32()
}
