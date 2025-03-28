package termprinter

import (
	"fmt"

	"github.com/August-Brandt/EgoLottery/gitstats"

	"github.com/NimbleMarkets/ntcharts/linechart/streamlinechart"
)

func PrintGraph(repos []*gitstats.Repo) error {
	// Count the collective number of commits on each day
	days := make(map[int]int)
	for _, repo := range repos {
		for daysAgo, value := range repo.Commits {
			days[daysAgo] += value
		}
	}
	fmt.Println(days)
	// Find max number of commits
	max := -1
	max_day := -1
	for day, numCommits := range days{
		if day > max_day {
			max_day = day
		}
		if numCommits > max {
			max = numCommits
		}
	}
	fmt.Printf("Max %d\n", max)

	slc := streamlinechart.New(max_day, max)
	for i := 0; i <= max_day; i++ {
		commits, err  := days[i]
		if !err {
			slc.Push(0)
		} else {
			slc.Push(float64(commits))
		}
	}
	slc.Draw()
    // for _, v := range days {
    //     slc.Push(float64(v))
    // }
    // slc.Draw()

    fmt.Println(slc.View())
	return nil
}
