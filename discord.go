package main

import (
	"crypto/sha1"
	"encoding/hex"
	"encoding/json"
	"log"
	"os"

	"github.com/bwmarrin/discordgo"
)

var s *discordgo.Session

type Configuration struct {
	Discord_bot_auth      string
	Dc_createsync_prompts string
}

type AtSync_info struct {
	dc_channel_id   string
	dc_channel_info string
	dc_creator_name string
	wx_group_name   string
	hashString      string
}

var AtSync_msg AtSync_info

func init_dc() {
	// 打开文件
	config_file, err := os.Open("config.json")
	if err != nil {
		log.Println("Failed to open config file : ", err)
		return
	}

	// 关闭文件
	defer config_file.Close()

	//NewDecoder创建一个从file读取并解码json对象的*Decoder，解码器有自己的缓冲，并可能超前读取部分json数据。
	decoder := json.NewDecoder(config_file)

	conf := Configuration{}
	//Decode从输入流读取下一个json编码值并保存在v指向的值里
	decoder.Decode(&conf)

	s, err = discordgo.New("Bot " + conf.Discord_bot_auth)
	if err != nil {
		log.Println("Failed to create Discord session: ", err)
		return
	}

	// 注册斜杆命令
	var (
		commands = []*discordgo.ApplicationCommand{
			{
				Name:        "createsync",
				Description: "创建微信到当前频道的同步",
			},
		}
		commandHandlers = map[string]func(s *discordgo.Session, i *discordgo.InteractionCreate){
			"createsync": func(s *discordgo.Session, i *discordgo.InteractionCreate) {
				// 取当前频道ID的Sha1值[:6]
				hasher := sha1.New()
				hasher.Write([]byte(i.ChannelID))
				hash := hasher.Sum(nil)
				hashString := hex.EncodeToString(hash)[:6]

				AtSync_msg.dc_channel_id = i.ChannelID
				AtSync_msg.dc_channel_info = "https://discord.com/channels/" + i.GuildID + "/" + i.ChannelID
				AtSync_msg.dc_creator_name = ""
				AtSync_msg.wx_group_name = ""
				AtSync_msg.hashString = hashString

				log.Println("dc_channel_id dc_channel_info", AtSync_msg.dc_channel_id, AtSync_msg.dc_channel_info)

				s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
					Type: discordgo.InteractionResponseChannelMessageWithSource,
					Data: &discordgo.InteractionResponseData{
						Flags:   discordgo.MessageFlagsEphemeral,
						Content: conf.Dc_createsync_prompts + hashString,
					},
				})
			},
		}
	)

	s.AddHandler(func(s *discordgo.Session, i *discordgo.InteractionCreate) {
		if h, ok := commandHandlers[i.ApplicationCommandData().Name]; ok {
			h(s, i)
		}
	})

	// 建立连接
	err = s.Open()
	if err != nil {
		log.Println("Failed to open connection: ", err)
	}

	log.Println("注册斜杆命令...")
	registeredCommands := make([]*discordgo.ApplicationCommand, len(commands))
	for i, v := range commands {
		cmd, err := s.ApplicationCommandCreate(s.State.User.ID, "", v)
		if err != nil {
			log.Panicf("Cannot create '%v' command: %v", v.Name, err)
		}
		registeredCommands[i] = cmd
	}

}

func discord_connection_check() {
	// 检查discord连接是否正常，否则重新连接
	if s == nil || !s.DataReady {
		init_dc()
	}
}

func discord_send_text(content string, dc_channel_id string) {

	// 发送消息
	_, err := s.ChannelMessageSend(dc_channel_id, content)
	if err != nil {
		log.Println("Error sending text: ", err)
	}
}

func discord_send_file(content string, name string, path string, dc_channel_id string) {

	file, err := os.Open(path)
	if err != nil {
		log.Println("Error opening image file: ", err)
		return
	}
	// 发送消息
	_, err = s.ChannelFileSendWithMessage(dc_channel_id, content, name, file)
	if err != nil {
		log.Println("Error sending image: ", err)
	}

}
