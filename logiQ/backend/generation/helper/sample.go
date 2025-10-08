package helper

import "math/rand/v2"

// SampleWeighted 辅助函数 从一个权重分布中随机抽样一个值。
func SampleWeighted[T comparable](dist map[T]float64, rng *rand.Rand) T {
	var zero T
	if len(dist) == 0 {
		return zero
	}

	var total float64
	for _, w := range dist {
		if w > 0 {
			total += w
		}
	}
	if total <= 0 {
		return zero
	}

	var roll float64
	if rng != nil {
		roll = rng.Float64() * total
	} else {
		roll = rand.Float64() * total
	}

	var (
		cumulative  float64
		fallback    T
		hasFallback bool
	)
	for val, w := range dist {
		if w <= 0 {
			continue
		}
		if !hasFallback {
			fallback = val
			hasFallback = true
		}
		cumulative += w
		if roll < cumulative {
			return val
		}
	}

	if hasFallback {
		return fallback
	}
	return zero
}
