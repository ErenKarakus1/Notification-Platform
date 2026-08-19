package repository

import (
	"context"
	"errors"

	"github.com/ErenKarakus1/Notification-Platform/notification-service/internal/model"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

var ErrNotificationNotFound = errors.New("notification not found")

const createNotificationQuery = `
	INSERT INTO notifications (
		id,
		customer_id,
		recipient,
		channel,
		subject,
		body,
		status
	)
	VALUES ($1,$2,$3,$4,$5,$6,$7)
`

const getNotificationByIDQuery = `
	SELECT 
		id,
		customer_id,
		recipient,
		channel,
		subject,
		body,
		status,
		created_at
	FROM notifications
	WHERE id=$1
	AND customer_id=$2
`

const updateNotificationStatusQuery = `
	UPDATE notifications
	SET status=$1
	WHERE id=$2
`

func CreateNotification(ctx context.Context, pool *pgxpool.Pool, notification model.Notification) error {
	_, err := pool.Exec(
		ctx,
		createNotificationQuery,
		notification.ID,
		notification.CustomerID,
		notification.Recipient,
		notification.Channel,
		notification.Subject,
		notification.Body,
		notification.Status,
	)
	if err != nil {
		return errors.New("internal server error")
	}
	return nil
}

func GetNotificationByID(ctx context.Context, pool *pgxpool.Pool, notificationID uuid.UUID, customerID uuid.UUID) (model.Notification, error) {
	var notification model.Notification
	err := pool.QueryRow(
		ctx,
		getNotificationByIDQuery,
		notificationID,
		customerID,
	).Scan(
		&notification.ID,
		&notification.CustomerID,
		&notification.Recipient,
		&notification.Channel,
		&notification.Subject,
		&notification.Body,
		&notification.Status,
		&notification.CreatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return model.Notification{}, ErrNotificationNotFound
		}
		return model.Notification{}, errors.New("internal server error")
	}
	return notification, nil
}

func UpdateNotificationStatus(ctx context.Context, pool *pgxpool.Pool, notificationID uuid.UUID, status string) error {
	result, err := pool.Exec(
		ctx,
		updateNotificationStatusQuery,
		status,
		notificationID,
	)
	if err != nil {
		return errors.New("internal server error")
	}

	if result.RowsAffected() == 0 {
		return ErrNotificationNotFound
	}
	return nil
}
