package main

import (
    "net"
    "fmt"
    "strings"
    "encoding/xml"
    "io/ioutil"
    "os"
	"path/filepath"
	"sync"
	"usr"
	"encoding/json"
)

const(
	Msgtype_log="0"
	Msgtype_state="1"
	Msgtype_users="2"
	Msgtype_info="3"
	Msgtype_info_1="4"
	Msgtype_line="5"
)

const(
	TcpType="tcp4"
	LocalAddr="127.0.0.1:8880"
)

type OutError struct{
	string
}

func(e *OutError) Error() string{
	return e.string
}

var userArrs usr.Users;
var connMaps map[string]net.Conn
 
func main() {
	
	fmt.Println("初始化用户信息")
	
	connMaps=map[string]net.Conn{}
		
	//获取用户列表存储文件
	filepath:=getXmlFile()
	
	//读取所有用户列表，用户验证登录和推送用户朋友信息
	users,err:=getXmlContent(filepath)
	
	if err!=nil{
		return
	}	
	
	userArrs=users
	userArrs.RMutex=new(sync.RWMutex)
	
	fmt.Println("创建监听服务对象")
	
    //创建一个TCP服务端
    tcpaddr,_:= net.ResolveTCPAddr(TcpType, LocalAddr);
    
    //创建一个监听对象
	listenTcp,err:=net.ListenTCP(TcpType, tcpaddr);
	
	//监听对象创建成功
	if(err==nil){
		
		fmt.Println("监听对象创建成功")
			
		//通过循环的方式等待用户的连接
		for{
				
			fmt.Println("等待用户连接......")
			
			//阻塞状态,持续等待用户的连接
			conn,conerr:=listenTcp.Accept()
					
			//如果有用户成功连接
			if(conerr==nil){		
				addr:=conn.RemoteAddr().String()
				if _,ok:=connMaps[addr];!ok{
					connMaps[addr]=conn
					
					fmt.Println(conn.RemoteAddr().String()+" 连接成功")
					//异步方式读取客户端信息
					go ReadFromTCP(conn);
				}
				
			}		
		}
	}
	
}

//读取客户端信息
func ReadFromTCP(conn net.Conn){	       
        var condion bool=true
        for condion{
	        
	        data,ln:=read(conn)
        
	        if ln > 0 {
				logValue:=string(data[0:1]) 
		    	valueByte:=data[1:ln]
		    	
				switch logValue{
					case Msgtype_log:
						if user,ok:=checkUser(valueByte);ok{
							
							fmt.Println("验证成功......")
							
							//验证成功，返回successful
							succValue:="1successful"		
							sendToCLient(conn,[]byte(succValue))	
							
							//更新用户信息
							userArrs.UpUser(user,conn,true)	
							
							//获取用户好友列
							users:=userArrs.GetFriends(user)
							fmt.Println(users)
							//发送好友列表
							userByte,_:=WriteToByte(users)
							userContent:=append([]byte("2"),userByte...)							
							sendToCLient(conn,userContent)
							
							//将该用户的登录状态发送给已经登录的好友
							for i:=0;i<len(users);i++{
								if users[i].Linestate{
									if userConn,ok:=connMaps[users[i].Address];ok{
										svalue:=fmt.Sprintf("%s%s,%s,%s", Msgtype_line,user.Name,user.Address,"1")
										fmt.Println(svalue)
										sendToCLient(userConn,[]byte(svalue))
									}
								}
							}								
												
						}
					
					case Msgtype_state:
						condion=false
					case Msgtype_users:
						condion=false
					default:
						condion=false								
					}			
	        }else{
		        condion=false
	        }
        }
}

func sendToCLient(conn net.Conn,value []byte){
	conn.Write(value)
}

func checkUser(valueByte []byte)(*usr.User,bool){
	
	values:=string(valueByte)
	valuesArr:=strings.Split(values, ",")
	if len(valuesArr)==3{
		user:=userArrs.CheckUser(valuesArr[0],valuesArr[1])
		if user!=nil{		
			return user,true
		}						
	}		
	return nil,false
}

func read(conn net.Conn)([]byte,int){
	bufsize:=256
	buf:=make([]byte,bufsize)
	var bufs []byte
	ln:=0
	for{
		n,err:=conn.Read(buf)
		if n>0{
			ln+=n
		}
		
		bufs=append(bufs,buf...)
		
		if n<bufsize{
			break
		}
		
		if err!=nil{
			break
		}
	}
	
	return bufs,ln
}



func WriteToByte(users []usr.User)([]byte, error){
	return json.Marshal(users)
	
}

func ReadFromByte(b []byte)([]usr.User,error){
	var users []usr.User
	err:=json.Unmarshal(b,  &users)
	return users,err
}

//读取xml配置文件
func getXmlContent(xmlfile string)(users usr.Users,err error){
	content, inerr := ioutil.ReadFile(xmlfile)
	if inerr!=nil{
		err=inerr
		fmt.Println("open xml file err",inerr)
		return
	}
	
	var result usr.Users
	inerr = xml.Unmarshal(content, &result)
	if inerr!=nil{
		err=inerr
		fmt.Println("read xml file to Users err",inerr)
		return
	}

	return result,nil
}

//获取xml文件目录
func getXmlFile()string{
	dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	var filepath string;
	if err == nil {
		index:=strings.LastIndex(dir, "\\")
		filepath=string([]byte(dir)[:index+1])+"src\\xml\\users.xml"
	}
	
	return filepath
}
