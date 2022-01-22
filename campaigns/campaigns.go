package campaigns

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
	"strings"
	"time"

	"github.com/eserilev/utilities.winc.services/winc_csv"
)

func Upload() {

}

func BatchUpload(filePath string) {
	campaigns := winc_csv.ReadCsv(filePath)
	for _, campaign := range campaigns[1:] {
		CreateCampaignJSON(campaign)
	}
	UploadPendingJSON()
}

func UploadPendingJSON() {
	minStartDate := new(time.Time)
	maxEndDate := new(time.Time)
	first := true
	campaign := new(Campaign)
	files, err := ioutil.ReadDir(pendingJsonPath)
	if err != nil {
		log.Fatal(err)
	}

	for _, file := range files {

		fileContent := GetFileContent(file, pendingJsonPath)
		if fileContent != nil {
			json.Unmarshal(fileContent, &campaign)
			start, end := UpdateContent(*campaign)
			if first {
				first = false
				*minStartDate = start
				*maxEndDate = end
			} else {
				if minStartDate.After(start) {
					*minStartDate = start
				}
				if maxEndDate.Before(end) {
					*maxEndDate = end
				}
			}
		}
		os.Rename(pendingJsonPath+file.Name(), archiveJsonPath+file.Name())
	}

}

func UpdateContent(campaign Campaign) (time.Time, time.Time) {
	const layoutISO = "2006-01-02"
	startDate, err := time.Parse(layoutISO, campaign.StartDate)
	if err != nil {
		log.Fatal(err)
	}

	endDate, err := time.Parse(layoutISO, campaign.EndDate)
	if err != nil {
		log.Fatal(err)
	}

	for d := startDate; d.After(endDate) == false; d = d.AddDate(0, 0, 1) {
		pathArray := CreateFilePathArray(contentRootPath, d, campaign)
		EnsurePathExists(pathArray)
		if campaign.Replace {
			if campaign.Default {
				UpdateDefault(campaign, pathArray)
			} else {
				UpdateCampaign(campaign, pathArray)
			}
		}
	}

	return startDate, endDate
}

func UpdateDefault(campaign Campaign, pathArray [5]string) {
	campaignFile := new(CampaignFile)
	campaignFilePath := strings.Join(pathArray[0:], "/")
	campaignFilePath = campaignFilePath + "/index.json"

	campaignFileBytes, err := ioutil.ReadFile(campaignFilePath)
	if err != nil {
		campaignFileBytes = CreateNewCampaignFile(campaignFilePath)
	}

	json.Unmarshal(campaignFileBytes, &campaignFile)

	campaignFile.V = campaign.Content.V
	campaignFile.B = campaign.Content.B
	campaignFile.C = campaign.Content.C
	campaignFile.P = campaign.Content.P

	campaignFileJson, err := json.MarshalIndent(campaignFile, "", "\t")

	if err != nil {
		log.Fatal(err)
	}

	err = ioutil.WriteFile(campaignFilePath, campaignFileJson, 0644)
	if err != nil {
		log.Fatal(err)
	}
}

func UpdateCampaign(campaign Campaign, pathArray [5]string) {

}

func CreateCampaignJSON(record []string) {
	campaign := CreateCampaign(record)
	fileName := campaign.StartDate + "-" + campaign.EndDate + "-" + campaign.Campaign + "-" + campaign.Status + ".json"

	content, err := json.MarshalIndent(campaign, "", "\t")
	if err != nil {
		log.Fatal(err)
	}

	err = ioutil.WriteFile(pendingJsonPath+fileName, content, 0644)
	if err != nil {
		log.Fatal(err)
	}
}
