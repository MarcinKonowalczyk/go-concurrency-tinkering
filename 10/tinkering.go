package main

import (
	"context"
	"log"
	"math/rand/v2"
	"time"

	"go-concurrency-tinkering/10/pb"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	pbtime "google.golang.org/protobuf/types/known/timestamppb"
)

// Timestamp converts time.Time to protobuf *Timestamp
func Timestamp(t time.Time) *pbtime.Timestamp {
    return &pbtime.Timestamp {
        Seconds: t.Unix(),
        Nanos: int32(t.Nanosecond()),
    }
}

func dummyData() []*pb.Metric {
    const size = 1000
    out := make([]*pb.Metric, size)
    t := time.Date(2020, 5, 22, 14, 13, 11, 0, time.UTC)
    for i := 0; i < size; i++ {
         m := pb.Metric{
              Time: Timestamp(t),
              Name: "CPU",
              // Normally we're below 40% CPU utilization
              Value: rand.Float64() * 40,
         }
         out[i] = &m
         t = t.Add(time.Second)
    }
    // Create some outliers
    out[7].Value = 97.3
    out[113].Value = 92.1
    out[835].Value = 93.2
    return out
}

func main() {
	addr := "localhost:9999"

	options := grpc.WithTransportCredentials(insecure.NewCredentials())
	conn, err := grpc.NewClient(addr, options)
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	client := pb.NewOutliersClient(conn)
	req := pb.OutliersRequest{
		Metrics: dummyData(),
	}

	resp, err := client.Detect(context.Background(), &req)
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("outliers at: %v", resp.Indices)
}