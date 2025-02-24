package services

import (
	"fmt"
	"os"
	"time"

	"github.com/davidalvarez305/yd_cocktails/constants"
	"github.com/davidalvarez305/yd_cocktails/database"
	"github.com/davidalvarez305/yd_cocktails/helpers"
	"github.com/davidalvarez305/yd_cocktails/models"
	"github.com/google/uuid"
)

const (
	audioTranscriptionS3Path = "uploads/audio/"
	textTranscriptionS3Path  = "uploads/transcription/"
)

func TranscribePhoneCall(phoneCall models.PhoneCall) error {
	// Download the file from Twilio
	audioFileName := uuid.New().String() + ".mp3"
	localAudioFilePath := constants.LOCAL_FILES_DIR + audioFileName

	err := DownloadFileFromTwilio(phoneCall.RecordingURL, localAudioFilePath)
	if err != nil {
		fmt.Printf("ERROR DOWNLOADING AUDIO FILE: %+v\n", err)
		return err
	}

	// Get file info before opening
	audioFileInfo, err := os.Stat(localAudioFilePath)
	if err != nil {
		fmt.Printf("ERROR GETTING FILE INFO: %+v\n", err)
		return err
	}

	// Open the file before uploading to S3
	audioFile, err := os.Open(localAudioFilePath)
	if err != nil {
		fmt.Printf("ERROR OPENING FILE: %+v\n", err)
		return err
	}
	defer audioFile.Close()

	// Upload audio file to S3
	audioS3FilePath := audioTranscriptionS3Path + audioFileName
	err = UploadFileToS3(audioFile, audioFileInfo.Size(), audioS3FilePath)
	if err != nil {
		fmt.Printf("ERROR UPLOADING AUDIO TO S3: %+v\n", err)
		return err
	}

	// Transcribe audio file
	audioFileS3URL := "s3://" + constants.AWSS3BucketName + "/" + audioS3FilePath
	transcriptionID, transcriptionText, err := TranscribeAudio(audioFileS3URL)
	if err != nil {
		fmt.Printf("ERROR TRANSCRIBING AUDIO: %+v\n", err)
		return err
	}

	// Upload transcription text to S3
	transcriptionFileName := uuid.New().String() + ".txt"

	// Define the local file path where the transcription text will be saved
	localTranscriptionTextPath := constants.LOCAL_FILES_DIR + transcriptionFileName

	// Create the local file and write the transcription text to it
	transcriptionFile, err := os.Create(localTranscriptionTextPath)
	if err != nil {
		fmt.Printf("ERROR CREATING LOCAL FILE: %+v\n", err)
		return err
	}
	defer transcriptionFile.Close()

	_, err = transcriptionFile.Write([]byte(transcriptionText))
	if err != nil {
		fmt.Printf("ERROR WRITING TO LOCAL FILE: %+v\n", err)
		return err
	}

	// Open the file for uploading to S3
	transcriptionFileToUpload, err := os.Open(localTranscriptionTextPath)
	if err != nil {
		fmt.Printf("ERROR OPENING FILE: %+v\n", err)
		return err
	}
	defer transcriptionFileToUpload.Close()

	// Determine file size
	transcriptionFileInfo, err := transcriptionFileToUpload.Stat()
	if err != nil {
		fmt.Printf("ERROR GETTING FILE INFO: %+v\n", err)
		return err
	}
	transcriptionFileSize := transcriptionFileInfo.Size()

	// S3 file path
	s3TranscriptionFilePath := textTranscriptionS3Path + transcriptionFileName

	// Upload the file to S3
	err = UploadFileToS3(transcriptionFileToUpload, transcriptionFileSize, s3TranscriptionFilePath)
	if err != nil {
		fmt.Printf("ERROR UPLOADING TRANSCRIPTION TO S3: %+v\n", err)
		return err
	}

	// Save transcription in DB
	transcription := models.PhoneCallTranscription{
		PhoneCallID: phoneCall.PhoneCallID,
		ExternalID:  transcriptionID,
		Text:        transcriptionText,
		AudioURL:    audioFileName,
		TextURL:     transcriptionFileName,
	}

	err = database.CreatePhoneCallTranscription(transcription)
	if err != nil {
		fmt.Printf("ERROR SAVING TRANSCRIPTION: %+v\n", err)
		return err
	}

	err = SummarizePhoneCall(phoneCall, transcriptionText)
	if err != nil {
		fmt.Printf("ERROR SAVING TRANSCRIPTION: %+v\n", err)
		return err
	}

	err = helpers.DeleteFilesInDirectory(constants.LOCAL_FILES_DIR)
	if err != nil {
		fmt.Printf("ERROR DELETING LOCAL FILES: %+v\n", err)
		return err
	}

	return nil
}

func SummarizePhoneCall(phoneCall models.PhoneCall, transcriptionText string) error {
	summary, err := GetOpenAICompletionsResponse(fmt.Sprintf("This was a sales call for a bartending service. Summarize the key points in the following text: %s", transcriptionText))
	if err != nil {
		fmt.Printf("ERROR GENERATING CALL SUMMARY: %+v\n", err)
		return err
	}

	crmUserPhoneNumber := phoneCall.CallFrom
	if phoneCall.IsInbound {
		crmUserPhoneNumber = phoneCall.CallTo
	}

	userId, err := database.GetUserIDFromPhoneNumber(crmUserPhoneNumber)
	if err != nil {
		fmt.Printf("ERROR GETTING USER ID FROM PHONE NUMBER: %+v\n", err)
		return err
	}

	leadPhoneNumber := phoneCall.CallTo
	if phoneCall.IsInbound {
		leadPhoneNumber = phoneCall.CallFrom
	}

	leadId, err := database.GetLeadIDFromPhoneNumber(leadPhoneNumber)
	if err != nil {
		fmt.Printf("ERROR GETTING LEAD ID FROM PHONE NUMBER: %+v\n", err)
		return err
	}

	leadNote := models.LeadNote{
		Note:          summary,
		LeadID:        leadId,
		DateAdded:     time.Now().Unix(),
		AddedByUserID: userId,
	}

	err = database.CreateLeadNote(leadNote)
	if err != nil {
		fmt.Printf("ERROR SAVING TRANSCRIPTION: %+v\n", err)
		return err
	}

	return nil
}

func checkPhoneCallTranscription() {
	phoneCalls, err := database.GetPhoneCallsWithoutTranscription()
	if err != nil {
		fmt.Printf("ERROR GETTING PHONE CALLS WITHOUT TRANSCRIPTION: %+v\n", err)
		return
	}

	// Not sure if it would be a good idea to do this concurrently...
	for _, phoneCall := range phoneCalls {
		err = TranscribePhoneCall(phoneCall)

		if err != nil {
			continue
		}
	}
}

func StartTranscriptionService() {
	go func() {
		for {
			checkPhoneCallTranscription()

			// Sleep for one minute before the next run
			time.Sleep(5 * time.Minute)
		}
	}()
}
