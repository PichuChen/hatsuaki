# Hatsuaki (初秋) 
是農曆的七月，也是這個專案第一版正式上線的月份。

這個專案的目的在於提供一個方便測試的 ActivityPub 實作，加速開發者對於 ActivityPub 的理解與開發速度。

## 特色
- 使用 JSON 作為資料庫格式，方便查看與修改
- 使用 Golang 開發，可以快速的建立與部署
- 使用純 JavaScript 進行開發，不需要額外的環境

## 使用方式

### 材料
1. 域名 (Domain) 一份
2. 可以接受聆聽 443 Port 的 IP 位址 一組
  或是可以利用 ngrok 來建立一個臨時的 HTTPS 連線

### 步驟
1. 下載二進位檔案
2. 執行

執行後會在該目錄中建立 `config.json`

他的長相會類似這樣
```json
{
  "domain": "example.com",
  "listen_address": ":8083"
}
```
請把這邊的 domain 改成你的域名

3. 重新執行

接著可以先透過 `http://localhost:8083` 來進行第一次的測試。
確定可以看到網頁之後就可以透過 nginx 或是其他的方式來將流量導到這個服務上。


