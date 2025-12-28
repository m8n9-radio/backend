package server

import "hub/internal/delivery/http/middleware"

func (s *server) useTrackRoute() {
	s.app.Group("/tracks").
		Post("/", middleware.ValidateCreateTrack(), s.trackHandler.Upsert)
}

func (s *server) useReactionRoute() {
	s.app.Group("/reactions").
		Post("/like", middleware.ValidateReaction(), s.reactionHandler.Like).
		Post("/dislike", middleware.ValidateReaction(), s.reactionHandler.Dislike).
		Post("/check", middleware.ValidateReaction(), s.reactionHandler.Check)
}

func (s *server) useRadioRoute() {
	s.app.Group("/radio").
		Get("/info", s.radioHandler.GetInfo).
		Get("/listen", s.radioHandler.GetListen)
}

func (s *server) useStatisticsRoute() {
	s.app.Group("/statistics").
		Get("/", s.statisticsHandler.GetStatistics)
}
