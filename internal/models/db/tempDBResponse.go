package db

import "github.com/jackc/pgtype"

type ResponseTemp struct {
	ProducerID        pgtype.UUID
	ProducerName      pgtype.Text
	ProducerGender    pgtype.Text
	ProducerBirthdate pgtype.Text
	ActorID           pgtype.UUID
	ActorName         pgtype.Text
	ActorGender       pgtype.Text
	ActorBirthdate    pgtype.Text
	ActorRole         pgtype.Text
}
