package repository

import (
	"context"

	"github.com/dejandjenic/go-gin-sample/application/configuration"
	"github.com/dejandjenic/go-gin-sample/application/entities"

	"cloud.google.com/go/firestore"
	"github.com/rs/zerolog/log"
	"google.golang.org/api/iterator"
)

type TodoRepository struct {
	Db     *firestore.Client
	Config configuration.Config
}

type ITodoRepository interface {
	CreateTodo(ctx context.Context, entity entities.Todo) (string, error)
	UpdateTodo(ctx context.Context, id string, item entities.Todo) error
	DeleteTodo(ctx context.Context, id string) error
	ListTodo(ctx context.Context) ([]entities.Todo, error)
	GetDetail(ctx context.Context, id string) (entities.Todo, error)
}

func (repository TodoRepository) CreateTodo(ctx context.Context, entity entities.Todo) (string, error) {
	docRef := repository.Db.Collection(repository.Config.FirestorePrefix + "Gos").NewDoc()
	_, err := docRef.Set(ctx, entity)
	log.Info().Msgf("%v", docRef)
	id := docRef.ID
	return id, err
}

func (r TodoRepository) UpdateTodo(ctx context.Context, id string, item entities.Todo) error {
	snapshot := r.Db.Collection(r.Config.FirestorePrefix + "Gos").Doc(id)

	//_, err := snapshot.Update(ctx,
	// []firestore.Update{
	// {
	// 	Path:  "name",
	// 	Value: "test",
	// },
	//	}
	//)

	_, err := snapshot.Set(ctx, item)

	return err
}

func (r TodoRepository) DeleteTodo(ctx context.Context, id string) error {
	snapshot := r.Db.Collection(r.Config.FirestorePrefix + "Gos").Doc(id)

	_, err := snapshot.Delete(ctx)

	return err
}

func (r TodoRepository) ListTodo(ctx context.Context) ([]entities.Todo, error) {
	data := []entities.Todo{}
	iter := r.Db.Collection(r.Config.FirestorePrefix + "Gos").Documents(ctx)
	for {
		doc, err := iter.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return nil, err
		}
		var d entities.Todo
		doc.DataTo(&d)
		d.ID = doc.Ref.ID
		data = append(data, d)
	}
	log.Info().Any("data", data).Msg("repository return")
	return data, nil
}

func (r TodoRepository) GetDetail(ctx context.Context, id string) (entities.Todo, error) {
	snapshot := r.Db.Collection(r.Config.FirestorePrefix + "Gos").Doc(id)

	var item entities.Todo
	a, err := snapshot.Get(ctx)
	if err != nil {
		return entities.Todo{}, err
	}

	if a.Exists() {
		a.DataTo(&item)
		item.ID = a.Ref.ID
		return item, nil
	}

	return entities.Todo{}, nil

}
