package main

import (
	"context"
	"fmt"
	corev1 "k8s.io/api/core/v1"
	v12 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

func main() {
	// 本地加载 kubeconfig 配置文件
	config, err := clientcmd.BuildConfigFromFlags("", clientcmd.RecommendedHomeFile)
	if err != nil {
		panic(err)
	}
	// 参考： /api/v1/namespaces/{namespace}/pods
	config.APIPath = "/api"
	// pod 的 group 是空字符串
	config.GroupVersion = &corev1.SchemeGroupVersion
	// 指定序列化工具
	config.NegotiatedSerializer = scheme.Codecs

	// 根据配置信息得到 restClient 实例
	restClient, err := rest.RESTClientFor(config)
	if err != nil {
		panic(err)
	}
	// 保存 pod 结果的数据结构实例
	pods := &corev1.PodList{}

	err = restClient.
		Get().
		// 指定 namespace，参考： /api/v1/namespaces/{namespace}/pods
		Namespace("kube-system").
		// 查找多个 pod，参考： /api/v1/namespaces/{namespace}/pods
		Resource("Pods").
		// 指定大小限制和序列化工具
		VersionedParams(&v12.ListOptions{Limit: 100}, scheme.ParameterCodec).
		// 请求
		Do(context.TODO()).
		// 结果存入 result
		Into(pods)
	if err != nil {
		panic(err)
	}
	fmt.Printf("namespce\t status\t\t name\n")
	for _, pod := range pods.Items {
		fmt.Printf("%v\t %v\t %v\n", pod.Namespace, pod.Status.Phase, pod.Name)
	}
}
