package tools

import (
	"encoding/json"
	"os"
	"sync"

	"github.com/anthropics/anthropic-sdk-go"
)

// Prices in USD per million tokens for claude-sonnet-5-5 (Anthropic first-party API).
// Cache writes use the 5-minute TTL rate (1.25x input).
const (
	priceInputPerMTok      = 2.00
	priceOutputPerMTok     = 10.00
	priceCacheWritePerMTok = 2.50
	priceCacheReadPerMTok  = 0.20
)

// UsageCounts holds the token counts and number of API calls for one role (agent or reviewer).
type UsageCounts struct {
	Turns            int     `json:"turns"` // number of model calls
	InputTokens      int64   `json:"input_tokens"`
	OutputTokens     int64   `json:"output_tokens"`
	CacheReadTokens  int64   `json:"cache_read_tokens"`
	CacheWriteTokens int64   `json:"cache_write_tokens"`
	CostUSD          float64 `json:"cost_usd"`
}

func (c *UsageCounts) add(u anthropic.Usage) {
	c.Turns++
	c.InputTokens += u.InputTokens
	c.OutputTokens += u.OutputTokens
	c.CacheReadTokens += u.CacheReadInputTokens
	c.CacheWriteTokens += u.CacheCreationInputTokens
	c.CostUSD = (float64(c.InputTokens)*priceInputPerMTok +
		float64(c.OutputTokens)*priceOutputPerMTok +
		float64(c.CacheWriteTokens)*priceCacheWritePerMTok +
		float64(c.CacheReadTokens)*priceCacheReadPerMTok) / 1e6
}

func (c UsageCounts) plus(o UsageCounts) UsageCounts {
	return UsageCounts{
		Turns:            c.Turns + o.Turns,
		InputTokens:      c.InputTokens + o.InputTokens,
		OutputTokens:     c.OutputTokens + o.OutputTokens,
		CacheReadTokens:  c.CacheReadTokens + o.CacheReadTokens,
		CacheWriteTokens: c.CacheWriteTokens + o.CacheWriteTokens,
		CostUSD:          c.CostUSD + o.CostUSD,
	}
}

// SessionStats is what gets written to the stats file at the end of a session.
type SessionStats struct {
	Model    string          `json:"model"`
	Features map[string]bool `json:"features"` // which optional features were on, e.g. review
	Reviews  int             `json:"reviews"`  // number of review_work calls
	Agent    UsageCounts     `json:"agent"`
	Reviewer UsageCounts     `json:"reviewer"`
	Total    UsageCounts     `json:"total"`
}

// UsageTracker counts turns and tokens for a whole session, including the reviewer.
type UsageTracker struct {
	mu       sync.Mutex
	model    string
	features map[string]bool
	reviews  int
	agent    UsageCounts
	reviewer UsageCounts
}

// Usage is the tracker for this process. The agent and the reviewer both record into it.
var Usage = &UsageTracker{
	model:    string(anthropic.BetaManagedAgentsModelClaudeSonnet5_5),
	features: map[string]bool{},
}

// SetFeature records whether an optional feature was on for this session.
func (t *UsageTracker) SetFeature(name string, on bool) {
	t.mu.Lock()
	defer t.mu.Unlock()
	t.features[name] = on
}

func (t *UsageTracker) RecordAgent(u anthropic.Usage) {
	t.mu.Lock()
	defer t.mu.Unlock()
	t.agent.add(u)
}

func (t *UsageTracker) RecordReviewer(u anthropic.Usage) {
	t.mu.Lock()
	defer t.mu.Unlock()
	t.reviewer.add(u)
}

func (t *UsageTracker) RecordReview() {
	t.mu.Lock()
	defer t.mu.Unlock()
	t.reviews++
}

func (t *UsageTracker) Stats() SessionStats {
	t.mu.Lock()
	defer t.mu.Unlock()
	return SessionStats{
		Model:    t.model,
		Features: t.features,
		Reviews:  t.reviews,
		Agent:    t.agent,
		Reviewer: t.reviewer,
		Total:    t.agent.plus(t.reviewer),
	}
}

// WriteFile writes the session stats as JSON to path.
func (t *UsageTracker) WriteFile(path string) error {
	data, err := json.MarshalIndent(t.Stats(), "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(path, append(data, '\n'), 0644)
}
