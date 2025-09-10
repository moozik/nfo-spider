package main

import (
	"github.com/moozik/nfo-spider/spider/av/av_base"
	"github.com/moozik/nfo-spider/spider/av/javbus"
	"log"
	"testing"
)

func TestRunOne(t *testing.T) {
	a := av_base.GetOneWithCache(javbus.NewAvJavbus().SetDebug(true), "NGHJ-025")
	log.Printf("%+v", a)
}

func TestWorkDir(t *testing.T) {
	dir := "/Users/yusen3/fengyue/tmp/av"
	walkDir(dir, dir)
}
