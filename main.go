package main

import (
	"strings"
	//"io/ioutil"
	"os/signal"
	"syscall"
	//"github.com/jonas747/dca"
	//"encoding/binary"
	//"github.com/davecgh/go-spew/spew"
	//"net/http"
	"fmt"
	"github.com/bwmarrin/discordgo"
	"os"
)
var (
	Token string  = "NzA5MDc0NDI4NjkyMDA0OTE1.Xrg1WA.Jut-QCgs6fAbFcTQ4uU3wgdL5GE"
)
var stopChannel chan bool
func main(){
	dg, err := discordgo.New("Bot "+ Token)
	if err!=nil {
		fmt.Println("error !gg!", err)
		return
	}
	stopChannel = make(chan bool)
	dg.AddHandler(messageCreate)
	err = dg.Open()
	if err!=nil {
		fmt.Print("open Web_socket Error gg!!", err)
		return
	}
	fmt.Println("Bot is now running.  Press CTRL-C to exit.")
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM)
	<-sc
	fmt.Println("asdasd")
	dg.Close()
}
func music_load(){



}
func messageCreate(s *discordgo.Session, m *discordgo.MessageCreate){
	if (m.Author.ID==s.State.User.ID) {
		return
	}
	if (m.Content[0]=='!'){
		msg_string := string(m.Content)
		bot_func := strings.Split(msg_string[1:len(msg_string)]," ")
		switch bot_func[0] {
		case "play":
			straem_to_discord(s, m, bot_func[1])
		}
	}
	if (m.Content=="Снюс"){
		s.ChannelMessageSend(m.ChannelID, "Чесвин")
	}
	if (m.Content=="Чесвин"){
		s.ChannelMessageSend(m.ChannelID, "Снюс")
	}
}
