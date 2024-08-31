package services

type CoreService interface {
	PromptGemini(prompt string) (string, error)
}
