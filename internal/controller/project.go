package controller

import (
	"context"
	"google.golang.org/grpc/codes"

	"github.com/go-ozzo/ozzo-validation/is"
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/opentracing/opentracing-go"

	"google.golang.org/grpc/status"
	pb "learning/grpc-project-service/api/gen/go/platform/v1"
)

type ProjectController interface {
	pb.ProjectAPIServer
}

func (c controller) CreateProject(ctx context.Context, req *pb.CreateProjectRequest) (*pb.CreateProjectResponse, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "Controller::CreateProject")
	defer span.Finish()

	err := validation.ValidateStruct(req, validation.Field(&req.Name, validation.Required))
	if err != nil {
		return nil, status.New(codes.InvalidArgument, err.Error()).Err()
	}

	return c.service.CreateProject(ctx, req)
}

func (c controller) ListProjects(ctx context.Context, req *pb.ListProjectsRequest) (*pb.ListProjectsResponse, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "Controller::ListProjects")
	defer span.Finish()

	err := validation.ValidateStruct(req,
		validation.Field(&req.Offset, validation.Min(0)),
		validation.Field(&req.Limit, validation.Min(0), validation.Max(100)),
	)
	if err != nil {
		return nil, status.New(codes.InvalidArgument, err.Error()).Err()
	}

	return c.service.ListProjects(ctx, req)
}

func (c controller) GetProject(ctx context.Context, req *pb.GetProjectRequest) (*pb.GetProjectResponse, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "Controller::GetProject")
	defer span.Finish()

	err := validation.ValidateStruct(req, validation.Field(&req.ProjectId, validation.Required, is.UUID))
	if err != nil {
		return nil, status.New(codes.InvalidArgument, err.Error()).Err()
	}

	return c.service.GetProject(ctx, req)
}

func (c controller) UpdateProject(ctx context.Context, req *pb.UpdateProjectRequest) (*pb.UpdateProjectResponse, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "Controller::UpdateProject")
	defer span.Finish()

	err := validation.ValidateStruct(req, validation.Field(&req.ProjectId, validation.Required, is.UUID))
	if err != nil {
		return nil, status.New(codes.InvalidArgument, err.Error()).Err()
	}

	return c.service.UpdateProject(ctx, req)
}

func (c controller) DeleteProject(ctx context.Context, req *pb.DeleteProjectRequest) (*pb.DeleteProjectResponse, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "Controller::DeleteProject")
	defer span.Finish()

	err := validation.ValidateStruct(req, validation.Field(&req.ProjectId, validation.Required, is.UUID))
	if err != nil {
		return nil, status.New(codes.InvalidArgument, err.Error()).Err()
	}

	return c.service.DeleteProject(ctx, req)
}
