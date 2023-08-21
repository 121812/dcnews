package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"os"

	_ "github.com/go-sql-driver/mysql"
)

var db *sql.DB

type Conf_connect_db struct {
	Mysql_host     string
	Mysql_port     string
	Mysql_db       string
	Mysql_user     string
	Mysql_password string
}

type Wechat_chat_log struct {
	Time         string
	Send_user    string
	Send_content string
	Send_group   string
}

type DCNews_info struct {
	dc_channel_id   string
	dc_channel_info string
}

func init_db() {

	// 打开文件
	config_file, _ := os.Open("config.json")

	// 关闭文件
	defer config_file.Close()

	//NewDecoder创建一个从file读取并解码json对象的*Decoder，解码器有自己的缓冲，并可能超前读取部分json数据。
	decoder := json.NewDecoder(config_file)

	conf := Conf_connect_db{}
	//Decode从输入流读取下一个json编码值并保存在v指向的值里
	decoder.Decode(&conf)
	database_connect_str := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", conf.Mysql_user, conf.Mysql_password, conf.Mysql_host, conf.Mysql_port, conf.Mysql_db)
	db, _ = sql.Open("mysql", database_connect_str)

	//设置数据库最大连接数
	db.SetConnMaxLifetime(1024)
	//设置上数据库最大闲置连接数
	db.SetMaxIdleConns(64)
	//验证连接
	if err := db.Ping(); err != nil {
		fmt.Println("open database fail")
		return
	}
}

func insert_wechat_chat_log(wechat_chat_log Wechat_chat_log) bool {
	//准备sql语句
	insert_sql := fmt.Sprintf("INSERT INTO `wechat_chat_log`(`time`, `send_user`, `send_content`, `send_group`) VALUES (\"%s\", \"%s\", \"%s\", \"%s\");", wechat_chat_log.Time, wechat_chat_log.Send_user, wechat_chat_log.Send_content, wechat_chat_log.Send_group)
	_, err := db.Exec(insert_sql)
	if err != nil {
		fmt.Println("Failed to execute SQL statement:", err)
		return false
	}

	return true
}

// 对dcnews进行判断，正在则返回channel id，负责返回false
func judge_dcnews_state(sendgr string) (DCNews_info, error) {
	var DCNews_info DCNews_info
	select_sql := "select dc_channel_id, dc_channel_info from dc_wx_association_table where wx_group = ?"
	err := db.QueryRow(select_sql, sendgr).Scan(&DCNews_info.dc_channel_id, &DCNews_info.dc_channel_info)
	if err != nil {
		fmt.Println("Failed to execute SQL statement:", err)
		return DCNews_info, err
	}
	return DCNews_info, nil
}
