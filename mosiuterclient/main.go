package main

import (
	"fmt"
	"log"
	"net/url"
	"os"

	"github.com/gorilla/websocket"
	"gopkg.in/yaml.v3"
)

type Proxy struct {
	Name   string `yaml:"name"`
	Server string `yaml:"server"`
	Port   int    `yaml:"port"`
}
type Config struct {
	Proxies []Proxy `yaml:"proxies"`
}

func loadConfig(path string) (*Config, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var cfg Config
	decoder := yaml.NewDecoder(file)
	err = decoder.Decode(&cfg)
	if err != nil {
		return nil, err
	}

	return &cfg, nil
}
func buildWSURL(p Proxy) string {
	return fmt.Sprintf("wss://%s:%d/ws", p.Server, p.Port)
}

func main() {
	println("Welcome to MoSIUTER.")
	println("If you saw this , you client is running and active.")
	cfg, err := loadConfig("config.yaml")
	if err != nil {
		log.Fatalf("讀取配置失敗: %v", err)
	}

	if len(cfg.Proxies) == 0 {
		log.Fatal("沒有發現任何代理節點！")
	}

	// 選擇第一個代理進行連接
	proxy := cfg.Proxies[0]
	wsURL := buildWSURL(proxy)

	log.Printf("連接到 [%s] %s", proxy.Name, wsURL)
	u, err := url.Parse(wsURL)
	if err != nil {
		log.Fatalf("URL 解析錯誤: %v", err)
	}

	c, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
	if err != nil {
		log.Fatalf("連接失敗: %v", err)
	}
	defer c.Close()

	log.Println("連線成功，發送測試訊息")
	err = c.WriteMessage(websocket.TextMessage, []byte("你好，自由世界！"))
	if err != nil {
		log.Printf("寫入錯誤: %v", err)
		return
	}

	_, message, err := c.ReadMessage()
	if err != nil {
		log.Printf("讀取錯誤: %v", err)
		return
	}

	log.Printf("收到回應: %s", message)
}
