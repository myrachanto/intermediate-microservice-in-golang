package repository

import (
	"fmt"
	"strings"
	"log"
	"os"
	"github.com/joho/godotenv"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"github.com/myrachanto/allmicro/gormmicro/categorymicroservice/httperors"
	"github.com/myrachanto/allmicro/gormmicro/categorymicroservice/model"
	"github.com/myrachanto/allmicro/gormmicro/categorymicroservice/support"
)

var (
	Sqlrepository sqlrepository = sqlrepository{}
)

///curtesy to gorm
type sqlrepository struct{}
func init(){
	Getconnected()
}

func Getconnected() (GormDB *gorm.DB, err *httperors.HttpError) {
	err2 := godotenv.Load()
	if err2 != nil {
		log.Fatal("Error loading .env file in routes")
	}
	dbuser := os.Getenv("DbUsername")
	DbName := os.Getenv("DbName")
	dbURI := dbuser+"@/"+DbName+"?charset=utf8&parseTime=True&loc=Local"
	GormDB, err1 := gorm.Open("mysql", dbURI)
	if err1 != nil {
		return nil, httperors.NewNotFoundError("No Mysql db connection")
	}

	GormDB.AutoMigrate(&model.Category{})
	GormDB.AutoMigrate(&model.MajorCategory{})
	GormDB.AutoMigrate(&model.Category{})
	GormDB.AutoMigrate(&model.SubCategory{})
	return GormDB, nil
}
func DbClose(GormDB *gorm.DB) {
	defer GormDB.Close()
}
func (repository sqlrepository) Create(category *model.Category) (*model.Category, *httperors.HttpError) {
	if err := category.Validate(); err != nil {
		return nil, err
	}
	GormDB, err1 := Getconnected()
	if err1 != nil {
		return nil, err1
	}
	GormDB.Create(&category)
	DbClose(GormDB)
	return category, nil
}
func (repository sqlrepository) GetOne(id int) (*model.Category, *httperors.HttpError) {
	ok := repository.ProductUserExistByid(id)
	if !ok {
		return nil, httperors.NewNotFoundError("category with that id does not exists!")
	}
	category := model.Category{}
	GormDB, err1 := Getconnected()
	if err1 != nil {
		return nil, err1
	}
	
	GormDB.Model(&category).Where("id = ?", id).First(&category)
	DbClose(GormDB)
	
	return &category, nil
}

func (repository sqlrepository) GetAll(categorys []model.Category,search *support.Search) ([]model.Category, *httperors.HttpError) {
	results, err1 := repository.Search(search, categorys)
	if err1 != nil {
			return nil, err1
		}
	return results, nil
}

func (repository sqlrepository) Update(id int, category *model.Category) (*model.Category, *httperors.HttpError) {
	ok := repository.ProductUserExistByid(id)
	if !ok {
		return nil, httperors.NewNotFoundError("category with that id does not exists!")
	}
	
	GormDB, err1 := Getconnected()
	if err1 != nil {
		return nil, err1
	}
	Category := model.Category{}
	acategory := model.Category{}
	
	GormDB.Model(&Category).Where("id = ?", id).First(&acategory)
	if category.Name  == "" {
		category.Name = acategory.Name
	}
	if category.Title  == "" {
		category.Title = acategory.Title
	}
	if category.Description  == "" {
		category.Description = acategory.Description
	}
	GormDB.Model(&Category).Where("id = ?", id).First(&Category).Update(&acategory)
	
	DbClose(GormDB)

	return category, nil
}
func (repository sqlrepository) Delete(id int) (*httperors.HttpSuccess, *httperors.HttpError) {
	ok := repository.ProductUserExistByid(id)
	if !ok {
		return nil, httperors.NewNotFoundError("Product with that id does not exists!")
	}
	category := model.Category{}
	GormDB, err1 := Getconnected()
	if err1 != nil {
		return nil, err1
	}
	GormDB.Model(&category).Where("id = ?", id).First(&category)
	GormDB.Delete(category)
	DbClose(GormDB)
	return httperors.NewSuccessMessage("deleted successfully"), nil
}
func (repository sqlrepository)ProductUserExistByid(id int) bool {
	category := model.Category{}
	GormDB, err1 := Getconnected()
	if err1 != nil {
		return false
	}
	if GormDB.First(&category, "id =?", id).RecordNotFound(){
	   return false
	}
	DbClose(GormDB)
	return true
	
}


func (repository sqlrepository) Search(Ser *support.Search, categorys []model.Category)([]model.Category, *httperors.HttpError){
	GormDB, err1 := Getconnected()
	if err1 != nil {
		return nil, err1
	}
	category := model.Category{}
	switch(Ser.Search_operator){
	case "all":
		q := GormDB.Model(&category).Find(&categorys)
		///////////////////////////////////////////////////////////////////////////////////////////////////////
		///////////////find some other paginator more effective one///////////////////////////////////////////
		p := paginator.New(adapter.NewGORMAdapter(q), Ser.Per_page)
		p.SetPage(1)
		
		if err3 := p.Results(&categorys); err3 != nil {
			return nil, httperors.NewNotFoundError("something went wrong paginating!")
		}
	break;
	case "equal_to":
		q := GormDB.Where(Ser.Search_column+" "+Operator[Ser.Search_operator]+"?", Ser.Search_query_1).Find(&categorys);
		p := paginator.New(adapter.NewGORMAdapter(q), Ser.Per_page)
		p.SetPage(1)
		
		if err3 := p.Results(&categorys); err3 != nil {
			return nil, httperors.NewNotFoundError("something went wrong paginating!")
		}
	break;
	case "not_equal_to":
		q := GormDB.Where(Ser.Search_column+" "+Operator[Ser.Search_operator]+"?", Ser.Search_query_1).Find(&categorys);	
		p := paginator.New(adapter.NewGORMAdapter(q), Ser.Per_page)
		p.SetPage(1)
		
		if err3 := p.Results(&categorys); err3 != nil {
			return nil, httperors.NewNotFoundError("something went wrong paginating!")
		}
	break;
	case "less_than" :
		q := GormDB.Where(Ser.Search_column+" "+Operator[Ser.Search_operator]+"?", Ser.Search_query_1).Find(&categorys);	
		p := paginator.New(adapter.NewGORMAdapter(q), Ser.Per_page)
		p.SetPage(1)
		
		if err3 := p.Results(&categorys); err3 != nil {
			return nil, httperors.NewNotFoundError("something went wrong paginating!")
		}
	break;
	case "greater_than":
		q := GormDB.Where(Ser.Search_column+" "+Operator[Ser.Search_operator]+"?", Ser.Search_query_1).Find(&categorys);	
		p := paginator.New(adapter.NewGORMAdapter(q), Ser.Per_page)
		p.SetPage(1)
		
		if err3 := p.Results(&categorys); err3 != nil {
			return nil, httperors.NewNotFoundError("something went wrong paginating!")
		}
	break;
	case "less_than_or_equal_to":
		q := GormDB.Where(Ser.Search_column+" "+Operator[Ser.Search_operator]+"?", Ser.Search_query_1).Find(&categorys);	
		p := paginator.New(adapter.NewGORMAdapter(q), Ser.Per_page)
		p.SetPage(1)
		
		if err3 := p.Results(&categorys); err3 != nil {
			return nil, httperors.NewNotFoundError("something went wrong paginating!")
		}
	break;
	case "greater_than_ro_equal_to":
		q := GormDB.Where(Ser.Search_column+" "+Operator[Ser.Search_operator]+"?", Ser.Search_query_1).Find(&categorys);	
		p := paginator.New(adapter.NewGORMAdapter(q), Ser.Per_page)
		p.SetPage(1)
		
		if err3 := p.Results(&categorys); err3 != nil {
			return nil, httperors.NewNotFoundError("something went wrong paginating!")
		}
	break;
		 case "in":
			// db.Where("name IN (?)", []string{"myrachanto", "anto"}).Find(&users)
		s := strings.Split(Ser.Search_query_1,",")
		fmt.Println(s)
		q := GormDB.Where(Ser.Search_column+" "+Operator[Ser.Search_operator]+"(?)", s).Find(&categorys);
		p := paginator.New(adapter.NewGORMAdapter(q), Ser.Per_page)
		p.SetPage(1)
		
		if err3 := p.Results(&categorys); err3 != nil {
			return nil, httperors.NewNotFoundError("something went wrong paginating!")
		}
		break;
	 case "not_in":
			//db.Not("name", []string{"jinzhu", "jinzhu 2"}).Find(&users)
		s := strings.Split(Ser.Search_query_1,",")
		q := GormDB.Not(Ser.Search_column, s).Find(&categorys);
		p := paginator.New(adapter.NewGORMAdapter(q), Ser.Per_page)
		p.SetPage(1)
		
		if err3 := p.Results(&categorys); err3 != nil {
			return nil, httperors.NewNotFoundError("something went wrong paginating!")
		}
	// break;
	case "like":
		q := GormDB.Where(Ser.Search_column+" "+Operator[Ser.Search_operator]+"?", "%"+Ser.Search_query_1+"%").Find(&categorys);
		p := paginator.New(adapter.NewGORMAdapter(q), Ser.Per_page)
		p.SetPage(1)
		
		if err3 := p.Results(&categorys); err3 != nil {
			return nil, httperors.NewNotFoundError("something went wrong paginating!")
		}
	break;
	case "between":
		//db.Where("name BETWEEN ? AND ?", "lastWeek, today").Find(&users)
		q := GormDB.Where(Ser.Search_column+" "+Operator[Ser.Search_operator]+"? AND ?", Ser.Search_query_1, Ser.Search_query_2).Find(&categorys);
		p := paginator.New(adapter.NewGORMAdapter(q), Ser.Per_page)
		p.SetPage(1)
		
		if err3 := p.Results(&categorys); err3 != nil {
			return nil, httperors.NewNotFoundError("something went wrong paginating!")
		}
	   break;
	default:
	return nil, httperors.NewNotFoundError("check your operator!")
	}
	categoryRepo.DbClose(GormDB)
	
	return categorys, nil
}