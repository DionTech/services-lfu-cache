package lfu

import (
	"reflect"
	"testing"
)

func TestInit(t *testing.T) {
	// No need to check whether `recover()` is nil. Just turn off the panic to can later test a wrong min size.
	defer func() { recover() }()

	myStore := Init(500000)

	isCache := func(t interface{}) bool {
		switch t.(type) {
		case *Cache:
			return true
		default:
			return false
		}
	}

	if !isCache(&myStore) {
		t.Fatalf("wrong init type")
	}

	if myStore.size != 500000 {
		t.Fatalf("wrong cache size %d", myStore.size)
	}

	//we ware testing the panic: when code after Init() can be reached, something went wrong, so we manually quit with t.Fatalf()
	Init(5)
	t.Fatalf("there is a minimum of 500000")
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

func TestReduce(t *testing.T) {
	LFU.Put("foo", "bar")
	LFU.Put("foo1", "Lorem ipsum dolor sit amet, consetetur sadipscing elitr, sed diam nonumy eirmod tempor invidunt ut labore et dolore magna aliquyam erat, sed diam voluptua. At vero eos et accusam et justo duo dolores et ea rebum. Stet clita kasd gubergren, no sea takimata sanctus est Lorem ipsum dolor sit amet. Lorem ipsum dolor sit amet, consetetur sadipscing elitr, sed diam nonumy eirmod tempor invidunt ut labore et dolore magna aliquyam erat, sed diam voluptua. At vero eos et accusam et justo duo dolores et ea rebum. Stet clita kasd gubergren, no sea takimata sanctus est Lorem ipsum dolor sit amet.")
	LFU.Put("foo2", "bar ist aber nicht immer das beste")

	max := LFU.heap

	LFU.Put("foo3", "Lorem ipsum dolor sit amet, consetetur sadipscing elitr, sed diam nonumy eirmod tempor invidunt ut labore et dolore magna aliquyam erat, sed diam voluptua. At vero eos et accusam et justo duo dolores et ea rebum. Stet clita kasd gubergren, no sea takimata sanctus est Lorem ipsum dolor sit amet. Lorem ipsum dolor sit amet, consetetur sadipscing elitr, sed diam nonumy eirmod tempor invidunt ut labore et dolore magna aliquyam erat, sed diam voluptua. At vero eos et accusam et justo duo dolores et ea rebum. Stet clita kasd gubergren, no sea takimata sanctus est Lorem ipsum dolor sit amet.")
	LFU.Put("foo4", "Lorem ipsum dolor sit amet, consetetur sadipscing elitr, sed diam nonumy eirmod tempor invidunt ut labore et dolore magna aliquyam erat, sed diam voluptua. At vero eos et accusam et justo duo dolores et ea rebum. Stet clita kasd gubergren, no sea takimata sanctus est Lorem ipsum dolor sit amet. Lorem ipsum dolor sit amet, consetetur sadipscing elitr, sed diam nonumy eirmod tempor invidunt ut labore et dolore magna aliquyam erat, sed diam voluptua. At vero eos et accusam et justo duo dolores et ea rebum. Stet clita kasd gubergren, no sea takimata sanctus est Lorem ipsum dolor sit amet.")
	LFU.Put("foo5", "bar ist aber nicht immer das beste")
	LFU.Put("foo6", "Lorem ipsum dolor sit amet, consetetur sadipscing elitr, sed diam nonumy eirmod tempor invidunt ut labore et dolore magna aliquyam erat, sed diam voluptua. At vero eos et accusam et justo duo dolores et ea rebum. Stet clita kasd gubergren, no sea takimata sanctus est Lorem ipsum dolor sit amet. Lorem ipsum dolor sit amet, consetetur sadipscing elitr, sed diam nonumy eirmod tempor invidunt ut labore et dolore magna aliquyam erat, sed diam voluptua. At vero eos et accusam et justo duo dolores et ea rebum. Stet clita kasd gubergren, no sea takimata sanctus est Lorem ipsum dolor sit amet.")

	LFU.tracking["foo"] = Tracking{hits: 200, heap: 16}
	LFU.tracking["foo1"] = Tracking{hits: 21, heap: 16}
	LFU.tracking["foo2"] = Tracking{hits: 2000, heap: 16}
	LFU.tracking["foo3"] = Tracking{hits: 21, heap: 16}
	LFU.tracking["foo4"] = Tracking{hits: 23, heap: 16}
	LFU.tracking["foo5"] = Tracking{hits: 2, heap: 16}
	LFU.tracking["foo6"] = Tracking{hits: 18, heap: 16}

	keys := LFU.getSortedTrackingKeys()
	awaited := []string{"foo2", "foo", "foo4", "foo1", "foo3", "foo6", "foo5"}

	if !reflect.DeepEqual(keys, awaited) {
		t.Fatalf("not well ordered %v, awaited %v", keys, awaited)
	}

	//three items must be still available
	LFU.reduce(max)

	if LFU.heap >= max {
		t.Fatalf("not correct reducing: %d, max: %d", LFU.heap, max)
	}

	if len(LFU.items) == 7 {
		t.Fatal("no items reduced")
	}
}

func BenchmarkPutITem(b *testing.B) {

}

func BenchmarkGetItem(b *testing.B) {

}

func BenchmarkForgetITem(b *testing.B) {

}
