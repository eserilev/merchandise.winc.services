package campaigns

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
	"path"
	"strings"
	"time"

	winc_csv "github.com/eserilev/utilities.winc.services/winc_csv"
	winc_s3 "github.com/eserilev/utilities.winc.services/winc_s3"
)

var fileUpdateSet map[string][]byte

func BatchUpload(filePath string) {

	fileUpdateSet = make(map[string][]byte)

	if filePath != "" {
		campaigns := winc_csv.ReadCsv(filePath)
		for _, campaign := range campaigns[1:] {
			CreateCampaignJSON(campaign)
		}
	}

	UpdateCampaignFiles()

	fileName := path.Base(filePath)
	// unused files, let's just not do this at the moment
	// UploadCampaignFilesToS3()
	os.Rename(filePath, archiveCsvPath+fileName)
}

func UploadCampaignFilesToS3() {
	for fileName, fileContent := range fileUpdateSet {
		fullPath := strings.Split(fileName, "./")
		f := fileName

		if len(fullPath) > 1 {
			f = fullPath[1]
		}

		winc_s3.UploadFile("winc-origin-content-develop", f, fileContent)
	}
}

func UpdateCampaignFiles() {
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
			start, end := UpdateCampaignFileContent(*campaign)
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

func UpdateCampaignFileContent(campaign Campaign) (time.Time, time.Time) {
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
	campaignFilePath := CreateCampaignFilePath(pathArray)

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

	fileUpdateSet[campaignFilePath] = campaignFileJson

	err = ioutil.WriteFile(campaignFilePath, campaignFileJson, 0644)
	if err != nil {
		log.Fatal(err)
	}
}

func UpdateCampaign(campaign Campaign, pathArray [5]string) {
	campaignFile := new(CampaignFile)
	campaignContent := new(CampaignContent)
	campaignFilePath := CreateCampaignFilePath(pathArray)

	campaignFileBytes, err := ioutil.ReadFile(campaignFilePath)
	if err != nil {
		campaignFileBytes = CreateNewCampaignFile(campaignFilePath)
	}

	json.Unmarshal(campaignFileBytes, &campaignFile)

	campaignContent.V = campaign.Content.V
	campaignContent.B = campaign.Content.B
	campaignContent.C = campaign.Content.C
	campaignContent.P = campaign.Content.P

	campaignFile.Campaigns[campaign.Campaign] = *campaignContent

	campaignFileJson, err := json.MarshalIndent(campaignFile, "", "\t")
	if err != nil {
		log.Fatal(err)
	}

	fileUpdateSet[campaignFilePath] = campaignFileJson

	err = ioutil.WriteFile(campaignFilePath, campaignFileJson, 0644)
	if err != nil {
		log.Fatal(err)
	}
}

func CreateCampaignJSON(record []string) {
	campaign := CreateCampaign(record)
	fileName := campaign.StartDate + "-" + campaign.EndDate + "-" + campaign.Campaign + "-" + campaign.Status + ".json"

	if !CampaignContentSpellCheck(campaign) {
		log.Fatal("Spellcheck FAILED!: " + campaign.Campaign)
	}

	content, err := json.MarshalIndent(campaign, "", "\t")
	if err != nil {
		log.Fatal(err)
	}

	err = ioutil.WriteFile(pendingJsonPath+fileName, content, 0644)
	if err != nil {
		log.Fatal(err)
	}
}
