package bot

import (
	"bufio"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/signal"
	"regexp"
	"strconv"
	"strings"
	"syscall"
	"time"

	"vrc_bot/vrcapi"
	"github.com/bwmarrin/discordgo"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

func Init() {
	//dg, err := discordgo.New("Bot " + viper.GetString("discordbot.token"))
	// Create a new Discord session using email+password.
	dg, err := discordgo.New(viper.GetString("discordbot.email"), viper.GetString("discordbot.password"))

	if err != nil {
		log.Error("error creating Discord session,", err)
		return
	}
	// Register the messageCreate func as a callback for MessageCreate events.
	dg.AddHandler(messageCreate)
	// Open a websocket connection to Discord and begin listening.
	err = dg.Open()
	if err != nil {
		log.Error("error opening connection,", err)
		return
	}
	// Wait here until CTRL-C or other term signal is received.
	log.Info("Bot is now running.  Press CTRL-C to exit.")
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-sc
	// Cleanly close down the Discord session.
	dg.Close()
}

// This function will be called (due to AddHandler above) every time a new
// message is created on any channel that the autenticated bot has access to.
func messageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {
	// Ignore all messages created by the bot itself
	// This isn't required in this specific example but it's a good practice.
	if m.Author.ID == s.State.User.ID {
		return
	}

	log.Info(fmt.Sprintf("%s %s messaged the bot", m.Author.Username, m.Author.ID))

	isDm, _ := ComesFromDM(s, m)

	if isDm {
		//vrc
		if strings.Contains(m.Content, "!export") {
			s.ChannelMessageSend(m.ChannelID, "Processing...")
			var re = regexp.MustCompile(`(?m)\"(.*?)\"`)
			var str = m.Content
			match := re.FindAllString(str, -1)
			fmt.Println(len(match))
			if len(match) > 1 {
				username := strings.Replace(match[0], "\"", "", -1)
				password := strings.Replace(match[1], "\"", "", -1)
				friendIds := vrcapi.ExportFriends(username, password)
				s.ChannelMessageSend(m.ChannelID, "Parsing Friends...")
				path := fmt.Sprintf("friends_%s_%s.txt", m.Author.ID, strconv.FormatInt(time.Now().UTC().UnixNano(), 10))
				f, err := os.Create("temp/" + path)
				if err != nil {
					fmt.Println(err)
					f.Close()
					return
				}

				for _, v := range friendIds {
					fmt.Fprintln(f, v)
					if err != nil {
						fmt.Println(err)
						return
					}
				}
				err = f.Close()
				if err != nil {
					fmt.Println(err)
					return
				}
				r, _ := os.Open("temp/" + path)
				s.ChannelFileSendWithMessage(m.ChannelID, "Backup this file to import your friends", path, r)
			} else {
				s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("@%s, please send your details in this format : !export \"<username>\" \"<password>\"", m.Author.Username))
			}
		}

		if strings.Contains(m.Content, "!import") && len(m.Attachments) > 0 {
			s.ChannelMessageSend(m.ChannelID, "Processing...")
			var re = regexp.MustCompile(`(?m)\"(.*?)\"`)
			var str = m.Content
			match := re.FindAllString(str, -1)
			if len(match) > 1 {
				username := strings.Replace(match[0], "\"", "", -1)
				password := strings.Replace(match[1], "\"", "", -1)

				s.ChannelMessageSend(m.ChannelID, "Parsing file: "+m.Attachments[0].Filename)
				tempFileName := fmt.Sprintf("temp/import_%s_%s.txt", m.Author.ID, strconv.FormatInt(time.Now().UTC().UnixNano(), 10))
				if err := DownloadFile(tempFileName, m.Attachments[0].URL); err != nil {
					panic(err)
				}

				file, err := os.Open(tempFileName)
				if err != nil {
					log.Fatal(err)
				}
				defer file.Close()

				scanner := bufio.NewScanner(file)
				for scanner.Scan() {
					vrcapi.SendFriendRequest(username, password, scanner.Text())
				}

				if err := scanner.Err(); err != nil {
					log.Fatal(err)
				}
				s.ChannelMessageSend(m.ChannelID, "Added all usr_id's")
			} else {
				s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("@%s, please attach the exported list and send your details in this format : !import \"<username>\" \"<password>\"", m.Author.Username))
			}
		}
	}
}

// ComesFromDM returns true if a message comes from a DM channel
func ComesFromDM(s *discordgo.Session, m *discordgo.MessageCreate) (bool, error) {
	channel, err := s.State.Channel(m.ChannelID)
	if err != nil {
		if channel, err = s.Channel(m.ChannelID); err != nil {
			return false, err
		}
	}

	return channel.Type == discordgo.ChannelTypeDM, nil
}

// DownloadFile will download a url to a local file. It's efficient because it will
// write as it downloads and not load the whole file into memory.
func DownloadFile(filepath string, url string) error {

	// Get the data
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// Create the file
	out, err := os.Create(filepath)
	if err != nil {
		return err
	}
	defer out.Close()

	// Write the body to file
	_, err = io.Copy(out, resp.Body)
	return err
}
