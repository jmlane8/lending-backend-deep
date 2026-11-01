package pg

import (
	"context"
	"errors"
	"fmt"

	"lending-backend-deep/internal/validation"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

// CustomerStore handles persistence for the customer and stored_photo tables.
// Every method runs inside the caller's transaction.
type CustomerStore struct{}

// NewCustomerStore returns a new CustomerStore.
func NewCustomerStore() *CustomerStore {
	return &CustomerStore{}
}

// InsertPhotos stores the given photos and returns their IDs in the same order.
func (s *CustomerStore) InsertPhotos(ctx context.Context, tx pgx.Tx, photos ...validation.Photo) ([]uuid.UUID, error) {
	ids := make([]uuid.UUID, 0, len(photos))
	for _, photo := range photos {
		id, err := uuid.NewV7()
		if err != nil {
			return nil, fmt.Errorf("CustomerStore.InsertPhotos: id: %w", err)
		}
		if _, err := tx.Exec(ctx, `
			INSERT INTO stored_photo (photo_id, content_type, content)
			VALUES ($1, $2, $3)
		`, id, photo.ContentType, photo.Content); err != nil {
			return nil, fmt.Errorf("CustomerStore.InsertPhotos: %w", err)
		}
		ids = append(ids, id)
	}
	return ids, nil
}

// UpsertCustomer inserts the customer, or updates the row with the same national
// ID, and returns its ID.  A customer's previous photos are soft-deleted; their
// content stays until retention rules are defined.
func (s *CustomerStore) UpsertCustomer(ctx context.Context, tx pgx.Tx, c validation.Customer, idPhoto, selfiePhoto uuid.UUID) (uuid.UUID, error) {
	var previousIDPhoto, previousSelfie uuid.UUID
	err := tx.QueryRow(ctx, `
		SELECT government_id_photo, selfie_photo
		FROM customer
		WHERE national_id_number = $1
		FOR UPDATE
	`, c.NationalIDNumber).Scan(&previousIDPhoto, &previousSelfie)
	replacing := err == nil
	if err != nil && !errors.Is(err, pgx.ErrNoRows) {
		return uuid.Nil, fmt.Errorf("CustomerStore.UpsertCustomer: previous photos: %w", err)
	}

	newID, err := uuid.NewV7()
	if err != nil {
		return uuid.Nil, fmt.Errorf("CustomerStore.UpsertCustomer: id: %w", err)
	}

	var customerID uuid.UUID
	err = tx.QueryRow(ctx, `
		INSERT INTO customer (
			customer_id, national_id_number, full_legal_name, government_id_photo,
			government_id_expiry_date, selfie_photo, phone_number, gender,
			date_of_birth, place_of_birth, marital_status, occupation,
			work_duration_years, monthly_income, number_of_dependents, email,
			domicile_street, domicile_district, domicile_city, domicile_province
		) VALUES (
			$1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13,
			$14::text::numeric, $15, $16, $17, $18, $19, $20
		)
		ON CONFLICT (national_id_number) DO UPDATE
		SET full_legal_name = EXCLUDED.full_legal_name,
		    government_id_photo = EXCLUDED.government_id_photo,
		    government_id_expiry_date = EXCLUDED.government_id_expiry_date,
		    selfie_photo = EXCLUDED.selfie_photo,
		    phone_number = EXCLUDED.phone_number,
		    gender = EXCLUDED.gender,
		    date_of_birth = EXCLUDED.date_of_birth,
		    place_of_birth = EXCLUDED.place_of_birth,
		    marital_status = EXCLUDED.marital_status,
		    occupation = EXCLUDED.occupation,
		    work_duration_years = EXCLUDED.work_duration_years,
		    monthly_income = EXCLUDED.monthly_income,
		    number_of_dependents = EXCLUDED.number_of_dependents,
		    email = EXCLUDED.email,
		    domicile_street = EXCLUDED.domicile_street,
		    domicile_district = EXCLUDED.domicile_district,
		    domicile_city = EXCLUDED.domicile_city,
		    domicile_province = EXCLUDED.domicile_province,
		    updated_at = NOW()
		RETURNING customer_id
	`, newID, c.NationalIDNumber, c.FullLegalName, idPhoto,
		c.GovernmentIDExpiryDate, selfiePhoto, c.PhoneNumber, c.Gender,
		c.DateOfBirth, c.PlaceOfBirth, c.MaritalStatus, c.Occupation,
		c.WorkDurationYears, c.MonthlyIncome, c.NumberOfDependents, c.Email,
		c.DomicileStreet, c.DomicileDistrict, c.DomicileCity, c.DomicileProvince,
	).Scan(&customerID)
	if err != nil {
		return uuid.Nil, fmt.Errorf("CustomerStore.UpsertCustomer: %w", err)
	}

	if replacing {
		if _, err := tx.Exec(ctx, `
			UPDATE stored_photo
			SET deleted = TRUE, deleted_at = NOW(), deleted_by = 'system', updated_at = NOW()
			WHERE photo_id IN ($1, $2)
		`, previousIDPhoto, previousSelfie); err != nil {
			return uuid.Nil, fmt.Errorf("CustomerStore.UpsertCustomer: soft-delete photos: %w", err)
		}
	}
	return customerID, nil
}
