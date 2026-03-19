package user

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"github.com/ulbwa/medincident-command-service/internal/common/authorization"
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

func (m *MockRepo) GetByID(ctx context.Context, id int64) (*model.User, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.User), args.Error(1)
}

func (m *MockRepo) ExistsByID(ctx context.Context, id int64) (bool, error) {
	args := m.Called(ctx, id)
	return args.Bool(0), args.Error(1)
}

func (m *MockRepo) Save(ctx context.Context, user *model.User) error {
	args := m.Called(ctx, user)
	return args.Error(0)
}

type MockIdentityProvider struct {
	mock.Mock
}

func (m *MockIdentityProvider) Get(ctx context.Context, id int64) (*Identity, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*Identity), args.Error(1)
}

func (m *MockIdentityProvider) UpdateHuman(ctx context.Context, id int64, human *IdentityHuman) (*Identity, error) {
	args := m.Called(ctx, id, human)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*Identity), args.Error(1)
}

type MockTokenIntrospector struct {
	mock.Mock
}

func (m *MockTokenIntrospector) Introspect(ctx context.Context, bearerToken string) (*authorization.AccessClaims, error) {
	args := m.Called(ctx, bearerToken)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*authorization.AccessClaims), args.Error(1)
}

// Tests

func TestService_GrantAdminRole(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	userID := int64(1 << 23)
	actorID := int64(2 << 23)

	t.Run("Success", func(t *testing.T) {
		mockTxFactory := new(MockTxFactory)
		mockTx := new(MockTx)
		mockDispatcher := new(MockEventDispatcher)
		mockRepo := new(MockRepo)
		mockIDP := new(MockIdentityProvider)
		mockTokenIntrospector := new(MockTokenIntrospector)

		svc, err := NewService(mockTxFactory, mockDispatcher, mockIDP, mockTokenIntrospector, mockRepo)
		require.NoError(t, err)

		un, _ := model.NewUserName("Test", "User", nil)
		user, _ := model.NewUser(userID, un)
		user.PopEvents() // clear creation event

		actorUn, _ := model.NewUserName("Actor", "User", nil)
		actor, _ := model.NewUser(actorID, actorUn)
		oldAdminTime := time.Now().Add(-100 * time.Hour)
		actor.AdminRole = &model.AdminRole{
			GrantedAt: oldAdminTime,
			GrantedBy: int64(3 << 23),
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
		mockTokenIntrospector := new(MockTokenIntrospector)

		svc, _ := NewService(mockTxFactory, mockDispatcher, mockIDP, mockTokenIntrospector, mockRepo)

		actorUn, _ := model.NewUserName("Actor", "User", nil)
		actor, _ := model.NewUser(actorID, actorUn)
		// Not admin

		mockTxFactory.On("Begin", ctx).Return(ctx, mockTx, nil)
		mockRepo.On("GetByID", ctx, actorID).Return(actor, nil)
		mockTx.On("Rollback", mock.Anything).Return(nil)
		mockTx.On("Close").Return(nil)

		err := svc.GrantAdminRole(ctx, &GrantAdminRoleRequest{
			ActorID:      actorID,
			TargetUserID: userID,
		})
		assert.ErrorIs(t, err, errors.ErrAdminRoleGrantForbidden)
	})

	t.Run("Unauthorized_NewAdmin", func(t *testing.T) {
		mockTxFactory := new(MockTxFactory)
		mockTx := new(MockTx)
		mockDispatcher := new(MockEventDispatcher)
		mockRepo := new(MockRepo)
		mockIDP := new(MockIdentityProvider)
		mockTokenIntrospector := new(MockTokenIntrospector)

		svc, _ := NewService(mockTxFactory, mockDispatcher, mockIDP, mockTokenIntrospector, mockRepo)

		actorUn, _ := model.NewUserName("Actor", "User", nil)
		actor, _ := model.NewUser(actorID, actorUn)
		recentAdminTime := time.Now()
		actor.AdminRole = &model.AdminRole{
			GrantedAt: recentAdminTime,
			GrantedBy: int64(3 << 23),
		}

		mockTxFactory.On("Begin", ctx).Return(ctx, mockTx, nil)
		mockRepo.On("GetByID", ctx, actorID).Return(actor, nil)
		mockTx.On("Rollback", mock.Anything).Return(nil)
		mockTx.On("Close").Return(nil)

		err := svc.GrantAdminRole(ctx, &GrantAdminRoleRequest{
			ActorID:      actorID,
			TargetUserID: userID,
		})
		assert.ErrorIs(t, err, errors.ErrAdminRoleGrantForbidden)
	})
}

func TestService_RevokeAdminRole(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	userID := int64(1 << 23)
	actorID := int64(2 << 23)

	t.Run("Success", func(t *testing.T) {
		mockTxFactory := new(MockTxFactory)
		mockTx := new(MockTx)
		mockDispatcher := new(MockEventDispatcher)
		mockRepo := new(MockRepo)
		mockIDP := new(MockIdentityProvider)
		mockTokenIntrospector := new(MockTokenIntrospector)

		svc, err := NewService(mockTxFactory, mockDispatcher, mockIDP, mockTokenIntrospector, mockRepo)
		require.NoError(t, err)

		un, _ := model.NewUserName("Test", "User", nil)
		user, _ := model.NewUser(userID, un)
		userAdminTime := time.Now()
		user.AdminRole = &model.AdminRole{
			GrantedAt: userAdminTime,
			GrantedBy: int64(3 << 23),
		}
		user.PopEvents()

		actorUn, _ := model.NewUserName("Actor", "User", nil)
		actor, _ := model.NewUser(actorID, actorUn)
		actorAdminTime := time.Now().Add(-100 * time.Hour)
		actor.AdminRole = &model.AdminRole{
			GrantedAt: actorAdminTime,
			GrantedBy: int64(4 << 23),
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
		mockTokenIntrospector := new(MockTokenIntrospector)

		svc, _ := NewService(mockTxFactory, mockDispatcher, mockIDP, mockTokenIntrospector, mockRepo)

		actorUn, _ := model.NewUserName("Actor", "User", nil)
		actor, _ := model.NewUser(actorID, actorUn)
		// Not admin

		mockTxFactory.On("Begin", ctx).Return(ctx, mockTx, nil)
		mockRepo.On("GetByID", ctx, actorID).Return(actor, nil)
		mockTx.On("Rollback", mock.Anything).Return(nil)
		mockTx.On("Close").Return(nil)

		err := svc.RevokeAdminRole(ctx, &RevokeAdminRoleRequest{
			ActorID:      actorID,
			TargetUserID: userID,
		})
		assert.ErrorIs(t, err, errors.ErrAdminRoleGrantForbidden)
	})

	t.Run("Unauthorized_NewAdmin_Revoke", func(t *testing.T) {
		mockTxFactory := new(MockTxFactory)
		mockTx := new(MockTx)
		mockDispatcher := new(MockEventDispatcher)
		mockRepo := new(MockRepo)
		mockIDP := new(MockIdentityProvider)
		mockTokenIntrospector := new(MockTokenIntrospector)

		svc, _ := NewService(mockTxFactory, mockDispatcher, mockIDP, mockTokenIntrospector, mockRepo)

		actorUn, _ := model.NewUserName("Actor", "User", nil)
		actor, _ := model.NewUser(actorID, actorUn)
		recentAdminTime := time.Now()
		actor.AdminRole = &model.AdminRole{
			GrantedAt: recentAdminTime,
			GrantedBy: int64(3 << 23),
		}

		mockTxFactory.On("Begin", ctx).Return(ctx, mockTx, nil)
		mockRepo.On("GetByID", ctx, actorID).Return(actor, nil)
		mockTx.On("Rollback", mock.Anything).Return(nil)
		mockTx.On("Close").Return(nil)

		err := svc.RevokeAdminRole(ctx, &RevokeAdminRoleRequest{
			ActorID:      actorID,
			TargetUserID: userID,
		})
		assert.ErrorIs(t, err, errors.ErrAdminRoleGrantForbidden)
	})
}

func TestNewService_NilDependencies(t *testing.T) {
	t.Parallel()

	newValidDeps := func() (persistence.TransactionFactory, outbox.EventDispatcher, IdentityProvider, authorization.AccessTokenIntrospector, Repository) {
		return new(MockTxFactory), new(MockEventDispatcher), new(MockIdentityProvider), new(MockTokenIntrospector), new(MockRepo)
	}

	tests := []struct {
		name        string
		nilArg      string
		expectedErr string
	}{
		{
			name:        "NilTxFactory",
			nilArg:      "txFactory",
			expectedErr: "transaction factory is required",
		},
		{
			name:        "NilEventDispatcher",
			nilArg:      "eventDispatcher",
			expectedErr: "event dispatcher is required",
		},
		{
			name:        "NilIdentityProvider",
			nilArg:      "identityProvider",
			expectedErr: "identity provider is required",
		},
		{
			name:        "NilTokenIntrospector",
			nilArg:      "tokenIntrospector",
			expectedErr: "token introspector is required",
		},
		{
			name:        "NilRepository",
			nilArg:      "repo",
			expectedErr: "repository is required",
		},
	}

	for _, tc := range tests {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			txFactory, eventDispatcher, identityProvider, tokenIntrospector, repo := newValidDeps()

			switch tc.nilArg {
			case "txFactory":
				txFactory = nil
			case "eventDispatcher":
				eventDispatcher = nil
			case "identityProvider":
				identityProvider = nil
			case "tokenIntrospector":
				tokenIntrospector = nil
			case "repo":
				repo = nil
			}

			svc, err := NewService(txFactory, eventDispatcher, identityProvider, tokenIntrospector, repo)
			require.Error(t, err)
			require.Nil(t, svc)
			require.EqualError(t, err, tc.expectedErr)
		})
	}
}
