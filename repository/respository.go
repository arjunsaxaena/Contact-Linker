package repository

import (
	"database/sql"
	"time"

	"github.com/arjunsaxaena/Moonrider-Assignment/model"
	"github.com/huandu/go-sqlbuilder"
)

func Create(db *sql.DB, contact *model.Contact) error {
	sb := sqlbuilder.PostgreSQL.NewInsertBuilder()
	sb.InsertInto("contacts")
	sb.Cols("email", "phone_number", "linked_id", "link_precedence", "created_at", "updated_at")
	sb.Values(contact.Email, contact.PhoneNumber, contact.LinkedID, contact.LinkPrecedence, time.Now(), time.Now())

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
	}
	return contacts, nil
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
