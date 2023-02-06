/*
   Copyright The containerd Authors.

   Licensed under the Apache License, Version 2.0 (the "License");
   you may not use this file except in compliance with the License.
   You may obtain a copy of the License at

       http://www.apache.org/licenses/LICENSE-2.0

   Unless required by applicable law or agreed to in writing, software
   distributed under the License is distributed on an "AS IS" BASIS,
   WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
   See the License for the specific language governing permissions and
   limitations under the License.
*/

package typeurl

import (
	"bytes"
	"reflect"
	"testing"
	"time"

	"github.com/gogo/protobuf/proto"
	gogotypes "github.com/gogo/protobuf/types"
	"google.golang.org/protobuf/types/known/anypb"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type test struct {
	Name string
	Age  int
}

func clear() {
	registry = make(map[reflect.Type]string)
}

var _ Any = &gogotypes.Any{}
var _ Any = &anypb.Any{}

func TestRegisterPointerGetPointer(t *testing.T) {
	clear()
	expected := "test"
	Register(&test{}, "test")

	url, err := TypeURL(&test{})
	if err != nil {
		t.Fatal(err)
	}
	if url != expected {
		t.Fatalf("expected %q but received %q", expected, url)
	}
}

func TestMarshal(t *testing.T) {
	clear()
	expected := "test"
	Register(&test{}, "test")

	v := &test{
		Name: "koye",
		Age:  6,
	}
	any, err := MarshalAny(v)
	if err != nil {
		t.Fatal(err)
	}
	if any.GetTypeUrl() != expected {
		t.Fatalf("expected %q but received %q", expected, any.GetTypeUrl())
	}

	// marshal it again and make sure we get the same thing back.
	newany, err := MarshalAny(any)
	if err != nil {
		t.Fatal(err)
	}

	val := any.GetValue()
	newval := newany.GetValue()

	// Ensure pointer to same exact slice
	newval[0] = val[0] ^ 0xff

	if !bytes.Equal(newval, val) {
		t.Fatalf("expected to get back same object: %v != %v", newany, any)
	}

}

func TestMarshalUnmarshal(t *testing.T) {
	clear()
	Register(&test{}, "test")

	v := &test{
		Name: "koye",
		Age:  6,
	}
	any, err := MarshalAny(v)
	if err != nil {
		t.Fatal(err)
	}
	nv, err := UnmarshalAny(any)
	if err != nil {
		t.Fatal(err)
	}
	td, ok := nv.(*test)
	if !ok {
		t.Fatal("expected value to cast to *test")
	}
	if td.Name != "koye" {
		t.Fatal("invalid name")
	}
	if td.Age != 6 {
		t.Fatal("invalid age")
	}
}

func TestMarshalUnmarshalTo(t *testing.T) {
	clear()
	Register(&test{}, "test")

	in := &test{
		Name: "koye",
		Age:  6,
	}
	any, err := MarshalAny(in)
	if err != nil {
		t.Fatal(err)
	}
	out := &test{}
	err = UnmarshalTo(any, out)
	if err != nil {
		t.Fatal(err)
	}
	if out.Name != "koye" {
		t.Fatal("invalid name")
	}
	if out.Age != 6 {
		t.Fatal("invalid age")
	}
}

type test2 struct {
	Name string
}

func TestUnmarshalToInvalid(t *testing.T) {
	clear()
	Register(&test{}, "test1")
	Register(&test2{}, "test2")

	in := &test{
		Name: "koye",
		Age:  6,
	}
	any, err := MarshalAny(in)
	if err != nil {
		t.Fatal(err)
	}

	out := &test2{}
	err = UnmarshalTo(any, out)
	if err == nil || err.Error() != `can't unmarshal type "test1" to output "test2"` {
		t.Fatalf("unexpected result: %+v", err)
	}
}

func TestIs(t *testing.T) {
	clear()
	Register(&test{}, "test")

	v := &test{
		Name: "koye",
		Age:  6,
	}
	any, err := MarshalAny(v)
	if err != nil {
		t.Fatal(err)
	}
	if !Is(any, &test{}) {
		t.Fatal("Is(any, test{}) should be true")
	}
}

func TestRegisterDiffUrls(t *testing.T) {
	clear()
	defer func() {
		if err := recover(); err == nil {
			t.Error("registering the same type with different urls should panic")
		}
	}()
	Register(&test{}, "test")
	Register(&test{}, "test", "two")
}

func TestUnmarshalNil(t *testing.T) {
	var pba *anypb.Any // This is nil.
	var a Any = pba    // This is typed nil.

	if pba != nil {
		t.Fatal("pbany must be nil")
	}
	if a == nil {
		t.Fatal("nilany must not be nil")
	}

	actual, err := UnmarshalAny(a)
	if err != nil {
		t.Fatal(err)
	}

	if actual != nil {
		t.Fatalf("expected nil, got %v", actual)
	}
}

func TestCheckNil(t *testing.T) {
	var a *anyType

	actual := a.GetValue()
	if actual != nil {
		t.Fatalf("expected nil, got %v", actual)
	}
}

func TestProtoFallback(t *testing.T) {
	expected := time.Now()
	b, err := proto.Marshal(timestamppb.New(expected))
	if err != nil {
		t.Fatal(err)
	}
	x, err := UnmarshalByTypeURL("type.googleapis.com/google.protobuf.Timestamp", b)
	if err != nil {
		t.Fatal(err)
	}
	ts, ok := x.(*timestamppb.Timestamp)
	if !ok {
		t.Fatalf("failed to convert %+v to Timestamp", x)
	}
	if expected.Sub(ts.AsTime()) != 0 {
		t.Fatalf("expected %+v but got %+v", expected, ts.AsTime())
	}
}
