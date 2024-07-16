package services

import (
	"github.com/poligonoio/vega-core/internal/models"
)

type ActivityService interface {
	CreateActivity(activity models.Activity) error
}
