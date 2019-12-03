package chat

import (
	"crypto/rand"
	"encoding/binary"
	"github.com/gocql/gocql"
	"log"
	"strconv"
	"time"
)

type Chat struct {
	Name     string `form:"name" json:"name"`
	Time     int64  `form:"time" json:"time"`
	Chatroom string `form:"chatroom" json:"chatroom"`
	Comment  string `form:"comment" json:"comment"`
}

// Insert test data
func InsertData(session *gocql.Session, chatData *Chat) string {
	log.Printf("Insert test data....", chatData)
	if err := session.Query(`INSERT INTO chat (name,time,chatroom,comment) VALUES (?,?,?,?)`,
		chatData.Name,
		chatData.Time,
		chatData.Chatroom,
		chatData.Comment).Exec(); err != nil {
		log.Fatal(err)
	}
	log.Println("Insert test data done!")
	return "insert done"
}

//Get Insert data
func SelectTestData(session *gocql.Session, chatData *Chat) Chat {
	log.Println("Select insert test data....")
	var selectChatData Chat
	if err := session.Query(`SELECT name,time,chatroom,comment FROM chat where name = ?`,
		chatData.Name).Consistency(gocql.One).Scan(
		&selectChatData.Name,
		&selectChatData.Time,
		&selectChatData.Chatroom,
		&selectChatData.Comment); err != nil {
		log.Fatal(err)
	}
	log.Println("Insert Data:", selectChatData)
	log.Println("Select insert test data done!")
	return selectChatData
}

// list all chat
func AllSelectData(session *gocql.Session) []Chat {
	// list all chat
	log.Println("Select all table data...")
	var ChatData Chat
	selectAllChatData := []Chat{}
	iter := session.Query(`SELECT name,time,chatroom,comment FROM chat`).Iter()
	for iter.Scan(
		&ChatData.Name,
		&ChatData.Time,
		&ChatData.Chatroom,
		&ChatData.Comment) {
		log.Println("All Chat:", ChatData)
		selectAllChatData = append(selectAllChatData, ChatData)
	}
	if err := iter.Close(); err != nil {
		log.Fatal(err)
	}
	time.Sleep(1000)
	log.Println("Select all table data Done!")
	return selectAllChatData
}

// generate random test data
func GenerateChatData() Chat {
	now := time.Now()

	chatData := Chat{}
	randomString := random()
	chatData.Name = "oranie-" + randomString
	chatData.Time = now.UnixNano()
	//chatData.Chatroom = "game_room-" + randomString
	chatData.Chatroom = "game_room-" + "oranie"
	chatData.Comment = "test comment : " + now.String()

	return chatData
}

//generate random string: 4char
func random() string {
	var n uint64
	binary.Read(rand.Reader, binary.LittleEndian, &n)
	rand_string := strconv.FormatUint(n, 36)
	return rand_string[:4]
}
