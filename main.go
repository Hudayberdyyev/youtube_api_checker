package main

import (
	"fmt"
	"github.com/kkdai/youtube/v2"
	"io"
	"log"
	"os"
	"runtime"
	"strings"
	"sync/atomic"
)

const MAXN = int(10)

func main() {
	log.SetFlags(log.LstdFlags | log.Llongfile)
	bufferChan := make(chan struct{}, runtime.GOMAXPROCS(0))
	client := youtube.Client{}
	var errCount atomic.Int32
	for i := 0; i < MAXN; i++ {
		bufferChan <- struct{}{}
		go func(iteration int) {
			videoID := "BaW_jenozKc"
			video, err := client.GetVideo(videoID)
			if err != nil {
				errCount.Add(1)
				log.Printf("%d) error when retrieving video information from API: %v\n", iteration, err)
			} else {
				log.Printf("%d) successful retrieve video information from API: %s\n", iteration, video.ID)
			}

			<-bufferChan
		}(i)
		runtime.Gosched()
	}

	for len(bufferChan) > 0 {
	}

	log.Printf("error count = %d\n", errCount.Load())
}

func ExampleClient() {
	videoID := "BaW_jenozKc"
	client := youtube.Client{}

	video, err := client.GetVideo(videoID)
	if err != nil {
		panic(err)
	}

	formats := video.Formats.WithAudioChannels() // only get videos with audio
	stream, _, err := client.GetStream(video, &formats[0])
	if err != nil {
		panic(err)
	}

	file, err := os.Create("video.mp4")
	if err != nil {
		panic(err)
	}
	defer file.Close()

	_, err = io.Copy(file, stream)
	if err != nil {
		panic(err)
	}
}

func ExamplePlaylist() {
	playlistID := "PLQZgI7en5XEgM0L1_ZcKmEzxW1sCOVZwP"
	client := youtube.Client{}

	playlist, err := client.GetPlaylist(playlistID)
	if err != nil {
		panic(err)
	}

	/* ----- Enumerating playlist videos ----- */
	header := fmt.Sprintf("Playlist %s by %s", playlist.Title, playlist.Author)
	println(header)
	println(strings.Repeat("=", len(header)) + "\n")

	for k, v := range playlist.Videos {
		fmt.Printf("(%d) %s - '%s'\n", k+1, v.Author, v.Title)
	}

	/* ----- Downloading the 1st video ----- */
	entry := playlist.Videos[0]
	video, err := client.VideoFromPlaylistEntry(entry)
	if err != nil {
		panic(err)
	}
	// Now it's fully loaded.

	fmt.Printf("Downloading %s by '%s'!\n", video.Title, video.Author)

	stream, _, err := client.GetStream(video, &video.Formats[0])
	if err != nil {
		panic(err)
	}

	file, err := os.Create("video.mp4")

	if err != nil {
		panic(err)
	}

	defer file.Close()
	_, err = io.Copy(file, stream)

	if err != nil {
		panic(err)
	}

	println("Downloaded /video.mp4")
}
