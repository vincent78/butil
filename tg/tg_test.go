package tg

import (
	"context"
	"strings"
	"testing"
	"time"

	"gopkg.in/telebot.v3"
)

func TestSplitRunes(t *testing.T) {
	text := "你好🙂世界"
	parts := SplitRunes(text, 2)
	if len(parts) != 3 {
		t.Fatalf("expected 3 parts, got %d", len(parts))
	}
	if got := strings.Join(parts, ""); got != text {
		t.Fatalf("split should keep original text, got %q", got)
	}
	for _, part := range parts {
		if RuneLen(part) > 2 {
			t.Fatalf("part is longer than limit: %q", part)
		}
	}
}

func TestBuildMessageBatchesMergesShortMessages(t *testing.T) {
	items := []Message{
		{Node: "server1", CreatedAt: "2026-07-16 10:00:00", Msg: "price changed 1"},
		{Node: "server1", CreatedAt: "2026-07-16 10:00:01", Msg: "price changed 2"},
	}

	batches := BuildMessageBatches(items, "2026-07-16 10:00:03", 500, LongMessageOverhead)
	if len(batches) != 1 {
		t.Fatalf("expected 1 merged batch, got %d", len(batches))
	}
	if !strings.Contains(batches[0], "count: 2") {
		t.Fatalf("batch should contain merged count, got %q", batches[0])
	}
	if !strings.Contains(batches[0], "price changed 1") || !strings.Contains(batches[0], "price changed 2") {
		t.Fatalf("batch should contain both messages, got %q", batches[0])
	}
}

func TestBuildMessageBatchesSplitsLongMessage(t *testing.T) {
	items := []Message{
		{Node: "server1", CreatedAt: "2026-07-16 10:00:00", Msg: strings.Repeat("长", 900)},
	}

	batches := BuildMessageBatches(items, "2026-07-16 10:00:03", 300, LongMessageOverhead)
	if len(batches) < 2 {
		t.Fatalf("expected long message to be split into multiple batches, got %d", len(batches))
	}
	for _, batch := range batches {
		if RuneLen(batch) > 300 {
			t.Fatalf("batch is longer than limit: %d", RuneLen(batch))
		}
	}
	if !strings.Contains(strings.Join(batches, "\n"), "part:") {
		t.Fatalf("split batch should contain part marker, got %q", strings.Join(batches, "\n"))
	}
}

func TestNewHTTPClientWithoutProxy(t *testing.T) {
	client, err := NewHTTPClient("", 1500*time.Millisecond)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if client == nil {
		t.Fatal("expected http client")
	}
	if client.Timeout != 1500*time.Millisecond {
		t.Fatalf("unexpected timeout: %v", client.Timeout)
	}
	if client.Transport != nil {
		t.Fatalf("expected default transport without proxy, got %T", client.Transport)
	}
}

func TestNewHTTPClientWithProxy(t *testing.T) {
	client, err := NewHTTPClient("http://127.0.0.1:7890", 0)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if client == nil || client.Transport == nil {
		t.Fatal("expected http client with proxy transport")
	}
	if client.Timeout != DefaultHTTPTimeout {
		t.Fatalf("expected default timeout %v, got %v", DefaultHTTPTimeout, client.Timeout)
	}
}

func TestNewHTTPClientRejectsInvalidProxy(t *testing.T) {
	_, err := NewHTTPClient("://bad-proxy", time.Second)
	if err == nil {
		t.Fatal("expected invalid proxy error")
	}
	if !strings.Contains(err.Error(), "parse telegram proxy url") {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestNewRobotRequiresToken(t *testing.T) {
	_, err := NewRobot(RobotConfig{Offline: true}, nil)
	if err == nil {
		t.Fatal("expected empty token error")
	}
}

func TestNewRobotOfflineRegistersHandlers(t *testing.T) {
	called := false
	robot, err := NewRobot(RobotConfig{
		Token:       "123456:test-token",
		Offline:     true,
		PollTimeout: 2 * time.Second,
		HTTPTimeout: 3 * time.Second,
	}, func(bot *telebot.Bot) {
		called = bot != nil
	})
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if robot == nil || robot.Bot == nil {
		t.Fatal("expected robot and bot")
	}
	if !called {
		t.Fatal("expected register callback to be called")
	}
	if robot.Bot.Token != "123456:test-token" {
		t.Fatalf("unexpected bot token: %s", robot.Bot.Token)
	}
}

func TestStartRobotSkipsDebugBeforeTokenValidation(t *testing.T) {
	robot, err := StartRobot(context.Background(), RobotConfig{
		Debug:         true,
		SkipWhenDebug: true,
	}, nil)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if robot != nil {
		t.Fatal("expected nil robot when skipped in debug mode")
	}
}

func TestNotifierUsesInjectedSender(t *testing.T) {
	sent := make(chan string, 1)
	notifier := NewNotifier(NotifierConfig{
		Token:        "token",
		ChatID:       "chat",
		Node:         "server1",
		SendInterval: 10 * time.Millisecond,
		Now: func() string {
			return "2026-07-16 10:00:00"
		},
		Sender: func(_ context.Context, text string) error {
			sent <- text
			return nil
		},
	})

	notifier.Notify("hello")

	select {
	case text := <-sent:
		if !strings.Contains(text, "hello") {
			t.Fatalf("expected sent text to contain message, got %q", text)
		}
	case <-time.After(time.Second):
		t.Fatal("expected notifier to send queued message")
	}
}
