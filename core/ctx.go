package core

import (
	"context"
	"fmt"
	"strings"

	"github.com/google/uuid"
	"golang.org/x/exp/slices"
)

type AuthorizationOption interface {
	Evaluate(Ctx) *Exception
}

type HasPermission struct {
	Permission string
}

type HasPermissions struct {
	All         bool
	Permissions []string
}

type HasRole string

func (hr HasRole) Evaluate(appCtx Ctx) *Exception {
	requiredRole := string(hr)

	for _, v := range appCtx.Roles() {
		if strings.EqualFold(v, requiredRole) {
			return nil
		}
	}

	return Forbidden(WithMessage(fmt.Sprintf("user does not have required role: %s", requiredRole)))
}

func (hp HasPermission) Evaluate(appCtx Ctx) *Exception {
	for _, v := range appCtx.Permissions() {
		if strings.EqualFold(v, hp.Permission) {
			return nil
		}
	}

	return Forbidden(WithMessage(fmt.Sprintf("user does not have required permission: %s", hp.Permission)))
}

type HasRoles struct {
	All   bool
	Roles []string
}

func (hr HasRoles) Evaluate(appCtx Ctx) *Exception {
	roles := appCtx.Roles()
	requiredRoles := hr.Roles

	var checks []bool
	for _, r := range requiredRoles {
		checks = append(checks, slices.Contains(roles, r))
	}

	hasAll := !slices.Contains(checks, false)
	hasOne := slices.Contains(checks, true)

	if hr.All && !hasAll {
		return Forbidden(WithMessage(fmt.Sprintf("user does not have required roles: %s", strings.Join(requiredRoles, ","))))
	} else if !hr.All && !hasOne {
		return Forbidden(WithMessage(fmt.Sprintf("user does not have required roles: %s", strings.Join(requiredRoles, ","))))
	}

	return nil
}

func (hp HasPermissions) Evaluate(appCtx Ctx) *Exception {
	permissions := appCtx.Permissions()
	requiredPermissions := hp.Permissions

	var checks []bool
	for _, p := range requiredPermissions {
		checks = append(checks, slices.Contains(permissions, p))
	}

	hasAll := !slices.Contains(checks, false)
	hasOne := slices.Contains(checks, true)

	if hp.All && !hasAll {
		return Forbidden(WithMessage(fmt.Sprintf("user does not have required permissions: %s", strings.Join(requiredPermissions, ","))))
	} else if !hp.All && !hasOne {
		return Forbidden(WithMessage(fmt.Sprintf("user does not have required permissions: %s", strings.Join(requiredPermissions, ","))))
	}

	return nil
}

type Ctx interface {
	ID() string
	UserID() string
	Authorize(...AuthorizationOption) *Exception

	Roles() []string
	Permissions() []string

	SetUserID(string) Ctx
	SetRoles([]string) Ctx
	SetPermissions([]string) Ctx
}

type UserCtx struct {
	ContextID        string   `json:"contextId"`
	UserIDField      string   `json:"userId"`
	ProviderField    string   `json:"provider"`
	RolesField       []string `json:"roles"`
	PermissionsField []string `json:"permissions"`
}

// ID returns the current context ID
func (ctx *UserCtx) ID() string {
	return ctx.ContextID
}

func (ctx *UserCtx) Authorize(opts ...AuthorizationOption) *Exception {
	for _, opt := range opts {
		if err := opt.Evaluate(ctx); err != nil {
			return err
		}
	}

	return nil
}

// UserID return the current user
func (ctx *UserCtx) UserID() string {
	return ctx.UserIDField
}

func (ctx *UserCtx) SetUserID(userID string) Ctx {
	ctx.UserIDField = userID
	return ctx
}

func (ctx *UserCtx) SetRoles(roles []string) Ctx {
	ctx.RolesField = roles
	return ctx
}

func (ctx *UserCtx) SetPermissions(permissions []string) Ctx {
	ctx.PermissionsField = permissions
	return ctx
}

func (ctx *UserCtx) Roles() []string {
	return ctx.RolesField
}

func (ctx *UserCtx) Permissions() []string {
	return ctx.PermissionsField
}

func NewUserCtx(userID string, roles, permissions []string) Ctx {
	return &UserCtx{
		ContextID:        uuid.New().String(),
		UserIDField:      userID,
		RolesField:       roles,
		PermissionsField: permissions,
	}
}

type userCtxKey struct{}

// FromContext returns the Ctx stored in a context.
func FromContext(ctx context.Context) (Ctx, bool) {
	v, ok := ctx.Value(userCtxKey{}).(Ctx)
	return v, ok
}

// NewContext returns a copy of parent context with the given Ctx attached with it.
func NewContext(parent context.Context, v Ctx) context.Context {
	return context.WithValue(parent, userCtxKey{}, v)
}
