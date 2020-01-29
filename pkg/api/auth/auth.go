package auth

import (
	"fmt"
	"github.com/labstack/echo"
	"github.com/simonhylander/gorsk/pkg/utl/model"
	"net/http"
)

// Custom errors
var (
	ErrInvalidCredentials = echo.NewHTTPError(http.StatusUnauthorized, "Username or password does not exist")
)

// Authenticate tries to authenticate the user provided by username and password
func (a *Auth) Authenticate(c echo.Context, user, pass string) (*gorsk.AuthToken, error) {
	u, err := a.udb.FindByUsername(a.db, user)
	if err != nil {
		return nil, err
	}

	if !a.sec.HashMatchesPassword(u.Password, pass) {
		return nil, ErrInvalidCredentials
	}

	if !u.Active {
		return nil, gorsk.ErrUnauthorized
	}

	token, expire, err := a.tg.GenerateToken(u)
	if err != nil {
		return nil, gorsk.ErrUnauthorized
	}

	fmt.Println(token)

	u.UpdateLastLogin(a.sec.Token(token))

	if err := a.udb.Update(a.db, u); err != nil {
		return nil, err
	}

	return &gorsk.AuthToken{Token: token, Expires: expire, RefreshToken: u.Token}, nil
}

// Refresh refreshes jwt token and puts new claims inside
func (a *Auth) Refresh(c echo.Context, token string) (*gorsk.RefreshToken, error) {
	user, err := a.udb.FindByToken(a.db, token)

	if err != nil {
		return nil, echo.NewHTTPError(http.StatusNotFound, "User not found")
	}

	token, expire, err := a.tg.GenerateToken(user)

	if err != nil {
		return nil, err
	}

	return &gorsk.RefreshToken{Token: token, Expires: expire}, nil
}

// Me returns info about currently logged user
func (a *Auth) Me(c echo.Context) (*gorsk.User, error) {
	au := a.rbac.User(c)
	return a.udb.View(a.db, au.ID)
}