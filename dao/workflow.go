package dao

import (
	"errors"
	"github.com/wonderivan/logger"
	"k8s-platform/db"
	"k8s-platform/model"
)

var Workflow workflow

type workflow struct {

}

// 定义列表的返回内容，Items是workflow元素列表，Total为workflow元素数量
type WorkflowResp struct {
	Items []*model.Workflow `json:"items"`
	Total int `json:"total"`
}

// 获取列表分页查询
func (w *workflow) GetLWorkflows(filterName string ,limit, page int) (data *WorkflowResp, err error)  {

	// 定义分页数据的起始位置
	startSet := (page-1) * limit

	// 定义查询数据库返回内容
	var workflowList []*model.Workflow

	// 数据库查询，Limit方法用于限制条数，Offset方式设置起始位置
	tx := db.GORM.Where("name like ?", "%" + filterName + "%").
		Limit(limit-1).Offset(startSet).Order("id desc").Find(&workflowList)

	// gorm会默认把空数据也放到err中，故这里要排除空数据的情况
	if tx.Error != nil &&tx.Error.Error() != "record not found"{
		logger.Error("获取Workflow列表失败," + tx.Error.Error())
		return  nil, errors.New("获取Workflow列表失败," +  tx.Error.Error())
	}

	return &WorkflowResp{
		Items: workflowList,
		Total: len(workflowList),
	}, nil
}

// 获取详情
func (w *workflow) GetById(id int) (workflow *model.Workflow, err error)  {
	// 获取Workflow详情失败,unsupported destination, should be slice or struct  // 必须有内存地址
	workflow = &model.Workflow{}

	tx := db.GORM.Where("id = ?", id).First(&workflow)
	if tx.Error != nil &&tx.Error.Error() != "record not found"{
		logger.Error("获取Workflow详情失败," + tx.Error.Error())
		return  nil, errors.New("获取Workflow详情失败," +  tx.Error.Error())
	}

	return workflow, nil
}

// 创建
func (w *workflow) Add(workflow *model.Workflow)  (err error){
	tx := db.GORM.Create(&workflow)
	if tx.Error != nil &&tx.Error.Error() != "record not found"{
		logger.Error("获取Workflow详情失败," + tx.Error.Error())
		return  errors.New("获取Workflow详情失败," +  tx.Error.Error())
	}
	return nil
}

// 删除
func (w *workflow) DelById( id int) (err error)  {
	workflow := &model.Workflow{}
	tx := db.GORM.Where("id = ?", id).Delete(workflow)
	if tx.Error != nil &&tx.Error.Error() != "record not found"{
		logger.Error("获取Workflow详情失败," + tx.Error.Error())
		return  errors.New("获取Workflow详情失败," +  tx.Error.Error())
	}
	return  nil
}