package transport

import (
	"net/http"
	"strconv"

	"github.com/simonhylander/gorsk/pkg/api/user"

	gorsk "github.com/simonhylander/gorsk/pkg/utl/model"

	"github.com/labstack/echo"
)

// HTTP represents user http service
type HTTP struct {
	svc user.Service
}

// NewHTTP creates new user http service
func NewHTTP(svc user.Service, er *echo.Group) {
	h := HTTP{svc}
	ur := er.Group("/users")
	ur.GET("", h.list)
	ur.GET("/:id", h.view)
	ur.POST("", h.create)
	ur.PATCH("/:id", h.update) // TODO PUT
	ur.DELETE("/:id", h.delete)
}

// Custom errors
var (
	ErrPasswordsNotMaching = echo.NewHTTPError(http.StatusBadRequest, "Passwords do not match")
)

// User create request
// swagger:model userCreate
type createReq struct {
	FirstName       string `json:"first_name" validate:"required"`
	LastName        string `json:"last_name" validate:"required"`
	Username        string `json:"username" validate:"required,min=3,alphanum"`
	Password        string `json:"password" validate:"required,min=8"`
	PasswordConfirm string `json:"password_confirm" validate:"required"`
	Email           string `json:"email" validate:"required,email"`

	CompanyID  int              `json:"company_id" validate:"required"`
	LocationID int              `json:"location_id" validate:"required"`
	RoleID     gorsk.AccessRole `json:"role_id" validate:"required"`
}

func (h *HTTP) create(c echo.Context) error {
	r := new(createReq)

	if err := c.Bind(r); err != nil {

		return err
	}

	if r.Password != r.PasswordConfirm {
		return ErrPasswordsNotMaching
	}

	if r.RoleID < gorsk.SuperAdminRole || r.RoleID > gorsk.UserRole {
		return gorsk.ErrBadRequest
	}

	usr, err := h.svc.Create(c, gorsk.User{
		Username:   r.Username,
		Password:   r.Password,
		Email:      r.Email,
		FirstName:  r.FirstName,
		LastName:   r.LastName,
		CompanyID:  r.CompanyID,
		LocationID: r.LocationID,
		RoleID:     r.RoleID,
	})

	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, usr)
}

type listResponse struct {
	Users []gorsk.User `json:"users"`
	Page  int          `json:"page"`
}

func (h *HTTP) list(c echo.Context) error {
	p := new(gorsk.PaginationReq)
	if err := c.Bind(p); err != nil {
		return err
	}

	result, err := h.svc.List(c, p.Transform())

	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, listResponse{result, p.Page})
}

func (h *HTTP) view(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return gorsk.ErrBadRequest
	}

	result, err := h.svc.View(c, id)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, result)
}

// User update request
type updateReq struct {
	ID        int    `json:"-"`
	FirstName string `json:"first_name,omitempty" validate:"omitempty,min=2"`
	LastName  string `json:"last_name,omitempty" validate:"omitempty,min=2"`
	Mobile    string `json:"mobile,omitempty"`
	Phone     string `json:"phone,omitempty"`
	Address   string `json:"address,omitempty"`
}

func (h *HTTP) update(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return gorsk.ErrBadRequest
	}

	req := new(updateReq)
	if err := c.Bind(req); err != nil {
		return err
	}

	usr, err := h.svc.Update(c, &user.Update{
		ID:        id,
		FirstName: req.FirstName,
		LastName:  req.LastName,
		Mobile:    req.Mobile,
		Phone:     req.Phone,
		Address:   req.Address,
	})

	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, usr)
}

func (h *HTTP) delete(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return gorsk.ErrBadRequest
	}

	if err := h.svc.Delete(c, id); err != nil {
		return err
	}

	return c.NoContent(http.StatusOK)
}