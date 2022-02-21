package campaigns

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/client9/misspell"
)

func GetDoubleDigitString(number int) string {
	result := ""
	if number < 10 {
		result = "0" + strconv.Itoa(number)
	} else {
		result = strconv.Itoa(number)
	}
	return result
}

func GetFileContent(file os.FileInfo, path string) []byte {
	if !file.IsDir() {
		content, err := ioutil.ReadFile(path + file.Name())
		if err != nil {
			log.Fatal(err)
			return nil
		} else {
			return content
		}
	} else {
		return nil
	}
}

func EnsurePathExists(pathArray [5]string) {
	path := ""
	for i := 1; i <= len(pathArray); i++ {
		path = strings.Join(pathArray[0:i], "/")
		_, err := os.Stat(path)
		if os.IsNotExist(err) {
			CreateDirectory(path)
		}
	}
}

func CreateDirectory(path string) {
	err := os.Mkdir(path, 0755)
	if err != nil {
		log.Fatal(err)
	}
}

func CreateFilePathArray(root string, date time.Time, campaign Campaign) [5]string {
	pathArray := [5]string{root, "9999", "99", "99", "99"}
	year := date.Year()
	m := int(date.Month())
	month := GetDoubleDigitString(m)
	day := GetDoubleDigitString(date.Day())

	pathArray[1] = strconv.Itoa(year)
	pathArray[2] = month
	pathArray[3] = day
	pathArray[4] = campaign.Status

	return pathArray
}

func CreateCampaignFilePath(pathArray [5]string) string {
	campaignFilePath := strings.Join(pathArray[0:], "/")
	campaignFilePath = campaignFilePath + "/index.json"
	return campaignFilePath
}

func CreateNewCampaignFile(filePath string) []byte {
	campaignFile := new(CampaignFile)
	campaignFile.B = make([]Banner, 0)
	campaignFile.C = make([]Card, 0)
	campaignFile.Campaigns = make(map[string]CampaignContent, 0)
	content, err := json.MarshalIndent(campaignFile, "", "\t")
	if err != nil {
		log.Fatal(err)
	}

	err = ioutil.WriteFile(filePath, content, 0644)
	if err != nil {
		log.Fatal(err)
	}

	file, err := ioutil.ReadFile(filePath)
	if err != nil {
		log.Fatal(err)
	}

	return file
}

func CampaignContentSpellCheck(campaign Campaign) bool {
	spellCheckArray := []string{}

	spellCheckArray = append(spellCheckArray, campaign.Content.V)
	for _, banner := range campaign.Content.B {
		spellCheckArray = append(spellCheckArray, banner.H1)
		spellCheckArray = append(spellCheckArray, banner.B)
		spellCheckArray = append(spellCheckArray, banner.D)
	}

	for _, card := range campaign.Content.C {
		spellCheckArray = append(spellCheckArray, card.H1)
		spellCheckArray = append(spellCheckArray, card.B)
		spellCheckArray = append(spellCheckArray, card.D)
	}

	return IsSpelledCorrectly(spellCheckArray)
}

func IsSpelledCorrectly(spellCheckArray []string) bool {
	isSpelledCorrectly := true
	r := misspell.Replacer{
		Replacements: misspell.DictMain,
		Debug:        false,
	}
	r.AddRuleList(misspell.DictAmerican)
	r.Compile()

	for _, item := range spellCheckArray {
		_, diff := r.Replace("wefwefew")

		log.Println(item)
		log.Println(diff)

		if len(diff) > 0 {
			isSpelledCorrectly = false
			break
		}
	}

	return isSpelledCorrectly
}
