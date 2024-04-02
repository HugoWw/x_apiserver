package main

import (
	"flag"
	"fmt"
	"github.com/HugoWw/x_apiserver/pkg/apiserver/resources"
	_ "github.com/HugoWw/x_apiserver/pkg/resource"
	restfulspec "github.com/emicklei/go-restful-openapi/v2"
	"github.com/go-openapi/spec"
	"os"
	"path"
	"reflect"
	"regexp"
	"strings"
)

func main() {

	var docDir string
	flag.StringVar(&docDir, "doc-dir", "./", "swagger document outout dir")
	flag.Parse()

	b, err := buildSwag()
	if err != nil {
		fmt.Printf("[error] build swagger %v error: %v\n", docDir, err)
		os.Exit(1)
	}

	if err := genericSwagDocFile(b, docDir); err != nil {
		fmt.Printf("[error] generate swagger json file to %v error: %v\n", docDir, err)
		os.Exit(1)
	}

	fmt.Printf("[info] success to generate swagger json file to %v\n", docDir)
}

func buildSwag() ([]byte, error) {
	config := restfulspec.Config{
		Host:                          "127.0.0.1",
		APIPath:                       "/apidoc.json",
		WebServices:                   resources.Default.RegisteredWebServices(),
		DisableCORS:                   false,
		PostBuildSwaggerObjectHandler: enrichSwaggerObject,
		ModelTypeNameHandler:          enrichModelTypeName,
	}

	swg := restfulspec.BuildSwagger(config)
	b, err := swg.MarshalJSON()
	if err != nil {
		return nil, err
	}

	return b, nil
}

func enrichModelTypeName(st reflect.Type) (string, bool) {
	var output string
	key := st.String() //v1.APIResponse[github.com.HugoWw/x_apiserver/pkg/resource/v1.AuthData]

	pattern := `\[(.*?)\]`
	re := regexp.MustCompile(pattern)
	match := re.FindStringSubmatch(key)

	if len(match) > 0 {
		// get submatch
		// Example:
		// match[1] = github.com.HugoWw/x_apiserver/pkg/resource/v1.AuthData
		content := match[1]
		end := strings.LastIndex(content, "/")
		replaceContent := content[end+1:]

		output = re.ReplaceAllString(key, "["+replaceContent+"]")
	} else {
		return "", false
	}

	return output, true
}

func genericSwagDocFile(b []byte, filePath string) error {
	var filename = "api-swagger.json"
	jsonFile := path.Join(filePath, filename)

	f, err := os.Create(jsonFile)
	if err != nil {
		return err
	}

	defer f.Close()

	if _, err := f.Write(b); err != nil {
		return err
	}

	return nil
}

func enrichSwaggerObject(swo *spec.Swagger) {
	swo.Info = &spec.Info{
		InfoProps: spec.InfoProps{
			Title:       "X-APIServer",
			Description: "X-APIServer Swagger Doc",
			Contact: &spec.ContactInfo{
				ContactInfoProps: spec.ContactInfoProps{
					Name:  "gouhuan",
					Email: "670195398@qq.com",
					URL:   "",
				},
			},
			License: &spec.License{
				LicenseProps: spec.LicenseProps{
					Name: "MIT",
					URL:  "http://mit.org",
				},
			},
			Version: "1.0.0",
		},
	}
	//swo.Tags = []spec.Tag{spec.Tag{TagProps: spec.TagProps{
	//	Name:        "users",
	//	Description: "Managing users"}}}
}
