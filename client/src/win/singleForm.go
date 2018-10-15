package win

import (
 	"github.com/lxn/walk"
	. "github.com/lxn/walk/declarative"
	//"strings"
		"fmt"
		"netconn"
		//"net"
		//"usr"
)

var mapWins map[string] *SingleForm

func init(){
	mapWins=map[string] *SingleForm{}
}

type SingleForm struct{
	singleForm *walk.MainWindow
	txtMsgs *walk.TextEdit
	txtSend *walk.TextEdit
	btnSend *walk.PushButton
	
	localuserName string
	touserName string
	toAddr string	
}

func (s *SingleForm) InitInfo(localname,toname,toaddr string){
	s.localuserName=localname
	s.touserName=toname	
	s.toAddr=toaddr		
}

func (s *SingleForm) CreateForm(){
	s.singleForm=new(walk.MainWindow)
	singleForm:=MainWindow{
		Title:   "单聊",
		MinSize: Size{600, 400},
		Layout:  VBox{},
		AssignTo:&s.singleForm,		
		Children: []Widget{
			TextEdit{
				AssignTo: &s.txtMsgs, 
				MinSize: Size{600, 300},
			},
            HSplitter{
				MaxSize: Size{30, 300},
                Children: []Widget{
                    //按钮
                    PushButton{
                        StretchFactor:2,
                        Text: "发送",
                        OnClicked: s.sendClick,
                    },
 
                    TextEdit{
                        StretchFactor: 8,
                        AssignTo: &s.txtSend, 
                        ReadOnly:      false,
                    },
                },
            },
		},
	}
	
	singleForm.Create()
}

func (s *SingleForm) Run(){
	s.singleForm.Run()
}

func (s *SingleForm) Show(){
	s.singleForm.Show()
}

func (s *SingleForm) sendClick(){
	
	info:=fmt.Sprintf("%s%s,%s", netconn.Msgtype_info_1,s.localuserName,s.txtSend.Text())
	netconn.WriteToUDP(info,s.toAddr)

	s.SetInfo(s.localuserName,s.txtSend.Text())
}

func (s *SingleForm) SetInfo(name,content string){
	if s.txtMsgs.Text()==""{
		s.txtMsgs.SetText(fmt.Sprintf("%s:\r\n%s\r\n", name,content))
	}else{
		s.txtMsgs.SetText(fmt.Sprintf("%s\r\n%s say:\r\n%s\r\n", s.txtMsgs.Text(),name,content))
	} 
	
}

func (s *SingleForm) closing(canceled *bool, reason walk.CloseReason){
	delete(mapWins,s.touserName)	
}

