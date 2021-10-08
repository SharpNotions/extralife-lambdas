package main

import (
	"encoding/json"
	"fmt"

	"github.com/aws/aws-lambda-go/lambda"
)

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

	participants, err := ListTeamLeaderboard(participant.TeamId)
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

	item := ParticipantResult{}

	item.Name = participant.Info.DisplayName
	item.FundraisingGoal = participant.Info.FundraisingGoal
	item.LastUpdated = participant.LastUpdated
	item.Links = participant.Info.Links
	item.NumDonations = participant.Info.NumDonations
	item.SumDonations = participant.Info.SumDonations
	item.Team = BuildTeamResult(team, participants)

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
