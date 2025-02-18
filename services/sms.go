package services

import (
	"fmt"
	"time"

	"github.com/davidalvarez305/yd_cocktails/constants"
	"github.com/davidalvarez305/yd_cocktails/database"
)

func checkSMS() {
	unreadMessages, err := database.GetUnreadMessagesInLast5Minutes()

	if err != nil {
		fmt.Printf("ERROR TRYING TO GET UNREAD MESSAGES IN LAST 5 MINUTES: %+v\n", err)
		return
	}

	if unreadMessages > 0 {
		subject := fmt.Sprintf("%d UNREAD MESSAGES", unreadMessages)
		recipients := []string{constants.CompanyEmail}

		body := fmt.Sprintf("Content-Type: text/html; charset=UTF-8\r\n%s", "LINK: "+fmt.Sprintf("%s/crm/message", constants.RootDomain))
		err = SendGmail(recipients, subject, constants.CompanyEmail, body)
		if err != nil {
			fmt.Printf("ERROR SENDING UNREAD MESSAGES NOTIFICATION EMAIL: %+v\n", err)
			return
		}
	}
}

func StartSMSNotifcationChecker() {
	go func() {
		for {
			checkSMS()

			time.Sleep(5 * time.Minute)
		}
	}()
}
