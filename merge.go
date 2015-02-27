package main

import (
	"fmt"
	"reflect"
)

func Merge(op_a, op_b interface{}) interface{} {
	av := reflect.ValueOf(op_a)
	bv := reflect.ValueOf(op_b)
	ans := reflect.New(av.Type()).Elem()
	merge(ans, av, bv)
	return ans.Interface()
}

func merge(res, op_av, op_bv reflect.Value) {
	if op_av.Kind() != op_bv.Kind() {
		// panic
	}
	switch op_av.Kind() {
	case reflect.Int, reflect.Int64:
		fmt.Println("can", res.CanSet())
		res.SetInt(op_av.Int() + op_bv.Int())
	case reflect.Map:
		fmt.Println("can", res.CanSet())
		res.Set(reflect.MakeMap(op_av.Type()))
		keys := op_av.MapKeys()
		for _, k := range op_bv.MapKeys() {
			check := 0
			for _, v := range keys {
				if k == v {
					check = 1
				}
			}
			if check == 0 {
				keys = append(keys, k)
			}
		}
		for _, k := range keys {
			av := op_av.MapIndex(k)
			bv := op_bv.MapIndex(k)
			r := reflect.New(av.Type()).Elem()
			merge(r, av, bv)
			res.SetMapIndex(k, r)
		}
	}
}

func main() {
	x := make(map[int]map[int]int)
	y := make(map[int]int)
	y[1] = 1
	x[1] = y
	fmt.Println("x", x)
	fmt.Println("y", y)
	fmt.Println("2x", Merge(x, x))
}
