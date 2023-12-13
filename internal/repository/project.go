package repository

//
//import (
//	"context"
//	"github.com/opentracing/opentracing-go"
//	"go.mongodb.org/mongo-driver/bson"
//	"go.mongodb.org/mongo-driver/mongo/options"
//	"google.golang.org/grpc/codes"
//	"google.golang.org/grpc/status"
//	"learning/grpc-project-service/internal/model"
//)
//
//const collectionProject = "project"
//
//type ProjectRepository interface {
//	CreateProject(context.Context, *model.Project) error
//	GetProject(context.Context, bson.M) (*model.Project, error)
//	UpdateProject(ctx context.Context, filter bson.M, update bson.M) error
//	DeleteProject(ctx context.Context, filter bson.M) error
//	ListProjects(context.Context, *model.ListProjectsFilter) ([]*model.Project, int64, error)
//}
//
//func (r *repository) UpdateProject(ctx context.Context, filter bson.M, update bson.M) error {
//	span, ctx := opentracing.StartSpanFromContext(ctx, "Repository::UpdateProject")
//	defer span.Finish()
//
//	result, err := r.MongoDatabase(ctx).Collection(collectionProject).UpdateOne(ctx, filter, update)
//	if result != nil && result.MatchedCount == 0 {
//		return status.Errorf(codes.NotFound, "not found")
//	}
//	return err
//}
//
//func (r *repository) DeleteProject(ctx context.Context, filter bson.M) error {
//	span, ctx := opentracing.StartSpanFromContext(ctx, "Repository::DeleteProject")
//	defer span.Finish()
//
//	result, err := r.MongoDatabase(ctx).Collection(collectionProject).DeleteMany(ctx, filter)
//	if result != nil && result.DeletedCount == 0 {
//		return status.Error(codes.NotFound, "not found")
//	}
//	return err
//}
//
//func (r *repository) GetProject(ctx context.Context, filter bson.M) (project *model.Project, err error) {
//	span, ctx := opentracing.StartSpanFromContext(ctx, "Repository::GetProject")
//	defer span.Finish()
//
//	err = r.MongoDatabase(ctx).Collection(collectionProject).FindOne(ctx, filter).Decode(&project)
//	return
//}
//
//func (r *repository) CreateProject(ctx context.Context, project *model.Project) error {
//	span, ctx := opentracing.StartSpanFromContext(ctx, "Repository::CreateProject")
//	defer span.Finish()
//
//	_, err := r.MongoDatabase(ctx).Collection(collectionProject).InsertOne(ctx, *project)
//	return err
//}
//
//func (r *repository) ListProjects(ctx context.Context, filter *model.ListProjectsFilter) ([]*model.Project, int64, error) {
//	span, ctx := opentracing.StartSpanFromContext(ctx, "Repository::ListProjects")
//	defer span.Finish()
//
//	count, err := r.MongoDatabase(ctx).Collection(collectionProject).CountDocuments(ctx, filter.GetFilter())
//	if err != nil || count == 0 {
//		return nil, 0, err
//	}
//
//	cursor, err := r.MongoDatabase(ctx).Collection(collectionProject).Find(ctx, filter.GetFilter(),
//		options.Find().SetSort(filter.GetSort()).SetSkip(filter.Offset).SetLimit(filter.Limit))
//	if err != nil {
//		return nil, 0, err
//	}
//
//	var projects []*model.Project
//	if err = cursor.All(ctx, &projects); err != nil {
//		return nil, 0, err
//	}
//
//	return projects, count, nil
//}
