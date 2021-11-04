package lfu

import (
	"reflect"
	"testing"
	"time"
)

var myCache = Init(500000)

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
	actualHeap := myCache.heap

	put, err := myCache.Put("foo", "bar")

	if err != nil {
		t.Fatalf("%v", err)
	}

	if !put {
		t.Fatalf("Cache key not putted")
	}

	if len(myCache.items) != 1 {
		t.Fatalf("%d is wrong items length", len(myCache.items))
	}

	if len(myCache.tracking) != 1 {
		t.Fatalf("%d is wrong tracking length", len(myCache.items))
	}

	tracking, exists := myCache.tracking["foo"]

	if !exists {
		t.Fatalf("tracking key not exists")
	}

	if tracking.heap == 0 {
		t.Fatal("missed heap allocation")
	}

	if actualHeap >= myCache.heap {
		t.Fatalf("heap is not calculated right; before: %d, after: %d", actualHeap, myCache.heap)
	}

	if tracking.lastUpdatedAt.IsZero() {
		t.Fatalf("last updated at is not done")
	}
}

func TestGetItem(t *testing.T) {
	lastUpdatedAt := myCache.tracking["foo"].lastUpdatedAt

	item, err := myCache.Get("foo")

	if err != nil {
		t.Fatalf("%v", err)
	}

	if item != "bar" {
		t.Fatalf("%v is wrong item", item)
	}

	if track := myCache.tracking["foo"].hits; track != 1 {
		t.Fatalf("tracking not working; expected 1, got %d", track)
	}

	if lastUpdatedAt.Equal(myCache.tracking["foo"].lastUpdatedAt) {
		t.Fatal("last updated at is not actualized on get")
	}

	//manually delete the track item to have coverage at this if statement
	delete(myCache.tracking, "foo")

	myCache.Get("foo")
	if track := myCache.tracking["foo"].hits; track != 1 {
		t.Fatalf("tracking not working; expected 1, got %d", track)
	}

	_, err = myCache.Get("wrong")

	if err == nil {
		t.Fatal("missing error on wrong get")
	}

	if _, exists := myCache.tracking["wrong"]; exists {
		t.Fatal("wrong tracking")
	}
}

func TestForgetITem(t *testing.T) {
	myCache.Put("foo", "bar")

	actualHeap := myCache.heap

	err := myCache.Forget("foo")

	if err != nil {
		t.Fatalf("%v", err)
	}

	if len(myCache.items) != 0 {
		t.Fatalf("%d is wrong items length", len(myCache.items))
	}

	if actualHeap <= myCache.heap {
		t.Fatalf("there is n heap reducing; before: %d, after: %d", actualHeap, myCache.heap)
	}

	//we must also the part when item not exists
	err = myCache.Forget("notexists")

	if err != nil {
		t.Fatalf("%v", err)
	}
}

func TestReduce(t *testing.T) {
	myCache.Put("foo", "bar")
	myCache.Put("foo1", "Lorem ipsum dolor sit amet, consetetur sadipscing elitr, sed diam nonumy eirmod tempor invidunt ut labore et dolore magna aliquyam erat, sed diam voluptua. At vero eos et accusam et justo duo dolores et ea rebum. Stet clita kasd gubergren, no sea takimata sanctus est Lorem ipsum dolor sit amet. Lorem ipsum dolor sit amet, consetetur sadipscing elitr, sed diam nonumy eirmod tempor invidunt ut labore et dolore magna aliquyam erat, sed diam voluptua. At vero eos et accusam et justo duo dolores et ea rebum. Stet clita kasd gubergren, no sea takimata sanctus est Lorem ipsum dolor sit amet.")
	myCache.Put("foo2", "bar ist aber nicht immer das beste")

	max := myCache.heap

	myCache.Put("foo3", "Lorem ipsum dolor sit amet, consetetur sadipscing elitr, sed diam nonumy eirmod tempor invidunt ut labore et dolore magna aliquyam erat, sed diam voluptua. At vero eos et accusam et justo duo dolores et ea rebum. Stet clita kasd gubergren, no sea takimata sanctus est Lorem ipsum dolor sit amet. Lorem ipsum dolor sit amet, consetetur sadipscing elitr, sed diam nonumy eirmod tempor invidunt ut labore et dolore magna aliquyam erat, sed diam voluptua. At vero eos et accusam et justo duo dolores et ea rebum. Stet clita kasd gubergren, no sea takimata sanctus est Lorem ipsum dolor sit amet.")
	myCache.Put("foo4", "Lorem ipsum dolor sit amet, consetetur sadipscing elitr, sed diam nonumy eirmod tempor invidunt ut labore et dolore magna aliquyam erat, sed diam voluptua. At vero eos et accusam et justo duo dolores et ea rebum. Stet clita kasd gubergren, no sea takimata sanctus est Lorem ipsum dolor sit amet. Lorem ipsum dolor sit amet, consetetur sadipscing elitr, sed diam nonumy eirmod tempor invidunt ut labore et dolore magna aliquyam erat, sed diam voluptua. At vero eos et accusam et justo duo dolores et ea rebum. Stet clita kasd gubergren, no sea takimata sanctus est Lorem ipsum dolor sit amet.")
	myCache.Put("foo5", "bar ist aber nicht immer das beste")
	myCache.Put("foo6", "Lorem ipsum dolor sit amet, consetetur sadipscing elitr, sed diam nonumy eirmod tempor invidunt ut labore et dolore magna aliquyam erat, sed diam voluptua. At vero eos et accusam et justo duo dolores et ea rebum. Stet clita kasd gubergren, no sea takimata sanctus est Lorem ipsum dolor sit amet. Lorem ipsum dolor sit amet, consetetur sadipscing elitr, sed diam nonumy eirmod tempor invidunt ut labore et dolore magna aliquyam erat, sed diam voluptua. At vero eos et accusam et justo duo dolores et ea rebum. Stet clita kasd gubergren, no sea takimata sanctus est Lorem ipsum dolor sit amet.")

	now := time.Now()

	//here we want to test the sorting without having a difference at the times
	myCache.tracking["foo"] = Tracking{hits: 200, heap: 16, lastUpdatedAt: now}
	myCache.tracking["foo1"] = Tracking{hits: 21, heap: 16, lastUpdatedAt: now}
	myCache.tracking["foo2"] = Tracking{hits: 2000, heap: 16, lastUpdatedAt: now}
	myCache.tracking["foo3"] = Tracking{hits: 20, heap: 16, lastUpdatedAt: now}
	myCache.tracking["foo4"] = Tracking{hits: 23, heap: 16, lastUpdatedAt: now}
	myCache.tracking["foo5"] = Tracking{hits: 2, heap: 16, lastUpdatedAt: now}
	myCache.tracking["foo6"] = Tracking{hits: 18, heap: 16, lastUpdatedAt: now}

	keys := myCache.getSortedTrackingKeys("foo5")
	awaited := []string{"foo2", "foo", "foo4", "foo1", "foo3", "foo6"}

	if !reflect.DeepEqual(keys, awaited) {
		t.Fatalf("not well ordered %v, awaited %v", keys, awaited)
	}

	//now we use different time fr foo2 - the most hits but lastUpdated was a few days before, so we will get the worst rate and so the last place in order
	myCache.tracking["foo"] = Tracking{hits: 200, heap: 16, lastUpdatedAt: now}
	myCache.tracking["foo1"] = Tracking{hits: 21, heap: 16, lastUpdatedAt: now}
	myCache.tracking["foo2"] = Tracking{hits: 2000, heap: 16, lastUpdatedAt: now.Add(100 * -24 * time.Hour)}
	myCache.tracking["foo3"] = Tracking{hits: 20, heap: 16, lastUpdatedAt: now}
	myCache.tracking["foo4"] = Tracking{hits: 23, heap: 16, lastUpdatedAt: now}
	myCache.tracking["foo5"] = Tracking{hits: 2, heap: 16, lastUpdatedAt: now}
	myCache.tracking["foo6"] = Tracking{hits: 18, heap: 16, lastUpdatedAt: now}

	keys = myCache.getSortedTrackingKeys("foo5")
	awaited = []string{"foo", "foo4", "foo1", "foo3", "foo6", "foo2"}

	if !reflect.DeepEqual(keys, awaited) {
		t.Fatalf("not well ordered %v, awaited %v", keys, awaited)
	}

	//three items must be still available
	myCache.reduce(max, "foo5")

	if myCache.heap >= max {
		t.Fatalf("not correct reducing: %d, max: %d", myCache.heap, max)
	}

	if len(myCache.items) == 7 {
		t.Fatal("no items reduced")
	}
}

func BenchmarkPutITem(b *testing.B) {

}

func BenchmarkGetItem(b *testing.B) {

}

func BenchmarkForgetITem(b *testing.B) {

}
