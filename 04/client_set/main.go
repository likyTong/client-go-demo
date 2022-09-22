package main

import (
	"context"
	"flag"
	"fmt"
	appsv1 "k8s.io/api/apps/v1"
	v1 "k8s.io/api/core/v1"
	v12 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/utils/pointer"
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
		createDeployment(clientset)
	}
}

func clean(clientset *kubernetes.Clientset) {
	//clientset.AppsV1().Deployments(NAMESPAECE).Delete(context.TODO(), DEPLOYMENT_NAME)
}

// 新建 namespace
func createNamespace(clientSet *kubernetes.Clientset) {
	namespaceClient := clientSet.CoreV1().Namespaces()

	namespace := &v1.Namespace{
		ObjectMeta: v12.ObjectMeta{
			Name: NAMESPAECE,
		},
	}

	result, err := namespaceClient.Create(context.TODO(), namespace, v12.CreateOptions{})
	if err != nil {
		panic(err.Error())
	}
	fmt.Printf("create namespace %s \n", result.GetName())
}

// 新建 deployment
func createDeployment(clientSet *kubernetes.Clientset) {
	deploymentClient := clientSet.AppsV1().Deployments(NAMESPAECE)
	deployment := &appsv1.Deployment{
		ObjectMeta: v12.ObjectMeta{
			Name: DEPLOYMENT_NAME,
		},
		Spec: appsv1.DeploymentSpec{
			Replicas: pointer.Int32Ptr(2),
			Selector: &v12.LabelSelector{MatchLabels: map[string]string{"app": "tomcat"}},
			Template: v1.PodTemplateSpec{
				ObjectMeta: v12.ObjectMeta{Labels: map[string]string{"app": "tomcat"}}, Spec: v1.PodSpec{Containers: []v1.Container{{Name: "http", Image: "tocat:8.0.18-jre8", ImagePullPolicy: "IfNotPresent", Ports: []v1.ContainerPort{{Name: "http", Protocol: v1.ProtocolSCTP, ContainerPort: 8080}}}}}},
		},
	}
	result, err := deploymentClient.Create(context.TODO(), deployment, v12.CreateOptions{})
	if err != nil {
		panic(err)
	}
	fmt.Printf("create deployment %s\n", result.GetName())
}
