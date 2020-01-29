package car

import (
	"github.com/go-pg/pg"
	"github.com/go-pg/pg/orm"
	gorsk "github.com/simonhylander/gorsk/pkg/utl/model"
)

// represents car application service
type Car struct {
	db   *pg.DB
	repository CarRepository
}

// represents car repository interface
type CarRepository interface {
	List(orm.DB, *gorsk.Pagination) ([]gorsk.Car, error)
}
