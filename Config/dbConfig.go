package Config

import (
	"database/sql"
	"errors"
	"fmt"
	"gorm.io/driver/mysql"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var (
	DB    *gorm.DB
	sqlDb *sql.DB
	//WithUps = SqlArg{ArgId: 1}
	//WithJobs = SqlArg{ArgId: 2}
	//WithFanjus = SqlArg{ArgId: 3}
	//WithGroups = SqlArg{ArgId: 4}
)

func dbInit() {
	var err error
	conn := ""
	if CoreConfig.DBConfig.DBType == "mysql" {
		conn = fmt.Sprintf("%v:%v@tcp(%v:%v)/%v?parseTime=True&loc=Local", CoreConfig.DBConfig.DBUserName, CoreConfig.DBConfig.DBPassword, CoreConfig.DBConfig.DBIP, CoreConfig.DBConfig.DBPort, CoreConfig.DBConfig.DBName)
		DB, err = gorm.Open(mysql.Open(conn), &gorm.Config{})
	} else if CoreConfig.DBConfig.DBType == "sqlite3" {
		conn = fmt.Sprintf("./%v.db", CoreConfig.DBConfig.DBName)
		DB, err = gorm.Open(sqlite.Open(conn), &gorm.Config{})
	} else {
		panic(errors.New("not supported database adapter"))
	}
	if err != nil {
		panic(err)
	}
	sqlDb, err = DB.DB()
	if err == nil {
		sqlDb.SetMaxIdleConns(10)
		sqlDb.SetMaxOpenConns(100)
	}
	//err = DB.AutoMigrate(&QQGroup{}, &SQLUp{}, &SQLFanju{},&SQLJob{})
	//if err != nil {
	//	log.Println(err)
	//}
}

//type QQGroup struct {
//	GroupID			   int64	`gorm:"primaryKey"`
//	Enable             bool
//	AdminUin           int64
//	Menu               string
//	MenuKeyWord        string
//	ShutUpWord         string
//	ShutUpTime         int
//	JoinVerifyTime     int
//	JoinAutoShutUpTime int
//	Zan                bool
//	SignIn             bool
//	Bili               bool
//	BiliUps            []*SQLUp `gorm:"many2many:up_groups;foreignKey:GroupID"`
//	Fanjus             []*SQLFanju `gorm:"many2many:fanju_groups;foreignKey:GroupID"`
//	Welcome            string
//	JoinVerifyType     int
//	Job                []*SQLJob `gorm:"many2many:job_groups;foreignKey:GroupID"`
//}
//type SQLUp struct {
//	Name    string
//	Created int64
//	UserId  int64 `gorm:"primaryKey"`
//	Groups  []*QQGroup `gorm:"many2many:up_groups;foreignKey:UserId"`
//}
//type SQLJob struct {
//	gorm.Model
//	Cron    string
//	Type    int
//	Title   string
//	Content string
//	Groups  []*QQGroup `gorm:"many2many:job_groups;"`
//}
//type SQLFanju struct {
//	Title  string
//	Id     int64
//	UserId int64
//	Groups  []*QQGroup `gorm:"many2many:fanju_groups;foreignKey:Id"`
//}
//type SqlArg struct {
//	ArgId int
//}
//func SubUp(groupId,upId int64) {
//
//}
//func DelUpSub(groupId int64,upId int64) error {
//	if upId == 0 || groupId == 0 {
//		return errors.New("参数非法")
//	}
//	c := DB
//	if CoreConfig.Debug {
//		c = c.Debug()
//	}
//	err := c.Model(&QQGroup{GroupID: groupId}).Association("BiliUps").Delete(SQLUp{UserId: upId})
//	if err != nil {
//		return err
//	}
//	if DB.Model(&SQLUp{UserId: upId}).Association("Groups").Count() == 0 {
//		err = DB.Delete(&SQLUp{},upId).Error
//	}
//	return err
//}
//func GetGroupConfig(groupId int64,args ... SqlArg) (g QQGroup,e error) {
//	c := DB
//	if CoreConfig.Debug {
//		c = c.Debug()
//	}
//	for _,v := range args {
//		switch v.ArgId {
//		case 1:
//			c = c.Preload("BiliUps")
//		case 2:
//			c = c.Preload("Job")
//		case 3:
//			c = c.Preload("Fanjus")
//		}
//	}
//	e = c.Where("group_id = ?",groupId).First(&g).Error
//	return
//}
//func GetBiliUps(userId int64,args ... SqlArg) (g SQLUp,e error) {
//	c := DB
//	if CoreConfig.Debug {
//		c = c.Debug()
//	}
//	for _,v := range args {
//		switch v.ArgId {
//		case 4:
//			c = c.Preload("Groups")
//		}
//	}
//	e = c.Where("user_id = ?",userId).First(&g).Error
//	return
//}
