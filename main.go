package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"runtime"
	"strconv"

	sigar "github.com/cloudfoundry/gosigar"
)

type response struct {
	Cpu			int
	MemoryUsage	struct {
		TotalAlloc	float64
		Sys			float64
	}
	MemoryTotal	float64
	FsUsage		float64
	ExecPath	string
}

func format(val uint64) float64 {
	return float64(val) / 1024 / 1024
}

func roundAndParseFloat(u uint64) (float64, error) {
	f, err := strconv.ParseFloat(fmt.Sprintf(
		"%.2f",
		format(u),
	), 64)

	return f, err
}

func main() {
	res := response{}
	var err error
	
	// TODO: Количество дескрипторов

	// Количество ядер процессора
	res.Cpu = runtime.NumCPU()

	// Потребление памяти программой
	var mStats runtime.MemStats
	runtime.ReadMemStats(&mStats)
	res.MemoryUsage.TotalAlloc, err = roundAndParseFloat(mStats.TotalAlloc)
	if err != nil {
		log.Println(err)
	}
	res.MemoryUsage.Sys, err = roundAndParseFloat(mStats.Sys)
	if err != nil {
		log.Println(err)
	}
	
	// Общий объём ОЗУ
	mem := sigar.Mem{}
	if err := mem.Get(); err != nil {
		log.Println(err)
	}
	res.MemoryTotal, err = roundAndParseFloat(mem.Total)
	if err != nil {
		log.Println(err)
	}

	// Объём корневого раздела
	fsu := sigar.FileSystemUsage{}
	if err := fsu.Get("/"); err != nil {
		log.Println(err)
	}
	res.FsUsage, err = roundAndParseFloat(fsu.Total)
	if err != nil {
		log.Println(err)
	}

	// Путь до исполняемого файла
	res.ExecPath, err = os.Executable()
	if err != nil {
		log.Println(err)
	}
	
	resJson, err := json.Marshal(res)
	if err != nil {
		log.Println(err)
	}

	fmt.Println(string(resJson))
}