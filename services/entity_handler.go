package services

import (
	"context"
	"grpc_with_some_methods/protos/go"

	"gorm.io/gorm"
)

type EntityHandler struct {
	_go.UnimplementedEntityServiceServer
	DB *gorm.DB
}

func (eh *EntityHandler) Create(ctx context.Context, entity *_go.Entity) (*_go.Entity, error) {
	if result := eh.DB.WithContext(ctx).Create(entity); result.Error != nil {
		return nil, result.Error
	}
	return entity, nil
}

func (eh *EntityHandler) Update(ctx context.Context, entity *_go.Entity) (*_go.Entity, error) {
	if result := eh.DB.WithContext(ctx).Save(entity); result.Error != nil {
		return nil, result.Error
	}
	return entity, nil
}

func (eh *EntityHandler) Delete(ctx context.Context, entityID *_go.EntityID) (*_go.DeleteResponse, error) {
	result := eh.DB.WithContext(ctx).Delete(&_go.Entity{}, entityID.Id)
	if result.Error != nil {
		return nil, result.Error
	}
	return &_go.DeleteResponse{Success: result.RowsAffected > 0}, nil
}

func (eh *EntityHandler) Get(ctx context.Context, entityID *_go.EntityID) (*_go.Entity, error) {
	entity := &_go.Entity{}
	if result := eh.DB.WithContext(ctx).First(entity, entityID.Id); result.Error != nil {
		return nil, result.Error
	}
	return entity, nil
}

func (eh *EntityHandler) Search(searchRequest *_go.SearchRequest, stream _go.EntityService_SearchServer) error {
	var entities []_go.Entity
	if result := eh.DB.Where("name LIKE ?", "%"+searchRequest.Query+"%").Find(&entities); result.Error != nil {
		return result.Error
	}

	for _, entity := range entities {
		if err := stream.Send(&entity); err != nil {
			return err
		}
	}

	return nil
}
