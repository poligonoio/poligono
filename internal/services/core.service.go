package services

import (
	"github.com/poligonoio/vega-core/internal/models"
)

type CoreService interface {
	PromptGemini(prompt string) (models.QueryResult, error)
}
