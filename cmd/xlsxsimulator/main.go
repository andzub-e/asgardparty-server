package main

import (
	"bitbucket.org/electronicjaw/asgardparty-server/internal/game"
	"bitbucket.org/electronicjaw/asgardparty-server/internal/simulation"
	"bitbucket.org/electronicjaw/asgardparty-server/pkg/rng"
	"bitbucket.org/electronicjaw/asgardparty-server/utils"
	"fmt"
	"os"
	"path/filepath"
	"time"
)

func main() {
	utils.PatchRand(rng.NewMockClient())

	conf := Config{Count: 1000 * 1000 * 10, SimulationCount: 1}

	for i := 0; i < conf.SimulationCount; i++ {
		Sim(conf)
	}
}

func Sim(conf Config) {
	simulator := simulation.NewSimulatorService()

	finRep := simulator.
		WithProgressListener(func(percent float64) {
			fmt.Printf("Processing: %.0f%%\n", percent*100)
		}).
		Simulate("Asgard", conf.Count, 100)

	reportPages := []utils.Page{
		{
			Name:  "Report",
			Table: utils.Transpose(utils.ExtractTable([]*simulation.Result{finRep}, "xlsx")),
		}}

	excel, err := utils.ExportMultiPageXLSX(reportPages)
	if err != nil {
		fmt.Printf("simulate: %v\n", err)
	}

	abs, err := filepath.Abs("results")
	if err != nil {
		panic(err)
	}

	file, err := os.Create(filepath.Join(abs, fmt.Sprintf("report-rtp-%v-%v.xlsx", finRep.RTP[:len(finRep.RTP)-1], time.Now().UnixNano())))
	if err != nil {
		fmt.Printf("simulate: %v\n", err)
	}

	if err = excel.Write(file); err != nil {
		fmt.Printf("simulate: %v\n", err)
	}
}

type Config struct {
	Count           int           `json:"count"`
	SimulationCount int           `json:"simulation_count"`
	Figures         []game.Figure `json:"figures"`
}
