package main

import "github.com/bwmarrin/discordgo"
func err_msg(s *discordgo.Session, channelID, msg string){
	s.ChannelMessageSend(channelID, msg)
}
func queue_msg(s *discordgo.Session, channelID, song_name, autrhor string){
	msg := "Add to queue " + song_name + ". Added By @" + autrhor
	s.ChannelMessageSend(channelID, msg)
}
func playing_msg(s *discordgo.Session, channelID, song_name, autrhor string){
	msg := "Playing " + song_name + ". Added By @" + autrhor
	s.ChannelMessageSend(channelID, msg)
}
func pause_msg(s *discordgo.Session, channelID string){
	s.ChannelMessageSend(channelID, "Pause_stream")
}
func stop_msg(s *discordgo.Session, channelID string){
	s.ChannelMessageSend(channelID, "Stream_stopped")
}