package server

import "hub/internal/delivery/http/middleware"

func (s *server) useTrackRoute() {
	s.app.
		Group("/tracks").
		Post("/",
			middleware.ValidateCreateTrack(),
			s.trackHandler.Upsert,
		)
}

func (s *server) useReactionRoute() {
	reactions := s.app.Group("/reactions")

	reactions.Post("/like",
		middleware.ValidateReaction(),
		s.reactionHandler.Like,
	)

	reactions.Post("/dislike",
		middleware.ValidateReaction(),
		s.reactionHandler.Dislike,
	)

	reactions.Post("/check",
		middleware.ValidateReaction(),
		s.reactionHandler.Check,
	)
}

func (s *server) useRadioRoute() {
	s.app.Get("/radio/info", s.radioHandler.GetInfo)
}
