package main

import (
	"encoding/json"
	"fmt"
	"github.com/martinlindhe/notify"
	_ "github.com/martinlindhe/notify"
	"github.com/robfig/cron/v3"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
)

type CovidCentresResponse struct {
	Centers []struct {
		CenterID     int    `json:"center_id"`
		Name         string `json:"name"`
		StateName    string `json:"state_name"`
		DistrictName string `json:"district_name"`
		BlockName    string `json:"block_name"`
		Pincode      int    `json:"pincode"`
		Lat          int    `json:"lat"`
		Long         int    `json:"long"`
		From         string `json:"from"`
		To           string `json:"to"`
		FeeType      string `json:"fee_type"`
		Sessions     []struct {
			SessionID         string   `json:"session_id"`
			Date              string   `json:"date"`
			AvailableCapacity int      `json:"available_capacity"`
			MinAgeLimit       int      `json:"min_age_limit"`
			Vaccine           string   `json:"vaccine"`
			Slots             []string `json:"slots"`
		} `json:"sessions"`
	} `json:"centers"`
}

func main() {
	if len(os.Args) > 2 {
		cowinArgs := os.Args[0:]
		searchBy, _ := strconv.Atoi(cowinArgs[1])
		date:= cowinArgs[2]
		minAge, _ := strconv.Atoi(cowinArgs[3])
		thirdParam := cowinArgs[4]
		forever, _ := strconv.ParseBool(cowinArgs[5])
		var _defaultInterval = "2s"
		if forever && len(cowinArgs)>6 {
			_defaultInterval =cowinArgs[6]
		}
		cronScheduler := cron.New()

		if searchBy == 1 {
			if forever{
				cronScheduler.AddFunc("@every "+_defaultInterval, func() { findSlotsByPinCode(date, thirdParam, minAge) })

			}else {
				findSlotsByPinCode(date,thirdParam,minAge)
			}

		}else{
			if forever{
				cronScheduler.AddFunc("@every "+_defaultInterval, func() { findSlotsByDistrictId(date,thirdParam,minAge) })

			}else{
				findSlotsByPinCode(date,thirdParam,minAge)
			}

		}
		if forever{
			cronScheduler.Start()

		}else {
			cronScheduler.Stop()
		}

	}else{
		print("Please give Arguments of this format - <searchBy (pincode = 1, district ID = 2)> <date (DD-MM-YYYY)> <minimumAge (18/45)> <pincode/districtID> <run_forever (true/false)> [interval]\n\nSee the README for more information.")
	}


	select {}
}

func findSlotsByDistrictId(date string, districtId string, minAge int) {
	getDetailsFromCowin("https://cdn-api.co-vin.in/api/v2/appointment/sessions/public/calendarByDistrict?district_id="+districtId+"&date="+date,minAge)
}

func findSlotsByPinCode(date string, pincode string, minAge int) {
	getDetailsFromCowin("https://cdn-api.co-vin.in/api/v2/appointment/sessions/public/calendarByPin?pincode="+pincode+"&date="+date,minAge)
}

func getDetailsFromCowin(url string,minAge int)  {
	cowinResponse := CovidCentresResponse{}
	var DefaultClient = &http.Client{}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil { log.Fatal(err) }
	req.Header.Add("User-Agent", "Mozilla/5.0 (X11; Linux x86_64; rv:88.0) Gecko/20100101 Firefox/88.0")
	response, err := DefaultClient.Do(req)
	if err!=nil {
		fmt.Println("Could not perform request to API. Check your internet connection/settings.")
	}else{

		responseData, err1 := ioutil.ReadAll(response.Body)
		fmt.Printf("%s\n", responseData)
		if err1==nil {
			err3 :=json.Unmarshal(responseData,&cowinResponse)
			if err3 ==nil{
				processResponseAndAlertIfPresent(cowinResponse,minAge)

			}else{
				fmt.Println("Error Processing JSON returned from API. Possible Cloudfront block.")
			}
		}else{
			fmt.Println("Unreadable response returned by API.")
		}

	}

}

func processResponseAndAlertIfPresent(response CovidCentresResponse, age int) {
	totalNoOfCenters:=len(response.Centers)
	if totalNoOfCenters>0 {
		for i := 0; i < totalNoOfCenters; i++{
			var center = response.Centers[i]
			var noOfSessionsInCenter=len(center.Sessions)
			fmt.Printf("Center Name %s \n", center.Name)
			if noOfSessionsInCenter>0 {
				for j:=0;j<noOfSessionsInCenter;j++{

					var minimumAgeLimit=center.Sessions[j].MinAgeLimit
					var availability=center.Sessions[j].AvailableCapacity
					if minimumAgeLimit== age && availability>5 {
						alertMessage:= fmt.Sprintf("Center: %s , Date: %s , Slot Info: %s Count: %d", center.Name, center.Sessions[j].Date, center.Sessions[j].Slots,center.Sessions[j].AvailableCapacity)
						fmt.Printf("Date %s \n", center.Sessions[j].Date)
						fmt.Printf("Slot Information %s \n",center.Sessions[j].Slots)
						fmt.Printf("Count Available: %d Minimum Age: %d\n",availability,minimumAgeLimit)
						notify.Alert("Slot Alert", "Slot Available", alertMessage, "https://png.pngtree.com/element_pic/17/04/07/628c04fea84856c8d04b3878eb989009.jpg")
					}else{
						fmt.Printf("No Slots Available for given date | pincode | district !  \n")
					}
				}
			}
			fmt.Printf("--------------------\n")
		}
	}else{
		fmt.Println("No Centers Found")
	}
}
