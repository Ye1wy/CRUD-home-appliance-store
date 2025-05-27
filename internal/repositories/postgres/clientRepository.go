package postgres

import (
	crud_errors "CRUD-HOME-APPLIANCE-STORE/internal/errors"
	"CRUD-HOME-APPLIANCE-STORE/internal/model/domain"
	"CRUD-HOME-APPLIANCE-STORE/pkg/logger"
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

type ClientRepo struct {
	*basePostgresRepository
}

func NewClientRepository(db DB, log *logger.Logger) *ClientRepo {
	baseRepo := newBasePostgresRepository(db, log)
	log.Debug("client repo is created")
	return &ClientRepo{baseRepo}
}

func (r *ClientRepo) Create(ctx context.Context, client domain.Client) error {
	op := "repositories.postgres.clientRepository.Create"
	sqlStatement := "INSERT INTO client(name, surname, birthday, gender, address_id) VALUES (@clientName, @clientSurname, @clientBirthday, @clientGender, @clientAddressId);"
	args := pgx.NamedArgs{
		"clientName":      client.Name,
		"clientSurname":   client.Surname,
		"clientBirthday":  client.Birthday,
		"clientGender":    client.Gender,
		"clientAddressId": client.Address.Id,
	}

	_, err := r.db.Exec(ctx, sqlStatement, args)
	if err != nil {
		r.logger.Debug("failed to create Client", logger.Err(err), "op", op)
		return fmt.Errorf("%s: unable to insert row: %v", op, err)
	}

	return nil
}

func (r *ClientRepo) GetAll(ctx context.Context, limit, offset int) ([]domain.Client, error) {
	op := "repositories.postgres.clientRepository.GetAll"
	sqlStatement := `SELECT 
		c.id AS client_id,
		c.name,
		c.surname,
		c.birthday,
		c.gender,
		c.registration_date,
		a.id AS address_id,
		a.country,
		a.city,
		a.street
		FROM client c
		LEFT JOIN address a ON c.address_id = a.id
		LIMIT @limit OFFSET @offset;`
	args := pgx.NamedArgs{
		"limit":  limit,
		"offset": offset,
	}

	rows, err := r.db.Query(ctx, sqlStatement, args)
	if err != nil {
		r.logger.Debug("unable to query client: %v", logger.Err(err), "op", op)
		return nil, fmt.Errorf("%s: query error: %v", op, err)
	}
	defer rows.Close()

	var clients []domain.Client

	for rows.Next() {
		var client domain.Client

		err := rows.Scan(
			&client.Id,
			&client.Name,
			&client.Surname,
			&client.Birthday,
			&client.Gender,
			&client.RegistrationDate,
			&client.Address.Id,
			&client.Address.Country,
			&client.Address.City,
			&client.Address.Street,
		)
		if err != nil {
			r.logger.Debug("failed binding data", logger.Err(err), "op", op)
			return nil, fmt.Errorf("%s: failed to bind data: %v", op, err)
		}

		clients = append(clients, client)
	}

	if len(clients) == 0 {
		r.logger.Debug("clients not found", "op", op)
		return nil, fmt.Errorf("%s: %w", op, crud_errors.ErrNotFound)
	}

	return clients, nil
}

func (r *ClientRepo) GetByNameAndSurname(ctx context.Context, name, surname string) ([]domain.Client, error) {
	op := "repositories.postgres.clientRepository.GetByNameAndSurname"
	sqlStatement := `SELECT 
		c.id,
		c.name,
		c.surname,
		c.birthday,
		c.gender,
		c.registration_date,
		a.id,
		a.country,
		a.city,
		a.street
		FROM client c
		LEFT JOIN address a ON c.address_id = a.id
		WHERE c.name = @clientName AND c.surname = @clientSurname;`
	args := pgx.NamedArgs{
		"clientName":    name,
		"clientSurname": surname,
	}

	rows, err := r.db.Query(ctx, sqlStatement, args)
	if err != nil {
		r.logger.Debug("failed get clients by name and surname", logger.Err(err), "op", op)
		return nil, fmt.Errorf("%s: query error: %v", op, err)
	}
	defer rows.Close()

	var clients []domain.Client

	for rows.Next() {
		var client domain.Client

		err := rows.Scan(
			&client.Id,
			&client.Name,
			&client.Surname,
			&client.Birthday,
			&client.Gender,
			&client.RegistrationDate,
			&client.Address.Id,
			&client.Address.Country,
			&client.Address.City,
			&client.Address.Street,
		)
		if err != nil {
			r.logger.Debug("failed binding data", logger.Err(err), "op", op)
			return nil, fmt.Errorf("%s: failed to bind data: %v", op, err)
		}

		clients = append(clients, client)
	}

	if len(clients) == 0 {
		r.logger.Debug("clients not found", "op", op)
		return nil, fmt.Errorf("%s: %w", op, crud_errors.ErrNotFound)
	}

	return clients, nil
}

func (r *ClientRepo) GetById(ctx context.Context, id uuid.UUID) (*domain.Client, error) {
	op := "repositories.postgres.clientRepository.GetById"
	sqlStatement := `SELECT
	c.id,
	c.name,
	c.surname,
	c.birthday,
	c.gender,
	a.id,
	a.country,
	a.city,
	a.street
	FROM client c
	LEFT JOIN address a ON c.address_id = a.id
	WHERE c.id = @id;`
	arg := pgx.NamedArgs{"id": id}

	var client domain.Client
	row := r.db.QueryRow(ctx, sqlStatement, arg)
	err := row.Scan(
		&client.Id,
		&client.Name,
		&client.Surname,
		&client.Birthday,
		&client.Gender,
		&client.Address.Id,
		&client.Address.Country,
		&client.Address.City,
		&client.Address.Street,
	)

	if errors.Is(err, pgx.ErrNoRows) {
		r.logger.Debug("client not found", "op", op)
		return nil, fmt.Errorf("%s: %w", op, crud_errors.ErrNotFound)
	}

	if err != nil {
		r.logger.Debug("scan unable", logger.Err(err), "op", op)
		return nil, fmt.Errorf("%s: %v", op, err)
	}

	return &client, nil
}

func (r *ClientRepo) UpdateAddress(ctx context.Context, id, address uuid.UUID) error {
	op := "repositories.postgres.clientRepository.Update"
	sqlStatement := "UPDATE client SET address_id=@addressId WHERE id=@id"
	arg := pgx.NamedArgs{
		"id":        id,
		"addressId": address,
	}

	tag, err := r.db.Exec(ctx, sqlStatement, arg)
	if tag.RowsAffected() == 0 {
		r.logger.Debug("client not found", "op", op)
		return fmt.Errorf("%s: %w", op, crud_errors.ErrNotFound)
	}

	if err != nil {
		r.logger.Debug("failed execution update query", logger.Err(err), "op", op)
		return fmt.Errorf("%s: failed exec query: %v", op, err)
	}

	return nil
}

func (r *ClientRepo) Delete(ctx context.Context, id uuid.UUID) error {
	op := "repositories.postgres.clientRepository.Delete"
	sqlStatement := "DELETE FROM client WHERE id=@id"
	arg := pgx.NamedArgs{
		"id": id,
	}

	_, err := r.db.Exec(ctx, sqlStatement, arg)
	if err != nil {
		r.logger.Debug("error in exec delete request to data base", "op", op)
		return fmt.Errorf("%s: %v", op, err)
	}

	return nil
}
