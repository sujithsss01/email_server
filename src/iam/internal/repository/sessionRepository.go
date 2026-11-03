package repository

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/sujithsss01/email_server/iam/internal/apperrors"
	"github.com/sujithsss01/email_server/iam/internal/model"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type SessionRepository struct {
	db *gorm.DB
}

func (r *SessionRepository) StoreRefreshToken(ctx context.Context, userId uuid.UUID, existingSessionId *uuid.UUID, refreshToken string) (string, error) {
	var sessionId uuid.UUID
	var err error
	var emptyStr string
	if existingSessionId != nil {
		sessionId = *existingSessionId
	} else {
		sessionId, err = uuid.NewV7()
		if err != nil {
			return emptyStr, apperrors.ErrUnexpectedInternalFailure
		}
	}
	return r.storeRefreshToken(
		r.db,
		ctx,
		userId,
		sessionId,
		refreshToken,
	)
}
func (r *SessionRepository) storeRefreshToken(db *gorm.DB, ctx context.Context, userId uuid.UUID, sessionId uuid.UUID, refreshToken string) (string, error) {
	var emptyStr string
	id, err := uuid.NewV7()
	if err != nil {
		return emptyStr, apperrors.ErrUnexpectedInternalFailure
	}
	now := time.Now()
	session := &model.Session{
		ID:           id,
		SessionId:    sessionId,
		UserId:       userId,
		RefreshToken: refreshToken,
		CreatedAt:    now,
		ExpiresAt:    now.Add(72 * time.Hour),
	}
	err = gorm.G[model.Session](db).Create(ctx, session)
	if err != nil {
		return emptyStr, err
	}
	return sessionId.String(), nil
}

func (r *SessionRepository) RotateRefreshToken(ctx context.Context, sessionId uuid.UUID, oldRefreshToken string, newRefreshToken string) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		var session model.Session
		err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).WithContext(ctx).Where("session_id=? AND refresh_token=?", sessionId, oldRefreshToken).First(&session).Error
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return apperrors.ErrInvalidCredentials
			}
			return err
		}
		if session.RevokedAt != nil || !session.ExpiresAt.After(time.Now()) {
			return apperrors.ErrInvalidCredentials
		}
		err = r.revokeRefreshToken(tx, ctx, oldRefreshToken)
		if err != nil {
			return err
		}
		_, err = r.storeRefreshToken(tx, ctx, session.UserId, session.SessionId, newRefreshToken)
		return err
	})
}

func (r *SessionRepository) RevokeRefreshToken(ctx context.Context, refreshToken string) error {
	return r.revokeRefreshToken(r.db, ctx, refreshToken)
}

func (r *SessionRepository) revokeRefreshToken(db *gorm.DB, ctx context.Context, refreshToken string) error {
	result := db.WithContext(ctx).Model(&model.Session{}).
		Where("refresh_token=?", refreshToken).
		Update("revoked_at", time.Now())
	if err := result.Error; err != nil {
		return err
	}
	if result.RowsAffected != 1 {
		return apperrors.ErrInvalidCredentials
	}
	return nil
}

func (r *SessionRepository) RevokeSession(ctx context.Context, sessionId uuid.UUID) error {
	return r.revokeSession(r.db, ctx, sessionId)
}

func (r *SessionRepository) revokeSession(db *gorm.DB, ctx context.Context, sessionId uuid.UUID) error {
	result := db.WithContext(ctx).Model(&model.Session{}).
		Where("session_id=?", sessionId).
		Update("revoked_at", time.Now())
	return result.Error
}

func (r *SessionRepository) RevokeTokensByUser(ctx context.Context, userId uuid.UUID) error {
	return r.revokeTokenByUser(r.db, ctx, userId)
}

func (r *SessionRepository) revokeTokenByUser(db *gorm.DB, ctx context.Context, userId uuid.UUID) error {
	result := db.WithContext(ctx).Model(&model.Session{}).
		Where("user_id=?", userId).
		Update("revoked_at", time.Now())
	return result.Error
}

func (r *SessionRepository) FetchSessionDetailsByToken(ctx context.Context, userId uuid.UUID, refreshToken string) (*model.Session, error) {
	var session model.Session
	result := r.db.WithContext(ctx).Where("refresh_token=? AND user_id = ?", refreshToken, userId).First(&session)
	if result.Error != nil {
		return nil, result.Error
	}
	return &session, nil
}

func NewSessionRepository(db *gorm.DB) *SessionRepository {
	return &SessionRepository{
		db: db,
	}
}
