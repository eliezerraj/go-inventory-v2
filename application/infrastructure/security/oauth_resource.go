package security

type OAuthResource struct {
	Resource               string   `json:"resource"`
	AuthorizationServers   []string `json:"authorization_servers"`
	ScopesSupported        []string `json:"scopes_supported"`
	RolesSupported         []string `json:"roles_supported"`
	BearerMethodsSupported []string `json:"bearer_methods_supported"`
}

func NewOAuthResource(resourceMetadata, authorizationServer string) *OAuthResource {
	return &OAuthResource{
		Resource:             resourceMetadata,
		AuthorizationServers: []string{authorizationServer},
		ScopesSupported: []string{
			"inventory:read",
			"inventory:create",
			"inventory:update",
			"product:update",
			"inventory:patch",
			"price:patch",
		},
		RolesSupported: []string{
			"inventory-admin-role",
			"inventory-all-role",
		},
		BearerMethodsSupported: []string{
			"header",
		},
	}
}