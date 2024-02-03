package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"regexp"
	"strings"
	"time"

	"github.com/eatmoreapple/openwechat"
)

func main() {
	// 初始化数据库连接
	init_db()
	// Discord 初始化
	init_dc()

	// 读取配置文件
	config_file, _ := os.Open("config.json")
	defer config_file.Close()
	decoder := json.NewDecoder(config_file)
	type Configuration struct {
		Wechat_group_name string
		Static_path       string
	}
	conf := Configuration{}
	decoder.Decode(&conf)

	// 初始化日志文件
	logfile_path := conf.Static_path + "logfile.log"
	logfile, err := os.Create(logfile_path)
	if err != nil {
		log.Fatal(err)
	}
	defer logfile.Close()
	log.SetOutput(logfile)

	// 桌面模式
	bot := openwechat.DefaultBot(openwechat.Desktop)
	reloadStorage := openwechat.NewFileHotReloadStorage("storage.json")
	defer reloadStorage.Close()
	bot.PushLogin(reloadStorage, openwechat.NewRetryLoginOption())

	self, err := bot.GetCurrentUser()
	if err != nil {
		log.Println(err)
		return
	}
	groups, err := self.Groups()
	fmt.Println(groups, err)

	bot.MessageHandler = func(msg *openwechat.Message) {
		fmt.Println(msg.Content)
		sender, _ := msg.SenderInGroup()
		sendgr, _ := msg.Sender()
		sender_content := msg.Content
		year, month, day := time.Now().Date()
		hour, min, sec := time.Now().Hour(), time.Now().Minute(), time.Now().Second()
		cur_time := fmt.Sprintf("%d-%02d-%02d-%02d-%02d-%02d", year, month, day, hour, min, sec)

		// 群聊天记录转录
		if msg.IsSendByGroup() && (msg.IsText() || msg.IsPicture()) {
			// 微信群名是否配置在数据库中
			DCNews_info, err := judge_dcnews_state(sendgr.NickName)
			// 如果不在，则...
			if err != nil {
				// 判断dc是否执行同步命令 和 发送内容是否为同步码
				if AtSync_msg.dc_channel_id != "" && sender_content == AtSync_msg.hashString {
					log.Println(sendgr.NickName, AtSync_msg.dc_channel_id, AtSync_msg.dc_channel_info)
					// 插入配置项
					insert_sync_dcCommandTarget(sendgr.NickName, AtSync_msg.dc_channel_id, AtSync_msg.dc_channel_info)
					// 清空变量
					AtSync_msg = AtSync_info{}
				}
				return
			}

			// 检查discord连接是否正常，否则重新连接
			discord_connection_check()

			// 消息发送人
			sender_name := sender.DisplayName
			if sender.DisplayName == "" {
				sender_name = sender.NickName
				print(sender_name)
			}

			// 群名 emoji 表情清除
			// discord 中 markdown []() 标签，不支持icon
			icon_str := regexp.MustCompile(` ?[\x{1F600}-\x{1F64F}\x{1F300}-\x{1F5FF}\x{1F680}-\x{1F6FF}\x{1F700}-\x{1F77F}\x{1F780}-\x{1F7FF}\x{1F800}-\x{1F8FF}\x{1F900}-\x{1F9FF}\x{1FA00}-\x{1FA6F}\x{2600}-\x{26FF}\x{2700}-\x{27BF}] ?`)
			sendgr_name := icon_str.ReplaceAllString(sendgr.NickName, "")

			if msg.IsPicture() {
				save_path := fmt.Sprintf("%s%s.jpg", conf.Static_path, cur_time)
				msg.SaveFileToLocal(save_path)
				discord_text_msg := fmt.Sprintf("> [%s](%s) - %s:\n", sendgr_name, DCNews_info.dc_channel_info, sender_name)
				discord_send_file(discord_text_msg, cur_time+".jpg", save_path, DCNews_info.dc_channel_id)
				return
			}
			fmt.Println(sender, err, sender_content, sendgr)

			// 格式化文本
			// 每行添加 >
			format_content := "> " + strings.ReplaceAll(sender_content, "\n", "\n> ")
			// 引用符合替换
			format_content = strings.ReplaceAll(format_content, "- - - - - - - - - - - - - - -", "-----------------------------")

			discord_text_msg := fmt.Sprintf("> [%s](%s) - %s:\n%s", sendgr_name, DCNews_info.dc_channel_info, sender_name, format_content)
			discord_send_text(discord_text_msg, DCNews_info.dc_channel_id)

			fmt.Println(*sendgr)
			fmt.Println(format_content)
			fmt.Println(sender_name, sender.NickName, sender.RemarkName, sender.DisplayName)
		}
	}

	bot.Block()
	// 关闭 Discord 连接
	s.Close()
}
