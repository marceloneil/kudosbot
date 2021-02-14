package repository

import (
	"kudosbot/model"
	"kudosbot/table"
	"time"

	"github.com/google/uuid"

	"github.com/jmoiron/sqlx"
	"github.com/lumina-tech/gooq/pkg/gooq"
)

const DayDuration = time.Hour * 24
const WeekDuration = time.Hour * 24 * 7
const MonthDuration = time.Hour * 24 * 30

type Leaderboard struct {
	Count      int    `db:"count"`
	ReceiverID string `db:"receiver_id"`
}

type KudosRepository struct {
	db *sqlx.DB
}

func NewKudosRepository(
	db *sqlx.DB,
) *KudosRepository {
	return &KudosRepository{
		db: db,
	}
}

func (repo *KudosRepository) UpsertKudos(
	timestamp time.Time,
	senderID string, receiverID string, eventID string, reaction string,
) (*model.Kudos, error) {
	stmt := gooq.InsertInto(table.Kudos).
		Set(table.Kudos.ID, uuid.New()).
		Set(table.Kudos.Timestamp, timestamp).
		Set(table.Kudos.SenderID, senderID).
		Set(table.Kudos.ReceiverID, receiverID).
		Set(table.Kudos.EventID, eventID).
		Set(table.Kudos.Reaction, reaction).
		OnConflictDoNothing().
		Returning(table.Kudos.Asterisk)
	return table.Kudos.ScanRow(repo.db, stmt)
}

func (repo *KudosRepository) CountKudosGivenForToday(
	senderID string,
) (int, error) {
	todayDate := time.Now().Truncate(DayDuration)
	stmt := gooq.Select(gooq.Count(table.Kudos.Asterisk)).
		From(table.Kudos).Where(
		table.Kudos.Timestamp.IsGt(todayDate),
		table.Kudos.SenderID.IsEq(senderID))
	return gooq.ScanCount(repo.db, stmt)
}

func (repo *KudosRepository) GetLeaderboard(
	duration *time.Duration,
) ([]Leaderboard, error) {
	var conditions []gooq.Expression
	if duration != nil {
		switch *duration {
		case DayDuration:
			date := time.Now().Truncate(DayDuration)
			conditions = append(conditions, table.Kudos.Timestamp.IsGt(date))
		case WeekDuration:
			date := time.Now().Truncate(WeekDuration)
			conditions = append(conditions, table.Kudos.Timestamp.IsGt(date))
		case MonthDuration:
			date := time.Now().Truncate(MonthDuration)
			conditions = append(conditions, table.Kudos.Timestamp.IsGt(date))
		}
	}
	stmt := gooq.Select(
		gooq.Count(table.Kudos.Asterisk).As("count"),
		table.Kudos.ReceiverID,
	).From(table.Kudos).Where(conditions...).
		GroupBy(table.Kudos.ReceiverID).
		OrderBy(gooq.String("count").Desc()).
		Limit(10)

	var results []Leaderboard
	if err := gooq.ScanRows(repo.db, stmt, &results); err != nil {
		return nil, err
	}
	return results, nil
}
