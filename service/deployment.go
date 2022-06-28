package service

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"github.com/wonderivan/logger"
	"io"
	"k8s-platform/config"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/json"
)

var Pod pod

type pod struct {}

// 定义列表的返回内容;Items是pod元素列表，Total是元素数量
type PodsResp struct {
	Total int `json:"total"`
	Items []corev1.Pod `json:"items"`
}


// 获取pod列表，支持过滤、排序、分页
func (p *pod) GetPods(filterName, namespace string, limit , page int) (podsResp *PodsResp, err error)  {

	//context.TODO()用于声明一个空的context上下文，用于List方法内设置这个请求的超时（源码），这里 的常用用法
	//metav1.ListOptions{}用于过滤List数据，如使用label，field等
	//kubectl get services --all-namespaces --field-seletor metadata.namespace != default
	podList, err := K8s.ClientSet.CoreV1().Pods(namespace).List(context.TODO(), metav1.ListOptions{})
	if err != nil{
		// 打印日志给自己看，排错使用
		logger.Error("获取pod列表失败，" + err.Error())
		// 返回给上一层，最终返回给前端，前端打印出来这个error
		return nil, errors.New("获取Pod列表失败，"+ err.Error())
	}

	// 实例化dataSelector结构体，组装数据
	selectTableData :=&dataSelector{
		GenericDataList: p.toCells(podList.Items),
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
	pods := p.fromCells(data.GenericDataList)

	// 数据处理后的数据与原始数据的比较
	// 处理后的数据
	fmt.Println("处理后的数据")
	for _, pod := range  pods{
		fmt.Println(pod.Name, pod.CreationTimestamp.Time)
	}
	// 原始数据
	fmt.Println("原始数据")
	for _, pod := range  podList.Items{
		fmt.Println(pod.Name, pod.CreationTimestamp.Time)
	}

	return &PodsResp{
		Total: total,
		Items: pods,
	}, nil
	

}



// 获取pods列表，支持过滤、排序、分页
// 类型转换的方法 corev1.Pod --> DataCell, DataCell --> corev1.Pod
func (p *pod) toCells(pods []corev1.Pod)  []DataCell {
	cells := make([]DataCell, len(pods))
	for i := range pods{
		// todo // 这里怎么转换的？？
		cells[i] = podCell(pods[i])
	}
	return cells
}

func (p *pod) fromCells(cells []DataCell) []corev1.Pod  {
	pods := make([]corev1.Pod, len(cells))
	for i := range cells{
		// cells[i].(podCell) 是将DataCell类型转换成podCell
		pods[i] = corev1.Pod(cells[i].(podCell))
	}
	return pods

}


// 获取pod详情
func (p *pod) GetPodDetail(podName, namespace string) (pod *corev1.Pod, err error)  {
	pod, err = K8s.ClientSet.CoreV1().Pods(namespace).Get(context.TODO(), podName, metav1.GetOptions{})
	if  err!= nil{
		logger.Error("获取pod详情失败," + err.Error())
		return nil, errors.New("获取pod详情失败," + err.Error())
	}
	return pod, nil
}


// 删除pod
func (p *pod) DeletePod(podName, namespace string) ( err error)  {
	err = K8s.ClientSet.CoreV1().Pods(namespace).Delete(context.TODO(), podName, metav1.DeleteOptions{})
	if  err!= nil{
		logger.Error("删除pod失败," + err.Error())
		return  errors.New("获取pod详情失败," + err.Error())
	}
	return nil
}

// 更新pod
func (p *pod) UpdatePod(podName, namespace, content string) (res *corev1.Pod, err error)  {

	var pod = &corev1.Pod{}
	err = json.Unmarshal([]byte(content), pod)
	if err != nil{
		logger.Error("反序列化失败，"+err.Error())
		return nil, errors.New("反序列化失败，"+err.Error())
	}

	newPod, err := K8s.ClientSet.CoreV1().Pods(namespace).Update(context.TODO(), pod, metav1.UpdateOptions{})
	if  err!= nil{
		logger.Error("更新pod失败," + err.Error())
		return nil,errors.New("更新pod失败," + err.Error())
	}

	return newPod, nil
}

// 获取pod的容器名列表
func (p *pod) GetPodContainer(podName, namespace string) (containers []string, err error)  {

	pod, err := p.GetPodDetail(podName, namespace)
	if err!= nil{
		logger.Error("获取pod详情," + err.Error())
		return containers, errors.New("获取pod详情失败," + err.Error())
	}

	for _, item := range pod.Spec.Containers {
		containers = append(containers, item.Name)
	}

	return containers, nil
}

// 获取pod内容器日志
func (p *pod) GetPodLog(containerName, podName, namespace string) (Log string , err error)  {
	// 设置日志的配置，容器名，获取的内容的配置
	lineLimit := int64(config.PodLogTailLine)
	option := &corev1.PodLogOptions{
		Container:                    containerName,
		TailLines:                    &lineLimit,
	}
	// 获取一个request实例
	req := K8s.ClientSet.CoreV1().Pods(namespace).GetLogs(podName, option)
	// 发起stream连接，得到Response.body
	podLogs, err := req.Stream(context.TODO())
	if err!= nil{
		logger.Error(errors.New("获取PodLog失败, " + err.Error()))
		return "", errors.New("获取PodLog失败, " + err.Error())
	}
	defer podLogs.Close()
	// 将response body写入到缓冲区，目的是为了转换成string类型
	buf := new(bytes.Buffer)
	_,err = io.Copy(buf,podLogs)
	if err != nil {
		logger.Error(errors.New("复制PodLog失败, " + err.Error()))
		return "", errors.New("复制PodLog失败, " + err.Error())
	}

	return buf.String(), nil
}

// 获取每个namespace的pod数量
type PodsNp struct {
	Namespace string
	PodNum int
}

func(p *pod) GetPodNumPerNp() (podsNps []*PodsNp, err error) {
	//获取namespace列表
	namespaceList, err := K8s.ClientSet.CoreV1().Namespaces().List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		return nil, err
	}

	for _, namespace := range namespaceList.Items {
		//获取pod列表
		podList, err := K8s.ClientSet.CoreV1().Pods(namespace.Name).List(context.TODO(), metav1.ListOptions{})
		if err != nil {
			return nil, err
		}
		//组装数据
		podsNp := &PodsNp{
			Namespace: namespace.Name,
			PodNum: len(podList.Items),
			}
		//添加到podsNps数组中
		podsNps = append(podsNps, podsNp) }
	return podsNps, nil
}