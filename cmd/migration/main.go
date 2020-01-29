package main

import (
	"fmt"
	"log"
	"strings"

	"github.com/simonhylander/gorsk/pkg/utl/model"
	"github.com/simonhylander/gorsk/pkg/utl/secure"

	"github.com/go-pg/pg"
	"github.com/go-pg/pg/orm"
)

func main() {
	dbInsert := `INSERT INTO public.companies VALUES (1, now(), now(), NULL, 'admin_company', true);
	INSERT INTO public.locations VALUES (1, now(), now(), NULL, 'admin_location', true, 'admin_address', 1);
	INSERT INTO public.roles VALUES (100, 100, 'SUPER_ADMIN');
	INSERT INTO public.roles VALUES (110, 110, 'ADMIN');
	INSERT INTO public.roles VALUES (120, 120, 'COMPANY_ADMIN');
	INSERT INTO public.roles VALUES (130, 130, 'LOCATION_ADMIN');
	INSERT INTO public.roles VALUES (200, 200, 'USER');
	INSERT INTO public.cars VALUES (1, now(), now(), null, 'Big blue car', 1);`
	var psn = `postgres://simon:docker@localhost:5432/salestrigger?sslmode=disable`
	queries := strings.Split(dbInsert, ";")

	//ParseURL parses an URL into options that can be used to connect to PostgreSQL.
	options, err := pg.ParseURL(psn)
	checkErr(err)

	fmt.Println(psn)

	db := pg.Connect(options)
	_, err = db.Exec("SELECT 1")

	fmt.Println(err)

	checkErr(err)

	createSchema(db, &gorsk.Company{}, &gorsk.Location{}, &gorsk.Role{}, &gorsk.User{}, &gorsk.Car{})

	fmt.Println(dbInsert)

	for _, v := range queries[0 : len(queries)-1] {
		_, err := db.Exec(v)
		checkErr(err)
	}

	sec := secure.New(1, nil)

	userInsert := `INSERT INTO public.users (id, created_at, updated_at, first_name, last_name, username, password, email, active, role_id, company_id, location_id) VALUES (1, now(),now(),'Admin', 'Admin', 'admin', '%s', 'johndoe@mail.com', true, 100, 1, 1);`
	_, err = db.Exec(fmt.Sprintf(userInsert, sec.Hash("admin")))
	checkErr(err)
}

func checkErr(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func createSchema(db *pg.DB, models ...interface{}) {
	for _, model := range models {
		checkErr(db.CreateTable(model, &orm.CreateTableOptions{
			FKConstraints: true,
		}))
	}
}