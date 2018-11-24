package main

import (
	"fmt"

	nestedset "gitlab.com/sulthonzh/scraperpath-nested-set"
	"gitlab.com/sulthonzh/scraperpath-nested-set/examples/db"
	"gitlab.com/sulthonzh/scraperpath-nested-set/examples/utils"
)

func main() {
	migrate()
}

func migrate() {
	db, err := db.InitDB()
	if err != nil {
		utils.ExitOnFailure(err)
	}

	categories := nestedset.Node{}
	categories.SetTableName("scraper_paths")

	fmt.Println(categories.TableName())
	err = db.Set("gorm:table_options", "ENGINE=InnoDB").AutoMigrate(
		&categories,
	).Error
	if err != nil {
		utils.ExitOnFailure(err)
	}

	err = db.Exec(`
		INSERT INTO scraper_paths (id,type,name,lft,rgt,data) VALUES 
		('5439f1ff-cb6a-4793-acbf-4d80ede90187',1,'description',5,6,'{"field": "---", "target": "desc"}'),
		('628fbfb9-87b3-4d29-bda2-50671001fd51',1,'name',3,4,'{"field": "---", "target": "name"}'),
		('7f9960e8-1f9a-4ee3-a204-9e801535d96a',1,'title',10,11,NULL),
		('806dd400-0d44-4854-8448-34e9d3f35266',1,'listings',1,12,'{"field": "---", "target": "listings"}'),
		('8d5fceba-2ee3-46c3-8eaf-f06cdb4fd68e',1,'price',7,8,'{"field": "---", "target": "price"}'),
		('fddfd29f-c19f-4112-b556-29ca9d132140',1,'url',2,9,'{"field": "---", "target": "url"}');
`).Error
	if err != nil {
		utils.ExitOnFailure(err)
	}
	fmt.Println("migrate: completed")
}
