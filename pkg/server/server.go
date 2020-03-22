package server

import (
	"context"
	"errors"
	"fmt"
	"net"

	"github.com/covidhub/covid-socialgraph/pkg/database"
	pb "github.com/covidhub/covid-socialgraph/pkg/server/socialgraph"
	"github.com/golang/protobuf/ptypes"
	"github.com/golang/protobuf/ptypes/wrappers"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

type SocialgraphServer struct {
	pb.UnimplementedSocialGraphServer

	url  string
	opts []grpc.ServerOption

	server *grpc.Server
	db     *database.Client
	logger *zap.SugaredLogger
}

func New(url string, tls bool, certFile, keyFile string, db *database.Client, logger *zap.SugaredLogger) (*SocialgraphServer, error) {
	socialgraphServer := &SocialgraphServer{
		url:    url,
		db:     db,
		logger: logger,
	}

	if tls {
		if certFile == "" {
			return nil, errors.New("No TLS cert file set")
		}
		if keyFile == "" {
			return nil, errors.New("No TLS cert file set")
		}
		creds, err := credentials.NewServerTLSFromFile(certFile, keyFile)
		if err != nil {
			return nil, fmt.Errorf("Failed to generate credentials %v", err)
		}
		socialgraphServer.opts = []grpc.ServerOption{grpc.Creds(creds)}
	}

	return socialgraphServer, nil
}

func (s *SocialgraphServer) Serve() {
	lis, err := net.Listen("tcp", s.url)
	if err != nil {
		s.logger.Fatalw(
			"Failed to listen",
			"url", s.url,
			"err", err,
		)
	}

	s.server = grpc.NewServer(s.opts...)

	pb.RegisterSocialGraphServer(s.server, s)

	s.server.Serve(lis)
}

func (s *SocialgraphServer) Close() {
	s.server.GracefulStop()
}

func (s *SocialgraphServer) GetProfile(context.Context, *wrappers.StringValue) (*pb.Profile, error) {
	profile := &pb.Profile{
		Id:          "id1234",
		Age:         42,
		AgeGroup:    pb.Profile_MIDDLE_AGED_FOURTIES,
		CovidStatus: pb.Profile_NOT_INFECTED,
		MedicalConditions: []pb.Profile_MedicalConditions{
			pb.Profile_HEART_DISEASE,
			pb.Profile_RESPIRATORY_DISEASE,
			pb.Profile_AUTOIMMUNE_DISEASE,
			pb.Profile_DIABETES,
		},
	}
	return profile, nil
}

func (s *SocialgraphServer) AddProfile(ctx context.Context, profile *pb.Profile) (*wrappers.StringValue, error) {
	return &wrappers.StringValue{Value: "id1234"}, nil
}

func (s *SocialgraphServer) UpdateProfile(ctx context.Context, profile *pb.Profile) (*wrappers.StringValue, error) {
	return &wrappers.StringValue{Value: "id1234"}, nil
}

func (s *SocialgraphServer) GetContacts(id *wrappers.StringValue, stream pb.SocialGraph_GetContactsServer) error {
	pubProfile := &pb.PublicProfile{
		Id:         "id5678",
		FirstName:  "Max",
		LastName:   "Mustermann",
		MiddleName: "Middlename",
	}

	contact := &pb.Contact{
		UserId:         "id1234",
		ContactWith:    pubProfile,
		TimeOfContact:  ptypes.TimestampNow(),
		PointOfContact: &pb.Point{Longitude: 1.0, Latitude: 1.0},
	}

	if err := stream.Send(contact); err != nil {
		return err
	}

	return nil
}

func (s *SocialgraphServer) AddContact(ctx context.Context, contact *pb.Contact) (*pb.PublicProfile, error) {
	pubProfile := &pb.PublicProfile{
		Id:         "id1234",
		FirstName:  "Max",
		LastName:   "Mustermann",
		MiddleName: "Middlename",
	}
	return pubProfile, nil
}

func (s *SocialgraphServer) GetBonds(id *wrappers.StringValue, stream pb.SocialGraph_GetBondsServer) error {
	pubProfile := &pb.PublicProfile{
		Id:         "id1234",
		FirstName:  "Max",
		LastName:   "Mustermann",
		MiddleName: "Middlename",
	}

	if err := stream.Send(pubProfile); err != nil {
		return err
	}

	return nil
}

func (s *SocialgraphServer) AddBond(ctx context.Context, id *wrappers.StringValue) (*pb.PublicProfile, error) {
	pubProfile := &pb.PublicProfile{
		Id:         "id1234",
		FirstName:  "Max",
		LastName:   "Mustermann",
		MiddleName: "Middlename",
	}

	return pubProfile, nil
}

func (s *SocialgraphServer) DeleteBond(ctx context.Context, id *wrappers.StringValue) (*pb.PublicProfile, error) {
	pubProfile := &pb.PublicProfile{
		Id:         "id1234",
		FirstName:  "Max",
		LastName:   "Mustermann",
		MiddleName: "Middlename",
	}

	return pubProfile, nil
}

func (s *SocialgraphServer) GetReport(ctx context.Context, id *wrappers.StringValue) (*pb.Report, error) {
	infection := &pb.Infection{
		Type:       pb.Infection_CONTACT,
		ReportedAt: ptypes.TimestampNow(),
	}

	recommendation := &pb.Recommendation{}

	report := &pb.Report{
		Infection:      []*pb.Infection{infection},
		Recommendation: recommendation,
	}

	return report, nil
}

func (s *SocialgraphServer) GetSummary(ctx context.Context, id *wrappers.StringValue) (*pb.Summary, error) {
	summary := &pb.Summary{
		NumContacts:              42,
		TotalInfectedContacts:    9,
		RecentlyInfectedContacts: 2,
		NumCloseContacts:         6,
		InfectedCloseContacts:    0,
	}
	return summary, nil
}
