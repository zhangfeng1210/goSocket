package usr

import(
	"sync"
	"net"
)

type Users struct{
	User []User
	RMutex *sync.RWMutex
}

type User struct{
	Name string       //用户名
	Pwd string        //密码
	Sex string       //性别
	Age int          //年龄
	Address string     //地址
	Linestate bool    //在线状态
	Friends Friends  //朋友列表
}

type Friends struct{
	Friend []string
}

func(u *Users) AddUser(user *User){
	u.RMutex.Lock()
	u.User=append(u.User,*user)
	u.RMutex.Unlock()
}

func(u *Users) UpUser(user *User,conn net.Conn,state bool){
	u.RMutex.Lock()
	user.Linestate=state
	user.Address=conn.RemoteAddr().String()
	u.RMutex.Unlock()
}

func(u *Users) GetUserByName(name string)*User{
	for i:=0;i<len(u.User);i++{
		if(u.User[i].Name==name){
			return &u.User[i]
		}
	}
	
	return nil
}

func(u *Users) CheckUser(name string,pwd string)(*User){
	for i:=0;i<len(u.User);i++{
		if(u.User[i].Name==name&&u.User[i].Pwd==pwd){
			return &u.User[i]
		}
	}
	
	return nil
}

func(u *Users) GetFriends(user *User)[]User{
	
	var userfriends []User
	
	friends:=user.Friends.Friend
	for i:=0;i<len(friends);i++{
		userfriend:=u.GetUserByName(friends[i])
		userfriends=append(userfriends,*userfriend)
	}
	
	return userfriends

}






