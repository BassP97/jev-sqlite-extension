package jev

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
)

var URL = "https://api.typesafe.ai/v1/systemone"

func Classify(ctx context.Context, apiKey, text, question string, options map[string]string) (string, float64, error) {
	body, err := json.Marshal(map[string]any{
		"model": "jev-latest",
		"state": text,
		"questions": map[string]any{
			"q": map[string]any{"type": "choice", "instructions": question, "criteria": options},
		},
	})
	if err != nil {
		return "", 0, err
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, URL, bytes.NewReader(body))
	if err != nil {
		return "", 0, err
	}
	req.Header.Set("Authorization", "Bearer "+apiKey)
	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", 0, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return "", 0, fmt.Errorf("jev: HTTP %d", resp.StatusCode)
	}

	var out struct {
		Answers struct {
			Q struct {
				Choice     string  `json:"choice"`
				Confidence float64 `json:"confidence"`
			} `json:"q"`
		} `json:"answers"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&out); err != nil {
		return "", 0, err
	}
	return out.Answers.Q.Choice, out.Answers.Q.Confidence, nil
}
