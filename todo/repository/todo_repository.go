package repository

import (
	"errors"
	"os"

	"github.com/google/uuid"
	r "gopkg.in/rethinkdb/rethinkdb-go.v6"

	"go-rethinkdb/todo/models"
	"go-rethinkdb/utils"
)

// RethinkDBTodoRepository as provider
type RethinkDBTodoRepository interface {
	FindAll(keyword string, limit int, offset int) ([]*models.Todo, error)
	CountFindAll(keyword string) (int, error)
	FindById(id string) (*models.Todo, error)
	CountFindByID(id string) (int, error)
	Store(value *models.Todo) (*models.Todo, error)
	Update(id string, value *models.Todo) (*models.Todo, error)
	Delete(id string) error
}

// RethinkDBTodoRepositoryImpl represent RethinkDBTodoRepository implementation
type RethinkDBTodoRepositoryImpl struct {
	session *r.Session
}

// NewRethinkDBTodoRepository will create an object that represent the RethinkDBTodoRepository interface
func NewRethinkDBTodoRepository(session *r.Session) RethinkDBTodoRepository {
	return &RethinkDBTodoRepositoryImpl{
		session: session,
	}
}

// FindAll - find all todo
func (m *RethinkDBTodoRepositoryImpl) FindAll(keyword string, limit int, offset int) ([]*models.Todo, error) {
	cur, err := r.DB(os.Getenv("DB_NAME")).Table("todo").
		Filter(
			r.Row.Field("title").Match(keyword),
		).
		OrderBy(r.Desc("created_at")).
		Skip(offset).
		Limit(limit).
		Run(m.session)

	if err != nil {
		return nil, err
	}
	defer cur.Close()

	var results []*models.Todo
	err = cur.All(&results)
	if err != nil {
		return nil, err
	}

	return results, nil
}

// CountFindAll - count find all todo
func (m *RethinkDBTodoRepositoryImpl) CountFindAll(keyword string) (int, error) {
	cur, err := r.DB(os.Getenv("DB_NAME")).Table("todo").
		Count().
		Run(m.session)
	if err != nil {
		return 0, err
	}

	var row interface{}
	err = cur.One(&row)
	if err == r.ErrEmptyResult {
		return 0, err
	}
	if err != nil {
		return 0, err
	}

	count := int(row.(float64))

	return count, nil
}

// FindById - find todo by id
func (m *RethinkDBTodoRepositoryImpl) FindById(id string) (*models.Todo, error) {
	cur, err := r.DB(os.Getenv("DB_NAME")).Table("todo").
		Filter(r.Row.Field("id").Eq(id)).
		Run(m.session)
	if err != nil {
		return nil, err
	}
	defer cur.Close()

	var elem models.Todo
	err = cur.One(&elem)
	if err != nil {
		return nil, err
	}

	return &elem, nil
}

// CountFindByID - find count todo by id
func (m *RethinkDBTodoRepositoryImpl) CountFindByID(id string) (int, error) {
	cur, err := r.DB(os.Getenv("DB_NAME")).Table("todo").
		Filter(r.Row.Field("id").Eq(id)).
		Count().
		Run(m.session)
	if err != nil {
		return 0, err
	}

	var row interface{}
	err = cur.One(&row)
	if err == r.ErrEmptyResult {
		return 0, err
	}
	if err != nil {
		return 0, err
	}

	count := int(row.(float64))

	if count <= 0 {
		return 0, errors.New(utils.ErrREQLNotFound)
	}

	return count, nil
}

// Store - store todo
func (m *RethinkDBTodoRepositoryImpl) Store(value *models.Todo) (*models.Todo, error) {
	timeNow := utils.GetTimeNow()

	result := &models.Todo{
		ID:          uuid.New().String(),
		Title:       value.Title,
		Description: value.Description,
		CreatedAt:   timeNow,
		UpdatedAt:   timeNow,
	}

	_, err := r.DB(os.Getenv("DB_NAME")).Table("todo").
		Insert(map[string]interface{}{
			"id":          result.ID,
			"title":       result.Title,
			"description": result.Description,
			"created_at":  result.CreatedAt,
			"updated_at":  result.UpdatedAt,
		}).RunWrite(m.session)
	if err != nil {
		return nil, err
	}

	return result, nil
}

// Update - update todo by id
func (m *RethinkDBTodoRepositoryImpl) Update(id string, value *models.Todo) (*models.Todo, error) {
	timeNow := utils.GetTimeNow()

	result := &models.Todo{
		Title:       value.Title,
		Description: value.Description,
		UpdatedAt:   timeNow,
	}

	_, err := r.DB(os.Getenv("DB_NAME")).Table("todo").
		Filter(r.Row.Field("id").Eq(id)).
		Update(map[string]interface{}{
			"title":       result.Title,
			"description": result.Description,
			"updated_at":  result.UpdatedAt,
		}).RunWrite(m.session)
	if err != nil {
		return nil, err
	}

	return nil, nil
}

// Delete - delete todo by id
func (m *RethinkDBTodoRepositoryImpl) Delete(id string) error {
	_, err := r.DB(os.Getenv("DB_NAME")).Table("todo").
		Filter(r.Row.Field("id").Eq(id)).
		Delete().
		RunWrite(m.session)
	if err != nil {
		return err
	}

	return nil
}
