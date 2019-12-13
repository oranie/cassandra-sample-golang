package chat

import (
	"crypto/rand"
	"encoding/binary"
	"log"
	"strconv"
	"time"

	"github.com/gocql/gocql"
)

type Comment struct {
	Name     string `form:"name" json:"name"`
	Time     int64  `form:"time" json:"time"`
	Chatroom string `form:"chatroom" json:"chatroom"`
	Comment  string `form:"comment" json:"comment"`
}

type Comments struct {
	Response []Comment `form:"name" json:"response"`
}

// Insert test data
func InsertData(session *gocql.Session, chatData *Comment) *Comment {
	log.Printf("Insert test data....%v %v", chatData, chatData)
	if err := session.Query(`INSERT INTO chat (name,time,chatroom,comment) VALUES (?,?,?,?)`,
		chatData.Name,
		chatData.Time,
		chatData.Chatroom,
		chatData.Comment+time.Now().String()).Exec(); err != nil {
		log.Println(err)
	}
	log.Println("Insert test data done!")
	return chatData
}

//Get Insert data
func SelectTestData(session *gocql.Session, chatData *Comment) Comment {
	log.Println("Select insert test data....")
	var selectChatData Comment
	if err := session.Query(`SELECT name,time,chatroom,comment FROM chat where name = ?`,
		chatData.Name).Consistency(gocql.One).Scan(
		&selectChatData.Name,
		&selectChatData.Time,
		&selectChatData.Chatroom,
		&selectChatData.Comment); err != nil {
		log.Println(err)
	}
	log.Println("Insert Data:", selectChatData)
	log.Println("Select insert test data done!")
	return selectChatData
}

func ChatroomLatestData(session *gocql.Session, chatroom string) []Comment {
	// list all chat
	//log.Println("Select Latest table data...")
	var ChatData Comment
	selectAllChatData := []Comment{}
	iter := session.Query(`SELECT name,time,chatroom,comment FROM chat WHERE chatroom = ? LIMIT 25 ALLOW FILTERING  `, chatroom).Iter()
	for iter.Scan(
		&ChatData.Name,
		&ChatData.Time,
		&ChatData.Chatroom,
		&ChatData.Comment) {
		selectAllChatData = append(selectAllChatData, ChatData)
	}
	//log.Println("selectAllChatData : ", selectAllChatData)

	if err := iter.Close(); err != nil {
		log.Println(err)
	}

	log.Println("Select all table data Done!")
	return selectAllChatData
}

func ChatroomAllData(session *gocql.Session, chatroom string) []Comment {
	// list all chat
	//log.Println("Select all table CQL...")
	var ChatData Comment
	selectAllChatData := []Comment{}
	iter := session.Query(`SELECT name,time,chatroom,comment FROM chat WHERE chatroom = ? DESC ALLOW FILTERING `, chatroom).Iter()
	for iter.Scan(
		&ChatData.Name,
		&ChatData.Time,
		&ChatData.Chatroom,
		&ChatData.Comment) {
		selectAllChatData = append(selectAllChatData, ChatData)
	}
	//log.Println("selectAllChatData : ", selectAllChatData)

	if err := iter.Close(); err != nil {
		log.Println(err)
	}

	log.Println("Select all table data Done!")
	return selectAllChatData
}

// list all chat
func AllSelectData(session *gocql.Session) []Comment {
	// list all chat
	log.Println("Select all table data...")
	var ChatData Comment
	selectAllChatData := []Comment{}
	iter := session.Query(`SELECT name,time,chatroom,comment FROM chat`).Iter()
	for iter.Scan(
		&ChatData.Name,
		&ChatData.Time,
		&ChatData.Chatroom,
		&ChatData.Comment) {
		selectAllChatData = append(selectAllChatData, ChatData)
	}

	if err := iter.Close(); err != nil {
		log.Println(err)
	}

	log.Println("Select all table data Done!")
	return selectAllChatData
}

// generate random test data
func GenerateChatData() Comment {
	now := time.Now()

	chatData := Comment{}
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
	err := binary.Read(rand.Reader, binary.LittleEndian, &n)
	if err != nil {
		log.Println("generate random data error", err)
	}
	rand_string := strconv.FormatUint(n, 36)
	return rand_string[:4]
}
