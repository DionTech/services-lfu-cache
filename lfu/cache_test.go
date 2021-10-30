package lfu

import (
	"testing"
)

func TestInit(t *testing.T) {
	myStore := Init(5)

	isCache := func(t interface{}) bool {
		switch t.(type) {
		case Cache:
			return true
		default:
			return false
		}
	}

	if !isCache(myStore) {
		t.Fatalf("wrong init type")
	}

	if myStore.size != 5 {
		t.Fatalf("wrong cache size %d", myStore.size)
	}
}

func TestPutITem(t *testing.T) {
	put, err := LFU.Put("foo", "bar")

	if err != nil {
		t.Fatalf("%v", err)
	}

	if !put {
		t.Fatalf("Cache key not putted")
	}

	if len(LFU.items) != 1 {
		t.Fatalf("%d is wrong items length", len(LFU.items))
	}
}

func TestGetItem(t *testing.T) {
	item, err := LFU.Get("foo")

	if err != nil {
		t.Fatalf("%v", err)
	}

	if item != "bar" {
		t.Fatalf("%v is wrong item", item)
	}

	_, err = LFU.Get("wrong")

	if err == nil {
		t.Fatalf("missing error on wrong get")
	}
}

func TestForgetITem(t *testing.T) {
	err := LFU.Forget("foo")

	if err != nil {
		t.Fatalf("%v", err)
	}

	if len(LFU.items) != 0 {
		t.Fatalf("%d is wrong items length", len(LFU.items))
	}

	//we must also the part when item not exists
	err = LFU.Forget("notexists")

	if err != nil {
		t.Fatalf("%v", err)
	}
}

func BenchmarkPutITem(b *testing.B) {

}

func BenchmarkGetItem(b *testing.B) {

}

func BenchmarkForgetITem(b *testing.B) {

}
