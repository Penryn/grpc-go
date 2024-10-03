/*
 *
 * 版权所有 2015 gRPC 作者。
 *
 * 根据 Apache 许可证 2.0 版（“许可证”）许可；
 * 除非遵守许可证，否则您不得使用此文件。
 * 您可以在以下网址获取许可证副本：
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * 除非适用法律要求或书面同意，否则按“原样”分发的许可证软件
 * 没有任何明示或暗示的担保或条件。
 * 请参阅许可证以了解管理权限和限制的特定语言。
 *
 */

// main 包实现了一个简单的 gRPC 客户端，演示如何使用 gRPC-Go 库
// 执行一元、客户端流、服务器流和全双工 RPC。
//
// 它与 route guide 服务交互，其定义可以在 routeguide/route_guide.proto 中找到。
package main

import (
	"context"
	"flag"
	"io"
	"log"
	"math/rand"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/examples/data"
	pb "google.golang.org/grpc/examples/route_guide/routeguide"
)

var (
	tls                = flag.Bool("tls", false, "如果为 true，则使用 TLS 连接，否则使用普通 TCP")
	caFile             = flag.String("ca_file", "", "包含 CA 根证书文件的文件")
	serverAddr         = flag.String("addr", "localhost:50051", "服务器地址，格式为 host:port")
	serverHostOverride = flag.String("server_host_override", "x.test.example.com", "用于验证 TLS 握手返回的主机名的服务器名称")
)

// printFeature 获取给定点的特征。
func printFeature(client pb.RouteGuideClient, point *pb.Point) {
	log.Printf("获取点 (%d, %d) 的特征", point.Latitude, point.Longitude)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	feature, err := client.GetFeature(ctx, point)
	if err != nil {
		log.Fatalf("client.GetFeature 失败: %v", err)
	}
	log.Println(feature)
}

// printFeatures 列出给定矩形范围内的所有特征。
func printFeatures(client pb.RouteGuideClient, rect *pb.Rectangle) {
	log.Printf("查找 %v 范围内的特征", rect)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	stream, err := client.ListFeatures(ctx, rect)
	if err != nil {
		log.Fatalf("client.ListFeatures 失败: %v", err)
	}
	for {
		feature, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatalf("client.ListFeatures 失败: %v", err)
		}
		log.Printf("特征: 名称: %q, 点:(%v, %v)", feature.GetName(),
			feature.GetLocation().GetLatitude(), feature.GetLocation().GetLongitude())
	}
}

// runRecordRoute 发送一系列点到服务器，并期望从服务器获得 RouteSummary。
func runRecordRoute(client pb.RouteGuideClient) {
	// 创建随机数量的随机点
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	pointCount := int(r.Int31n(100)) + 2 // 至少遍历两个点
	var points []*pb.Point
	for i := 0; i < pointCount; i++ {
		points = append(points, randomPoint(r))
	}
	log.Printf("遍历 %d 个点。", len(points))
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	stream, err := client.RecordRoute(ctx)
	if err != nil {
		log.Fatalf("client.RecordRoute 失败: %v", err)
	}
	for _, point := range points {
		if err := stream.Send(point); err != nil {
			log.Fatalf("client.RecordRoute: stream.Send(%v) 失败: %v", point, err)
		}
	}
	reply, err := stream.CloseAndRecv()
	if err != nil {
		log.Fatalf("client.RecordRoute 失败: %v", err)
	}
	log.Printf("路线总结: %v", reply)
}

// runRouteChat 在发送各种位置的笔记时接收一系列路线笔记。
func runRouteChat(client pb.RouteGuideClient) {
	notes := []*pb.RouteNote{
		{Location: &pb.Point{Latitude: 0, Longitude: 1}, Message: "第一条消息"},
		{Location: &pb.Point{Latitude: 0, Longitude: 2}, Message: "第二条消息"},
		{Location: &pb.Point{Latitude: 0, Longitude: 3}, Message: "第三条消息"},
		{Location: &pb.Point{Latitude: 0, Longitude: 1}, Message: "第四条消息"},
		{Location: &pb.Point{Latitude: 0, Longitude: 2}, Message: "第五条消息"},
		{Location: &pb.Point{Latitude: 0, Longitude: 3}, Message: "第六条消息"},
	}
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	stream, err := client.RouteChat(ctx)
	if err != nil {
		log.Fatalf("client.RouteChat 失败: %v", err)
	}
	waitc := make(chan struct{})
	go func() {
		for {
			in, err := stream.Recv()
			if err == io.EOF {
				// 读取完成。
				close(waitc)
				return
			}
			if err != nil {
				log.Fatalf("client.RouteChat 失败: %v", err)
			}
			log.Printf("在点(%d, %d)收到消息 %s", in.Location.Latitude, in.Location.Longitude, in.Message)
		}
	}()
	for _, note := range notes {
		if err := stream.Send(note); err != nil {
			log.Fatalf("client.RouteChat: stream.Send(%v) 失败: %v", note, err)
		}
	}
	stream.CloseSend()
	<-waitc
}

func randomPoint(r *rand.Rand) *pb.Point {
	lat := (r.Int31n(180) - 90) * 1e7
	long := (r.Int31n(360) - 180) * 1e7
	return &pb.Point{Latitude: lat, Longitude: long}
}

func main() {
	flag.Parse()
	var opts []grpc.DialOption
	if *tls {
		if *caFile == "" {
			*caFile = data.Path("x509/ca_cert.pem")
		}
		creds, err := credentials.NewClientTLSFromFile(*caFile, *serverHostOverride)
		if err != nil {
			log.Fatalf("创建 TLS 凭证失败: %v", err)
		}
		opts = append(opts, grpc.WithTransportCredentials(creds))
	} else {
		opts = append(opts, grpc.WithTransportCredentials(insecure.NewCredentials()))
	}

	conn, err := grpc.Dial(*serverAddr, opts...)
	if err != nil {
		log.Fatalf("拨号失败: %v", err)
	}
	defer conn.Close()
	client := pb.NewRouteGuideClient(conn)

	// 查找有效特征
	printFeature(client, &pb.Point{Latitude: 409146138, Longitude: -746188906})

	// 特征缺失。
	printFeature(client, &pb.Point{Latitude: 0, Longitude: 0})

	// 查找 40, -75 和 42, -73 之间的特征。
	printFeatures(client, &pb.Rectangle{
		Lo: &pb.Point{Latitude: 400000000, Longitude: -750000000},
		Hi: &pb.Point{Latitude: 420000000, Longitude: -730000000},
	})

	// 记录路线
	runRecordRoute(client)

	// 路线聊天
	runRouteChat(client)
}
