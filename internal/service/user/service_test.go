package user

import (
	"context"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"github.com/ulbwa/medincident-command-service/internal/common/errors"
	"github.com/ulbwa/medincident-command-service/internal/common/outbox"
	"github.com/ulbwa/medincident-command-service/internal/common/persistence"
	"github.com/ulbwa/medincident-command-service/internal/model"
)

// Mocks

type MockTxFactory struct {
	mock.Mock
}

func (m *MockTxFactory) Begin(ctx context.Context) (context.Context, persistence.Transaction, error) {
	args := m.Called(ctx)
	return args.Get(0).(context.Context), args.Get(1).(persistence.Transaction), args.Error(2)
}

type MockTx struct {
	mock.Mock
}

func (m *MockTx) Close() error {
	return m.Called().Error(0)
}

func (m *MockTx) Savepoint(ctx context.Context, name string) error {
	return m.Called(ctx, name).Error(0)
}

func (m *MockTx) Commit(ctx context.Context) error {
	return m.Called(ctx).Error(0)
}

func (m *MockTx) Rollback(ctx context.Context) error {
	return m.Called(ctx).Error(0)
}

func (m *MockTx) RollbackSavepoint(ctx context.Context, name string) error {
	return m.Called(ctx, name).Error(0)
}

type MockEventDispatcher struct {
	mock.Mock
}

func (m *MockEventDispatcher) Dispatch(ctx context.Context, tx persistence.Transaction, sources ...outbox.EventSource) error {
	callArgs := []interface{}{ctx, tx}
	for _, s := range sources {
		callArgs = append(callArgs, s)
	}
	args := m.Called(callArgs...)
	return args.Error(0)
}

type MockRepo struct {
	mock.Mock
}

func (m *MockRepo) GetByID(ctx context.Context, id uuid.UUID) (*model.User, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.User), args.Error(1)
}

func (m *MockRepo) ExistsByID(ctx context.Context, id uuid.UUID) (bool, error) {
	args := m.Called(ctx, id)
	return args.Bool(0), args.Error(1)
}

func (m *MockRepo) ExistsByIdentityID(ctx context.Context, identityID string) (bool, error) {
	args := m.Called(ctx, identityID)
	return args.Bool(0), args.Error(1)
}

func (m *MockRepo) Save(ctx context.Context, user *model.User) error {
	args := m.Called(ctx, user)
	return args.Error(0)
}

type MockIdentityProvider struct {
	mock.Mock
}

func (m *MockIdentityProvider) Get(ctx context.Context, identityID string) (*Identity, error) {
	args := m.Called(ctx, identityID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*Identity), args.Error(1)
}

func (m *MockIdentityProvider) UpdateHuman(ctx context.Context, identityID string, human *IdentityHuman) (*Identity, error) {
	args := m.Called(ctx, identityID, human)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*Identity), args.Error(1)
}

func (m *MockIdentityProvider) UpdateUserMetadata(ctx context.Context, identityID string, metadata *IdentityUserMetadata) error {
	args := m.Called(ctx, identityID, metadata)
	return args.Error(0)
}

// Tests

func TestService_Create(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	const identityID = "zitadel|123456789"

	t.Run("Success", func(t *testing.T) {
		mockTxFactory := new(MockTxFactory)
		mockTx := new(MockTx)
		mockDispatcher := new(MockEventDispatcher)
		mockRepo := new(MockRepo)
		mockIDP := new(MockIdentityProvider)

		svc, err := NewService(mockTxFactory, mockDispatcher, mockIDP, mockRepo)
		require.NoError(t, err)

		humanProfile := &IdentityHuman{GivenName: "Test", FamilyName: "User"}
		identity := &Identity{ID: identityID, Human: humanProfile, IsActive: true}

		// Synchronous call inside Create.
		mockIDP.On("Get", ctx, identityID).Return(identity, nil).Once()
		// Background goroutines fire after the transaction commits with derived contexts.
		// Return an error so syncHumanIdentity exits without calling UpdateHuman.
		mockIDP.On("Get", mock.Anything, identityID).Return(nil, errors.ErrInvalidRequest).Maybe()
		mockIDP.On("UpdateUserMetadata", mock.Anything, identityID, mock.Anything).Return(nil).Maybe()

		mockTxFactory.On("Begin", ctx).Return(ctx, mockTx, nil)
		mockRepo.On("ExistsByIdentityID", ctx, identityID).Return(false, nil)
		mockRepo.On("Save", ctx, mock.MatchedBy(func(u *model.User) bool {
			return u.IdentityID == identityID &&
				u.ID != uuid.Nil &&
				u.Name.GivenName == "Test" &&
				u.Name.FamilyName == "User"
		})).Return(nil)
		mockDispatcher.On("Dispatch", ctx, mockTx, mock.AnythingOfType("*model.User")).Return(nil)
		mockTx.On("Commit", ctx).Return(nil)
		mockTx.On("Close").Return(nil)

		resp, err := svc.Create(ctx, &CreateUserRequest{
			IdentityID: identityID,
			GivenName:  "Test",
			FamilyName: "User",
		})
		require.NoError(t, err)
		require.NotNil(t, resp)
		require.NotNil(t, resp.User)
		assert.Equal(t, identityID, resp.User.IdentityID)
		assert.NotEqual(t, uuid.Nil, resp.User.ID)

		mockTxFactory.AssertExpectations(t)
		mockRepo.AssertExpectations(t)
		mockDispatcher.AssertExpectations(t)
		mockTx.AssertExpectations(t)
	})

	t.Run("IdentityNotFound", func(t *testing.T) {
		mockTxFactory := new(MockTxFactory)
		mockDispatcher := new(MockEventDispatcher)
		mockRepo := new(MockRepo)
		mockIDP := new(MockIdentityProvider)

		svc, _ := NewService(mockTxFactory, mockDispatcher, mockIDP, mockRepo)

		mockIDP.On("Get", ctx, identityID).Return(nil, errors.ErrIdentityNotFound)

		_, err := svc.Create(ctx, &CreateUserRequest{
			IdentityID: identityID,
			GivenName:  "Test",
			FamilyName: "User",
		})
		assert.ErrorIs(t, err, errors.ErrIdentityNotFound)
	})

	t.Run("IdentityAlreadyLinked", func(t *testing.T) {
		mockTxFactory := new(MockTxFactory)
		mockDispatcher := new(MockEventDispatcher)
		mockRepo := new(MockRepo)
		mockIDP := new(MockIdentityProvider)

		svc, _ := NewService(mockTxFactory, mockDispatcher, mockIDP, mockRepo)

		existingUserID := uuid.MustParse("01010101-0101-7101-8101-010101010101")
		humanProfile := &IdentityHuman{GivenName: "Test", FamilyName: "User"}
		identity := &Identity{
			ID:       identityID,
			Human:    humanProfile,
			IsActive: true,
			Metadata: &IdentityUserMetadata{UserID: &existingUserID},
		}
		mockIDP.On("Get", ctx, identityID).Return(identity, nil)

		_, err := svc.Create(ctx, &CreateUserRequest{
			IdentityID: identityID,
			GivenName:  "Test",
			FamilyName: "User",
		})
		assert.ErrorIs(t, err, errors.ErrIdentityAlreadyLinked)
	})

	t.Run("IdentityNotHuman", func(t *testing.T) {
		mockTxFactory := new(MockTxFactory)
		mockDispatcher := new(MockEventDispatcher)
		mockRepo := new(MockRepo)
		mockIDP := new(MockIdentityProvider)

		svc, _ := NewService(mockTxFactory, mockDispatcher, mockIDP, mockRepo)

		// Service account or incomplete registration — no Human profile
		identity := &Identity{ID: identityID, Human: nil, IsActive: true}
		mockIDP.On("Get", ctx, identityID).Return(identity, nil)

		_, err := svc.Create(ctx, &CreateUserRequest{
			IdentityID: identityID,
			GivenName:  "Test",
			FamilyName: "User",
		})
		assert.ErrorIs(t, err, errors.ErrIdentityNotHuman)
	})

	t.Run("UserAlreadyExists", func(t *testing.T) {
		mockTxFactory := new(MockTxFactory)
		mockTx := new(MockTx)
		mockDispatcher := new(MockEventDispatcher)
		mockRepo := new(MockRepo)
		mockIDP := new(MockIdentityProvider)

		svc, _ := NewService(mockTxFactory, mockDispatcher, mockIDP, mockRepo)

		humanProfile := &IdentityHuman{GivenName: "Test", FamilyName: "User"}
		identity := &Identity{ID: identityID, Human: humanProfile, IsActive: true}

		mockIDP.On("Get", ctx, identityID).Return(identity, nil)
		mockTxFactory.On("Begin", ctx).Return(ctx, mockTx, nil)
		mockRepo.On("ExistsByIdentityID", ctx, identityID).Return(true, nil)
		mockTx.On("Rollback", mock.Anything).Return(nil)
		mockTx.On("Close").Return(nil)

		_, err := svc.Create(ctx, &CreateUserRequest{
			IdentityID: identityID,
			GivenName:  "Test",
			FamilyName: "User",
		})
		assert.ErrorIs(t, err, errors.ErrUserAlreadyExists)
	})
}

func TestService_GrantAdminRole(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	userID := uuid.MustParse("01010101-0101-7101-8101-010101010101")
	actorID := uuid.MustParse("02020202-0202-7202-8202-020202020202")
	granterID := uuid.MustParse("03030303-0303-7303-8303-030303030303")

	t.Run("Success", func(t *testing.T) {
		mockTxFactory := new(MockTxFactory)
		mockTx := new(MockTx)
		mockDispatcher := new(MockEventDispatcher)
		mockRepo := new(MockRepo)
		mockIDP := new(MockIdentityProvider)

		svc, err := NewService(mockTxFactory, mockDispatcher, mockIDP, mockRepo)
		require.NoError(t, err)

		un, _ := model.NewUserName("Test", "User", nil)
		user, _ := model.NewUser(userID, "idp|user", un)
		user.PopEvents() // clear creation event

		actorUn, _ := model.NewUserName("Actor", "User", nil)
		actor, _ := model.NewUser(actorID, "idp|actor", actorUn)
		oldAdminTime := time.Now().Add(-100 * time.Hour)
		actor.AdminRole = &model.AdminRole{
			GrantedAt: oldAdminTime,
			GrantedBy: granterID,
		}
		actor.PopEvents()

		mockTxFactory.On("Begin", ctx).Return(ctx, mockTx, nil)
		mockRepo.On("GetByID", ctx, actorID).Return(actor, nil)
		mockRepo.On("GetByID", ctx, userID).Return(user, nil)
		mockRepo.On("Save", ctx, mock.MatchedBy(func(u *model.User) bool {
			return u.ID == userID && u.AdminRole != nil
		})).Return(nil)
		mockDispatcher.On("Dispatch", ctx, mockTx, user).Return(nil)
		mockTx.On("Commit", ctx).Return(nil)
		mockTx.On("Close").Return(nil)

		err = svc.GrantAdminRole(ctx, &GrantAdminRoleRequest{
			ActorID:      actorID,
			TargetUserID: userID,
		})
		require.NoError(t, err)

		mockTxFactory.AssertExpectations(t)
		mockRepo.AssertExpectations(t)
		mockDispatcher.AssertExpectations(t)
		mockTx.AssertExpectations(t)
	})

	t.Run("Unauthorized_NotAdmin", func(t *testing.T) {
		mockTxFactory := new(MockTxFactory)
		mockTx := new(MockTx)
		mockDispatcher := new(MockEventDispatcher)
		mockRepo := new(MockRepo)
		mockIDP := new(MockIdentityProvider)

		svc, _ := NewService(mockTxFactory, mockDispatcher, mockIDP, mockRepo)

		actorUn, _ := model.NewUserName("Actor", "User", nil)
		actor, _ := model.NewUser(actorID, "idp|actor", actorUn)
		// Not admin

		mockTxFactory.On("Begin", ctx).Return(ctx, mockTx, nil)
		mockRepo.On("GetByID", ctx, actorID).Return(actor, nil)
		mockTx.On("Rollback", mock.Anything).Return(nil)
		mockTx.On("Close").Return(nil)

		err := svc.GrantAdminRole(ctx, &GrantAdminRoleRequest{
			ActorID:      actorID,
			TargetUserID: userID,
		})
		assert.ErrorIs(t, err, errors.ErrAdminRoleGrantActorNotAdmin)
	})

	t.Run("Unauthorized_NewAdmin", func(t *testing.T) {
		mockTxFactory := new(MockTxFactory)
		mockTx := new(MockTx)
		mockDispatcher := new(MockEventDispatcher)
		mockRepo := new(MockRepo)
		mockIDP := new(MockIdentityProvider)

		svc, _ := NewService(mockTxFactory, mockDispatcher, mockIDP, mockRepo)

		actorUn, _ := model.NewUserName("Actor", "User", nil)
		actor, _ := model.NewUser(actorID, "idp|actor", actorUn)
		recentAdminTime := time.Now()
		actor.AdminRole = &model.AdminRole{
			GrantedAt: recentAdminTime,
			GrantedBy: granterID,
		}

		mockTxFactory.On("Begin", ctx).Return(ctx, mockTx, nil)
		mockRepo.On("GetByID", ctx, actorID).Return(actor, nil)
		mockTx.On("Rollback", mock.Anything).Return(nil)
		mockTx.On("Close").Return(nil)

		err := svc.GrantAdminRole(ctx, &GrantAdminRoleRequest{
			ActorID:      actorID,
			TargetUserID: userID,
		})
		assert.ErrorIs(t, err, errors.ErrAdminRoleGrantInsufficientTenure)
	})
}

func TestService_RevokeAdminRole(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	userID := uuid.MustParse("01010101-0101-7101-8101-010101010101")
	actorID := uuid.MustParse("02020202-0202-7202-8202-020202020202")
	granterID := uuid.MustParse("03030303-0303-7303-8303-030303030303")
	otherGranterID := uuid.MustParse("04040404-0404-7404-8404-040404040404")

	t.Run("Success", func(t *testing.T) {
		mockTxFactory := new(MockTxFactory)
		mockTx := new(MockTx)
		mockDispatcher := new(MockEventDispatcher)
		mockRepo := new(MockRepo)
		mockIDP := new(MockIdentityProvider)

		svc, err := NewService(mockTxFactory, mockDispatcher, mockIDP, mockRepo)
		require.NoError(t, err)

		un, _ := model.NewUserName("Test", "User", nil)
		user, _ := model.NewUser(userID, "idp|user", un)
		userAdminTime := time.Now()
		user.AdminRole = &model.AdminRole{
			GrantedAt: userAdminTime,
			GrantedBy: granterID,
		}
		user.PopEvents()

		actorUn, _ := model.NewUserName("Actor", "User", nil)
		actor, _ := model.NewUser(actorID, "idp|actor", actorUn)
		actorAdminTime := time.Now().Add(-100 * time.Hour)
		actor.AdminRole = &model.AdminRole{
			GrantedAt: actorAdminTime,
			GrantedBy: otherGranterID,
		}
		actor.PopEvents()

		mockTxFactory.On("Begin", ctx).Return(ctx, mockTx, nil)
		mockRepo.On("GetByID", ctx, actorID).Return(actor, nil)
		mockRepo.On("GetByID", ctx, userID).Return(user, nil)
		mockRepo.On("Save", ctx, mock.MatchedBy(func(u *model.User) bool {
			return u.ID == userID && u.AdminRole == nil
		})).Return(nil)
		mockDispatcher.On("Dispatch", ctx, mockTx, user).Return(nil)
		mockTx.On("Commit", ctx).Return(nil)
		mockTx.On("Close").Return(nil)

		err = svc.RevokeAdminRole(ctx, &RevokeAdminRoleRequest{
			ActorID:      actorID,
			TargetUserID: userID,
		})
		require.NoError(t, err)

		mockTxFactory.AssertExpectations(t)
		mockRepo.AssertExpectations(t)
		mockDispatcher.AssertExpectations(t)
		mockTx.AssertExpectations(t)
	})

	t.Run("Unauthorized_NotAdmin_Revoke", func(t *testing.T) {
		mockTxFactory := new(MockTxFactory)
		mockTx := new(MockTx)
		mockDispatcher := new(MockEventDispatcher)
		mockRepo := new(MockRepo)
		mockIDP := new(MockIdentityProvider)

		svc, _ := NewService(mockTxFactory, mockDispatcher, mockIDP, mockRepo)

		actorUn, _ := model.NewUserName("Actor", "User", nil)
		actor, _ := model.NewUser(actorID, "idp|actor", actorUn)
		// Not admin

		mockTxFactory.On("Begin", ctx).Return(ctx, mockTx, nil)
		mockRepo.On("GetByID", ctx, actorID).Return(actor, nil)
		mockTx.On("Rollback", mock.Anything).Return(nil)
		mockTx.On("Close").Return(nil)

		err := svc.RevokeAdminRole(ctx, &RevokeAdminRoleRequest{
			ActorID:      actorID,
			TargetUserID: userID,
		})
		assert.ErrorIs(t, err, errors.ErrAdminRoleGrantActorNotAdmin)
	})

	t.Run("Unauthorized_NewAdmin_Revoke", func(t *testing.T) {
		mockTxFactory := new(MockTxFactory)
		mockTx := new(MockTx)
		mockDispatcher := new(MockEventDispatcher)
		mockRepo := new(MockRepo)
		mockIDP := new(MockIdentityProvider)

		svc, _ := NewService(mockTxFactory, mockDispatcher, mockIDP, mockRepo)

		actorUn, _ := model.NewUserName("Actor", "User", nil)
		actor, _ := model.NewUser(actorID, "idp|actor", actorUn)
		recentAdminTime := time.Now()
		actor.AdminRole = &model.AdminRole{
			GrantedAt: recentAdminTime,
			GrantedBy: granterID,
		}

		mockTxFactory.On("Begin", ctx).Return(ctx, mockTx, nil)
		mockRepo.On("GetByID", ctx, actorID).Return(actor, nil)
		mockTx.On("Rollback", mock.Anything).Return(nil)
		mockTx.On("Close").Return(nil)

		err := svc.RevokeAdminRole(ctx, &RevokeAdminRoleRequest{
			ActorID:      actorID,
			TargetUserID: userID,
		})
		assert.ErrorIs(t, err, errors.ErrAdminRoleGrantInsufficientTenure)
	})
}
