package main

import (
	"context"
	"fmt"
	"log"

	"github.com/maximhq/bifrost/core/schemas"

	bifrost "github.com/maximhq/bifrost/core"
)

// bifrostOpenAIAccount implements schemas.Account for Bifrost SDK (OpenAI via gateway).
type bifrostOpenAIAccount struct{}

func (a *bifrostOpenAIAccount) GetConfiguredProviders() ([]schemas.ModelProvider, error) {
	return []schemas.ModelProvider{schemas.OpenAI}, nil
}

func (a *bifrostOpenAIAccount) GetKeysForProvider(ctx context.Context, provider schemas.ModelProvider) ([]schemas.Key, error) {
	if provider == schemas.OpenAI {
		return []schemas.Key{{
			Value:  schemas.EnvVar{Val: OPENAI_API_KEY},
			Models: []string{},
			Weight: 1.0,
		}}, nil
	}
	return nil, fmt.Errorf("provider %s not supported", provider)
}

func (a *bifrostOpenAIAccount) GetConfigForProvider(provider schemas.ModelProvider) (*schemas.ProviderConfig, error) {
	if provider == schemas.OpenAI {
		return &schemas.ProviderConfig{
			NetworkConfig: schemas.NetworkConfig{
				BaseURL: "http://localhost:8080/openai",
			},
			ConcurrencyAndBufferSize: schemas.DefaultConcurrencyAndBufferSize,
		}, nil
	}
	return nil, fmt.Errorf("provider %s not supported", provider)
}

func bifrostOpenAIMessage(model string) {
	client, initErr := bifrost.Init(context.Background(), schemas.BifrostConfig{
		Account: &bifrostOpenAIAccount{},
	})
	if initErr != nil {
		log.Fatal(initErr)
	}
	defer client.Shutdown()

	messages := []schemas.ChatMessage{
		{
			Role: schemas.ChatMessageRoleUser,
			Content: &schemas.ChatMessageContent{
				ContentStr: schemas.Ptr(initPrompt),
			},
		},
	}

	stream, bErr := client.ChatCompletionStreamRequest(
		schemas.NewBifrostContext(context.Background(), schemas.NoDeadline),
		&schemas.BifrostChatRequest{
			Provider: schemas.OpenAI,
			Model:    model,
			Input:    messages,
		},
	)
	if bErr != nil {
		log.Fatal(bErr)
	}

	for chunk := range stream {
		if chunk.BifrostError != nil {
			msg := "stream error"
			if chunk.BifrostError.Error != nil && chunk.BifrostError.Error.Message != "" {
				msg = chunk.BifrostError.Error.Message
			}
			extra := chunk.BifrostError.ExtraFields
			if extra.Provider != "" || extra.ModelRequested != "" {
				msg = fmt.Sprintf("%s (provider=%s model=%s)", msg, extra.Provider, extra.ModelRequested)
			}
			if chunk.BifrostError.StatusCode != nil {
				msg = fmt.Sprintf("%s [status=%d]", msg, *chunk.BifrostError.StatusCode)
			}
			log.Fatalf("Bifrost SDK stream error: %s", msg)
		}
		if chunk.BifrostChatResponse != nil && len(chunk.BifrostChatResponse.Choices) > 0 {
			choice := chunk.BifrostChatResponse.Choices[0]
			if choice.ChatStreamResponseChoice != nil &&
				choice.ChatStreamResponseChoice.Delta != nil &&
				choice.ChatStreamResponseChoice.Delta.Content != nil {
				fmt.Print(*choice.ChatStreamResponseChoice.Delta.Content)
			}
			if chunk.BifrostChatResponse.Usage != nil {
				u := chunk.BifrostChatResponse.Usage
				fmt.Printf("\n\n[usage] prompt_tokens=%d  completion_tokens=%d  total_tokens=%d\n",
					u.PromptTokens, u.CompletionTokens, u.TotalTokens)
			}
		}
	}
}
