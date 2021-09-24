package main

import (
	"fmt"
	"os"
	"strconv"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
)

type Link struct {
	Type string `json:"type"`
	Link string `json:"link"`
}

type ParticipantInfo struct {
	DisplayName     string  `json:"displayName"`
	FundraisingGoal float64 `json:"fundraisingGoal"`
	SumDonations    float64 `json:"sumDonations"`
	NumDonations    int     `json:"numDonations"`
	Links           []Link  `json:"links"`
}
type Participant struct {
	ParticipantId int             `json:"participantId"`
	TeamId        int             `json:"teamId"`
	Info          ParticipantInfo `json:"info"`
	LastUpdated   string          `json:"lastUpdated"`
}

type TeamInfo struct {
	DisplayName     string  `json:"displayName"`
	FundraisingGoal float64 `json:"fundraisingGoal"`
	SumDonations    float64 `json:"sumDonations"`
	NumDonations    int     `json:"numDonations"`
	NumParticipants int     `json:"numParticipants"`
	Links           []Link  `json:"links"`
}

type Team struct {
	TeamId      int      `json:"teamId"`
	Info        TeamInfo `json:"info"`
	LastUpdated string   `json:"lastUpdated"`
}

func GetParticipantById(participantId string) (Participant, error) {
	// Build the Dynamo client object
	sess := session.Must(session.NewSession())
	svc := dynamodb.New(sess)
	item := Participant{}

	// Perform the query
	fmt.Println("Trying to read from table: ", os.Getenv("PARTICIPANTS_TABLE"))
	result, err := svc.GetItem(&dynamodb.GetItemInput{
		TableName: aws.String(os.Getenv("PARTICIPANTS_TABLE")),
		Key: map[string]*dynamodb.AttributeValue{
			"participantId": {
				N: aws.String(participantId),
			},
		},
	})
	if err != nil {
		fmt.Println(err.Error())
		return item, err
	}

	// Unmarshall the result in to an Item
	err = dynamodbattribute.UnmarshalMap(result.Item, &item)
	if err != nil {
		fmt.Println(err.Error())
		return item, err
	}

	return item, nil
}

func GetTeamById(teamId int) (Team, error) {
	// Build the Dynamo client object
	sess := session.Must(session.NewSession())
	svc := dynamodb.New(sess)
	item := Team{}

	// Perform the query
	fmt.Println("Trying to read from table: ", os.Getenv("TEAMS_TABLE"))
	result, err := svc.GetItem(&dynamodb.GetItemInput{
		TableName: aws.String(os.Getenv("TEAMS_TABLE")),
		Key: map[string]*dynamodb.AttributeValue{
			"teamId": {
				N: aws.String(strconv.Itoa(teamId)),
			},
		},
	})
	if err != nil {
		fmt.Println(err.Error())
		return item, err
	}

	// Unmarshall the result in to an Item
	err = dynamodbattribute.UnmarshalMap(result.Item, &item)
	if err != nil {
		fmt.Println(err.Error())
		return item, err
	}

	return item, nil
}

func ListTeams() ([]Team, error) {
	// Build the Dynamo client object
	sess := session.Must(session.NewSession())
	svc := dynamodb.New(sess)

	// Perform the query
	fmt.Println("Trying to read list from table: ", os.Getenv("TEAMS_TABLE"))
	result, err := svc.Scan(&dynamodb.ScanInput{
		TableName: aws.String(os.Getenv("TEAMS_TABLE")),
	})
	if err != nil {
		fmt.Println(err.Error())
		return nil, err
	}
	var items []Team = make([]Team, 0)
	for _, i := range result.Items {
		item := Team{}

		err = dynamodbattribute.UnmarshalMap(i, &item)

		if err != nil {
			fmt.Println(err.Error())
		}
		items = append(items, item)
	}

	return items, nil
}

func UpdateParticipant(participant Participant) (Participant, error) {
	// Build the Dynamo client object
	sess := session.Must(session.NewSession())
	svc := dynamodb.New(sess)
	updateParticipant := Participant{}

	item, err := dynamodbattribute.ConvertToMap(participant)
	if err != nil {
		return updateParticipant, err
	}
	// Perform the query
	fmt.Println("Trying to update participant to table: ", os.Getenv("PARTICIPANTS_TABLE"), participant)
	result, err := svc.PutItem(&dynamodb.PutItemInput{
		TableName: aws.String(os.Getenv("PARTICIPANTS_TABLE")),
		Item:      item,
	})
	if err != nil {
		fmt.Println(err.Error())
		return updateParticipant, err
	}

	// Unmarshall the result in to an Item
	err = dynamodbattribute.UnmarshalMap(result.Attributes, &updateParticipant)
	if err != nil {
		fmt.Println(err.Error())
		return updateParticipant, err
	}

	return updateParticipant, nil
}

func UpdateTeam(team Team) (Team, error) {
	// Build the Dynamo client object
	sess := session.Must(session.NewSession())
	svc := dynamodb.New(sess)
	updateTeam := Team{}

	item, err := dynamodbattribute.ConvertToMap(team)
	if err != nil {
		return updateTeam, err
	}
	// Perform the query
	fmt.Println("Trying to update Team to table: ", os.Getenv("TEAMS_TABLE"), team)
	result, err := svc.PutItem(&dynamodb.PutItemInput{
		TableName: aws.String(os.Getenv("TEAMS_TABLE")),
		Item:      item,
	})
	if err != nil {
		fmt.Println(err.Error())
		return updateTeam, err
	}

	// Unmarshall the result in to an Item
	err = dynamodbattribute.UnmarshalMap(result.Attributes, &updateTeam)
	if err != nil {
		fmt.Println(err.Error())
		return updateTeam, err
	}

	return updateTeam, nil
}
