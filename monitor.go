package main

import (
	"context"
	"errors"
	"fmt"
	"os"
	"runtime"
	"time"

	"github.com/allegro/bigcache/v3"
	cacheStoreCache "github.com/eko/gocache/cache"
	cacheStore "github.com/eko/gocache/store"
	"github.com/patrickmn/go-cache"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func monitorRuntime() {
	m := &runtime.MemStats{}
	f, err := os.Create(fmt.Sprintf("mmem_%s.csv"))
	if err != nil {
		panic(err)
	}
	f.WriteString("Time;Allocated;Total Allocated; System Memory;Num Gc;Heap Allocated;Heap System;Heap Objects;Heap Released;\n")
	for {
		runtime.ReadMemStats(m)
		f.WriteString(fmt.Sprintf("%s;%d;%d;%d;%d;%d;%d;%d;%d;\n", m.Alloc, m.TotalAlloc, m.Sys, m.NumGC, m.HeapAlloc, m.HeapSys, m.HeapObjects, m.HeapReleased))
		time.Sleep(5 * time.Second)
	}
}