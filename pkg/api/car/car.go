package car

import (
	"github.com/labstack/echo"
	gorsk "github.com/simonhylander/gorsk/pkg/utl/model"
)

// List returns list of users
func (car *Car) List(c echo.Context, p *gorsk.Pagination) ([]gorsk.Car, error) {
	return car.repository.List(car.db, p)
}

/*func (u *User) List(c echo.Context, p *gorsk.Pagination) ([]gorsk.User, error) {
	au := u.rbac.User(c)
	q, err := query.List(au)
	if err != nil {
		return nil, err
	}
	return u.udb.List(u.db, q, p)
}*/