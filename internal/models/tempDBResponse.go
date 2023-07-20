package models

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

func NewResponseTemp(pID pgtype.UUID, pName, pGender, pBirthdate pgtype.Text,
	aID pgtype.UUID, aName, aGender, aBirthdate, aRole pgtype.Text) *ResponseTemp {
	return &ResponseTemp{ProducerID: pID,
		ProducerName:      pName,
		ProducerGender:    pGender,
		ProducerBirthdate: pBirthdate,
		ActorID:           aID,
		ActorName:         aName,
		ActorGender:       aGender,
		ActorRole:         aRole,
		ActorBirthdate:    aBirthdate,
	}
}
