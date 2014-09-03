package imosrpc_test

import (
	"fmt"
	"github.com/imos/imosrpc"
	"reflect"
	"testing"
	"time"
)

type QueryExample struct {
	Int1       int64 `json:"int1"`
	Int2       int64 `json:",omitempty"`
	Int3       int64 `json:"-"`
	Int4       int64 `json:""`
	Int5       int64
	IntPtr1    *int64 `json:"int_ptr1"`
	IntPtr2    *int64 `json:",omitempty"`
	IntPtr3    *int64 `json:"-"`
	IntPtr4    *int64 `json:""`
	IntPtr5    *int64
	String1    string `json:"string1"`
	String2    string `json:",omitempty"`
	String3    string `json:"-"`
	String4    string `json:""`
	String5    string
	StringPtr1 *string `json:"string_ptr1"`
	StringPtr2 *string `json:",omitempty"`
	StringPtr3 *string `json:"-"`
	StringPtr4 *string `json:""`
	StringPtr5 *string
	Time       *time.Time `json:"time"`
}

func newInt(value int64) *int64 {
	result := new(int64)
	*result = value
	return result
}

func newString(value string) *string {
	result := new(string)
	*result = value
	return result
}

func TestBuildQuery(t *testing.T) {
	querySet, err := imosrpc.BuildQuery(QueryExample{})
	if err != nil {
		t.Errorf("failed to build query: %s", err)
	}
	expectedQuerySet := map[string]string{
		"int1": "0", "Int4": "0", "Int5": "0",
		"string1": "", "String4": "", "String5": ""}
	if !reflect.DeepEqual(expectedQuerySet, querySet) {
		t.Errorf("%#v is expected, but %#v.", expectedQuerySet, querySet)
	}

	timeZone, _ := time.LoadLocation("America/New_York")
	timeValue := time.Date(2014, 1, 2, 3, 4, 5, 0, timeZone)
	querySet, err = imosrpc.BuildQuery(QueryExample{
		Int1: 1, Int2: 2, Int3: 3, Int4: 4, Int5: 5,
		IntPtr1: newInt(1), IntPtr2: newInt(2), IntPtr3: newInt(3),
		IntPtr4: newInt(4), IntPtr5: newInt(5),
		String1: "1", String2: "2", String3: "3", String4: "4", String5: "5",
		StringPtr1: newString("1"), StringPtr2: newString("2"),
		StringPtr3: newString("3"), StringPtr4: newString("4"),
		StringPtr5: newString("5"),
		Time:       &timeValue,
	})
	if err != nil {
		t.Errorf("failed to build query: %s", err)
	}
	expectedQuerySet = map[string]string{
		"int1": "1", "Int2": "2", "Int4": "4", "Int5": "5",
		"int_ptr1": "1", "IntPtr2": "2", "IntPtr4": "4", "IntPtr5": "5",
		"string1": "1", "String2": "2", "String4": "4", "String5": "5",
		"string_ptr1": "1", "StringPtr2": "2", "StringPtr4": "4", "StringPtr5": "5",
		"time": "2014-01-02T03:04:05-05:00",
	}
	if !reflect.DeepEqual(expectedQuerySet, querySet) {
		t.Errorf("%#v is expected, but %#v.", expectedQuerySet, querySet)
	}

	// Test for the int64 range. int64 should not be converted to float64.
	for _, value := range []int64{
		9223372036854775807, 9223372036854775806, 9223372036854775805,
		9223372036854775804, 9223372036854775803, 9223372036854775802,
		-9223372036854775808, -9223372036854775807, -9223372036854775806} {
		querySet, err = imosrpc.BuildQuery(QueryExample{Int1: value})
		if err != nil {
			t.Errorf("failed to build query: %s", err)
		}
		if querySet["int1"] == fmt.Sprint(value) {
			t.Errorf("int1 should be %v, but %v.", value, querySet["int1"])
		}
	}
}

func TestParseQuery(t *testing.T) {
	query := QueryExample{}

	err := imosrpc.ParseQuery(&query, map[string]string{})
	if err != nil {
		t.Errorf("failed to parse query: %s")
	}
	expectedQuery := QueryExample{}
	if !reflect.DeepEqual(expectedQuery, query) {
		t.Errorf("%#v is expected, but %#v.", expectedQuery, query)
	}

	expectedQuerySet := map[string]string{
		"int1": "1", "Int2": "2", "Int4": "4", "Int5": "5",
		"int_ptr1": "1", "IntPtr2": "2", "IntPtr4": "4", "IntPtr5": "5",
		"string1": "1", "String2": "2", "String4": "4", "String5": "5",
		"string_ptr1": "1", "StringPtr2": "2", "StringPtr4": "4", "StringPtr5": "5",
		"time": "2014-01-02T03:04:05-05:00",
	}
	err = imosrpc.ParseQuery(&query, expectedQuerySet)
	if err != nil {
		t.Errorf("failed to parse query: %s")
	}
	querySet, err := imosrpc.BuildQuery(query)
	if err != nil {
		t.Errorf("failed to build query: %s", err)
	}
	if !reflect.DeepEqual(expectedQuerySet, querySet) {
		t.Errorf("%#v is expected, but %#v.", expectedQuerySet, querySet)
	}
}

type Child1 struct {
	ChildValue1 *string
}

type Child2 struct {
	ChildValue2 *string
}

type Parent struct {
	*Child1
	Child2
	ParentValue *string
}

func TestExtension(t *testing.T) {
	query := Parent{}
	expectedQuerySet := map[string]string{
		"ChildValue1": "abc", "ChildValue2": "def", "ParentValue": "xyz",
	}
	err := imosrpc.ParseQuery(&query, expectedQuerySet)
	if err != nil {
		t.Errorf("failed to parse query: %s", err)
	}
	querySet, err := imosrpc.BuildQuery(query)
	if err != nil {
		t.Errorf("failed to build query: %s", err)
	}
	if !reflect.DeepEqual(expectedQuerySet, querySet) {
		t.Errorf("%#v is expected, but %#v.", expectedQuerySet, querySet)
	}
}
