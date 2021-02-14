package kudos

import (
	"log"

	"github.com/slack-go/slack"
	"github.com/slack-go/slack/slackevents"
	"github.com/slack-go/slack/socketmode"
)

func (service *Service) handleEventAPIEvents(
	evt *socketmode.Event,
) {
	eventsAPIEvent, ok := evt.Data.(slackevents.EventsAPIEvent)
	if !ok {
		log.Printf("Ignored %+v\n", evt)
		return
	}
	switch eventsAPIEvent.Type {
	case slackevents.CallbackEvent:
		eventData := eventsAPIEvent.Data.(*slackevents.EventsAPICallbackEvent)
		innerEvent := eventsAPIEvent.InnerEvent
		switch ev := innerEvent.Data.(type) {
		// this is when someone is talking to me
		case *slackevents.AppMentionEvent:
			service.handleCommands(ev)
		case *slackevents.ReactionAddedEvent:
			service.handleReaction(eventData.EventID, ev)
		case *slack.ChannelCreatedEvent:
			service.handleChannelCreated(ev)
		}
	default:
		service.socketClient.Debugf("unsupported Events API event received")
	}
	service.socketClient.Ack(*evt.Request)
}

func (service *Service) handleInteractionEvents(
	evt *socketmode.Event,
) {
	callback, ok := evt.Data.(slack.InteractionCallback)
	if !ok {
		log.Printf("Ignored %+v\n", evt)
		return
	}

	log.Printf("Interaction received: %+v\n", callback)
	service.socketClient.Ack(*evt.Request)
}

func (service *Service) handleSlashCommandEvents(
	evt *socketmode.Event,
) {
	cmd, ok := evt.Data.(slack.SlashCommand)
	if !ok {
		log.Printf("Ignored %+v\n", evt)
		return
	}
	log.Printf("Slash command received: %+v", cmd)
	service.socketClient.Ack(*evt.Request)
}
