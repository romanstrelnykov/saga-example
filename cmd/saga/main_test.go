package main

import (
	"fmt"
	"saga/orchestrator"
	"saga/storage/inmemory"
	"saga/tasks"
	"testing"
	"time"
)

func TestParallel(t *testing.T) {
	jobsSequence := tasks.PrepareSequence()
	fmt.Println(jobsSequence)

	storage := inmemory.NewStorage()
	o := orchestrator.Init(storage, jobsSequence)

	t.Run("1", func(t *testing.T) {
		t.Parallel()
		err := o.StartProcess("111")
		if err != nil {
			t.Error(err)
		}
	})
	t.Run("2", func(t *testing.T) {
		t.Parallel()
		err := o.StartProcess("222")
		if err != nil {
			t.Error(err)
		}
	})
	t.Run("3", func(t *testing.T) {
		t.Parallel()
		err := o.StartProcess("333")
		if err != nil {
			t.Error(err)
		}
	})

}

func TestSameRequestID(t *testing.T) {
	jobsSequence := tasks.PrepareSequence()
	fmt.Println(jobsSequence)

	storage := inmemory.NewStorage()
	o := orchestrator.Init(storage, jobsSequence)
	t.Run("1", func(t *testing.T) {
		t.Parallel()
		err := o.StartProcess("111")
		if err != nil {
			t.Error(err)
		}
	})

	t.Run("2", func(t *testing.T) {
		t.Parallel()
		time.Sleep(3 * time.Second)
		err := o.StartProcess("111")
		if err == nil {
			t.Errorf("should throw an error")
		}
	})
}

func TestStopBefore5(t *testing.T) {
	jobsSequence := tasks.PrepareSequence()
	fmt.Println(jobsSequence)

	storage := inmemory.NewStorage()
	o := orchestrator.Init(storage, jobsSequence)
	t.Run("1", func(t *testing.T) {
		t.Parallel()
		err := o.StartProcess("111")
		if err == nil {
			t.Errorf("should throw an error")
		}
	})

	t.Run("2", func(t *testing.T) {
		t.Parallel()
		time.Sleep(3 * time.Second)
		err := o.StopProcess("111")
		if err != nil {
			t.Error(err)
		}
	})

}

func TestStopAfter5(t *testing.T) {
	jobsSequence := tasks.PrepareSequence()
	fmt.Println(jobsSequence)

	storage := inmemory.NewStorage()
	o := orchestrator.Init(storage, jobsSequence)

	t.Run("1", func(t *testing.T) {
		t.Parallel()
		err := o.StartProcess("111")
		if err == nil {
			t.Errorf("should throw an error")
		}
	})

	t.Run("2", func(t *testing.T) {
		t.Parallel()
		time.Sleep(7 * time.Second)
		err := o.StopProcess("111")
		if err != nil {
			t.Error(err)
		}
	})

}
