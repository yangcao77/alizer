package main

import (
	"fmt"
	"reflect"

	"github.com/devfile/registry-support/index/generator/schema"
	registryLibrary "github.com/devfile/registry-support/registry-library/library"
	"github.com/redhat-developer/alizer/go/pkg/apis/model"
	"github.com/redhat-developer/alizer/go/pkg/apis/recognizer"
)

func main() {
	devfileRegistryURL := "https://registry.devfile.io"
	//path := "devfile-sample-java-springboot-basic-main"
	//path := "httpd-shell-master"
	path := "dotnet-hello-main"
	alizerDevfileTypes, err := getAlizerDevfileTypes(devfileRegistryURL)
	if err != nil {
		fmt.Printf("Error: %v", err)
		return
	} else {
		fmt.Printf("alizerDevfileTypes is %v \n", alizerDevfileTypes)
	}
	results, err := recognizer.DetectComponents("/Users/stephanie/Downloads/" + path)
	if err != nil {
		fmt.Printf("Error: %v", err)
		return
	}
	fmt.Printf("lenth of results %v \n", len(results))
	for _, element := range results {
		fmt.Println("**************")
		fmt.Printf("component %v,path: %v \n languages: %v \n ports: %v \n", element.Name, element.Path, element.Languages, element.Ports)
		fmt.Println("**************")
		for _, language := range element.Languages {
			if language.CanBeComponent {
				// if we get one language analysis that can be a component
				// we can then determine a devfile from the registry and return

				// The highest rank is the most suggested component. priorty: configuration file > high %

				index, err := recognizer.SelectDevFileFromTypes(element.Path, alizerDevfileTypes)

				detectedType := alizerDevfileTypes[index]
				fmt.Printf("detectedType: %v \n", detectedType)
				// fmt.Printf("index is: %v, detectedType is %v \n", index, detectedType)
				if err != nil && err.Error() != fmt.Sprintf("No valid devfile found for project in %s", element.Path) {
					// No need to check for err, if a path does not have a detected devfile, ignore err
					// if a dir can be a component but we get an unrelated err, err out
					fmt.Printf("Error: %v", err)
					return
				} else if !reflect.DeepEqual(detectedType, model.DevFileType{}) {
					detectedDevfileEndpoint := devfileRegistryURL + "/devfiles/" + detectedType.Name
					// devfileBytes, err = CurlEndpoint(detectedDevfileEndpoint)
					// if err != nil {
					// 	return nil, "", "", err
					// }

					// if len(devfileBytes) > 0 {
					// 	return devfileBytes, detectedDevfileEndpoint, detectedType.Name, nil
					// }
					fmt.Printf("detectedDevfileEndpoint: %v \n", detectedDevfileEndpoint)
				}
			}
		}
	}

}

// getAlizerDevfileTypes gets the Alizer devfile types for a specified registry
func getAlizerDevfileTypes(registryURL string) ([]model.DevFileType, error) {
	types := []model.DevFileType{}
	registryIndex, err := registryLibrary.GetRegistryIndex(registryURL, registryLibrary.RegistryOptions{
		Telemetry: registryLibrary.TelemetryData{},
	}, schema.SampleDevfileType)
	if err != nil {
		return nil, fmt.Errorf("cannot get registry index: %v", err)
	}

	for _, index := range registryIndex {
		types = append(types, model.DevFileType{
			Name:        index.Name,
			Language:    index.Language,
			ProjectType: index.ProjectType,
			Tags:        index.Tags,
		})
	}

	return types, nil
}
