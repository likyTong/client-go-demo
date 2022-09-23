package main

import (
	"fmt"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/discovery"
	"k8s.io/client-go/tools/clientcmd"
)

func main() {
	config, err := clientcmd.BuildConfigFromFlags("", clientcmd.RecommendedHomeFile)
	if err != nil {
		panic(err)
	}

	// 新建 discoveryClient 实例
	discoveryClient, err := discovery.NewDiscoveryClientForConfig(config)
	if err != nil {
		panic(err.Error())
	}

	ApiGroup, APIResourceListSlice, err := discoveryClient.ServerGroupsAndResources()
	if err != nil {
		panic(err.Error())
	}

	// 看 Group 信息
	fmt.Printf("APIGroup: \n\n  %v \n\n\n\n", ApiGroup)

	for _, singleAPIResourceList := range APIResourceListSlice {
		groupVersionStr := singleAPIResourceList.GroupVersion
		gv, err := schema.ParseGroupVersion(groupVersionStr)
		if err != nil {
			panic(err.Error())
		}

		fmt.Println("******************************")
		fmt.Printf("GV string [%v]\nGV struct [%#v]\nresources:\n\n", groupVersionStr, gv)

		for _, singleAPIResource := range singleAPIResourceList.APIResources {
			fmt.Printf("%v\n", singleAPIResource.Name)
		}
	}

}
