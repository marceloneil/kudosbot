package table

import (
	"kudosbot/model"

	"github.com/lumina-tech/gooq/pkg/gooq"
	"gopkg.in/guregu/null.v3"
)

type kudosConstraints struct {
	KudosPkey                         gooq.DatabaseConstraint
	KudosSenderIDReceiverIDEventIDKey gooq.DatabaseConstraint
}

type kudos struct {
	gooq.TableImpl
	Asterisk   gooq.StringField
	ID         gooq.UUIDField
	Timestamp  gooq.TimeField
	SenderID   gooq.StringField
	ReceiverID gooq.StringField
	EventID    gooq.StringField
	Reaction   gooq.StringField

	Constraints *kudosConstraints
}

func newKudosConstraints() *kudosConstraints {
	constraints := &kudosConstraints{}
	constraints.KudosPkey = gooq.DatabaseConstraint{
		Name:      "kudos_pkey",
		Predicate: null.NewString("", false),
	}
	constraints.KudosSenderIDReceiverIDEventIDKey = gooq.DatabaseConstraint{
		Name:      "kudos_sender_id_receiver_id_event_id_key",
		Predicate: null.NewString("", false),
	}
	return constraints
}

func newKudos() *kudos {
	instance := &kudos{}
	instance.Initialize("public", "kudos")
	instance.Asterisk = gooq.NewStringField(instance, "*")
	instance.ID = gooq.NewUUIDField(instance, "id")
	instance.Timestamp = gooq.NewTimeField(instance, "timestamp")
	instance.SenderID = gooq.NewStringField(instance, "sender_id")
	instance.ReceiverID = gooq.NewStringField(instance, "receiver_id")
	instance.EventID = gooq.NewStringField(instance, "event_id")
	instance.Reaction = gooq.NewStringField(instance, "reaction")
	instance.Constraints = newKudosConstraints()
	return instance
}

func (t *kudos) As(alias string) *kudos {
	instance := newKudos()
	instance.TableImpl = *instance.TableImpl.As(alias)
	return instance
}

func (t *kudos) GetColumns() []gooq.Expression {
	return []gooq.Expression{
		t.ID,
		t.Timestamp,
		t.SenderID,
		t.ReceiverID,
		t.EventID,
		t.Reaction,
	}
}

func (t *kudos) ScanRow(
	db gooq.DBInterface, stmt gooq.Fetchable,
) (*model.Kudos, error) {
	result := model.Kudos{}
	if err := gooq.ScanRow(db, stmt, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

func (t *kudos) ScanRows(
	db gooq.DBInterface, stmt gooq.Fetchable,
) ([]model.Kudos, error) {
	results := []model.Kudos{}
	if err := gooq.ScanRows(db, stmt, &results); err != nil {
		return nil, err
	}
	return results, nil
}

var Kudos = newKudos()
