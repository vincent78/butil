package tg

/*
Telegram 通用能力说明

限制依据：
  - Telegram Bot API 的 sendMessage.text 限制为 1-4096 个字符（entities 解析后）。
    官方文档：https://core.telegram.org/bots/api#sendmessage
  - Telegram Bot FAQ 建议同一 chat 不要超过 1 条消息/秒；group 场景不超过 20 条消息/分钟；
    超过限制后会收到 429 Too Many Requests。
    官方 FAQ：https://core.telegram.org/bots/faq#my-bot-is-hitting-limits-how-do-i-avoid-this

当前能力：
  - SendText：直接发送一条 TG 文本消息。
  - Notifier：先记录消息，再进入内存队列，由后台 worker 按频率合并/分段发送。
  - Robot：封装 telebot 的 HTTP client/proxy、bot 创建、handler 注册、异步启动和 debug 跳过。
  - SplitRunes / BuildMessageBatches：提供可复用、可单测的分段和合并逻辑。
*/

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net"
	"net/http"
	"net/url"
	"strings"
	"sync"
	"time"

	log "github.com/vincent78/butil/logger/logger4"
	"github.com/vincent78/butil/net/gchttp"
	"gopkg.in/telebot.v3"
)

const (
	MaxMessageChars     = 4096
	SafeMessageChars    = 3800
	DefaultQueueSize    = 1024
	DefaultSendInterval = 3 * time.Second
	LongMessageOverhead = 256
	DefaultPollTimeout  = 10 * time.Second
	DefaultHTTPTimeout  = 60 * time.Second
)

type Message struct {
	Node      string
	Msg       string
	CreatedAt string
}

type SendConfig struct {
	Token       string
	ChatID      string
	Proxy       string
	HTTPTimeout time.Duration
	Logger      *log.NormalLogger
}

type NotifierConfig struct {
	Token               string
	ChatID              string
	Proxy               string
	Node                string
	Debug               bool
	SkipWhenDebug       bool
	QueueSize           int
	SendInterval        time.Duration
	MaxMessageChars     int
	SafeMessageChars    int
	LongMessageOverhead int
	HTTPTimeout         time.Duration
	Now                 func() string
	Logger              *log.NormalLogger
	Sender              func(ctx context.Context, text string) error
}

type Notifier struct {
	cfg   NotifierConfig
	queue chan Message
	once  sync.Once
}

type RobotConfig struct {
	Token         string
	Proxy         string
	Debug         bool
	PollTimeout   time.Duration
	HTTPTimeout   time.Duration
	StartAsync    bool
	SkipWhenDebug bool
	Offline       bool
	Logger        *log.NormalLogger
}

type RobotRegister func(bot *telebot.Bot)

type Robot struct {
	Bot *telebot.Bot
}

func ManagerURL(token string) string {
	return fmt.Sprintf("https://api.telegram.org/bot%s/sendMessage", token)
}

func NewNotifier(cfg NotifierConfig) *Notifier {
	cfg = normalizeNotifierConfig(cfg)
	return &Notifier{
		cfg:   cfg,
		queue: make(chan Message, cfg.QueueSize),
	}
}

func (n *Notifier) Notify(msg string) {
	if n == nil {
		return
	}
	item := Message{
		Node:      n.cfg.Node,
		Msg:       msg,
		CreatedAt: n.cfg.Now(),
	}

	n.logger().Info("telegram message generated",
		log.String("node", item.Node),
		log.Int("chars", RuneLen(item.Msg)),
		log.String("message", item.Msg),
	)

	if n.cfg.SkipWhenDebug && n.cfg.Debug {
		n.logger().Info("telegram message skipped in debug mode",
			log.String("node", item.Node),
		)
		return
	}

	if n.cfg.Token == "" || n.cfg.ChatID == "" {
		n.logger().Warn("telegram message skipped because token or chat_id is empty",
			log.String("node", item.Node),
		)
		return
	}

	n.Enqueue(item)
}

func (n *Notifier) Enqueue(item Message) {
	if n == nil {
		return
	}
	n.start()

	select {
	case n.queue <- item:
		return
	default:
	}

	select {
	case dropped := <-n.queue:
		n.logger().Warn("telegram queue is full, dropped oldest message",
			log.String("node", dropped.Node),
			log.String("created_at", dropped.CreatedAt),
			log.Int("chars", RuneLen(dropped.Msg)),
		)
	default:
	}

	select {
	case n.queue <- item:
	default:
		n.logger().Error("telegram queue is full, dropped newest message",
			log.String("node", item.Node),
			log.String("created_at", item.CreatedAt),
			log.Int("chars", RuneLen(item.Msg)),
		)
	}
}

func (n *Notifier) start() {
	n.once.Do(func() {
		go n.run()
	})
}

func (n *Notifier) run() {
	ticker := time.NewTicker(n.cfg.SendInterval)
	defer ticker.Stop()

	pending := make([]Message, 0, 16)
	for {
		select {
		case item := <-n.queue:
			pending = append(pending, item)
		case <-ticker.C:
			if len(pending) == 0 {
				continue
			}
			batches := BuildMessageBatches(pending, n.cfg.Now(), n.cfg.SafeMessageChars, n.cfg.LongMessageOverhead)
			pending = pending[:0]

			for i, text := range batches {
				if i > 0 {
					time.Sleep(n.cfg.SendInterval)
				}
				if err := n.cfg.Sender(context.Background(), text); err != nil {
					n.logger().Error("telegram message send failed",
						log.Int("chars", RuneLen(text)),
						log.String("error", err.Error()),
					)
					continue
				}
				n.logger().Info("telegram message sent",
					log.Int("chars", RuneLen(text)),
				)
			}
		}
	}
}

func BuildMessageBatches(items []Message, sentAt string, maxChars int, longMessageOverhead int) []string {
	if len(items) == 0 {
		return nil
	}
	if maxChars <= 0 || maxChars > MaxMessageChars {
		maxChars = SafeMessageChars
	}
	if longMessageOverhead <= 0 {
		longMessageOverhead = LongMessageOverhead
	}

	node := batchNodeName(items)
	entries := expandEntries(items, maxChars, longMessageOverhead)
	batches := make([][]string, 0, 1)
	current := make([]string, 0, len(entries))

	for _, entry := range entries {
		next := append(append([]string{}, current...), entry)
		if len(current) > 0 && RuneLen(FormatBatch(node, sentAt, next)) > maxChars {
			batches = append(batches, current)
			current = []string{entry}
			continue
		}
		current = next
	}
	if len(current) > 0 {
		batches = append(batches, current)
	}

	texts := make([]string, 0, len(batches))
	for _, batch := range batches {
		text := FormatBatch(node, sentAt, batch)
		if RuneLen(text) <= maxChars {
			texts = append(texts, text)
			continue
		}
		texts = append(texts, SplitRunes(text, maxChars)...)
	}
	return texts
}

func FormatBatch(node string, sentAt string, entries []string) string {
	return fmt.Sprintf("from: %v\ntime(utc): %v\ncount: %d\ncontext:\n%v",
		node,
		sentAt,
		len(entries),
		strings.Join(entries, "\n---\n"),
	)
}

func SplitRunes(text string, limit int) []string {
	if limit <= 0 {
		limit = SafeMessageChars
	}
	if text == "" {
		return []string{""}
	}

	runes := []rune(text)
	if len(runes) <= limit {
		return []string{text}
	}

	parts := make([]string, 0, len(runes)/limit+1)
	for start := 0; start < len(runes); start += limit {
		end := start + limit
		if end > len(runes) {
			end = len(runes)
		}
		parts = append(parts, string(runes[start:end]))
	}
	return parts
}

func RuneLen(text string) int {
	return len([]rune(text))
}

func SendText(ctx context.Context, cfg SendConfig, text string) error {
	if cfg.Token == "" {
		return errors.New("telegram token is empty")
	}
	if cfg.ChatID == "" {
		return errors.New("telegram chat_id is empty")
	}
	if ctx == nil {
		ctx = context.Background()
	}

	url := ManagerURL(cfg.Token)
	payload, err := json.Marshal(map[string]interface{}{
		"chat_id": cfg.ChatID,
		"text":    text,
	})
	if err != nil {
		return err
	}

	proxyClient := gchttp.NewHttpClient(url,
		gchttp.WithProxy(cfg.Proxy),
		gchttp.WithTimeout(httpTimeout(cfg.HTTPTimeout)),
		gchttp.WithLogger(normalLogger(cfg.Logger)),
	)
	task := gchttp.NewTask(ctx, url, nil)
	task.Method = "POST"
	task.Body = string(payload)

	resp := gchttp.Done(task, proxyClient)
	if !resp.IsSuccess() {
		return resp.Error()
	}

	jsonResp := resp.GetJson()
	if jsonResp.Exists() && jsonResp.Get("ok").Exists() && !jsonResp.Get("ok").Bool() {
		return fmt.Errorf("telegram api error code=%d description=%s retry_after=%d",
			jsonResp.Get("error_code").Int(),
			jsonResp.Get("description").String(),
			jsonResp.Get("parameters.retry_after").Int(),
		)
	}
	return nil
}

func StartRobot(ctx context.Context, cfg RobotConfig, register RobotRegister) (*Robot, error) {
	if cfg.SkipWhenDebug && cfg.Debug {
		normalLogger(cfg.Logger).Info("telegram robot skipped in debug mode")
		return nil, nil
	}

	robot, err := NewRobot(cfg, register)
	if err != nil {
		return nil, err
	}

	if cfg.StartAsync {
		go robot.Bot.Start()
	} else {
		robot.Bot.Start()
	}

	if ctx != nil {
		go func() {
			<-ctx.Done()
			robot.Stop()
		}()
	}

	normalLogger(cfg.Logger).Info("telegram robot started")
	return robot, nil
}

func NewRobot(cfg RobotConfig, register RobotRegister) (*Robot, error) {
	if cfg.Token == "" {
		return nil, errors.New("telegram robot token is empty")
	}

	httpClient, err := NewHTTPClient(cfg.Proxy, cfg.HTTPTimeout)
	if err != nil {
		return nil, err
	}

	pollTimeout := cfg.PollTimeout
	if pollTimeout <= 0 {
		pollTimeout = DefaultPollTimeout
	}

	bot, err := telebot.NewBot(telebot.Settings{
		Token:   cfg.Token,
		Poller:  &telebot.LongPoller{Timeout: pollTimeout},
		Client:  httpClient,
		Offline: cfg.Offline,
	})
	if err != nil {
		return nil, err
	}

	if register != nil {
		register(bot)
	}

	return &Robot{Bot: bot}, nil
}

func NewHTTPClient(proxyConfig string, timeout time.Duration) (*http.Client, error) {
	timeout = httpTimeout(timeout)

	client := &http.Client{Timeout: timeout}
	if proxyConfig == "" {
		return client, nil
	}

	proxyURL, err := url.Parse(proxyConfig)
	if err != nil {
		return nil, fmt.Errorf("parse telegram proxy url: %w", err)
	}

	client.Transport = &http.Transport{
		Proxy: http.ProxyURL(proxyURL),
		DialContext: (&net.Dialer{
			Timeout:   30 * time.Second,
			KeepAlive: 30 * time.Second,
		}).DialContext,
	}
	return client, nil
}

func (r *Robot) Stop() {
	if r == nil || r.Bot == nil {
		return
	}
	r.Bot.Stop()
}

func expandEntries(items []Message, maxChars int, longMessageOverhead int) []string {
	entries := make([]string, 0, len(items))
	bodyLimit := maxChars - longMessageOverhead
	if bodyLimit < 1000 {
		bodyLimit = maxChars
	}

	for _, item := range items {
		msgParts := SplitRunes(item.Msg, bodyLimit)
		for i, part := range msgParts {
			partMark := ""
			if len(msgParts) > 1 {
				partMark = fmt.Sprintf(" part:%d/%d", i+1, len(msgParts))
			}
			entries = append(entries, fmt.Sprintf("created(utc): %s%s\n%s", item.CreatedAt, partMark, part))
		}
	}
	return entries
}

func batchNodeName(items []Message) string {
	node := items[0].Node
	for _, item := range items[1:] {
		if item.Node != node {
			return "mixed"
		}
	}
	return node
}

func normalizeNotifierConfig(cfg NotifierConfig) NotifierConfig {
	if cfg.QueueSize <= 0 {
		cfg.QueueSize = DefaultQueueSize
	}
	if cfg.SendInterval <= 0 {
		cfg.SendInterval = DefaultSendInterval
	}
	if cfg.SafeMessageChars <= 0 || cfg.SafeMessageChars > MaxMessageChars {
		cfg.SafeMessageChars = SafeMessageChars
	}
	if cfg.MaxMessageChars <= 0 || cfg.MaxMessageChars > MaxMessageChars {
		cfg.MaxMessageChars = MaxMessageChars
	}
	if cfg.LongMessageOverhead <= 0 {
		cfg.LongMessageOverhead = LongMessageOverhead
	}
	if cfg.Now == nil {
		cfg.Now = func() string {
			return time.Now().UTC().Format(time.RFC3339)
		}
	}
	if cfg.Logger == nil {
		cfg.Logger = log.DefaultNormalLogger()
	}
	if cfg.Sender == nil {
		sendCfg := SendConfig{
			Token:       cfg.Token,
			ChatID:      cfg.ChatID,
			Proxy:       cfg.Proxy,
			HTTPTimeout: cfg.HTTPTimeout,
			Logger:      cfg.Logger,
		}
		cfg.Sender = func(ctx context.Context, text string) error {
			return SendText(ctx, sendCfg, text)
		}
	}
	return cfg
}

func normalLogger(logger *log.NormalLogger) *log.NormalLogger {
	if logger != nil {
		return logger
	}
	return log.DefaultNormalLogger()
}

func (n *Notifier) logger() *log.NormalLogger {
	return normalLogger(n.cfg.Logger)
}

func httpTimeout(timeout time.Duration) time.Duration {
	if timeout <= 0 {
		return DefaultHTTPTimeout
	}
	return timeout
}
