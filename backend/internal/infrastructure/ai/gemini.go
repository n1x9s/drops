package ai

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/n1x9s/second-brain/backend/internal/domain"
)

const geminiBaseURL = "https://generativelanguage.googleapis.com/v1beta/models"

type GeminiProvider struct {
	apiKey         string
	model          string
	embeddingModel string
	client         *http.Client
}

func NewGeminiProvider(apiKey string, model string, embeddingModel string) GeminiProvider {
	return GeminiProvider{
		apiKey:         apiKey,
		model:          model,
		embeddingModel: embeddingModel,
		client:         &http.Client{Timeout: 12 * time.Second},
	}
}

func (p GeminiProvider) EnrichMemory(ctx context.Context, text string) (domain.Enrichment, error) {
	if p.apiKey == "" {
		return domain.Enrichment{}, domain.ErrProviderEmpty
	}
	prompt := `Return compact JSON for a personal memory. Schema:
{"summary":"one sentence","category":"Work|Learning|Personal|Projects|Meetings|Ideas","tags":["lowercase"]}
Text: ` + text

	raw, err := p.generate(ctx, prompt)
	if err != nil {
		return domain.Enrichment{}, err
	}
	var out struct {
		Summary  string   `json:"summary"`
		Category string   `json:"category"`
		Tags     []string `json:"tags"`
	}
	if err := json.Unmarshal([]byte(extractJSON(raw)), &out); err != nil {
		local := domain.LocalEnrichMemory(text)
		return local, nil
	}
	return domain.Enrichment{Summary: out.Summary, Category: domain.Category(out.Category), Tags: out.Tags}, nil
}

func (p GeminiProvider) ExtractTask(ctx context.Context, text string) (domain.TaskExtraction, error) {
	if p.apiKey == "" {
		return domain.TaskExtraction{}, domain.ErrProviderEmpty
	}
	prompt := `Return compact JSON for a task command. Dates must be RFC3339 or null. Schema:
{"title":"task title","due_at":null,"priority":"low|medium|high|urgent","tags":["lowercase"]}
Text: ` + text

	raw, err := p.generate(ctx, prompt)
	if err != nil {
		return domain.TaskExtraction{}, err
	}
	var out struct {
		Title    string   `json:"title"`
		DueAt    *string  `json:"due_at"`
		Priority string   `json:"priority"`
		Tags     []string `json:"tags"`
	}
	if err := json.Unmarshal([]byte(extractJSON(raw)), &out); err != nil {
		local := domain.LocalExtractTask(text)
		return local, nil
	}
	task := domain.TaskExtraction{Title: out.Title, Priority: domain.Priority(out.Priority), Tags: out.Tags}
	if out.DueAt != nil && *out.DueAt != "" {
		if parsed, err := time.Parse(time.RFC3339, *out.DueAt); err == nil {
			task.DueAt = &parsed
		}
	}
	return task, nil
}

func (p GeminiProvider) Embed(ctx context.Context, text string) ([]float32, error) {
	if p.apiKey == "" {
		return nil, domain.ErrProviderEmpty
	}
	body := map[string]any{
		"model": "models/" + p.embeddingModel,
		"content": map[string]any{
			"parts": []map[string]string{{"text": text}},
		},
		"outputDimensionality": 768,
	}
	payload, _ := json.Marshal(body)
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, fmt.Sprintf("%s/%s:embedContent?key=%s", geminiBaseURL, p.embeddingModel, p.apiKey), bytes.NewReader(payload))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	resp, err := p.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode >= 300 {
		return nil, fmt.Errorf("gemini embed failed: %s", resp.Status)
	}
	var out struct {
		Embedding struct {
			Values []float32 `json:"values"`
		} `json:"embedding"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&out); err != nil {
		return nil, err
	}
	return out.Embedding.Values, nil
}

func (p GeminiProvider) generate(ctx context.Context, prompt string) (string, error) {
	body := map[string]any{
		"contents": []map[string]any{
			{"parts": []map[string]string{{"text": prompt}}},
		},
		"generationConfig": map[string]any{
			"temperature":      0.2,
			"responseMimeType": "application/json",
		},
	}
	payload, _ := json.Marshal(body)
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, fmt.Sprintf("%s/%s:generateContent?key=%s", geminiBaseURL, p.model, p.apiKey), bytes.NewReader(payload))
	if err != nil {
		return "", err
	}
	req.Header.Set("Content-Type", "application/json")
	resp, err := p.client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	if resp.StatusCode >= 300 {
		return "", fmt.Errorf("gemini generate failed: %s", resp.Status)
	}
	var out struct {
		Candidates []struct {
			Content struct {
				Parts []struct {
					Text string `json:"text"`
				} `json:"parts"`
			} `json:"content"`
		} `json:"candidates"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&out); err != nil {
		return "", err
	}
	if len(out.Candidates) == 0 || len(out.Candidates[0].Content.Parts) == 0 {
		return "", fmt.Errorf("gemini response had no candidates")
	}
	return out.Candidates[0].Content.Parts[0].Text, nil
}

func extractJSON(raw string) string {
	raw = strings.TrimSpace(raw)
	raw = strings.TrimPrefix(raw, "```json")
	raw = strings.TrimPrefix(raw, "```")
	raw = strings.TrimSuffix(raw, "```")
	return strings.TrimSpace(raw)
}
