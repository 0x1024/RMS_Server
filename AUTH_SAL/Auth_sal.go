package AUTH_SAL

import (
	"RMS_Server/Public"
	"fmt"
	_ "github.com/lib/pq"
	"github.com/xormplus/xorm"
	"time"
)

//user manage
type Um_index struct {
	Pid     int64     "1,Pid"
	Name    string    "2,Name"
	Passwd  string    "3,Passwd"
	Level   string    "4,Level"
	Tags    string    "5,Tags"
	Created time.Time "6,Created"
	Updated time.Time "7,Updated"
}

var AuthEng *xorm.Engine

func AuthDB_Init() {
	var err error

	//=====================================================================
	//open db
	//AuthEng, err = xorm.NewPostgreSQL("postgres://rms:123@10.1.11.151:5432/AuthDB?sslmode=disable")
	//if err != nil {
	//	log.Panic(err)
	//}
	<-Public.Signal

	defer fmt.Printf("db_init closed")
	defer AuthEng.Close()

	//=====================================================================
	//en func
	//AuthEng.ShowSQL(true)
	//AuthEng.Logger().SetLevel(core.LOG_DEBUG)
	//AuthEng.ShowExecTime(true)

	//=====================================================================
	//check table
	err = AuthEng.CreateTables(new(Um_index))
	if err != nil {
		panic(err)
	}

	for true {
	}
}
