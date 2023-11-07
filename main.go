package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"runtime"

	sigar "github.com/cloudfoundry/gosigar"
)

type response struct {
	Cpu			int
	MemoryUsage	struct {
		TotalAlloc	string // TODO: float64
		Sys			string // TODO: float64
	}
	MemoryTotal	string // TODO: float64
	FsUsage		string // TODO: float64
	ExecPath	string
}

func format(val uint64) float64 {
	return float64(val) / 1024 / 1024
}

func main() {
	res := response{}
	
	// TODO: Количество дескрипторов

	// Количество ядер процессора
	res.Cpu = runtime.NumCPU()

	// Потребление памяти программой
	var mStats runtime.MemStats
	runtime.ReadMemStats(&mStats)
	res.MemoryUsage.TotalAlloc = fmt.Sprintf("%.2f", format(mStats.TotalAlloc))
	res.MemoryUsage.Sys = fmt.Sprintf("%.2f", format(mStats.Sys))
	
	// Общий объём ОЗУ
	mem := sigar.Mem{}
	if err := mem.Get(); err != nil {
		log.Println(err)
	}
	res.MemoryTotal = fmt.Sprintf("%.2f", format(mem.Total))

	// Объём корневого раздела
	fsu := sigar.FileSystemUsage{}
	if err := fsu.Get("/"); err != nil {
		log.Println(err)
	}
	res.FsUsage = fmt.Sprintf("%.2f", format(fsu.Total))

	// Путь до исполняемого файла
	ex, err := os.Executable()
	if err != nil {
		log.Println(err)
	}
	res.ExecPath = filepath.Dir(ex)
	
	fmt.Println(
		res.Cpu,
		res.MemoryUsage.Sys,
		res.MemoryUsage.TotalAlloc,
		res.MemoryTotal,
		res.FsUsage,
		res.ExecPath,
	)
}