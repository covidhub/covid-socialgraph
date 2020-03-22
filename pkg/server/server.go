package server

import (
	"context"
	"errors"
	"fmt"
	"net"

	pb "github.com/covidhub/covid-socialgraph/pkg/server/socialgraph"
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
	logger *zap.SugaredLogger
}

func New(host string, port int, tls bool, certFile, keyFile string, logger *zap.SugaredLogger) (*SocialgraphServer, error) {
	socialgraphServer := &SocialgraphServer{
		url:    fmt.Sprintf("%s:%d", host, port),
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
	return nil, nil
}

func (s *SocialgraphServer) AddProfile(context.Context, *pb.Profile) (*wrappers.StringValue, error) {
	return nil, nil
}

func (s *SocialgraphServer) UpdateProfile(context.Context, *pb.Profile) (*wrappers.StringValue, error) {
	return nil, nil
}

func (s *SocialgraphServer) GetContacts(*wrappers.StringValue, pb.SocialGraph_GetContactsServer) error {
	return nil
}

func (s *SocialgraphServer) AddContact(context.Context, *pb.Contact) (*pb.PublicProfile, error) {
	return nil, nil
}

func (s *SocialgraphServer) GetBonds(*wrappers.StringValue, pb.SocialGraph_GetBondsServer) error {
	return nil
}

func (s *SocialgraphServer) AddBond(context.Context, *wrappers.StringValue) (*pb.PublicProfile, error) {
	return nil, nil
}

func (s *SocialgraphServer) DeleteBond(context.Context, *wrappers.StringValue) (*pb.PublicProfile, error) {
	return nil, nil
}

func (s *SocialgraphServer) GetReport(context.Context, *wrappers.StringValue) (*pb.Report, error) {
	return nil, nil
}

func (s *SocialgraphServer) GetSummary(context.Context, *wrappers.StringValue) (*pb.Summary, error) {
	return nil, nil
}
