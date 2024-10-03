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
 * 除非适用法律要求或书面同意，否则按“原样”分发的软件
 * 不提供任何明示或暗示的担保或条件。
 * 请参阅许可证了解管理权限和限制的具体语言。
 *
 */

// main 包实现了 Greeter 服务的客户端。
package main

import (
	"context"
	"flag"
	"log"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	pb "google.golang.org/grpc/examples/helloworld/helloworld"
)

const (
	defaultName = "grpc"
)

var (
	addr = flag.String("addr", "localhost:50051", "连接的地址")
	name = flag.String("name", defaultName, "问候的名字")
)

func main() {
	flag.Parse()
	// 建立与服务器的连接。
	// 这里我们使用 grpc.WithInsecure() 选项，这样我们就可以不使用 TLS 连接。
	conn, err := grpc.Dial(*addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("连接失败: %v", err)
	}
	defer conn.Close()
	// 创建一个 Greeter 客户端。
	c := pb.NewGreeterClient(conn)

	// 联系服务器并打印其响应。
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	r, err := c.SayHello(ctx, &pb.HelloRequest{Name: *name})
	if err != nil {
		log.Fatalf("无法问候: %v", err)
	}
	log.Printf("问候: %s", r.GetMessage())

	r, err = c.SayHelloAgain(ctx, &pb.HelloRequest{Name: *name})
	if err != nil {
		log.Fatalf("无法再次问候: %v", err)
	}
	log.Printf("问候: %s", r.GetMessage())
}
