package kudos

import (
	"fmt"
	"os"

	"golang.org/x/net/context"

	"github.com/slack-go/slack/socketmode"
)

func (service *Service) start(
	ctx context.Context,
) error {
	if err := service.getUsers(); err != nil {
		return err
	}
	return nil
}

func (service *Service) Run() {
	go func() {
		for evt := range service.socketClient.Events {
			switch evt.Type {
			case socketmode.EventTypeConnecting:
				fmt.Println("Connecting to Slack with Socket Mode...")
			case socketmode.EventTypeConnectionError:
				fmt.Println("Connection failed. Retrying later...")
			case socketmode.EventTypeConnected:
				fmt.Println("Connected to Slack with Socket Mode.")
			case socketmode.EventTypeEventsAPI:
				service.handleEventAPIEvents(&evt)
			case socketmode.EventTypeInteractive:
				service.handleInteractionEvents(&evt)
			case socketmode.EventTypeSlashCommand:
				service.handleSlashCommandEvents(&evt)
			case socketmode.EventTypeHello:
			default:
				fmt.Fprintf(os.Stderr, "Unexpected event type received: %s\n", evt.Type)
			}
		}
	}()
	service.socketClient.Run()
}
