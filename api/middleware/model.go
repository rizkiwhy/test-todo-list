package middleware

const (
	MissingAuthHeaderTitleMessage    = "Missing authorization header"
	MissingAuthHeaderErrorMessage    = "Authorization header is required"
	ErrInvalidSigningMethodMessage   = "invalid signing method"
	ErrInvalidAuthHeaderTitleMessage = "Invalid authorization header"
	ErrInvalidTokenMessage           = "Invalid token"
	ErrInvalidTokenClaimsMessage     = "Invalid token claims"
	ErrFetchJWTPayloadErrorMessage   = "Failed to fetch JWT payload from cache"
)
