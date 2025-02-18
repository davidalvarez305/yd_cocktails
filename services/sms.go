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

		body := fmt.Sprintf("Content-Type: text/html; charset=UTF-8\r\n%s", fmt.Sprintf(`
		<!DOCTYPE html>	
		<html>
				<head>Unread Messages E-mail Notification</head>
				<body>
					<p>You have <strong>%d</strong> unread messages in the last 5 minutes.</p>
					<p><a href="%s/crm/message">Click here to view your messages</a></p>
				</body>
			</html>`, unreadMessages, constants.RootDomain))

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
