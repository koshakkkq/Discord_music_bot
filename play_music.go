package main

import (
	"github.com/rs/zerolog/log"
	"layeh.com/gopus"
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
func straem_to_discord(s *discordgo.Session, m *discordgo.MessageCreate, url string){
	channel, err := s.State.Channel(m.ChannelID)
	if err!=nil{
		fmt.Println(err)
	}
	guild,_ := s.State.Guild(channel.GuildID)
	for _, vs := range guild.VoiceStates{
		if (vs.UserID == m.Author.ID){
			URL_deocde(s, guild.ID, vs.ChannelID, url)
		}
	}
}
func get_url (videoURL string) (string){
	client := ytdl.Client{
		HTTPClient: http.DefaultClient,
		Logger:     log.Logger,
	}
	c :=  context.Background()
	vid, err := client.GetVideoInfo(c, videoURL)
	if err != nil {
		fmt.Println("Failed to get video info")
	}

	format := vid.Formats.Extremes(ytdl.FormatAudioBitrateKey, true)[0]
	DownloadURL, _ := client.GetDownloadURL(c, vid, format)
	return DownloadURL.String()
}
func URL_deocde(s *discordgo.Session, guildID,channelID, videoURL string){

}
func SendPCM(vc *discordgo.VoiceConnection, pcm <-chan []int16){
	if pcm==nil {
		return
	}
	opusEncoder, _ := gopus.NewEncoder(frameRate, channels, gopus.Audio)

	for {
		recv, ok := <-pcm
		if !ok {
			fmt.Println("PCM Channel closed", nil)
			return
		}
		opus, err := opusEncoder.Encode(recv, frameSize, maxBytes)
		if err != nil {
			fmt.Println("Encoding Error", err)
			return
		}
		vc.OpusSend <- opus
	}
}