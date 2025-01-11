package switchbot

import (
	"context"
	"fmt"
	"log"
	"log/slog"
	"time"

	_ "embed"

	"github.com/ericdaugherty/alexa-skills-kit-golang"
	"gopkg.in/yaml.v3"
)

//go:embed config/config.yaml
var configData []byte

type Mapping struct {
	Device   string     `yaml:"device"`
	Action   string     `yaml:"action"`
	Category string     `yaml:"category,omitempty"`
	DeviceID string     `yaml:"deviceID"`
	Commands []Commands `yaml:"commands"`
}

type Commands struct {
	Command     string `yaml:"command"`
	Parameter   string `yaml:"parameter"`
	CommandType string `yaml:"commandType"`
}

type Config struct {
	Mappings []Mapping `yaml:"mappings"`
}

type Handler struct {
	Config Config
	Client *client
}

func NewHandler() *Handler {
	var config Config
	err := yaml.Unmarshal(configData, &config)
	if err != nil {
		log.Fatalf("switchbot handler config read error: %v", err)
	}
	return &Handler{
		Client: NewClient(),
		Config: config,
	}
}

func (s Handler) OnSessionStarted(ctx context.Context, req *alexa.Request, session *alexa.Session, ctx2 *alexa.Context, res *alexa.Response) error {
	slog.InfoContext(ctx, "OnSessionStarted", "req", req, "session", session, "ctx", ctx2)
	return nil
}
func (s Handler) OnLaunch(ctx context.Context, req *alexa.Request, session *alexa.Session, ctx2 *alexa.Context, res *alexa.Response) error {
	slog.InfoContext(ctx, "OnLaunch", "req", req, "session", session, "ctx", ctx2)
	res.SetOutputText("スイッチボットを操作します。指示をしてください")
	res.ShouldSessionEnd = false
	return nil
}
func (s Handler) OnIntent(ctx context.Context, req *alexa.Request, session *alexa.Session, ctx2 *alexa.Context, res *alexa.Response) error {
	slog.InfoContext(ctx, "OnIntent", "req", req, "session", session, "ctx", ctx2)

	switch req.Intent.Name {
	case "SmartHomeIntent":
		slog.InfoContext(ctx, "SmartHomeIntent triggered")
		var deviceName, categoryName, actionName string
		for k, v := range req.Intent.Slots {
			if v.Resolutions == nil {
				// categoryに設定されてない場合
				continue
			}
			val := v.Resolutions.ResolutionsPerAuthority[0].Values[0].Value.Name
			switch k {
			case "device":
				deviceName = val
			case "category":
				categoryName = val
			case "action":
				actionName = val
			}
		}

		for _, m := range s.Config.Mappings {
			if m.Device == deviceName && m.Category == categoryName && m.Action == actionName {
				slog.InfoContext(ctx, "Matched", "device", m.Device, "category", m.Category, "action", m.Action)

				for _, v := range m.Commands {
					err := s.Client.ExecCommand(m.DeviceID, v)
					if err != nil {
						slog.ErrorContext(ctx, "failed to exec command", "command", v, "err", err)
						res.SetOutputText(fmt.Sprintf("コマンドの実行に失敗しました。ログを確認してください。: %v", err))
						return nil
					}
					time.Sleep(100 * time.Millisecond)
				}
				res.SetOutputText("スイッチボットを操作しました。続けて操作する場合は指示をしてください。")
				res.ShouldSessionEnd = false
			}
		}

	default:
		res.SetOutputText(fmt.Sprintf("一致するインテントがありませんでした。呼び出されたインテントは%sです", req.Intent.Name))
		return nil
	}
	return nil
}
func (s Handler) OnSessionEnded(ctx context.Context, req *alexa.Request, session *alexa.Session, ctx2 *alexa.Context, res *alexa.Response) error {
	slog.InfoContext(ctx, "OnSessionEnded", "req", req, "session", session, "ctx", ctx2)
	res.SetOutputText("スイッチボット操作を終了します")
	res.ShouldSessionEnd = true
	return nil
}
