package main

import (
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/slack-go/slack"
	"github.com/slack-go/slack/socketmode"
)

type TimeTrackerCommand struct {
	client     *socketmode.Client
	controller *Controller
}

const (
	month    = "200601"
	workTime = "20060102 15:04"
)

func NewTimeTrackerCommand(c *socketmode.Client, db *sqlx.DB) *TimeTrackerCommand {
	controller := NewController(db)
	return &TimeTrackerCommand{
		client:     c,
		controller: controller,
	}
}

func (tt *TimeTrackerCommand) Do(evt socketmode.Event, cmd slack.SlashCommand) error {
	log.Printf("slack cmd text:%s", cmd.Text)
	subcmd := strings.Split(cmd.Text, " ")
	if len(subcmd) == 0 {
		// show worktime
	}
	blocks := []slack.Block{}
	switch subcmd[0] {
	case "start":
		blocks = tt.Start()
	case "stop":
		blocks = tt.Stop(subcmd)
	case "show":
		blocks = tt.Show(subcmd)
	case "detail":
		blocks = tt.Details(subcmd)
	case "rstart":
		blocks = tt.RestStart()
	case "rstop":
		blocks = tt.RestStop()
	case "add":
		blocks = tt.Add(subcmd)
	default:
		blocks = tt.Help()
	}
	tt.client.Ack(*evt.Request, nil)
	tt.client.PostMessage(cmd.ChannelID, slack.MsgOptionBlocks(blocks...))
	return nil
}
func (tt *TimeTrackerCommand) Start() []slack.Block {
	if err := tt.controller.Start(attendance); err != nil {
		return tt.makePayload(fmt.Sprintf("err:%v", err))
	}
	return tt.makePayload("start tracking")
}

func (tt *TimeTrackerCommand) Stop(subCmd []string) []slack.Block {
	if len(subCmd) > 2 {
		return tt.makePayload("invalid command")
	}
	finishedAt := time.Now()
	if len(subCmd) > 1 {
		dayStr := subCmd[1]
		var err error
		finishedAt, err = time.Parse(time.RFC3339, dayStr)
		if err != nil {
			return tt.makePayload(fmt.Sprintf("stop is error by parse finished time,err:%v", err))
		}
	}
	tracks, err := tt.controller.Stop(finishedAt.UTC())
	if err != nil {
		return tt.makePayload(fmt.Sprintf("failed to save finish,err:%v", err))
	}
	var message string
	if len(tracks) == 1 {
		message = showWorkTime(tracks[0])
	} else {
		message = stopMultipleMessage(tracks)
	}
	return tt.makePayload(message)
}

func (tt *TimeTrackerCommand) RestStart() []slack.Block {
	if err := tt.controller.Start(rest); err != nil {
		return tt.makePayload(fmt.Sprintf("failed to start rest,err:%v", err))
	}
	return tt.makePayload("rest start")
}
func (tt *TimeTrackerCommand) RestStop() []slack.Block {
	if err := tt.controller.StopRest(time.Now()); err != nil {
		return tt.makePayload(fmt.Sprintf("failed to finish rest,err:%v", err))
	}
	return tt.makePayload("stop rest")
}

func (tt *TimeTrackerCommand) Add(subCmd []string) []slack.Block {

	if len(subCmd) != 3 {
		return tt.makePayload("invalid args...")
	}
	start, err := time.ParseInLocation(workTime, subCmd[1], jst)
	if err != nil {
		return tt.makePayload(fmt.Sprintf("invalid start_at format ex.) %s", workTime))
	}
	finishedAt, err := time.ParseInLocation(workTime, subCmd[2], jst)
	if err != nil {
		return tt.makePayload(fmt.Sprintf("invalid finish_at format ex.) %s", workTime))
	}
	if err := tt.controller.Add(start.UTC(), finishedAt.UTC()); err != nil {
		return tt.makePayload("invalid finished_at format,")
	}
	return tt.makePayload(fmt.Sprintf("add success, start:%s end:%s", start.String(), finishedAt.String()))
}

func (tt *TimeTrackerCommand) Show(subCmd []string) []slack.Block {
	showDate := time.Now()
	if len(subCmd) == 2 {
		var err error
		showDate, err = time.Parse(month, subCmd[1])
		if err != nil {
			return tt.makePayload(fmt.Sprintf("stop is error by parse finished time,err:%v", err))
		}
	}
	showDateJST := toJST(showDate)
	b, e := BeginAndLateDateInMonth(showDateJST)
	worktime, err := tt.controller.ShowWorkTimeIn(b.UTC(), e.UTC())
	if err != nil {
		return tt.makePayload(fmt.Sprintf("failed to calculate:%v", err))
	}
	return tt.makePayload(fmt.Sprintf("worktime:`%s`", worktime.String()))
}

func (tt *TimeTrackerCommand) Details(subCmd []string) []slack.Block {
	showDate := time.Now()
	if len(subCmd) == 2 {
		var err error
		showDate, err = time.Parse(month, subCmd[1])
		if err != nil {
			return tt.makePayload(fmt.Sprintf("stop is error by parse finished time,err:%v", err))
		}
	}
	showDateJST := toJST(showDate)
	b, e := BeginAndLateDateInMonth(showDateJST)
	tracks, err := tt.controller.ShowIn(b.UTC(), e.UTC())
	if err != nil {
		return tt.makePayload(fmt.Sprintf("failed to get show in:%v", err))
	}
	return tt.makePayload(showWorkTimes(showDateJST, tracks))
}

const helpMessage = "勤怠開始\n `/tm start`\n" +
	"勤怠終了\n `/tm stop [default now or yyyymm]` \n" +
	"休憩開始\n `/tm rstart ` \n" +
	"休憩終了\n `/tm rstop ` \n" +
	"勤怠時間表示\n `/tm show ` \n" +
	"詳細表示\n `/tm detail`" +
	"勤怠追加\n `/tm add [yyyymmdd HH:MM]`" +
	"ヘルプ\n `/tm help ` \n"

func (tt *TimeTrackerCommand) Help() []slack.Block {
	return tt.makePayload(helpMessage)
}

func (tt *TimeTrackerCommand) makePayload(markdownText string) []slack.Block {
	return []slack.Block{
		slack.NewSectionBlock(
			&slack.TextBlockObject{
				Type: slack.MarkdownType,
				Text: markdownText,
			},
			nil,
			nil,
		),
	}
}
