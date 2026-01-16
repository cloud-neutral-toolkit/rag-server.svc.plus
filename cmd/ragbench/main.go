package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"rag-server/internal/bench"
)

func main() {
	in := flag.String("in", "queries.yaml", "path to queries.yaml")
	api := flag.String("api", "", "override API base (optional, otherwise from yaml)")
	k := flag.Int("k", 0, "override K (optional, otherwise from yaml)")
	out := flag.String("out", "report.md", "output markdown report")
	parallel := flag.Int("parallel", 16, "concurrency")
	timeoutMS := flag.Int("timeout_ms", 8000, "per request timeout ms")
	flag.Parse()

	cfg, err := bench.LoadConfig(*in)
	if err != nil {
		log.Fatal(err)
	}
	if *api != "" {
		cfg.APIBase = *api
	}
	if *k > 0 {
		cfg.K = *k
	}

	res := bench.RunSuite(cfg, bench.Options{
		Parallel:  *parallel,
		TimeoutMS: *timeoutMS,
	})

	md := bench.RenderMarkdown(cfg, res)
	if err := os.WriteFile(*out, []byte(md), 0644); err != nil {
		log.Fatal(err)
	}

	fmt.Println("✅ RAG benchmark done. Report:", *out)
	fmt.Printf("Hit@%d: %.2f%%  Recall@%d: %.2f%%  MRR: %.4f  nDCG@%d: %.4f  P95: %dms  Errors: %d\n",
		cfg.K, 100*res.Metrics.HitAtK, cfg.K, 100*res.Metrics.RecallAtK,
		res.Metrics.MRR, cfg.K, res.Metrics.NDCGAtK, res.Latency.P95, res.Errors)
}
