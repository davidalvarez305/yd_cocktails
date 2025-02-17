package database

import (
	"database/sql"
	"fmt"
	"strconv"
	"time"

	"github.com/davidalvarez305/yd_cocktails/constants"
	"github.com/davidalvarez305/yd_cocktails/models"
	"github.com/davidalvarez305/yd_cocktails/types"
	"github.com/davidalvarez305/yd_cocktails/utils"
	"github.com/google/uuid"
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
		INSERT INTO lead (full_name, phone_number, created_at, message, opt_in_text_messaging, email, lead_status_id, next_action_id)
		VALUES ($1, $2, to_timestamp($3)::timestamptz AT TIME ZONE 'America/New_York', $4, $5, $6, $7, $8)
		RETURNING lead_id
	`)
	if err != nil {
		return leadID, fmt.Errorf("error preparing lead statement: %w", err)
	}
	defer leadStmt.Close()

	message := utils.CreateNullString(quoteForm.Message)

	err = leadStmt.QueryRow(
		utils.CreateNullString(quoteForm.FullName),
		utils.CreateNullString(quoteForm.PhoneNumber),
		utils.CreateNullInt64(quoteForm.CreatedAt),
		message,
		utils.CreateNullBool(quoteForm.OptInTextMessaging),
		utils.CreateNullString(quoteForm.Email),
		constants.NewLeadStatusID,
		constants.InitialContactActionID,
	).Scan(&leadID)
	if err != nil {
		return leadID, fmt.Errorf("error inserting lead: %w", err)
	}

	marketingStmt, err := tx.Prepare(`
		INSERT INTO lead_marketing (lead_id, source, medium, channel, landing_page, keyword, referrer, click_id, campaign_id, ad_campaign, ad_group_id, ad_group_name, ad_set_id, ad_set_name, ad_id, ad_headline, language, user_agent, button_clicked, ip, external_id, google_client_id, csrf_secret, facebook_click_id, facebook_client_id, longitude, latitude, instant_form_lead_id, instant_form_id, instant_form_name)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17, $18, $19, $20, $21, $22, $23, $24, $25, $26, $27, $28, $29, $30)
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
		utils.CreateNullInt64(quoteForm.InstantFormLeadID),
		utils.CreateNullInt64(quoteForm.InstantFormID),
		utils.CreateNullString(quoteForm.InstantFormName),
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
		INSERT INTO message (external_id, text, date_created, text_from, text_to, is_inbound, is_read)
		VALUES ($1, $2, to_timestamp($3)::timestamptz AT TIME ZONE 'America/New_York', $4, $5, $6, $7)
	`)
	if err != nil {
		return fmt.Errorf("error preparing statement: %w", err)
	}
	defer stmt.Close()

	_, err = stmt.Exec(msg.ExternalID, msg.Text, msg.DateCreated, msg.TextFrom, msg.TextTo, msg.IsInbound, msg.IsRead)
	if err != nil {
		return fmt.Errorf("error executing statement: %w", err)
	}

	return nil
}

func SetSMSToRead(messageId int) error {
	stmt, err := DB.Prepare(`
		UPDATE message
		SET is_read = true
		WHERE message_id = $1
	`)
	if err != nil {
		return fmt.Errorf("error preparing statement: %w", err)
	}
	defer stmt.Close()

	_, err = stmt.Exec(messageId)
	if err != nil {
		return fmt.Errorf("error executing statement: %w", err)
	}

	return nil
}

func SavePhoneCall(phoneCall models.PhoneCall) error {
	stmt, err := DB.Prepare(`
		INSERT INTO phone_call (external_id, call_duration, date_created, call_from, call_to, is_inbound, recording_url, status) VALUES ($1, $2, to_timestamp($3)::timestamptz AT TIME ZONE 'America/New_York', $4, $5, $6, $7, $8)`)
	if err != nil {
		return fmt.Errorf("error preparing statement: %w", err)
	}
	defer stmt.Close()

	var callDuration sql.NullInt64

	if phoneCall.CallDuration != 0 {
		callDuration = sql.NullInt64{Int64: int64(phoneCall.CallDuration), Valid: true}
	}

	_, err = stmt.Exec(
		phoneCall.ExternalID,
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

	query := `SELECT l.lead_id, l.full_name, l.phone_number, 
			l.created_at, lm.language, li.interest, ls.status, na.action, 
			COUNT(*) OVER() AS total_rows
		FROM lead AS l
		JOIN lead_marketing AS lm ON lm.lead_id = l.lead_id
		LEFT JOIN lead_interest AS li ON li.lead_interest_id = l.lead_interest_id
		LEFT JOIN lead_status AS ls ON ls.lead_status_id = l.lead_status_id
		LEFT JOIN next_action AS na ON na.next_action_id = l.next_action_id
		WHERE 
			(
				$5::TEXT IS NOT NULL 
				AND (
					l.search_vector @@ plainto_tsquery('english', $5::TEXT)
					OR l.full_name ILIKE '%' || $5 || '%'
					OR l.phone_number ILIKE '%' || $5 || '%'
				)
			)
			OR 
			(
				$5 IS NULL 
				AND (
					($6::INTEGER IS NOT NULL AND ls.lead_status_id = $6::INTEGER) 
					OR 
					($6::INTEGER IS NULL AND (ls.lead_status_id IS DISTINCT FROM $3::INTEGER OR ls.lead_status_id IS NULL))
				)
				AND 
				(
					($7::INTEGER IS NOT NULL AND li.lead_interest_id = $7::INTEGER) 
					OR 
					($7::INTEGER IS NULL AND (li.lead_interest_id IS DISTINCT FROM $4::INTEGER OR li.lead_interest_id IS NULL))
				)
				AND 
				($8::INTEGER IS NULL OR na.next_action_id = $8::INTEGER)
			)
		ORDER BY l.created_at DESC
		LIMIT $1 OFFSET $2;`

	var offset int

	// Handle pagination
	if params.PageNum != nil {
		pageNum, err := strconv.Atoi(*params.PageNum)
		if err != nil {
			return nil, 0, fmt.Errorf("could not convert page num: %w", err)
		}
		offset = (pageNum - 1) * int(constants.LeadsPerPage)
	}

	rows, err := DB.Query(query, constants.LeadsPerPage, offset, constants.ArchivedLeadStatusID, constants.NoInterestLeadInterestID,
		utils.CreateNullString(params.Search),
		utils.CreateNullInt(params.LeadStatusID),
		utils.CreateNullInt(params.LeadInterestID),
		utils.CreateNullInt(params.NextActionID))
	if err != nil {
		return nil, 0, fmt.Errorf("error executing query: %w", err)
	}
	defer rows.Close()

	var totalRows int
	for rows.Next() {
		var lead types.LeadList
		var createdAt time.Time
		var language, nextAction, leadInterest, leadStatus sql.NullString

		err := rows.Scan(&lead.LeadID,
			&lead.FullName,
			&lead.PhoneNumber,
			&createdAt,
			&language,
			&leadInterest,
			&leadStatus,
			&nextAction,
			&totalRows)
		if err != nil {
			return nil, 0, fmt.Errorf("error scanning row: %w", err)
		}
		lead.CreatedAt = utils.FormatTimestampEST(createdAt.Unix())

		if language.Valid {
			lead.Language = language.String
		}

		if nextAction.Valid {
			lead.NextAction = nextAction.String
		}

		if leadInterest.Valid {
			lead.LeadInterest = leadInterest.String
		}

		if leadStatus.Valid {
			lead.LeadStatus = leadStatus.String
		}

		leads = append(leads, lead)
	}

	if err := rows.Err(); err != nil {
		return nil, 0, fmt.Errorf("error iterating rows: %w", err)
	}

	return leads, totalRows, nil
}

func GetReferrals() ([]types.Referral, error) {
	var referrals []types.Referral

	query := `SELECT l.lead_id, l.full_name
		FROM lead AS l
		ORDER BY l.created_at ASC`

	rows, err := DB.Query(query)
	if err != nil {
		return nil, fmt.Errorf("error executing query: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var referral types.Referral

		err := rows.Scan(&referral.LeadID, &referral.FullName)
		if err != nil {
			return nil, fmt.Errorf("error scanning row: %w", err)
		}

		referrals = append(referrals, referral)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating rows: %w", err)
	}

	return referrals, nil
}

func GetLeadDetails(leadID string) (types.LeadDetails, error) {
	query := `SELECT l.lead_id,
	l.full_name,
	l.phone_number,
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
	lm.instant_form_lead_id,
	lm.instant_form_id,
	lm.instant_form_name,
	lm.referral_lead_id,
	li.lead_interest_id,
	ls.lead_status_id,
	na.next_action_id
	FROM lead l
	JOIN lead_marketing lm ON l.lead_id = lm.lead_id
	LEFT JOIN lead_interest li ON l.lead_interest_id = li.lead_interest_id
	LEFT JOIN lead_status ls ON l.lead_status_id = ls.lead_status_id
	LEFT JOIN next_action na ON l.next_action_id = na.next_action_id
	WHERE l.lead_id = $1`

	var leadDetails types.LeadDetails

	row := DB.QueryRow(query, leadID)

	var adCampaign, medium, source, referrer, landingPage, ip, keyword, channel, language, email, facebookClickId, facebookClientId sql.NullString
	var message, externalId, userAgent, clickId, googleClientId sql.NullString
	var campaignId, instantFormleadId, instantFormId, referralLeadId, leadInterestId, leadStatusId, nextActionId sql.NullInt64

	var buttonClicked, instantFormName sql.NullString

	err := row.Scan(
		&leadDetails.LeadID,
		&leadDetails.FullName,
		&leadDetails.PhoneNumber,
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
		&instantFormleadId,
		&instantFormId,
		&instantFormName,
		&referralLeadId,
		&leadInterestId,
		&leadStatusId,
		&nextActionId,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return leadDetails, fmt.Errorf("no lead found with ID %s", leadID)
		}
		return leadDetails, fmt.Errorf("error scanning row: %w", err)
	}

	// Map the nullable fields to your struct
	if leadInterestId.Valid {
		leadDetails.LeadInterestID = int(leadInterestId.Int64)
	}
	if leadStatusId.Valid {
		leadDetails.LeadStatusID = int(leadStatusId.Int64)
	}
	if nextActionId.Valid {
		leadDetails.NextActionID = int(nextActionId.Int64)
	}
	if referralLeadId.Valid {
		leadDetails.ReferralLeadID = int(referralLeadId.Int64)
	}
	if buttonClicked.Valid {
		leadDetails.ButtonClicked = buttonClicked.String
	}
	if instantFormleadId.Valid {
		leadDetails.InstantFormLeadID = instantFormleadId.Int64
	}
	if instantFormId.Valid {
		leadDetails.InstantFormID = instantFormId.Int64
	}
	if instantFormName.Valid {
		leadDetails.InstantFormName = instantFormName.String
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

func GetConversionReporting(leadID int) (types.ConversionReporting, error) {
	query := `WITH referral_lead AS (
		SELECT 
			l.lead_id, 
			l.phone_number, 
			l.email, 
			l.first_name, 
			l.last_name, 
			l.city, 
			l.state, 
			l.zip_code
		FROM lead l
		JOIN lead_marketing lm ON l.lead_id = lm.lead_id
		WHERE lm.referral_lead_id = $1
		FROM lead_marketing
		WHERE lead_id = $1
	),
	lead_events AS (
		SELECT SUM(e.amount::NUMERIC + e.tip::NUMERIC) AS total_revenue
		FROM event e
		WHERE e.lead_id = $1
	),
	referral_events AS (
		SELECT SUM(e.amount::NUMERIC + e.tip::NUMERIC) AS total_revenue
		FROM event e
		WHERE e.lead_id IN (
			SELECT lm1.lead_id
			FROM lead_marketing lm1
			WHERE lm1.referral_lead_id = (SELECT referral_lead_id FROM referral_lead)
		)
	)
	SELECT 
		COALESCE(referral_lead.lead_id, l.lead_id) AS lead_id,
		COALESCE(referral_lead.phone_number, l.phone_number) AS phone_number,
		COALESCE(referral_lead.email, l.email) AS email,
		COALESCE(referral_lead.first_name, l.first_name) AS first_name,
		COALESCE(referral_lead.last_name, l.last_name) AS last_name,
		COALESCE(referral_lead.city, l.city) AS city,
		COALESCE(referral_lead.state, l.state) AS state,
		COALESCE(referral_lead.zip_code, l.zip_code) AS zip_code,
		COALESCE(referral_events.total_revenue, 0) + COALESCE(lead_events.total_revenue, 0) AS total_revenue
	FROM lead l
	LEFT JOIN referral_lead ON TRUE
	LEFT JOIN lead_events ON TRUE
	LEFT JOIN referral_events ON TRUE
	WHERE l.lead_id = $1;`

	var conversionReporting types.ConversionReporting

	row := DB.QueryRow(query, leadID)

	var adCampaign, landingPage, ip, email, facebookClickId, facebookClientId sql.NullString
	var externalId, userAgent, clickId, googleClientId sql.NullString
	var campaignId, instantFormleadId, eventId sql.NullInt64
	var revenue sql.NullFloat64

	err := row.Scan(
		&conversionReporting.LeadID,
		&conversionReporting.PhoneNumber,
		&adCampaign,
		&landingPage,
		&ip,
		&email,
		&facebookClickId,
		&facebookClientId,
		&externalId,
		&userAgent,
		&clickId,
		&googleClientId,
		&campaignId,
		&instantFormleadId,
		&eventId,
		&revenue,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return conversionReporting, fmt.Errorf("no lead found with ID %v", leadID)
		}
		return conversionReporting, fmt.Errorf("error scanning row: %w", err)
	}

	if revenue.Valid {
		conversionReporting.Revenue = revenue.Float64
	}

	if instantFormleadId.Valid {
		conversionReporting.InstantFormLeadID = instantFormleadId.Int64
	}

	if clickId.Valid {
		conversionReporting.ClickID = clickId.String
	}

	if googleClientId.Valid {
		conversionReporting.GoogleClientID = googleClientId.String
	}

	if externalId.Valid {
		conversionReporting.ExternalID = externalId.String
	}

	if userAgent.Valid {
		conversionReporting.UserAgent = userAgent.String
	}

	if facebookClickId.Valid {
		conversionReporting.FacebookClickID = facebookClickId.String
	}

	if facebookClientId.Valid {
		conversionReporting.FacebookClientID = facebookClientId.String
	}

	if email.Valid {
		conversionReporting.Email = email.String
	}

	if adCampaign.Valid {
		conversionReporting.CampaignName = adCampaign.String
	}
	if campaignId.Valid {
		conversionReporting.CampaignID = campaignId.Int64
	}
	if eventId.Valid {
		conversionReporting.EventID = int(eventId.Int64)
	}

	if landingPage.Valid {
		conversionReporting.LandingPage = landingPage.String
	}

	if ip.Valid {
		conversionReporting.IP = ip.String
	}

	return conversionReporting, nil
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

func UpdateLead(form types.UpdateLeadForm) error {
	if form.LeadID == nil {
		return fmt.Errorf("lead_id cannot be nil")
	}

	query := `
		UPDATE lead
		SET full_name = COALESCE($2, full_name), 
		    phone_number = COALESCE($3, phone_number), 
		    email = $4,
			stripe_customer_id = $5,
			lead_interest_id = $6,
			lead_status_id = $7,
			next_action_id = $8
		WHERE lead_id = $1
	`

	args := []interface{}{
		*form.LeadID,
		utils.CreateNullString(form.FullName),
		utils.CreateNullString(form.PhoneNumber),
		utils.CreateNullString(form.Email),
		utils.CreateNullString(form.StripeCustomerID),
		utils.CreateNullInt(form.LeadInterestID),
		utils.CreateNullInt(form.LeadStatusID),
		utils.CreateNullInt(form.NextActionID),
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
		    language = $10,
			referral_lead_id = $11
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
		utils.CreateNullInt(form.ReferralLeadID),
	}

	_, err := DB.Exec(query, args...)
	if err != nil {
		return fmt.Errorf("failed to update lead marketing: %v", err)
	}

	return nil
}

func GetForwardPhoneNumber(to, from string) (string, error) {
	var forwardPhoneNumber string

	query := `SELECT u.forward_phone_number FROM "user" AS u WHERE u.phone_number IN ($1, $2) LIMIT 1`
	row := DB.QueryRow(query, to, from)

	err := row.Scan(&forwardPhoneNumber)
	if err != nil {
		return forwardPhoneNumber, err
	}

	return "1" + forwardPhoneNumber, nil
}

func GetPhoneCallBySID(sid string) (models.PhoneCall, error) {
	var phoneCall models.PhoneCall

	stmt, err := DB.Prepare(`SELECT phone_call_id, external_id, call_duration, date_created, call_from, call_to, is_inbound, recording_url, status FROM phone_call WHERE external_id = $1`)
	if err != nil {
		return phoneCall, err
	}
	defer stmt.Close()

	var callDuration sql.NullInt64
	var recordingUrl, status sql.NullString

	row := stmt.QueryRow(sid)

	var dateCreated time.Time

	err = row.Scan(
		&phoneCall.PhoneCallID,
		&phoneCall.ExternalID,
		&callDuration,
		&dateCreated,
		&phoneCall.CallFrom,
		&phoneCall.CallTo,
		&phoneCall.IsInbound,
		&recordingUrl,
		&status,
	)
	if err != nil {
		return phoneCall, err
	}

	if callDuration.Valid {
		phoneCall.CallDuration = int(callDuration.Int64)
	}

	if recordingUrl.Valid {
		phoneCall.RecordingURL = recordingUrl.String
	}

	if status.Valid {
		phoneCall.Status = status.String
	}

	phoneCall.DateCreated = dateCreated.Unix()

	return phoneCall, nil
}

func UpdatePhoneCall(phoneCall models.PhoneCall) error {
	query := `
		UPDATE phone_call SET
			call_duration = $1,
			recording_url = $2,
			status = $3
		WHERE external_id = $4`

	args := []interface{}{
		utils.CreateNullInt(&phoneCall.CallDuration),
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
	var csrfSecret, externalId sql.NullString

	err := row.Scan(
		&session.SessionID,
		&userID,
		&csrfSecret,
		&externalId,
		&dateCreated,
		&dateExpires,
	)
	if err != nil {
		return session, err
	}

	if userID.Valid {
		session.UserID = int(userID.Int32)
	}

	if csrfSecret.Valid {
		session.CSRFSecret = csrfSecret.String
	}

	if externalId.Valid {
		session.ExternalID = externalId.String
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
        SET external_id = $1,
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

func ArchiveLead(id int) error {
	sqlStatement := `
        UPDATE lead
		SET lead_status_id = $2
		WHERE lead_id = $1
    `
	_, err := DB.Exec(sqlStatement, id, constants.ArchivedLeadStatusID)
	if err != nil {
		return err
	}

	return nil
}

func CreateEvent(form types.EventForm) error {
	query := `
		INSERT INTO event (
			bartender_id, lead_id, event_type_id, venue_type_id, street_address, city, zip_code,
			start_time, end_time, date_created, date_paid, amount, tip, guests
		)
		VALUES (
			$1, $2, $3, $4, $5, $6, $7,
			to_timestamp($8)::timestamptz AT TIME ZONE 'America/New_York',
			to_timestamp($9)::timestamptz AT TIME ZONE 'America/New_York',
			to_timestamp($10)::timestamptz AT TIME ZONE 'America/New_York',
			to_timestamp($11)::timestamptz AT TIME ZONE 'America/New_York',
			$12, $13, $14
		)
	`

	_, err := DB.Exec(
		query,
		utils.CreateNullInt(form.BartenderID),
		utils.CreateNullInt(form.LeadID),
		utils.CreateNullInt(form.EventTypeID),
		utils.CreateNullInt(form.VenueTypeID),
		utils.CreateNullString(form.StreetAddress),
		utils.CreateNullString(form.City),
		utils.CreateNullString(form.ZipCode),
		utils.CreateNullInt64(form.StartTime),
		utils.CreateNullInt64(form.EndTime),
		utils.CreateNullInt64(form.DateCreated),
		utils.CreateNullInt64(form.DatePaid),
		utils.CreateNullFloat64(form.Amount),
		utils.CreateNullFloat64(form.Tip),
		utils.CreateNullInt(form.Guests),
	)
	if err != nil {
		return fmt.Errorf("error inserting event data: %w", err)
	}

	return nil
}

func UpdateEvent(form types.EventForm) error {
	query := `
		UPDATE event
		SET 
			bartender_id = COALESCE($2, bartender_id),
			event_type_id = COALESCE($3, event_type_id),
			venue_type_id = COALESCE($4, venue_type_id),
			street_address = COALESCE($5, street_address),
			city = COALESCE($6, city),
			zip_code = COALESCE($7, zip_code),
			start_time = COALESCE(to_timestamp($8)::timestamptz AT TIME ZONE 'America/New_York', start_time),
			end_time = COALESCE(to_timestamp($9)::timestamptz AT TIME ZONE 'America/New_York', end_time),
			date_paid = to_timestamp($10)::timestamptz AT TIME ZONE 'America/New_York',
			amount = $11,
			tip = $12,
			guests = $13
		WHERE event_id = $1
	`

	_, err := DB.Exec(
		query,
		utils.CreateNullInt(form.EventID),
		utils.CreateNullInt(form.BartenderID),
		utils.CreateNullInt(form.EventTypeID),
		utils.CreateNullInt(form.VenueTypeID),
		utils.CreateNullString(form.StreetAddress),
		utils.CreateNullString(form.City),
		utils.CreateNullString(form.ZipCode),
		utils.CreateNullInt64(form.StartTime),
		utils.CreateNullInt64(form.EndTime),
		utils.CreateNullInt64(form.DatePaid),
		utils.CreateNullFloat64(form.Amount),
		utils.CreateNullFloat64(form.Tip),
		utils.CreateNullInt(form.Guests),
	)
	if err != nil {
		return fmt.Errorf("error updating event data: %w", err)
	}

	return nil
}

func DeleteEvent(id int) error {
	sqlStatement := `
        DELETE FROM event WHERE event_id = $1
    `
	_, err := DB.Exec(sqlStatement, id)
	if err != nil {
		return err
	}

	return nil
}

func GetEventList(leadId int) ([]types.EventList, error) {
	var events []types.EventList

	rows, err := DB.Query(`
		SELECT 
			e.event_id,
			e.lead_id,
			e.amount::NUMERIC + e.tip::NUMERIC,
			l.full_name,
			CONCAT(b.first_name, ' ', b.last_name),
			et.name,
			vt.name,
			e.guests,
			e.start_time,
			e.end_time
		FROM event AS e
		JOIN lead AS l ON l.lead_id = e.lead_id
		LEFT JOIN "user" AS b ON b.user_id = e.bartender_id
		JOIN event_type AS et ON et.event_type_id = e.event_type_id
		JOIN venue_type AS vt ON vt.venue_type_id = e.venue_type_id
		WHERE e.lead_id = $1
		ORDER BY e.date_created ASC;
	`, leadId)
	if err != nil {
		return events, fmt.Errorf("error executing query: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var event types.EventList
		var eventStart, eventEnd time.Time

		err := rows.Scan(
			&event.EventID,
			&event.LeadID,
			&event.Amount,
			&event.LeadName,
			&event.Bartender,
			&event.EventType,
			&event.VenueType,
			&event.Guests,
			&eventStart,
			&eventEnd,
		)
		if err != nil {
			return events, fmt.Errorf("error scanning row: %w", err)
		}

		event.EventTime = fmt.Sprintf(
			"%s - %s",
			utils.FormatTimestampEST(eventStart.Unix()),
			utils.FormatTimestampEST(eventEnd.Unix()),
		)

		events = append(events, event)
	}

	if err := rows.Err(); err != nil {
		return events, fmt.Errorf("error iterating rows: %w", err)
	}

	return events, nil
}

func GetEventDetails(eventId string) (models.Event, error) {
	query := `SELECT 
		event_id, 
		bartender_id, 
		lead_id,
		street_address,
		city,
		zip_code,
		start_time,
		end_time,
		date_created,
		date_paid,
		amount::NUMERIC,
		tip::NUMERIC,
		event_type_id,
		venue_type_id,
		guests
	FROM event 
	WHERE event_id = $1`

	var eventDetails models.Event

	// Declare nullable SQL variables for fields that might be NULL in the database
	var streetAddress, city, zipCode sql.NullString
	var startTime, endTime, dateCreated, datePaid sql.NullTime
	var amount, tip sql.NullFloat64
	var bartenderID, eventTypeID, venueTypeID, guests sql.NullInt64

	row := DB.QueryRow(query, eventId)

	err := row.Scan(
		&eventDetails.EventID,
		&bartenderID,
		&eventDetails.LeadID,
		&streetAddress,
		&city,
		&zipCode,
		&startTime,
		&endTime,
		&dateCreated,
		&datePaid,
		&amount,
		&tip,
		&eventTypeID,
		&venueTypeID,
		&guests,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return eventDetails, fmt.Errorf("no event found with ID %s", eventId)
		}
		return eventDetails, fmt.Errorf("error scanning row: %w", err)
	}

	// Map nullable SQL variables to the Event struct
	if streetAddress.Valid {
		eventDetails.StreetAddress = streetAddress.String
	}
	if city.Valid {
		eventDetails.City = city.String
	}
	if zipCode.Valid {
		eventDetails.ZipCode = zipCode.String
	}
	if startTime.Valid {
		eventDetails.StartTime = startTime.Time.Unix()
	}
	if endTime.Valid {
		eventDetails.EndTime = endTime.Time.Unix()
	}
	if dateCreated.Valid {
		eventDetails.DateCreated = dateCreated.Time.Unix()
	}
	if datePaid.Valid {
		eventDetails.DatePaid = datePaid.Time.Unix()
	}
	if amount.Valid {
		eventDetails.Amount = amount.Float64
	}
	if tip.Valid {
		eventDetails.Tip = tip.Float64
	}
	if bartenderID.Valid {
		eventDetails.BartenderID = int(bartenderID.Int64)
	}
	if eventTypeID.Valid {
		eventDetails.EventTypeID = int(eventTypeID.Int64)
	}
	if venueTypeID.Valid {
		eventDetails.VenueTypeID = int(venueTypeID.Int64)
	}
	if guests.Valid {
		eventDetails.Guests = int(guests.Int64)
	}

	return eventDetails, nil
}

func GetUsers() ([]models.User, error) {
	var users []models.User

	stmt, err := DB.Prepare(`SELECT user_id, username, phone_number, password, user_role_id, first_name, last_name FROM "user"`)
	if err != nil {
		return nil, fmt.Errorf("error preparing statement: %w", err)
	}
	defer stmt.Close()

	rows, err := stmt.Query()
	if err != nil {
		return nil, fmt.Errorf("error executing query: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var user models.User
		err := rows.Scan(
			&user.UserID,
			&user.Username,
			&user.PhoneNumber,
			&user.Password,
			&user.UserRoleID,
			&user.FirstName,
			&user.LastName,
		)
		if err != nil {
			return nil, fmt.Errorf("error scanning row: %w", err)
		}

		users = append(users, user)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error with rows: %w", err)
	}

	return users, nil
}

func IsPhoneNumberInDB(phoneNumber string) (bool, error) {
	var exists bool
	err := DB.QueryRow("SELECT EXISTS(SELECT 1 FROM lead WHERE phone_number = $1)", phoneNumber).Scan(&exists)
	return exists, err
}

func CreateLeadQuote(form types.LeadQuoteForm) error {
	query := `
		INSERT INTO quote (
			lead_id, 
			number_of_bartenders, 
			guests, 
			hours, 
			event_type_id, 
			venue_type_id, 
			event_date, 
			external_id
		)
		VALUES (
			$1, $2, $3, $4, 
			$5, $6, to_timestamp($7)::timestamptz AT TIME ZONE 'America/New_York', 
			$8
		);
	`

	_, err := DB.Exec(
		query,
		utils.CreateNullInt(form.LeadID),
		utils.CreateNullInt(form.NumberOfBartenders),
		utils.CreateNullInt(form.Guests),
		utils.CreateNullInt(form.Hours),
		utils.CreateNullInt(form.EventTypeID),
		utils.CreateNullInt(form.VenueTypeID),
		utils.CreateNullInt64(form.EventDate),
		uuid.New().String(),
	)
	if err != nil {
		return fmt.Errorf("error inserting lead quote data: %w", err)
	}

	return nil
}

func GetLeadQuotes(leadId int) ([]types.LeadQuoteList, error) {
	var leads []types.LeadQuoteList

	query := `SELECT q.quote_id, e.name, v.name, q.event_date, q.guests, 
		SUM(qs.units * qs.price_per_unit::NUMERIC) AS total_price, 
		q.lead_id
	FROM quote AS q
	LEFT JOIN event_type AS e ON q.event_type_id = e.event_type_id
	LEFT JOIN venue_type AS v ON q.venue_type_id = v.venue_type_id
	LEFT JOIN quote_service qs ON qs.quote_id = q.quote_id
	WHERE q.lead_id = $1
	GROUP BY q.quote_id, e.name, v.name, q.event_date, q.guests, q.lead_id
	ORDER BY q.event_date ASC;`

	rows, err := DB.Query(query, leadId)
	if err != nil {
		return nil, fmt.Errorf("error executing query: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var lead types.LeadQuoteList
		var eventDate time.Time
		var venueType, eventType sql.NullString
		var guests sql.NullInt64
		var amount sql.NullFloat64

		err := rows.Scan(&lead.QuoteID,
			&venueType,
			&eventType,
			&eventDate,
			&guests,
			&amount,
			&lead.LeadID)
		if err != nil {
			return nil, fmt.Errorf("error scanning row: %w", err)
		}
		lead.EventDate = utils.FormatTimestampEST(eventDate.Unix())

		if venueType.Valid {
			lead.VenueType = venueType.String
		}

		if eventType.Valid {
			lead.EventType = eventType.String
		}

		if guests.Valid {
			lead.Guests = int(guests.Int64)
		}

		if amount.Valid {
			lead.Amount = amount.Float64
		}

		leads = append(leads, lead)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating rows: %w", err)
	}

	return leads, nil
}

func GetLeadQuoteDetails(quoteId string) (models.Quote, error) {
	query := `SELECT 
		lead_id,
		number_of_bartenders,
		guests,
		hours,
		event_type_id,
		venue_type_id,
		event_date,
		quote_id
	FROM quote 
	WHERE quote_id = $1`

	var quoteDetails models.Quote

	var leadID, bartenders, guests, hours, eventTypeID, venueTypeID sql.NullInt64
	var eventDate sql.NullTime

	row := DB.QueryRow(query, quoteId)

	err := row.Scan(
		&leadID,
		&bartenders,
		&guests,
		&hours,
		&eventTypeID,
		&venueTypeID,
		&eventDate,
		&quoteDetails.QuoteID,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return quoteDetails, fmt.Errorf("no quote found with ID %s", quoteId)
		}
		return quoteDetails, fmt.Errorf("error scanning row: %w", err)
	}

	if leadID.Valid {
		quoteDetails.LeadID = int(leadID.Int64)
	}
	if bartenders.Valid {
		quoteDetails.NumberOfBartenders = int(bartenders.Int64)
	}
	if guests.Valid {
		quoteDetails.Guests = int(guests.Int64)
	}
	if hours.Valid {
		quoteDetails.Hours = int(hours.Int64)
	}
	if eventTypeID.Valid {
		quoteDetails.EventTypeID = int(eventTypeID.Int64)
	}
	if venueTypeID.Valid {
		quoteDetails.VenueTypeID = int(venueTypeID.Int64)
	}
	if eventDate.Valid {
		quoteDetails.EventDate = eventDate.Time.Unix()
	}

	return quoteDetails, nil
}

func UpdateLeadQuote(form types.LeadQuoteForm) error {
	query := `
		UPDATE quote
		SET 
			number_of_bartenders = COALESCE($2, number_of_bartenders),
			guests = COALESCE($3, guests),
			hours = COALESCE($4, hours),
			event_type_id = COALESCE($5, event_type_id),
			venue_type_id = COALESCE($6, venue_type_id),
			event_date = COALESCE(to_timestamp($7)::timestamp AT TIME ZONE 'America/New_York', event_date)
		WHERE quote_id = $1
	`

	_, err := DB.Exec(
		query,
		utils.CreateNullInt(form.QuoteID),
		utils.CreateNullInt(form.NumberOfBartenders),
		utils.CreateNullInt(form.Guests),
		utils.CreateNullInt(form.Hours),
		utils.CreateNullInt(form.EventTypeID),
		utils.CreateNullInt(form.VenueTypeID),
		utils.CreateNullInt64(form.EventDate),
	)
	if err != nil {
		return fmt.Errorf("error updating lead quote data: %w", err)
	}

	return nil
}

func AssignStripeCustomerIDToLead(stripeCustomerId string, leadId int) error {
	if leadId == 0 {
		return fmt.Errorf("lead_id cannot be nil")
	}

	query := `
		UPDATE lead
		SET stripe_customer_id = $2
		WHERE lead_id = $1
	`
	_, err := DB.Exec(query, leadId, stripeCustomerId)
	if err != nil {
		return fmt.Errorf("failed to assign stripe customer id to lead: %v", err)
	}

	return nil
}

func CreateQuoteInvoice(stripeInvoiceId, invoiceUrl string, quoteId, invoiceTypeId int, dueDate int64) error {
	query := `
		INSERT INTO invoice (stripe_invoice_id, quote_id, invoice_type_id, url, due_date, date_created, invoice_status_id)
		VALUES ($1, $2, $3, $4, to_timestamp($5)::timestamptz AT TIME ZONE 'America/New_York', to_timestamp($6)::timestamptz AT TIME ZONE 'America/New_York', $7);
	`
	_, err := DB.Exec(query, stripeInvoiceId, quoteId, invoiceTypeId, invoiceUrl, dueDate, time.Now().Unix(), constants.OpenInvoiceStatusID)
	if err != nil {
		return fmt.Errorf("failed to create stripe invoice: %v", err)
	}

	return nil
}

func UpdateInvoiceStatus(stripeInvoiceId string, invoiceStatusId int) error {
	query := `
		UPDATE invoice
		SET invoice_status_id = $2
		WHERE stripe_invoice_id = $1
	`
	_, err := DB.Exec(query, stripeInvoiceId, invoiceStatusId)
	if err != nil {
		return fmt.Errorf("failed to update invoice status: %v", err)
	}

	return nil
}

func GetLeadQuoteInvoiceDetails(leadID, quoteId string) (types.QuoteDetails, error) {
	query := `SELECT 
		l.lead_id,
		l.full_name,
		l.phone_number,
		l.email,
		l.stripe_customer_id,
		q.event_date AT TIME ZONE 'America/New_York' AT TIME ZONE 'UTC' AS event_date_utc,
		q.external_id,
		i.invoice_id,
		SUM(qs.units * qs.price_per_unit::NUMERIC)
	FROM lead l
	JOIN quote AS q ON q.lead_id = l.lead_id
	LEFT JOIN quote_service qs ON qs.quote_id = q.quote_id
	LEFT JOIN invoice AS i ON i.quote_id = q.quote_id
	WHERE l.lead_id = $1 AND q.quote_id = $2
	GROUP BY l.lead_id, l.full_name, l.phone_number, l.email, l.stripe_customer_id, q.event_date, q.external_id, i.invoice_id`

	var quote types.QuoteDetails

	row := DB.QueryRow(query, leadID, quoteId)

	var email, stripeCustomerId sql.NullString
	var eventDate time.Time
	var invoiceId sql.NullInt64

	err := row.Scan(
		&quote.LeadID,
		&quote.FullName,
		&quote.PhoneNumber,
		&email,
		&stripeCustomerId,
		&eventDate,
		&quote.ExternalID,
		&invoiceId,
		&quote.Amount,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return quote, fmt.Errorf("no lead found with ID %s", leadID)
		}
		return quote, fmt.Errorf("error scanning row: %w", err)
	}

	if invoiceId.Valid {
		quote.InvoiceID = int(invoiceId.Int64)
	}
	if email.Valid {
		quote.Email = email.String
	}
	if stripeCustomerId.Valid {
		quote.StripeCustomerID = stripeCustomerId.String
	}

	quote.EventDate = eventDate.Unix()

	return quote, nil
}

func GetExternalQuoteDetails(externalQuoteId string) (types.ExternalQuoteDetails, error) {
	query := `SELECT 
		q.quote_id,
		number_of_bartenders,
		guests,
		hours,
		e.name AS event_type,
		v.name AS venue_type,
		event_date,
		SUM(qs.units * qs.price_per_unit::NUMERIC) AS full_amount,
		l.full_name,
		l.phone_number,
		l.email,
		i.url AS deposit_invoice_url,
		SUM(qs.units * qs.price_per_unit::NUMERIC) * it.amount_percentage AS deposit_amount,
		SUM(qs.units * qs.price_per_unit::NUMERIC) * (1 - it.amount_percentage) AS remaining_amount,
		(
			SELECT i2.url
			FROM invoice AS i2
			JOIN invoice_type AS it2 ON it2.invoice_type_id = i2.invoice_type_id
			JOIN invoice_status AS stat2 ON stat2.invoice_status_id = i2.invoice_status_id AND i2.invoice_status_id = $2
			WHERE i2.quote_id = q.quote_id AND i2.invoice_type_id = $4
			LIMIT 1
		) AS full_invoice_url,
		(
			SELECT i3.url
			FROM invoice AS i3
			JOIN invoice_type AS it3 ON it3.invoice_type_id = i3.invoice_type_id
			JOIN invoice_status AS stat3 ON stat3.invoice_status_id = i3.invoice_status_id AND i3.invoice_status_id = $2
			WHERE i3.quote_id = q.quote_id AND i3.invoice_type_id = $5
			LIMIT 1
		) AS remaining_invoice_url,
		EXISTS (
			SELECT 1 FROM invoice AS deposit_invoice
			JOIN invoice_type AS deposit_invoice_type ON deposit_invoice_type.invoice_type_id = deposit_invoice.invoice_type_id AND deposit_invoice.invoice_type_id = $3
			JOIN invoice_status AS deposit_invoice_status ON deposit_invoice_status.invoice_status_id = deposit_invoice.invoice_status_id AND deposit_invoice.invoice_status_id = $6
			WHERE deposit_invoice.quote_id = q.quote_id
		) AS is_deposit_paid
	FROM quote AS q
	LEFT JOIN event_type AS e ON q.event_type_id = e.event_type_id
	LEFT JOIN venue_type AS v ON q.venue_type_id = v.venue_type_id
	JOIN lead AS l ON q.lead_id = l.lead_id
	JOIN invoice AS i ON i.quote_id = q.quote_id
	JOIN invoice_type AS it ON it.invoice_type_id = i.invoice_type_id AND it.invoice_type_id = $3
	LEFT JOIN quote_service qs ON qs.quote_id = q.quote_id
	WHERE q.external_id = $1
	GROUP BY q.quote_id, number_of_bartenders, guests, hours, e.name, v.name, event_date, 
			l.full_name, l.phone_number, l.email, i.url, it.amount_percentage, i.date_created
	ORDER BY i.date_created DESC
	LIMIT 1;`

	var quoteDetails types.ExternalQuoteDetails

	var bartenders, guests, hours sql.NullInt64
	var eventDate sql.NullTime
	var eventType, venueType, email sql.NullString

	row := DB.QueryRow(query, externalQuoteId, constants.OpenInvoiceStatusID, constants.DepositInvoiceTypeID, constants.FullInvoiceTypeID, constants.RemainingInvoiceTypeID, constants.PaidInvoiceStatusID)

	err := row.Scan(
		&quoteDetails.QuoteID,
		&bartenders,
		&guests,
		&hours,
		&eventType,
		&venueType,
		&eventDate,
		&quoteDetails.Amount,
		&quoteDetails.FullName,
		&quoteDetails.PhoneNumber,
		&email,
		&quoteDetails.DepositInvoiceURL,
		&quoteDetails.Deposit,
		&quoteDetails.RemainingAmount,
		&quoteDetails.FullInvoiceURL,
		&quoteDetails.RemainingInvoiceURL,
		&quoteDetails.IsDepositPaid,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return quoteDetails, fmt.Errorf("no quote found with external id: %s", externalQuoteId)
		}
		return quoteDetails, fmt.Errorf("error scanning row: %w", err)
	}

	if bartenders.Valid {
		quoteDetails.NumberOfBartenders = int(bartenders.Int64)
	}
	if guests.Valid {
		quoteDetails.Guests = int(guests.Int64)
	}
	if hours.Valid {
		quoteDetails.Hours = int(hours.Int64)
	}
	if eventType.Valid {
		quoteDetails.EventType = eventType.String
	}
	if venueType.Valid {
		quoteDetails.VenueType = venueType.String
	}
	if eventDate.Valid {
		quoteDetails.EventDate = utils.FormatTimestamp(eventDate.Time.Unix())
	}
	if email.Valid {
		quoteDetails.Email = email.String
	}

	return quoteDetails, nil
}

func GetInvoiceByStripeInvoiceID(stripeInvoiceId string) (models.Invoice, error) {
	query := `SELECT i.invoice_id,
	i.quote_id,
	i.date_paid,
	i.due_date,
	i.invoice_type_id,
	i.url,
	i.stripe_invoice_id
	FROM invoice i
	WHERE i.stripe_invoice_id = $1`

	var invoice models.Invoice

	row := DB.QueryRow(query, stripeInvoiceId)

	var datePaid, dueDate time.Time
	var url sql.NullString

	err := row.Scan(
		&invoice.InvoiceID,
		&invoice.QuoteID,
		&datePaid,
		&dueDate,
		&invoice.InvoiceTypeID,
		&url,
		&invoice.StripeInvoiceID,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return invoice, fmt.Errorf("no invoice found with Stripe Invoice ID %s", stripeInvoiceId)
		}
		return invoice, fmt.Errorf("error scanning row: %w", err)
	}

	if url.Valid {
		invoice.URL = url.String
	}

	invoice.DatePaid = datePaid.Unix()
	invoice.DueDate = dueDate.Unix()

	return invoice, nil
}

func SetInvoiceStatusToPaid(stripeInvoiceId string, datePaid int64) error {
	query := `
		UPDATE invoice
		SET date_paid = $2,
		invoice_status_id = $3
		WHERE stripe_invoice_id = $1
	`
	_, err := DB.Exec(query, stripeInvoiceId, datePaid, constants.PaidInvoiceStatusID)
	if err != nil {
		return fmt.Errorf("failed to assign stripe customer id to lead: %v", err)
	}

	return nil
}

func GetQuoteDetailsByStripeInvoiceID(stripeInvoiceId string) (types.InvoiceQuoteDetails, error) {
	query := `SELECT 
		l.lead_id,
		q.event_type_id,
		q.venue_type_id,
		SUM(qs.units * qs.price_per_unit::NUMERIC),
		q.guests,
		l.phone_number,
		q.event_date,
		q.quote_id
	FROM invoice i
	JOIN quote q ON i.stripe_invoice_id = q.stripe_invoice_id
	JOIN lead l ON l.lead_id = q.lead_id
	LEFT JOIN quote_service qs ON qs.quote_id = q.quote_id
	WHERE i.stripe_invoice_id = $1
	GROUP BY l.lead_id, q.event_type_id, q.venue_type_id, q.guests, l.phone_number, q.event_date, q.quote_id;`

	var invoiceQuoteDetails types.InvoiceQuoteDetails

	row := DB.QueryRow(query, stripeInvoiceId)

	var eventTypeId, venueTypeId sql.NullInt64
	var eventDate time.Time

	err := row.Scan(
		&invoiceQuoteDetails.LeadID,
		&eventTypeId,
		&venueTypeId,
		&invoiceQuoteDetails.Amount,
		&invoiceQuoteDetails.Guests,
		&invoiceQuoteDetails.PhoneNumber,
		&eventDate,
		&invoiceQuoteDetails.QuoteID,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return invoiceQuoteDetails, fmt.Errorf("no lead found with stripeInvoiceId: %s", stripeInvoiceId)
		}
		return invoiceQuoteDetails, fmt.Errorf("error scanning row: %w", err)
	}

	if eventTypeId.Valid {
		invoiceQuoteDetails.EventTypeID = int(eventTypeId.Int64)
	}

	if venueTypeId.Valid {
		invoiceQuoteDetails.VenueTypeID = int(venueTypeId.Int64)
	}

	invoiceQuoteDetails.EventDate = eventDate.Unix()

	return invoiceQuoteDetails, nil
}

func GetLeadQuoteInvoices(quoteId int) ([]types.LeadQuoteInvoice, error) {
	var leadQuoteInvoices []types.LeadQuoteInvoice

	query := `SELECT 
		i.stripe_invoice_id, 
		l.stripe_customer_id, 
		SUM(qs.units * qs.price_per_unit::NUMERIC),
		i.due_date,
		CASE 
			WHEN i.invoice_type_id = 1 THEN 0.25
			WHEN i.invoice_type_id = 2 THEN 0.75
			ELSE 1.00
		END AS invoice_multiplier,
		i.invoice_type_id
	FROM quote AS q
	LEFT JOIN quote_service qs ON qs.quote_id = q.quote_id
	JOIN invoice AS i ON i.quote_id = q.quote_id
	JOIN lead AS l ON l.lead_id = q.lead_id
	WHERE q.quote_id = $1 AND i.invoice_status_id = $2
	GROUP BY 
		i.stripe_invoice_id, 
		l.stripe_customer_id, 
		i.due_date, 
		i.invoice_type_id;`

	rows, err := DB.Query(query, quoteId, constants.OpenInvoiceStatusID)
	if err != nil {
		return leadQuoteInvoices, fmt.Errorf("error executing query: %w", err)
	}
	defer rows.Close()

	var invoiceDueDate time.Time

	for rows.Next() {
		var leadQuoteInvoice types.LeadQuoteInvoice
		var amount sql.NullFloat64

		err := rows.Scan(
			&leadQuoteInvoice.StripeInvoiceID,
			&leadQuoteInvoice.StripeCustomerID,
			&amount,
			&invoiceDueDate,
			&leadQuoteInvoice.InvoiceTypeMultiplier,
			&leadQuoteInvoice.InvoiceTypeID,
		)
		if err != nil {
			return leadQuoteInvoices, fmt.Errorf("error scanning row: %w", err)
		}

		if amount.Valid {
			leadQuoteInvoice.Amount = amount.Float64
		}

		leadQuoteInvoice.DueDate = invoiceDueDate.Unix()

		leadQuoteInvoices = append(leadQuoteInvoices, leadQuoteInvoice)
	}

	if err := rows.Err(); err != nil {
		return leadQuoteInvoices, fmt.Errorf("error iterating rows: %w", err)
	}

	return leadQuoteInvoices, nil
}

func GetInvoiceTypes() ([]models.InvoiceType, error) {
	var invoiceTypes []models.InvoiceType

	rows, err := DB.Query(`SELECT invoice_type_id, type, amount_percentage FROM "invoice_type"`)
	if err != nil {
		return invoiceTypes, fmt.Errorf("error executing query: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var invType models.InvoiceType
		err := rows.Scan(&invType.InvoiceTypeID, &invType.Type, &invType.AmountPercentage)
		if err != nil {
			return invoiceTypes, fmt.Errorf("error scanning row: %w", err)
		}
		invoiceTypes = append(invoiceTypes, invType)
	}

	if err := rows.Err(); err != nil {
		return invoiceTypes, fmt.Errorf("error iterating rows: %w", err)
	}

	return invoiceTypes, nil
}

func SetOpenInvoicesToVoid(quoteId int) error {
	query := `
		UPDATE invoice AS i
		SET i.invoice_status_id = $2
		FROM quote AS q
		WHERE q.quote_id = i.quote_id
		AND q.quote_id = $1
		AND i.invoice_status_id <> $3;
	`
	_, err := DB.Exec(query, quoteId, constants.VoidInvoiceStatusID, constants.PaidInvoiceStatusID)
	if err != nil {
		return fmt.Errorf("failed to assign stripe customer id to lead: %v", err)
	}

	return nil
}

func GetLeadStatusList() ([]models.LeadStatus, error) {
	var leadStatuses []models.LeadStatus

	rows, err := DB.Query(`SELECT lead_status_id, status FROM "lead_status"`)
	if err != nil {
		return leadStatuses, fmt.Errorf("error executing query: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var ls models.LeadStatus
		err := rows.Scan(&ls.LeadStatusID, &ls.Status)
		if err != nil {
			return leadStatuses, fmt.Errorf("error scanning row: %w", err)
		}
		leadStatuses = append(leadStatuses, ls)
	}

	if err := rows.Err(); err != nil {
		return leadStatuses, fmt.Errorf("error iterating rows: %w", err)
	}

	return leadStatuses, nil
}

func GetLeadInterestList() ([]models.LeadInterest, error) {
	var leadInterests []models.LeadInterest

	rows, err := DB.Query(`SELECT lead_interest_id, interest FROM "lead_interest"`)
	if err != nil {
		return leadInterests, fmt.Errorf("error executing query: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var li models.LeadInterest
		err := rows.Scan(&li.LeadInterestID, &li.Interest)
		if err != nil {
			return leadInterests, fmt.Errorf("error scanning row: %w", err)
		}
		leadInterests = append(leadInterests, li)
	}

	if err := rows.Err(); err != nil {
		return leadInterests, fmt.Errorf("error iterating rows: %w", err)
	}

	return leadInterests, nil
}

func GetNextActionList() ([]models.NextAction, error) {
	var nextActions []models.NextAction

	rows, err := DB.Query(`SELECT next_action_id, action FROM "next_action"`)
	if err != nil {
		return nextActions, fmt.Errorf("error executing query: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var na models.NextAction
		err := rows.Scan(&na.NextActionID, &na.Action)
		if err != nil {
			return nextActions, fmt.Errorf("error scanning row: %w", err)
		}
		nextActions = append(nextActions, na)
	}

	if err := rows.Err(); err != nil {
		return nextActions, fmt.Errorf("error iterating rows: %w", err)
	}

	return nextActions, nil
}

func DeleteService(id int) error {
	sqlStatement := `
        DELETE FROM service WHERE service_id = $1
    `
	_, err := DB.Exec(sqlStatement, id)
	if err != nil {
		return err
	}

	return nil
}

func GetServicesList(pageNum int) ([]models.Service, int, error) {
	var services []models.Service
	var totalRows int

	offset := (pageNum - 1) * int(constants.LeadsPerPage)

	rows, err := DB.Query(`SELECT service_id, service, suggested_price::NUMERIC
			COUNT(*) OVER() AS total_rows
			FROM "service"
			OFFSET $1
			LIMIT $2`, offset, constants.LeadsPerPage)
	if err != nil {
		return services, totalRows, fmt.Errorf("error executing query: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var service models.Service
		var suggestedPrice sql.NullFloat64
		err := rows.Scan(&service.ServiceID, &service.Service, &suggestedPrice, &totalRows)
		if err != nil {
			return services, totalRows, fmt.Errorf("error scanning row: %w", err)
		}
		if suggestedPrice.Valid {
			service.SuggestedPrice = suggestedPrice.Float64
		}
		services = append(services, service)
	}

	if err := rows.Err(); err != nil {
		return services, totalRows, fmt.Errorf("error iterating rows: %w", err)
	}

	return services, totalRows, nil
}

func GetServices() ([]models.Service, error) {
	var services []models.Service

	rows, err := DB.Query(`SELECT service_id, service, suggested_price::NUMERIC FROM "service";`)
	if err != nil {
		return services, fmt.Errorf("error executing query: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var service models.Service
		var suggestedPrice sql.NullFloat64
		err := rows.Scan(&service.ServiceID, &service.Service, &suggestedPrice)
		if err != nil {
			return services, fmt.Errorf("error scanning row: %w", err)
		}
		if suggestedPrice.Valid {
			service.SuggestedPrice = suggestedPrice.Float64
		}
		services = append(services, service)
	}

	if err := rows.Err(); err != nil {
		return services, fmt.Errorf("error iterating rows: %w", err)
	}

	return services, nil
}

func CreateService(form types.ServiceForm) error {
	stmt, err := DB.Prepare(`
		INSERT INTO service (service, suggested_price) VALUES ($1, $2)
	`)
	if err != nil {
		return fmt.Errorf("error preparing statement: %w", err)
	}
	defer stmt.Close()

	_, err = stmt.Exec(utils.CreateNullString(form.Service), utils.CreateNullFloat64(form.SuggestedPrice))
	if err != nil {
		return fmt.Errorf("error executing statement: %w", err)
	}

	return nil
}

func UpdateService(form types.ServiceForm) error {
	stmt, err := DB.Prepare(`
		UPDATE service
		SET service = COALESCE($1, service), suggested_price = $2
		WHERE service_id = $3
	`)
	if err != nil {
		return fmt.Errorf("error preparing statement: %w", err)
	}
	defer stmt.Close()

	_, err = stmt.Exec(utils.CreateNullString(form.Service), utils.CreateNullFloat64(form.SuggestedPrice), utils.CreateNullInt(form.ServiceID))
	if err != nil {
		return fmt.Errorf("error executing statement: %w", err)
	}

	return nil
}

func GetQuoteServices(quoteId int) ([]types.QuoteServiceList, error) {
	var quoteServiceList []types.QuoteServiceList

	query := `SELECT qs.quote_id, 
		qs.service_id, 
		s.service, 
		qs.units, 
		qs.price_per_unit::NUMERIC, 
		(qs.price_per_unit::NUMERIC * qs.units),
		qs.quote_service_id
	FROM quote_service AS qs
	JOIN service AS s ON qs.service_id = s.service_id
	WHERE qs.quote_id = $1;`

	rows, err := DB.Query(query, quoteId)
	if err != nil {
		return nil, fmt.Errorf("error executing query: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var service types.QuoteServiceList

		err := rows.Scan(&service.QuoteID,
			&service.ServiceID,
			&service.Service,
			&service.Units,
			&service.PricePerUnit,
			&service.Total,
			&service.QuoteServiceID)
		if err != nil {
			return nil, fmt.Errorf("error scanning row: %w", err)
		}

		quoteServiceList = append(quoteServiceList, service)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating rows: %w", err)
	}

	return quoteServiceList, nil
}

func DeleteQuoteService(id int) error {
	sqlStatement := `
        DELETE FROM quote_service WHERE quote_service_id = $1
    `
	_, err := DB.Exec(sqlStatement, id)
	if err != nil {
		return err
	}

	return nil
}

func CreateQuoteService(form types.QuoteServiceForm) error {
	stmt, err := DB.Prepare(`
		INSERT INTO quote_service (service_id, quote_id, units, price_per_unit) VALUES ($1, $2, $3, $4)
	`)
	if err != nil {
		return fmt.Errorf("error preparing statement: %w", err)
	}
	defer stmt.Close()

	_, err = stmt.Exec(
		utils.CreateNullInt(form.ServiceID),
		utils.CreateNullInt(form.QuoteID),
		utils.CreateNullInt(form.Units),
		utils.CreateNullFloat64(form.PricePerUnit),
	)
	if err != nil {
		return fmt.Errorf("error executing statement: %w", err)
	}

	return nil
}

func UpdateQuoteService(form types.QuoteServiceForm) error {
	stmt, err := DB.Prepare(`
		UPDATE quote_service
		SET price_per_unit = COALESCE($1, price_per_unit),
		units = COALESCE($2, units)
		WHERE quote_service_id = $3
	`)
	if err != nil {
		return fmt.Errorf("error preparing statement: %w", err)
	}
	defer stmt.Close()

	_, err = stmt.Exec(
		utils.CreateNullFloat64(form.PricePerUnit),
		utils.CreateNullInt(form.Units),
		utils.CreateNullInt(form.QuoteServiceID),
	)
	if err != nil {
		return fmt.Errorf("error executing statement: %w", err)
	}

	return nil
}

func GetQuoteIDByQuoteServiceID(quoteServiceId int) (int, error) {
	var quoteId int

	query := `SELECT qs.quote_id FROM quote_service AS qs WHERE qs.quote_service_id = $1;`

	rows, err := DB.Query(query, quoteServiceId)
	if err != nil {
		return quoteId, err
	}
	defer rows.Close()

	for rows.Next() {
		err := rows.Scan(&quoteId)
		if err != nil {
			return quoteId, err
		}
	}

	if err := rows.Err(); err != nil {
		return quoteId, err
	}

	return quoteId, err
}

func CheckQuoteHasInvoiceID(quote int) (bool, error) {
	var hasInvoice bool

	query := `SELECT EXISTS(SELECT 1 FROM invoice AS i WHERE i.quote_id = $1 AND i.stripe_invoice_id IS NOT NULL)`

	err := DB.QueryRow(query, quote).Scan(&hasInvoice)
	if err != nil {
		return false, err
	}

	return hasInvoice, nil
}

func GetMessagesByLeadID(leadId int) ([]types.FrontendMessage, error) {
	var messages []types.FrontendMessage

	query := `SELECT l.full_name as client_name,
	CONCAT(u.first_name, ' ', u.last_name) as user_name,
	m.text,
	m.date_created,
	m.is_inbound,
	m.message_id,
	l.lead_id,
	m.is_read
	FROM "message" AS m
	JOIN "lead" AS l ON l.phone_number IN (m.text_from, m.text_to)
	JOIN "user" AS u  ON u.phone_number IN (m.text_from, m.text_to)
	WHERE l.lead_id = $1
	ORDER BY m.date_created ASC;`

	rows, err := DB.Query(query, leadId)
	if err != nil {
		fmt.Printf("%+v\n", err)
		return messages, err
	}
	defer rows.Close()

	for rows.Next() {
		var dateCreated time.Time

		var message types.FrontendMessage
		err := rows.Scan(
			&message.ClientName,
			&message.UserName,
			&message.Message,
			&dateCreated,
			&message.IsInbound,
			&message.MessageID,
			&message.LeadID,
			&message.IsRead,
		)
		if err != nil {
			fmt.Printf("%+v\n", err)
			return messages, err
		}

		message.DateCreated = utils.FormatTimestamp(dateCreated.Unix())
		messages = append(messages, message)
	}

	if err = rows.Err(); err != nil {
		return messages, err
	}

	return messages, nil
}

func CreateLeadNote(note models.LeadNote) error {
	stmt, err := DB.Prepare(`
		INSERT INTO lead_note (note, lead_id, date_added, added_by_user_id)
		VALUES ($1, $2, to_timestamp($3)::timestamptz AT TIME ZONE 'America/New_York', $4)
	`)
	if err != nil {
		return fmt.Errorf("error preparing statement: %w", err)
	}
	defer stmt.Close()

	leadID := utils.CreateNullInt(&note.LeadID)

	_, err = stmt.Exec(note.Note, leadID, note.DateAdded, note.AddedByUserID)
	if err != nil {
		return fmt.Errorf("error executing statement: %w", err)
	}

	return nil
}

func GetLeadNotesByLeadID(leadId int) ([]types.FrontendNote, error) {
	var notes []types.FrontendNote

	query := `SELECT u.username,
	n.note,
	n.date_added
	FROM "lead_note" AS n
	JOIN "user" AS u ON u.user_id = n.added_by_user_id
	WHERE n.lead_id = $1
	ORDER BY n.date_added DESC`

	rows, err := DB.Query(query, leadId)
	if err != nil {
		fmt.Printf("%+v\n", err)
		return notes, err
	}
	defer rows.Close()

	for rows.Next() {
		var dateAdded time.Time

		var note types.FrontendNote
		err := rows.Scan(
			&note.UserName,
			&note.Note,
			&dateAdded,
		)
		if err != nil {
			fmt.Printf("%+v\n", err)
			return notes, err
		}

		note.DateAdded = utils.FormatTimestamp(dateAdded.Unix())
		notes = append(notes, note)
	}

	if err = rows.Err(); err != nil {
		return notes, err
	}

	return notes, nil
}

func GetUsersWithMessages() ([]types.FrontendMessage, error) {
	var messages []types.FrontendMessage

	query := `SELECT DISTINCT ON (l.lead_id) 
		l.lead_id, l.full_name
	FROM "message" AS m
	JOIN "lead" AS l ON l.phone_number IN (m.text_from, m.text_to)
	JOIN "user" AS u  ON u.phone_number IN (m.text_from, m.text_to)
	ORDER BY l.lead_id, m.date_created ASC;`

	rows, err := DB.Query(query)
	if err != nil {
		return messages, err
	}
	defer rows.Close()

	for rows.Next() {
		var message types.FrontendMessage
		err := rows.Scan(
			&message.LeadID,
			&message.ClientName,
		)
		if err != nil {
			fmt.Printf("%+v\n", err)
			return messages, err
		}

		messages = append(messages, message)
	}

	if err = rows.Err(); err != nil {
		return messages, err
	}

	return messages, nil
}
