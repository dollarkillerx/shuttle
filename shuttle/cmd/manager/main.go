package main

import (
	"flag"
	"fmt"
	"log"
	"net"
	"net/http"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
	"google.dev/google/common/pkg/logger"
	"google.dev/google/shuttle/core/app/manager/conf"
	"google.dev/google/shuttle/core/app/manager/generated"
	"google.dev/google/shuttle/core/app/manager/middlewares"
	"google.dev/google/shuttle/core/app/manager/resolvers"
	"google.dev/google/shuttle/core/app/manager/rpc"
	"google.dev/google/shuttle/core/app/manager/storage/simple"
	"google.dev/google/shuttle/core/app/manager/utils"
	"google.dev/google/shuttle/proto/manager"
	"google.dev/google/socks5_discovery/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

var configFilename string
var configDirs string

func init() {
	const (
		defaultConfigFilename = "config"
		configUsage           = "Name of the config file, without extension"
		defaultConfigDirs     = "configs"
		configDirUsage        = "Directories to search for config file, separated by ','"
	)
	flag.StringVar(&configFilename, "c", defaultConfigFilename, configUsage)
	flag.StringVar(&configFilename, "config", defaultConfigFilename, configUsage)
	flag.StringVar(&configDirs, "cPath", defaultConfigDirs, configDirUsage)
}

// 用户管理， 流量控制
func main() {
	managerServiceInit()

	router := chi.NewRouter()
	router.Use(middlewares.Cors())
	if !conf.CONFIG.EnablePlayground {
		router.Use(middlewares.Safety())
	}

	router.Use(middlewares.Context())

	router.Post("/upload", utils.UploadFile)
	router.Handle("/static/*", http.StripPrefix("/static", http.FileServer(http.Dir("static"))))

	router.Get("/healthcheck", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("ack"))
	})

	if conf.CONFIG.EnablePlayground {
		router.Handle("/playground", playground.Handler("GraphQL playground", "/graphql"))

		log.Println("playground is enabled!")
		log.Printf("connect to http://%s/playground for GraphQL playground", conf.CONFIG.GraphQLListenAddress)
	}

	//将所有的rpcClient放入 gqlgen的Resolver
	newSimple, err := simple.NewSimple(conf.CONFIG.PostgresConfiguration)
	if err != nil {
		log.Fatalln(err)
	}

	// Socks5DiscoveryAddress
	conn, err := grpc.Dial(conf.CONFIG.Socks5DiscoveryAddress, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		panic(err)
	}
	s5Discovery := proto.NewSocks5DiscoveryClient(conn)

	cf := generated.Config{
		Resolvers: resolvers.NewResolver(newSimple, s5Discovery),
	}

	cf.Directives.HasLogined = middlewares.HasLoginFunc

	graphQLServer := handler.NewDefaultServer(generated.NewExecutableSchema(cf))

	graphQLServer.SetRecoverFunc(middlewares.RecoverFunc)
	graphQLServer.SetErrorPresenter(middlewares.MiddleError)

	router.Handle("/graphql", graphQLServer)

	log.Printf("connect to http://%s/graphql", conf.CONFIG.GraphQLListenAddress)

	server := rpc.NewGRPCServer(newSimple)
	listen, err := net.Listen("tcp", conf.CONFIG.RPCListenAddress)
	if err != nil {
		panic(err)
	}
	grpc := grpc.NewServer()
	manager.RegisterGuardLinkManagerServer(grpc, server)
	go func() {
		log.Println("GRPC: ", conf.CONFIG.RPCListenAddress)
		grpc.Serve(listen)
	}()

	if err := http.ListenAndServe(conf.CONFIG.GraphQLListenAddress, router); err != nil {
		log.Fatalln(err)
	}
}

func managerServiceInit() {
	utils.InitJWT()
	flag.Parse()

	// Setting up configurations
	err := conf.InitConfiguration(configFilename, configDirs)
	if err != nil {
		panic(fmt.Errorf("Error parsing config, %s", err))
	}

	// init logger
	logger.InitLogger(conf.CONFIG.LoggerConfig)
}
