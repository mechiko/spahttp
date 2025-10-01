package domain

// contextKey is a custom type used for defining strongly typed keys in context to avoid collision and improve clarity.
type contextKey string

// IsAuthenticatedContextKey is a context key used to store and retrieve authentication status in a strongly typed manner.
const IsAuthenticatedContextKey = contextKey("isAuthenticated")
