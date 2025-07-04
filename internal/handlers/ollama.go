package handlers

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
)

func OllamaQuery(query string, modelName string, ollamaServer string) (string, error) {
	url := fmt.Sprintf("%s/api/generate", ollamaServer)
	payload := map[string]interface{}{
		"model":  modelName,
		"prompt": query,
		"stream": false,
	}

	jsonPayload, err := json.Marshal(payload)
	if err != nil {
		return "", fmt.Errorf("failed to marshal payload: %w", err)
	}

	resp, err := http.Post(url, "application/json", strings.NewReader(string(jsonPayload)))
	if err != nil {
		return "", fmt.Errorf("failed to make HTTP request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		errorBody, _ := io.ReadAll(resp.Body)
		return "", fmt.Errorf("HTTP request failed with status %d: %s", resp.StatusCode, string(errorBody))
	}

	scanner := bufio.NewScanner(resp.Body)
	var fullResponse strings.Builder

	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			continue
		}

		var response map[string]interface{}
		if err := json.Unmarshal([]byte(line), &response); err != nil {
			continue
		}

		if responseText, ok := response["response"].(string); ok {
			fullResponse.WriteString(responseText)
		}

		if done, ok := response["done"].(bool); ok && done {
			break
		}
	}

	if err := scanner.Err(); err != nil {
		return "", fmt.Errorf("failed to read response: %w", err)
	}

	result := fullResponse.String()
	if result == "" {
		return "", fmt.Errorf("no response text received from Ollama")
	}

	return result, nil
}

// OllamaQueryStream streams responses from Ollama in real-time
func OllamaQueryStream(w http.ResponseWriter, query string, modelName string, ollamaServer string) error {
	url := fmt.Sprintf("%s/api/generate", ollamaServer)
	payload := map[string]interface{}{
		"model":  modelName,
		"prompt": query,
		"stream": true,
	}

	jsonPayload, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("failed to marshal payload: %w", err)
	}

	resp, err := http.Post(url, "application/json", strings.NewReader(string(jsonPayload)))
	if err != nil {
		return fmt.Errorf("failed to make HTTP request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		errorBody, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("HTTP request failed with status %d: %s", resp.StatusCode, string(errorBody))
	}

	// Set headers for Server-Sent Events
	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	// Create a flusher to ensure data is sent immediately
	flusher, ok := w.(http.Flusher)
	if !ok {
		return fmt.Errorf("streaming not supported")
	}

	scanner := bufio.NewScanner(resp.Body)

	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			continue
		}

		var response map[string]interface{}
		if err := json.Unmarshal([]byte(line), &response); err != nil {
			continue
		}

		if responseText, ok := response["response"].(string); ok {
			// Send the response chunk as a Server-Sent Event
			// Don't escape HTML here - let the client handle it
			fmt.Fprintf(w, "data: %s\n\n", responseText)
			flusher.Flush()
		}

		if done, ok := response["done"].(bool); ok && done {
			// Send end signal
			fmt.Fprintf(w, "data: [DONE]\n\n")
			flusher.Flush()
			break
		}
	}

	if err := scanner.Err(); err != nil {
		return fmt.Errorf("failed to read response: %w", err)
	}

	return nil
}
