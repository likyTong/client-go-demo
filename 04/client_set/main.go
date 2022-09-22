package main

import (
	"context"
	"flag"
	"fmt"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
)

const (
	NAMESPAECE      = "test-clientset"
	DEPLOYMENT_NAME = "client-test-deployment"
	SERVICE_NAME    = "client-test-service"
)

func main() {
	operate := flag.String("operate", "create", "operate type: create or clean")
	flag.Parse()

	config, err := clientcmd.BuildConfigFromFlags("", clientcmd.RecommendedHomeFile)
	if err != nil {
		panic(err.Error())
	}

	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		panic(err.Error())
	}
	fmt.Printf("operation is %v\n", operate)

	if "clean" == *operate {
		clean(clientset)
	} else {
		createNamespace(clientset)
	}
}

func clean(clientset *kubernetes.Clientset) {
	//clientset.AppsV1().Deployments(NAMESPAECE).Delete(context.TODO(), DEPLOYMENT_NAME)
}

// 新建 namespace
func createNamespace(clientSet *kubernetes.Clientset) {
	namespaceClient := clientSet.CoreV1().Namespaces()

	namespace := &corev1.Namespace{
		ObjectMeta: metav1.ObjectMeta{
			Name: NAMESPAECE,
		},
	}

	result, err := namespaceClient.Create(context.TODO(), namespace, metav1.CreateOptions{})
	if err != nil {
		panic(err.Error())
	}
	fmt.Printf("create namespace %s \n", result.GetName())
}

// 创建 service
func createService(clientSet *kubernetes.Clientset) {
	serviceClient := clientSet.CoreV1().Services(NAMESPAECE)

	service := &corev1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name: SERVICE_NAME,
		},
		Spec: corev1.ServiceSpec{
			Ports: []corev1.ServicePort{
				{
					Name:     "http",
					Port:     8080,
					NodePort: 30080,
				},
			},
			Selector: map[string]string{
				"app": "tomcat",
			},
			Type: corev1.ServiceTypeNodePort,
		},
	}
	result, err := serviceClient.Create(context.TODO(), service, metav1.CreateOptions{})
	if err != nil {
		panic(err.Error())
	}
	fmt.Printf("Create service %s\n", result.GetName())
}
