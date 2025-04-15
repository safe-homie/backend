package notify

import (
	"fmt"

	expo "github.com/oliveroneill/exponent-server-sdk-golang/sdk"
	"github.com/safe-homie/backend/internal/domain"
	"github.com/safe-homie/backend/internal/store"
)

type NotifyService interface {
	Notify(token string, ntf domain.Notification) error
	SaveToken(req *domain.SaveTokenRequest) (*store.Token, error)
}

type expoNotifer struct {
	store store.Store
}

func NewService(store store.Store) NotifyService {
	return &expoNotifer{store: store}
}

func (ex *expoNotifer) SaveToken(req *domain.SaveTokenRequest) (*store.Token, error) {
	save := store.Token{
		Token: req.Token,
	}
	tokenDB, err := ex.store.SaveToken(&save)
	if err != nil {
		return nil, err
	}
	return tokenDB, nil

}

func (ex *expoNotifer) Notify(token string, ntf domain.Notification) error {
	pushToken, err := expo.NewExponentPushToken(token)
	if err != nil {
		return nil
	}
	client := expo.NewPushClient(nil)
	response, err := client.Publish(&expo.PushMessage{
		To:       []expo.ExponentPushToken{pushToken},
		Body:     ntf.Body,
		Data:     ntf.Data,
		Sound:    "default",
		Title:    ntf.Title,
		Priority: expo.HighPriority,
	})
	if err != nil {
		return err
	}
	if response.ValidateResponse() != nil {
		fmt.Println(response.PushMessage.To, "failed")
	}
	return nil
}
