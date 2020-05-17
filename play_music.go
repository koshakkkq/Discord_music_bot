package main

import (
	"encoding/binary"
	"github.com/jonas747/dca"
	"github.com/rs/zerolog/log"
	"net/http"

	//"github.com/jonas747/dca"
	//"encoding/binary"
	//"github.com/davecgh/go-spew/spew"
	"github.com/rylio/ytdl"
	//"net/http"
	"context"
	"fmt"
	"github.com/bwmarrin/discordgo"
)

func straem_to_discord(s *discordgo.Session, m *discordgo.MessageCreate, url string) {
	channel, err := s.State.Channel(m.ChannelID)
	if err != nil {
		fmt.Println(err)
	}
	guild, _ := s.State.Guild(channel.GuildID)
	for _, vs := range guild.VoiceStates {
		if vs.UserID == m.Author.ID {
			URL_deocde(s, guild.ID, vs.ChannelID, url)
		}
	}
}
func get_url(videoURL string) string {
	client := ytdl.Client{
		HTTPClient: http.DefaultClient,
		Logger:     log.Logger,
	}
	c := context.Background()
	vid, err := client.GetVideoInfo(c, videoURL)
	if err != nil {
		fmt.Println("Failed to get video info")
	}

	format := vid.Formats.Extremes(ytdl.FormatAudioBitrateKey, true)[0]
	DownloadURL, _ := client.GetDownloadURL(c, vid, format)
	return DownloadURL.String()
}
func URL_deocde(s *discordgo.Session, guildID, channelID, videoURL string) {
	options := dca.StdEncodeOptions
	options.RawOutput = true
	options.Bitrate = 96
	options.Application = "lowdelay"
	url := get_url(videoURL)
	dca.EncodeFile(url, options)
	vc, _ := s.ChannelVoiceJoin(guildID, channelID, false, false)
	defer vc.Disconnect()
	encode, _ := dca.EncodeFile(url, options)
	buffer := make([][]byte, 0)
	defer encode.Cleanup()
	for {
		var sz_frame int16
		err := binary.Read(encode, binary.LittleEndian, &sz_frame)
		if err != nil {
			break
		}
		Inbuf := make([]byte, sz_frame)
		_ = binary.Read(encode, binary.LittleEndian, &Inbuf)
		buffer = append(buffer, Inbuf)
	}
	SendPCM(vc, buffer)
}
func SendPCM(vc *discordgo.VoiceConnection, buffer [][]byte) {
	for _, buff := range buffer {
		vc.OpusSend <- buff
	}
	return
}
