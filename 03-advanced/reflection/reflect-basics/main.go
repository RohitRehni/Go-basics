// Topic: reflect-basics | Concept: reflection | Stage: 03-advanced
// Run:   go run ./03-advanced/reflection/reflect-basics
//
// Explore:
// * reflect.TypeOf and reflect.ValueOf inspect values at runtime
// * Kind (struct, slice, int...) vs Type (main.User)
// * Walking struct fields and reading struct tags (how encoding/json works)
// * Setting a value needs a pointer: reflect.ValueOf(&x).Elem()
// * Calling a method by name
// * Reflection is slow and loses type safety: prefer generics or interfaces
package main

import (
	"fmt"
	"reflect"
)

type User struct {
	Name  string `json:"name" validate:"required"`
	Email string `json:"email" validate:"email"`
	Age   int    `json:"age"`
	notes string
}

func (u User) Hello(greeting string) string {
	return greeting + ", " + u.Name
}

func main() {
	fmt.Println("Topic: reflection")

	for _, v := range []any{42, "go", 3.5, []int{1}, map[string]int{}, User{}, &User{}} {
		t := reflect.TypeOf(v)
		fmt.Printf("  %-22v type=%-16v kind=%v\n", v, t, t.Kind())
	}

	u := User{Name: "Asha", Email: "asha@example.com", Age: 30, notes: "private"}
	t := reflect.TypeOf(u)
	v := reflect.ValueOf(u)
	fmt.Println("struct fields of", t.Name()+":")
	for i := range t.NumField() {
		f := t.Field(i)
		value := "(unexported)"
		if f.IsExported() {
			value = fmt.Sprint(v.Field(i).Interface())
		}
		fmt.Printf("  %-6s %-6s json=%-6q validate=%-10q value=%s\n",
			f.Name, f.Type, f.Tag.Get("json"), f.Tag.Get("validate"), value)
	}

	// Setting through a pointer
	ptr := reflect.ValueOf(&u).Elem()
	ptr.FieldByName("Age").SetInt(31)
	fmt.Println("age after reflect SetInt:", u.Age)
	fmt.Println("can set via plain value?", reflect.ValueOf(u).Field(2).CanSet())

	// Calling a method by name
	m := reflect.ValueOf(u).MethodByName("Hello")
	out := m.Call([]reflect.Value{reflect.ValueOf("Namaste")})
	fmt.Println("MethodByName(\"Hello\"):", out[0].String())
}
