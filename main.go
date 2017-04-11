package main

import (
	"RMS_Server/DB_SAL"
	"RMS_Server/Public"
	"RMS_Server/WEB_IO"
	"fmt"
	"golang.org/x/net/websocket"
	"time"
)

//start up without gui option
//go build -ldflags "-s -w -H=windowsgui"

func main() {
	//cmd:=exec.Command("notepad", "123" )
	////var out bytes.Buffer
	////cmd.Stdout=&out
	//err:=cmd.Run()
	//if err == nil {}

	bits := 1024
	if err := WEB_IO.GenRsaKey(bits); err != nil {
		print("密钥文件生成失败！")
	}
	print("密钥文件生成成功！")

	var c, d []byte
	c, _ = WEB_IO.RsaEncrypt([]byte("hello world111111111111111111111111111111111111111111111111111111111111111"))

	a := time.Now().UnixNano()
	for i := 0; i < 1000; i++ {
		d, _ = WEB_IO.RsaDecrypt(c)
	}
	b := time.Now().UnixNano()
	fmt.Println("\n\ntime     ", a, b, (b - a))

	fmt.Println()
	fmt.Printf("%s", d)

	Public.LoginUser = make(map[*websocket.Conn]*Public.LoginType)

	fmt.Println("logintyoe %q：", Public.LoginUser)
	fmt.Println("loginsrt %+v：", Public.LoginUser)

	//go AUTH_SAL.AuthDB_Init()
	go DB_SAL.DB_Init()
	go WEB_IO.Http_init()

	for true {
		time.Sleep(time.Second * 60)
	}
}
