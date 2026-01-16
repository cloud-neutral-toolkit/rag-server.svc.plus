package bench

import (
	"math"
	"sort"
)

type Metrics struct {
	HitAtK    float64
	RecallAtK float64
	MRR       float64
	NDCGAtK   float64
}

type LatencyStats struct{ P50, P95 int }

func CalcMetrics(cases []PerCase, k int) Metrics {
	var hitN, total int
	var rrSum float64
	var dcgSum, idcgSum float64
	var recallSum float64

	for _, c := range cases {
		if c.Error != nil {
			continue
		}
		total++

		// expected set
		exp := c.Expected
		if len(exp) == 0 {
			continue
		}

		// Hit@K
		hit := caseHitAtK(c)
		if hit {
			hitN++
		}

		// Recall@K: (# 前K内命中标准片段) / (标准片段总数)
		matched := 0
		for _, h := range c.Hits {
			if _, ok := exp[h.ID]; ok {
				matched++
			}
		}
		recall := float64(matched) / float64(len(exp))
		if recall > 1 {
			recall = 1
		}
		recallSum += recall

		// Reciprocal Rank
		rank := firstHitRank(c)
		if rank > 0 {
			rrSum += 1.0 / float64(rank)
		}

		// nDCG@K （二值相关性）
		dcg := 0.0
		for i, h := range c.Hits {
			if _, ok := exp[h.ID]; ok {
				rel := 1.0
				dcg += (math.Pow(2, rel) - 1) / math.Log2(float64(i+2)) // i从0
			}
		}
		// IDCG：把rel=1的若干个放在最前面
		g := minInt(len(exp), k)
		idcg := 0.0
		for i := 0; i < g; i++ {
			idcg += (math.Pow(2, 1.0) - 1) / math.Log2(float64(i+2))
		}
		if idcg == 0 {
			idcg = 1
		} // 防零
		dcgSum += dcg
		idcgSum += idcg
	}

	if total == 0 {
		return Metrics{}
	}
	return Metrics{
		HitAtK:    float64(hitN) / float64(total),
		RecallAtK: recallSum / float64(total),
		MRR:       rrSum / float64(total),
		NDCGAtK:   (dcgSum / idcgSum),
	}
}

func calcLatency(cases []PerCase) LatencyStats {
	var xs []int
	for _, c := range cases {
		if c.Error == nil {
			xs = append(xs, int(c.Latency))
		}
	}
	if len(xs) == 0 {
		return LatencyStats{}
	}
	sort.Ints(xs)
	p50 := xs[len(xs)*50/100]
	p95 := xs[len(xs)*95/100]
	return LatencyStats{P50: p50, P95: p95}
}

func caseHitAtK(c PerCase) bool {
	for _, h := range c.Hits {
		if _, ok := c.Expected[h.ID]; ok {
			return true
		}
	}
	return false
}
func firstHitRank(c PerCase) int {
	for i, h := range c.Hits {
		if _, ok := c.Expected[h.ID]; ok {
			return i + 1
		}
	}
	return 0
}

func minInt(a, b int) int {
	if a < b {
		return a
	}
	return b
}
