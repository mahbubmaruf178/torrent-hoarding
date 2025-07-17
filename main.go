package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/anacrolix/torrent"
)

var cfg Config

func main() {
	data, err := os.ReadFile("config.json")
	if err != nil {
		panic(err)
	}

	if err := json.Unmarshal(data, &cfg); err != nil {
		panic(err)
	}

	tg, found := GetTelegramByName(&cfg, "telegram")
	if !found {
		fmt.Println("Telegram config not found")
		return
	}

	uploader := NewUploader(tg.APIID, tg.APIHash, tg.DownloadPath, tg.SessionName, tg.ChannelID)

	go func() {
		if err := uploader.Connect(); err != nil {
			log.Fatalf("Error connecting to Telegram: %v", err)
		}
		fmt.Println("Connected to Telegram successfully!")
	}()
	defer func() {
		uploader.Disconnect()
		fmt.Println("Disconnected from Telegram.")
	}()

	clientConfig := torrent.NewDefaultClientConfig()
	clientConfig.Seed = false
	clientConfig.DataDir = tg.DownloadPath

	client, err := torrent.NewClient(clientConfig)
	if err != nil {

		log.Fatalf("Error creating torrent client: %v", err)
	}
	defer client.Close()

	for {
		uploadloop(client, uploader)
	}

}
func uploadloop(client *torrent.Client, uploader *uploaderengin) {
	var megnet string

	fmt.Println("input your manget url :- ")

	fmt.Scanln(&megnet)

	t, err := client.AddMagnet(megnet)
	if err != nil {
		log.Fatalf("Error adding magnet: %v", err)
	}

	fmt.Println("Fetching torrent metadata...")
	<-t.GotInfo()
	fmt.Println("Metadata received!")

	t.DownloadAll()

	fmt.Println("\n=== Torrent Information ===")
	fmt.Printf("Name: %s\n", t.Name())
	fmt.Printf("Size: %.2f MB\n", float64(t.Length())/(1024*1024))
	fmt.Printf("Files: %d\n", len(t.Files()))

	for _, file := range t.Files() {

		file.SetPriority(torrent.PiecePriorityHigh)
		downloadFile(*file)

		if err := uploader.Upload(file.Path()); err != nil {
			log.Printf("ERROR: Upload failed for %s: %v", file.Path(), err)
		}

	}
}

type DownloadProgress struct {

	// For calculating speed
	BytesDownloadedAtLastUpdate int64
	LastProgress                float64
	LastUpdatedTime             time.Time
}

func downloadFile(file torrent.File) {
	fullsize := file.Length()
	fmt.Printf("Downloading %s: %.2f MB\n", file.Path(), float64(fullsize)/(1024*1024))
	stime := time.Now()
	dp := DownloadProgress{LastUpdatedTime: stime, BytesDownloadedAtLastUpdate: 0}

	for {

		progress := 100 * (float64(file.BytesCompleted()) / float64(file.Length()))
		now := time.Now()
		if now.Sub(dp.LastUpdatedTime) >= time.Second {
			// Calculate download speed in bytes per second
			downloadedSinceLastUpdate := file.BytesCompleted() - dp.BytesDownloadedAtLastUpdate
			timeSinceLastUpdate := now.Sub(dp.LastUpdatedTime).Seconds()
			speed := float64(downloadedSinceLastUpdate) / timeSinceLastUpdate

			fmt.Printf("\rDownloading %s: %.2f%% - full size %.2f MB -- downloaded %.2f MB , speed %.2f KB/s, elap %s \n", file.DisplayPath(), progress, float64(fullsize)/(1024*1024), float64(file.BytesCompleted())/(1024*1024), speed/1024, time.Since(stime).Round(time.Second))
			dp.LastUpdatedTime = now
			dp.BytesDownloadedAtLastUpdate = file.BytesCompleted()
		}

		if file.BytesCompleted() == fullsize {
			break
		}
	}
}
