package netconn

import (
    "net"
    "fmt"
    "encoding/json"
    "usr"
    "strings"
    "sync"
)

const(
	TcpType="tcp4"
	UdpType="udp"
)

const(
	Msgtype_log="0"
	Msgtype_state="1"
	Msgtype_users="2"
	Msgtype_info="3"
	Msgtype_info_1="4"
	Msgtype_line="5"
)

type UdpConns struct{
	rMutex *sync.RWMutex
	connMaps map[string] *net.UDPConn
}

func(u *UdpConns) add(address string,conn *net.UDPConn){
	
	u.rMutex.Lock()
	
	_,ok:=u.connMaps[address]
	if !ok{
		u.connMaps[address]=conn
	}

	u.rMutex.Unlock()
}

func(u *UdpConns) get(address string)*net.UDPConn{
	
	u.rMutex.RLock()
	
	defer u.rMutex.RUnlock()
	
	return u.connMaps[address]
}

func(u *UdpConns) del(address string){
	
	u.rMutex.RLock()
	
	defer u.rMutex.RUnlock()
	
	delete(u.connMaps,address)
}

var tcpConn net.Conn
var udpConn *net.UDPConn
var localAddress,targetAddress string


var FunWin func(interface{})
var FunUpList func(name,addr,state string)
var FunSetInfo func(name,info string)
var FunSetSingleInfo func(usrname,info string)

func InitAddress(ladr,tadr string){
	localAddress=ladr
	targetAddress=tadr
}

func ConnToTCP()(net.Conn,error){
	//创建请求服务器
    tcpaddrserver, _ := net.ResolveTCPAddr(TcpType, targetAddress);
    
    //本机客户端地址
    tcpaddrlocal, _ := net.ResolveTCPAddr(TcpType, localAddress);

    //连接请求的服务器
    dialTcp, err := net.DialTCP(TcpType, tcpaddrlocal, tcpaddrserver);
    
    if err==nil{
	    tcpConn=dialTcp
    }
        
	return dialTcp, err
}

//读取数据
func ReadFromTCP(conn net.Conn){
	var condion bool=true
    for condion{
		data,ln:=ReadBuffer(conn)		
	    if ln > 0 {
	    	logValue:=string(data[0:1]) 
	    	value:=data[1:ln]
			switch logValue{
				case Msgtype_state:
					valueStr:=string(value)
					if valueStr!="successful"{
						condion=false
					}
				case Msgtype_users:
					users,err:=ReadFromByte(value)
					if err==nil{
						FunWin(users)
					}
				case Msgtype_line:
					lineValue:=string(value) 
					lineSplit:=strings.Split(lineValue, ",")
					if len(lineSplit)==3{
						FunUpList(lineSplit[0],lineSplit[1],lineSplit[2])
					}						
				default:
					condion=false									
			}
		}	    	
	}
    
    conn.Close()
}

func ReadTCPAsync(conn net.Conn){
	go ReadFromTCP(conn)
}

//发送数据
func WriteToTCP(conn net.Conn,value []byte){
	conn.Write(value)
}

//启动UDP监听
func ListenUDP(){
	//创建请求服务器
	udpaddrserver, _ := net.ResolveUDPAddr(UdpType, localAddress);
	
	conn, err := net.ListenUDP(UdpType, udpaddrserver)
	
	if err==nil{
		udpConn=conn		
		go ReadFromUDP(udpConn)
	}
}

//读取UDP数据
func ReadFromUDP(conn *net.UDPConn){
	var condion bool=true
    for condion{
		data,ln:=ReadBufferUDP(conn)
	    if ln > 0 {
			logValue:=string(data[0:1]) 
	    	value:=data[1:ln]
	    	switch logValue{
				case Msgtype_info:
					valueStr:=string(value)
					valueArr:=strings.Split(valueStr, ",")
					if len(valueArr)==2{
						FunSetInfo(valueArr[0],valueArr[1])
					}
				case Msgtype_info_1:
					valueStr:=string(value)
					valueArr:=strings.Split(valueStr, ",")
					if len(valueArr)==2{
						go FunSetSingleInfo(valueArr[0],valueArr[1])
					}								
				default:
					fmt.Println("value error："+string(value))									
				}		
	    }else{
		    condion=false
		}
	}
    
    conn.Close()	
}

//发送UDP数据
func WriteToUDP(info,addr string){
	udpaddr, _ := net.ResolveUDPAddr(UdpType, addr);
	udpConn.WriteTo([]byte(info),udpaddr)
}

//读取TCP缓存数据
func ReadBuffer(conn net.Conn)([]byte,int){
		
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

//读取UDP缓存数据
func ReadBufferUDP(conn *net.UDPConn)([]byte,int){		
	bufsize:=256
	buf:=make([]byte,bufsize)
	var bufs []byte
	ln:=0
	for{
		n,err,_:=conn.ReadFromUDP(buf)
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
	
	return buf,ln
}


func setUserState(name,linestate string){
	
}

//JSON序列化
func WriteToByte(users []usr.User)([]byte, error){
	return json.Marshal(users)
	
}

//反序列化
func ReadFromByte(b []byte)([]usr.User,error){
	var users []usr.User
	err:=json.Unmarshal(b,  &users)
	return users,err
}




