package database

import (
	"database/sql"
	"errors"
	"fmt"
	"strconv"
	"time"

	"github.com/davidalvarez305/yd_cocktails/constants"
	"github.com/davidalvarez305/yd_cocktails/models"
	"github.com/davidalvarez305/yd_cocktails/types"
	"github.com/davidalvarez305/yd_cocktails/utils"
)

func InsertCSRFToken(token models.CSRFToken) error {
	stmt, err := DB.Prepare(`INSERT INTO "csrf_token" ("expiry_time", "token", "is_used") VALUES(to_timestamp($1)::timestamptz AT TIME ZONE 'America/New_York', $2, $3)`)
	if err != nil {
		return fmt.Errorf("error preparing statement: %w", err)
	}
	defer stmt.Close()

	_, err = stmt.Exec(token.ExpiryTime, token.Token, token.IsUsed)
	if err != nil {
		return fmt.Errorf("error executing statement: %w", err)
	}

	return nil
}

func CheckIsTokenUsed(decryptedToken string) (bool, error) {
	var isUsed bool

	stmt, err := DB.Prepare(`SELECT is_used FROM "csrf_token" WHERE "token" = $1`)
	if err != nil {
		return isUsed, fmt.Errorf("error preparing statement: %w", err)
	}
	defer stmt.Close()

	row := stmt.QueryRow(decryptedToken)

	err = row.Scan(&isUsed)
	if err != nil {
		return isUsed, fmt.Errorf("error scanning row: %w", err)
	}

	return isUsed, nil
}

func CreateLeadAndMarketing(quoteForm types.QuoteForm) (int, error) {
	var leadID int
	tx, err := DB.Begin()
	if err != nil {
		return leadID, fmt.Errorf("error starting transaction: %w", err)
	}
	defer tx.Rollback()

	leadStmt, err := tx.Prepare(`
		INSERT INTO lead (first_name, last_name, phone_number, created_at, event_type_id, venue_type_id, message, opt_in_text_messaging, email, guests)
		VALUES ($1, $2, $3, to_timestamp($4)::timestamptz AT TIME ZONE 'America/New_York', $5, $6, $7, $8, $9, $10)
		RETURNING lead_id
	`)
	if err != nil {
		return leadID, fmt.Errorf("error preparing lead statement: %w", err)
	}
	defer leadStmt.Close()

	createdAt, err := utils.GetCurrentTimeInEST()
	if err != nil {
		return leadID, fmt.Errorf("error getting time as EST: %w", err)
	}

	message := utils.CreateNullString(quoteForm.Message)
	eventTypeId := utils.CreateNullInt(quoteForm.EventType)
	venueTypeId := utils.CreateNullInt(quoteForm.VenueType)

	err = leadStmt.QueryRow(
		utils.CreateNullString(quoteForm.FirstName),
		utils.CreateNullString(quoteForm.LastName),
		utils.CreateNullString(quoteForm.PhoneNumber),
		createdAt,
		eventTypeId,
		venueTypeId,
		message,
		utils.CreateNullBool(quoteForm.OptInTextMessaging),
		utils.CreateNullString(quoteForm.Email),
		utils.CreateNullInt(quoteForm.Guests),
	).Scan(&leadID)
	if err != nil {
		return leadID, fmt.Errorf("error inserting lead: %w", err)
	}

	marketingStmt, err := tx.Prepare(`
		INSERT INTO lead_marketing (lead_id, source, medium, channel, landing_page, keyword, referrer, click_id, campaign_id, ad_campaign, ad_group_id, ad_group_name, ad_set_id, ad_set_name, ad_id, ad_headline, language, user_agent, button_clicked, ip, external_id, google_client_id, csrf_secret, facebook_click_id, facebook_client_id, longitude, latitude)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17, $18, $19, $20, $21, $22, $23, $24, $25, $26, $27)
	`)
	if err != nil {
		return leadID, fmt.Errorf("error preparing marketing statement: %w", err)
	}
	defer marketingStmt.Close()

	_, err = marketingStmt.Exec(
		leadID,
		utils.CreateNullString(quoteForm.Source),
		utils.CreateNullString(quoteForm.Medium),
		utils.CreateNullString(quoteForm.Channel),
		utils.CreateNullString(quoteForm.LandingPage),
		utils.CreateNullString(quoteForm.Keyword),
		utils.CreateNullString(quoteForm.Referrer),
		utils.CreateNullString(quoteForm.ClickID),
		utils.CreateNullInt64(quoteForm.CampaignID),
		utils.CreateNullString(quoteForm.AdCampaign),
		utils.CreateNullInt64(quoteForm.AdGroupID),
		utils.CreateNullString(quoteForm.AdGroupName),
		utils.CreateNullInt64(quoteForm.AdSetID),
		utils.CreateNullString(quoteForm.AdSetName),
		utils.CreateNullInt64(quoteForm.AdID),
		utils.CreateNullInt64(quoteForm.AdHeadline),
		utils.CreateNullString(quoteForm.Language),
		utils.CreateNullString(quoteForm.UserAgent),
		utils.CreateNullString(quoteForm.ButtonClicked),
		utils.CreateNullString(quoteForm.IP),
		utils.CreateNullString(quoteForm.ExternalID),
		utils.CreateNullString(quoteForm.GoogleClientID),
		utils.CreateNullString(quoteForm.CSRFSecret),
		utils.CreateNullString(quoteForm.FacebookClickID),
		utils.CreateNullString(quoteForm.FacebookClientID),
		utils.CreateNullString(quoteForm.Longitude),
		utils.CreateNullString(quoteForm.Latitude),
	)
	if err != nil {
		return leadID, fmt.Errorf("error inserting marketing data: %w", err)
	}

	err = tx.Commit()
	if err != nil {
		return leadID, fmt.Errorf("error committing transaction: %w", err)
	}

	return leadID, nil
}

func MarkCSRFTokenAsUsed(token string) error {
	stmt, err := DB.Prepare(`UPDATE "csrf_token" SET "is_used" = true WHERE "token" = $1`)
	if err != nil {
		return fmt.Errorf("error preparing statement: %w", err)
	}
	defer stmt.Close()

	_, err = stmt.Exec(token)
	if err != nil {
		return fmt.Errorf("error executing statement: %w", err)
	}

	fmt.Println("CSRFToken marked as used successfully")
	return nil
}

func SaveSMS(msg models.Message) error {
	stmt, err := DB.Prepare(`
		INSERT INTO message (external_id, user_id, lead_id, text, date_created, text_from, text_to, is_inbound)
		VALUES ($1, $2, $3, $4, to_timestamp($5)::timestamptz AT TIME ZONE 'America/New_York', $6, $7, $8)
	`)
	if err != nil {
		return fmt.Errorf("error preparing statement: %w", err)
	}
	defer stmt.Close()

	var leadID sql.NullInt64
	if msg.LeadID != 0 {
		leadID = sql.NullInt64{Int64: int64(msg.LeadID), Valid: true}
	}

	_, err = stmt.Exec(msg.ExternalID, msg.UserID, leadID, msg.Text, msg.DateCreated, msg.TextFrom, msg.TextTo, msg.IsInbound)
	if err != nil {
		return fmt.Errorf("error executing statement: %w", err)
	}

	return nil
}

func SavePhoneCall(phoneCall models.PhoneCall) error {
	stmt, err := DB.Prepare(`
		INSERT INTO phone_call (
			external_id, user_id, lead_id, call_duration,
			date_created, call_from, call_to, is_inbound,
			recording_url, status
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)`)
	if err != nil {
		return fmt.Errorf("error preparing statement: %w", err)
	}
	defer stmt.Close()

	var leadID, callDuration sql.NullInt64

	if phoneCall.LeadID != 0 {
		leadID = sql.NullInt64{Int64: int64(phoneCall.LeadID), Valid: true}
	}
	if phoneCall.CallDuration != 0 {
		callDuration = sql.NullInt64{Int64: int64(phoneCall.CallDuration), Valid: true}
	}

	_, err = stmt.Exec(
		phoneCall.ExternalID,
		phoneCall.UserID,
		leadID,
		callDuration,
		phoneCall.DateCreated,
		phoneCall.CallFrom,
		phoneCall.CallTo,
		phoneCall.IsInbound,
		phoneCall.RecordingURL,
		phoneCall.Status,
	)
	if err != nil {
		return fmt.Errorf("error executing statement: %w", err)
	}

	return nil
}

func GetUserIDFromPhoneNumber(from string) (int, error) {
	var userId int

	stmt, err := DB.Prepare(`SELECT "user_id" FROM "user" WHERE "phone_number" = $1`)
	if err != nil {
		return userId, fmt.Errorf("error preparing statement: %w", err)
	}
	defer stmt.Close()

	row := stmt.QueryRow(from)

	err = row.Scan(&userId)
	if err != nil {
		return userId, fmt.Errorf("error scanning row: %w", err)
	}

	return userId, nil
}

func GetLeadByStripeCustomerID(stripeCustomerID string) (int, error) {
	var userId int

	stmt, err := DB.Prepare(`SELECT lead_id FROM lead WHERE stripe_customer_id = $1`)
	if err != nil {
		return userId, fmt.Errorf("error preparing statement: %w", err)
	}
	defer stmt.Close()

	row := stmt.QueryRow(stripeCustomerID)

	err = row.Scan(&userId)
	if err != nil {
		return userId, fmt.Errorf("error scanning row: %w", err)
	}

	return userId, nil
}

func GetConversionLeadInfo(leadId int) (types.ConversionLeadInfo, error) {
	var leadConversionInfo types.ConversionLeadInfo

	stmt, err := DB.Prepare(`SELECT l.lead_id, l.created_at, et.name, vt.name, l.guests
	FROM "lead" AS l
	JOIN event_type  AS et ON et.event_type_id = l.event_type_id
	JOIN venue_type AS vt ON vt.venue_type_id  = l.venue_type_id 
	WHERE l.lead_id = $1;`)

	if err != nil {
		return leadConversionInfo, fmt.Errorf("error preparing statement: %w", err)
	}
	defer stmt.Close()

	row := stmt.QueryRow(leadId)

	var createdAt time.Time
	err = row.Scan(&leadConversionInfo.LeadID,
		&createdAt,
		&leadConversionInfo.EventType,
		&leadConversionInfo.VenueType,
		&leadConversionInfo.Guests,
	)
	if err != nil {
		return leadConversionInfo, fmt.Errorf("error scanning row: %w", err)
	}

	leadConversionInfo.CreatedAt = createdAt.Unix()

	return leadConversionInfo, nil
}

func GetPhoneNumberFromUserID(userID int) (string, error) {
	var phoneNumber string

	stmt, err := DB.Prepare(`SELECT "phone_number" FROM "user" WHERE "user_id" = $1`)
	if err != nil {
		return phoneNumber, fmt.Errorf("error preparing statement: %w", err)
	}
	defer stmt.Close()

	row := stmt.QueryRow(userID)

	err = row.Scan(&phoneNumber)
	if err != nil {
		return phoneNumber, fmt.Errorf("error scanning row: %w", err)
	}

	return phoneNumber, nil
}

func GetUserById(id int) (models.User, error) {
	var user models.User

	stmt, err := DB.Prepare(`SELECT user_id, username, password, user_role_id, phone_number, first_name, last_name FROM "user" WHERE "user_id" = $1`)
	if err != nil {
		return user, fmt.Errorf("error preparing statement: %w", err)
	}
	defer stmt.Close()

	row := stmt.QueryRow(id)

	err = row.Scan(&user.UserID, &user.Username, &user.Password, &user.UserRoleID, &user.PhoneNumber, &user.FirstName, &user.LastName)
	if err != nil {
		return user, fmt.Errorf("error scanning row: %w", err)
	}

	return user, nil
}

func GetUserByUsername(username string) (models.User, error) {
	var user models.User

	stmt, err := DB.Prepare(`SELECT user_id, username, password, user_role_id, phone_number, first_name, last_name FROM "user" WHERE "username" = $1`)
	if err != nil {
		return user, fmt.Errorf("error preparing statement: %w", err)
	}
	defer stmt.Close()

	row := stmt.QueryRow(username)

	err = row.Scan(&user.UserID, &user.Username, &user.Password, &user.UserRoleID, &user.PhoneNumber, &user.FirstName, &user.LastName)
	if err != nil {
		return user, fmt.Errorf("error scanning row: %w", err)
	}

	return user, nil
}

func GetEventTypes() ([]models.EventType, error) {
	var eventTypes []models.EventType

	rows, err := DB.Query(`SELECT event_type_id, name FROM "event_type"`)
	if err != nil {
		return eventTypes, fmt.Errorf("error executing query: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var et models.EventType
		err := rows.Scan(&et.EventTypeID, &et.Name)
		if err != nil {
			return eventTypes, fmt.Errorf("error scanning row: %w", err)
		}
		eventTypes = append(eventTypes, et)
	}

	if err := rows.Err(); err != nil {
		return eventTypes, fmt.Errorf("error iterating rows: %w", err)
	}

	return eventTypes, nil
}

func GetVenueTypes() ([]models.VenueType, error) {
	var venueTypes []models.VenueType

	rows, err := DB.Query(`SELECT venue_type_id, name FROM "venue_type"`)
	if err != nil {
		return venueTypes, fmt.Errorf("error executing query: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var vt models.VenueType
		err := rows.Scan(&vt.VenueTypeID, &vt.Name)
		if err != nil {
			return venueTypes, fmt.Errorf("error scanning row: %w", err)
		}
		venueTypes = append(venueTypes, vt)
	}

	if err := rows.Err(); err != nil {
		return venueTypes, fmt.Errorf("error iterating rows: %w", err)
	}

	return venueTypes, nil
}

func GetLeadList(params types.GetLeadsParams) ([]types.LeadList, int, error) {
	var leads []types.LeadList

	query := `SELECT l.lead_id, l.first_name, l.last_name, l.phone_number, 
		l.created_at, et.name AS event_type, vt.name AS venue_type, lm.language, l.event_type_id, l.venue_type_id, l.guests,
		COUNT(*) OVER() AS total_rows
		FROM lead AS l
		JOIN event_type AS et ON et.event_type_id = l.event_type_id
		JOIN venue_type AS vt ON vt.venue_type_id = l.venue_type_id
		JOIN lead_marketing AS lm ON lm.lead_id = l.lead_id
		WHERE (et.event_type_id = $1 OR $1 IS NULL) 
		AND (vt.venue_type_id = $2 OR $2 IS NULL)
		ORDER BY l.created_at ASC
		LIMIT $3
		OFFSET $4`

	var offset int

	// Handle pagination
	if params.PageNum != nil {
		pageNum, err := strconv.Atoi(*params.PageNum)
		if err != nil {
			return nil, 0, fmt.Errorf("could not convert page num: %w", err)
		}
		offset = (pageNum - 1) * int(constants.LeadsPerPage)
	}

	rows, err := DB.Query(query, params.EventType, params.VenueType, constants.LeadsPerPage, offset)
	if err != nil {
		return nil, 0, fmt.Errorf("error executing query: %w", err)
	}
	defer rows.Close()

	var totalRows int
	for rows.Next() {
		var lead types.LeadList
		var createdAt time.Time
		var eventType, venueType sql.NullString
		var eventTypeId, venueTypeId, guests sql.NullInt64

		err := rows.Scan(&lead.LeadID,
			&lead.FirstName,
			&lead.LastName,
			&lead.PhoneNumber,
			&createdAt,
			&eventType,
			&venueType,
			&lead.Language,
			&eventTypeId,
			&venueTypeId,
			&guests,
			&totalRows)
		if err != nil {
			return nil, 0, fmt.Errorf("error scanning row: %w", err)
		}
		lead.CreatedAt = utils.FormatTimestampEST(createdAt.Unix())

		if eventTypeId.Valid {
			lead.EventTypeID = int(eventTypeId.Int64)
		}
		if venueTypeId.Valid {
			lead.VenueTypeID = int(venueTypeId.Int64)
		}
		if eventType.Valid {
			lead.EventType = eventType.String
		}
		if venueType.Valid {
			lead.VenueType = venueType.String
		}
		if guests.Valid {
			lead.Guests = int(guests.Int64)
		}

		leads = append(leads, lead)
	}

	if err := rows.Err(); err != nil {
		return nil, 0, fmt.Errorf("error iterating rows: %w", err)
	}

	return leads, totalRows, nil
}

func GetLeadDetails(leadID string) (types.LeadDetails, error) {
	query := `SELECT l.lead_id,
	l.first_name,
	l.last_name,
	l.phone_number,
	et.name,
	vt.name,
	lm.ad_campaign,
	lm.medium,
	lm.source,
	lm.referrer,
	lm.landing_page,
	lm.ip,
	lm.keyword,
	lm.channel,
	lm.language,
	l.message,
	l.email,
	lm.facebook_click_id,
	lm.facebook_client_id,
	lm.external_id,
	lm.user_agent,
	lm.click_id,
	lm.google_client_id,
	lm.button_clicked,
	lm.campaign_id,
	l.guests
	FROM lead l
	JOIN event_type et ON l.event_type_id = et.event_type_id
	JOIN venue_type vt ON l.venue_type_id = vt.venue_type_id
	JOIN lead_marketing lm ON l.lead_id = lm.lead_id
	WHERE l.lead_id = $1`

	var leadDetails types.LeadDetails

	row := DB.QueryRow(query, leadID)

	var adCampaign, medium, source, referrer, landingPage, ip, keyword, channel, language, email, facebookClickId, facebookClientId sql.NullString
	var eventType, venueType, message, externalId, userAgent, clickId, googleClientId sql.NullString
	var campaignId, guests sql.NullInt64

	var buttonClicked sql.NullString

	err := row.Scan(
		&leadDetails.LeadID,
		&leadDetails.FirstName,
		&leadDetails.LastName,
		&leadDetails.PhoneNumber,
		&eventType,
		&venueType,
		&adCampaign,
		&medium,
		&source,
		&referrer,
		&landingPage,
		&ip,
		&keyword,
		&channel,
		&language,
		&message,
		&email,
		&facebookClickId,
		&facebookClientId,
		&externalId,
		&userAgent,
		&clickId,
		&googleClientId,
		&buttonClicked,
		&campaignId,
		&guests,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return leadDetails, fmt.Errorf("no lead found with ID %s", leadID)
		}
		return leadDetails, fmt.Errorf("error scanning row: %w", err)
	}

	// Map the nullable fields to your struct
	if guests.Valid {
		leadDetails.Guests = int(guests.Int64)
	}
	if buttonClicked.Valid {
		leadDetails.ButtonClicked = buttonClicked.String
	}

	if clickId.Valid {
		leadDetails.ClickID = clickId.String
	}

	if googleClientId.Valid {
		leadDetails.GoogleClientID = googleClientId.String
	}

	if externalId.Valid {
		leadDetails.ExternalID = externalId.String
	}

	if userAgent.Valid {
		leadDetails.UserAgent = userAgent.String
	}

	if facebookClickId.Valid {
		leadDetails.FacebookClickID = facebookClickId.String
	}

	if facebookClientId.Valid {
		leadDetails.FacebookClientID = facebookClientId.String
	}

	if email.Valid {
		leadDetails.Email = email.String
	}

	if eventType.Valid {
		leadDetails.EventType = eventType.String
	}

	if venueType.Valid {
		leadDetails.VenueType = venueType.String
	}

	if adCampaign.Valid {
		leadDetails.CampaignName = adCampaign.String
	}
	if campaignId.Valid {
		leadDetails.CampaignID = campaignId.Int64
	}

	if medium.Valid {
		leadDetails.Medium = medium.String
	}

	if source.Valid {
		leadDetails.Source = source.String
	}

	if referrer.Valid {
		leadDetails.Referrer = referrer.String
	}

	if landingPage.Valid {
		leadDetails.LandingPage = landingPage.String
	}

	if ip.Valid {
		leadDetails.IP = ip.String
	}

	if keyword.Valid {
		leadDetails.Keyword = keyword.String
	}

	if channel.Valid {
		leadDetails.Channel = channel.String
	}

	if language.Valid {
		leadDetails.Language = language.String
	}

	if message.Valid {
		leadDetails.Message = message.String
	}

	return leadDetails, nil
}

func GetLeadIDFromPhoneNumber(from string) (int, error) {
	var leadId int

	stmt, err := DB.Prepare(`SELECT "lead_id" FROM "lead" WHERE "phone_number" = $1`)
	if err != nil {
		return leadId, fmt.Errorf("error preparing statement: %w", err)
	}
	defer stmt.Close()

	row := stmt.QueryRow(from)
	err = row.Scan(&leadId)
	if err != nil {
		return leadId, fmt.Errorf("error scanning row: %w", err)
	}

	return leadId, nil
}

func GetLeadIDFromIncomingTextMessage(from string) (int, error) {
	var leadId int

	stmt, err := DB.Prepare(`SELECT l.lead_id FROM "lead" AS l WHERE l.phone_number = $1`)
	if err != nil {
		return leadId, err
	}
	defer stmt.Close()

	row := stmt.QueryRow(from)

	var forwardingLeadID sql.NullInt64
	err = row.Scan(&forwardingLeadID)
	if err != nil && err != sql.ErrNoRows {
		return leadId, err
	}

	if forwardingLeadID.Valid {
		leadId = int(forwardingLeadID.Int64)
	}

	return leadId, nil
}

func UpdateLead(form types.UpdateLeadForm) error {
	if form.LeadID == nil {
		return fmt.Errorf("lead_id cannot be nil")
	}

	query := `
		UPDATE lead
		SET first_name = COALESCE($2, first_name), 
		    last_name = COALESCE($3, last_name), 
		    phone_number = COALESCE($4, phone_number), 
		    event_type_id = COALESCE($5, event_type_id), 
		    venue_type_id = COALESCE($6, venue_type_id), 
		    guests = COALESCE($7, guests)
		WHERE lead_id = $1
	`

	args := []interface{}{
		*form.LeadID,
		utils.CreateNullString(form.FirstName),
		utils.CreateNullString(form.LastName),
		utils.CreateNullString(form.PhoneNumber),
		utils.CreateNullInt(form.EventType),
		utils.CreateNullInt(form.VenueType),
		utils.CreateNullInt(form.Guests),
	}

	_, err := DB.Exec(query, args...)
	if err != nil {
		return fmt.Errorf("failed to update lead: %v", err)
	}

	return nil
}

func UpdateLeadMarketing(form types.UpdateLeadMarketingForm) error {
	if form.LeadID == nil {
		return fmt.Errorf("lead_id cannot be nil")
	}

	query := `
		UPDATE lead_marketing
		SET ad_campaign = $2, 
		    medium = $3, 
		    source = $4, 
		    referrer = $5, 
		    landing_page = $6,
		    ip = $7, 
		    keyword = $8, 
		    channel = $9, 
		    language = $10
		WHERE lead_id = $1
	`

	args := []interface{}{
		*form.LeadID,
		utils.CreateNullString(form.CampaignName),
		utils.CreateNullString(form.Medium),
		utils.CreateNullString(form.Source),
		utils.CreateNullString(form.Referrer),
		utils.CreateNullString(form.LandingPage),
		utils.CreateNullString(form.IP),
		utils.CreateNullString(form.Keyword),
		utils.CreateNullString(form.Channel),
		utils.CreateNullString(form.Language),
	}

	_, err := DB.Exec(query, args...)
	if err != nil {
		return fmt.Errorf("failed to update lead marketing: %v", err)
	}

	return nil
}

func GetForwardPhoneNumber(to, from string) (types.IncomingPhoneCallForwarding, error) {
	var forwardingCall types.IncomingPhoneCallForwarding

	stmt, err := DB.Prepare(`SELECT u.first_name, u.user_id FROM "user" AS u WHERE u.phone_number = $1`)
	if err != nil {
		return forwardingCall, err
	}
	defer stmt.Close()

	row := stmt.QueryRow(to)

	err = row.Scan(&forwardingCall.FirstName, &forwardingCall.UserID)
	if err != nil {
		return forwardingCall, err
	}

	stmt, err = DB.Prepare(`SELECT l.lead_id FROM "lead" AS l WHERE l.phone_number = $1`)
	if err != nil {
		return forwardingCall, err
	}
	defer stmt.Close()

	row = stmt.QueryRow(from)

	var leadID sql.NullInt64
	err = row.Scan(&leadID)
	if err != nil && err != sql.ErrNoRows {
		return forwardingCall, err
	}

	if leadID.Valid {
		forwardingCall.LeadID = int(leadID.Int64)
	} else {
		forwardingCall.LeadID = 0
	}

	switch forwardingCall.FirstName {
	case "Yovana":
		forwardingCall.ForwardPhoneNumber = "+1" + constants.YovaPhoneNumber
	case "David":
		forwardingCall.ForwardPhoneNumber = "+1" + constants.DavidPhoneNumber
	default:
		return forwardingCall, errors.New("no matching phone number")
	}

	return forwardingCall, nil
}

func GetPhoneCallBySID(sid string) (models.PhoneCall, error) {
	var phoneCall models.PhoneCall

	stmt, err := DB.Prepare(`SELECT phone_call_id, external_id, user_id, lead_id, call_duration, date_created, call_from, call_to, is_inbound, recording_url, status FROM phone_call WHERE external_id = $1`)
	if err != nil {
		return phoneCall, err
	}
	defer stmt.Close()

	var leadID, callDuration sql.NullInt64
	var recordingUrl sql.NullString

	row := stmt.QueryRow(sid)

	err = row.Scan(
		&phoneCall.PhoneCallID,
		&phoneCall.ExternalID,
		&phoneCall.UserID,
		&leadID,
		&callDuration,
		&phoneCall.DateCreated,
		&phoneCall.CallFrom,
		&phoneCall.CallTo,
		&phoneCall.IsInbound,
		&recordingUrl,
		&phoneCall.Status,
	)
	if err != nil {
		return phoneCall, err
	}

	if leadID.Valid {
		phoneCall.LeadID = int(leadID.Int64)
	}

	if callDuration.Valid {
		phoneCall.CallDuration = int(callDuration.Int64)
	}

	if recordingUrl.Valid {
		phoneCall.RecordingURL = recordingUrl.String
	}

	return phoneCall, nil
}

func UpdatePhoneCall(phoneCall models.PhoneCall) error {
	query := `
		UPDATE phone_call SET
			user_id = $1,
			lead_id = $2,
			call_duration = COALESCE($3, call_duration),
			date_created = $4,
			call_from = $5,
			call_to = $6,
			is_inbound = $7,
			recording_url = COALESCE($8, recording_url),
			status = COALESCE($9, status)
		WHERE external_id = $10`

	args := []interface{}{
		phoneCall.UserID,
		phoneCall.LeadID,
		utils.CreateNullInt(&phoneCall.CallDuration),
		phoneCall.DateCreated,
		phoneCall.CallFrom,
		phoneCall.CallTo,
		phoneCall.IsInbound,
		utils.CreateNullString(&phoneCall.RecordingURL),
		utils.CreateNullString(&phoneCall.Status),
		phoneCall.ExternalID,
	}

	_, err := DB.Exec(query, args...)
	if err != nil {
		return fmt.Errorf("error updating phone call: %w", err)
	}

	return nil
}

func GetSession(userKey string) (models.Session, error) {
	var session models.Session
	sqlStatement := `
        SELECT session_id, user_id, csrf_secret, external_id, date_created, date_expires
        FROM sessions
        WHERE csrf_secret = $1
    `
	row := DB.QueryRow(sqlStatement, userKey)

	var dateCreated, dateExpires time.Time
	var userID sql.NullInt32

	err := row.Scan(
		&session.SessionID,
		&userID,
		&session.CSRFSecret,
		&session.ExternalID,
		&dateCreated,
		&dateExpires,
	)
	if err != nil {
		return session, err
	}

	if userID.Valid {
		session.UserID = int(userID.Int32)
	}

	session.DateCreated = dateCreated.Unix()
	session.DateExpires = dateExpires.Unix()

	return session, nil
}

func CreateSession(session models.Session) error {
	sqlStatement := `
        INSERT INTO sessions (csrf_secret, external_id, date_created, date_expires)
        VALUES ($1, $2, to_timestamp($3)::timestamptz AT TIME ZONE 'America/New_York', to_timestamp($4)::timestamptz AT TIME ZONE 'America/New_York')
    `

	_, err := DB.Exec(sqlStatement,
		session.CSRFSecret,
		session.ExternalID,
		session.DateCreated,
		session.DateExpires,
	)

	if err != nil {
		return err
	}

	return nil
}

func UpdateSession(session models.Session) error {
	sqlStatement := `
        UPDATE sessions
        SET external_id = COALESCE($1, external_id),
            user_id = COALESCE($2, user_id)
        WHERE csrf_secret = $3
    `

	args := []interface{}{
		utils.CreateNullString(&session.ExternalID),
		utils.CreateNullInt(&session.UserID),
		session.CSRFSecret,
	}

	_, err := DB.Exec(sqlStatement, args...)
	if err != nil {
		return fmt.Errorf("error updating session: %w", err)
	}

	return nil
}

func DeleteSession(secret string) error {
	sqlStatement := `
        DELETE FROM sessions WHERE csrf_secret = $1
    `
	_, err := DB.Exec(sqlStatement, secret)
	if err != nil {
		return err
	}

	return nil
}

func DeleteLead(id int) error {
	sqlStatement := `
        DELETE FROM lead WHERE lead_id = $1
    `
	_, err := DB.Exec(sqlStatement, id)
	if err != nil {
		return err
	}

	return nil
}

func CreateEstimate(form types.EstimateForm, price float64) error {
	query := `
		INSERT INTO estimate (guests, hours, package_type_id, alcohol_segment_id, 
                             will_provide_liquor, will_provide_beer_and_wine, 
                             will_provide_mixers, will_provide_juices, 
                             will_provide_soft_drinks, will_provide_cups, 
                             will_provide_ice, will_require_glassware, 
                             will_require_bar, num_bars, price, date_created, lead_id)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, to_timestamp($16)::timestamptz, $17)
	`

	_, err := DB.Exec(
		query,
		utils.CreateNullInt(form.Guests),
		utils.CreateNullInt(form.Hours),
		utils.CreateNullInt(form.PackageTypeID),
		utils.CreateNullInt(form.AlcoholSegmentID),
		utils.CreateNullBool(form.WillProvideLiquor),
		utils.CreateNullBool(form.WillProvideBeerAndWine),
		utils.CreateNullBool(form.WillProvideMixers),
		utils.CreateNullBool(form.WillProvideJuices),
		utils.CreateNullBool(form.WillProvideSoftDrinks),
		utils.CreateNullBool(form.WillProvideCups),
		utils.CreateNullBool(form.WillProvideIce),
		utils.CreateNullBool(form.WillRequireGlassware),
		utils.CreateNullBool(form.WillRequireBar),
		utils.CreateNullInt(form.NumBars),
		price,
		time.Now().Unix(),
		utils.CreateNullInt(form.LeadID),
	)
	if err != nil {
		return fmt.Errorf("error inserting estimate data: %w", err)
	}

	return nil
}

func UpdateEstimate(form types.EstimateForm, price float64) error {
	query := `
		UPDATE estimate
		SET 
		    guests = COALESCE($2, guests),
		    hours = COALESCE($3, hours),
		    package_type_id = COALESCE($4, package_type_id),
		    alcohol_segment_id = COALESCE($5, alcohol_segment_id),
		    will_provide_liquor = COALESCE($6, will_provide_liquor),
		    will_provide_beer_and_wine = COALESCE($7, will_provide_beer_and_wine),
		    will_provide_mixers = COALESCE($8, will_provide_mixers),
		    will_provide_juices = COALESCE($9, will_provide_juices),
		    will_provide_soft_drinks = COALESCE($10, will_provide_soft_drinks),
		    will_provide_cups = COALESCE($11, will_provide_cups),
		    will_provide_ice = COALESCE($12, will_provide_ice),
		    will_require_glassware = COALESCE($13, will_require_glassware),
		    will_require_bar = COALESCE($14, will_require_bar),
		    num_bars = COALESCE($15, num_bars),
		    price = COALESCE($16, price),
		    date_updated = COALESCE(to_timestamp($17)::timestamptz, date_updated),
		    lead_id = COALESCE($18, lead_id)
		WHERE estimate_id = $1
	`
	_, err := DB.Exec(
		query,
		utils.CreateNullInt(form.EstimateID),
		utils.CreateNullInt(form.Guests),
		utils.CreateNullInt(form.Hours),
		utils.CreateNullInt(form.PackageTypeID),
		utils.CreateNullInt(form.AlcoholSegmentID),
		utils.CreateNullBool(form.WillProvideLiquor),
		utils.CreateNullBool(form.WillProvideBeerAndWine),
		utils.CreateNullBool(form.WillProvideMixers),
		utils.CreateNullBool(form.WillProvideJuices),
		utils.CreateNullBool(form.WillProvideSoftDrinks),
		utils.CreateNullBool(form.WillProvideCups),
		utils.CreateNullBool(form.WillProvideIce),
		utils.CreateNullBool(form.WillRequireGlassware),
		utils.CreateNullBool(form.WillRequireBar),
		utils.CreateNullInt(form.NumBars),
		price,
		time.Now().Unix(),
		utils.CreateNullInt(form.LeadID),
	)
	if err != nil {
		return fmt.Errorf("error updating estimate data: %w", err)
	}

	return nil
}

func AssignStripeCustomerToLead(stripeCustomerID string, leadId int) error {
	stmt, err := DB.Prepare(`UPDATE lead SET stripe_customer_id = $1 WHERE lead_id = $2`)
	if err != nil {
		return fmt.Errorf("error preparing statement: %w", err)
	}
	defer stmt.Close()

	_, err = stmt.Exec(stripeCustomerID, leadId)
	if err != nil {
		return fmt.Errorf("error executing statement: %w", err)
	}

	return nil
}

func AssignStripeInvoiceToEstimate(stripeInvoiceId string, estimateId int) error {
	stmt, err := DB.Prepare(`UPDATE estimate SET stripe_invoice_id = $1 WHERE estimate_id = $2`)
	if err != nil {
		return fmt.Errorf("error preparing statement: %w", err)
	}
	defer stmt.Close()

	_, err = stmt.Exec(stripeInvoiceId, estimateId)
	if err != nil {
		return fmt.Errorf("error executing statement: %w", err)
	}

	return nil
}

func GetAlcoholSegments() ([]models.AlcoholSegment, error) {
	var alcoholSegments []models.AlcoholSegment

	rows, err := DB.Query(`SELECT alcohol_segment_id, name, price_modification FROM "alcohol_segment"`)
	if err != nil {
		return alcoholSegments, fmt.Errorf("error executing query: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var segment models.AlcoholSegment
		err := rows.Scan(&segment.AlcoholSegmentID, &segment.Name, &segment.PriceModification)
		if err != nil {
			return alcoholSegments, fmt.Errorf("error scanning row: %w", err)
		}
		alcoholSegments = append(alcoholSegments, segment)
	}

	if err := rows.Err(); err != nil {
		return alcoholSegments, fmt.Errorf("error iterating rows: %w", err)
	}

	return alcoholSegments, nil
}

func GetPackageTypes() ([]models.PackageType, error) {
	var packageTypes []models.PackageType

	rows, err := DB.Query(`SELECT package_type_id, name, price_modification FROM "package_type"`)
	if err != nil {
		return packageTypes, fmt.Errorf("error executing query: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var pkg models.PackageType
		err := rows.Scan(&pkg.PackageTypeID, &pkg.Name, &pkg.PriceModification)
		if err != nil {
			return packageTypes, fmt.Errorf("error scanning row: %w", err)
		}
		packageTypes = append(packageTypes, pkg)
	}

	if err := rows.Err(); err != nil {
		return packageTypes, fmt.Errorf("error iterating rows: %w", err)
	}

	return packageTypes, nil
}

func GetBookingList(leadId int) ([]types.BookingList, error) {
	var bookings []types.BookingList

	rows, err := DB.Query(`
		SELECT 
			b.booking_id,
			b.lead_id,
			CONCAT(b.street_address, ', ', b.city, ', ', b.state, ', ', b.postal_code) AS address,
			b.start_time,
			b.end_time,
			CONCAT(u.first_name, ' ', u.last_name) AS bartender,
			e.price::NUMERIC
		FROM booking AS b
		JOIN estimate AS e ON e.estimate_id = b.estimate_id
		JOIN "user" AS u ON u.user_id = b.bartender_id
		WHERE b.lead_id = $1
		ORDER BY b.start_time ASC;
	`, leadId)
	if err != nil {
		return bookings, fmt.Errorf("error executing query: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var booking types.BookingList

		var startTime, endTime time.Time

		err := rows.Scan(
			&booking.BookingID,
			&booking.LeadID,
			&booking.Address,
			&booking.StartTime,
			&booking.EndTime,
			&booking.Bartender,
			&booking.Price,
		)
		if err != nil {
			return bookings, fmt.Errorf("error scanning row: %w", err)
		}

		booking.StartTime = utils.FormatTimestamp(startTime.Unix())
		booking.EndTime = utils.FormatTimestamp(endTime.Unix())

		bookings = append(bookings, booking)
	}

	if err := rows.Err(); err != nil {
		return bookings, fmt.Errorf("error iterating rows: %w", err)
	}

	return bookings, nil
}

func GetEstimateList(leadId int) ([]types.EstimatesList, error) {
	var estimates []types.EstimatesList

	rows, err := DB.Query(`
		SELECT 
			e.estimate_id,
			e.lead_id,
			e.date_created,
			e.price::NUMERIC,
			e.stripe_invoice_id,
			e.status
		FROM estimate AS e
		WHERE eb.lead_id = $1
		ORDER BY e.date_created ASC;
	`, leadId)
	if err != nil {
		return estimates, fmt.Errorf("error executing query: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var estimate types.EstimatesList

		var dateCreated time.Time

		err := rows.Scan(
			&estimate.EstimateID,
			&estimate.LeadID,
			&dateCreated,
			&estimate.Price,
			&estimate.StripeInvoiceID,
			&estimate.Status,
		)
		if err != nil {
			return estimates, fmt.Errorf("error scanning row: %w", err)
		}

		estimate.DateCreated = utils.FormatDateMMDDYYYY(dateCreated.Unix())

		estimates = append(estimates, estimate)
	}

	if err := rows.Err(); err != nil {
		return estimates, fmt.Errorf("error iterating rows: %w", err)
	}

	return estimates, nil
}

func CreateBooking(form types.BookingForm) error {
	stmt, err := DB.Prepare(`
		INSERT INTO booking (
			estimate_id,
			street_address,
			city,
			state,
			postal_code,
			country,
			start_time,
			end_time,
			bartender_id,
			lead_id
		) VALUES ($1, $2, $3, $4, $5, $6, to_timestamp($7)::timestamptz AT TIME ZONE 'America/New_York', to_timestamp($8)::timestamptz AT TIME ZONE 'America/New_York', $9, $10)
	`)
	if err != nil {
		return fmt.Errorf("error preparing statement: %w", err)
	}
	defer stmt.Close()

	_, err = stmt.Exec(
		utils.CreateNullInt(form.EstimateID),
		utils.CreateNullString(form.StreetAddress),
		utils.CreateNullString(form.City),
		utils.CreateNullString(form.State),
		utils.CreateNullString(form.PostalCode),
		utils.CreateNullString(form.Country),
		utils.CreateNullInt64(form.StartTime),
		utils.CreateNullInt64(form.EndTime),
		utils.CreateNullInt(form.BartenderID),
		utils.CreateNullInt(form.LeadID),
	)
	if err != nil {
		return fmt.Errorf("error executing statement: %w", err)
	}

	return nil
}

func UpdateBooking(form types.BookingForm) error {
	stmt, err := DB.Prepare(`
		UPDATE booking
		SET 
			estimate_id = COALESCE($2, estimate_id),
			street_address = COALESCE($3, street_address),
			city = COALESCE($4, city),
			state = COALESCE($5, state),
			postal_code = COALESCE($6, postal_code),
			country = COALESCE($7, country),
			start_time = COALESCE(to_timestamp($8) AT TIME ZONE 'America/New_York', start_time),
			end_time = COALESCE(to_timestamp($9) AT TIME ZONE 'America/New_York', end_time),
			bartender_id = COALESCE($10, bartender_id),
			lead_id = COALESCE($11, lead_id)
		WHERE booking_id = $1;
	`)
	if err != nil {
		return fmt.Errorf("error preparing statement: %w", err)
	}
	defer stmt.Close()

	_, err = stmt.Exec(
		utils.CreateNullInt(form.BookingID),
		utils.CreateNullInt(form.EstimateID),
		utils.CreateNullString(form.StreetAddress),
		utils.CreateNullString(form.City),
		utils.CreateNullString(form.State),
		utils.CreateNullString(form.PostalCode),
		utils.CreateNullString(form.Country),
		utils.CreateNullInt64(form.StartTime),
		utils.CreateNullInt64(form.EndTime),
		utils.CreateNullInt(form.BartenderID),
		utils.CreateNullInt(form.LeadID),
	)
	if err != nil {
		return fmt.Errorf("error executing statement: %w", err)
	}

	return nil
}
