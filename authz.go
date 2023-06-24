package gokeycloak

import (
	"context"
	"net/http"
	"net/url"

	"github.com/go-resty/resty/v2"
	"github.com/pkg/errors"
)

func (g *GoKeycloak) getRequestingParty(ctx context.Context, token, realm string, options RequestingPartyTokenOptions, res interface{}) (*resty.Response, error) {
	return g.GetRequestWithBearerAuth(ctx, token).
		SetFormData(options.FormData()).
		SetFormDataFromValues(url.Values{"permission": PStringSlice(options.Permissions)}).
		SetResult(&res).
		Post(g.getRealmURL(realm, g.Config.openIDConnect, "token"))
}

// GetRequestingPartyPermissions returns a requesting party permissions granted by the server
func (g *GoKeycloak) GetRequestingPartyPermissions(ctx context.Context, token, realm string, options RequestingPartyTokenOptions) (int, *[]RequestingPartyPermission, error) {
	const errMessage = "could not get requesting party token"

	var res []RequestingPartyPermission

	options.ResponseMode = StringP("permissions")

	resp, err := g.getRequestingParty(ctx, token, realm, options, &res)
	if err := checkForError(resp, err, errMessage); err != nil {
		return resp.StatusCode(), nil, err
	}
	return resp.StatusCode(), &res, nil
}

// GetRequestingPartyPermissionDecision returns a requesting party permission decision granted by the server
func (g *GoKeycloak) GetRequestingPartyPermissionDecision(ctx context.Context, token, realm string, options RequestingPartyTokenOptions) (int, *RequestingPartyPermissionDecision, error) {
	const errMessage = "could not get requesting party token"

	var res RequestingPartyPermissionDecision

	options.ResponseMode = StringP("decision")

	resp, err := g.getRequestingParty(ctx, token, realm, options, &res)
	if err := checkForError(resp, err, errMessage); err != nil {
		return resp.StatusCode(), nil, err
	}

	return resp.StatusCode(), &res, nil
}

// -----------
// Realm Roles
// -----------

// CreateRealmRole creates a role in a realm
func (g *GoKeycloak) CreateRealmRole(ctx context.Context, token string, realm string, role Role) (int, string, error) {
	const errMessage = "could not create realm role"

	resp, err := g.GetRequestWithBearerAuth(ctx, token).
		SetBody(role).
		Post(g.getAdminRealmURL(realm, "roles"))

	if err := checkForError(resp, err, errMessage); err != nil {
		return resp.StatusCode(), "", err
	}

	return resp.StatusCode(), getID(resp), nil
}

// GetRealmRole returns a role from a realm by role's name
func (g *GoKeycloak) GetRealmRole(ctx context.Context, token, realm, roleName string) (int, *Role, error) {
	const errMessage = "could not get realm role"

	var result Role

	resp, err := g.GetRequestWithBearerAuth(ctx, token).
		SetResult(&result).
		Get(g.getAdminRealmURL(realm, "roles", roleName))

	if err = checkForError(resp, err, errMessage); err != nil {
		return resp.StatusCode(), nil, err
	}

	return resp.StatusCode(), &result, nil
}

// GetRealmRoleByID returns a role from a realm by role's ID
func (g *GoKeycloak) GetRealmRoleByID(ctx context.Context, token, realm, roleID string) (int, *Role, error) {
	const errMessage = "could not get realm role"

	var result Role
	resp, err := g.GetRequestWithBearerAuth(ctx, token).
		SetResult(&result).
		Get(g.getAdminRealmURL(realm, "roles-by-id", roleID))

	if err := checkForError(resp, err, errMessage); err != nil {
		return resp.StatusCode(), nil, err
	}

	return resp.StatusCode(), &result, nil
}

// GetRealmRoles get all roles of the given realm.
func (g *GoKeycloak) GetRealmRoles(ctx context.Context, token, realm string, params GetRoleParams) (int, []*Role, error) {
	const errMessage = "could not get realm roles"

	var result []*Role
	queryParams, err := GetQueryParams(params)
	if err != nil {
		return http.StatusInternalServerError, nil, errors.Wrap(err, errMessage)
	}

	resp, err := g.GetRequestWithBearerAuth(ctx, token).
		SetResult(&result).
		SetQueryParams(queryParams).
		Get(g.getAdminRealmURL(realm, "roles"))

	if err := checkForError(resp, err, errMessage); err != nil {
		return resp.StatusCode(), nil, err
	}

	return resp.StatusCode(), result, nil
}

// GetRealmRolesByUserID returns all roles assigned to the given user
func (g *GoKeycloak) GetRealmRolesByUserID(ctx context.Context, token, realm, userID string) (int, []*Role, error) {
	const errMessage = "could not get realm roles by user id"

	var result []*Role
	resp, err := g.GetRequestWithBearerAuth(ctx, token).
		SetResult(&result).
		Get(g.getAdminRealmURL(realm, "users", userID, "role-mappings", "realm"))

	if err = checkForError(resp, err, errMessage); err != nil {
		return resp.StatusCode(), nil, err
	}

	return resp.StatusCode(), result, nil
}

// GetRealmRolesByGroupID returns all roles assigned to the given group
func (g *GoKeycloak) GetRealmRolesByGroupID(ctx context.Context, token, realm, groupID string) (int, []*Role, error) {
	const errMessage = "could not get realm roles by group id"

	var result []*Role
	resp, err := g.GetRequestWithBearerAuth(ctx, token).
		SetResult(&result).
		Get(g.getAdminRealmURL(realm, "groups", groupID, "role-mappings", "realm"))

	if err = checkForError(resp, err, errMessage); err != nil {
		return resp.StatusCode(), nil, err
	}

	return resp.StatusCode(), result, nil
}

// UpdateRealmRole updates a role in a realm
func (g *GoKeycloak) UpdateRealmRole(ctx context.Context, token, realm, roleName string, role Role) (int, error) {
	const errMessage = "could not update realm role"

	resp, err := g.GetRequestWithBearerAuth(ctx, token).
		SetBody(role).
		Put(g.getAdminRealmURL(realm, "roles", roleName))

	return resp.StatusCode(), checkForError(resp, err, errMessage)
}

// UpdateRealmRoleByID updates a role in a realm by role's ID
func (g *GoKeycloak) UpdateRealmRoleByID(ctx context.Context, token, realm, roleID string, role Role) (int, error) {
	const errMessage = "could not update realm role"

	resp, err := g.GetRequestWithBearerAuth(ctx, token).
		SetBody(role).
		Put(g.getAdminRealmURL(realm, "roles-by-id", roleID))

	return resp.StatusCode(), checkForError(resp, err, errMessage)
}

// DeleteRealmRole deletes a role in a realm by role's name
func (g *GoKeycloak) DeleteRealmRole(ctx context.Context, token, realm, roleName string) (int, error) {
	const errMessage = "could not delete realm role"

	resp, err := g.GetRequestWithBearerAuth(ctx, token).
		Delete(g.getAdminRealmURL(realm, "roles", roleName))

	return resp.StatusCode(), checkForError(resp, err, errMessage)
}

// AddRealmRoleToUser adds realm-level role mappings
func (g *GoKeycloak) AddRealmRoleToUser(ctx context.Context, token, realm, userID string, roles []Role) (int, error) {
	const errMessage = "could not add realm role to user"

	resp, err := g.GetRequestWithBearerAuth(ctx, token).
		SetBody(roles).
		Post(g.getAdminRealmURL(realm, "users", userID, "role-mappings", "realm"))

	return resp.StatusCode(), checkForError(resp, err, errMessage)
}

// DeleteRealmRoleFromUser deletes realm-level role mappings
func (g *GoKeycloak) DeleteRealmRoleFromUser(ctx context.Context, token, realm, userID string, roles []Role) (int, error) {
	const errMessage = "could not delete realm role from user"

	resp, err := g.GetRequestWithBearerAuth(ctx, token).
		SetBody(roles).
		Delete(g.getAdminRealmURL(realm, "users", userID, "role-mappings", "realm"))

	return resp.StatusCode(), checkForError(resp, err, errMessage)
}

// AddRealmRoleToGroup adds realm-level role mappings
func (g *GoKeycloak) AddRealmRoleToGroup(ctx context.Context, token, realm, groupID string, roles []Role) (int, error) {
	const errMessage = "could not add realm role to group"

	resp, err := g.GetRequestWithBearerAuth(ctx, token).
		SetBody(roles).
		Post(g.getAdminRealmURL(realm, "groups", groupID, "role-mappings", "realm"))

	return resp.StatusCode(), checkForError(resp, err, errMessage)
}

// DeleteRealmRoleFromGroup deletes realm-level role mappings
func (g *GoKeycloak) DeleteRealmRoleFromGroup(ctx context.Context, token, realm, groupID string, roles []Role) (int, error) {
	const errMessage = "could not delete realm role from group"

	resp, err := g.GetRequestWithBearerAuth(ctx, token).
		SetBody(roles).
		Delete(g.getAdminRealmURL(realm, "groups", groupID, "role-mappings", "realm"))

	return resp.StatusCode(), checkForError(resp, err, errMessage)
}

// AddRealmRoleComposite adds a role to the composite.
func (g *GoKeycloak) AddRealmRoleComposite(ctx context.Context, token, realm, roleName string, roles []Role) (int, error) {
	const errMessage = "could not add realm role composite"

	resp, err := g.GetRequestWithBearerAuth(ctx, token).
		SetBody(roles).
		Post(g.getAdminRealmURL(realm, "roles", roleName, "composites"))

	return resp.StatusCode(), checkForError(resp, err, errMessage)
}

// DeleteRealmRoleComposite deletes a role from the composite.
func (g *GoKeycloak) DeleteRealmRoleComposite(ctx context.Context, token, realm, roleName string, roles []Role) (int, error) {
	const errMessage = "could not delete realm role composite"

	resp, err := g.GetRequestWithBearerAuth(ctx, token).
		SetBody(roles).
		Delete(g.getAdminRealmURL(realm, "roles", roleName, "composites"))

	return resp.StatusCode(), checkForError(resp, err, errMessage)
}

// GetCompositeRealmRoles returns all realm composite roles associated with the given realm role
func (g *GoKeycloak) GetCompositeRealmRoles(ctx context.Context, token, realm, roleName string) (int, []*Role, error) {
	const errMessage = "could not get composite realm roles by role"

	var result []*Role
	resp, err := g.GetRequestWithBearerAuth(ctx, token).
		SetResult(&result).
		Get(g.getAdminRealmURL(realm, "roles", roleName, "composites"))

	if err = checkForError(resp, err, errMessage); err != nil {
		return resp.StatusCode(), nil, err
	}

	return resp.StatusCode(), result, nil
}

// GetCompositeRolesByRoleID returns all realm composite roles associated with the given client role
func (g *GoKeycloak) GetCompositeRolesByRoleID(ctx context.Context, token, realm, roleID string) (int, []*Role, error) {
	const errMessage = "could not get composite client roles by role id"

	var result []*Role
	resp, err := g.GetRequestWithBearerAuth(ctx, token).
		SetResult(&result).
		Get(g.getAdminRealmURL(realm, "roles-by-id", roleID, "composites"))

	if err = checkForError(resp, err, errMessage); err != nil {
		return resp.StatusCode(), nil, err
	}

	return resp.StatusCode(), result, nil
}

// GetCompositeRealmRolesByRoleID returns all realm composite roles associated with the given client role
func (g *GoKeycloak) GetCompositeRealmRolesByRoleID(ctx context.Context, token, realm, roleID string) (int, []*Role, error) {
	const errMessage = "could not get composite client roles by role id"

	var result []*Role
	resp, err := g.GetRequestWithBearerAuth(ctx, token).
		SetResult(&result).
		Get(g.getAdminRealmURL(realm, "roles-by-id", roleID, "composites", "realm"))

	if err = checkForError(resp, err, errMessage); err != nil {
		return resp.StatusCode(), nil, err
	}

	return resp.StatusCode(), result, nil
}

// GetCompositeRealmRolesByUserID returns all realm roles and composite roles assigned to the given user
func (g *GoKeycloak) GetCompositeRealmRolesByUserID(ctx context.Context, token, realm, userID string) (int, []*Role, error) {
	const errMessage = "could not get composite client roles by user id"

	var result []*Role
	resp, err := g.GetRequestWithBearerAuth(ctx, token).
		SetResult(&result).
		Get(g.getAdminRealmURL(realm, "users", userID, "role-mappings", "realm", "composite"))

	if err = checkForError(resp, err, errMessage); err != nil {
		return resp.StatusCode(), nil, err
	}

	return resp.StatusCode(), result, nil
}

// GetCompositeRealmRolesByGroupID returns all realm roles and composite roles assigned to the given group
func (g *GoKeycloak) GetCompositeRealmRolesByGroupID(ctx context.Context, token, realm, groupID string) (int, []*Role, error) {
	const errMessage = "could not get composite client roles by user id"

	var result []*Role
	resp, err := g.GetRequestWithBearerAuth(ctx, token).
		SetResult(&result).
		Get(g.getAdminRealmURL(realm, "groups", groupID, "role-mappings", "realm", "composite"))

	if err = checkForError(resp, err, errMessage); err != nil {
		return resp.StatusCode(), nil, err
	}

	return resp.StatusCode(), result, nil
}

// GetAvailableRealmRolesByUserID returns all available realm roles to the given user
func (g *GoKeycloak) GetAvailableRealmRolesByUserID(ctx context.Context, token, realm, userID string) (int, []*Role, error) {
	const errMessage = "could not get available client roles by user id"

	var result []*Role
	resp, err := g.GetRequestWithBearerAuth(ctx, token).
		SetResult(&result).
		Get(g.getAdminRealmURL(realm, "users", userID, "role-mappings", "realm", "available"))

	if err = checkForError(resp, err, errMessage); err != nil {
		return resp.StatusCode(), nil, err
	}

	return resp.StatusCode(), result, nil
}

// GetAvailableRealmRolesByGroupID returns all available realm roles to the given group
func (g *GoKeycloak) GetAvailableRealmRolesByGroupID(ctx context.Context, token, realm, groupID string) (int, []*Role, error) {
	const errMessage = "could not get available client roles by user id"

	var result []*Role
	resp, err := g.GetRequestWithBearerAuth(ctx, token).
		SetResult(&result).
		Get(g.getAdminRealmURL(realm, "groups", groupID, "role-mappings", "realm", "available"))

	if err = checkForError(resp, err, errMessage); err != nil {
		return resp.StatusCode(), nil, err
	}

	return resp.StatusCode(), result, nil
}

func (g *GoKeycloak) EvaluatePermission(ctx context.Context, userToken, realm, audience, response_mode string, permissions []string) (int, *JWT, error) {
	var permission_token_grant string = "urn:ietf:params:oauth:grant-type:uma-ticket"
	var options RequestingPartyTokenOptions = RequestingPartyTokenOptions{
  	GrantType: &permission_token_grant,
		Audience: &audience,
		ResponseMode: &response_mode,
		Permissions: &permissions,		
	}
	
	statusCode, jwt, err := g.GetRequestingPartyToken(ctx, userToken, realm, options)
	if err != nil {
		return statusCode, nil, err
	}

	return statusCode, jwt, nil
}
