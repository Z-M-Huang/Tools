package db

import (
	"github.com/Z-M-Huang/Tools/data"
	"github.com/Z-M-Huang/Tools/utils"
	"github.com/jinzhu/gorm"

	_ "github.com/jinzhu/gorm/dialects/mssql" //supporting packages
	_ "github.com/jinzhu/gorm/dialects/mysql"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
)

//dbContext database connection
var dbContext *gorm.DB

func init() {
	var err error
	dbContext, err = gorm.Open(data.Config.DatabaseConfig.Driver, data.Config.DatabaseConfig.ConnectionString)
	if err != nil {
		utils.Logger.Sugar().Fatalf("failed to open database %s", err.Error())
	}
	dbContext.AutoMigrate(&User{}, &Application{})
}

//DoTransaction do transaction
func DoTransaction(fc func(tx *gorm.DB) error) error {
	return dbContext.Transaction(fc)
}
