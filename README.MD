# DCNews
> 该工具能够即时捕捉微信群中的聊天消息，并将其自动同步到预设的Discord频道中

## 功能特性
* 文本/图片/接龙 消息实时同步至discord
* 支持一同步多或者多同步一
* PIN 功能 (开发中)

## 安装部署
> 实现过程：当收到一条新的群消息后，查询数据库对应表项，获取消息同步至哪个 DC 频道并发送消息，所以要配置 config.json 文件，包括 dc 机器人 key 和数据库连接方式，导入表结构，配置好表，程序即可开始同步
1. 下载依赖包：
    ```shell
    go mod download
    ```
2. 配置文件 
    ```go
    // 打开 config.json 配置 Discord 机器人key，静态文件存放路径，数据库连接方式
    {
    "Discord_bot_auth": "",
    "Static_path": "",
    "Mysql_host":"",
    "Mysql_port":"",
    "Mysql_db":"",
    "Mysql_user":"",
    "Mysql_password":""
    }
    ```
3. 导入表结构
    ```
    mysql -u -p
    source dc_wx_association_table.sql
    ```

4. 构建项目：
    ```shell
    go build
    ```
5. 运行项目:
    ```
    chmod 744 ./dcnews
    ./dcnews
    ```
6. 项目日志：
    ```shell
    logfile.log
    ```
7. 通过数据库关联群组对应频道
    ```
    create_time	添加时间
    dc_user	dc用户名称
    wx_user	微信用户名称
    wx_group	微信群组名称 必填
    dc_channel_id	dc频道ID 必填
    dc_channel_info	dc频道介绍 必填，需为URL
    remark	备注
    ```


## 在线托管&用户体验群
> 本项目支持在线托管

> ![Group](/Group.png)


