package main

import (
	"fmt"

	"taobao/jackbao/user/domain/service"

	"taobao/jackbao/user/common"
	"taobao/jackbao/user/domain/repository"
	"taobao/jackbao/user/handler"
	user "taobao/jackbao/user/proto/user"

	"github.com/micro/go-micro/v2"
	log "github.com/micro/go-micro/v2/logger"
	"github.com/micro/go-micro/v2/registry"
	"github.com/micro/go-plugins/registry/consul/v2"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
)

func main() {
	//配置中心
	consulConfig, err := common.GetConsulConfig("127.0.0.1", 8500, "/micro/config")
	if err != nil {
		log.Error(err)
	}
	//注册中心
	consulRegistry := consul.NewRegistry(func(options *registry.Options) {
		options.Addrs = []string{
			"127.0.0.1:8500",
		}
	})

	//1 创建服务
	srv := micro.NewService(
		micro.Name("go.micro.service.user"),
		micro.Version("latest"),
		//设置地址和需要暴露的端口
		micro.Address("127.0.0.1:8081"),
		//添加consul作为注册中心
		micro.Registry(consulRegistry),
	)
	//2初始化服务
	srv.Init()
	//获取mysql的配置
	mysqlInfo := common.GetMysqlFromConsul(consulConfig, "mysql")
	//3. 初始化中间件 创建数据库连接
	//dsn := "root:root@tcp(localhost:3306)/micro?charset=utf8mb4&parseTime=True&loc=Local"
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		mysqlInfo.User,
		mysqlInfo.Pwd,
		mysqlInfo.Host,
		mysqlInfo.Port,
		mysqlInfo.Database)
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{NamingStrategy: schema.NamingStrategy{SingularTable: true}})
	if err != nil {
		log.Error(err)
	}
	rp := repository.NewUserRepository(db)
	_ = rp.InitTable()
	//4.创建服务实例
	userDataService := service.NewUserDataService(rp)
	//5.注册handler
	err = user.RegisterUserHandler(srv.Server(), &handler.User{UserDataService: userDataService})
	if err != nil {
		fmt.Println(err)
	}
	//6. Run service
	if err := srv.Run(); err != nil {
		fmt.Println()
	}
}
