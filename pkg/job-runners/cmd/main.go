package main

import (
	"fmt"
	"github.com/go-gota/gota/dataframe"
	"github.com/kniren/gota/series"
	"github.com/mmacheerpuppy/GoPlayground/pkg/job-runners/internal/interfaces/runners"
	"github.com/mmacheerpuppy/GoPlayground/pkg/job-runners/internal/processors"
	"log"
	"os"
)

func main() {
	toStringTransformations := []runners.ToStringJob{
		func() (string, error) {
			return "s1", nil
		},
		func() (string, error) {
			return "s2" + " Manipulated", nil
		},
		func() (string, error) {
			return "s3", nil
		},
	}
	batchStringProcessor := processors.NewBatchStringProcessor()
	batchStringProcessor.AddJobs(toStringTransformations)
	jobResults, err := batchStringProcessor.Run()
	if err != nil {
		log.Fatal("Parallel execution failed with: " + err.Error())
	}

	df := dataframe.New(series.New(jobResults, series.String, "Results"))
	fmt.Println(df)
	os.Exit(1)
}
