package win

import (
 	"github.com/lxn/walk"
 	"fmt"
)

type EnvItem struct {
    name string
    address string     //地址
	linestate string    //在线状态
}

//列表数据模型
type EnvModel struct {
    //继承ListModelBase
    walk.ListModelBase
 
    //环境变量数集合
    items []EnvItem
}

//列表数据模型的工厂方法
func NewEnvModel(envitems []EnvItem) *EnvModel { 
    m:=&EnvModel{items:envitems}
    return m
}

//初始化model
func CreateModel() *EnvModel { 
    m:=&EnvModel{items:make([]EnvItem,0)}
    return m
}

//列表的系统回调方法：获得listbox的数据长度
func (m *EnvModel) ItemCount() int {
    return len(m.items)
}
 
//列表的系统回调方法：根据序号获得数据
func (m *EnvModel) Value(index int) interface{} {
    return fmt.Sprintf("%s(%s)",m.items[index].name,m.items[index].linestate)
}
