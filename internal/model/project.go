package model

import (
	uuid "github.com/satori/go.uuid"
	"go.mongodb.org/mongo-driver/bson"
	pb "learning/grpc-project-service/api/gen/go/platform/v1"
)

type Project struct {
	ID       uuid.UUID `bson:"_id"`
	Name     string    `bson:"name"`
	UserName string    `bson:"userName"`
}

func NewProject(req *pb.CreateProjectRequest) (*Project, error) {
	project := &Project{
		ID:   uuid.NewV4(),
		Name: req.GetName(),
	}
	return project, nil
}

func (p *Project) ToCreateProjectResponse() (*pb.CreateProjectResponse, error) {
	return &pb.CreateProjectResponse{Project: p.ToAPI()}, nil
}

func (p *Project) ToGetProjectResponse() (*pb.GetProjectResponse, error) {
	return &pb.GetProjectResponse{Project: p.ToAPI()}, nil
}

func (p *Project) ToUpdateProjectResponse() (*pb.UpdateProjectResponse, error) {
	return &pb.UpdateProjectResponse{Project: p.ToAPI()}, nil
}

func (p *Project) ToAPI() *pb.Project {
	return &pb.Project{
		Id:   p.ID.String(),
		Name: p.Name,
	}
}

type ListProjectsFilter struct {
	Filter
	Name    string
	OrderBy []string
}

func NewListProjectFilter(req *pb.ListProjectsRequest) *ListProjectsFilter {
	l := &ListProjectsFilter{Filter: DefaultFilter()}

	if req.Offset != 0 {
		l.Offset = req.Offset
	}

	if req.Limit != 0 {
		l.Limit = req.Limit
	}

	l.Name = req.Name
	l.OrderBy = req.OrderBy
	return l
}

func (l *ListProjectsFilter) GetFilter() bson.M {
	var ret = bson.M{}

	if len(l.Name) > 0 {
		ret["name"] = l.Name
	}
	return ret
}

func (l *ListProjectsFilter) GetSort() bson.M {
	var ret = bson.M{}
	if len(l.OrderBy) > 0 {

	}
	return ret
}

type PageProjects struct {
	Pagination
	Elements []*Project
}

func NewPageProjects(projects []*Project, count, limit, offset int64) *PageProjects {
	return &PageProjects{
		Pagination: Pagination{
			Limit:  limit,
			Offset: offset,
			Count:  count,
		},
		Elements: projects,
	}
}

func (p *PageProjects) ToListProjectsAPI() (*pb.ListProjectsResponse, error) {
	projects := &pb.ListProjectsResponse{
		Count:    p.Count,
		Offset:   p.Offset,
		Limit:    p.Limit,
		Elements: make([]*pb.Project, 0),
	}
	for _, e := range p.Elements {
		projects.Elements = append(projects.Elements, e.ToAPI())
	}
	return projects, nil
}
