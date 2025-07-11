package termprinter

import (
	"fmt"
	"os"

	"github.com/August-Brandt/EgoLottery/config"
	"github.com/August-Brandt/EgoLottery/gitstats"

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
		chart := createChart(config.TimeAgo, 2, map[int]int{})
		fmt.Println(chart.View())
		return nil
	}

	// Find max number of commits
	max := -1
	max_time := -1
	for time, numCommits := range commitGroups{
		if time > max_time {
			max_time = time
		}
		if numCommits > max {
			max = numCommits
		}
	}

	chart := createChart(config.TimeAgo, max, commitGroups)
	if output == "" {
		fmt.Println(chart.View())
	} else {
		err := os.WriteFile(output, []byte(chart.View()), 0666)
		if err != nil {
			panic(err)
		}
	}
	return nil
}

func createChart(maxtime, maxval int, commits map[int]int) streamlinechart.Model {
	slc := streamlinechart.New(maxtime+len(fmt.Sprintf("%d", maxval))+1, maxval)
	for i := 0; i <= maxtime; i++ {
		commits, err  := commits[i]
		if !err {
			slc.Push(0)
		} else {
			slc.Push(float64(commits))
		}
	}
	slc.Draw()

	return slc
}