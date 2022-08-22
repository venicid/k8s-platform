package service

import (
	"context"
	"errors"
	"fmt"
	"github.com/wonderivan/logger"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
	"k8s.io/apimachinery/pkg/util/json"
	"strconv"
	"time"
)

var Deployment deployment

type deployment struct {}

/**
// 获取pods列表，支持过滤、排序、分页
// 类型转换的方法 corev1.Pod --> DataCell, DataCell --> corev1.Pod
*/
func (d *deployment) toCells(deployments []appsv1.Deployment)  []DataCell {
	cells := make([]DataCell, len(deployments))
	for i := range deployments{
		// todo // 这里怎么转换的？？
		cells[i] = deploymentCell(deployments[i])
	}
	return cells
}

func (d *deployment) fromCells(cells []DataCell) []appsv1.Deployment  {
	deployments := make([]appsv1.Deployment, len(cells))
	for i := range cells{
		// cells[i].(podCell) 是将DataCell类型转换成podCell
		deployments[i] = appsv1.Deployment(cells[i].(deploymentCell))
	}
	return deployments

}

/**
增删改查
 */
// 定义列表的返回内容;Items是deployment元素列表，Total是元素数量
type DeploymentsResp struct {
	Total int `json:"total"`
	Items []appsv1.Deployment `json:"items"`
}

//定义DeployCreate结构体，用于创建deployment需要的参数属性的定义
type DeploymentCreate struct {
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
}

//定义DeploysNp类型，用于返回namespace中deployment的数量
type DeploysNp struct {
	Namespace string `json:"namespace"`
	DeploymentNum int `json:"deployment_num"`
}


// 获取deployment列表，支持过滤、排序、分页
func (d *deployment) GetDeployments(filterName, namespace string, limit , page int) (deploymentsResp *DeploymentsResp, err error)  {

	//context.TODO()用于声明一个空的context上下文，用于List方法内设置这个请求的超时（源码），这里 的常用用法
	//metav1.ListOptions{}用于过滤List数据，如使用label，field等
	//kubectl get services --all-namespaces --field-seletor metadata.namespace != default
	deploymentList, err := K8s.ClientSet.AppsV1().Deployments(namespace).List(context.TODO(), metav1.ListOptions{})
	if err != nil{
		// 打印日志给自己看，排错使用
		logger.Error("获取deployments列表失败，" + err.Error())
		// 返回给上一层，最终返回给前端，前端打印出来这个error
		return nil, errors.New("获取deployments列表失败，"+ err.Error())
	}

	// 实例化dataSelector结构体，组装数据
	selectTableData :=&dataSelector{
		GenericDataList: d.toCells(deploymentList.Items),
		DataSelect:      &DataSelectQuery{
			Filter:   &FilterQuery{Name: filterName},
			Paginate: &PaginateQuery{
				Limit: limit,
				Page:  page,
			},
		},
	}
	// 先过滤
	filtered := selectTableData.Filter()
	total := len(filtered.GenericDataList)
	// 再排序和分页
	data := filtered.Sort().Paginate()
	// 将DataCell类型转换成pod
	deployments := d.fromCells(data.GenericDataList)

	// 数据处理后的数据与原始数据的比较
	// 处理后的数据
	fmt.Println("处理后的数据")
	for _, deployment := range  deployments{
		fmt.Println(deployment.Name, deployment.CreationTimestamp.Time)
	}
	// 原始数据
	fmt.Println("原始数据")
	for _, deployment := range  deploymentList.Items{
		fmt.Println(deployment.Name, deployment.CreationTimestamp.Time)
	}

	return &DeploymentsResp{
		Total: total,
		Items: deployments,
	}, nil


}

//获取deployment详情
func(d *deployment) GetDeploymentDetail(deploymentName, namespace string) (deployment *appsv1.Deployment, err error) {
	deployment, err = K8s.ClientSet.AppsV1().Deployments(namespace).Get(context.TODO(), deploymentName, metav1.GetOptions{})
	if err != nil {
		logger.Error(errors.New("获取Deployment详情失败, " + err.Error()))
		return nil, errors.New("获取Deployment详情失败, " + err.Error())
	}

	return deployment, nil
}

//设置deployment副本数
func (d *deployment) ScaleDeployment(deploymentName, namespace string, scaleNum int) (replica int32, err error)  {

	// 获取autoscalingv1.Scale类型的对象，能点出当前的副本数
	scale, err := K8s.ClientSet.AppsV1().Deployments(namespace).GetScale(context.TODO(), deploymentName, metav1.GetOptions{})
	if err!= nil{
		logger.Error("获取Deployment副本数信息失败," + err.Error())
		return 0, errors.New("获取Deployment副本数信息失败," + err.Error())
	}

	// 修改副本数
	scale.Spec.Replicas = int32(scaleNum)

	// 更新副本数，传入scale对象
	newScale, err := K8s.ClientSet.AppsV1().Deployments(namespace).UpdateScale(context.TODO(), deploymentName, scale, metav1.UpdateOptions{})
	if err!= nil{
		logger.Error("更新Deployment副本数信息失败," + err.Error())
		return 0, errors.New("更新Deployment副本数信息失败," + err.Error())
	}

	return newScale.Spec.Replicas, nil
}

// 创建deployment,接收DeploymentCreate对象
func (d *deployment) CreateDeployment(data *DeploymentCreate)(err error)  {
	// 将data中的数据组装成appsv1.Deployment对象
	deployment := &appsv1.Deployment{

		// ObjectMeta 定义资源名称,命名空间，label
		ObjectMeta: metav1.ObjectMeta{
			Name:                       data.Name,
			Namespace:                  data.Namespace,
			Labels:                     data.Label,
		},

		//Spec中定义副本数、选择器、以及pod属性
		Spec:       appsv1.DeploymentSpec{
			Replicas:                &data.Replicas,
			Selector:                &metav1.LabelSelector{
				MatchLabels:      data.Label,
			},
			Template:                corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{  //定义pod名和标签
					Name:                     data.Name,
					Labels:                     data.Label,
				},
				Spec:       corev1.PodSpec{
					Containers:[]corev1.Container{
						{
							Name: data.Name,
							Image : data.Image,
							Ports: []corev1.ContainerPort{
								{
									Name: "http",
									Protocol: corev1.ProtocolTCP,
									ContainerPort: 80,
								},
							},
						},
					},

				},
			},
		},

		//Status定义资源的运行状态，这里由于是新建，传入空的appsv1.DeploymentStatus{}对象即可
		Status:     appsv1.DeploymentStatus{},
	}

	// 判断健康检查功能是否打开，若打开，则增加健康检查功能
	if data.HealthCheck {

		//设置第一个容器的ReadinessProbe，因为我们pod中只有一个容器，所以直接使用index 0即可 //若pod中有多个容器，则这里需要使用for循环去定义了
		deployment.Spec.Template.Spec.Containers[0].ReadinessProbe = &corev1.Probe{
			Handler:             corev1.Handler{
				HTTPGet: &corev1.HTTPGetAction{
					Path:        data.HealthPath,
					Port:        intstr.IntOrString{
						Type : 0,
						IntVal: data.ContainerPort,
					},
				},
			},
			InitialDelaySeconds: 5, //初始化等待时间
			TimeoutSeconds:      5,
			PeriodSeconds:       5,  //执行间隔
		}

		// Liveness
		deployment.Spec.Template.Spec.Containers[0].LivenessProbe = &corev1.Probe{
			Handler:             corev1.Handler{
				HTTPGet: &corev1.HTTPGetAction{
					Path:        data.HealthPath,
					Port:        intstr.IntOrString{
						Type : 0,
						IntVal: data.ContainerPort,
					},
				},
			},
			InitialDelaySeconds: 15,
			TimeoutSeconds:      5,
			PeriodSeconds:       5,
		}
	}

	// 定义容器的limit与request资源
	deployment.Spec.Template.Spec.Containers[0].Resources.Limits =
		map[corev1.ResourceName]resource.Quantity{
			corev1.ResourceCPU: resource.MustParse(data.Cpu),
			corev1.ResourceMemory: resource.MustParse(data.Memory),
		}
	deployment.Spec.Template.Spec.Containers[0].Resources.Requests =
		map[corev1.ResourceName]resource.Quantity{
			corev1.ResourceCPU: resource.MustParse(data.Cpu),
			corev1.ResourceMemory: resource.MustParse(data.Memory),
		}

	//调用sdk创建deployment
	_, err = K8s.ClientSet.AppsV1().Deployments(data.Namespace).Create(context.TODO(), deployment, metav1.CreateOptions{})
	if err != nil {
		logger.Error(errors.New("创建Deployment失败, " + err.Error()))
		return errors.New("创建Deployment失败, " + err.Error())
	}
	return nil
}

//更新deployment
func(d *deployment) UpdateDeployment(namespace, content string) (err error) {
	var deploy = &appsv1.Deployment{}

	err = json.Unmarshal([]byte(content), deploy)
	if err != nil {
		logger.Error(errors.New("反序列化失败, " + err.Error()))
		return errors.New("反序列化失败, " + err.Error())
	}

	_, err = K8s.ClientSet.AppsV1().Deployments(namespace).Update(context.TODO(), deploy, metav1.UpdateOptions{})
	if err != nil {
		logger.Error(errors.New("更新Deployment失败, " + err.Error()))
		return errors.New("更新Deployment失败, " + err.Error())
	}
	return nil
}

// 删除Deployment
func (d *deployment) DeleteDeployment(deploymentName, namespace string) ( err error)  {
	err = K8s.ClientSet.AppsV1().Deployments(namespace).Delete(context.TODO(), deploymentName, metav1.DeleteOptions{})
	if  err!= nil{
		logger.Error("删除deployment失败," + err.Error())
		return  errors.New("删除deployment失败," + err.Error())
	}
	return nil
}

// 重启deployment
func (d *deployment) RestartDeployment(deploymentName , namespace string) (err error)  {

	//此功能等同于一下kubectl命令
	//kubectl deployment ${service} -p \
	//'{"spec":{"template":{"spec":{"containers":[{"name":"'"${service}"'","env": [{"name":"RESTART_","value":"'$(date +%s)'"}]}]}}}}'

	// 使用patchData Map组装数据
	patchData := map[string]interface{}{
		"spec": map[string]interface{}{
			"template": map[string]interface{}{
				"containers": []map[string]interface{}{
					{
						"name": deploymentName,
						"env": []map[string]string{
							{
								"name": "RESTART_",
								"value": strconv.FormatInt(time.Now().Unix(), 10),
							},
						},
					},
				},
			},
		},
	}

	// 序列化为字节，因为patch只接受字节类型的参数
	patchByte, err := json.Marshal(patchData)
	if  err!= nil{
		logger.Error("json序列化失败," + err.Error())
		return  errors.New("json序列化失败," + err.Error())
	}

	// 调用patch方法更新deployment
	_, err = K8s.ClientSet.AppsV1().Deployments(namespace).Patch(context.TODO(), deploymentName,
		"application/strategic-merge-patch+json", patchByte, metav1.PatchOptions{})
	if  err!= nil{
		logger.Error("重启Deployment失败," + err.Error())
		return  errors.New("重启Deployment失败," + err.Error())
	}

	return nil
}

//获取每个namespace的deployment数量
func(d *deployment) GetDeployNumPerNp() (deploysNps []*DeploysNp, err error) {
	namespaceList, err := K8s.ClientSet.CoreV1().Namespaces().List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		return nil, err
	}
	for _, namespace := range namespaceList.Items {
		deploymentList, err := K8s.ClientSet.AppsV1().Deployments(namespace.Name).List(context.TODO(), metav1.ListOptions{})
		if err != nil {
			return nil, err
		}

		deploysNp := &DeploysNp{
			Namespace: namespace.Name,
			DeploymentNum:    len(deploymentList.Items),
		}

		deploysNps = append(deploysNps, deploysNp)
	}
	return deploysNps, nil
}
