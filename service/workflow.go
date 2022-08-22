package service

import (
	"k8s-platform/dao"
	"k8s-platform/model"
)

// workflow的类型分为三种，ClusterIp，NodePort，Ingress

var Workflow workflow

type workflow struct {}

type WorkflowCreate struct {
	Name string `json:"name"`
	Namespace string `json:"namespace"`
	Replicas int32 `json:"replicas"`
	Image string `json:"image"`
	Label map[string]string `json:"label"`
	Cpu string `json:"cpu"`
	Memory string `json:"memory"`
	ContainerPort int32 `json:"container_port"`
	HealthCheck bool `json:"health_check"`
	HealthPath string `json:"health_path"`
	Type string `json:"type"`
	Port int32 `json:"port"`
	NodePort int32 `json:"node_port"`
	Hosts map[string][]*HttpPath `json:"hosts"`
}

// 创建workflow
func (w *workflow) CreateWorkFlow(data *WorkflowCreate) (err error)  {
	// 定义ingress名字
	var ingressName string
	if data.Type == "Ingress"{
		ingressName = getIngressName(data.Name)
	}

	// workflow数据落库
	workflow := &model.Workflow{
		Name:       data.Name,
		Namespace:  data.Namespace,
		Replicas:   data.Replicas,
		Deployment: data.Name,
		Service:    getServiceName(data.Name),
		Ingress:    ingressName,
		Type:       data.Type,
	}
	err = dao.Workflow.Add(workflow)
	if err != nil{
		return err
	}

	// 创建k8s资源
	err = createWorkflowRes(data)
	if err != nil{
		return err
	}
	return
}

// 创建k8s资源 deployment service ingress
func createWorkflowRes(data *WorkflowCreate) (err error){
	// 创建deployment
	dc := &DeploymentCreate{
		Name:          data.Name,
		Namespace:     data.Namespace,
		Replicas:      data.Replicas,
		Image:         data.Image,
		Label:         data.Label,
		Cpu:           data.Cpu,
		Memory:        data.Memory,
		ContainerPort: data.ContainerPort,
		HealthCheck:   data.HealthCheck,
		HealthPath:    data.HealthPath,
	}
	err = Deployment.CreateDeployment(dc)
	if err != nil{
		return err
	}

	// 创建service
	var serviceType string
	if data.Type != "Ingress"{
		serviceType = data.Type
	}else{
		serviceType = "ClusterIP"
	}
	sc := &ServiceCreate{
		Name:          getServiceName(data.Name),
		Namespace:     data.Namespace,
		Type:          serviceType,
		ContainerPort: data.ContainerPort,
		Port:          data.Port,
		NodePort:      data.NodePort,
		Label:         data.Label,
	}
	/**
	区别：
		1. err的作用域不同，只作用if与内部
		2. 需要data的，使用第一种，需要抛出去
	 */
	if err := Servicev1.CreateService(sc); err!=nil{
		return err
	}

	// 创建ingress
	var ic *IngressCreate
	if data.Type == "Ingress"{
		ic = &IngressCreate{
			Name:      getIngressName(data.Name),
			Namespace: data.Namespace,
			Label:     data.Label,
			Hosts:     data.Hosts,
		}
	}
	err = Ingress.CreateIngress(ic)
	if err != nil{
		return err
	}

	return
}

// 删除workflow
func (w *workflow) DeleteById(id int) (err error)  {
	// 获取数据库数据
	workflow,err  := dao.Workflow.GetById(id)
	if err != nil{
		return err
	}
	// 删除k8s资源
	err = deleteWorkflowRes(workflow)
	if err != nil{
		return err
	}
	// 删除数据库数据
	err = dao.Workflow.DelById(id)
	if err != nil{
		return err
	}
	return
}

// 删除k8s资源 deployment service  ingress
func deleteWorkflowRes(workflow *model.Workflow) (err error)  {

	err = Deployment.DeleteDeployment(workflow.Name, workflow.Namespace)
	if err != nil{
		return err
	}

	err = Servicev1.DeleteService(getServiceName(workflow.Name), workflow.Namespace)
	if err != nil{
		return err
	}

	if workflow.Type == "Ingress" {
		err = Ingress.DeleteIngress(getIngressName(workflow.Name), workflow.Namespace)
		if err != nil{
			return err
		}
	}

	return nil
}


// 查看workflows
func (w *workflow) GetWorkflows(filterName, namespace string, page, limit int) (data *dao.WorkflowResp, err error)  {
	data, err = dao.Workflow.GetLWorkflows(filterName, namespace, page, limit)
	if err != nil{
		return nil, err
	}
	return data, nil
}

// 查看workflow详情
func (w *workflow) GetById(id int) (workflow *model.Workflow, err error)  {
	 workflow, err = dao.Workflow.GetById(id)
	if err != nil{
		return nil, err
	}
	return workflow, nil
}


//workflow名字转换成service名字，添加-svc后缀
func getServiceName(workflowName string) (serviceName string)  {
	return workflowName + "-svc"
}

//workflow名字转换成ingress名字，添加-ing后缀
func getIngressName(workflowName string) (ingressName string)  {
	return workflowName + "-ing"

}