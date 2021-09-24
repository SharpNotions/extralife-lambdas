package main

import (
	"encoding/json"
	"fmt"

	"github.com/aws/aws-lambda-go/lambda"
)

type GetResult struct {
	LastUpdated     string  `json:"lastUpdated"`
	Name            string  `json:"displayName"`
	FundraisingGoal float64 `json:"fundraisingGoal"`
	Links           []Link  `json:"links"`
	SumDonations    float64 `json:"sumDonations"`
	NumDonations    int     `json:"numDonations"`
	Team            struct {
		Name            string  `json:"name"`
		NumParticipants int     `json:"numParticipants"`
		FundraisingGoal float64 `json:"fundraisingGoal"`
		Links           []Link  `json:"links"`
		SumDonations    float64 `json:"sumDonations"`
		NumDonations    int     `json:"numDonations"`
	} `json:"team"`
}

// Handler is our lambda handler invoked by the `lambda.Start` function call
func Handler(ctx Request) (Response, error) {
	// Make the call to the DAO with params found in the path
	fmt.Println("Path vars: ", ctx.PathParameters["participantId"])
	participant, err := GetParticipantById(ctx.PathParameters["participantId"])
	if err != nil {
		fmt.Println("Failed to find Participant, ", err)
		body, _ := json.Marshal(map[string]interface{}{
			"message": "Can not find Participant",
		})
		return Response{Body: string(body), StatusCode: 404,
			Headers: map[string]string{
				"Content-Type": "application/json",
			}}, err

	}

	team, err := GetTeamById(participant.TeamId)
	if err != nil {
		fmt.Println("Failed to find Team, ", err)
		body, _ := json.Marshal(map[string]interface{}{
			"message": "Can not find Team",
		})
		return Response{Body: string(body), StatusCode: 404,
			Headers: map[string]string{
				"Content-Type": "application/json",
			}}, err

	}

	item := GetResult{}

	item.Name = participant.Info.DisplayName
	item.FundraisingGoal = participant.Info.FundraisingGoal
	item.LastUpdated = participant.LastUpdated
	item.Links = participant.Info.Links
	item.NumDonations = participant.Info.NumDonations
	item.SumDonations = participant.Info.SumDonations

	item.Team.FundraisingGoal = team.Info.FundraisingGoal
	item.Team.Links = team.Info.Links
	item.Team.Name = team.Info.DisplayName
	item.Team.NumDonations = team.Info.NumDonations
	item.Team.NumParticipants = team.Info.NumParticipants
	item.Team.SumDonations = team.Info.SumDonations

	// Log and return result
	jsonItem, _ := json.Marshal(item)
	stringItem := string(jsonItem) + "\n"
	fmt.Println("Found item: ", stringItem)
	return Response{Body: stringItem, StatusCode: 200,
		Headers: map[string]string{
			"Content-Type": "application/json",
		}}, nil
}

func main() {
	lambda.Start(Handler)
}
