package config

import (
	"fmt"
	"text/template"

	"github.com/go-sql-driver/mysql"
	"github.com/spf13/viper"
	tb "gopkg.in/telebot.v3"
)

type RunType string


var (
	version = "dev"
	commit  = "none"
	date    = "unknown"
	ProjectName          string = "toshiki-rssbot"
	BotToken             string
	Socks5               string
	TelegraphToken       []string
	TelegraphAccountName string
	TelegraphAuthorName  string = "toshiki-rssbot"
	TelegraphAuthorURL   string

	// EnableTelegraph 是否启用telegraph
	EnableTelegraph       bool = false
	PreviewText           int  = 0
	DisableWebPagePreview bool = false
	mysqlConfig           *mysql.Config
	SQLitePath            string
	EnableMysql           bool = false

	// UpdateInterval RSS fetching interval
	UpdateInterval int = 10

	// ErrorThreshold Error threshold for RSS source fetching
	ErrorThreshold uint = 100

	// MessageTpl RSS update push template
	MessageTpl *template.Template

	// MessageMode Telegram message rendering mode
	MessageMode tb.ParseMode

	// TelegramEndpoint Telegram bot server address, default is empty
	TelegramEndpoint string = tb.DefaultApiURL

	// UserAgent User-Agent
	UserAgent string

	// RunMode Running mode Release / Debug
	RunMode RunType = ReleaseMode

	// AllowUsers Users allowed to use the bot
	AllowUsers []int64

	// DBLogMode Whether to print database logs
	DBLogMode bool = false
)

const (
	defaultMessageTplMode = tb.ModeHTML
	defaultMessageTpl     = `<b>{{.SourceTitle}}</b>{{ if .PreviewText }}
---------- Preview ----------
{{.PreviewText}}
-----------------------------
{{- end}}{{if .EnableTelegraph}}
{{.ContentTitle}} <a href="{{.TelegraphURL}}">Telegraph</a> | <a href="{{.RawLink}}">Original</a>
{{- else }}
<a href="{{.RawLink}}">{{.ContentTitle}}</a>
{{- end }}
{{.Tags}}
`
	defaultMessageMarkdownTpl = `** {{.SourceTitle}} **{{ if .PreviewText }}
---------- Preview ----------
{{.PreviewText}}
-----------------------------
{{- end}}{{if .EnableTelegraph}}
{{.ContentTitle}} [Telegraph]({{.TelegraphURL}}) | [Original]({{.RawLink}})
{{- else }}
[{{.ContentTitle}}]({{.RawLink}})
{{- end }}
{{.Tags}}
`
	TestMode    RunType = "Test"
	ReleaseMode RunType = "Release"
)

type TplData struct {
	SourceTitle     string
	ContentTitle    string
	RawLink         string
	PreviewText     string
	TelegraphURL    string
	Tags            string
	EnableTelegraph bool
}

func AppVersionInfo() (s string) {
	s = fmt.Sprintf("version %v, commit %v, built at %v", version, commit, date)
	return
}

// GetString get string config value by key
func GetString(key string) string {
	var value string
	if viper.IsSet(key) {
		value = viper.GetString(key)
	}

	return value
}

func GetMysqlDSN() string {
	return mysqlConfig.FormatDSN()
}
