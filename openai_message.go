package main

import (
	"context"
	"fmt"
	"log"

	"github.com/openai/openai-go/v3"
	openaioption "github.com/openai/openai-go/v3/option"
)

func openAIMessage(model string) {
	client := openai.NewClient(
		// openaioption.WithAPIKey(os.Getenv("OPENAI_API_KEY")),
		openaioption.WithAPIKey("dummy-api-key"),
		openaioption.WithBaseURL("http://localhost:8080/openai/v1/"),
	)

	ctx := context.Background()
	stream := client.Chat.Completions.NewStreaming(ctx, openai.ChatCompletionNewParams{
		Model:     model,
		MaxTokens: openai.Int(1024),
		Messages: []openai.ChatCompletionMessageParamUnion{
			openai.UserMessage(initPrompt),
		},
	})

	for stream.Next() {
		chunk := stream.Current()
		if len(chunk.Choices) > 0 && chunk.Choices[0].Delta.Content != "" {
			fmt.Print(chunk.Choices[0].Delta.Content)
		}
		if chunk.Usage.PromptTokens != 0 || chunk.Usage.CompletionTokens != 0 {
			u := chunk.Usage
			pd := u.PromptTokensDetails
			cd := u.CompletionTokensDetails
			fmt.Printf("\n\n[usage]\n")
			fmt.Printf("  prompt_tokens=%d  completion_tokens=%d  total_tokens=%d\n",
				u.PromptTokens, u.CompletionTokens, u.TotalTokens)
			fmt.Printf("  cached_tokens=%d  prompt_audio_tokens=%d\n",
				pd.CachedTokens, pd.AudioTokens)
			fmt.Printf("  completion_details: reasoning=%d  audio=%d  accepted_prediction=%d  rejected_prediction=%d\n",
				cd.ReasoningTokens, cd.AudioTokens, cd.AcceptedPredictionTokens, cd.RejectedPredictionTokens)
		}
	}

	if err := stream.Err(); err != nil {
		log.Fatal(err)
	}
}
