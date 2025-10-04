package spaserver

import (
	"net/http"
	"spahttp/domain"
)

// IsAuthenticated checks
// if a user is authenticated by verifying the presence of session or authentication data in the request.
func (s *Server) IsAuthenticated(r *http.Request) bool {
	isAuthenticated, ok := r.Context().Value(domain.IsAuthenticatedContextKey).(bool)
	if !ok {
		return false
	}
	return isAuthenticated
}

// GetAuthenticatedUserId retrieves the authenticated user ID from the session data using the HTTP request context.
func (s *Server) GetAuthenticatedUserId(r *http.Request) int {
	userId := s.sessionManager.GetInt(r.Context(), "authenticatedUserID")
	return userId
}

// GetAuthenticatedUserName retrieves the authenticated user's name from the session data using the HTTP request context.
func (s *Server) GetAuthenticatedUserName(r *http.Request) string {
	userName := s.sessionManager.GetString(r.Context(), "authenticatedUsername")
	return userName
}
