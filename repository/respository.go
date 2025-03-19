package repository

import (
	"database/sql"
	"time"

	"github.com/arjunsaxaena/Moonrider-Assignment/model"
	"github.com/huandu/go-sqlbuilder"
)

func Create(db *sql.DB, contact *model.Contact) error {
	email := ""
	phoneNumber := ""
	if contact.Email != nil {
		email = *contact.Email
	}
	if contact.PhoneNumber != nil {
		phoneNumber = *contact.PhoneNumber
	}

	existingContacts, err := Get(db, email, phoneNumber)
	if err != nil {
		return err
	}

	var primaryContactID *string
	for _, existing := range existingContacts {
		if existing.LinkPrecedence == "primary" {
			primaryContactID = &existing.ID
			break
		}
	}

	if primaryContactID != nil {
		phoneExists := false
		emailExists := false
		for _, existing := range existingContacts {
			if existing.PhoneNumber != nil && *existing.PhoneNumber == phoneNumber {
				phoneExists = true
			}
			if existing.Email != nil && *existing.Email == email {
				emailExists = true
			}
		}

		if !phoneExists {
			contact.LinkedID = primaryContactID
			contact.LinkPrecedence = "secondary"
		} else if !emailExists {
			contact.LinkedID = primaryContactID
			contact.LinkPrecedence = "secondary"
		} else {
			return nil
		}
	} else {
		contact.LinkPrecedence = "primary"
	}

	sb := sqlbuilder.PostgreSQL.NewInsertBuilder()
	sb.InsertInto("contacts")
	sb.Cols("email", "phone_number", "linked_id", "link_precedence", "created_at", "updated_at")
	sb.Values(email, phoneNumber, contact.LinkedID, contact.LinkPrecedence, time.Now(), time.Now())

	query, args := sb.Build()

	return db.QueryRow(query+" RETURNING id", args...).Scan(&contact.ID)
}

func Get(db *sql.DB, email, phoneNumber string) ([]model.Contact, error) {
	sb := sqlbuilder.PostgreSQL.NewSelectBuilder()

	sb.Select("*").
		From("contacts").
		Where(
			sb.Or(
				sb.Equal("email", email),
				sb.Equal("phone_number", phoneNumber),
			),
			sb.IsNull("deleted_at"),
		)

	query, args := sb.Build()

	rows, err := db.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var contacts []model.Contact
	var primaryContactID *string

	for rows.Next() {
		var contact model.Contact
		err := rows.Scan(
			&contact.ID, &contact.Email, &contact.PhoneNumber,
			&contact.LinkedID, &contact.LinkPrecedence,
			&contact.CreatedAt, &contact.UpdatedAt, &contact.DeletedAt,
		)
		if err != nil {
			return nil, err
		}

		contacts = append(contacts, contact)

		if contact.LinkPrecedence == "primary" {
			primaryContactID = &contact.ID
		}
		if contact.LinkedID != nil && primaryContactID == nil {
			primaryContactID = contact.LinkedID
		}
	}

	if primaryContactID == nil {
		return contacts, nil
	}

	sb = sqlbuilder.PostgreSQL.NewSelectBuilder()
	sb.Select("*").
		From("contacts").
		Where(
			sb.Or(
				sb.Equal("id", *primaryContactID),
				sb.Equal("linked_id", *primaryContactID),
			),
			sb.IsNull("deleted_at"),
		)

	query, args = sb.Build()
	rows, err = db.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var allContacts []model.Contact
	for rows.Next() {
		var contact model.Contact
		err := rows.Scan(
			&contact.ID, &contact.Email, &contact.PhoneNumber,
			&contact.LinkedID, &contact.LinkPrecedence,
			&contact.CreatedAt, &contact.UpdatedAt, &contact.DeletedAt,
		)
		if err != nil {
			return nil, err
		}
		allContacts = append(allContacts, contact)
	}

	return allContacts, nil
}

func Update(db *sql.DB, contact *model.Contact) error {
	sb := sqlbuilder.NewUpdateBuilder()

	sb.Update("contacts")
	sb.Set(
		sb.Assign("email", contact.Email),
		sb.Assign("phone_number", contact.PhoneNumber),
		sb.Assign("linked_id", contact.LinkedID),
		sb.Assign("link_precedence", contact.LinkPrecedence),
		sb.Assign("updated_at", time.Now()),
	)
	sb.Where(sb.Equal("id", contact.ID))

	query, args := sb.Build()

	_, err := db.Exec(query, args...)
	return err
}

func Delete(db *sql.DB, id string) error {
	sb := sqlbuilder.NewUpdateBuilder()

	sb.Update("contacts")
	sb.Set(sb.Assign("deleted_at", time.Now()))
	sb.Where(sb.Equal("id", id))

	query, args := sb.Build()

	_, err := db.Exec(query, args...)
	return err
}
