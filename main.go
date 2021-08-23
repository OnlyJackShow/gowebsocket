/*
* @Author: Hifun
* @Date: 2020/1/6 17:21
 */
package main

import (
	"fmt"
	"github.com/gorilla/websocket"
	"net/http"
	"socket/bs"
	"time"
)


var (
	upgrader = websocket.Upgrader{
		// 允许跨域
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}
)



func wsHandler(w http.ResponseWriter, r *http.Request) {
	var (
		wsConn *websocket.Conn
		err    error
		data   []byte
		conn   *bs.Connection
	)
	//todo 通过Upgrade方法把http转化成websocket
	if wsConn, err = upgrader.Upgrade(w, r, nil); err != nil {
		return
	}
	if conn, err = bs.InitConnection(wsConn); err != nil {
		goto ERR
	}

	//todo 开启一个携程，创建一个心跳包
	go func() {
		var (
			err error
		)
		for {
			if err = conn.WriteMessage([]byte(fmt.Sprintf("%s----%s","心跳",time.Now().Format("2006-01-02 03:04:05")))); err != nil {
				return
			}
			time.Sleep(2 * time.Second)
		}
	}()

	//todo 读写数据操作
	for {
		if data, err = conn.ReadMessage(); err != nil {
			goto ERR
		}
		if err = conn.WriteMessage(data); err != nil {
			goto ERR
		}
	}

ERR:
	conn.Close()
}

func main() {
	http.HandleFunc("/ws", wsHandler)

	http.ListenAndServe(":7777", nil)
}
