package main

import (
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/aws/aws-lambda-go/lambda"
)

// Handler is our lambda handler invoked by the `lambda.Start` function call
func Handler(ctx Request) (Response, error) {
	// Make the call to the DAO with params found in the path
	fmt.Println("Path vars: ", ctx.PathParameters["teamId"])
	teamId, err := strconv.Atoi(ctx.PathParameters["teamId"])
	if err != nil {
		fmt.Println("Failed to find convert value to int", err)
		body, _ := json.Marshal(map[string]interface{}{
			"message": "Failed to find convert value to int",
		})
		return Response{Body: string(body), StatusCode: 400,
			Headers: map[string]string{
				"Content-Type": "application/json",
			}}, err

	}
	team, err := GetTeamById(teamId)
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
	participants, err := ListTeamLeaderboard(teamId)
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

	item := BuildTeamResult(team, participants)

	// Log and return result
	jsonItem, _ := json.Marshal(item)
	stringItem := string(jsonItem) + "\n"
	fmt.Println("Found Team: ", stringItem)
	return Response{Body: stringItem, StatusCode: 200,
		Headers: map[string]string{
			"Content-Type": "application/json",
		}}, nil
}

func main() {
	lambda.Start(Handler)
}
