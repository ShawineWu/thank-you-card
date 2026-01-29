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

// buildPromptForLanguage builds a prompt in the specified language
func buildPromptForLanguage(language, recipientName, senderName, senderDepartment, userInput string, valueNames []string) string {
	// Build values text and sender info based on language
	var valuesText, senderInfo string
	var promptTemplate string

	// Language-specific configurations
	switch language {
	case "zh": // Chinese
		if len(valueNames) > 0 {
			valuesText = fmt.Sprintf("，并且要体现以下公司价值观：%s", joinStringsForDeepSeek(valueNames, "、"))
		}
		if senderName != "" {
			if senderDepartment != "" {
				senderInfo = fmt.Sprintf("发送人：%s（%s部门）\n", senderName, senderDepartment)
			} else {
				senderInfo = fmt.Sprintf("发送人：%s\n", senderName)
			}
		}
		promptTemplate = `请根据以下信息，生成一段温暖、真诚、充满感激之情的感谢文本（彩虹屁风格），用于感谢同事。
%s接收人：%s
用户输入的关键信息：%s%s

要求：
1. 语言要温暖、真诚、充满感激之情
2. 要体现对接收人的认可和赞美
3. 如果提到了公司价值观，要自然地融入文本中
4. 长度控制在300-600字之间（不要超过600字）
5. 必须使用中文（简体中文）生成感谢文本
6. 语气要自然、亲切，不要太正式
7. 小红书式风格，可以适当使用emoji表情符号

请直接输出感谢文本，不要包含其他说明文字。`

	case "ja": // Japanese
		if len(valueNames) > 0 {
			valuesText = fmt.Sprintf("、以下の会社の価値観を反映してください：%s", joinStringsForDeepSeek(valueNames, "、"))
		}
		if senderName != "" {
			if senderDepartment != "" {
				senderInfo = fmt.Sprintf("送信者：%s（%s部門）\n", senderName, senderDepartment)
			} else {
				senderInfo = fmt.Sprintf("送信者：%s\n", senderName)
			}
		}
		promptTemplate = `以下の情報に基づいて、同僚への感謝のメッセージ（温かく、誠実で、感謝の気持ちに満ちた）を生成してください。
%s受信者：%s
ユーザーが入力した重要な情報：%s%s

要件：
1. 温かく、誠実で、感謝の気持ちに満ちた言語を使用する
2. 受信者への認識と賞賛を示す
3. 会社の価値観が言及されている場合は、自然にテキストに組み込む
4. 長さは300-600文字の間（600文字を超えない）
5. 日本語で感謝のテキストを生成する必要があります。他の言語は使用しないでください。
6. トーンは自然で親しみやすく、あまりフォーマルすぎない
7. カジュアルなソーシャルメディアスタイルで、適切な絵文字を使用

感謝のテキストのみを直接出力してください。追加の説明は含めないでください。`

	case "ko": // Korean
		if len(valueNames) > 0 {
			valuesText = fmt.Sprintf(", 다음 회사 가치를 반영해야 합니다: %s", joinStringsForDeepSeek(valueNames, ", "))
		}
		if senderName != "" {
			if senderDepartment != "" {
				senderInfo = fmt.Sprintf("발신자: %s (%s 부서)\n", senderName, senderDepartment)
			} else {
				senderInfo = fmt.Sprintf("발신자: %s\n", senderName)
			}
		}
		promptTemplate = `다음 정보를 기반으로 동료에게 따뜻하고 진심 어린 감사 메시지를 생성하세요.
%s수신자: %s
사용자가 입력한 주요 정보: %s%s

요구사항:
1. 언어는 따뜻하고 진심 어리며 감사로 가득해야 합니다
2. 수신자에 대한 인정과 찬사를 보여주세요
3. 회사 가치가 언급된 경우 자연스럽게 텍스트에 통합하세요
4. 길이는 300-600자 사이여야 합니다 (600자를 초과하지 마세요)
5. 한국어로 감사 텍스트를 생성해야 합니다. 다른 언어를 사용하지 마세요.
6. 톤은 자연스럽고 친근하며 너무 격식적이지 않아야 합니다
7. 캐주얼한 소셜 미디어 스타일로 적절한 이모지를 사용하세요

감사 텍스트만 직접 출력하세요. 추가 설명을 포함하지 마세요.`

	case "fr": // French
		if len(valueNames) > 0 {
			valuesText = fmt.Sprintf(", et doit refléter les valeurs suivantes de l'entreprise: %s", joinStringsForDeepSeek(valueNames, ", "))
		}
		if senderName != "" {
			if senderDepartment != "" {
				senderInfo = fmt.Sprintf("Expéditeur: %s (Département %s)\n", senderName, senderDepartment)
			} else {
				senderInfo = fmt.Sprintf("Expéditeur: %s\n", senderName)
			}
		}
		promptTemplate = `Veuillez générer un message de remerciement chaleureux, sincère et reconnaissant (dans un style positif et appréciatif) pour un collègue basé sur les informations suivantes:
%sDestinataire: %s
Informations clés de l'utilisateur: %s%s

Exigences:
1. Le langage doit être chaleureux, sincère et plein de gratitude
2. Montrer la reconnaissance et l'appréciation pour le destinataire
3. Si des valeurs de l'entreprise sont mentionnées, les intégrer naturellement dans le texte
4. La longueur doit être entre 300-600 mots (ne pas dépasser 600 mots)
5. Vous DEVEZ générer le texte de remerciement en français. N'utilisez pas d'autres langues.
6. Le ton doit être naturel et amical, pas trop formel
7. Utilisez un style décontracté de médias sociaux avec des emojis appropriés

Veuillez sortir uniquement le texte de remerciement directement, sans explications supplémentaires.`

	case "de": // German
		if len(valueNames) > 0 {
			valuesText = fmt.Sprintf(" und sollte die folgenden Unternehmenswerte widerspiegeln: %s", joinStringsForDeepSeek(valueNames, ", "))
		}
		if senderName != "" {
			if senderDepartment != "" {
				senderInfo = fmt.Sprintf("Absender: %s (%s Abteilung)\n", senderName, senderDepartment)
			} else {
				senderInfo = fmt.Sprintf("Absender: %s\n", senderName)
			}
		}
		promptTemplate = `Bitte generieren Sie eine warme, aufrichtige und dankbare Dankesnachricht (in einem positiven, wertschätzenden Stil) für einen Kollegen basierend auf den folgenden Informationen:
%sEmpfänger: %s
Wichtige Informationen des Benutzers: %s%s

Anforderungen:
1. Die Sprache sollte warm, aufrichtig und voller Dankbarkeit sein
2. Anerkennung und Wertschätzung für den Empfänger zeigen
3. Wenn Unternehmenswerte erwähnt werden, diese natürlich in den Text einbauen
4. Die Länge sollte zwischen 300-600 Wörtern liegen (600 Wörter nicht überschreiten)
5. Sie MÜSSEN den Dankestext auf Deutsch generieren. Verwenden Sie keine anderen Sprachen.
6. Der Ton sollte natürlich und freundlich sein, nicht zu formal
7. Verwenden Sie einen lockeren Social-Media-Stil mit angemessenen Emojis

Bitte geben Sie nur den Dankestext direkt aus, ohne zusätzliche Erklärungen.`

	case "es": // Spanish
		if len(valueNames) > 0 {
			valuesText = fmt.Sprintf(", y debe reflejar los siguientes valores de la empresa: %s", joinStringsForDeepSeek(valueNames, ", "))
		}
		if senderName != "" {
			if senderDepartment != "" {
				senderInfo = fmt.Sprintf("Remitente: %s (Departamento %s)\n", senderName, senderDepartment)
			} else {
				senderInfo = fmt.Sprintf("Remitente: %s\n", senderName)
			}
		}
		promptTemplate = `Por favor, genera un mensaje de agradecimiento cálido, sincero y agradecido (en un estilo positivo y apreciativo) para un colega basado en la siguiente información:
%sDestinatario: %s
Información clave del usuario: %s%s

Requisitos:
1. El lenguaje debe ser cálido, sincero y lleno de gratitud
2. Mostrar reconocimiento y aprecio por el destinatario
3. Si se mencionan valores de la empresa, incorporarlos naturalmente en el texto
4. La longitud debe estar entre 300-600 palabras (no exceder 600 palabras)
5. DEBES generar el texto de agradecimiento en español. No uses otros idiomas.
6. El tono debe ser natural y amigable, no demasiado formal
7. Usa un estilo casual de redes sociales con emojis apropiados

Por favor, genera solo el texto de agradecimiento directamente, sin explicaciones adicionales.`

	case "pt": // Portuguese
		if len(valueNames) > 0 {
			valuesText = fmt.Sprintf(", e deve refletir os seguintes valores da empresa: %s", joinStringsForDeepSeek(valueNames, ", "))
		}
		if senderName != "" {
			if senderDepartment != "" {
				senderInfo = fmt.Sprintf("Remetente: %s (Departamento %s)\n", senderName, senderDepartment)
			} else {
				senderInfo = fmt.Sprintf("Remetente: %s\n", senderName)
			}
		}
		promptTemplate = `Por favor, gere uma mensagem de agradecimento calorosa, sincera e grata (em um estilo positivo e apreciativo) para um colega com base nas seguintes informações:
%sDestinatário: %s
Informações-chave do usuário: %s%s

Requisitos:
1. A linguagem deve ser calorosa, sincera e cheia de gratidão
2. Mostrar reconhecimento e apreço pelo destinatário
3. Se os valores da empresa forem mencionados, incorporá-los naturalmente no texto
4. O comprimento deve estar entre 300-600 palavras (não exceder 600 palavras)
5. Você DEVE gerar o texto de agradecimento em português. Não use outros idiomas.
6. O tom deve ser natural e amigável, não muito formal
7. Use um estilo casual de mídia social com emojis apropriados

Por favor, gere apenas o texto de agradecimento diretamente, sem explicações adicionais.`

	case "it": // Italian
		if len(valueNames) > 0 {
			valuesText = fmt.Sprintf(", e deve riflettere i seguenti valori aziendali: %s", joinStringsForDeepSeek(valueNames, ", "))
		}
		if senderName != "" {
			if senderDepartment != "" {
				senderInfo = fmt.Sprintf("Mittente: %s (Dipartimento %s)\n", senderName, senderDepartment)
			} else {
				senderInfo = fmt.Sprintf("Mittente: %s\n", senderName)
			}
		}
		promptTemplate = `Per favore, genera un messaggio di ringraziamento caloroso, sincero e grato (in uno stile positivo e apprezzativo) per un collega basato sulle seguenti informazioni:
%sDestinatario: %s
Informazioni chiave dell'utente: %s%s

Requisiti:
1. Il linguaggio deve essere caloroso, sincero e pieno di gratitudine
2. Mostrare riconoscimento e apprezzamento per il destinatario
3. Se vengono menzionati i valori aziendali, incorporarli naturalmente nel testo
4. La lunghezza deve essere tra 300-600 parole (non superare 600 parole)
5. DEVI generare il testo di ringraziamento in italiano. Non usare altre lingue.
6. Il tono deve essere naturale e amichevole, non troppo formale
7. Usa uno stile casual dei social media con emoji appropriati

Per favore, genera solo il testo di ringraziamento direttamente, senza spiegazioni aggiuntive.`

	case "ru": // Russian
		if len(valueNames) > 0 {
			valuesText = fmt.Sprintf(", и должно отражать следующие корпоративные ценности: %s", joinStringsForDeepSeek(valueNames, ", "))
		}
		if senderName != "" {
			if senderDepartment != "" {
				senderInfo = fmt.Sprintf("Отправитель: %s (Отдел %s)\n", senderName, senderDepartment)
			} else {
				senderInfo = fmt.Sprintf("Отправитель: %s\n", senderName)
			}
		}
		promptTemplate = `Пожалуйста, создайте теплое, искреннее и благодарное сообщение (в позитивном, признательном стиле) для коллеги на основе следующей информации:
%sПолучатель: %s
Ключевая информация пользователя: %s%s

Требования:
1. Язык должен быть теплым, искренним и полным благодарности
2. Показать признание и благодарность получателю
3. Если упоминаются корпоративные ценности, естественно включить их в текст
4. Длина должна быть от 300-600 слов (не превышать 600 слов)
5. Вы ДОЛЖНЫ создать текст благодарности на русском языке. Не используйте другие языки.
6. Тон должен быть естественным и дружелюбным, не слишком формальным
7. Используйте непринужденный стиль социальных сетей с подходящими эмодзи

Пожалуйста, выведите только текст благодарности напрямую, без дополнительных объяснений.`

	default: // English (default)
		if len(valueNames) > 0 {
			valuesText = fmt.Sprintf(", and should reflect the following company values: %s", joinStringsForDeepSeek(valueNames, ", "))
		}
		if senderName != "" {
			if senderDepartment != "" {
				senderInfo = fmt.Sprintf("Sender: %s (%s Department)\n", senderName, senderDepartment)
			} else {
				senderInfo = fmt.Sprintf("Sender: %s\n", senderName)
			}
		}
		promptTemplate = `Please generate a warm, sincere, and grateful thank you message (in a positive, appreciative style) for a colleague based on the following information:
%sRecipient: %s
User's key information: %s%s

Requirements:
1. The language should be warm, sincere, and full of gratitude
2. Show recognition and appreciation for the recipient
3. If company values are mentioned, naturally incorporate them into the text
4. Length should be between 300-600 words (do not exceed 600 words)
5. You MUST generate the thank you text in English. Do not use other languages.
6. Tone should be natural and friendly, not too formal
7. Use a casual, social media style with appropriate emojis

Please output only the thank you text directly, without any additional explanations.`
	}

	return fmt.Sprintf(promptTemplate, senderInfo, recipientName, userInput, valuesText)
}

// GenerateRecognitionText generates a recognition text using DeepSeek API
func (s *deepseekService) GenerateRecognitionText(ctx context.Context, recipientName string, senderName string, senderDepartment string, userInput string, valueNames []string, language string) (string, error) {
	// Default to English if language is empty or invalid
	if language == "" {
		language = "en"
	}
	// Supported languages: zh, en, ja, ko, fr, de, es, pt, it, ru
	supportedLanguages := map[string]bool{
		"zh": true, // Chinese
		"en": true, // English
		"ja": true, // Japanese
		"ko": true, // Korean
		"fr": true, // French
		"de": true, // German
		"es": true, // Spanish
		"pt": true, // Portuguese
		"it": true, // Italian
		"ru": true, // Russian
	}
	if !supportedLanguages[language] {
		language = "en" // Default to English for unsupported languages
	}

	log.Printf("INFO: Generating recognition text for recipient: %s, sender: %s (%s), input: %s, values: %v, language: %s\n",
		recipientName, senderName, senderDepartment, userInput, valueNames, language)

	// Build language-specific prompt
	prompt := buildPromptForLanguage(language, recipientName, senderName, senderDepartment, userInput, valueNames)

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
	supportedLanguages := map[string]bool{
		"zh": true, "en": true, "ja": true, "ko": true, "fr": true,
		"de": true, "es": true, "pt": true, "it": true, "ru": true,
	}
	if !supportedLanguages[language] {
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
