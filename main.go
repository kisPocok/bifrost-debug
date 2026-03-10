package main

import (
	"fmt"
)

const initPrompt = `Mi az a makos guba? 1 mondatban válaszolj!`
const ANTHROPIC_API_KEY = ""
const OPENAI_API_KEY = ""

func main() {
	fmt.Println("OpenAI Message... ")
	openAIMessage("gpt-4o-mini")
	// [usage]
	// prompt_tokens=24  completion_tokens=50  total_tokens=74
	// cached_tokens=0  prompt_audio_tokens=0

	fmt.Println("\n\nAnthropic Message... ")
	anthropicMessage("claude-sonnet-4-6")
	// [usage] input_tokens=0  output_tokens=65  cache_creation=0  cache_read=0

	fmt.Println("\n\nBifrost SDK (Anthropic) Message... ")
	bifrostAnthropicMessage("claude-sonnet-4-6")
	// [usage] prompt_tokens=26  completion_tokens=72  total_tokens=98

	fmt.Println("\n\nBifrost SDK (OpenAI) Message... ")
	bifrostOpenAIMessage("gpt-4o-mini")
	// [usage] prompt_tokens=24  completion_tokens=50  total_tokens=74

	fmt.Println("\n\nOpenAI SDK + Bifrost + Anthropic... ")
	openAIWithBifrostAnthropicMessage("anthropic/claude-sonnet-4-6")
	// [usage]
	// prompt_tokens=26  completion_tokens=74  total_tokens=100
	// cached_tokens=0  prompt_audio_tokens=0

	fmt.Println("\n\nAnthropic SDK + Bifrost + OpenAI... ")
	anthropicWithBifrostOpenAIMessage("openai/gpt-4o-mini")
	// [usage] input_tokens=0  output_tokens=49  cache_creation=0  cache_read=0

	fmt.Println("\n\nDone")
}
