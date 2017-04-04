package main

import (
	"fmt"

	"github.com/jinzhu/gorm"
	_ "github.com/lib/pq"
)

// SurveyWorkgroup is survey_workgroup
type SurveyWorkgroup struct {
	CaseID int
	WorkID int
}

// TableName setting table name is survey_workgroup
func (c SurveyWorkgroup) TableName() string {
	return "survey_workgroup"
}

func main() {

	db, err := gorm.Open("postgres", "host=192.168.99.100 user=admin password=testtest dbname=enet sslmode=disable")

	if err != nil {
		panic(err)
	}

	db.DB()
	// Enable Logger
	db.LogMode(true)

	err = db.DB().Ping()
	data := SurveyWorkgroup{}

	db.Debug().Select("case_id, work_id").First(&data)
	fmt.Println(data)
}
