// Copyright 2022 CloudWeGo Authors.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"sync"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/app/server"
	"github.com/cloudwego/hertz/pkg/app/server/registry"
	"github.com/cloudwego/hertz/pkg/common/utils"
	"github.com/cloudwego/hertz/pkg/protocol/consts"
	consulapi "github.com/hashicorp/consul/api"
	"github.com/hertz-contrib/registry/consul"
)

var (
	wg      sync.WaitGroup
	localIP = "192.168.0.98"
)

func main() {
	config := consulapi.DefaultConfig()
	config.Address = "127.0.0.1:8500"
	consulClient, err := consulapi.NewClient(config)
	if err != nil {
		log.Fatal(err)
		return
	}

	wg.Add(2)

	tags := map[string]string{
		"app":     "wattapi",
		"version": "v1",
		"env":     "dev",
		"role":    "nginx",
	}

	go func() {
		defer wg.Done()
		addr := net.JoinHostPort(localIP, "8886")
		r := consul.NewConsulRegister(consulClient)
		tags["app"] = "demo1"
		h := server.Default(
			server.WithHostPorts(addr),
			server.WithRegistry(r, &registry.Info{
				ServiceName: "demo1",
				Addr:        utils.NewNetAddr("tcp", addr),
				Weight:      10,
				Tags:        tags,
			}),
		)

		h.GET("/api/demo1/ping", func(c context.Context, ctx *app.RequestContext) {
			fmt.Println("call from demo1")
			ctx.JSON(consts.StatusOK, utils.H{"ping": "pong1"})
		})
		h.Spin()
	}()
	go func() {
		defer wg.Done()
		addr := net.JoinHostPort(localIP, "8887")
		r := consul.NewConsulRegister(consulClient)
		tags["app"] = "demo2"
		var h = server.Default(
			server.WithHostPorts(addr),
			server.WithRegistry(r, &registry.Info{
				ServiceName: "demo2",
				Addr:        utils.NewNetAddr("tcp", addr),
				Weight:      10,
				Tags:        tags,
			}),
		)
		h.GET("/api/demo2/ping", func(c context.Context, ctx *app.RequestContext) {
			fmt.Println("call from demo2")
			ctx.JSON(consts.StatusOK, utils.H{"ping": "pong2"})
		})
		h.Spin()
	}()

	wg.Wait()
}
