package main

import (
	"saga/api/websvc"
	"saga/orchestrator"
	"saga/storage/inmemory"
	"saga/tasks"
)

const addr = ":8080"

func main() {
	jobsSequence := tasks.PrepareSequence()

	storage := inmemory.NewStorage()
	o := orchestrator.Init(storage, jobsSequence)

	cfg := websvc.NewCfg(addr, o)
	websvc.Init(cfg)

}
