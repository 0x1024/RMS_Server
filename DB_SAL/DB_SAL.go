package DB_SAL

import (
	"RMS_Server/AUTH_SAL"
	"RMS_Server/Public"
	"encoding/json"
	"errors"
	"fmt"
	_ "github.com/lib/pq"
	log "github.com/sirupsen/logrus"
	"github.com/xormplus/xorm"
	"golang.org/x/net/websocket"
	"reflect"
	"strconv"
	"strings"
	"time"
)

var engine *xorm.Engine

type Pd_index struct {
	Pid     uint64    "1,Pid"
	Dtype   string    "2,Dtype"
	Client  string    "3,Client"
	Tags    string    "4,Tags"
	Passwd  string    "5,Passwd"
	Created time.Time "6,Created"
	Updated time.Time "7,Updated"
}

//user manage
type Um_index struct {
	Pid     uint64 "1,Pid"
	Uid     uint64
	Name    string    "2,Name"
	Passwd  string    "3,Passwd"
	Role    string    "4,Level"
	Tags    string    "5,Tags"
	Created time.Time "6,Created"
	Updated time.Time "7,Updated"
	Stamp   uint64
	Jail    time.Duration
}

//client group
type Customer struct {
	Cid     uint64 "1,Cid"
	Uid     []uint64
	Pid     []uint64
	Created time.Time "6,Created"
	Updated time.Time "7,Updated"
}

//role group
type Role_group struct {
	Rid     uint64 "1,Rid"
	Name    string "json:name"
	Priv    uint
	Wlist   string //string={'"asdf","fda","fff" '}
	Blist   string
	Created time.Time "6,Created"
	Updated time.Time "7,Updated"
}

const (
	OP_Null   = 0
	OP_Ping   = 1 << 1
	OP_Read   = 1 << 2
	OP_Write  = 1 << 3
	OP_Create = 1 << 4
	OP_Delete = 1 << 5
	OP_Manage = 1 << 6
	OP_SysLv0 = 1 << 7
	OP_SysLv1 = 1 << 8
	OP_SysLv2 = 1 << 9
	OP_SysLv3 = 1 << 10
)

//type Privilege struct {
//	Read	bool
//	Create	bool
//	Update 	bool
//	Delete 	bool
//	ManPri	bool
//	SysLv	int
//}

func DB_Init() {
	var err error

	//=====================================================================
	//open db
	engine, err = xorm.NewPostgreSQL("postgres://rms:123@10.1.11.151:5432/RMSDB?sslmode=disable")
	if err != nil {
		log.Panic(err)
	}

	AUTH_SAL.AuthEng, err = xorm.NewPostgreSQL("postgres://rms:123@10.1.11.151:5432/AuthDB?sslmode=disable")
	if err != nil {
		log.Panic(err)
	}

	defer fmt.Printf("db_init closed")
	defer engine.Close()
	defer AUTH_SAL.AuthEng.Close()

	//=====================================================================
	//en func
	engine.ShowSQL(true)
	//	engine.Logger().SetLevel(core.LOG_DEBUG)
	engine.ShowExecTime(true)

	//	AUTH_SAL.AuthEng.Logger().SetLevel(core.LOG_DEBUG)

	AUTH_SAL.AuthEng.ShowSQL(true)
	//=====================================================================
	//check table
	err = engine.CreateTables(new(Pd_index))
	if err != nil {
		panic(err)
	}

	err = AUTH_SAL.AuthEng.CreateTables(new(Um_index))
	if err != nil {
		panic(err)
	}
	err = AUTH_SAL.AuthEng.CreateTables(new(Customer))
	if err != nil {
		panic(err)
	}
	err = AUTH_SAL.AuthEng.CreateTables(new(Role_group))
	if err != nil {
		panic(err)
	}

	for true {
	}
}

type cmd struct {
	Cmd  string      `json:"cmd"`
	Data interface{} `json:"data"`
}

func ReqProcess(ws *websocket.Conn, dat string) {
	var err error
	var send cmd

	var dats map[string]string

	var rrr string
	rrr = strings.TrimRight(dat, "\x00")

	if err := json.Unmarshal([]byte(rrr), &dats); err == nil {
		//fmt.Printf("\r\nReqProcess   %q \r\n", dats) //debug
		//fmt.Println("cmd: ", dats["cmd"])            //debug
	}

	Senders := new(Public.Senders)
	Senders.Ws = ws

	//auth login=====================================
	if Public.LoginUser[ws].Logined == false {

		if dats["cmd"] == "auth_req" {

			auth_tmp := new(Um_index)
			_, err = AUTH_SAL.AuthEng.Where("name=?", dats["user"]).Get(auth_tmp)

			if err != nil {
				fmt.Println(err)
				send.Cmd = "auth_failed"
				send.Data = ""

			} else if (auth_tmp.Name == "") || (auth_tmp.Name == "0") {
				send.Cmd = "auth_name_fault"
				send.Data = ""

			} else if (dats["pswd"] == auth_tmp.Passwd) && (dats["pswd"] != "") {
				send.Cmd = "auth_ok"
				send.Data = ""
				Public.LoginUser[ws].Logined = true

				rolegroup := new(Role_group)

				Public.LoginUser[ws].Role = auth_tmp.Role
				_, err = AUTH_SAL.AuthEng.Table("role_group").Where("name=?", auth_tmp.Role).Get(rolegroup)
				fmt.Println(rolegroup)
				Public.LoginUser[ws].Priv = rolegroup.Priv
				Public.LoginUser[ws].Wlist = rolegroup.Wlist
				Public.LoginUser[ws].Blist = rolegroup.Blist

			} else if dats["cmd"] == "HB" {
				Public.LoginUser[ws].HBLife = 0
			} else {
				send.Cmd = "auth_pwd_fault"
				send.Data = ""

			}

			rec, _ := json.Marshal(send)
			fmt.Printf("json %q \r\n==========%q\r\n  \r\n", rec, send)
			data_tmp := string(rec)
			Senders.Dat = data_tmp
			Public.DB2Ret <- Senders
		}
		//end of    auth login=====================================
	} else {

		//fmt.Print("priv   ")
		//fmt.Printf("%b", Public.LoginUser[ws].Priv)
		switch dats["cmd"] {
		case "req":
			fmt.Printf("\npriv  %X  ", Public.LoginUser[ws].Priv)
			if Public.LoginUser[ws].Priv&OP_Read != 0 {
				sa1 := new(Pd_index)
				re, _ := strconv.Atoi(dats["pid"])
				_, err = engine.Where("pid=?", re).Get(sa1)
				if err != nil {
					fmt.Println(err)
				}

				send.Cmd = "data_single"
				send.Data = sa1
				rec, _ := json.Marshal(send)
				data_tmp := string(rec)

				Senders.Dat = data_tmp
				Public.DB2Ret <- Senders
			} else {
				send.Cmd = dats["cmd"]
				authAct_NoPermition(Senders, send)
			}
		case "all":
			if Public.LoginUser[ws].Priv&OP_Read != 0 {
				sa2 := new([]Pd_index)
				err = engine.Find(sa2)

				send.Cmd = "data_all"
				send.Data = sa2

				rec, _ := json.Marshal(send)

				data_tmp := string(rec)

				Senders.Dat = data_tmp
				Public.DB2Ret <- Senders

			} else {
				send.Cmd = dats["cmd"]
				authAct_NoPermition(Senders, send)
			}
		case "comitone":
			if Public.LoginUser[ws].Priv&OP_Write != 0 {
				var recs int64
				//prepare data
				delete(dats, "cmd")
				result := &Pd_index{}
				fmt.Printf("%q \r\n\n", dats) //debug
				err = FillStruct(dats, result)

				//check exsit
				sa1 := new(Pd_index)
				_, err = engine.Where("pid=?", result.Pid).Get(sa1)
				if sa1.Pid == 0 {
					//no item,insert new one
					recs, err = engine.InsertOne(result)
				} else {
					//no err means there is item
					recs, err = engine.Update(result, &Pd_index{Pid: result.Pid})

				}

				fmt.Println(recs, err) //debug

				send.Cmd = "respond"
				send.Data = nil

				rec, _ := json.Marshal(send)

				data_tmp := string(rec)

				Senders.Dat = data_tmp
				Public.DB2Ret <- Senders

			} else {
				send.Cmd = dats["cmd"]
				authAct_NoPermition(Senders, send)
			}
		case "update":
			if Public.LoginUser[ws].Priv&OP_Write != 0 {
				delete(dats, "cmd")
				result := &Pd_index{}
				fmt.Printf("%q \r\n\n", dats) //debug

				err = FillStruct(dats, result)
				recs, err := engine.Update(result, &Pd_index{Pid: result.Pid})
				fmt.Println(recs, err) //debug

				send.Cmd = "respond"
				send.Data = nil

				rec, _ := json.Marshal(send)

				data_tmp := string(rec)

				Senders.Dat = data_tmp
				Public.DB2Ret <- Senders

			} else {
				send.Cmd = dats["cmd"]
				authAct_NoPermition(Senders, send)
			}
		case "delete_id":
			if Public.LoginUser[ws].Priv&OP_Delete != 0 {

				delete(dats, "cmd")
				result := &Pd_index{}
				fmt.Printf("%q \r\n\n", dats) //debug

				err = FillStruct(dats, result)

				n, err := engine.Delete(result)
				if err != nil {
					fmt.Println(n, err)
				}

				send.Cmd = "respond"
				send.Data = nil

				rec, _ := json.Marshal(send)

				data_tmp := string(rec)

				Senders.Dat = data_tmp
				Public.DB2Ret <- Senders

			} else {
				send.Cmd = dats["cmd"]
				authAct_NoPermition(Senders, send)
			}
		case "HB":
			if Public.LoginUser[ws].Priv&OP_Ping != 0 {

				Public.LoginUser[ws].HBLife = 0

			}
		default:
			send.Cmd = "respond"
			send.Data = nil

			rec, _ := json.Marshal(send)

			data_tmp := string(rec)
			Senders.Dat = data_tmp
			Public.DB2Ret <- Senders

		} //end of switch
	} //end of  if logined

}

func FillStruct(data map[string]string, obj interface{}) error {
	for k, v := range data {
		err := SetField(obj, k, v)
		if err != nil {
			return err
		}
	}
	return nil
}

//用map的值替换结构的值
func SetField(obj interface{}, name string, value interface{}) error {

	name = strings.Title(name)
	structValue := reflect.ValueOf(obj).Elem()        //结构体属性值
	structFieldValue := structValue.FieldByName(name) //结构体单个属性值

	if !structFieldValue.IsValid() {
		return fmt.Errorf("No such field: %s in obj", name)
	}

	if !structFieldValue.CanSet() {
		return fmt.Errorf("Cannot set %s field value", name)
	}

	structFieldType := structFieldValue.Type() //结构体的类型
	val := reflect.ValueOf(value)              //map值的反射值

	var err error
	if structFieldType != val.Type() {
		val, err = TypeConversion(fmt.Sprintf("%v", value), structFieldValue.Type().Name()) //类型转换
		if err != nil {
			return err
		}
	}

	structFieldValue.Set(val)
	return nil
}

//类型转换
func TypeConversion(value string, ntype string) (reflect.Value, error) {
	if ntype == "string" {
		return reflect.ValueOf(value), nil
	} else if ntype == "time.Time" {
		t, err := time.ParseInLocation("2006-01-02 15:04:05", value, time.Local)
		return reflect.ValueOf(t), err
	} else if ntype == "Time" {
		t, err := time.ParseInLocation("2006-01-02 15:04:05", value, time.Local)
		return reflect.ValueOf(t), err
	} else if ntype == "int" {
		i, err := strconv.Atoi(value)
		return reflect.ValueOf(i), err
	} else if ntype == "int8" {
		i, err := strconv.ParseInt(value, 10, 64)
		return reflect.ValueOf(int8(i)), err
	} else if ntype == "int32" {
		i, err := strconv.ParseInt(value, 10, 64)
		return reflect.ValueOf(int64(i)), err
	} else if ntype == "int64" {
		i, err := strconv.ParseInt(value, 10, 64)
		return reflect.ValueOf(i), err
	} else if ntype == "float32" {
		i, err := strconv.ParseFloat(value, 64)
		return reflect.ValueOf(float32(i)), err
	} else if ntype == "float64" {
		i, err := strconv.ParseFloat(value, 64)
		return reflect.ValueOf(i), err
	}

	//else if .......增加其他一些类型的转换

	return reflect.ValueOf(value), errors.New("未知的类型：" + ntype)
}

func authAct_NoPermition(Senders *Public.Senders, send cmd) {
	send.Data = fmt.Sprintf("%s,%s", send.Cmd, "No Permittion")
	send.Cmd = ""
	rec, _ := json.Marshal(send)
	data_tmp := string(rec)

	Senders.Dat = data_tmp
	Public.DB2Ret <- Senders
}
