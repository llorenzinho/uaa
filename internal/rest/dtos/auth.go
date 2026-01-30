package dtos

type AuthorizationCodeQueryParams struct {
	ResponseType string `form:"response_type" binding:"required"`
	CLientId     string `form:"client_id" binding:"required"`
	RedirectUri  string `form:"redirect_uri" binding:"required"`
	Scope        string `form:"scope"`
	State        string `form:"state" binding:"required"` // CSRF token
}
