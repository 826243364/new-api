package aws

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"strings"

	"github.com/QuantumNous/new-api/common"
	"github.com/QuantumNous/new-api/dto"
	"github.com/QuantumNous/new-api/logger"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/bedrockruntime"
	bedrockruntimeTypes "github.com/aws/aws-sdk-go-v2/service/bedrockruntime/types"
)

type AwsClaudeRequest struct {
	// AnthropicVersion should be "bedrock-2023-05-31"
	AnthropicVersion string              `json:"anthropic_version"`
	AnthropicBeta    json.RawMessage     `json:"anthropic_beta,omitempty"`
	System           any                 `json:"system,omitempty"`
	Messages         []dto.ClaudeMessage `json:"messages"`
	MaxTokens        uint                `json:"max_tokens,omitempty"`
	Temperature      *float64            `json:"temperature,omitempty"`
	TopP             float64             `json:"top_p,omitempty"`
	TopK             int                 `json:"top_k,omitempty"`
	StopSequences    []string            `json:"stop_sequences,omitempty"`
	Tools            any                 `json:"tools,omitempty"`
	ToolChoice       any                 `json:"tool_choice,omitempty"`
	Thinking         *dto.Thinking       `json:"thinking,omitempty"`
	OutputConfig     json.RawMessage     `json:"output_config,omitempty"`
	//Metadata         json.RawMessage     `json:"metadata,omitempty"`
}

func formatRequest(requestBody io.Reader, requestHeader http.Header) (*AwsClaudeRequest, error) {
	var awsClaudeRequest AwsClaudeRequest
	err := common.DecodeJson(requestBody, &awsClaudeRequest)
	if err != nil {
		return nil, err
	}
	awsClaudeRequest.AnthropicVersion = "bedrock-2023-05-31"

	// check header anthropic-beta
	anthropicBetaValues := requestHeader.Get("anthropic-beta")
	if len(anthropicBetaValues) > 0 {
		var tempArray []string
		tempArray = strings.Split(anthropicBetaValues, ",")
		if len(tempArray) > 0 {
			betaJson, err := json.Marshal(tempArray)
			if err != nil {
				return nil, err
			}
			awsClaudeRequest.AnthropicBeta = betaJson
		}
	}
	logger.LogJson(context.Background(), "json", awsClaudeRequest)
	return &awsClaudeRequest, nil
}

// NovaMessage Nova模型使用messages-v1格式
type NovaMessage struct {
	Role    string        `json:"role"`
	Content []NovaContent `json:"content"`
}

type NovaContent struct {
	Text string `json:"text"`
}

type NovaRequest struct {
	SchemaVersion   string               `json:"schemaVersion"`             // 请求版本，例如 "1.0"
	Messages        []NovaMessage        `json:"messages"`                  // 对话消息列表
	InferenceConfig *NovaInferenceConfig `json:"inferenceConfig,omitempty"` // 推理配置，可选
}

type NovaInferenceConfig struct {
	MaxTokens     int      `json:"maxTokens,omitempty"`     // 最大生成的 token 数
	Temperature   float64  `json:"temperature,omitempty"`   // 随机性 (默认 0.7, 范围 0-1)
	TopP          float64  `json:"topP,omitempty"`          // nucleus sampling (默认 0.9, 范围 0-1)
	TopK          int      `json:"topK,omitempty"`          // 限制候选 token 数 (默认 50, 范围 0-128)
	StopSequences []string `json:"stopSequences,omitempty"` // 停止生成的序列
}

// 转换OpenAI请求为Nova格式
func convertToNovaRequest(req *dto.GeneralOpenAIRequest) *NovaRequest {
	novaMessages := make([]NovaMessage, len(req.Messages))
	for i, msg := range req.Messages {
		novaMessages[i] = NovaMessage{
			Role:    msg.Role,
			Content: []NovaContent{{Text: msg.StringContent()}},
		}
	}

	novaReq := &NovaRequest{
		SchemaVersion: "messages-v1",
		Messages:      novaMessages,
	}

	// 设置推理配置
	if (req.MaxTokens != nil && *req.MaxTokens != 0) || (req.Temperature != nil && *req.Temperature != 0) || (req.TopP != nil && *req.TopP != 0) || (req.TopK != nil && *req.TopK != 0) || req.Stop != nil {
		novaReq.InferenceConfig = &NovaInferenceConfig{}
		if req.MaxTokens != nil && *req.MaxTokens != 0 {
			novaReq.InferenceConfig.MaxTokens = int(*req.MaxTokens)
		}
		if req.Temperature != nil && *req.Temperature != 0 {
			novaReq.InferenceConfig.Temperature = *req.Temperature
		}
		if req.TopP != nil && *req.TopP != 0 {
			novaReq.InferenceConfig.TopP = *req.TopP
		}
		if req.TopK != nil && *req.TopK != 0 {
			novaReq.InferenceConfig.TopK = *req.TopK
		}
		if req.Stop != nil {
			if stopSequences := parseStopSequences(req.Stop); len(stopSequences) > 0 {
				novaReq.InferenceConfig.StopSequences = stopSequences
			}
		}
	}

	return novaReq
}

// parseStopSequences 解析停止序列，支持字符串或字符串数组
func parseStopSequences(stop any) []string {
	if stop == nil {
		return nil
	}

	switch v := stop.(type) {
	case string:
		if v != "" {
			return []string{v}
		}
	case []string:
		return v
	case []interface{}:
		var sequences []string
		for _, item := range v {
			if str, ok := item.(string); ok && str != "" {
				sequences = append(sequences, str)
			}
		}
		return sequences
	}
	return nil
}

// convertToConverseMessages converts OpenAI messages to Bedrock Converse format.
// System messages are extracted as separate system blocks.
func convertToConverseMessages(messages []dto.Message) ([]bedrockruntimeTypes.Message, []bedrockruntimeTypes.SystemContentBlock) {
	var converseMessages []bedrockruntimeTypes.Message
	var systemBlocks []bedrockruntimeTypes.SystemContentBlock

	for _, msg := range messages {
		if msg.Role == "system" {
			systemBlocks = append(systemBlocks, &bedrockruntimeTypes.SystemContentBlockMemberText{
				Value: msg.StringContent(),
			})
			continue
		}

		role := bedrockruntimeTypes.ConversationRoleUser
		if msg.Role == "assistant" {
			role = bedrockruntimeTypes.ConversationRoleAssistant
		}

		converseMessages = append(converseMessages, bedrockruntimeTypes.Message{
			Role: role,
			Content: []bedrockruntimeTypes.ContentBlock{
				&bedrockruntimeTypes.ContentBlockMemberText{
					Value: msg.StringContent(),
				},
			},
		})
	}
	return converseMessages, systemBlocks
}

func buildConverseInferenceConfig(req *dto.GeneralOpenAIRequest) *bedrockruntimeTypes.InferenceConfiguration {
	config := &bedrockruntimeTypes.InferenceConfiguration{}
	hasConfig := false

	if req.MaxTokens != nil && *req.MaxTokens > 0 {
		maxTokens := int32(*req.MaxTokens)
		config.MaxTokens = &maxTokens
		hasConfig = true
	}
	if req.Temperature != nil {
		temp := float32(*req.Temperature)
		config.Temperature = &temp
		hasConfig = true
	}
	if req.TopP != nil {
		topP := float32(*req.TopP)
		config.TopP = &topP
		hasConfig = true
	}
	if req.Stop != nil {
		if stops := parseStopSequences(req.Stop); len(stops) > 0 {
			config.StopSequences = stops
			hasConfig = true
		}
	}

	if !hasConfig {
		return nil
	}
	return config
}

// convertToConverseInput builds a Bedrock ConverseInput from an OpenAI chat request.
func convertToConverseInput(modelId string, req *dto.GeneralOpenAIRequest) *bedrockruntime.ConverseInput {
	messages, systemBlocks := convertToConverseMessages(req.Messages)
	input := &bedrockruntime.ConverseInput{
		ModelId:  aws.String(modelId),
		Messages: messages,
	}
	if len(systemBlocks) > 0 {
		input.System = systemBlocks
	}
	if config := buildConverseInferenceConfig(req); config != nil {
		input.InferenceConfig = config
	}
	return input
}

// convertToConverseStreamInput builds a Bedrock ConverseStreamInput from an OpenAI chat request.
func convertToConverseStreamInput(modelId string, req *dto.GeneralOpenAIRequest) *bedrockruntime.ConverseStreamInput {
	messages, systemBlocks := convertToConverseMessages(req.Messages)
	input := &bedrockruntime.ConverseStreamInput{
		ModelId:  aws.String(modelId),
		Messages: messages,
	}
	if len(systemBlocks) > 0 {
		input.System = systemBlocks
	}
	if config := buildConverseInferenceConfig(req); config != nil {
		input.InferenceConfig = config
	}
	return input
}
