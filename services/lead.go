package services

import (
	"fmt"
	"time"

	"github.com/davidalvarez305/yd_cocktails/constants"
	"github.com/davidalvarez305/yd_cocktails/database"
	"github.com/davidalvarez305/yd_cocktails/helpers"
	"github.com/davidalvarez305/yd_cocktails/types"
)

func checkSpreadsheets() {
	resp, err := GetDataFromSheets(constants.FacebookLeadsSpreadsheetID, constants.FacebookLeadsSpreadsheetRange)
	if err != nil {
		fmt.Printf("Unable to retrieve data from sheet: %v", err)
		return
	}

	var leads []types.FacebookInstantFormLead

	if len(resp.Values) > 0 {
		// Skip headers [1:]
		for _, row := range resp.Values[1:] {

			if len(row) < 16 {
				continue
			}

			lead := types.FacebookInstantFormLead{
				ID:               fmt.Sprintf("%v", row[0]),
				CreatedTime:      fmt.Sprintf("%v", row[1]),
				AdID:             fmt.Sprintf("%v", row[2]),
				AdName:           fmt.Sprintf("%v", row[3]),
				AdsetID:          fmt.Sprintf("%v", row[4]),
				AdsetName:        fmt.Sprintf("%v", row[5]),
				CampaignID:       fmt.Sprintf("%v", row[6]),
				CampaignName:     fmt.Sprintf("%v", row[7]),
				FormID:           fmt.Sprintf("%v", row[8]),
				FormName:         fmt.Sprintf("%v", row[9]),
				IsOrganic:        fmt.Sprintf("%v", row[10]),
				Platform:         fmt.Sprintf("%v", row[11]),
				FullName:         fmt.Sprintf("%v", row[12]),
				PhoneNumber:      fmt.Sprintf("%v", row[13]),
				EventDescription: fmt.Sprintf("%v", row[14]),
				Email:            fmt.Sprintf("%v", row[15]),
			}

			leads = append(leads, lead)
		}
	}

	for _, lead := range leads {
		form, err := helpers.MapInstantFormToQuoteForm(lead)
		if err != nil {
			continue
		}

		exists, err := database.IsPhoneNumberInDB(helpers.SafeString(form.PhoneNumber))

		// Skip if lead has already been saved before
		if exists || err != nil {
			continue
		}

		_, err = database.CreateLeadAndMarketing(form)

		if err != nil {
			fmt.Printf("ERROR CREATING FB LEAD FOR: %+v. MESSAGE: %+v\n", lead, err)
			continue
		}

		if !constants.Production {
			continue
		}

		for _, phoneNumber := range constants.NotificationSubscribers {

			var textMessageTemplateNotification = fmt.Sprintf(
				`NEW FACEBOOK LEAD:

				Phone: %s,
				Full Name: %s,
				Message: %s
			`, lead.PhoneNumber, lead.FullName, lead.EventDescription)

			_, err := SendTextMessage(phoneNumber, constants.CompanyPhoneNumber, textMessageTemplateNotification)

			fmt.Printf("ERROR SENDING FB LEAD AD NOTIFICATION MSG: %+v\n", err)
		}
	}
}

func archiveUnresponsiveLeads() {
	err := database.ArchivedLeadsWithLastContactOverTwoWeeks()

	if err != nil {
		fmt.Printf("ERROR ARCHIVING UNRESPONSIVE LEADS: %+v\n", err)
	}
}

func StartLeadChecker() {
	go func() {
		for {
			checkSpreadsheets()

			// Sleep for one minute before the next run
			time.Sleep(1 * time.Minute)
		}
	}()

	go func() {
		for {
			archiveUnresponsiveLeads()

			time.Sleep(24 * time.Hour)
		}
	}()
}
