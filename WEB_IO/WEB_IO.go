package WEB_IO

import (
	"RMS_Server/DB_SAL"
	"RMS_Server/Public"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/json"
	"encoding/pem"
	"errors"
	"fmt"
	log "github.com/sirupsen/logrus"
	"golang.org/x/net/websocket"
	"io"
	"io/ioutil"
	rand2 "math/rand"
	"net/http"
	"os"
	"runtime"
	"time"
)

func Http_init() {

	http.Handle("/", websocket.Handler(echoHandler))

	//no tls
	go http.ListenAndServe(":9003", nil)

	//tls addon test
	go http.ListenAndServeTLS(":9004", "sign.pem", "ssl.key", nil)

	for true {
	}
}

var decrypted string
var privateKey, publicKey []byte

// 加密
func RsaEncrypt(origData []byte) ([]byte, error) {
	block, _ := pem.Decode(publicKey)
	if block == nil {
		return nil, errors.New("public key error")
	}
	pubInterface, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		return nil, err
	}
	pub := pubInterface.(*rsa.PublicKey)
	return rsa.EncryptPKCS1v15(rand.Reader, pub, origData)
}

// 解密
func RsaDecrypt(ciphertext []byte) ([]byte, error) {
	block, _ := pem.Decode(privateKey)
	if block == nil {
		return nil, errors.New("private key error!")
	}
	priv, err := x509.ParsePKCS1PrivateKey(block.Bytes)
	if err != nil {
		return nil, err
	}
	return rsa.DecryptPKCS1v15(rand.Reader, priv, ciphertext)
}

//RSA公钥私钥产生
func GenRsaKey(bits int) error {
	// 生成私钥文件
	//var privateKey *rsa.PrivateKey
	//privateKey =new(rsa.PrivateKey)
	PKEY, err := rsa.GenerateKey(rand.Reader, bits)
	if err != nil {
		return err
	}
	derStream := x509.MarshalPKCS1PrivateKey(PKEY)
	block := &pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: derStream,
	}
	file, err := os.Create("d:\\private.pem")
	if err != nil {
		return err
	}
	privateKey = pem.EncodeToMemory(block)
	fmt.Println(privateKey)

	err = pem.Encode(file, block)
	if err != nil {
		return err
	}
	// 生成公钥文件
	PUKEY := &PKEY.PublicKey
	derPkix, err := x509.MarshalPKIXPublicKey(PUKEY)
	if err != nil {
		return err
	}
	block = &pem.Block{
		Type:  "PUBLIC KEY",
		Bytes: derPkix,
	}
	file, err = os.Create("d:\\public.pem")
	if err != nil {
		return err
	}
	publicKey = pem.EncodeToMemory(block)
	err = pem.Encode(file, block)
	if err != nil {
		return err
	}
	return nil
}

func load_keys() {
	var err error
	// flag.StringVar(&decrypted, "d", "", "加密过的数据")
	// flag.Parse()
	publicKey, err = ioutil.ReadFile("d:\\public.pem")
	if err != nil {
		os.Exit(-1)
	}
	privateKey, err = ioutil.ReadFile("d:\\private.pem")
	if err != nil {
		os.Exit(-1)
	}
}

func GenPPL(ws *websocket.Conn) {

	var nopass bool = true
	var tmp uint64
	fmt.Printf("ppl in  \r\n\n")
	for nopass {
		nopass = false
		tmp = rand2.Uint64()
		for _, v := range Public.LoginUser {
			//fmt.Printf("member %q,,%q  \r\nnn", n, v)
			if v.PplId == tmp {
				nopass = true
			}
		}
	}
	//fmt.Printf(" %s PPL is %d \r\n\n\n", ws, tmp)
	Public.LoginUser[ws].PplId = tmp
}

type cmd struct {
	Cmd  string      `json:"cmd"`
	Data interface{} `json:"data"`
}

func HB(ws *websocket.Conn) {
	Senders := new(Public.Senders)
	var send cmd
	send.Cmd = "HB"
	send.Data = ""
	rec, _ := json.Marshal(send)
	data_tmp := string(rec)
	fmt.Println("addr", ws.Request().RemoteAddr)
	for true {

		if Public.LoginUser[ws] != nil {
			Public.DB2Ret <- Senders
			Senders.Ws = ws
			Senders.Dat = data_tmp
			Public.LoginUser[ws].HBLife = Public.LoginUser[ws].HBLife + 1
			if Public.LoginUser[ws].HBLife > 10 {
				fmt.Println("closed by HB!! ")
				ws.Close()
			}
		}
		time.Sleep(time.Second * 10)
	}
}

func echoHandler(ws *websocket.Conn) {
	var err error
	var n int

	fmt.Println("echoHandler ws addr ", ws.Request().RemoteAddr)

	defer func() {
		if err := recover(); err != nil {
			strLog := "longweb:main recover error => " + fmt.Sprintln(err)
			os.Stdout.Write([]byte(strLog))
			log.Error(strLog)

			buf := make([]byte, 4096)
			n := runtime.Stack(buf, true)
			log.Error(string(buf[:n]))
			os.Stdout.Write(buf[:n])
		}
	}()

	defer ws.Close()
	go sender()
	go HB(ws)

	//bits := 1024
	//if err := GenRsaKey(bits); err != nil {
	//	log.Fatal("密钥文件生成失败！")
	//}
	//log.Println("密钥文件生成成功！")
	//
	//initData := "abcdefghijklmnopq"
	//init := []byte(initData)
	////load_keys()
	//
	//data, err := RsaEncrypt(init)
	//if err != nil {
	//	panic(err)
	//}
	//pre := time.Now()
	//origData, err := RsaDecrypt(data)
	//if err != nil {
	//	panic(err)
	//}
	//now := time.Now()
	//fmt.Println(now.Sub(pre))
	//fmt.Println(string(origData))

	//register current dialog
	if _, ok := Public.LoginUser[ws]; !ok {
		//fmt.Println("a,ok,ww",a,ok,Public.LoginUser[ws])
		Public.LoginUser[ws] = new(Public.LoginType)
		Public.LoginUser[ws].Name = "匿名"
		Public.LoginUser[ws].Handle = ws
		go GenPPL(ws)

	}
	fmt.Println("users %q：", Public.LoginUser)

	msg := make([]byte, 1024)

	for true {

		n, err = ws.Read(msg)
		if err != nil {
			fmt.Printf("errss %s\n", err)
			fmt.Print(err)
			switch {
			case err == io.EOF:
				delete(Public.LoginUser, ws)
				fmt.Println("\n\n\nusers %q：\n\n\n", Public.LoginUser)
				ws.Close()
				goto out
			default:
				log.Fatal(err)
			}
		}

		fmt.Printf("Receive: %s\n", msg[:n])
		DB_SAL.ReqProcess(ws, string(msg[:n]))

	}
out:
}

func sender() {

	for {
		rec := <-Public.DB2Ret
		fmt.Printf("sender to send :%s\r\n", rec)
		_, err := rec.Ws.Write([]byte(rec.Dat))
		if err != nil {
			fmt.Printf("sender err %s\n", err)
			//fmt.Print(err)
			switch {
			case err == io.EOF:
				goto Exit
			default:
				goto Exit
				log.Fatal("Fatal Err: %s \r\n", err)
			}
		}
	}
Exit:
}
