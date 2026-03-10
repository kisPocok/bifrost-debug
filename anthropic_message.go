package main

import (
	"context"
	"fmt"
	"log"

	"github.com/anthropics/anthropic-sdk-go"
	"github.com/anthropics/anthropic-sdk-go/option"
)

func anthropicMessage(model string) {
	client := anthropic.NewClient(
		option.WithAPIKey(ANTHROPIC_API_KEY),
		option.WithBaseURL("http://localhost:8080/anthropic"),
	)

	ctx := context.Background()
	stream := client.Messages.NewStreaming(ctx, anthropic.MessageNewParams{
		Model:     anthropic.Model(model),
		MaxTokens: 1024,
		Messages: []anthropic.MessageParam{
			anthropic.NewUserMessage(
				anthropic.NewTextBlock(initPrompt),
			),
		},
	})

	message := anthropic.Message{}
	for stream.Next() {
		event := stream.Current()
		message.Accumulate(event)

		switch delta := event.AsAny().(type) {
		case anthropic.ContentBlockDeltaEvent:
			switch d := delta.Delta.AsAny().(type) {
			case anthropic.TextDelta:
				fmt.Print(d.Text)
			}
		}
	}
	fmt.Printf("\n\n[usage] input_tokens=%d  output_tokens=%d  cache_creation=%d  cache_read=%d\n",
		message.Usage.InputTokens,
		message.Usage.OutputTokens,
		message.Usage.CacheCreationInputTokens,
		message.Usage.CacheReadInputTokens,
	)

	if err := stream.Err(); err != nil {
		log.Fatal(err)
	}
}
