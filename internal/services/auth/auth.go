package auth

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"time"

	"github.com/alexzanser/sso/internal/domain/models"
	"github.com/alexzanser/sso/internal/lib/jwt"
	"github.com/alexzanser/sso/internal/storage"
	"golang.org/x/crypto/bcrypt"
)

var (
	ErrInvalidCredentials = errors.New("invalid credentials")
	ErrInvalidUserID      = errors.New("invalid user ID")
	ErrUserAlreadyExists  = errors.New("user already exists")
)

type AuthService struct {
	log          *slog.Logger
	userSaver    UserSaver
	userProvider UserProvider
	appProvider  AppProvider
	tokenTTL     time.Duration
}

type UserSaver interface {
	SaveUser(ctx context.Context, email string, passHash []byte) (uid int64, err error)
}

type UserProvider interface {
	User(ctx context.Context, email string) (models.User, error)
	IsAdmin(ctx context.Context, userID int64) (bool, error)
}

type AppProvider interface {
	App(ctx context.Context, appID int32) (models.App, error)
}

// New creates a new Auth service instance with the provided dependencies.
func NewService(log *slog.Logger,
	userSaver UserSaver,
	userProvider UserProvider,
	appProvider AppProvider,
	tokenTTL time.Duration,
) *AuthService {
	return &AuthService{
		log:          log.With(slog.String("service", "auth")),
		userSaver:    userSaver,
		userProvider: userProvider,
		appProvider:  appProvider,
		tokenTTL:     tokenTTL,
	}
}

// RegisterNewUser registers a new user with the provided email and password.
func (a *AuthService) RegisterNewUser(ctx context.Context, email, password string) (int64, error) {
	const op = "auth.RegisterNewUser"

	log := a.log.With(
		slog.String("op", op),
		slog.String("email", email))

	passHash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		log.Error("internal error", slog.String("email", email), slog.Any("error", err))
		return 0, err
	}

	id, err := a.userSaver.SaveUser(ctx, email, passHash)
	if err != nil {
		if errors.Is(err, storage.ErrUserExist) {
			log.Warn("user already exists", slog.String("email", email))
			return 0, fmt.Errorf("%s: %w", op, ErrUserAlreadyExists)
		}
		log.Error("failed to save user", slog.Any("error", err))
		return 0, err
	}

	log.Info("Registering new user")
	return id, nil
}

// Login authenticates a user with the provided email and password, and returns a JWT token if successful.
func (a *AuthService) Login(ctx context.Context, email, password string, appID int32) (string, error) {
	const op = "auth.Login"

	log := a.log.With(
		slog.String("op", op),
		slog.String("email", email))

	user, err := a.userProvider.User(ctx, email)
	if err != nil {
		if errors.Is(err, storage.ErrUserNotFound) {
			log.Warn("user not found", slog.String("email", email))
			return "", fmt.Errorf("%s: %w", op, ErrInvalidCredentials)
		}

		log.Error("failed to get user", slog.Any("error", err))
		return "", err
	}

	if err = bcrypt.CompareHashAndPassword(user.PassHash, []byte(password)); err != nil {
		log.Error("invalid credetials", slog.Any("error", err))
		return "", err
	}

	app, err := a.appProvider.App(ctx, appID)
	if err != nil {
		log.Error("failed to get app", slog.Any("error", err))
		return "", err
	}

	token, err := jwt.NewToken(user, app, a.tokenTTL)
	if err != nil {
		log.Error("failed to generate token", slog.Any("error", err))
		return "", fmt.Errorf("%s: %w", op, err)
	}
	log.Info("User logged in successfully")

	return token, nil
}

// IsAdmin checks if the user with the given userID is an admin.
func (a *AuthService) IsAdmin(ctx context.Context, userID int64) (bool, error) {
	const op = "auth.IsAdmin"

	log := a.log.With(
		slog.String("op", op),
		slog.Int64("userID", userID))

	isAdmin, err := a.userProvider.IsAdmin(ctx, userID)
	if err != nil {
		if errors.Is(err, storage.ErrUserNotFound) {
			log.Warn("user not found", slog.Int64("userID", userID))
			return false, ErrInvalidUserID
		}
		return false, fmt.Errorf("%s: %w", op, ErrInvalidUserID)
	}

	log.Info("Checked admin status", slog.Bool("isAdmin", isAdmin))
	return isAdmin, nil
}
