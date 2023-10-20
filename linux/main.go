package main

/*
#include <stdio.h>
#include <stdlib.h>
*/
import "C"

// "Cだけ分けて書く必要がある。"

import (
	"bytes"
	"fmt"
	"log"
	"os"
	"unsafe"

	"golang.org/x/crypto/ssh"
)

// disposeなどの処理はtypescript側に書く。

// type Miyuu_SSH struct {
// 	ClientConfig *ssh.ClientConfig
// }

// ssh clientのラッパーを提供します。
// GCに回収されるのを拒否するためです。
type ClientWrapper struct {
	Client *ssh.Client // Closeすること。
}

type ClientConfigWrapper struct {
	ClientConfig *ssh.ClientConfig
}

// malloc時に必要な大きさ。
const ClientConfigWrapper_Size = (C.ulong)(unsafe.Sizeof(ClientConfigWrapper{}))

// SessionのWrapper
type SessionWrapper struct {
	Session *ssh.Session
}

const SessionWrapper_Size = (C.ulong)(unsafe.Sizeof(SessionWrapper{}))

func (sw SessionWrapper) SessionRun(cmd string) {

	// 対話形式の場合は
	// session.Stdout = os.stdout
	// session.Stderr = os.Stderr
	// と紐づける。
	var b bytes.Buffer
    sw.Session.Stdout = &b
	sw.Session.Run(cmd)

	fmt.Println(b.String())
}

//export SessionRun
func SessionRun(ptr unsafe.Pointer, raw_cmd *byte, cmd_len int64) {
	sw := (*SessionWrapper)(ptr)
	cmd := unsafe.String(raw_cmd, cmd_len)
	sw.SessionRun(cmd)
}

func (sw SessionWrapper) SessionClose() {
	sw.Session.Close()
}

//export SessionClose
func SessionClose(ptr unsafe.Pointer) {
	sw := (*SessionWrapper)(ptr)

	sw.SessionClose()
}

// Wrapperが持っているのをcloseする。
// 名前はDisposeの方が良いかも
func (cw ClientWrapper) ClientClose() {
	// Close何回されても良いように作る。
	cw.Client.Close()
	// Disposeの方が良いかも。
	// Freeな処理
}

//export ClientClose
func ClientClose(ptr unsafe.Pointer) {
	cw := (*ClientWrapper)(ptr)
	cw.ClientClose()
}

// malloc時に必要な大きさ。
const ClientWrapper_Size = (C.ulong)(unsafe.Sizeof(ClientWrapper{}))

func main() {
	// サクッと、対話形式でやるとか。
	// pythonのmainみたいな感じと思えば良い。
	// 必要ないなら。os.exit(0)を返すなど最低限の処理だけ
	// 書けば容量は取らない。
	os.Exit(0)
}

// malloc時に必要な大きさ。
const ClientConfig_Size = (C.ulong)(unsafe.Sizeof(ssh.ClientConfig{}))

// C言語のmallocのラッパーを提供する
// ちゃんとfreeを読んで解放してください。
func Malloc(size C.ulong) unsafe.Pointer{
	ptr := C.malloc(size)
	return ptr
}

//export InitConfig
func InitConfig(raw_user *byte, user_len int64, raw_pass *byte, pass_len int64) unsafe.Pointer {

	user := unsafe.String(raw_user, user_len)
	pass := unsafe.String(raw_pass, pass_len)

	ptr := Malloc(ClientConfigWrapper_Size)
	clientConfigWrapper := (*ClientConfigWrapper)(ptr)
	// ptr := Malloc(ClientConfig_Size)
	// clientConfig := (*ssh.ClientConfig)(ptr)

	clientConfigWrapper.ClientConfig = &ssh.ClientConfig{
		User:            user,
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),   // https://github.com/golang/go/issues/19767
		Auth: []ssh.AuthMethod{
				ssh.Password(pass),
		},
	}

	// clientConfig.User = user
	// clientConfig.Auth = []ssh.AuthMethod{
	// 	ssh.Password(pass),
	// }

	return ptr
}

// sshサーバーに接続することにより、Clientを返します。
//export Connect
func Connect(clientConfigWrapperPtr unsafe.Pointer, addr_octet uint8, addr_octet2 uint8, addr_octet3 uint8, addr_octet4 uint8, port uint16) unsafe.Pointer{
// func Connect(clientConfigPtr unsafe.Pointer, addr_octet uint8, addr_octet2 uint8, addr_octet3 uint8, addr_octet4 uint8, port uint16) unsafe.Pointer{

	clientConfigWrapper := (*ClientConfigWrapper)(clientConfigWrapperPtr)

	ip_port := fmt.Sprintf(
		"%d.%d.%d.%d:%d",
		addr_octet,
		addr_octet2,
		addr_octet3,
		addr_octet4,
		port,
	)
	// いつ解放されるかわからんのでラッパー
	client, err := ssh.Dial("tcp", ip_port, clientConfigWrapper.ClientConfig)
	if err != nil {
		log.Println(err)
	}
	ptr := Malloc(ClientWrapper_Size)

	clientWrapper := (*ClientWrapper)(ptr)

	clientWrapper.Client = client;

	// client.Close()
	// 
	return ptr
}

//export NewSession
func NewSession(clientWrapperPtr unsafe.Pointer) unsafe.Pointer{
	clientWrapper := (*ClientWrapper)(clientWrapperPtr)

	session, err := clientWrapper.Client.NewSession()
	if err != nil {
		log.Fatal("abcd")
	}
	ptr := Malloc(SessionWrapper_Size)

	sessionWrapper := (*SessionWrapper)(ptr)

	sessionWrapper.Session = session;

	return ptr
}

//export Free
func Free(ptr unsafe.Pointer) {
	C.free(ptr)
}

// defer conn.Close()

// session, err := conn.NewSession()
// if err != nil {
// 		log.Println(err)
// }
// defer session.Close()

// //Check whoami
// var b bytes.Buffer
// session.Stdout = &b
// remote_command := "/usr/bin/whoami"
// if err := session.Run(remote_command); err != nil {
// 		log.Fatal("Failed to run: " + err.Error())
// }
// log.Println(remote_command + ":" + b.String())
