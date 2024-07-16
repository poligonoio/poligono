package services

import (
	"github.com/poligonoio/vega-core/internal/models"
)

type CoreService interface {
	GeminiPrompt(prompt string) (models.GeminiQueryResult, error)
}
