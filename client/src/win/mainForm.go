package win

import (
 	"github.com/lxn/walk"
	. "github.com/lxn/walk/declarative"
	//"strings"
		"fmt"
		"netconn"
		//"net"
		"usr"
)

type MainForm struct{
	mainForm *walk.MainWindow
	listInfo *walk.ListBox
	txtMsgs *walk.TextEdit
	txtSend *walk.TextEdit
	btnSend *walk.PushButton	
	menuLogin *walk.Action
	menuCancel *walk.Action
	
	//listbox使用的数据
    model *EnvModel
    userName string
}

func (m *MainForm) CreateForm(){
	m.mainForm=new(walk.MainWindow)
	//m.model=NewEnvModel1()
	mainWindow:=MainWindow{
		Title:   "群聊",
		MinSize: Size{600, 400},
		Layout:  VBox{},
		AssignTo:&m.mainForm,
		//窗口菜单
        MenuItems: []MenuItem{
            //主菜单一
            Menu{
                Text: "操作",
 
                //菜单项
                Items: []MenuItem{
                    //菜单项一
                    Action{
                        AssignTo: &m.menuLogin,
                        Text:     "登录",
                        OnTriggered: m.menuLoginClick,
                    },
 
                    //分隔线
                    Separator{},
 
                    //菜单项二
                    Action{
                        //文本
                        Text: "退出",
 
                        //响应函数
                        OnTriggered: m.menuCancelClick,
                    },
                },
            },
        },
		
		Children: []Widget{
			HSplitter{
				 MinSize: Size{600, 300},
				Children: []Widget{
					ListBox{                        
						StretchFactor: 2,                        
						//赋值给myWindow.listBox                        
						AssignTo: &m.listInfo,                                             
						OnItemActivated: m.listInfoDbClick,               
					},
					TextEdit{ 
						StretchFactor: 8,
						MaxLength: 10,
						AssignTo: &m.txtMsgs,
					},
				},
			},
            HSplitter{
				MaxSize: Size{30, 300},
                Children: []Widget{
                    //按钮
                    PushButton{
                        StretchFactor:2,
                        Text: "发送",
                        OnClicked:m.sendClick,
                    },
 
                    TextEdit{
                        StretchFactor: 8,
                        AssignTo:      &m.txtSend,
                        //ReadOnly:      false,
                    },
                },
            },
		},
	}
	
	mainWindow.Create()
}

func (m *MainForm) Run(){
	m.mainForm.Run()
}

func (m *MainForm) Show(){
	m.mainForm.Show()
}

func (m *MainForm) menuLoginClick(){
	loginForm:=LoginForm{}
	loginForm.CreateForm()
	loginForm.SetOwner(m.mainForm)
	result:=loginForm.Show()
	
	if result==0{
		m.userName=loginForm.GetName()
		m.mainForm.SetTitle("当前用户:"+m.userName)
		
		targetAddr:="127.0.0.1:8880"
		netconn.InitAddress(loginForm.GetIpAddress(),targetAddr)
		conn,err:=netconn.ConnToTCP()
		
		if err==nil{				
			logValue:=fmt.Sprintf("%s%s,%s,%s", netconn.Msgtype_log,loginForm.GetName(),loginForm.GetPwd(),loginForm.GetIpAddress())
			fmt.Println(logValue)
			netconn.WriteToTCP(conn, []byte(logValue))
			
			netconn.FunWin=m.setListInfo
			netconn.FunUpList=m.updateList
			netconn.FunSetInfo=m.setInfo
			netconn.FunSetSingleInfo=m.runSingleWin
			netconn.ReadTCPAsync(conn)
			netconn.ListenUDP()
		}
	}
}

func (m *MainForm) menuCancelClick(){
	m.mainForm.Close()
}

func (m *MainForm) listInfoDbClick() {
	curModel:=m.model.items[m.listInfo.CurrentIndex()]
	if curModel.linestate=="在线"{
		tousrname:=curModel.name
		
		if win,ok:=mapWins[tousrname];ok{
			win.Show()
		}else{
			singleForm:=SingleForm{}
			singleForm.InitInfo(m.userName,tousrname, curModel.address)
			mapWins[tousrname]=&singleForm		
				
			singleForm.CreateForm()
			singleForm.Run()	
		}
	}		
}

func (m *MainForm) setListInfo(obj interface{}){
	if users,ok:=obj.([]usr.User);ok{
		
		fmt.Println("user ok")
		lenUsers:=len(users)
		m.model=&EnvModel{items:make([]EnvItem,lenUsers)}
		
		for i:=0;i<lenUsers;i++{			
			listItem:=EnvItem{users[i].Name,users[i].Address,getLineValue(users[i].Linestate)}
			m.model.items[i]=listItem
		}
		
		m.listInfo.SetModel(m.model)
	}

}

func (m *MainForm) updateList(name,addr,state string) {	
	for i:=0;i<len(m.model.items);i++{
		if m.model.items[i].name==name{
			m.model.items[i].address=addr
			m.model.items[i].linestate=getLineValue_1(state)
			m.listInfo.SetModel(m.model)
			break
		}
	}
		
}

func (m *MainForm) sendClick(){
	for i:=0;i<m.model.ItemCount();i++{
		if m.model.items[i].linestate=="在线"{
			info:=fmt.Sprintf("%s%s,%s", netconn.Msgtype_info,m.userName,m.txtSend.Text())
			netconn.WriteToUDP(info,m.model.items[i].address)
		}
	}

	m.setInfo(m.userName,m.txtSend.Text())
}

func (m *MainForm) setInfo(name,content string){
	if m.txtMsgs.Text()==""{
		m.txtMsgs.SetText(fmt.Sprintf("%s:\r\n%s\r\n", name,content))
	}else{
		m.txtMsgs.SetText(fmt.Sprintf("%s\r\n%s say:\r\n%s\r\n", m.txtMsgs.Text(),name,content))
	} 
	
}

func (m *MainForm) getAddr(usrname string)string{
	for i:=0;i<m.model.ItemCount();i++{
		if m.model.items[i].name==usrname{
			return m.model.items[i].address
		}
	}
	fmt.Println(m.model.items)
	return ""
}


func (m *MainForm) runSingleWin(usrname,content string){
	if win,ok:=mapWins[usrname];ok{
		walk.MsgBox(m.mainForm, usrname, content, walk.MsgBoxIconInformation)
		win.SetInfo(usrname, content)
	}else{
			
			addr:=m.getAddr(usrname)
								
			singleForm:=SingleForm{}
			singleForm.InitInfo(m.userName,usrname, addr)
			mapWins[usrname]=&singleForm
			
			singleForm.CreateForm()
			singleForm.SetInfo(usrname, content)
			singleForm.Run()						
	}
}


func getLineValue(linestate bool)string{
	var lineValue string="离线"
	if linestate{
		lineValue="在线"
	}
	
	return lineValue
}

func getLineValue_1(linestate string)string{
	var lineValue string="离线"
	if linestate=="1"{
		lineValue="在线"
	}
	
	return lineValue
}



