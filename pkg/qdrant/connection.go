/**
 * @Author Nil
 * @Description pkg/qdrant/connection.go
 * @Date 2023/3/30 11:19
 **/

package qdrant

import (
	"context"
	"github.com/ha5ky/hu5ky-bot/pkg/config"
	"github.com/ha5ky/hu5ky-bot/pkg/logger"
	pb "github.com/qdrant/go-client/qdrant"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type ClientSet struct {
	conn *grpc.ClientConn
	pb.QdrantClient
	pb.CollectionsClient
	pb.PointsClient
}

func NewClientSet(ctx context.Context) (clientSet *ClientSet, err error) {
	clientSet = new(ClientSet)
	clientSet.conn = new(grpc.ClientConn)
	clientSet.conn, err = NewConn()
	if err != nil {
		logger.Fatalf("connection error: %v", err)
		return
	}
	clientSet.QdrantClient, err = NewQdRantClient(ctx, clientSet.conn)
	if err != nil {
		logger.Fatalf("qdrantClient init error: %v", err)
		return
	}

	clientSet.CollectionsClient = NewCollectionClient(clientSet.conn)
	clientSet.PointsClient = NewPointsClient(clientSet.conn)
	return
}

func (c *ClientSet) ConnClose() error {
	return c.conn.Close()
}

func NewConn() (conn *grpc.ClientConn, err error) {
	addr := config.SysCache.DB.QdRant.Host + ":" + config.SysCache.DB.QdRant.Port
	conn, err = grpc.Dial(
		addr,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		logger.Errorf("did not connect: %v", err)
		return
	}
	return
}

func NewQdRantClient(ctx context.Context, conn *grpc.ClientConn) (qdrantClient pb.QdrantClient, err error) {
	qdrantClient = pb.NewQdrantClient(conn)
	healthCheckResult, err := qdrantClient.HealthCheck(ctx, &pb.HealthCheckRequest{})
	if err != nil {
		logger.Fatalf("Could not get health: %v", err)
		return
	}
	logger.Infof("QdRant version: %s", healthCheckResult.GetVersion())
	return qdrantClient, nil
}

func NewCollectionClient(conn *grpc.ClientConn) (collectionClient pb.CollectionsClient) {
	collectionClient = pb.NewCollectionsClient(conn)
	return
}

func NewPointsClient(conn *grpc.ClientConn) (pointsClient pb.PointsClient) {
	pointsClient = pb.NewPointsClient(conn)
	return
}
