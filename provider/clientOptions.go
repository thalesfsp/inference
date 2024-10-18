package provider

//////
// Vars, consts, and types.
//////

// ClientFunc allows to set provider options.
type ClientFunc func(o *ClientOptions) error

// ClientOptions for the provider.
type ClientOptions struct {
	// Endpoint to reach the provider.
	Endpoint string `json:"endpoint" validate:"required"`

	// Token to authenticate against the provider.
	Token string `json:"-"`
}

//////
// Exported built-in options.
//////

// WithEndpoint sets the provider endpoint.
func WithEndpoint(endpoint string) ClientFunc {
	return func(o *ClientOptions) error {
		if endpoint != "" {
			o.Endpoint = endpoint
		}

		return nil
	}
}

// WithToken sets the provider token.
func WithToken(token string) ClientFunc {
	return func(o *ClientOptions) error {
		if token != "" {
			o.Token = token
		}

		return nil
	}
}
