
#### consul 服务使用官方的 Consul 镜像(hashicorp/Consul)，监听在 8500 端口，并将其 UI 端口映射到宿主机的 8500 端口，以便可以访问 Consul 的 Web UI。
#### nginx 服务使用官方的 Nginx 镜像，并映射本地的 default.conf 文件和 templates 目录到容器内部，以便 Nginx 可以加载动态生成的配置文件。
#### consul-template 服务使用 HashiCorp 提供的 consul-template 镜像，配置了 Consul 地址和模板文件的路径，以便监听 Consul 中服务的变化，并根据模板文件生成 Nginx 配置，最后触发 Nginx 的重载操作。
#### vhost.conf: 默认的 Nginx 配置文件，该文件将由 consul-template 生成。
#### templates: 存放 Nginx 配置模板文件的目录。

#### 模板中，{{ range services }} 遍历了 Consul 中所有服务， {{if service .Name}} 则筛选出所有有效在consul注册的服务（你需要根据实际情况进行调整比如增加筛选条件，确保只选择你想要负载均衡的服务）。然后，将这些服务的地址和端口动态地写入 upstream 块中，最终生成完整的 Nginx 配置文件。

#### 启动
#### docker-compose up -d

#### Consul 将会监听服务的注册和注销，consul-template 将会根据模板文件动态生成 Nginx 配置，实现基于服务注册的动态负载均衡配置。新增服务只需在 Consul 中注册，并且模板文件会自动更新 Nginx 的配置，无需手动修改和重启 Nginx。

service-server目录下的main.go时用于测试的consul服务注册演示

All codes can be edit as you want, if you want update this code, please let me know, Thanks  
Code from clark zhu, contact email : chenyusolar@gmail.com
