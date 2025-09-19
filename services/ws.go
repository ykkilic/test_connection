package services

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"backend/db_utils"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"golang.org/x/crypto/ssh"
	"gorm.io/gorm"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

type SSHConfig struct {
	User     string
	Password string
	Host     string
	Port     string
}

func WSHandler(c *gin.Context, db *gorm.DB, sshConf SSHConfig, sessionID uint) {

	fmt.Println("HERERERERE")

	wsConn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Println("WebSocket upgrade error:", err)
		return
	}
	defer wsConn.Close()

	fmt.Println("SSH CONF: ", sshConf.User)

	sshClient, err := ssh.Dial("tcp", sshConf.Host+":"+sshConf.Port, &ssh.ClientConfig{
		User:            sshConf.User,
		Auth:            []ssh.AuthMethod{ssh.Password(sshConf.Password)},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	})
	if err != nil {
		log.Println("SSH connection error:", err)
		return
	}
	defer sshClient.Close()

	session, err := sshClient.NewSession()
	if err != nil {
		log.Println("SSH session error:", err)
		return
	}
	defer session.Close()

	modes := ssh.TerminalModes{
		ssh.ECHO:          1,
		ssh.TTY_OP_ISPEED: 14400,
		ssh.TTY_OP_OSPEED: 14400,
	}
	if err := session.RequestPty("xterm", 80, 40, modes); err != nil {
		log.Println("PTY request error:", err)
		return
	}

	stdinPipe, _ := session.StdinPipe()
	stdoutPipe, _ := session.StdoutPipe()

	if err := session.Shell(); err != nil {
		log.Println("SSH shell error:", err)
		return
	}

	// SSH stdout -> WebSocket + DB kaydı
	go func() {
		buf := make([]byte, 1024)
		for {
			n, err := stdoutPipe.Read(buf)
			if err != nil {
				break
			}
			data := buf[:n]
			wsConn.WriteMessage(websocket.TextMessage, data)

			// DB’ye kaydet
			event := db_utils.TerminalEvent{
				SessionID: sessionID,
				Type:      "stdout",
				Data:      string(data),
				Timestamp: time.Now().Unix(),
			}
			db.Create(&event)
		}
	}()

	// WebSocket -> SSH stdin + DB kaydı
	for {
		_, msg, err := wsConn.ReadMessage()
		if err != nil {
			break
		}
		stdinPipe.Write(msg)

		// DB’ye kaydet
		event := db_utils.TerminalEvent{
			SessionID: sessionID,
			Type:      "stdin",
			Data:      string(msg),
			Timestamp: time.Now().Unix(),
		}
		db.Create(&event)
	}
}
