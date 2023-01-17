package service

import "fmt"

type DataCellv2 interface {
	GetName() string
}

// dusSir 可以与string互相转换
type duSir string

// duSir，重写了GetName方法，实现了DataCellv2接口，使duSir和DataCellv2画上等号
func (d duSir) GetName() string {
	return "duSir"
}

type dataSelectorv2 struct {
	GenericData DataCellv2
	FilterName  string
	Limit       int
	Page        int
}

func (d *dataSelectorv2) aaa() {
	fmt.Println("aaa," + d.GenericData.GetName())
	fmt.Println(d.FilterName, d.Limit, d.Page)
}

func (d *dataSelectorv2) bbb() {
	fmt.Println("bbb," + d.GenericData.GetName())
	fmt.Println(d.FilterName, d.Limit, d.Page)
}

var Ds dataSelectorv2

func runDemo() {
	a := "ttt"
	fmt.Println(a)
	fmt.Println(duSir(a))

	selectableData := &dataSelectorv2{
		GenericData: duSir(a),
		FilterName:  "nnnnnnn",
		Limit:       10,
		Page:        1,
	}

	fmt.Println(selectableData)
}
