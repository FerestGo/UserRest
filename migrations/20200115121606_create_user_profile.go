package main

import (
	"github.com/go-pg/pg/orm"
	migrations "github.com/robinjoseph08/go-pg-migrations"
)

type User struct {
	tableName string `pg:"user" json:"-"` // nolint
	Id        int32
	Login     string `pg:",unique"`
	Password  string
	IsGuest   bool
	Email     string     `pg:",unique"`
	CreatedAt time.Time  `pg:",default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt time.Time  `pg:",default:CURRENT_TIMESTAMP" json:"updated_at"`
	DeletedAt *time.Time `pg:",soft_delete" json:"deleted_at,omitempty"`
}

func init() {
	up := func(db orm.DB) error {
		err := orm.CreateTable(db, &User{}, &orm.CreateTableOptions{FKConstraints: true, Varchar: 255})

		return err
	}

	down := func(db orm.DB) error {
		err = orm.DropTable(db, &User{}, &orm.DropTableOptions{})

		return err
	}

	opts := migrations.MigrationOptions{}

	migrations.Register("20200115121606_create_user_profile", up, down, opts)
}
