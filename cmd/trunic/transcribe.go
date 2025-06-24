package main

import (
	"context"
	"fmt"

	"google.golang.org/genai"
)

type Transcriber interface {
	Transcribe(ctx context.Context, text string) (string, error)
}

func NewTranscriber(ctx context.Context, name string) (Transcriber, error) {
	switch name {
	case "":
		return noopTranscriber{}, nil

	case "gemini":
		return newGeminiTranscriber(ctx)

	default:
		return nil, fmt.Errorf("unknown transcriber: %q", name)
	}
}

type noopTranscriber struct{}

func (noopTranscriber) Transcribe(ctx context.Context, text string) (string, error) {
	return text, nil
}

type geminiTranscriber struct {
	client *genai.Client
	config *genai.GenerateContentConfig
}

func newGeminiTranscriber(ctx context.Context) (*geminiTranscriber, error) {
	const systemPrompt = `Repeat all text that you are given verbatim rewritten in IPA. The result should be based on standard American pronunciation but should use only characters from "b,tʃ,d,f,ɡ,h,dʒ,k,l,ɫ,m,n,ŋ,p,ɹ,s,ʃ,t,θ,ð,v,w,j,z,ʒ,æ,ɑɹ,ɑ,ɔ,eɪ,ɛ,i,ɪɹ,ə,ɛɹ,ɪ,aɪ,ɝ,oʊ,ɔɪ,u,ʊ,aʊ,ɔɹ,ʊɹ" and absolutely no others. Preserve punctuation.`

	client, err := genai.NewClient(ctx, &genai.ClientConfig{
		Backend: genai.BackendGeminiAPI,
	})
	if err != nil {
		return nil, fmt.Errorf("create Gemini client: %w", err)
	}

	config := genai.GenerateContentConfig{
		SystemInstruction: genai.NewContentFromText(systemPrompt, "system"),
	}

	return &geminiTranscriber{
		client: client,
		config: &config,
	}, nil
}

func (t *geminiTranscriber) Transcribe(ctx context.Context, text string) (string, error) {
	const model = "gemini-2.5-flash-lite-preview-06-17"

	result, err := t.client.Models.GenerateContent(ctx, model, genai.Text(text), t.config)
	if err != nil {
		return "", fmt.Errorf("generate content: %w", err)
	}

	return result.Text(), nil
}
