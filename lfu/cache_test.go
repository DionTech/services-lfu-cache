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
	actualHeap := LFU.heap

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

	if len(LFU.tracking) != 1 {
		t.Fatalf("%d is wrong tracking length", len(LFU.items))
	}

	tracking, exists := LFU.tracking["foo"]

	if !exists {
		t.Fatalf("tracking key not exists")
	}

	if tracking.heap == 0 {
		t.Fatal("missed heap allocation")
	}

	if actualHeap >= LFU.heap {
		t.Fatalf("heap is not calculated right; before: %d, after: %d", actualHeap, LFU.heap)
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

	if track := LFU.tracking["foo"].hits; track != 1 {
		t.Fatalf("tracking not working; expected 1, got %d", track)
	}

	//manually delete the track item to have coverage at this if statement
	delete(LFU.tracking, "foo")

	LFU.Get("foo")
	if track := LFU.tracking["foo"].hits; track != 1 {
		t.Fatalf("tracking not working; expected 1, got %d", track)
	}

	_, err = LFU.Get("wrong")

	if err == nil {
		t.Fatal("missing error on wrong get")
	}

	if _, exists := LFU.tracking["wrong"]; exists {
		t.Fatal("wrong tracking")
	}
}

func TestForgetITem(t *testing.T) {
	actualHeap := LFU.heap

	err := LFU.Forget("foo")

	if err != nil {
		t.Fatalf("%v", err)
	}

	if len(LFU.items) != 0 {
		t.Fatalf("%d is wrong items length", len(LFU.items))
	}

	if actualHeap <= LFU.heap {
		t.Fatalf("there is n heap reducing; before: %d, after: %d", actualHeap, LFU.heap)
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
