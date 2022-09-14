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
	config, err := clientcmd.BuildConfigFromFlags("", clientcmd.RecommendedHomeFile)
	if err != nil {
		panic(err)
	}

	config.APIPath = "/api"
	config.GroupVersion = &corev1.SchemeGroupVersion
	config.NegotiatedSerializer = scheme.Codecs

	restClient, err := rest.RESTClientFor(config)
	if err != nil {
		panic(err)
	}
	pods := &corev1.PodList{}

	err = restClient.
		Get().
		Namespace("kube-system").
		Resource("Pods").
		VersionedParams(&v12.ListOptions{Limit: 100}, scheme.ParameterCodec).
		Do(context.TODO()).
		Into(pods)
	if err != nil {
		panic(err)
	}
	fmt.Printf("namespce\t status\t name\n")
	for _, pod := range pods.Items {
		fmt.Printf("%v\t %v\t %v\n", pod.Namespace, pod.Status.Phase, pod.Name)
	}
}
