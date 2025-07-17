package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/gotd/td/session"
	"github.com/gotd/td/telegram"
	"github.com/gotd/td/telegram/message"
	"github.com/gotd/td/telegram/message/styling"
	"github.com/gotd/td/telegram/uploader"
	"github.com/gotd/td/tg"
)

type uploaderengin struct {
	appid        int
	apphash      string
	downloadpath string
	sessionpath  string
	tgclient     *telegram.Client
	tgchannel    string
	ctx          context.Context
	cancel       context.CancelFunc // Add a CancelFunc to control the context
}

func NewUploader(appid int, apphash string, downloadpath, sessionpath, tgchannel string) *uploaderengin {
	sessionStorage := &session.FileStorage{Path: sessionpath}

	// Create a context that can be canceled to control the client's lifetime
	ctx, cancel := context.WithCancel(context.Background())

	tgclient := telegram.NewClient(appid, apphash, telegram.Options{
		SessionStorage: sessionStorage,
	})

	return &uploaderengin{
		appid:        appid,
		apphash:      apphash,
		downloadpath: downloadpath,
		sessionpath:  sessionpath,
		tgclient:     tgclient,
		tgchannel:    tgchannel,
		ctx:          ctx,    // Use the new context
		cancel:       cancel, // Store the cancel function
	}
}

func (up *uploaderengin) Connect() error {
	log.Println("Connecting to Telegram...")
	return up.tgclient.Run(up.ctx, func(ctx context.Context) error {
		log.Println("Telegram client connected.")
		// This function will keep the client running until up.ctx is canceled.
		<-ctx.Done()
		log.Println("Telegram client disconnected.")
		return up.ctx.Err() // Return the context error if canceled
	})
}

func (up *uploaderengin) Disconnect() {
	log.Println("Disconnecting from Telegram...")
	up.cancel() // Cancel the context to stop the client
}

func (up *uploaderengin) Uploadwithffmpeg(rpath string) error {
	npath := cfg.DownloadPath + "/" + rpath
	log.Printf("Starting upload process for: %s", npath)

	if _, err := os.Stat(npath); os.IsNotExist(err) {
		return fmt.Errorf("video file does not exist: %s", npath)
	}

	length, width, height, err := VideoInfo(npath)
	if err != nil {
		return fmt.Errorf("getting video info for %q: %w", npath, err)
	}

	videoDir := filepath.Dir(npath)
	baseName := filepath.Base(npath)
	frameDir := filepath.Join(videoDir, "frames_"+strings.TrimSuffix(baseName, filepath.Ext(baseName)))

	if err := os.MkdirAll(frameDir, 0755); err != nil {
		return fmt.Errorf("creating frame directory %q: %w", frameDir, err)
	}

	ss, err := ExtractRandomFramesWithDir(npath, length, 4, frameDir)
	if err != nil {
		return fmt.Errorf("extracting frames: %w", err)
	}
	if len(ss) == 0 {
		return fmt.Errorf("no frames extracted")
	}
	defer os.RemoveAll(frameDir)

	heightInt, _ := strconv.Atoi(height)
	widthInt, _ := strconv.Atoi(width)
	lengthFloat, _ := strconv.ParseFloat(strings.TrimSpace(length), 64)
	lengthDur := time.Duration(lengthFloat * float64(time.Second))
	// fileName := strings.TrimSuffix(baseName, filepath.Ext(baseName))

	api := tg.NewClient(up.tgclient)
	sender := message.NewSender(api)
	target := sender.Resolve(up.tgchannel)

	// Upload thumbnail
	thumbnailPath := ss[0]
	thumbFile, err := os.Open(thumbnailPath)
	if err != nil {
		return fmt.Errorf("open thumb: %w", err)
	}
	defer thumbFile.Close()
	stat, _ := thumbFile.Stat()
	thumb, err := uploader.NewUploader(api).
		WithProgress(NewProgressLogger("thumbnail")).
		Upload(up.ctx, uploader.NewUpload(filepath.Base(thumbnailPath), thumbFile, stat.Size()))
	if err != nil {
		return fmt.Errorf("upload thumb: %w", err)
	}

	// Upload video
	videoUploader := uploader.NewUploader(api)
	videoUpload, err := videoUploader.
		WithProgress(NewProgressLogger(npath)).
		WithPartSize(512*1024).
		FromPath(up.ctx, npath)
	if err != nil {
		return fmt.Errorf("upload video: %w", err)
	}

	doc := message.UploadedDocument(videoUpload, styling.Italic(rpath)).
		MIME("video/mp4").
		Thumb(thumb).
		Video().
		SupportsStreaming().
		Duration(lengthDur).
		Resolution(widthInt, heightInt)

	if _, err := target.Media(up.ctx, doc); err != nil {
		return fmt.Errorf("send video: %w", err)
	}

	// Upload screenshots as album
	if len(ss) > 0 {
		var inputMedias []message.MultiMediaOption
		for _, img := range ss {
			_, err := os.Stat(img)
			if err != nil {
				continue
			}
			photoUpload, err := uploader.NewUploader(api).
				WithProgress(NewProgressLogger(filepath.Base(img))).
				FromPath(up.ctx, img)
			if err != nil {
				log.Printf("warn: failed to upload %s: %v", img, err)
				continue
			}
			inputMedias = append(inputMedias, message.UploadedPhoto(photoUpload))
		}
		if len(inputMedias) > 1 {
			_, err = target.Album(up.ctx, inputMedias[0], inputMedias[1:]...)
		} else if len(inputMedias) == 1 {
			_, err = target.Media(up.ctx, inputMedias[0])
		}
		if err != nil {
			return fmt.Errorf("send album: %w", err)
		}
	}

	// Delete video
	// check config
	if cfg.DeleteAfterUpload {
		log.Printf("Deleting video file after upload: %s", npath)
		if err := os.Remove(npath); err != nil {
			log.Printf("warn: failed to delete video %s: %v", npath, err)
		}

	} else {
		log.Printf("Skipping deletion of video file after upload: %s", npath)
		return nil
	}
	return nil
	// log.Printf("Deleting video file after upload: %s", npath)
	// if err := os.Remove(npath); err != nil {
	// 	log.Printf("warn: failed to delete video %s: %v", npath, err)
	// }
	// return nil
}
func (up *uploaderengin) Upload(rpath string) error {
	npath := cfg.DownloadPath + "/" + rpath

	if _, err := os.Stat(npath); os.IsNotExist(err) {
		return fmt.Errorf("video file does not exist: %s", npath)
	}

	// Check if ffmpeg is installed
	if _, err := exec.LookPath("ffmpeg"); err == nil {
		return up.Uploadwithffmpeg(rpath)
	}

	// If ffmpeg is not installed, use the alternative upload method
	return up.UploadWithoffmpeg(rpath)
}

// if ffmpeg not installed, use this function
func (up *uploaderengin) UploadWithoffmpeg(rpath string) error {
	npath := cfg.DownloadPath + "/" + rpath

	if _, err := os.Stat(npath); os.IsNotExist(err) {
		return fmt.Errorf("video file does not exist: %s", npath)
	}
	api := tg.NewClient(up.tgclient)
	sender := message.NewSender(api)
	target := sender.Resolve(up.tgchannel)

	// Upload video
	videoUploader := uploader.NewUploader(api)
	videoUpload, err := videoUploader.
		WithProgress(NewProgressLogger(npath)).
		WithPartSize(512*1024).
		FromPath(up.ctx, npath)
	if err != nil {
		return fmt.Errorf("upload video: %w", err)
	}

	doc := message.UploadedDocument(videoUpload, styling.Italic(rpath)).
		MIME("video/mp4").
		Video().
		SupportsStreaming()

	if _, err := target.Media(up.ctx, doc); err != nil {
		return fmt.Errorf("send video: %w", err)
	}
	// Delete video
	// check config
	if cfg.DeleteAfterUpload {
		log.Printf("Deleting video file after upload: %s", npath)
		if err := os.Remove(npath); err != nil {
			log.Printf("warn: failed to delete video %s: %v", npath, err)
		}

	} else {
		log.Printf("Skipping deletion of video file after upload: %s", npath)
		return nil
	}
	return nil
}

// ExtractRandomFramesWithDir is a modified version that accepts a target directory
func ExtractRandomFramesWithDir(videoPath, duration string, count int, frameDir string) ([]string, error) {
	// This function should be implemented to extract frames to the specified directory
	// For now, assuming the original ExtractRandomFrames function exists and modifying it
	return ExtractRandomFrames(videoPath, duration, count)
}
