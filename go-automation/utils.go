package goautomation

import (
	"bytes"
	"io"
	"io/ioutil"

	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/util/yaml"
)

func ParseYaml(yamlPath string) ([]*unstructured.Unstructured, error) {
    yamlFile, err := ioutil.ReadFile(yamlPath)
    if err != nil {
        return nil, err
    }
    var objects []*unstructured.Unstructured
    decoder := yaml.NewYAMLOrJSONDecoder(bytes.NewReader(yamlFile), 4096)
    for {
        obj := &unstructured.Unstructured{}
        if err := decoder.Decode(obj); err != nil {
            if err == io.EOF {
                break
            }
            return nil, err
        }
        objects = append(objects, obj)
    }
    return objects, nil
}
