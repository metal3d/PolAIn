package api

import (
	"fmt"
	"testing"
)

func TestCall(t *testing.T) {
	// Test the Ask function
	prompt := "Write a long poem about AI in french"
	stream, history := Ask(prompt, nil, "openai")

	for chunk := range stream {
		if chunk.Choices[0].Delta.Content == "" {
			t.Error("Received empty chunk")
		}
		fmt.Printf("%s", chunk.Choices[0].Delta.Content)
	}
	fmt.Println("\nFinished receiving chunks")
	// ensure that the prompt is in the history
	if len(history) == 0 {
		t.Error("History is empty")
	}
	if history[0].Role != System {
		t.Error("First message in history is not the system message")
	}
	if history[1].Content != prompt {
		t.Errorf("First message in history is not the prompt: %s", history[0].Content)
	}
}

func TestIds(t *testing.T) {
	// ask 2 questions, the first chunks should have the same id
	prompt1 := "Write a long poem about AI in french"
	prompt2 := "Write a long poem about AI in english"
	stream1, _ := Ask(prompt1, nil, "openai")
	id1 := ""
	for chunk1 := range stream1 {
		if id1 == "" {
			id1 = chunk1.Id
			t.Log("Id1: ", id1)
		}
		if chunk1.Id != id1 {
			t.Errorf("First message in history is not the prompt: %s", chunk1.Id)
		}
	}
	id2 := ""
	stream2, _ := Ask(prompt2, nil, "openai")
	for chunk2 := range stream2 {
		if id2 == "" {
			id2 = chunk2.Id
			t.Log("Id2: ", id2)
		}
		if chunk2.Id != id2 {
			t.Errorf("First message in history is not the prompt: %s", chunk2.Id)
		}
	}
	if id1 == id2 {
		t.Errorf("Ids are the same: %s", id1)
	}
}
