package database

import (
	"fmt"

	"github.com/beego/beego/v2/client/orm"
	beego "github.com/beego/beego/v2/server/web"

	_ "github.com/lib/pq"
)

func Init() {
	registerDatabase()
	verifyConnection()
}

func registerDatabase() {
	if err := orm.RegisterDriver("postgres", orm.DRPostgres); err != nil {
		panic(fmt.Errorf("failed to register postgres driver: %w", err))
	}

	host := beego.AppConfig.DefaultString("database::host", "localhost")
	port := beego.AppConfig.DefaultString("database::port", "5432")
	name := beego.AppConfig.DefaultString("database::name", "rent_by_owner")
	user := beego.AppConfig.DefaultString("database::user", "postgres")
	password := beego.AppConfig.DefaultString("database::password", "postgres")
	sslMode := beego.AppConfig.DefaultString("database::sslmode", "disable")

	dataSource := fmt.Sprintf(
		"user=%s password=%s dbname=%s host=%s port=%s sslmode=%s",
		user,
		password,
		name,
		host,
		port,
		sslMode,
	)

	if err := orm.RegisterDataBase(
		"default",
		"postgres",
		dataSource,
	); err != nil {
		panic(fmt.Errorf("failed to register database: %w", err))
	}
}

func verifyConnection() {
	o := orm.NewOrm()

	var result int

	if err := o.Raw("SELECT 1").QueryRow(&result); err != nil {
		panic(fmt.Errorf("failed to connect to database: %w", err))
	}
}