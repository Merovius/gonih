package cmp_test

import (
	"fmt"

	"golang.org/x/exp/slices"
	"gonih.org/cmp"
	pb "gonih.org/cmp/internal/addresspb"
)

func Example_chain() {
	// Sort people first ascending by name, then descending by age.
	type Person struct {
		Name string
		Age  int
	}
	people := []Person{
		{"Paula", 42},
		{"Joziah", 23},
		{"Austin", 32},
		{"Caleb", 23},
		{"Paula", 37},
		{"Austin", 45},
	}

	// Note: This is golang.org/x/exp/slices. The slices package from Go 1.21
	// does not require the .Less.
	slices.SortFunc(people, cmp.Chain(
		cmp.By(func(p Person) string { return p.Name }),
		cmp.Reverse(cmp.By(func(p Person) int { return p.Age })),
	).Less)

	for _, p := range people {
		fmt.Println(p)
	}
	// Output:
	// {Austin 45}
	// {Austin 32}
	// {Caleb 23}
	// {Joziah 23}
	// {Paula 42}
	// {Paula 37}
}

func Example_proto() {
	// Protocol buffers generate getters for all fields. This plays well with
	// the By* helpers.
	// You can find the generated protobuf message used as an example in
	// gonih.org/cmp/internal/addresspb.

	// Sort numbers by type and then number. Use the string representation of
	// the type for sorting. nil messages should sort first.
	cmpNumber := cmp.PointerFunc(cmp.Chain(
		cmp.ByFunc((*pb.PhoneNumber).GetType, cmp.By(pb.PhoneType.String)),
		cmp.By((*pb.PhoneNumber).GetNumber),
	))
	// Sort people by name, email and then number. nil messages should sort first.
	cmpPerson := cmp.PointerFunc(cmp.Chain(
		cmp.By((*pb.Person).GetName),
		cmp.By((*pb.Person).GetEmail),
		cmp.ByFunc((*pb.Person).GetNumbers, cmp.SliceFunc[[]*pb.PhoneNumber](cmpNumber)),
	))
	var people []*pb.Person
	slices.SortFunc(people, cmpPerson.Less)
}
