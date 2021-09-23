package main

import (
	"log"
	"os"

	"context"
	"encoding/json"

	"net/http"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

// Response is of type APIGatewayProxyResponse since we're leveraging the
// AWS Lambda Proxy Request functionality (default behavior)
//
// https://serverless.com/framework/docs/providers/aws/events/apigateway/#lambda-proxy-integration
type Response events.APIGatewayProxyResponse

type TeamResponse struct {
	NumParticipants int     `json:"numParticipants"`
	FundraisingGoal float64 `json:"fundraisingGoal"`
	EventName       string  `json:"eventName"`
	Links           struct {
		Stream string `json:"stream"`
		Page   string `json:"page"`
	} `json:"links"`
	EventID            int     `json:"eventID"`
	SumDonations       float64 `json:"sumDonations"`
	CreatedDateUTC     string  `json:"createdDateUTC"`
	Name               string  `json:"name"`
	NumAwardedBadges   int     `json:"numAwardedBadges"`
	CaptainDisplayName string  `json:"captainDisplayName"`
	StreamIsLive       bool    `json:"streamIsLive"`
	AvatarImageURL     string  `json:"avatarImageURL"`
	TeamID             int     `json:"teamID"`
	SumPledges         float64 `json:"sumPledges"`
	NumDonations       int     `json:"numDonations"`
}

type ParticipantResponse struct {
	DisplayName     string  `json:"displayName"`
	FundraisingGoal float64 `json:"fundraisingGoal"`
	EventName       string  `json:"eventName"`
	Links           struct {
		Donate             string `json:"donate,omitempty"`
		Page               string `json:"page,omitempty"`
		Stream             string `json:"stream,omitempty"`
		FacebookFundraiser string `json:"facebookFundraiser,omitempty"`
	} `json:"links,omitempty"`
	EventID          int     `json:"eventID"`
	SumDonations     float64 `json:"sumDonations"`
	CreatedDateUTC   string  `json:"createdDateUTC"`
	NumAwardedBadges int     `json:"numAwardedBadges"`
	ParticipantID    int     `json:"participantID"`
	NumMilestones    int     `json:"numMilestones"`
	TeamName         string  `json:"teamName"`
	AvatarImageURL   string  `json:"avatarImageURL"`
	TeamID           int     `json:"teamID"`
	NumIncentives    int     `json:"numIncentives"`
	IsTeamCaptain    bool    `json:"isTeamCaptain"`
	SumPledges       float64 `json:"sumPledges"`
	NumDonations     int     `json:"numDonations"`
	StreamIsLive     bool    `json:"streamIsLive,omitempty"`
}

type ActivityResponse struct {
	Message        string  `json:"message,omitempty"`
	Amount         float64 `json:"amount"`
	CreatedDateUTC string  `json:"createdDateUTC"`
	Title          string  `json:"title"`
	ImageURL       string  `json:"imageURL"`
	Type           string  `json:"type"`
}

type Result struct {
	Team         TeamResponse
	Participants []ParticipantResponse
	Activity     []ActivityResponse
}

func FetchTeam(team string) (TeamResponse, error) {
	//Build The URL string
	URL := "https://www.extra-life.org/api/teams/" + team
	//We make HTTP request using the Get function
	resp, err := http.Get(URL)
	if err != nil {
		log.Fatal("Fetch Team Failed")
	}
	defer resp.Body.Close()
	//Create a variable of the same type as our model
	var response TeamResponse
	//Decode the data
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		log.Fatal("Team Response Failed to decode")
	}
	//Invoke the text output function & return it with nil as the error value
	return response, nil
}

func FetchTeamParticipants(team string) ([]ParticipantResponse, error) {
	//Build The URL string
	URL := "https://www.extra-life.org/api/teams/" + team + "/participants"
	//We make HTTP request using the Get function
	resp, err := http.Get(URL)
	if err != nil {
		log.Fatal("Fetch Team Participants Failed")
	}
	defer resp.Body.Close()
	//Create a variable of the same type as our model
	var response []ParticipantResponse
	//Decode the data
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		log.Fatal("Team Participants Response Failed to decode")
	}
	//Invoke the text output function & return it with nil as the error value
	return response, nil
}

func FetchTeamActivity(team string) ([]ActivityResponse, error) {
	//Build The URL string
	URL := "https://www.extra-life.org/api/teams/" + team + "/activity"
	//We make HTTP request using the Get function
	resp, err := http.Get(URL)
	if err != nil {
		log.Fatal("Fetch Team Activity Failed")
	}
	defer resp.Body.Close()
	//Create a variable of the same type as our model
	var response []ActivityResponse
	//Decode the data
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		log.Fatal("Team Activity Response Failed to decode")
	}
	//Invoke the text output function & return it with nil as the error value
	return response, nil
}

func FetchParticipant(participant string) (ParticipantResponse, error) {
	//Build The URL string
	URL := "https://www.extra-life.org/api/participants/" + participant
	//We make HTTP request using the Get function
	resp, err := http.Get(URL)
	if err != nil {
		log.Fatal("Fetch Participant Failed")
	}
	defer resp.Body.Close()
	//Create a variable of the same type as our model
	var response ParticipantResponse
	//Decode the data
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		log.Fatal("Participant Response Failed to decode")
	}
	//Invoke the text output function & return it with nil as the error value
	return response, nil
}

func FetchParticipantActivity(participant string) ([]ActivityResponse, error) {
	//Build The URL string
	URL := "https://www.extra-life.org/api/participants/" + participant + "/activity"
	//We make HTTP request using the Get function
	resp, err := http.Get(URL)
	if err != nil {
		log.Fatal("Fetch Participant Failed")
	}
	defer resp.Body.Close()
	//Create a variable of the same type as our model
	var response []ActivityResponse
	//Decode the data
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		log.Fatal("Participant Response Failed to decode")
	}
	//Invoke the text output function & return it with nil as the error value
	return response, nil
}

// Handler is our lambda handler invoked by the `lambda.Start` function call
func Handler(ctx context.Context) (Response, error) {
	teamId := os.Getenv("EXTRA_LIFE_TEAM_ID")

	//We Read the response body on the line below.
	team, err := FetchTeam(teamId)
	if err != nil {
		return Response{StatusCode: 404}, err
	}

	teamParticipants, err := FetchTeamParticipants(teamId)
	if err != nil {
		return Response{StatusCode: 404}, err
	}

	teamActivity, err := FetchTeamActivity(teamId)
	if err != nil {
		return Response{StatusCode: 404}, err
	}

	results := Result{Team: team, Participants: teamParticipants, Activity: teamActivity}
	b, err := json.Marshal(results)

	resp := Response{
		StatusCode:      200,
		IsBase64Encoded: false,
		Body:            string(b),
		Headers: map[string]string{
			"Content-Type": "application/json",
		},
	}

	return resp, nil
}

func main() {
	lambda.Start(Handler)
}
