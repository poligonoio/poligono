package services

import (
	"context"
	"encoding/json"

	"github.com/google/generative-ai-go/genai"
	"github.com/poligonoio/vega-core/pkg/logger"
	"github.com/poligonoio/vega-core/pkg/utils"
)

type Content struct {
	Parts []string `json:"Parts"`
	Role  string   `json:"Role"`
}
type Candidates struct {
	Content *Content `json:"Content"`
}
type ContentResponse struct {
	Candidates *[]Candidates `json:"Candidates"`
}

type CoreServiceImpl struct {
	ctx   context.Context
	model *genai.GenerativeModel
}

func NewCoreService(ctx context.Context, model *genai.GenerativeModel) CoreService {
	return &CoreServiceImpl{
		ctx:   ctx,
		model: model,
	}
}

func (self *CoreServiceImpl) PromptGemini(prompt string) (string, error) {
	resp, err := self.model.GenerateContent(self.ctx, genai.Text(prompt))
	if err != nil {
		logger.Error.Fatalf("Failed to generate content using Gemini: %v\n", err)
		return "", err
	}

	marshalResponse, _ := json.MarshalIndent(resp, "", "  ")

	var generateResponse ContentResponse
	if err := json.Unmarshal(marshalResponse, &generateResponse); err != nil {
		logger.Error.Fatalf("Failed to Unmarshal json from Gemini  %v\n", err)
		return "", err
	}

	// We can make multiple request to the users data source and return the one that seems more accurate
	for _, cad := range *generateResponse.Candidates {
		if cad.Content != nil {
			for _, part := range cad.Content.Parts {
				return utils.SanitizeQuery(part), nil
			}
		}
	}

	// improve this
	return "", nil
}
