package bench

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"sync"
	"time"

	"gopkg.in/yaml.v3"
)

type Config struct {
	APIBase string  `yaml:"api_base"`
	K       int     `yaml:"k"`
	Queries []QItem `yaml:"queries"`
}

type QItem struct {
	ID             string   `yaml:"id"`
	Query          string   `yaml:"query"`
	ExpectedDocIDs []string `yaml:"expected_doc_ids"`
}

func LoadConfig(path string) (*Config, error) {
	b, err := osReadFile(path)
	if err != nil {
		return nil, err
	}
	var c Config
	if err := yaml.Unmarshal(b, &c); err != nil {
		return nil, err
	}
	if c.K <= 0 {
		c.K = 5
	}
	return &c, nil
}

// Abstraction of /api/rag/query response
type QueryRequest struct {
	Question string `json:"question"`
	K        int    `json:"k"`
}
type Hit struct {
	ID      string  `json:"id"`
	Score   float64 `json:"score"`
	Snippet string  `json:"snippet,omitempty"`
}
type QueryResponse struct {
	Answer string `json:"answer"`
	Hits   []Hit  `json:"hits"`
}

type Options struct {
	Parallel  int
	TimeoutMS int
}

type PerCase struct {
	ID       string
	Query    string
	K        int
	Hits     []Hit
	Latency  int64 // ms
	Error    error
	Expected map[string]struct{}
}

type SuiteResult struct {
	Cases   []PerCase
	Metrics Metrics
	Latency LatencyStats
	Errors  int
}

func RunSuite(cfg *Config, opt Options) *SuiteResult {
	if opt.Parallel <= 0 {
		opt.Parallel = 16
	}
	if opt.TimeoutMS <= 0 {
		opt.TimeoutMS = 8000
	}

	wg := sync.WaitGroup{}
	sem := make(chan struct{}, opt.Parallel)
	out := make([]PerCase, len(cfg.Queries))

	for i := range cfg.Queries {
		wg.Add(1)
		sem <- struct{}{}
		go func(i int) {
			defer wg.Done()
			defer func() { <-sem }()
			q := cfg.Queries[i]
			out[i] = runOne(cfg.APIBase, q, cfg.K, opt.TimeoutMS)
		}(i)
	}
	wg.Wait()

	m := CalcMetrics(out, cfg.K)
	lat := calcLatency(out)
	errs := 0
	for _, c := range out {
		if c.Error != nil {
			errs++
		}
	}

	return &SuiteResult{Cases: out, Metrics: m, Latency: lat, Errors: errs}
}

func runOne(apiBase string, q QItem, k, timeoutMS int) PerCase {
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(timeoutMS)*time.Millisecond)
	defer cancel()

	expected := make(map[string]struct{}, len(q.ExpectedDocIDs))
	for _, id := range q.ExpectedDocIDs {
		expected[id] = struct{}{}
	}

	reqBody, _ := json.Marshal(QueryRequest{Question: q.Query, K: k})
	url := fmt.Sprintf("%s/api/rag/query", apiBase)

	t0 := time.Now()
	resp, err := httpDo(ctx, "POST", url, "application/json", bytes.NewReader(reqBody))
	lat := time.Since(t0).Milliseconds()
	if err != nil {
		return PerCase{ID: q.ID, Query: q.Query, K: k, Latency: lat, Error: err, Expected: expected}
	}
	defer resp.Body.Close()
	b, _ := io.ReadAll(resp.Body)
	if resp.StatusCode >= 300 {
		return PerCase{ID: q.ID, Query: q.Query, K: k, Latency: lat, Error: fmt.Errorf("http %d: %s", resp.StatusCode, truncate(string(b), 200)), Expected: expected}
	}
	var qr QueryResponse
	if err := json.Unmarshal(b, &qr); err != nil {
		return PerCase{ID: q.ID, Query: q.Query, K: k, Latency: lat, Error: err, Expected: expected}
	}
	// normalize: trim to K and remove dups preserving order
	hits := dedupTopK(qr.Hits, k)

	return PerCase{ID: q.ID, Query: q.Query, K: k, Hits: hits, Latency: lat, Expected: expected}
}

func dedupTopK(h []Hit, k int) []Hit {
	seen := map[string]struct{}{}
	out := make([]Hit, 0, k)
	for _, x := range h {
		if x.ID == "" {
			continue
		}
		if _, ok := seen[x.ID]; ok {
			continue
		}
		seen[x.ID] = struct{}{}
		out = append(out, x)
		if len(out) >= k {
			break
		}
	}
	return out
}

// tiny helpers (no external deps)
func osReadFile(p string) ([]byte, error) { return os.ReadFile(p) }
func httpDo(ctx context.Context, method, url, ctype string, body io.Reader) (*http.Response, error) {
	req, _ := http.NewRequestWithContext(ctx, method, url, body)
	if ctype != "" {
		req.Header.Set("Content-Type", ctype)
	}
	return http.DefaultClient.Do(req)
}
func truncate(s string, n int) string {
	if len(s) <= n {
		return s
	}
	return s[:n] + "..."
}

// Markdown rendering
func RenderMarkdown(cfg *Config, res *SuiteResult) string {
	// summary
	md := &bytes.Buffer{}
	fmt.Fprintf(md, "# RAG Benchmark Report\n\n")
	fmt.Fprintf(md, "- API: `%s`\n- K: `%d`\n- Cases: `%d`\n- Errors: `%d`\n\n", cfg.APIBase, cfg.K, len(res.Cases), res.Errors)

	fmt.Fprintf(md, "## Summary Metrics\n\n")
	fmt.Fprintf(md, "| Metric | Value |\n|---|---|\n")
	fmt.Fprintf(md, "| Hit@%d | %.2f%% |\n", cfg.K, 100*res.Metrics.HitAtK)
	fmt.Fprintf(md, "| Recall@%d | %.2f%% |\n", cfg.K, 100*res.Metrics.RecallAtK)
	fmt.Fprintf(md, "| MRR | %.4f |\n", res.Metrics.MRR)
	fmt.Fprintf(md, "| nDCG@%d | %.4f |\n", cfg.K, res.Metrics.NDCGAtK)
	fmt.Fprintf(md, "| P50 latency | %d ms |\n", res.Latency.P50)
	fmt.Fprintf(md, "| P95 latency | %d ms |\n\n", res.Latency.P95)

	// failures table
	fail := failedCases(res.Cases)
	if len(fail) > 0 {
		fmt.Fprintf(md, "## Failures (%d)\n\n", len(fail))
		fmt.Fprintf(md, "| ID | Error |\n|---|---|\n")
		for _, f := range fail {
			fmt.Fprintf(md, "| %s | %s |\n", f.ID, truncate(fmt.Sprintf("%v", f.Error), 180))
		}
		fmt.Fprintln(md)
	}

	// per-case top-K (optional; keep short)
	fmt.Fprintf(md, "## Per-case (Top-%d IDs)\n\n", cfg.K)
	fmt.Fprintf(md, "| ID | Hit@K | Latency(ms) | TopIDs |\n|---|---:|---:|---|\n")
	for _, c := range res.Cases {
		hit := caseHitAtK(c)
		topIDs := make([]string, 0, len(c.Hits))
		for _, h := range c.Hits {
			topIDs = append(topIDs, h.ID)
		}
		fmt.Fprintf(md, "| %s | %t | %d | %s |\n", c.ID, hit, c.Latency, backtick(join(topIDs, ", ")))
	}
	return md.String()
}

func failedCases(cs []PerCase) []PerCase {
	out := make([]PerCase, 0)
	for _, c := range cs {
		if c.Error != nil {
			out = append(out, c)
		}
	}
	return out
}
func backtick(s string) string { return "`" + s + "`" }
func join(a []string, sep string) string {
	switch len(a) {
	case 0:
		return ""
	case 1:
		return a[0]
	default:
		b := a[0]
		for i := 1; i < len(a); i++ {
			b += sep + a[i]
		}
		return b
	}
}
