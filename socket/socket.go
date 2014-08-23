package socket

import (
	"bytes"
	"code.google.com/p/goprotobuf/proto"
	"encoding/binary"
	"fmt"
	"net"
	"tcp-proto/client"
	"time"
)

var conn *net.TCPConn

const (
	RECV_BUF_LEN = 1024
)

func Start(c *net.TCPConn) {
	conn = c
	go sendHeartBeat()
	go sendTestData()
	readLoop()
}

func sendHeartBeat() {
	msg := new(client.SCS100002)
	for {
		time.Sleep(5 * time.Second)
		SendMsg(100002, msg)
	}
}

func sendTestData() {
	var num int64 = 0
	msg := new(client.SCS100001)

	for {
		msg.Num = proto.Int64(num)
		msg.Msg = new(client.TestData)
		msg.Msg.Time = proto.Int64(int64(time.Now().Second()))
		msg.Msg.State = proto.Bool(num%2 == 0)
		SendMsg(100001, msg)
		num++
		time.Sleep(time.Second)
	}
}

func readLoop() {
	defer conn.Close() // connection already closed by client
	var request []byte
	for {

		request = make([]byte, RECV_BUF_LEN) // clear last read content
		len, err := conn.Read(request)

		if err != nil {
			fmt.Println(err.Error())
			break
		}
		data := request[:len]

		buf := bytes.NewBuffer(data)
		var l uint32
		var head uint32
		errR := binary.Read(buf, binary.BigEndian, &l)
		checkError(errR)

		errR = binary.Read(buf, binary.BigEndian, &head)
		checkError(errR)

		var b []byte = make([]byte, l)
		errR = binary.Read(buf, binary.BigEndian, &b)
		checkError(errR)
		if head == 100001 {
			msg := new(client.SSC100001)
			msgErr := proto.Unmarshal(b, msg)
			checkError(msgErr)
			fmt.Println("接收解析", head, "->", msg)
		} else {
			fmt.Println("接收", head, "->", string(data))
		}
		// errR = binary.Read(buf, binary.BigEndian, &attempts)
		// checkError(errR)
	}
}

func SendMsg(head uint32, msg proto.Message) {
	data, err := proto.Marshal(msg)
	checkError(err)

	var buf bytes.Buffer
	var l uint32 = uint32(len(data))
	err = binary.Write(&buf, binary.BigEndian, l)
	checkError(err)

	err = binary.Write(&buf, binary.BigEndian, head)
	checkError(err)

	err = binary.Write(&buf, binary.BigEndian, data)
	checkError(err)

	fmt.Println("发送", head, "->", buf.String())
	conn.Write(buf.Bytes())
}

func checkError(err error) {
	if err != nil {
		panic(err.Error())
	}
}
