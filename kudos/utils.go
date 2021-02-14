package kudos

import (
	"fmt"

	"github.com/slack-go/slack"
)

func (service *Service) getUsers() error {
	users, err := service.apiClient.GetUsers()
	if err != nil {
		return err
	}
	usersMap := make(map[string]*slack.User)
	for index := range users {
		user := &users[index]
		usersMap[user.ID] = user
	}
	service.usersMap = usersMap
	return nil
}

func (service *Service) getUser(
	userID string,
) (*slack.User, error) {
	user, ok := service.usersMap[userID]
	if ok {
		return user, nil
	}
	// update user map, probably should lock but this is unlikely
	if err := service.getUsers(); err != nil {
		return nil, err
	}
	user, ok = service.usersMap[userID]
	if ok {
		return user, nil
	}
	return nil, fmt.Errorf("cannot find user id=%s", userID)
}
