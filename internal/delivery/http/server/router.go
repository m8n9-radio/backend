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
