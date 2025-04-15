package psgrep

import (
	"CRUD-HOME-APPLIANCE-STORE/internal/logger"
	"CRUD-HOME-APPLIANCE-STORE/internal/model/domain"
	"context"
	"errors"
	"fmt"
	"runtime"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

var (
	ErrClientNotFound = errors.New("client not found")
	ErrQueryExection  = errors.New("query execution error")
)

type Write interface {
	Create(ctx context.Context, client domain.Client) error
	UpdateAddress(ctx context.Context, id, address uuid.UUID) error
	Delete(ctx context.Context, id uuid.UUID) error
}

type clientRepo struct {
	*basePostgresRepository
}

func NewClientRepository(conn *pgx.Conn, log *logger.Logger) *clientRepo {
	baseRepo := newBasePostgresRepository(conn, log)
	log.Debug("Client repo is created")
	return &clientRepo{baseRepo}
}

func (r *clientRepo) Create(ctx context.Context, client domain.Client) error {
	op := "repositories.postgres.clientRepository.Create"
	sqlStatement := "INSERT INTO client(name, surname, birthday, gender, address_id) VALUES (@clientName, @clientSurname, @clientBirthday, @clientGender, @clientAddressId);"
	args := pgx.NamedArgs{
		"clientName":      client.Name,
		"clientSurname":   client.Surname,
		"clientBirthday":  client.Birthday,
		"clientGender":    client.Gender,
		"clientAddressId": client.AddressId,
	}

	r.logger.Debug("Birthday", "Time", client.Birthday)

	_, err := r.db.Exec(ctx, sqlStatement, args)
	if err != nil {
		r.logger.Debug("Failed to create Client", "op", op)
		return fmt.Errorf("Client Repository: unable to insert row: %v", err)
	}

	return nil
}

func (r *clientRepo) GetAll(ctx context.Context, limit, offset int) ([]domain.Client, error) {
	op := "repositories.postgres.clientRepository.GetAll"
	sqlStatement := "SELECT * FROM client LIMIT @limit OFFSET @offset"
	args := pgx.NamedArgs{
		"limit":  limit,
		"offset": offset,
	}

	rows, err := r.db.Query(ctx, sqlStatement, args)
	if err != nil {
		r.logger.Debug("Unable to query client: %v", logger.Err(err), "op", op)
		return nil, fmt.Errorf("Client Repository: Error to get all client: %v", err)
	}
	defer rows.Close()

	var clients []domain.Client

	for rows.Next() {
		var client domain.Client

		if err := rows.Scan(&client.Id, &client.Name, &client.Surname,
			&client.Birthday, &client.Gender, &client.RegistrationDate,
			&client.AddressId); err != nil {
			return nil, fmt.Errorf("Client Repository: Failed to bind data in GetAll: %v", err)
		}

		clients = append(clients, client)
	}

	if len(clients) == 0 {
		r.logger.Debug("Clients not found", "op", op)
		return nil, ErrClientNotFound
	}

	r.logger.Debug("Client taked", "Clients:", clients, "op", op)
	return clients, nil
}

func (r *clientRepo) GetByNameAndSurname(ctx context.Context, name, surname string) ([]domain.Client, error) {
	op := "repositories.postgres.clientRepository.GetByNameAndSurname"
	sqlStatement := "SELECT * FROM client WHERE name = @clientName AND surname = @clientSurname"
	args := pgx.NamedArgs{
		"clientName":    name,
		"clientSurname": surname,
	}

	rows, err := r.db.Query(ctx, sqlStatement, args)

	if err != nil {
		r.logger.Debug("Failed get clients by name and surname", logger.Err(err), "op", op)
		return nil, fmt.Errorf("Client Repository: unable to query client by name and surname: %v", err)
	}
	defer rows.Close()

	var clients []domain.Client

	for rows.Next() {
		var client domain.Client

		if err := rows.Scan(&client.Id, &client.Name, &client.Surname,
			&client.Birthday, &client.Gender,
			&client.RegistrationDate, &client.AddressId); err != nil {
			r.logger.Debug("Failed binding data", logger.Err(err), "op", op)
			return nil, fmt.Errorf("Client Repository: failed to bind data in GeyBy...: %v", err)
		}

		clients = append(clients, client)
	}

	if len(clients) == 0 {
		r.logger.Debug("Clients not found", "op", op)
		return nil, ErrClientNotFound
	}

	r.logger.Debug("All data retrieved", "op", op)
	return clients, nil
}

func (r *clientRepo) UpdateAddress(ctx context.Context, id, address uuid.UUID) error {
	op := "repositories.postgres.clientRepository.Update"
	sqlStatement := "UPDATE client SET address_id=@addressId WHERE id=@id"
	arg := pgx.NamedArgs{
		"id":        id,
		"addressId": address,
	}

	r.logger.Debug("Parameter", "id", id, "address id", address)

	_, err := r.db.Exec(ctx, sqlStatement, arg)
	if errors.Is(err, pgx.ErrNoRows) {
		r.logger.Debug("Client not found", "op", op)
		return ErrClientNotFound
	}

	if err != nil {
		r.logger.Debug("Failed execution update query", logger.Err(err), "op", op)
		return fmt.Errorf("Client Repository: failed exec update query %v", err)
	}

	r.logger.Debug("Data updated", "op", op)
	return nil
}

func (r *clientRepo) Delete(ctx context.Context, id uuid.UUID) error {
	op := "repositories.postgres.clientRepository.Delete"
	sqlStatement := "DELETE FROM client WHERE id=@id"
	arg := pgx.NamedArgs{
		"id": id,
	}

	_, err := r.db.Exec(ctx, sqlStatement, arg)
	if err != nil {
		r.logger.Debug("Error in exec delete request to data base", "op", op)
		return fmt.Errorf("Client Repository: Failed delete client: %v", err)
	}

	r.logger.Debug("Client is deleted", "op", op)
	return nil
}

func (r *clientRepo) UnitOfWork(ctx context.Context, fn func(Write) error) error {
	op := "repositories.postgres.clientRepository.UnitOfWork"
	trx, err := r.db.Begin(ctx)
	if err != nil {
		r.logger.Debug("Transaction begin error", logger.Err(err), "op", op)
		return fmt.Errorf("Client unit of work %v", err)
	}

	defer func(ctx context.Context) {
		if p := recover(); p != nil {
			r.logger.Error("Panic catch in unit of work, rolling back", "op", op)
			_ = trx.Rollback(ctx)

			switch e := p.(type) {
			case runtime.Error:
				r.logger.Error("Runtime panic", "error", e, "op", op)
				panic(e)
			case error:
				err = fmt.Errorf("panic err: %v", err)
				return
			default:
				r.logger.Error("Not handle error", "error", e, "op", op)
				panic(e)
			}
		}

		if err != nil {
			r.logger.Error("Transaction error!", logger.Err(err), "op", op)
			_ = trx.Rollback(ctx)
		} else {
			r.logger.Debug("Transaction commited", "op", op)
			_ = trx.Commit(ctx)
		}
	}(ctx)

	newStore := NewClientRepository(r.db, r.logger)
	return fn(newStore)
}
