// Package user provides the PostgreSQL-backed write repository for the User aggregate.
package user

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"

	errs "github.com/ulbwa/medincident-command-service/internal/common/errors"
	"github.com/ulbwa/medincident-command-service/internal/model"
	"github.com/ulbwa/medincident-command-service/internal/repository/postgres"
	"github.com/ulbwa/medincident-command-service/internal/repository/user/entity"
	serviceuser "github.com/ulbwa/medincident-command-service/internal/service/user"
)

// Compile-time check that userRepository satisfies the service port.
var _ serviceuser.Repository = (*userRepository)(nil)

// querier abstracts the read operations shared by *sqlx.DB and *sqlx.Tx.
type querier interface {
	GetContext(ctx context.Context, dest interface{}, query string, args ...interface{}) error
	SelectContext(ctx context.Context, dest interface{}, query string, args ...interface{}) error
}

// userRepository is the PostgreSQL implementation of service/user.Repository.
type userRepository struct {
	pool *sqlx.DB
}

// NewRepository returns a user write repository backed by the given pool.
func NewRepository(pool *sqlx.DB) *userRepository {
	return &userRepository{pool: pool}
}

// querier returns the active transaction from ctx when present, otherwise the
// pool. This allows read operations to participate in an ongoing transaction
// without coupling the caller to the postgres package.
func (r *userRepository) db(ctx context.Context) querier {
	if tx := postgres.TxFromContext(ctx); tx != nil {
		return tx
	}

	return r.pool
}

// GetByID loads the User aggregate together with all its Employments.
// Returns errs.ErrNotFound when no user exists with that id.
func (r *userRepository) GetByID(ctx context.Context, id uuid.UUID) (*model.User, error) {
	db := r.db(ctx)

	var rec entity.UserRecord

	const userQuery = `SELECT * FROM users WHERE id = $1`
	if err := db.GetContext(ctx, &rec, userQuery, id); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errs.ErrNotFound
		}

		return nil, fmt.Errorf("query user by id: %w", err)
	}

	var empRecs []entity.EmploymentRecord

	const empQuery = `SELECT * FROM employments WHERE user_id = $1 ORDER BY assigned_at`
	if err := db.SelectContext(ctx, &empRecs, empQuery, id); err != nil {
		return nil, fmt.Errorf("query employments for user %s: %w", id, err)
	}

	user, err := entity.ToUser(rec, empRecs)
	if err != nil {
		return nil, fmt.Errorf("restore user %s from storage: %w", id, err)
	}

	return user, nil
}

// ExistsByID reports whether a user with the given id exists.
func (r *userRepository) ExistsByID(ctx context.Context, id uuid.UUID) (bool, error) {
	db := r.db(ctx)

	var exists bool

	const query = `SELECT EXISTS(SELECT 1 FROM users WHERE id = $1)`
	if err := db.GetContext(ctx, &exists, query, id); err != nil {
		return false, fmt.Errorf("check user exists by id: %w", err)
	}

	return exists, nil
}

// ExistsByIdentityID reports whether a user with the given identity provider ID exists.
func (r *userRepository) ExistsByIdentityID(ctx context.Context, identityID string) (bool, error) {
	db := r.db(ctx)

	var exists bool

	const query = `SELECT EXISTS(SELECT 1 FROM users WHERE identity_id = $1)`
	if err := db.GetContext(ctx, &exists, query, identityID); err != nil {
		return false, fmt.Errorf("check user exists by identity id: %w", err)
	}

	return exists, nil
}

// Save persists the full User aggregate state: upserts the user row, upserts
// all current employments, and deletes any employments that are no longer
// present in the aggregate. Must be called within a transaction.
func (r *userRepository) Save(ctx context.Context, user *model.User) error {
	tx := postgres.TxFromContext(ctx)
	if tx == nil {
		return fmt.Errorf("save user: no active transaction in context")
	}

	if err := r.upsertUser(ctx, tx, user); err != nil {
		return err
	}

	empRecs := entity.FromEmployments(user.Employments)

	for _, emp := range empRecs {
		if err := r.upsertEmployment(ctx, tx, emp); err != nil {
			return err
		}
	}

	return r.deleteRemovedEmployments(ctx, tx, user.ID, empRecs)
}

const upsertUserQuery = `
INSERT INTO users (
    id, identity_id, given_name, family_name, middle_name,
    custom_given_name, custom_family_name, custom_middle_name,
    admin_role_granted_at, admin_role_granted_by
) VALUES (
    :id, :identity_id, :given_name, :family_name, :middle_name,
    :custom_given_name, :custom_family_name, :custom_middle_name,
    :admin_role_granted_at, :admin_role_granted_by
)
ON CONFLICT (id) DO UPDATE SET
    identity_id           = EXCLUDED.identity_id,
    given_name            = EXCLUDED.given_name,
    family_name           = EXCLUDED.family_name,
    middle_name           = EXCLUDED.middle_name,
    custom_given_name     = EXCLUDED.custom_given_name,
    custom_family_name    = EXCLUDED.custom_family_name,
    custom_middle_name    = EXCLUDED.custom_middle_name,
    admin_role_granted_at = EXCLUDED.admin_role_granted_at,
    admin_role_granted_by = EXCLUDED.admin_role_granted_by`

func (r *userRepository) upsertUser(ctx context.Context, tx *sqlx.Tx, user *model.User) error {
	rec := entity.FromUser(user)
	if _, err := tx.NamedExecContext(ctx, upsertUserQuery, rec); err != nil {
		return fmt.Errorf("upsert user %s: %w", user.ID, err)
	}

	return nil
}

const upsertEmploymentQuery = `
INSERT INTO employments (
    id, user_id, organization_id, clinic_id, department_id,
    position, assigned_at, deputy_user_id, vacation_starts_at, vacation_ends_at
) VALUES (
    :id, :user_id, :organization_id, :clinic_id, :department_id,
    :position, :assigned_at, :deputy_user_id, :vacation_starts_at, :vacation_ends_at
)
ON CONFLICT (id) DO UPDATE SET
    position           = EXCLUDED.position,
    deputy_user_id     = EXCLUDED.deputy_user_id,
    vacation_starts_at = EXCLUDED.vacation_starts_at,
    vacation_ends_at   = EXCLUDED.vacation_ends_at`

func (r *userRepository) upsertEmployment(ctx context.Context, tx *sqlx.Tx, rec entity.EmploymentRecord) error {
	if _, err := tx.NamedExecContext(ctx, upsertEmploymentQuery, rec); err != nil {
		return fmt.Errorf("upsert employment %s: %w", rec.ID, err)
	}

	return nil
}

// deleteRemovedEmployments removes employment rows that exist in the database
// for this user but are absent from the current aggregate state (i.e. dismissed).
func (r *userRepository) deleteRemovedEmployments(ctx context.Context, tx *sqlx.Tx, userID uuid.UUID, current []entity.EmploymentRecord) error {
	if len(current) == 0 {
		const query = `DELETE FROM employments WHERE user_id = $1`
		if _, err := tx.ExecContext(ctx, query, userID); err != nil {
			return fmt.Errorf("delete all employments for user %s: %w", userID, err)
		}

		return nil
	}

	ids := make([]uuid.UUID, len(current))
	for i, emp := range current {
		ids[i] = emp.ID
	}

	const query = `DELETE FROM employments WHERE user_id = $1 AND NOT (id = ANY($2::uuid[]))`
	if _, err := tx.ExecContext(ctx, query, userID, pq.Array(ids)); err != nil {
		return fmt.Errorf("delete removed employments for user %s: %w", userID, err)
	}

	return nil
}
