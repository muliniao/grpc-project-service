package service

import (
	"context"
	"github.com/opentracing/opentracing-go"
	uuid "github.com/satori/go.uuid"
	pb "learning/grpc-project-service/api/gen/go/platform/v1"
	"learning/grpc-project-service/internal/model"
)

type ProjectService interface {
	CreateProject(context.Context, *pb.CreateProjectRequest) (*pb.CreateProjectResponse, error)
	ListProjects(context.Context, *pb.ListProjectsRequest) (*pb.ListProjectsResponse, error)
	GetProject(context.Context, *pb.GetProjectRequest) (*pb.GetProjectResponse, error)
	UpdateProject(context.Context, *pb.UpdateProjectRequest) (*pb.UpdateProjectResponse, error)
	DeleteProject(context.Context, *pb.DeleteProjectRequest) (*pb.DeleteProjectResponse, error)
}

func (s *service) CreateProject(ctx context.Context, req *pb.CreateProjectRequest) (*pb.CreateProjectResponse, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "Service::CreateProject")
	defer span.Finish()

	project, err := model.NewProject(req)
	if err != nil {
		return nil, err
	}

	//if err = s.repository.CreateProject(ctx, project); err != nil {
	//	return nil, err
	//}
	//
	//project, err = s.repository.GetProject(ctx, util.WithID(project.ID))
	//if err != nil {
	//	return nil, err
	//}
	//
	//// TODO: call resource api, user_id can be decode from user token
	//resp, err := s.userResource.GetUser(ctx, &user.GetUserRequest{UserId: uuid.NewV4().String()})
	//if err != nil {
	//	return nil, err
	//}
	//
	//u := resp.GetUser()
	//if (u != nil) {
	//	project.UserName = u.Name
	//}

	return project.ToCreateProjectResponse()
}

func (s *service) ListProjects(ctx context.Context, req *pb.ListProjectsRequest) (*pb.ListProjectsResponse, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "Service::ListProjects")
	defer span.Finish()

	filter := model.NewListProjectFilter(req)

	//projects, count, err := s.repository.ListProjects(ctx, filter)
	//if err != nil {
	//	return nil, err
	//}

	projects := []*model.Project{
		{
			ID:   uuid.NewV4(),
			Name: "Project 1",
		},
	}

	count := int64(1)

	return model.NewPageProjects(projects, count, filter.Limit, filter.Offset).ToListProjectsAPI()

}

func (s *service) GetProject(ctx context.Context, req *pb.GetProjectRequest) (*pb.GetProjectResponse, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "Service::GetProject")
	defer span.Finish()

	//project, err := s.repository.GetProject(ctx, util.WithID(uuid.FromStringOrNil(req.ProjectId)))
	//if err != nil {
	//	return nil, err
	//}

	project := &model.Project{
		ID:   uuid.FromStringOrNil(req.ProjectId),
		Name: "Project 1",
	}

	return project.ToGetProjectResponse()
}

func (s *service) UpdateProject(ctx context.Context, req *pb.UpdateProjectRequest) (*pb.UpdateProjectResponse, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "Service::UpdateProject")
	defer span.Finish()

	//if err := s.repository.UpdateProject(ctx,
	//	util.WithID(uuid.FromStringOrNil(req.ProjectId)),
	//	util.WithUpdate(bson.M{"name": req.GetBody().Name}),
	//); err != nil {
	//	return nil, err
	//}
	//
	//project, err := s.repository.GetProject(ctx, util.WithID(uuid.FromStringOrNil(req.ProjectId)))
	//if err != nil {
	//	return nil, err
	//}

	project := &model.Project{
		ID:   uuid.FromStringOrNil(req.ProjectId),
		Name: "Project 1",
	}

	return project.ToUpdateProjectResponse()
}

func (s *service) DeleteProject(ctx context.Context, req *pb.DeleteProjectRequest) (*pb.DeleteProjectResponse, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "Service::DeleteProject")
	defer span.Finish()

	//if err := s.repository.DeleteProject(ctx, util.WithID(uuid.FromStringOrNil(req.ProjectId))); err != nil {
	//	return nil, err
	//}

	return new(pb.DeleteProjectResponse), nil
}
