package service

import (
	"Practice8/repository"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func TestGetUserByID(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := repository.NewMockUserRepository(ctrl)
	userService := NewUserService(mockRepo)

	user := &repository.User{ID: 1, Name: "Bakytzhan", Email: "a@b.com"}
	mockRepo.EXPECT().GetUserByID(1).Return(user, nil)

	result, err := userService.GetUserByID(1)

	assert.NoError(t, err)
	assert.Equal(t, user, result)
}

func TestCreateUser(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := repository.NewMockUserRepository(ctrl)
	userService := NewUserService(mockRepo)

	user := &repository.User{ID: 1, Name: "Bakytzhan", Email: "a@b.com"}
	mockRepo.EXPECT().CreateUser(user).Return(nil)

	err := userService.CreateUser(user)

	assert.NoError(t, err)
}

func TestRegisterUser_UserAlreadyExists(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := repository.NewMockUserRepository(ctrl)
	userService := NewUserService(mockRepo)

	user := &repository.User{ID: 1, Name: "A", Email: "a@b.com"}
	mockRepo.EXPECT().GetByEmail("a@b.com").Return(user, nil)

	err := userService.RegisterUser(user, "a@b.com")

	assert.Error(t, err)
	assert.EqualError(t, err, "user with this email already exists")
}

func TestRegisterUser_NewUserSuccess(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := repository.NewMockUserRepository(ctrl)
	userService := NewUserService(mockRepo)

	user := &repository.User{ID: 2, Name: "B", Email: "b@b.com"}
	mockRepo.EXPECT().GetByEmail("b@b.com").Return(nil, nil)
	mockRepo.EXPECT().CreateUser(user).Return(nil)

	err := userService.RegisterUser(user, "b@b.com")

	assert.NoError(t, err)
}

func TestRegisterUser_RepoErrorOnCreate(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := repository.NewMockUserRepository(ctrl)
	userService := NewUserService(mockRepo)

	user := &repository.User{ID: 3, Name: "C", Email: "c@b.com"}
	mockRepo.EXPECT().GetByEmail("c@b.com").Return(nil, nil)
	mockRepo.EXPECT().CreateUser(user).Return(errors.New("db error"))

	err := userService.RegisterUser(user, "c@b.com")

	assert.Error(t, err)
}

func TestUpdateUserName_EmptyName(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := repository.NewMockUserRepository(ctrl)
	userService := NewUserService(mockRepo)

	err := userService.UpdateUserName(1, "")

	assert.Error(t, err)
	assert.EqualError(t, err, "name cannot be empty")
}

func TestUpdateUserName_UserNotFound(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := repository.NewMockUserRepository(ctrl)
	userService := NewUserService(mockRepo)

	mockRepo.EXPECT().GetUserByID(1).Return(nil, errors.New("not found"))

	err := userService.UpdateUserName(1, "New Name")

	assert.Error(t, err)
}

func TestUpdateUserName_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := repository.NewMockUserRepository(ctrl)
	userService := NewUserService(mockRepo)

	user := &repository.User{ID: 1, Name: "Old Name", Email: "a@b.com"}
	mockRepo.EXPECT().GetUserByID(1).Return(user, nil)

	mockRepo.EXPECT().UpdateUser(gomock.AssignableToTypeOf(&repository.User{})).
		DoAndReturn(func(u *repository.User) error {
			assert.Equal(t, "New Name", u.Name)
			return nil
		})

	err := userService.UpdateUserName(1, "New Name")

	assert.NoError(t, err)
}

func TestUpdateUserName_UpdateFails(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := repository.NewMockUserRepository(ctrl)
	userService := NewUserService(mockRepo)

	user := &repository.User{ID: 1, Name: "Old Name", Email: "a@b.com"}
	mockRepo.EXPECT().GetUserByID(1).Return(user, nil)
	mockRepo.EXPECT().UpdateUser(gomock.Any()).Return(errors.New("update failed"))

	err := userService.UpdateUserName(1, "New Name")

	assert.Error(t, err)
}

func TestDeleteUser_AdminBlocked(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := repository.NewMockUserRepository(ctrl)
	userService := NewUserService(mockRepo)

	err := userService.DeleteUser(1)

	assert.Error(t, err)
	assert.EqualError(t, err, "it is not allowed to delete admin user")
}

func TestDeleteUser_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := repository.NewMockUserRepository(ctrl)
	userService := NewUserService(mockRepo)

	mockRepo.EXPECT().DeleteUser(2).Return(nil)

	err := userService.DeleteUser(2)

	assert.NoError(t, err)
}

func TestDeleteUser_RepoError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := repository.NewMockUserRepository(ctrl)
	userService := NewUserService(mockRepo)

	mockRepo.EXPECT().DeleteUser(2).Return(errors.New("delete failed"))

	err := userService.DeleteUser(2)

	assert.Error(t, err)
}
