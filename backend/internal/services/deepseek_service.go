package services

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
	"time"
)

// DeepSeekService handles AI text generation using DeepSeek API
type DeepSeekService interface {
	GenerateRecognitionText(ctx context.Context, recipientName string, senderName string, senderDepartment string, userInput string, valueNames []string, language string) (string, error)
}

type deepseekService struct {
	apiKey     string
	baseURL    string
	httpClient *http.Client
}

// NewDeepSeekService creates a new DeepSeek service
func NewDeepSeekService() DeepSeekService {
	apiKey := os.Getenv("DEEPSEEK_API_KEY")
	if apiKey == "" {
		// Return a no-op service if API key is not configured
		log.Printf("WARNING: DEEPSEEK_API_KEY not configured, using no-op service")
		return &noOpDeepSeekService{}
	}

	baseURL := os.Getenv("DEEPSEEK_BASE_URL")
	if baseURL == "" {
		baseURL = "https://api.deepseek.com"
	}

	log.Printf("INFO: DeepSeek service initialized with base URL: %s (API key configured: %s...)", baseURL, apiKey[:min(10, len(apiKey))])
	return &deepseekService{
		apiKey:  apiKey,
		baseURL: baseURL,
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

// Helper function for min
func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

// GenerateRecognitionText generates a recognition text using DeepSeek API
func (s *deepseekService) GenerateRecognitionText(ctx context.Context, recipientName string, senderName string, senderDepartment string, userInput string, valueNames []string, language string) (string, error) {
	// Default to English if language is empty or invalid
	if language == "" {
		language = "en"
	}
	// Normalize language code - only support Chinese and English
	if language != "zh" && language != "en" {
		language = "en" // Default to English for unsupported languages
	}

	log.Printf("INFO: Generating recognition text for recipient: %s, sender: %s (%s), input: %s, values: %v, language: %s\n",
		recipientName, senderName, senderDepartment, userInput, valueNames, language)

	// Build language-specific prompt
	var prompt string
	if language == "zh" {
		// Chinese prompt
		valuesText := ""
		if len(valueNames) > 0 {
			valuesText = fmt.Sprintf("，并且要体现以下公司价值观：%s", joinStringsForDeepSeek(valueNames, "、"))
		}

		senderInfo := ""
		if senderName != "" {
			if senderDepartment != "" {
				senderInfo = fmt.Sprintf("发送人：%s（%s部门）\n", senderName, senderDepartment)
			} else {
				senderInfo = fmt.Sprintf("发送人：%s\n", senderName)
			}
		}

		prompt = fmt.Sprintf(`请根据以下信息，生成一段温暖、真诚、充满感激之情的感谢文本（彩虹屁风格），用于感谢同事。
%s接收人：%s
用户输入的关键信息：%s%s

要求：
1. 语言要温暖、真诚、充满感激之情
2. 要体现对接收人的认可和赞美
3. 如果提到了公司价值观，要自然地融入文本中
4. 长度控制在300-600字之间（不要超过600字）
5. 必须使用中文（简体中文）生成感谢文本
6. 语气要自然、亲切，不要太正式
7. 使用 STAR 原则（不需要特殊标题，大致按该原则输出）

请直接输出感谢文本，不要包含其他说明文字。`, senderInfo, recipientName, userInput, valuesText)
	} else {
		// English prompt
		valuesText := ""
		if len(valueNames) > 0 {
			valuesText = fmt.Sprintf(", and should reflect the following company values: %s", joinStringsForDeepSeek(valueNames, ", "))
		}

		senderInfo := ""
		if senderName != "" {
			if senderDepartment != "" {
				senderInfo = fmt.Sprintf("Sender: %s (%s Department)\n", senderName, senderDepartment)
			} else {
				senderInfo = fmt.Sprintf("Sender: %s\n", senderName)
			}
		}

		prompt = fmt.Sprintf(`Please generate a warm, sincere, and grateful thank you message (in a positive, appreciative style) for a colleague based on the following information:
%sRecipient: %s
User's key information: %s%s

Requirements:
1. The language should be warm, sincere, and full of gratitude
2. Show recognition and appreciation for the recipient
3. If company values are mentioned, naturally incorporate them into the text
4. Length should be between 300-600 words (do not exceed 600 words)
5. You MUST generate the thank you text in English. Do not use Chinese.
6. Tone should be natural and friendly, not too formal
7. Use STAR principle (Situation, Task, Action, Result) - do not use special headings, just follow this principle in the output naturally

Please output only the thank you text directly, without any additional explanations.`, senderInfo, recipientName, userInput, valuesText)
	}

	// Build request payload
	requestBody := map[string]interface{}{
		"model": "deepseek-chat",
		"messages": []map[string]interface{}{
			{
				"role":    "user",
				"content": prompt,
			},
		},
		"max_tokens":  800, // Limit to ensure we don't exceed 2000 chars (roughly 1 token = 2-3 chars for Chinese)
		"temperature": 0.8,
	}

	jsonData, err := json.Marshal(requestBody)
	if err != nil {
		return "", fmt.Errorf("failed to marshal request: %w", err)
	}

	// Create HTTP request
	url := fmt.Sprintf("%s/chat/completions", s.baseURL)
	log.Printf("INFO: Calling DeepSeek API: %s\n", url)
	req, err := http.NewRequestWithContext(ctx, "POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return "", fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", s.apiKey))
	log.Printf("INFO: Request headers set, sending to DeepSeek API...\n")

	// Send request
	resp, err := s.httpClient.Do(req)
	if err != nil {
		log.Printf("ERROR: Failed to send request to DeepSeek API: %v\n", err)
		return "", fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	log.Printf("INFO: DeepSeek API response status: %d\n", resp.StatusCode)

	// Read response
	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Printf("ERROR: Failed to read response body: %v\n", err)
		return "", fmt.Errorf("failed to read response: %w", err)
	}

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		log.Printf("ERROR: DeepSeek API returned error status %d: %s\n", resp.StatusCode, string(bodyBytes))
		return "", fmt.Errorf("API returned status %d: %s", resp.StatusCode, string(bodyBytes))
	}

	log.Printf("INFO: DeepSeek API response received successfully, body length: %d\n", len(bodyBytes))

	// Parse response
	var response struct {
		Choices []struct {
			Message struct {
				Content string `json:"content"`
			} `json:"message"`
		} `json:"choices"`
		Error struct {
			Message string `json:"message"`
		} `json:"error"`
	}

	if err := json.Unmarshal(bodyBytes, &response); err != nil {
		return "", fmt.Errorf("failed to parse response: %w", err)
	}

	if response.Error.Message != "" {
		return "", fmt.Errorf("API error: %s", response.Error.Message)
	}

	if len(response.Choices) == 0 {
		return "", fmt.Errorf("no choices in response")
	}

	// Trim and limit the generated text to 2000 characters
	generatedText := strings.TrimSpace(response.Choices[0].Message.Content)
	if len(generatedText) > 2000 {
		generatedText = generatedText[:2000]
		log.Printf("WARNING: Generated text exceeded 2000 chars, truncated to 2000")
	}

	return generatedText, nil
}

// noOpDeepSeekService is a no-op implementation when API key is not configured
type noOpDeepSeekService struct{}

func (s *noOpDeepSeekService) GenerateRecognitionText(ctx context.Context, recipientName string, senderName string, senderDepartment string, userInput string, valueNames []string, language string) (string, error) {
	// Default to English if language is empty or invalid
	if language == "" {
		language = "en"
	}
	if language != "zh" && language != "en" {
		language = "en"
	}

	// Return a simple template when API is not configured
	// For simplicity, only handle Chinese and English in no-op mode
	if language == "zh" {
		valuesText := ""
		if len(valueNames) > 0 {
			valuesText = fmt.Sprintf("，特别是你在 %s 方面的表现", joinStringsForDeepSeek(valueNames, "、"))
		}
		senderPrefix := ""
		if senderName != "" {
			senderPrefix = fmt.Sprintf("来自 %s", senderName)
			if senderDepartment != "" {
				senderPrefix = fmt.Sprintf("来自 %s（%s部门）", senderName, senderDepartment)
			}
			senderPrefix += "："
		}
		return fmt.Sprintf("%s感谢 %s%s。%s", senderPrefix, recipientName, valuesText, userInput), nil
	} else {
		// Default to English for all other languages in no-op mode
		valuesText := ""
		if len(valueNames) > 0 {
			valuesText = fmt.Sprintf(", especially your performance in %s", joinStringsForDeepSeek(valueNames, ", "))
		}
		senderPrefix := ""
		if senderName != "" {
			senderPrefix = fmt.Sprintf("From %s", senderName)
			if senderDepartment != "" {
				senderPrefix = fmt.Sprintf("From %s (%s Department)", senderName, senderDepartment)
			}
			senderPrefix += ": "
		}
		return fmt.Sprintf("%sThank you %s%s. %s", senderPrefix, recipientName, valuesText, userInput), nil
	}
}

// Helper function to join strings with separator
func joinStringsForDeepSeek(strs []string, sep string) string {
	if len(strs) == 0 {
		return ""
	}
	if len(strs) == 1 {
		return strs[0]
	}
	result := strs[0]
	for i := 1; i < len(strs); i++ {
		result += sep + strs[i]
	}
	return result
}
