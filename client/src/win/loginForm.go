package win

import (
 	"github.com/lxn/walk"
	. "github.com/lxn/walk/declarative"
//	"strings"
		//"fmt"
)

type LoginForm struct{
	loginForm *walk.MainWindow
	txtName *walk.TextEdit
	txtPwd *walk.TextEdit
	txtIP *walk.TextEdit
	btnOK *walk.PushButton	
	btnCancle *walk.PushButton
	
	name string
	pwd string
	ipaddress string
}

func (l *LoginForm) CreateForm(){
	l.loginForm=new(walk.MainWindow)
	
	loginForm:=MainWindow{
		Title:   "登录",
		MinSize: Size{270, 190},
		Layout:  VBox{},
		AssignTo:&l.loginForm,

		Children: []Widget{
			HSplitter{
				MinSize: Size{270, 30},
				Children: []Widget{
					Label{StretchFactor: 2,Text: "账号:",},
					TextEdit{ AssignTo: &l.txtName,StretchFactor: 8,MaxLength: 20},
				},
			},
			HSplitter{
				MinSize: Size{270, 30},
				Children: []Widget{
					Label{StretchFactor: 2,Text: "密码:",},
					TextEdit{ AssignTo: &l.txtPwd,StretchFactor: 8,MaxLength: 20},
				},
			},
			HSplitter{
				MinSize: Size{270, 30},
				Children: []Widget{
					Label{StretchFactor: 2,Text: "IPAddress:",},
					TextEdit{ AssignTo: &l.txtIP,StretchFactor: 8,MaxLength: 20},
				},
			},
            HSplitter{
				MaxSize: Size{30, 300},
                Children: []Widget{
                    //按钮
                    PushButton{
                        StretchFactor:1,
                        Text: "确定",
                        OnClicked: l.DMine,
                    },
 
                   //按钮
                    PushButton{
                        StretchFactor:1,
                        Text: "取消",
                        OnClicked: l.Close,
                    },
                },
            },
		},
	}
	
	loginForm.Create()

}

func (l *LoginForm) Show()int{
	return l.loginForm.Run()
}

func (l *LoginForm) SetOwner(win *walk.MainWindow){
	l.loginForm.SetOwner(win)
}

func (l *LoginForm) GetName()string{
	return l.name
}

func (l *LoginForm) GetPwd()string{
	return l.pwd
}

func (l *LoginForm) GetIpAddress()string{
	return l.ipaddress
}

func (l *LoginForm) DMine(){

	if l.txtName.Text()==""||l.txtPwd.Text()==""{
		walk.MsgBox(l.loginForm, "提示", "账号或密码不能为空!", walk.MsgBoxIconInformation)
	}else{
		l.name=l.txtName.Text()
		l.pwd=l.txtPwd.Text()
		l.ipaddress=l.txtIP.Text()
		l.Close()
	}
}

func (l *LoginForm) Close(){
	l.loginForm.Close()
}





