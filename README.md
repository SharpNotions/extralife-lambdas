# extra life lamabdas

## Goal

Every Year Extra Life is DDoS'ed and data just is impossible to get. These lambdas have the goal to get the latest information from the Extra Life API and then store it in Dynamodb. Then when something asks for the data we can easily return it from the stored data.

## Why Go

Why Not?

## Plan

A few Lambdas and call it done

1. Update Data - this function will attempt to grab the latest information from Extra Life and store the data
1. Get By Participant Id - this function will get and format the latest data from Dynamodb for that Participant

### Data

At the moment the data _I think_ we need is something like this

```json
{
  "lastUpdated": "2021-09-20T17:07:25.3+0000",
  "displayName": "Shawn Carr",
  "fundraisingGoal": 1000.0,
  "links": {
    "donate": "https://www.extra-life.org/index.cfm?fuseaction=donorDrive.participant&participantID=451669#donate",
    "page": "https://www.extra-life.org/index.cfm?fuseaction=donorDrive.participant&participantID=451669",
    ...
  },
  "sumDonations": 350.0,
  "numDonations": 6,
  "team": {
    "name": "Sharp Notions, LLC",
    "numParticipants": 2,
    "fundraisingGoal": 2500.0,
    "links": {
      "stream": "https://player.twitch.tv/?channel=sharpnotions",
      "page": "https://www.extra-life.org/index.cfm?fuseaction=donorDrive.team&teamID=56470"
    },
    "sumDonations": 355.0,
    "numDonations": 7
  }
}

```