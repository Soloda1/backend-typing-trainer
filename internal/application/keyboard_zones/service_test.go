package keyboardzones_test

import (
	"context"
	"errors"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"

	keyboardzonesapp "backend-typing-trainer/internal/application/keyboard_zones"
	"backend-typing-trainer/internal/domain/models"
	"backend-typing-trainer/internal/utils"
	"backend-typing-trainer/mocks"
)

func newServiceWithMocks(t *testing.T, repo *mocks.KeyboardZonesRepository) *keyboardzonesapp.Service {
	t.Helper()

	log := mocks.NewLogger(t)
	log.EXPECT().With(mock.Anything).Return(log).Once()
	log.EXPECT().Warn(mock.Anything, mock.Anything).Maybe()
	log.EXPECT().Warn(mock.Anything, mock.Anything, mock.Anything).Maybe()
	log.EXPECT().Error(mock.Anything, mock.Anything).Maybe()
	log.EXPECT().Error(mock.Anything, mock.Anything, mock.Anything).Maybe()

	return keyboardzonesapp.NewService(repo, log)
}

func TestService_GetByID(t *testing.T) {
	repo := mocks.NewKeyboardZonesRepository(t)
	svc := newServiceWithMocks(t, repo)

	_, err := svc.GetByID(context.Background(), uuid.Nil)
	if !errors.Is(err, utils.ErrInvalidRequest) {
		t.Fatalf("expected invalid request, got %v", err)
	}
	repo.AssertNotCalled(t, "GetByID", mock.Anything, mock.Anything)

	id := uuid.MustParse("aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaaa")
	zone := &models.KeyboardZone{ID: id, Name: "en_red", Symbols: "q,a,z"}
	repo.EXPECT().GetByID(mock.Anything, id).Return(zone, nil).Once()

	got, err := svc.GetByID(context.Background(), id)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got == nil || got.ID != id {
		t.Fatalf("unexpected result: %+v", got)
	}

	repo.EXPECT().GetByID(mock.Anything, id).Return(nil, errors.New("db timeout")).Once()
	_, err = svc.GetByID(context.Background(), id)
	if err == nil {
		t.Fatal("expected wrapped repository error")
	}
}

func TestService_GetByName(t *testing.T) {
	repo := mocks.NewKeyboardZonesRepository(t)
	svc := newServiceWithMocks(t, repo)

	_, err := svc.GetByName(context.Background(), "   ")
	if !errors.Is(err, utils.ErrInvalidRequest) {
		t.Fatalf("expected invalid request, got %v", err)
	}
	repo.AssertNotCalled(t, "GetByName", mock.Anything, mock.Anything)

	zone := &models.KeyboardZone{ID: uuid.MustParse("bbbbbbbb-bbbb-bbbb-bbbb-bbbbbbbbbbbb"), Name: "en_blue", Symbols: "y,u"}
	repo.EXPECT().GetByName(mock.Anything, "en_blue").Return(zone, nil).Once()

	got, err := svc.GetByName(context.Background(), "  en_blue  ")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got == nil || got.Name != "en_blue" {
		t.Fatalf("unexpected result: %+v", got)
	}

	repo.EXPECT().GetByName(mock.Anything, "en_blue").Return(nil, utils.ErrNotFound).Once()
	_, err = svc.GetByName(context.Background(), "en_blue")
	if !errors.Is(err, utils.ErrNotFound) {
		t.Fatalf("expected not found, got %v", err)
	}
}

func TestService_List(t *testing.T) {
	repo := mocks.NewKeyboardZonesRepository(t)
	svc := newServiceWithMocks(t, repo)

	_, err := svc.List(context.Background(), 0, 0)
	if !errors.Is(err, utils.ErrInvalidRequest) {
		t.Fatalf("expected invalid request, got %v", err)
	}
	repo.AssertNotCalled(t, "List", mock.Anything, mock.Anything, mock.Anything)

	expected := []*models.KeyboardZone{{ID: uuid.MustParse("cccccccc-cccc-cccc-cccc-cccccccccccc"), Name: "en_red"}}
	repo.EXPECT().List(mock.Anything, 10, 0).Return(expected, nil).Once()

	got, err := svc.List(context.Background(), 10, 0)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(got) != 1 || got[0].ID != expected[0].ID {
		t.Fatalf("unexpected result: %+v", got)
	}
}
