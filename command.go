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

const month = "200601"

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
	payload := map[string]interface{}{}
	switch subcmd[0] {
	case "start":
		payload = tt.Start()
	case "stop":
		payload = tt.Stop(subcmd)
	case "show":
		payload = tt.Show(subcmd)
	case "detail":
		payload = tt.Show(subcmd)
	case "rstart":
		payload = tt.RestStart()
	case "rstop":
		payload = tt.RestStop()
	default:
		payload = tt.Help()
	}
	tt.client.Ack(*evt.Request, payload)
	return nil
}
func (tt *TimeTrackerCommand) Start() map[string]interface{} {
	if err := tt.controller.Start(attendance); err != nil {
		return tt.makePayload(fmt.Sprintf("err:%v", err))
	}
	return tt.makePayload("start tracking")
}

func (tt *TimeTrackerCommand) Stop(subCmd []string) map[string]interface{} {
	if len(subCmd) > 2 {
		return tt.makePayload("invalid command")
	}
	if len(subCmd) == 1 {
		_, err := tt.controller.Stop(time.Now())
		if err != nil {
			return tt.makePayload(fmt.Sprintf("failed to save finish,err:%v", err))
		}
		return tt.makePayload("finish success")
	}
	dayStr := subCmd[1]
	finishedAt, err := time.Parse(time.RFC3339, dayStr)
	if err != nil {
		return tt.makePayload(fmt.Sprintf("stop is error by parse finished time,err:%v", err))
	}
	_, err = tt.controller.Stop(finishedAt.UTC())
	if err != nil {
		return tt.makePayload(fmt.Sprintf("failed to save finish,err:%v", err))
	}
	return tt.makePayload("finish success")
}

func (tt *TimeTrackerCommand) RestStart() map[string]interface{} {
	if err := tt.controller.Start(rest); err != nil {
		return tt.makePayload(fmt.Sprintf("failed to start rest,err:%v", err))
	}
	return tt.makePayload("rest start")
}
func (tt *TimeTrackerCommand) RestStop() map[string]interface{} {
	if err := tt.controller.StopRest(time.Now()); err != nil {
		return tt.makePayload(fmt.Sprintf("failed to finish rest,err:%v", err))
	}
	return tt.makePayload("stop rest")
}

func (tt *TimeTrackerCommand) Show(subCmd []string) map[string]interface{} {
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
	return tt.makePayload(fmt.Sprintf("worktime:%v", worktime.String()))
}

const helpMessage = "勤怠開始\n `/tm start`\n" +
	"勤怠終了\n `/tm stop [default now or yyyymm]` \n" +
	"休憩開始\n `/tm rstart ` \n" +
	"休憩終了\n `/tm rstop ` \n" +
	"勤怠時間表示\n `/tm show ` \n" +
	"詳細表示\n `/tm detail`" +
	"ヘルプ\n `/tm help ` \n"

func (tt *TimeTrackerCommand) Help() map[string]interface{} {
	return tt.makePayload(helpMessage)
}

func (tt *TimeTrackerCommand) makePayload(markdownText string) map[string]interface{} {
	return map[string]interface{}{
		"blocks": []slack.Block{
			slack.NewSectionBlock(
				&slack.TextBlockObject{
					Type: slack.MarkdownType,
					Text: markdownText,
				},
				nil,
				nil,
			),
		}}
}
