package main

import (
	"encoding/binary"
	"github.com/jonas747/dca"
	"github.com/rs/zerolog/log"
	"github.com/rylio/ytdl"
	"net/http"

	"context"
	"fmt"
	"github.com/bwmarrin/discordgo"
)

func stop_stream(s *discordgo.Session, m *discordgo.MessageCreate){
	channel, err := s.State.Channel(m.ChannelID)
	if err != nil {
		fmt.Println(err)
	}
	guild, _ := s.State.Guild(channel.GuildID)
	for _, vs := range guild.VoiceStates {
		if vs.UserID == m.Author.ID {
			servers_vc[vs.GuildID].music_queue[0].stop = true
			stop_msg(s, m.ChannelID)
		}
	}
}
func pause_stream(s *discordgo.Session, m *discordgo.MessageCreate){
	channel, err := s.State.Channel(m.ChannelID)
	if err != nil {
		fmt.Println(err)
	}
	guild, _ := s.State.Guild(channel.GuildID)
	for _, vs := range guild.VoiceStates {
		if vs.UserID == m.Author.ID {
			servers_vc[vs.GuildID].music_queue[0].pause <- true
			pause_msg(s, m.ChannelID)
		}
	}
}
func straem_to_discord(s *discordgo.Session, m *discordgo.MessageCreate, url string) {
	channel, err := s.State.Channel(m.ChannelID)
	if err != nil {
		fmt.Println(err)
	}
	guild, _ := s.State.Guild(channel.GuildID)
	for _, vs := range guild.VoiceStates {
		if vs.UserID == m.Author.ID {
			mi := new(music_info)
			url, name := get_url(url)
			mi.songname = name
			mi.url = url
			mi.guildID = vs.GuildID
			mi.Autorname= m.Author.Username
			mi.channelID = vs.ChannelID
			mi.s = s
			mi.text_channelId = m.ChannelID
			mi.pause=make(chan bool ,1)
			queue_chan <- mi
		}
	}
}
func queue_way(){
	for {
		mi := <-queue_chan
		if mi==nil{
			continue
		}
		if servers_vc[mi.guildID]==nil{
			vc := get_vc_connection(mi.s, mi.guildID, mi.channelID)
			servers_vc[mi.guildID] = new(server_queue)
			servers_vc[mi.guildID].vc = vc
			servers_vc[mi.guildID].music_queue = append(servers_vc[mi.guildID].music_queue, mi)
			go play_on_server(mi.guildID)
		} else {
			queue_msg(servers_vc[mi.guildID].music_queue[0].s, mi.text_channelId, mi.songname,mi.Autorname)
			servers_vc[mi.guildID].music_queue = append(servers_vc[mi.guildID].music_queue, mi)
		}
	}
}
func play_on_server(guildID string){
	for {
		if len(servers_vc[guildID].music_queue) == 0{
			break
		}
		playing_msg(servers_vc[guildID].music_queue[0].s, servers_vc[guildID].music_queue[0].text_channelId, servers_vc[guildID].music_queue[0].songname,servers_vc[guildID].music_queue[0].Autorname)
		stop_status := play(servers_vc[guildID].music_queue[0].guildID,servers_vc[guildID].music_queue[0].channelID, servers_vc[guildID].music_queue[0].url, servers_vc[guildID].vc)
		if stop_status == true{
			break
		}
		servers_vc[guildID].music_queue=servers_vc[guildID].music_queue[1:]
	}
	servers_vc[guildID].vc.Disconnect()
	delete(servers_vc, guildID)
}
func get_url(videoURL string) (string, string) {
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
	return DownloadURL.String(), vid.Title
}
func play(guildID, channelID, url string, vc *discordgo.VoiceConnection) bool {
	options := dca.StdEncodeOptions
	options.RawOutput = true
	options.Bitrate = 96
	options.Application = "lowdelay"
	dca.EncodeFile(url, options)

	encode, _ := dca.EncodeFile(url, options)
	defer encode.Cleanup()
	for {
		select {
		case <-servers_vc[guildID].music_queue[0].pause:
			<-servers_vc[guildID].music_queue[0].pause
		default:
			if (servers_vc[guildID].music_queue[0].stop==true){
				return true
			}
			var sz_frame int16
			err := binary.Read(encode, binary.LittleEndian, &sz_frame)
			if err != nil {
				fmt.Println(err)
				return false
			}
			Inbuf := make([]byte, sz_frame)
			_ = binary.Read(encode, binary.LittleEndian, &Inbuf)
			vc.OpusSend <- Inbuf
		}
	}
	return true
}
func get_vc_connection(s *discordgo.Session, guildID, channelID string) (*discordgo.VoiceConnection){
	vc, _ := s.ChannelVoiceJoin(guildID, channelID, false, false)
	return vc
}