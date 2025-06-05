package database

import (
	"database/sql"
	"testing"

	"github.com/mobenaus/fc-pos-go-desafio-cleancode/internal/entity"
	"github.com/stretchr/testify/suite"

	// sqlite3

	_ "github.com/mattn/go-sqlite3"

	// migrations
	"github.com/golang-migrate/migrate/v4"
	sqlitemigrate "github.com/golang-migrate/migrate/v4/database/sqlite3"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

type OrderRepositoryTestSuite struct {
	suite.Suite
	Db *sql.DB
}

func (suite *OrderRepositoryTestSuite) SetupSuite() {
	db, err := sql.Open("sqlite3", ":memory:")
	suite.NoError(err)
	//db.Exec("CREATE TABLE orders (id varchar(255) NOT NULL, price float NOT NULL, tax float NOT NULL, final_price float NOT NULL, PRIMARY KEY (id))")

	instance, err := sqlitemigrate.WithInstance(db, &sqlitemigrate.Config{})
	if err != nil {
		panic(err)
	}

	// o file abaixo tem que ter o ../../.. para referenciar a raiz do projeto
	migration, err := migrate.NewWithDatabaseInstance("file://../../../migrations/", "mysql", instance)
	if err != nil {
		panic(err)
	}

	if err := migration.Up(); err != nil {
		panic(err)
	}

	suite.Db = db
}

func (suite *OrderRepositoryTestSuite) TearDownTest() {
	suite.Db.Close()
}

func TestSuite(t *testing.T) {
	suite.Run(t, new(OrderRepositoryTestSuite))
}

func (suite *OrderRepositoryTestSuite) TestGivenAnOrder_WhenSave_ThenShouldSaveOrder() {
	order, err := entity.NewOrder("123", 10.0, 2.0)
	suite.NoError(err)
	suite.NoError(order.CalculateFinalPrice())
	repo := NewOrderRepository(suite.Db)
	err = repo.Save(order)
	suite.NoError(err)

	var orderResult entity.Order
	err = suite.Db.QueryRow("Select id, price, tax, final_price from orders where id = ?", order.ID).
		Scan(&orderResult.ID, &orderResult.Price, &orderResult.Tax, &orderResult.FinalPrice)

	suite.NoError(err)
	suite.Equal(order.ID, orderResult.ID)
	suite.Equal(order.Price, orderResult.Price)
	suite.Equal(order.Tax, orderResult.Tax)
	suite.Equal(order.FinalPrice, orderResult.FinalPrice)
}
