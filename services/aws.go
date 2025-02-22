package services

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/transcribe"
	"github.com/aws/aws-sdk-go-v2/service/transcribe/types"
	"github.com/davidalvarez305/yd_cocktails/constants"
	"github.com/google/uuid"
)

func UploadFileToS3(file multipart.File, fileSize int64, s3FilePath string) error {
	cfg, err := config.LoadDefaultConfig(context.TODO(), config.WithRegion(constants.AWSRegion))
	if err != nil {
		return fmt.Errorf("failed to load AWS config: %w", err)
	}

	client := s3.NewFromConfig(cfg)

	buffer := make([]byte, fileSize)
	_, err = file.Read(buffer)
	if err != nil {
		return fmt.Errorf("failed to read file: %w", err)
	}
	fileType := http.DetectContentType(buffer)

	_, err = client.PutObject(context.TODO(), &s3.PutObjectInput{
		Bucket:        aws.String(constants.AWSS3BucketName),
		Key:           aws.String(s3FilePath),
		Body:          bytes.NewReader(buffer),
		ContentLength: &fileSize,
		ContentType:   aws.String(fileType),
	})
	if err != nil {
		return fmt.Errorf("failed to upload file to S3: %w", err)
	}

	return nil
}

func DownloadFileFromS3(s3FilePath, localFilePath string) (string, error) {
	cfg, err := config.LoadDefaultConfig(context.TODO(), config.WithRegion(constants.AWSRegion))
	if err != nil {
		return "", fmt.Errorf("failed to load AWS config: %w", err)
	}

	client := s3.NewFromConfig(cfg)

	result, err := client.GetObject(context.TODO(), &s3.GetObjectInput{
		Bucket: aws.String(constants.AWSS3BucketName),
		Key:    aws.String(s3FilePath),
	})
	if err != nil {
		return "", fmt.Errorf("failed to download file from S3: %w", err)
	}
	defer result.Body.Close()

	outFile, err := os.Create(localFilePath)
	if err != nil {
		return "", fmt.Errorf("failed to create local file: %w", err)
	}
	defer outFile.Close()

	_, err = io.Copy(outFile, result.Body)
	if err != nil {
		return "", fmt.Errorf("failed to save file locally: %w", err)
	}

	fmt.Printf("File downloaded successfully to %s\n", localFilePath)
	return localFilePath, nil
}

func TranscribeAudio(audioFileURL string) (string, string, error) {
	cfg, err := config.LoadDefaultConfig(context.TODO(), config.WithRegion(constants.AWSRegion))
	if err != nil {
		return "", "", fmt.Errorf("failed to load AWS config: %w", err)
	}

	client := transcribe.NewFromConfig(cfg)

	jobName := uuid.New().String()

	jobInput := &transcribe.StartTranscriptionJobInput{
		TranscriptionJobName:      aws.String(jobName),
		LanguageOptions:           []types.LanguageCode{types.LanguageCodeEsUs, types.LanguageCodeEnUs},
		IdentifyMultipleLanguages: aws.Bool(true),
		Media: &types.Media{
			MediaFileUri: aws.String(audioFileURL),
		},
		OutputBucketName: aws.String(constants.AWSS3BucketName),
		OutputKey:        aws.String("uploads/jobs/" + jobName + ".json"),
	}

	_, err = client.StartTranscriptionJob(context.TODO(), jobInput)
	if err != nil {
		return "", "", fmt.Errorf("failed to start transcription job: %w", err)
	}

	for {
		job, err := client.GetTranscriptionJob(context.TODO(), &transcribe.GetTranscriptionJobInput{
			TranscriptionJobName: aws.String(jobName),
		})
		if err != nil {
			return "", "", fmt.Errorf("failed to check job status: %w", err)
		}

		status := job.TranscriptionJob.TranscriptionJobStatus
		if status == types.TranscriptionJobStatusCompleted {
			transcriptURL := *job.TranscriptionJob.Transcript.TranscriptFileUri
			transcriptText, err := fetchTranscript(transcriptURL)
			if err != nil {
				return jobName, "", fmt.Errorf("failed to fetch transcript: %w", err)
			}
			return jobName, transcriptText, nil
		}

		if status == types.TranscriptionJobStatusFailed {
			return jobName, "", fmt.Errorf("transcription job failed")
		}

		time.Sleep(5 * time.Second)
	}
}

func fetchTranscript(transcriptURL string) (string, error) {
	resp, err := http.Get(transcriptURL)
	if err != nil {
		return "", fmt.Errorf("failed to fetch transcript file: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("non-200 response from transcript file: %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to read transcript file: %w", err)
	}

	var transcriptData struct {
		Results struct {
			Transcripts []struct {
				Transcript string `json:"transcript"`
			} `json:"transcripts"`
		} `json:"results"`
	}

	err = json.Unmarshal(body, &transcriptData)
	if err != nil {
		return "", fmt.Errorf("failed to parse transcript JSON: %w", err)
	}

	if len(transcriptData.Results.Transcripts) > 0 {
		return transcriptData.Results.Transcripts[0].Transcript, nil
	}

	return "", fmt.Errorf("no transcript found in response")
}
