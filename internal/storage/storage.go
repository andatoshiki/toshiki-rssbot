//go:generate mockgen -source=storage.go -destination=./mock/storage_mock.go -package=mock

package storage

import (
	"context"
	"errors"

	"github.com/andatoshiki/toshiki-rssbot/internal/model"
)

var (
	// ErrRecordNotFound record not found error
	ErrRecordNotFound = errors.New("record not found")
)

type Storage interface {
	Init(ctx context.Context) error
}

// User user storage interface
type User interface {
	Storage
	CreateUser(ctx context.Context, user *model.User) error
	GetUser(ctx context.Context, id int64) (*model.User, error)
}

// Source subscription source storage interface
type Source interface {
	Storage
	AddSource(ctx context.Context, source *model.Source) error
	GetSource(ctx context.Context, id uint) (*model.Source, error)
	GetSources(ctx context.Context) ([]*model.Source, error)
	GetSourceByURL(ctx context.Context, url string) (*model.Source, error)
	Delete(ctx context.Context, id uint) error
	UpsertSource(ctx context.Context, sourceID uint, newSource *model.Source) error
}

type SubscriptionSortType = int

const (
	SubscriptionSortTypeCreatedTimeDesc SubscriptionSortType = iota
)

type GetSubscriptionsOptions struct {
	Count    int // Number of items to retrieve, -1 to retrieve all
	Offset   int
	SortType SubscriptionSortType
}

type GetSubscriptionsResult struct {
	Subscriptions []*model.Subscribe
	HasMore       bool
}

type Subscription interface {
	Storage
	AddSubscription(ctx context.Context, subscription *model.Subscribe) error
	SubscriptionExist(ctx context.Context, userID int64, sourceID uint) (bool, error)
	GetSubscription(ctx context.Context, userID int64, sourceID uint) (*model.Subscribe, error)
	GetSubscriptionsByUserID(
		ctx context.Context, userID int64, opts *GetSubscriptionsOptions,
	) (*GetSubscriptionsResult, error)
	GetSubscriptionsBySourceID(
		ctx context.Context, sourceID uint, opts *GetSubscriptionsOptions,
	) (*GetSubscriptionsResult, error)
	CountSubscriptions(ctx context.Context) (int64, error)
	DeleteSubscription(ctx context.Context, userID int64, sourceID uint) (int64, error)
	CountSourceSubscriptions(ctx context.Context, sourceID uint) (int64, error)
	UpdateSubscription(
		ctx context.Context, userID int64, sourceID uint, newSubscription *model.Subscribe,
	) error
	UpsertSubscription(
		ctx context.Context, userID int64, sourceID uint, newSubscription *model.Subscribe,
	) error
}

type Content interface {
	Storage
	// AddContent adds a new article
	AddContent(ctx context.Context, content *model.Content) error
	// DeleteSourceContents deletes all articles of a subscription source and returns the number of deleted articles
	DeleteSourceContents(ctx context.Context, sourceID uint) (int64, error)
	// HashIDExist checks if an article with the given hash id already exists
	HashIDExist(ctx context.Context, hashID string) (bool, error)
}
