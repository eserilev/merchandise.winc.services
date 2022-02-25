package campaigns

import (
	"log"
	"strconv"
	"strings"
	"time"
)

const campaignNameColumn int = 0
const couponCodeColumn int = 1
const startDateColumn int = 10
const endDateColumn int = 11
const violatorCopyColumn int = 12

type Campaign struct {
	Delete    bool            `json:"delete"`
	Replace   bool            `json:"replace"`
	Default   bool            `json:"default"`
	Campaign  string          `json:"campaign"`
	StartDate string          `json:"startDate"`
	EndDate   string          `json:"endDate"`
	Status    string          `json:"status"`
	Content   CampaignContent `json:"content"`
}

type CampaignFile struct {
	Debug     string                     `json:"debug"`
	V         string                     `json:"v"`
	P         string                     `json:"p"`
	B         []Banner                   `json:"b"`
	C         []Card                     `json:"c"`
	Campaigns map[string]CampaignContent `json:"campaigns"`
}

type CampaignContent struct {
	V string   `json:"v"`
	P string   `json:"p"`
	B []Banner `json:"b"`
	C []Card   `json:"c"`
}

type Banner struct {
	H1 string            `json:"h1"`
	D  string            `json:"d"`
	B  string            `json:"b"`
	A  Action            `json:"a"`
	F  string            `json:"f"`
	T  int               `json:"t"`
	I  map[string]string `json:"i"`
}

type Card struct {
	H1 string            `json:"h1"`
	D  string            `json:"d"`
	B  string            `json:"b"`
	A  Action            `json:"a"`
	I  map[string]string `json:"i"`
}

type Action struct {
	M  string `json:"m"`
	Id string `json:"id"`
}

func CreateCampaign(record []string) Campaign {
	var campaign Campaign
	var campaignContent CampaignContent
	const layoutISO = "11/5/2021 0:00:00"

	if strings.Contains(record[couponCodeColumn], "default") {
		campaign.Default = true
		campaign.Campaign = "default"
	} else {
		campaign.Default = false
		campaign.Campaign = record[0]
	}

	startString := strings.Split(record[startDateColumn], "/")
	endString := strings.Split(record[endDateColumn], "/")
	sYear := strings.Split(startString[2], " ")[0]
	startYear, err := strconv.Atoi(sYear)
	if err != nil {
		log.Fatal(err)
	}

	startMonth, err := strconv.Atoi(startString[0])
	if err != nil {
		log.Fatal(err)
	}
	sMonth := time.Month(startMonth)

	startDay, err := strconv.Atoi(startString[1])
	if err != nil {
		log.Fatal(err)
	}
	start := time.Date(startYear, sMonth, startDay, 0, 0, 0, 0, time.UTC)

	eYear := strings.Split(endString[2], " ")[0]
	endYear, err := strconv.Atoi(eYear)
	if err != nil {
		log.Fatal(err)
	}

	endMonth, err := strconv.Atoi(endString[0])
	if err != nil {
		log.Fatal(err)
	}
	eMonth := time.Month(endMonth)

	endDay, err := strconv.Atoi(endString[1])
	if err != nil {
		log.Fatal(err)
	}
	end := time.Date(endYear, eMonth, endDay, 0, 0, 0, 0, time.UTC)

	campaign.StartDate = strconv.Itoa(start.Year()) + "-" + GetDoubleDigitString(int(start.Month())) + "-" + GetDoubleDigitString(start.Day())
	campaign.EndDate = strconv.Itoa(end.Year()) + "-" + GetDoubleDigitString(int(end.Month())) + "-" + GetDoubleDigitString(end.Day())
	campaign.Status = "0"
	campaign.Replace = true
	campaignContent.V = record[violatorCopyColumn]
	campaignContent.B = make([]Banner, 0)
	campaignContent.C = make([]Card, 0)
	campaignContent.P = record[couponCodeColumn]
	campaign.Content = CreateCampaignContent(record[violatorCopyColumn], record[couponCodeColumn])

	return campaign
}

func CreateCampaignContent(violator string, coupon string) CampaignContent {
	campaignContent := new(CampaignContent)
	campaignContent.V = violator
	campaignContent.B = make([]Banner, 0)
	campaignContent.C = make([]Card, 0)
	campaignContent.P = coupon
	return *campaignContent
}
