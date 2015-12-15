// ClientMsgProcess
package FDServer

import (
	"fmt"
	"io"
	"net"
)

type Client struct {
	AccountID     int
	Name          string
	Icon          int
	Level         int
	Area          int
	Viplevel      int
	IsLive        bool
	RecMsg        chan string
	AckMsg        chan string
	Server        *Server
	Client_Socket net.Conn
	UserRelation
}

type UserRelation struct {
	FriendInfo map[uint64]string
	IgnorID    map[uint64]string
}

func (this *Client) ClientMsgProcess() {
	fmt.Println("In ClientMsgProcess...")
	select {
	case buf := <-this.RecMsg:
		{
			data := []byte(buf)

			buf, Head := MsgJsonDecode(data)
			if buf == nil || Head == nil {
				return
			}
			this.ProcessMsg(buf, Head)

			fmt.Println("In <-this.RecMsg...")
			fmt.Println(Head, ":", buf)
			fmt.Println("Out <-this.RecMsg...")
		}
	case data := <-this.AckMsg:
		{
			fmt.Println("In <-this.AckMsg...")
			fmt.Println("data:", data)
			buf := []byte(data)
			n, err := this.Client_Socket.Write(buf)
			if err != nil {
				if err != io.EOF {
					return
				}
			}
			fmt.Println("Ack len: ", n)
			if n != 0 {
				this.Server.AckDataSize += uint64(n)
			}
			fmt.Println("In <-this.AckMsg...")
		}
	}
}

func (this *Client) SendDataToChann(buf []byte) {
	data := string(buf)
	this.AckMsg <- data
	this.ClientMsgProcess()
}

func (this *Client) SendToSlef(msg interface{}) {
	buf, err := MsgJsonEncode(msg)
	if err {
		return
	}
	this.SendDataToChann(buf)
}

func (this *Client) SendToOtherByAccount(account uint64, msg interface{}) {
	conn := this.Server.AccountList[account]
	cl := this.Server.UserList[conn]

	buf, err := MsgJsonEncode(msg)
	if err {
		return
	}

	cl.SendDataToChann(buf)
}

func (this *Client) SendToOtherByName(name string, msg interface{}) {
	conn := this.Server.NameList[name]
	cl := this.Server.UserList[conn]

	buf, err := MsgJsonEncode(msg)
	if err {
		return
	}

	cl.SendDataToChann(buf)
}

func (this *Client) ProcessMsg(data []byte, Head *MsgHead) {
	//length 取出后吧Msg的length设置为0然后从新获取然后作比较
	switch Head.Id {
	case E_NONE:
		{
			testMsg := new(TestMsg)
			MsgByteToJson(data, testMsg)

		}
	case E_HEARTBEAT:
		{
			heartBeat := new(HeartBeat)
			MsgByteToJson(data, heartBeat)
		}
	default:
		{

		}
	}
}
