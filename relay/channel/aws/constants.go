package aws

import "strings"

var awsModelIDMap = map[string]string{
	// ── Anthropic Claude ─────────────────────────────────────────────────────
	"claude-3-sonnet-20240229":   "anthropic.claude-3-sonnet-20240229-v1:0",
	"claude-3-opus-20240229":     "anthropic.claude-3-opus-20240229-v1:0",
	"claude-3-haiku-20240307":    "anthropic.claude-3-haiku-20240307-v1:0",
	"claude-3-5-sonnet-20240620": "anthropic.claude-3-5-sonnet-20240620-v1:0",
	"claude-3-5-sonnet-20241022": "anthropic.claude-3-5-sonnet-20241022-v2:0",
	"claude-3-5-haiku-20241022":  "anthropic.claude-3-5-haiku-20241022-v1:0",
	"claude-3-7-sonnet-20250219": "anthropic.claude-3-7-sonnet-20250219-v1:0",
	"claude-sonnet-4-20250514":   "anthropic.claude-sonnet-4-20250514-v1:0",
	"claude-opus-4-20250514":     "anthropic.claude-opus-4-20250514-v1:0",
	"claude-opus-4-1-20250805":   "anthropic.claude-opus-4-1-20250805-v1:0",
	"claude-sonnet-4-5-20250929": "anthropic.claude-sonnet-4-5-20250929-v1:0",
	"claude-sonnet-4-6":          "anthropic.claude-sonnet-4-6",
	"claude-haiku-4-5-20251001":  "anthropic.claude-haiku-4-5-20251001-v1:0",
	"claude-opus-4-5-20251101":   "anthropic.claude-opus-4-5-20251101-v1:0",
	"claude-opus-4-6":            "anthropic.claude-opus-4-6-v1",

	// ── Amazon Nova 1 (InvokeModel path, messages-v1 schema) ─────────────────
	// Friendly names intentionally include ":" so isNovaModel detects them.
	"nova-micro-v1:0":   "amazon.nova-micro-v1:0",
	"nova-lite-v1:0":    "amazon.nova-lite-v1:0",
	"nova-pro-v1:0":     "amazon.nova-pro-v1:0",
	"nova-premier-v1:0": "amazon.nova-premier-v1:0",
	"nova-canvas-v1:0":  "amazon.nova-canvas-v1:0",
	"nova-reel-v1:0":    "amazon.nova-reel-v1:0",
	"nova-reel-v1:1":    "amazon.nova-reel-v1:1",
	"nova-sonic-v1:0":   "amazon.nova-sonic-v1:0",

	// ── Amazon Nova 2 (Converse API) ─────────────────────────────────────────
	"nova-2-lite":  "amazon.nova-2-lite",
	"nova-2-sonic": "amazon.nova-2-sonic",

	// ── Meta Llama ───────────────────────────────────────────────────────────
	"llama3-8b":       "meta.llama3-8b-instruct-v1:0",
	"llama3-70b":      "meta.llama3-70b-instruct-v1:0",
	"llama3.1-8b":     "meta.llama3-1-8b-instruct-v1:0",
	"llama3.1-70b":    "meta.llama3-1-70b-instruct-v1:0",
	"llama3.1-405b":   "meta.llama3-1-405b-instruct-v1:0",
	"llama3.2-1b":     "meta.llama3-2-1b-instruct-v1:0",
	"llama3.2-3b":     "meta.llama3-2-3b-instruct-v1:0",
	"llama3.2-11b":    "meta.llama3-2-11b-instruct-v1:0",
	"llama3.2-90b":    "meta.llama3-2-90b-instruct-v1:0",
	"llama3.3-70b":    "meta.llama3-3-70b-instruct-v1:0",
	"llama4-scout":    "meta.llama4-scout-17b-16e-instruct-v1:0",
	"llama4-maverick": "meta.llama4-maverick-17b-128e-instruct-v1:0",

	// ── Mistral AI ───────────────────────────────────────────────────────────
	"mistral-7b":         "mistral.mistral-7b-instruct-v0:2",
	"mixtral-8x7b":       "mistral.mixtral-8x7b-instruct-v0:1",
	"mistral-large":      "mistral.mistral-large-2402-v1:0",
	"mistral-large-2407": "mistral.mistral-large-2407-v1:0",
	"mistral-small":      "mistral.mistral-small-2402-v1:0",
	"mistral-large-3":    "mistral.mistral-large-3",
	"ministral-3b":       "mistral.ministral-3b",
	"ministral-8b":       "mistral.ministral-3-8b",
	"ministral-14b":      "mistral.ministral-14b-3-0",
	"magistral-small":    "mistral.magistral-small-2509",
	"devstral-2":         "mistral.devstral-2-123b",
	"pixtral-large":      "mistral.pixtral-large",

	// ── Cohere ───────────────────────────────────────────────────────────────
	"command-r":      "cohere.command-r-v1:0",
	"command-r-plus": "cohere.command-r-plus-v1:0",

	// ── AI21 Labs ─────────────────────────────────────────────────────────────
	"jamba-1.5-mini":  "ai21.jamba-1-5-mini-v1:0",
	"jamba-1.5-large": "ai21.jamba-1-5-large-v1:0",

	// ── DeepSeek ─────────────────────────────────────────────────────────────
	"deepseek-r1":   "deepseek.r1-v1:0",
	"deepseek-v3-1": "deepseek.deepseek-v3-1",
	"deepseek-v3-2": "deepseek.deepseek-v3-2",

	// ── Google Gemma ─────────────────────────────────────────────────────────
	"gemma-3-4b":  "google.gemma-3-4b-it",
	"gemma-3-12b": "google.gemma-3-12b-it",
	"gemma-3-27b": "google.gemma-3-27b-pt",

	// ── MiniMax ──────────────────────────────────────────────────────────────
	"minimax-m2":   "minimax.minimax-m2",
	"minimax-m2.1": "minimax.minimax-m2-1",

	// ── Moonshot Kimi ─────────────────────────────────────────────────────────
	"kimi-k2":         "moonshot.kimi-k2-5",
	"kimi-k2-thinking": "moonshot.kimi-k2-thinking",

	// ── NVIDIA Nemotron ──────────────────────────────────────────────────────
	"nemotron-nano-9b":  "nvidia.nemotron-nano-9b-v2",
	"nemotron-nano-12b": "nvidia.nemotron-nano-12b-v2-vl-bf16",
	"nemotron-nano-30b": "nvidia.nemotron-nano-3-30b",

	// ── Qwen ─────────────────────────────────────────────────────────────────
	"qwen3-32b":        "qwen.qwen3-32b",
	"qwen3-235b":       "qwen.qwen3-235b-a22b-2507",
	"qwen3-coder-30b":  "qwen.qwen3-coder-30b-a3b-instruct",
	"qwen3-coder-480b": "qwen.qwen3-coder-480b-a35b-instruct",

	// ── Writer ───────────────────────────────────────────────────────────────
	"palmyra-x4": "writer.palmyra-x4-v1:0",
	"palmyra-x5": "writer.palmyra-x5-v1:0",

	// ── Z.AI (GLM) ───────────────────────────────────────────────────────────
	"glm-4.7":       "zai.glm-4-7",
	"glm-4.7-flash": "zai.glm-4-7-flash",
}

var awsModelCanCrossRegionMap = map[string]map[string]bool{
	"anthropic.claude-3-sonnet-20240229-v1:0": {
		"us": true,
		"eu": true,
		"ap": true,
	},
	"anthropic.claude-3-opus-20240229-v1:0": {
		"us": true,
	},
	"anthropic.claude-3-haiku-20240307-v1:0": {
		"us": true,
		"eu": true,
		"ap": true,
	},
	"anthropic.claude-3-5-sonnet-20240620-v1:0": {
		"us": true,
		"eu": true,
		"ap": true,
	},
	"anthropic.claude-3-5-sonnet-20241022-v2:0": {
		"us": true,
		"ap": true,
	},
	"anthropic.claude-3-5-haiku-20241022-v1:0": {
		"us": true,
	},
	"anthropic.claude-3-7-sonnet-20250219-v1:0": {
		"us": true,
		"ap": true,
		"eu": true,
	},
	"anthropic.claude-sonnet-4-20250514-v1:0": {
		"us": true,
		"ap": true,
		"eu": true,
	},
	"anthropic.claude-opus-4-20250514-v1:0": {
		"us": true,
	},
	"anthropic.claude-opus-4-1-20250805-v1:0": {
		"us": true,
	},
	"anthropic.claude-sonnet-4-5-20250929-v1:0": {
		"us": true,
		"ap": true,
		"eu": true,
	},
	"anthropic.claude-sonnet-4-6": {
		"us": true,
		"ap": true,
		"eu": true,
	},
	"anthropic.claude-opus-4-5-20251101-v1:0": {
		"us": true,
		"ap": true,
		"eu": true,
	},
	"anthropic.claude-opus-4-6-v1": {
		"us": true,
		"ap": true,
		"eu": true,
	},
	"anthropic.claude-haiku-4-5-20251001-v1:0": {
		"us": true,
		"ap": true,
		"eu": true,
	},
	// Nova models - all support three major regions
	"amazon.nova-micro-v1:0": {
		"us":   true,
		"eu":   true,
		"apac": true,
	},
	"amazon.nova-lite-v1:0": {
		"us":   true,
		"eu":   true,
		"apac": true,
	},
	"amazon.nova-pro-v1:0": {
		"us":   true,
		"eu":   true,
		"apac": true,
	},
	"amazon.nova-premier-v1:0": {
		"us": true,
	},
	"amazon.nova-canvas-v1:0": {
		"us":   true,
		"eu":   true,
		"apac": true,
	},
	"amazon.nova-reel-v1:0": {
		"us":   true,
		"eu":   true,
		"apac": true,
	},
	"amazon.nova-reel-v1:1": {
		"us": true,
	},
	"amazon.nova-sonic-v1:0": {
		"us":   true,
		"eu":   true,
		"apac": true,
	},
}

var awsRegionCrossModelPrefixMap = map[string]string{
	"us": "us",
	"eu": "eu",
	"ap": "apac",
}

var ChannelName = "aws"

// isNovaModel returns true for Nova 1 models that use the InvokeModel path with
// the messages-v1 schema. Nova 1 friendly names always contain a colon version
// suffix (e.g. "nova-micro-v1:0"), while Nova 2 and later do not.
func isNovaModel(modelId string) bool {
	return strings.Contains(modelId, "nova-") && strings.ContainsRune(modelId, ':')
}

// isClaudeModel returns true for Anthropic Claude models.
func isClaudeModel(modelId string) bool {
	return strings.Contains(modelId, "claude")
}

// isConverseModel returns true for models that should use the Bedrock Converse API
// (i.e. not Claude and not Nova, which have their own InvokeModel paths).
func isConverseModel(modelId string) bool {
	return !isClaudeModel(modelId) && !isNovaModel(modelId)
}
