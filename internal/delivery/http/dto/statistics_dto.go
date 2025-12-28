package dto

import "hub/internal/delivery/http/entity"

type StatisticCategory struct {
	Key         string          `json:"key"`
	Description string          `json:"description"`
	Icon        string          `json:"icon"`
	Tracks      []*entity.Track `json:"tracks"`
}

type StatisticsResponse struct {
	Statistics []*StatisticCategory `json:"statistics"`
}
