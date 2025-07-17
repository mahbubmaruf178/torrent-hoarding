package main

import (

	// "time" // time is not used, so it remains commented out

	"context"
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"os"
	"path/filepath"
	"strconv"
	"time"

	"github.com/gotd/td/telegram/uploader"
	ffmpeg "github.com/u2takey/ffmpeg-go"
)

func VideoInfo(filePath string) (length string, width string, height string, err error) {
	// Use ffprobe to get video information
	log.Printf("Probing video file: %s\n", filePath)
	probe, err := ffmpeg.Probe(filePath)
	log.Printf("Probing video file: %s\n", filePath)
	if err != nil {
		log.Printf("Error probing video file: %v\n", err)
		return "", "", "", err
	}

	// Parse the JSON output to extract video information
	// This is a simplified implementation - you may need to handle the JSON more carefully
	data := struct {
		Streams []struct {
			CodecType string `json:"codec_type"`
			Width     int    `json:"width"`
			Height    int    `json:"height"`
		} `json:"streams"`
		Format struct {
			Duration string `json:"duration"`
		} `json:"format"`
	}{}

	if err := json.Unmarshal([]byte(probe), &data); err != nil {
		return "", "", "", err
	}

	// Find the video stream to get width and height
	for _, stream := range data.Streams {
		if stream.CodecType == "video" {
			width = strconv.Itoa(stream.Width)
			height = strconv.Itoa(stream.Height)
			break
		}
	}

	// Get the duration from the format section
	length = data.Format.Duration

	return length, width, height, nil
}

func ExtractRandomFrames(videoFilepath string, length string, amount int) ([]string, error) {
	// Get correct output directory path
	framepath := "frames"
	log.Printf("Output directory for frames: %s\n", framepath)

	err := os.MkdirAll(framepath, 0755)
	if err != nil {
		log.Printf("Error creating output directory: %v\n", err)
		return nil, fmt.Errorf("failed to create output directory: %w", err)
	}

	filepathofframes := make([]string, amount)

	// Parse video duration
	duration, err := strconv.ParseFloat(length, 64)
	if err != nil {
		log.Printf("Invalid video duration: %v\n", err)
		return nil, err
	}

	for i := 0; i < amount; i++ {
		// Generate a timestamp between 1% and 95% of the video
		randomSeconds := 0.01*duration + rand.Float64()*(0.94*duration)
		hours := int(randomSeconds / 3600)
		minutes := int((randomSeconds - float64(hours*3600)) / 60)
		seconds := int(randomSeconds) % 60
		milliseconds := int((randomSeconds - float64(int(randomSeconds))) * 1000)

		timestamp := fmt.Sprintf("%02d:%02d:%02d.%03d", hours, minutes, seconds, milliseconds)
		formattedTimestamp := fmt.Sprintf("%02d-%02d-%02d_%03d", hours, minutes, seconds, milliseconds)

		outputFile := filepath.Join(framepath, fmt.Sprintf("frame_%d_%s.jpg", i, formattedTimestamp))
		filepathofframes[i] = outputFile

		err = ffmpeg.Input(videoFilepath, ffmpeg.KwArgs{"ss": timestamp}).
			Output(outputFile, ffmpeg.KwArgs{"vframes": "1", "update": "1"}).Run()

		if err != nil {
			log.Printf("Error extracting frame %d at %s: %v\n", i, timestamp, err)
			return nil, fmt.Errorf("ffmpeg error extracting frame %d: %w", i, err)
		}

		log.Printf("Extracted frame %d/%d at %s", i+1, amount, timestamp)
	}

	return filepathofframes, nil
}

// ProgressLogger implements uploader.Progress to track upload progress
type ProgressLogger struct {
	startTime time.Time
	fileName  string // Add filename for clearer logs
}

// NewProgressLogger creates a new ProgressLogger instance.
func NewProgressLogger(fileName string) *ProgressLogger {
	return &ProgressLogger{fileName: fileName}
}

// Chunk logs upload progress including estimated remaining time
func (p *ProgressLogger) Chunk(ctx context.Context, state uploader.ProgressState) error {
	// Initialize start time on the first chunk
	if p.startTime.IsZero() {
		p.startTime = time.Now()
	}

	uploadedBytes := state.Uploaded
	totalBytes := state.Total

	// Avoid division by zero if totalBytes is 0 (though unlikely for uploads)
	if totalBytes == 0 {
		log.Printf("Uploading %s: Progress: 0.00%% (0.00/0.00 MB)", p.fileName)
		return nil
	}

	uploadedMB := float64(uploadedBytes) / (1024 * 1024)
	totalMB := float64(totalBytes) / (1024 * 1024)
	percentage := float64(uploadedBytes) * 100 / float64(totalBytes)

	// Calculate elapsed time and speed
	elapsed := time.Since(p.startTime)
	speedBytesPerSec := float64(uploadedBytes) / elapsed.Seconds() // Bytes per second
	speedMBPerSec := speedBytesPerSec / (1024 * 1024)              // MB per second

	// Calculate remaining time
	remainingBytes := totalBytes - uploadedBytes
	remainingTime := time.Duration(0)
	if speedBytesPerSec > 0 { // Avoid division by zero if speed is 0
		remainingSeconds := float64(remainingBytes) / speedBytesPerSec
		remainingTime = time.Duration(remainingSeconds * float64(time.Second))
	}

	// Format remaining time for better readability (e.g., 1m30s)
	remainingTimeStr := formatDuration(remainingTime)

	// Use carriage return '\r' to update the line in the terminal
	fmt.Printf("\rUploading %s: %.2f%% (%.2f/%.2f MB) | Speed: %.2f MB/s | ETA: %s ",
		p.fileName, percentage, uploadedMB, totalMB, speedMBPerSec, remainingTimeStr)

	// Print a newline when upload is complete
	if uploadedBytes == totalBytes {
		fmt.Println() // Move to the next line after completion
	}

	return nil
}

// formatDuration formats time.Duration into a more readable string like "1m30s"
func formatDuration(d time.Duration) string {
	d = d.Round(time.Second)
	m := d / time.Minute
	d -= m * time.Minute
	s := d / time.Second
	if m > 0 {
		return fmt.Sprintf("%dm%ds", m, s)
	}
	return fmt.Sprintf("%ds", s)
}
