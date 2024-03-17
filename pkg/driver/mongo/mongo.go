package mongo

import (
	"context"

	"github.com/ReanSn0w/tk4go/pkg/tools"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type (
	Preferences struct {
		URI  string `long:"uri" env:"URI" description:"MongoDB URI" required:"true"`
		Name string `long:"name" env:"NAME" description:"MongoDB database name" required:"true"`
	}

	Mongo struct {
		client *mongo.Client
		name   string
	}
)

func Connect(ctx context.Context, logger tools.Logger, preferences Preferences) (*Mongo, error) {
	cl, err := mongo.Connect(ctx, options.Client().ApplyURI(preferences.URI))
	if err != nil {
		logger.Logf("[ERROR] failed to connect to MongoDB: %v", err)
		return nil, err
	}

	return &Mongo{
		client: cl,
		name:   preferences.Name,
	}, nil
}

// Operation executes a single operation on the database
func (m *Mongo) Operation(f func(db *mongo.Database) error) error {
	db := m.client.Database(m.name)
	return f(db)
}

// Transaction executes a series of actions in a single transaction
func (m *Mongo) Transaction(ctx context.Context, actions ...func(sessionCtx mongo.SessionContext, db *mongo.Database) error) error {
	session, err := m.client.StartSession()
	if err != nil {
		return err
	}

	defer session.EndSession(ctx)

	return mongo.WithSession(ctx, session, func(sessionCtx mongo.SessionContext) error {
		for _, action := range actions {
			err := action(sessionCtx, m.client.Database(m.name))
			if err != nil {
				return err
			}
		}

		return nil
	})
}

func (m *Mongo) Close(ctx context.Context) error {
	return m.client.Disconnect(ctx)
}
