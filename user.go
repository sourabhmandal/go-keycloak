package gokeycloak

import (
	"context"
	"net/http"

	"github.com/pkg/errors"
)

// -----
// Users
// -----

// CreateUser creates the given user in the given realm and returns it's userID
// Note: Keycloak has not documented what members of the User object are actually being accepted, when creating a user.
// Things like RealmRoles must be attached using followup calls to the respective functions.
func (g *GoKeycloak) CreateUser(ctx context.Context, token, realm string, user User) (int, string, error) {
	const errMessage = "could not create user"

	resp, err := g.GetRequestWithBearerAuth(ctx, token).
		SetBody(user).
		Post(g.getAdminRealmURL(realm, "users"))

	if err := checkForError(resp, err, errMessage); err != nil {
		return resp.StatusCode(), "", err
	}

	return resp.StatusCode(), getID(resp), nil
}

// DeleteUser delete a given user
func (g *GoKeycloak) DeleteUser(ctx context.Context, token, realm, userID string) (int, error) {
	const errMessage = "could not delete user"

	resp, err := g.GetRequestWithBearerAuth(ctx, token).
		Delete(g.getAdminRealmURL(realm, "users", userID))

	return resp.StatusCode(), checkForError(resp, err, errMessage)
}

// GetUserByID fetches a user from the given realm with the given userID
func (g *GoKeycloak) GetUserByID(ctx context.Context, accessToken, realm, userID string) (int, *User, error) {
	const errMessage = "could not get user by id"

	if userID == "" {
		return http.StatusBadRequest, nil, errors.Wrap(errors.New("userID shall not be empty"), errMessage)
	}

	var result User
	resp, err := g.GetRequestWithBearerAuth(ctx, accessToken).
		SetResult(&result).
		Get(g.getAdminRealmURL(realm, "users", userID))

	if err := checkForError(resp, err, errMessage); err != nil {
		return resp.StatusCode(), nil, err
	}

	return resp.StatusCode(), &result, nil
}

// GetUserCount gets the user count in the realm
func (g *GoKeycloak) GetUserCount(ctx context.Context, token string, realm string, params GetUsersParams) (int, int, error) {
	const errMessage = "could not get user count"

	var result int
	queryParams, err := GetQueryParams(params)
	if err != nil {
		return http.StatusBadRequest, 0, errors.Wrap(err, errMessage)
	}

	resp, err := g.GetRequestWithBearerAuth(ctx, token).
		SetResult(&result).
		SetQueryParams(queryParams).
		Get(g.getAdminRealmURL(realm, "users", "count"))

	if err := checkForError(resp, err, errMessage); err != nil {
		return resp.StatusCode(), -1, errors.Wrap(err, errMessage)
	}

	return resp.StatusCode(), result, nil
}

// GetUserGroups get all groups for user
func (g *GoKeycloak) GetUserGroups(ctx context.Context, token, realm, userID string, params GetGroupsParams) (int, []*Group, error) {
	const errMessage = "could not get user groups"

	var result []*Group
	queryParams, err := GetQueryParams(params)
	if err != nil {
		return http.StatusBadRequest, nil, errors.Wrap(err, errMessage)
	}

	resp, err := g.GetRequestWithBearerAuth(ctx, token).
		SetResult(&result).
		SetQueryParams(queryParams).
		Get(g.getAdminRealmURL(realm, "users", userID, "groups"))

	if err := checkForError(resp, err, errMessage); err != nil {
		return resp.StatusCode(), nil, err
	}

	return resp.StatusCode(), result, nil
}

// GetUsers get all users in realm
func (g *GoKeycloak) GetUsers(ctx context.Context, token, realm string, params GetUsersParams) (int, []*User, error) {
	const errMessage = "could not get users"

	var result []*User
	queryParams, err := GetQueryParams(params)
	if err != nil {
		return http.StatusBadRequest, nil, errors.Wrap(err, errMessage)
	}

	resp, err := g.GetRequestWithBearerAuth(ctx, token).
		SetResult(&result).
		SetQueryParams(queryParams).
		Get(g.getAdminRealmURL(realm, "users"))

	if err := checkForError(resp, err, errMessage); err != nil {
		return resp.StatusCode(), nil, err
	}

	return resp.StatusCode(), result, nil
}

// GetUsersByRoleName returns all users have a given role
func (g *GoKeycloak) GetUsersByRoleName(ctx context.Context, token, realm, roleName string, params GetUsersByRoleParams) (int, []*User, error) {
	const errMessage = "could not get users by role name"

	var result []*User
	queryParams, err := GetQueryParams(params)
	if err != nil {
		return http.StatusBadRequest, nil, errors.Wrap(err, errMessage)
	}

	resp, err := g.GetRequestWithBearerAuth(ctx, token).
		SetResult(&result).
		SetQueryParams(queryParams).
		Get(g.getAdminRealmURL(realm, "roles", roleName, "users"))

	if err := checkForError(resp, err, errMessage); err != nil {
		return resp.StatusCode(), nil, err
	}

	return resp.StatusCode(), result, nil
}

// GetUsersByClientRoleName returns all users have a given client role
func (g *GoKeycloak) GetUsersByClientRoleName(ctx context.Context, token, realm, idOfClient, roleName string, params GetUsersByRoleParams) (int, []*User, error) {
	const errMessage = "could not get users by client role name"

	var result []*User
	queryParams, err := GetQueryParams(params)
	if err != nil {
		return http.StatusBadRequest, nil, err
	}

	resp, err := g.GetRequestWithBearerAuth(ctx, token).
		SetResult(&result).
		SetQueryParams(queryParams).
		Get(g.getAdminRealmURL(realm, "clients", idOfClient, "roles", roleName, "users"))

	if err := checkForError(resp, err, errMessage); err != nil {
		return resp.StatusCode(), nil, err
	}

	return resp.StatusCode(), result, nil
}

// SetPassword sets a new password for the user with the given id. Needs elevated privileges
func (g *GoKeycloak) SetPassword(ctx context.Context, token, userID, realm, password string, temporary bool) (int, error) {
	const errMessage = "could not set password"

	requestBody := SetPasswordRequest{Password: &password, Temporary: &temporary, Type: StringP("password")}
	resp, err := g.GetRequestWithBearerAuth(ctx, token).
		SetBody(requestBody).
		Put(g.getAdminRealmURL(realm, "users", userID, "reset-password"))

	return resp.StatusCode(), checkForError(resp, err, errMessage)
}

// UpdateUser updates a given user
func (g *GoKeycloak) UpdateUser(ctx context.Context, token, realm string, user User) (int, error) {
	const errMessage = "could not update user"

	resp, err := g.GetRequestWithBearerAuth(ctx, token).
		SetBody(user).
		Put(g.getAdminRealmURL(realm, "users", PString(user.ID)))

	return resp.StatusCode(), checkForError(resp, err, errMessage)
}

// AddUserToGroup puts given user to given group
func (g *GoKeycloak) AddUserToGroup(ctx context.Context, token, realm, userID, groupID string) (int, error) {
	const errMessage = "could not add user to group"

	resp, err := g.GetRequestWithBearerAuth(ctx, token).
		Put(g.getAdminRealmURL(realm, "users", userID, "groups", groupID))

	return resp.StatusCode(), checkForError(resp, err, errMessage)
}

// DeleteUserFromGroup deletes given user from given group
func (g *GoKeycloak) DeleteUserFromGroup(ctx context.Context, token, realm, userID, groupID string) (int, error) {
	const errMessage = "could not delete user from group"

	resp, err := g.GetRequestWithBearerAuth(ctx, token).
		Delete(g.getAdminRealmURL(realm, "users", userID, "groups", groupID))

	return resp.StatusCode(), checkForError(resp, err, errMessage)
}

// GetUserSessions returns user sessions associated with the user
func (g *GoKeycloak) GetUserSessions(ctx context.Context, token, realm, userID string) (int, []*UserSessionRepresentation, error) {
	const errMessage = "could not get user sessions"

	var res []*UserSessionRepresentation
	resp, err := g.GetRequestWithBearerAuth(ctx, token).
		SetResult(&res).
		Get(g.getAdminRealmURL(realm, "users", userID, "sessions"))

	if err := checkForError(resp, err, errMessage); err != nil {
		return resp.StatusCode(), nil, err
	}

	return resp.StatusCode(), res, nil
}

// GetUserOfflineSessionsForClient returns offline sessions associated with the user and client
func (g *GoKeycloak) GetUserOfflineSessionsForClient(ctx context.Context, token, realm, userID, idOfClient string) (int, []*UserSessionRepresentation, error) {
	const errMessage = "could not get user offline sessions for client"

	var res []*UserSessionRepresentation
	resp, err := g.GetRequestWithBearerAuth(ctx, token).
		SetResult(&res).
		Get(g.getAdminRealmURL(realm, "users", userID, "offline-sessions", idOfClient))

	if err := checkForError(resp, err, errMessage); err != nil {
		return resp.StatusCode(), nil, err
	}

	return resp.StatusCode(), res, nil
}

// AddClientRolesToUser adds client-level role mappings
func (g *GoKeycloak) AddClientRolesToUser(ctx context.Context, token, realm, idOfClient, userID string, roles []Role) (int, error) {
	const errMessage = "could not add client role to user"

	resp, err := g.GetRequestWithBearerAuth(ctx, token).
		SetBody(roles).
		Post(g.getAdminRealmURL(realm, "users", userID, "role-mappings", "clients", idOfClient))

	return resp.StatusCode(), checkForError(resp, err, errMessage)
}

// AddClientRoleToUser adds client-level role mappings
//
// Deprecated: replaced by AddClientRolesToUser
func (g *GoKeycloak) AddClientRoleToUser(ctx context.Context, token, realm, idOfClient, userID string, roles []Role) (int, error) {
	return g.AddClientRolesToUser(ctx, token, realm, idOfClient, userID, roles)
}

// DeleteClientRolesFromUser adds client-level role mappings
func (g *GoKeycloak) DeleteClientRolesFromUser(ctx context.Context, token, realm, idOfClient, userID string, roles []Role) (int, error) {
	const errMessage = "could not delete client role from user"

	resp, err := g.GetRequestWithBearerAuth(ctx, token).
		SetBody(roles).
		Delete(g.getAdminRealmURL(realm, "users", userID, "role-mappings", "clients", idOfClient))

	return resp.StatusCode(), checkForError(resp, err, errMessage)
}

// DeleteClientRoleFromUser adds client-level role mappings
//
// Deprecated: replaced by DeleteClientRolesFrom
func (g *GoKeycloak) DeleteClientRoleFromUser(ctx context.Context, token, realm, idOfClient, userID string, roles []Role) (int, error) {
	return g.DeleteClientRolesFromUser(ctx, token, realm, idOfClient, userID, roles)
}

// GetUserFederatedIdentities gets all user federated identities
func (g *GoKeycloak) GetUserFederatedIdentities(ctx context.Context, token, realm, userID string) (int, []*FederatedIdentityRepresentation, error) {
	const errMessage = "could not get user federated identities"

	var res []*FederatedIdentityRepresentation
	resp, err := g.GetRequestWithBearerAuth(ctx, token).
		SetResult(&res).
		Get(g.getAdminRealmURL(realm, "users", userID, "federated-identity"))

	if err := checkForError(resp, err, errMessage); err != nil {
		return resp.StatusCode(), nil, err
	}

	return resp.StatusCode(), res, err
}

// CreateUserFederatedIdentity creates an user federated identity
func (g *GoKeycloak) CreateUserFederatedIdentity(ctx context.Context, token, realm, userID, providerID string, federatedIdentityRep FederatedIdentityRepresentation) (int, error) {
	const errMessage = "could not create user federeated identity"

	resp, err := g.GetRequestWithBearerAuth(ctx, token).
		SetBody(federatedIdentityRep).
		Post(g.getAdminRealmURL(realm, "users", userID, "federated-identity", providerID))

	return resp.StatusCode(), checkForError(resp, err, errMessage)
}

// DeleteUserFederatedIdentity deletes an user federated identity
func (g *GoKeycloak) DeleteUserFederatedIdentity(ctx context.Context, token, realm, userID, providerID string) (int, error) {
	const errMessage = "could not delete user federeated identity"

	resp, err := g.GetRequestWithBearerAuth(ctx, token).
		Delete(g.getAdminRealmURL(realm, "users", userID, "federated-identity", providerID))

	return resp.StatusCode(), checkForError(resp, err, errMessage)
}
