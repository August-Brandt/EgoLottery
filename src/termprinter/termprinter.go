package termprinter

import (
	"fmt"

	"github.com/August-Brandt/EgoLottery/gitstats"

	"github.com/NimbleMarkets/ntcharts/linechart/streamlinechart"
)

func PrintGraph(repos []*gitstats.Repo) error {
	// Count the collective number of commits on each day
	commitGroups := make(map[int]int)
	for _, repo := range repos {
		for timeAgo, value := range repo.Commits {
			commitGroups[timeAgo] += value
		}
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
	// fmt.Printf("Max %d\n", max)

	slc := streamlinechart.New(max_time+len(fmt.Sprintf("%d", max))+1, max)
	for i := 0; i <= max_time; i++ {
		commits, err  := commitGroups[i]
		if !err {
			slc.Push(0)
		} else {
			slc.Push(float64(commits))
		}
	}
	slc.Draw()

    fmt.Println(slc.View())
	return nil
}
