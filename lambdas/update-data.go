package main

import (
	"fmt"
	"log"
	"strconv"
	"time"

	"context"
	"encoding/json"

	"net/http"

	"github.com/aws/aws-lambda-go/lambda"
)

type TeamResponse struct {
	NumParticipants    int                    `json:"numParticipants"`
	FundraisingGoal    float64                `json:"fundraisingGoal"`
	EventName          string                 `json:"eventName"`
	Links              map[string]interface{} `json:"links"`
	EventID            int                    `json:"eventID"`
	SumDonations       float64                `json:"sumDonations"`
	CreatedDateUTC     string                 `json:"createdDateUTC"`
	Name               string                 `json:"name"`
	NumAwardedBadges   int                    `json:"numAwardedBadges"`
	CaptainDisplayName string                 `json:"captainDisplayName"`
	StreamIsLive       bool                   `json:"streamIsLive"`
	AvatarImageURL     string                 `json:"avatarImageURL"`
	TeamID             int                    `json:"teamID"`
	SumPledges         float64                `json:"sumPledges"`
	NumDonations       int                    `json:"numDonations"`
}

type ParticipantResponse struct {
	DisplayName      string                 `json:"displayName"`
	FundraisingGoal  float64                `json:"fundraisingGoal"`
	EventName        string                 `json:"eventName"`
	Links            map[string]interface{} `json:"links"`
	EventID          int                    `json:"eventID"`
	SumDonations     float64                `json:"sumDonations"`
	CreatedDateUTC   string                 `json:"createdDateUTC"`
	NumAwardedBadges int                    `json:"numAwardedBadges"`
	ParticipantID    int                    `json:"participantID"`
	NumMilestones    int                    `json:"numMilestones"`
	TeamName         string                 `json:"teamName"`
	AvatarImageURL   string                 `json:"avatarImageURL"`
	TeamID           int                    `json:"teamID"`
	NumIncentives    int                    `json:"numIncentives"`
	IsTeamCaptain    bool                   `json:"isTeamCaptain"`
	SumPledges       float64                `json:"sumPledges"`
	NumDonations     int                    `json:"numDonations"`
	StreamIsLive     bool                   `json:"streamIsLive,omitempty"`
}

func getTeamUrl(teamId int) string {
	return fmt.Sprint("https://www.extra-life.org/api/teams/", teamId)
}

func FetchTeam(team int) (TeamResponse, error) {
	//Build The URL string
	URL := getTeamUrl(team)
	//We make HTTP request using the Get function
	resp, err := http.Get(URL)
	if err != nil {
		log.Fatal("Fetch Team Failed", team, err, URL)
	}
	defer resp.Body.Close()
	//Create a variable of the same type as our model
	var response TeamResponse
	//Decode the data
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		log.Fatal("Team Response Failed to decode", team, err, URL)
	}
	//Invoke the text output function & return it with nil as the error value
	return response, nil
}

func FetchTeamParticipants(teamId int) ([]ParticipantResponse, error) {
	//Build The URL string
	URL := getTeamUrl(teamId) + "/participants"
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

func FetchParticipantActivity(participantId int) (Activity, error) {
	//Build The URL string
	URL := "https://www.extra-life.org/api/participants/" + strconv.Itoa(participantId) + "/activity"
	//We make HTTP request using the Get function
	resp, err := http.Get(URL)
	if err != nil {
		log.Fatal("Fetch Participant Activity Failed")
	}
	defer resp.Body.Close()
	//Create a variable of the same type as our model
	var response Activity
	//Decode the data
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		log.Fatal("Participant Response Failed to decode")
	}
	//Invoke the text output function & return it with nil as the error value
	return response, nil
}

func GetTeams() ([]Team, error) {
	items, err := ListTeams()
	return items, err
}

func ProcessTeam(team Team) error {
	//We Read the response body on the line below.
	teamResponse, err := FetchTeam(team.TeamId)
	if err != nil {
		return err
	}

	var links []Link = make([]Link, 0)

	for k, v := range teamResponse.Links {
		item := Link{
			Type: k,
			Link: v.(string),
		}
		links = append(links, item)
	}

	team.Info.DisplayName = teamResponse.Name
	team.Info.FundraisingGoal = teamResponse.FundraisingGoal
	team.Info.NumDonations = teamResponse.NumDonations
	team.Info.NumParticipants = teamResponse.NumParticipants
	team.Info.SumDonations = teamResponse.SumDonations
	team.Info.Links = links
	team.LastUpdated = time.Now().Format(time.RFC3339)

	_, err = UpdateTeam(team)
	if err != nil {
		return err
	}

	teamParticipants, err := FetchTeamParticipants(team.TeamId)
	if err != nil {
		return err
	}

	for _, p := range teamParticipants {

		activity, err := FetchParticipantActivity(p.ParticipantID)
		if err != nil {
			return err
		}

		item := Participant{}

		var links []Link = make([]Link, 0)
		for k, v := range p.Links {
			item := Link{
				Type: k,
				Link: v.(string),
			}
			links = append(links, item)
		}

		item.ParticipantId = p.ParticipantID
		item.TeamId = p.TeamID
		item.LastUpdated = time.Now().Format(time.RFC3339)
		item.Info.DisplayName = p.DisplayName
		item.Info.FundraisingGoal = p.FundraisingGoal
		item.Info.NumDonations = p.NumDonations
		item.Info.SumDonations = p.SumDonations
		item.Info.Links = links
		item.Activity = activity

		_, err = UpdateParticipant(item)
		if err != nil {
			return err
		}
	}

	return err
}

func UpdateData() error {
	teams, err := GetTeams()
	if err != nil {
		return err
	}

	for _, team := range teams {
		err = ProcessTeam(team)
		if err != nil {
			return err
		}
	}

	return nil
}

// Handler is our lambda handler invoked by the `lambda.Start` function call
func Handler(ctx context.Context) (Response, error) {
	err := UpdateData()
	if err != nil {
		fmt.Println("Failed to Update, ", err)
		body, _ := json.Marshal(map[string]interface{}{
			"message": "Failed Update",
		})
		return Response{Body: string(body), StatusCode: 500,
			Headers: map[string]string{
				"Content-Type": "application/json",
			}}, err
	}

	b, _ := json.Marshal(map[string]interface{}{
		"message": "Update Complete",
	})

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
