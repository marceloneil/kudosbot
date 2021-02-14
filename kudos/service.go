package kudos

import (
	"fmt"
	"kudosbot/pkg/config"
	"kudosbot/repository"
	"log"
	"math"
	"os"
	"strconv"
	"strings"
	"time"

	"go.uber.org/fx"
	"go.uber.org/zap"

	"github.com/slack-go/slack"
	"github.com/slack-go/slack/slackevents"
	"github.com/slack-go/slack/socketmode"
)

type Service struct {
	logger          *zap.SugaredLogger
	apiClient       *slack.Client
	socketClient    *socketmode.Client
	usersMap        map[string]*slack.User
	kudosRepository *repository.KudosRepository
}

func NewService(
	lc fx.Lifecycle,
	config *config.Config,
	kudosRepository *repository.KudosRepository,
) *Service {
	logger, err := zap.NewProduction()
	if err != nil {
		log.Fatal(err)
	}
	apiClient := slack.New(
		config.Server.SlackBotToken,
		slack.OptionDebug(true),
		slack.OptionLog(log.New(os.Stdout, "api: ", log.Lshortfile|log.LstdFlags)),
		slack.OptionAppLevelToken(config.Server.SlackAppToken),
	)
	socketClient := socketmode.New(
		apiClient,
		socketmode.OptionDebug(true),
		socketmode.OptionLog(log.New(os.Stdout, "socketmode: ", log.Lshortfile|log.LstdFlags)),
	)
	service := &Service{
		logger:          logger.Sugar(),
		apiClient:       apiClient,
		socketClient:    socketClient,
		kudosRepository: kudosRepository,
	}
	lc.Append(fx.Hook{
		OnStart: service.start,
	})
	return service
}

func (service *Service) handleReaction(
	eventID string, event *slackevents.ReactionAddedEvent,
) {
	// only care about dogecoin :)
	if event.Reaction != "dogecoin" {
		return
	}
	// cannot kudos yourselves
	if event.User == event.ItemUser {
		return
	}
	sendUser, err := service.getUser(event.User)
	if err != nil {
		service.logger.Errorf("cannot find user id=%s", event.User)
		return
	}
	receiveUser, err := service.getUser(event.ItemUser)
	if err != nil {
		service.logger.Errorf("cannot find user id=%s", event.ItemUser)
		return
	}
	if receiveUser.IsBot {
		return
	}
	numKudosRemaining, err := service.getNumKudosRemaining(event.User)
	if err != nil {
		return
	}
	// sender is out of kudos
	if numKudosRemaining <= 0 {
		return
	}

	timeFloat, err := strconv.ParseFloat(event.EventTimestamp, 64)
	if err != nil {
		return
	}
	eventTimestamp := time.Unix(int64(timeFloat), 0)
	_, err = service.kudosRepository.UpsertKudos(
		eventTimestamp, sendUser.ID, receiveUser.ID, eventID, event.Reaction)
	if err != nil {
		return
	}
}

func (service *Service) handleCommands(
	event *slackevents.AppMentionEvent,
) {
	messageTokens := strings.Split(event.Text, " ")
	command := messageTokens[1]
	if command == "leaderboard" {
		service.handleLeaderboardCommand(event)
	} else if command == "dogecoin" {
		service.handleDogeCommand(event)
	}
}

func (service *Service) handleDogeCommand(
	event *slackevents.AppMentionEvent,
) {
	numKudosRemaining, err := service.getNumKudosRemaining(event.User)
	if err != nil {
		return
	}
	message := fmt.Sprintf("You have %d dogecoin left to give today!", numKudosRemaining)
	_, _, err = service.apiClient.PostMessage(event.Channel,
		slack.MsgOptionText(message, false))
	if err != nil {
		fmt.Printf("failed posting message: %v", err)
	}
}

func (service *Service) handleLeaderboardCommand(
	event *slackevents.AppMentionEvent,
) {
	leaderboard, err := service.kudosRepository.GetLeaderboard(nil)
	if err != nil {
		service.logger.Errorf("fail to get leaderboard err=%v", err)
		return
	}
	leaderboardString := "#   Name                Doge\n"
	for index, board := range leaderboard {
		position := index + 1
		positionPadded := fmt.Sprintf("%d.", position)
		if len(positionPadded) < 3 {
			positionPadded += " "
		}
		user := service.usersMap[board.ReceiverID]
		namePadded := user.RealName
		if len(namePadded) < 19 {
			for len(namePadded) < 19 {
				namePadded += " "
			}
		} else if len(namePadded) > 19 {
			namePadded = namePadded[0:16] + "..."
		}
		leaderboardString += fmt.Sprintf("%s %s %d\n",
			positionPadded, namePadded, board.Count)
	}

	message := fmt.Sprintf("*Today's Leaderboard*\n```%s```", leaderboardString)
	_, _, err = service.apiClient.PostMessage(event.Channel,
		slack.MsgOptionText(message, false))
	if err != nil {
		service.logger.Errorf("fail to post message err=%v", err)
	}
}

func (service *Service) handleChannelCreated(
	event *slack.ChannelCreatedEvent,
) {
	_, _, _, err := service.apiClient.JoinConversation(event.Channel.ID)
	if err != nil {
		service.logger.Errorf("failed to join public channel err=%v", err)
	}
}

func (service *Service) getNumKudosRemaining(
	userID string,
) (int, error) {
	numDogeGiven, err := service.kudosRepository.CountKudosGivenForToday(userID)
	if err != nil {
		service.logger.Errorf("fail to count kudos given err=%v", err)
		return 0, err
	}
	return int(math.Max(float64(10-numDogeGiven), 0)), nil
}
