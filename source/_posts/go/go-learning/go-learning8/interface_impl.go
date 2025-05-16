package main

import (
	"fmt"
	"sort"
)

type SortableStrings []string
//type SortableStrings [3]string

//下面三个方法的声明，SortableStrings类型就已经是一个sort.Interface接口类型的实现了
//func (self *SortableStrings) Len() int {
func (self *SortableStrings) Len() int {
	return len(self)
}

//func (self *SortableStrings) Less(i, j int) bool {
func (self *SortableStrings) Less(i, j int) bool {
	return self[i] < self[j]
}

//func (self *SortableStrings) Swap(i, j int) {
func (self *SortableStrings) Swap(i, j int) {
	self[i], self[j] = self[j], self[i]
}

//将一个接口类型嵌入到另一个接口类型,变成Sortable
type Sortable interface {
	sort.Interface
	Sort()
}

//自定义数据类型SortableStrings也可以实现接口类型Sortable,只需声明如下方法
//func (self *SortableStrings) Sort() {
func (self *SortableStrings) Sort() {
	sort.Sort(self)
}

func main() {
	_, ok := interface{}(SortableStrings{}).(sort.Interface)
	fmt.Printf("ok:%v\n",ok)
	
	_, ok2 := interface{}(SortableStrings{}).(Sortable)
	fmt.Printf("ok2:%v\n",ok2)
	
	_, ok3 := interface{}(&SortableStrings{}).(Sortable)
	fmt.Printf("ok3:%v\n",ok3)
	
	ss := SortableStrings{"2", "3", "1"}
	ss.Sort()
	fmt.Printf("Sortable strings: %v\n", ss)

}