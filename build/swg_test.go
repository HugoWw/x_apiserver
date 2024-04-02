package main

import (
	"fmt"
	v1 "github.com/HugoWw/x_apiserver/pkg/resource/v1"
	"reflect"
	"testing"
)

func TestMains(t *testing.T) {

	st := v1.APIResponse[v1.AuthData]{}
	st_type := reflect.TypeOf(st)
	fmt.Printf("st_type.Name:%v, st_type.String:%v, st_type.Kind:%v\n", st_type.Name(), st_type.String(), st_type.Kind())
}
