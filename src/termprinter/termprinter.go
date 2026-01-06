package termprinter

import (
	"fmt"
	"os"

	"github.com/August-Brandt/EgoLottery/config"
	"github.com/August-Brandt/EgoLottery/gitstats"
	"golang.org/x/term"

	"github.com/NimbleMarkets/ntcharts/linechart/streamlinechart"
)

func PrintGraph(repos []*gitstats.Repo, config *config.Config, output string) error {
	// Count the collective number of commits on each time interval
	commitGroups := make(map[int]int)
	for _, repo := range repos {
		for timeAgo, value := range repo.Commits {
			commitGroups[timeAgo] += value
		}
	}
	if len(commitGroups) == 0 { // Fallback when no commits. Create empty chart
		charts := createChart(config.TimeAgo, 2, 0, map[int]int{})
		fmt.Println(charts[0].View())
		return nil
	}

	// Find max number of commits
	max := -1
	max_time := -1
	for time, numCommits := range commitGroups {
		if time > max_time {
			max_time = time
		}
		if numCommits > max {
			max = numCommits
		}
	}

	width, _, err := term.GetSize(0)
	if err != nil {
		return err
	}

	charts := createChart(config.TimeAgo, max, width, commitGroups)
	if output == "" {
		for _, chart := range charts {
			fmt.Println()
			fmt.Println(chart.View())
		}
	} else {
		file, err := os.Create(output)
		if err != nil {
			return err
		}
		for _, chart := range charts {
			_, err = file.Write([]byte(chart.View()))
			if err != nil {
				return err
			}
			_, err = file.Write([]byte("\n"))
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func createChart(maxtime, maxval, width int, commits map[int]int) []*streamlinechart.Model {
	sizeOfYLabels := len(fmt.Sprintf("%d", maxval)) + 1
	fullChartXDim := maxtime + sizeOfYLabels
	chartYDim := maxval
	fmt.Printf("ChartXDim: %d | width: %d\n", fullChartXDim, width)
	
	numOfCharts := fullChartXDim / width
	commitsPrChart := width - sizeOfYLabels

	charts := make([]*streamlinechart.Model, numOfCharts+1)
	for i := 0; i < numOfCharts; i++ { // Create all full sized charts
		chart := streamlinechart.New(width, chartYDim)
		for j := i * commitsPrChart; j < (i+1)*commitsPrChart; j++ {
			commit, exists := commits[j]
			if !exists {
				chart.Push(0)
			} else {
				chart.Push(float64(commit))
			}
		}
		chart.Draw()
		charts[i] = &chart
	}
	// Add chart with remainder
	fmt.Printf("ChartXDim: %d | width: %d | remainder: %d\n", fullChartXDim, width, fullChartXDim % width)
	if remainder := fullChartXDim % width; remainder != 0 {
		fullChartXDim = remainder + sizeOfYLabels
		chart := streamlinechart.New(fullChartXDim, chartYDim)
		for i := numOfCharts * commitsPrChart; i <= maxtime; i++ {
			commit, exists := commits[i]
			if !exists {
				chart.Push(0)
			} else {
				chart.Push(float64(commit))
			}
		}
		chart.Draw()
		charts[numOfCharts] = &chart
	}
	fmt.Println(charts)

	// slc := streamlinechart.New(chartXDim, chartYDim) // Set dimensions for linechart
	// for i := 0; i <= maxtime; i++ {
	// 	commits, err  := commits[i]
	// 	if !err {
	// 		slc.Push(0)
	// 	} else {
	// 		slc.Push(float64(commits))
	// 	}
	// }
	// slc.Draw()

	return charts
}
