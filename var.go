package main

import (
	//"gopkg.in/oleiade/lane.v1"
	"github.com/bwmarrin/discordgo"
)
const (
	channels  int = 2                   // 1 for mono, 2 for stereo
	frameRate int = 48000               // audio sampling rate
	frameSize int = 960                 // uint16 size of each audio frame
	maxBytes  int = (frameSize * 2) * 2 // max size of opus data
)
type server_queue struct {
	music_queue []*music_info
	vc *discordgo.VoiceConnection
}
type config struct {
	Token string
	PREFIX string
	YT_KEY string
 }
type id_video struct {
	kind string
	videoId string
}
type music_info struct {
	songname string
	Autorname string
	url string
	text_channelId string
	guildID string
	channelID string
	pause chan bool
	stop bool
	s *discordgo.Session
}
var (
	servers_vc = map[string]*server_queue{}
	queue_block= map[string]bool{}
	queue_chan = make(chan *music_info, 101)
)
