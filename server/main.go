package main

import (
	"context"
	"database/sql"
	"log"
	"net"
	"net/http"

	config "github.com/abdullahPrasetio/backend-school-grpc/common/config"
	"github.com/abdullahPrasetio/backend-school-grpc/common/models/proto"
	pb "github.com/abdullahPrasetio/backend-school-grpc/common/models/proto"
	"github.com/abdullahPrasetio/backend-school-grpc/connections"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"

	"github.com/golang/protobuf/ptypes/empty"
	"golang.org/x/crypto/bcrypt"
	"google.golang.org/grpc"
)



type UsersServer struct {
	db *sql.DB
	pb.UnimplementedUsersServer
}

func NewServer(db *sql.DB) (*UsersServer) {
	return &UsersServer{db:db} 
}

func(s *UsersServer)Register(ctx context.Context,param *proto.UserRegister)(*proto.UserWithoutPassword,error){
	password :=[]byte(param.Password)
	hashedPassword,err :=bcrypt.GenerateFromPassword(password,bcrypt.DefaultCost)
	if err != nil{
		panic (err.Error())
	}
	sql :="INSERT INTO users (first_name,last_name,email,password,phone,role) VALUES (?,?,?,?,?,?)"
	res,err:=s.db.ExecContext(ctx,sql,param.FirstName,param.LastName,param.Email,hashedPassword,param.Phone,"0")
	
	if err!=nil {
		log.Println("Error inserting user",err.Error())
	}
	id, err := res.LastInsertId()
    if err != nil {
        panic (err.Error())
    }
    user:=proto.UserWithoutPassword{
		Id:id,
		FirstName:param.FirstName,
		LastName:param.LastName,
		Email:param.Email,
		Phone:param.Phone,
	}
	log.Println("User",&user)
	return &user,nil
}

func (s *UsersServer)List(ctx context.Context,void *empty.Empty) (*proto.UserList, error) {
	sql:="SELECT id,first_name,last_name,email,phone,role FROM users"
	userList := proto.UserList{}
	rows,err:=s.db.QueryContext(ctx,sql)
	if err != nil {
		return &userList,err
	}
	defer rows.Close()
	for rows.Next() {
		user:=proto.UserWithoutPassword{}
		rows.Scan(&user.Id,&user.FirstName,&user.LastName,&user.Email,&user.Phone,&user.Role)
		log.Println(&user)
		userList.List = append(userList.List, &user)
	}
	
	return &userList, nil
}
func main() {
	srv := grpc.NewServer()
	// Conect DB
	db,err:=connections.NewConnection()
	if err!=nil {
		panic(err)
	}

    userSrv :=NewServer(db)
	
	go func ()  {
		// mux
		mux:=runtime.NewServeMux()
		// register
		pb.RegisterUsersHandlerServer(context.Background(),mux,userSrv)
		log.Fatalln(http.ListenAndServe("localhost:8082",mux))
	}()

    proto.RegisterUsersServer(srv, userSrv)

    log.Println("Starting RPC server at", config.SERVICE_USER_PORT)

  	l, err := net.Listen("tcp", config.SERVICE_USER_PORT)
	if err != nil {
	log.Fatalf("could not listen to %s: %v", config.SERVICE_USER_PORT, err)
	}

	log.Fatal(srv.Serve(l))

	
}