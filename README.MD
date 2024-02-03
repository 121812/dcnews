# <p align="center"> DCNews </p>
<p align="center">
  <a href="https://discord.com/invite/cAKzHRhwtZ"><img src="https://img.shields.io/badge/Discord-%E7%A4%BE%E5%8C%BA-blue" target="_blank" alt="Discord"></a>
  <a href="https://github.com/121812/dcnews#dcnews-%E7%A4%BE%E5%8C%BA"><img src="https://img.shields.io/badge/%E5%BE%AE%E4%BF%A1-%E4%BA%A4%E6%B5%81%E7%BE%A4-green" target="_blank" alt="微信交流群"></a>
  <a href="https://dcnews.gitbook.io/ru-he-shi-yong-tuo-guan-ban/xiao-bai-kuai-su-jie-ru-dcnews-tuo-guan-ban"><img src="https://img.shields.io/badge/%E6%89%98%E7%AE%A1%E7%89%88-gitbook-blue" target="_blank" alt="托管版"></a>
  <a href="https://dcnews.gitbook.io/ru-he-shi-yong-tuo-guan-ban/ru-he-wei-xiang-mu-jin-hang-gong-xian"><img src="https://img.shields.io/badge/%E5%A6%82%E4%BD%95%E8%B4%A1%E7%8C%AE-%E6%96%87%E6%A1%A3-green" target="_blank" alt="如何贡献"></a>
</p>
<p align="center"> 该工具能够即时捕捉微信群中的聊天消息，并将其自动同步到预设的Discord频道中 </p>

## 功能特性
* 文本/图片/接龙 消息实时同步至discord
* 支持一同步多或者多同步一
* PIN 功能 (开发中)
* DC斜杆命令配置同步项

## 安装部署
> 实现过程：当收到一条新的群消息后，查询数据库对应表项，获取消息同步至哪个 DC 频道并发送消息，所以要配置 config.json 文件，包括 dc 机器人 key 和数据库连接方式，导入表结构，配置好表，程序即可开始同步
### docker-compose 部署
1. 配置 docker-config.json 文件, 修改数据库密码，添加 Discord 机器人授权 token, 对应下面docker-compose.yaml文件
    ```go
    {
        "Discord_bot_auth": "",
        "Static_path": "/app/static/",
        "Mysql_host":"192.168.210.11",
        "Mysql_port":"3306",
        "Mysql_db":"dcnews",
        "Mysql_user":"root",
        "Mysql_password":"root",
        "Dc_createsync_prompts": "正在建立微信与Dc同步渠道...\nStep 1.请添加微信: \nStep 2.将该微信拉入目标微信群，等待30秒\nStep 3.在目标微信群输入同步码: "
    }
    ```
2. 在 docker-compose.yaml 中修改数据库映射端口和密码
    ```go
    mysql:
        image: mysql:5.7
        networks:
        dcnews_network:
            ipv4_address:
                192.168.210.11
        ports:
            - "23306:3306"
        environment:
            MYSQL_ROOT_PASSWORD: root
            MYSQL_DATABASE: dcnews
    ```

3. 启动docker compose
    ```
    docker compose up
    ```

4. 在日志中找到微信登录链接

5. 手动连接数据库，关联群组对应频道，即可开始同步
    ```go
    create_time	添加时间
    dc_user	dc用户名称
    wx_user	微信用户名称
    wx_group	微信群组名称 必填
    dc_channel_id	dc频道ID 必填
    dc_channel_info	dc频道介绍 必填，需为URL
    remark	备注
    ```

### 手动编译部署 
<details> 

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
    "Mysql_password":"",
    "Dc_createsync_prompts": "正在建立微信与Dc同步渠道...\nStep 1.请添加微信 \nStep 2.将该微信拉入目标微信群，等待30秒\nStep 3.在目标微信群输入同步码: "
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
</details>

## 使用DC斜杆命令创建同步
1. 在Dc中输入 /createsync 获取同步码
2. 在需要同步的微信群发送同步码，即可建立同步

