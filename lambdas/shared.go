package main

import (
	"sort"

	"github.com/aws/aws-lambda-go/events"
)

type Response events.APIGatewayProxyResponse

type Request events.APIGatewayProxyRequest

type TeamResult struct {
	Id              int                     `json:"id"`
	Name            string                  `json:"name"`
	NumParticipants int                     `json:"numParticipants"`
	FundraisingGoal float64                 `json:"fundraisingGoal"`
	Links           []Link                  `json:"links"`
	SumDonations    float64                 `json:"sumDonations"`
	NumDonations    int                     `json:"numDonations"`
	Leaderboard     []TeamLeaderboardResult `json:"leaderboard"`
}

type TeamLeaderboardResult struct {
	SumDonations float64 `json:"sumDonations"`
	Name         string  `json:"displayName"`
}

type ParticipantResult struct {
	LastUpdated     string     `json:"lastUpdated"`
	Name            string     `json:"displayName"`
	FundraisingGoal float64    `json:"fundraisingGoal"`
	Links           []Link     `json:"links"`
	SumDonations    float64    `json:"sumDonations"`
	NumDonations    int        `json:"numDonations"`
	Team            TeamResult `json:"team"`
	Activity        Activity   `json:"activity"`
}

type TeamLeaderboardSorter []TeamLeaderboardResult

func (a TeamLeaderboardSorter) Len() int           { return len(a) }
func (a TeamLeaderboardSorter) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a TeamLeaderboardSorter) Less(i, j int) bool { return a[i].SumDonations > a[j].SumDonations }

func BuildTeamResult(team Team, participants []Participant) TeamResult {
	results := TeamResult{}
	results.Id = team.TeamId
	results.FundraisingGoal = team.Info.FundraisingGoal
	results.Links = team.Info.Links
	results.Name = team.Info.DisplayName
	results.NumDonations = team.Info.NumDonations
	results.NumParticipants = team.Info.NumParticipants
	results.SumDonations = team.Info.SumDonations

	var items []TeamLeaderboardResult = make([]TeamLeaderboardResult, 0)
	for _, i := range participants {
		item := TeamLeaderboardResult{}
		item.Name = i.Info.DisplayName
		item.SumDonations = i.Info.SumDonations
		items = append(items, item)
	}

	sort.Sort(TeamLeaderboardSorter(items))

	if len(items) > 5 {
		items = items[0:5]
	}
	results.Leaderboard = items

	return results
}
