package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"net"

	_ "github.com/lib/pq"
	pb "github.com/tcc-grpc/v2/filmepb"
	"google.golang.org/grpc"
)

const (
	porta = ":50052"
)

func NewFilmeServiceServer() *FilmeServiceServer {
	return &FilmeServiceServer{}
}

type FilmeServiceServer struct {
	conn *sql.DB
	pb.UnimplementedMoviesServiceServer
}

func (server *FilmeServiceServer) Run() error {
	lis, err := net.Listen("tcp", porta)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	s := grpc.NewServer()
	pb.RegisterMoviesServiceServer(s, server)
	log.Printf("server listening at %v", lis.Addr())
	return s.Serve(lis)
}

func OpenConnection() *sql.DB {
	psqlInfo := fmt.Sprintf("host=%s user=%s "+
		"password=%s dbname=%s sslmode=require",
		host, user, password, dbname)

	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		panic(err)
	}

	err = db.Ping()
	if err != nil {
		panic(err)
	}

	return db
}

func (server *FilmeServiceServer) ListarFilmesGrpc(ctx context.Context, in *pb.ReadFilmesRequest) (*pb.ReadFilmesResponse, error) {

	var filmes_list *pb.ReadFilmesResponse = &pb.ReadFilmesResponse{}
	db := OpenConnection()

	rows, err := db.Query(`select distinct m.movieid, m.title, m.genres, r.userid, r.rating, r.timestamp, t.tag from movies m 
	inner join ratings r on r.movieid = m.movieid 
	inner join tags t on t.movieid = m.movieid
	order by movieid asc
	limit 1500`)
	if err != nil {
		return nil, err
	}
	defer db.Close()
	for rows.Next() {
		filme := pb.Filme{}
		err = rows.Scan(&filme.Movieid, &filme.Title, &filme.Genres, &filme.Userid, &filme.Rating, &filme.Timestamp, &filme.Tag)
		if err != nil {
			return nil, err
		}
		filmes_list.Filmes = append(filmes_list.Filmes, &filme)
	}
	defer rows.Close()
	// defer db.Close()

	return filmes_list, nil

}

const (
	host     = "bancotcc.postgres.database.azure.com"
	user     = "jefferson"
	password = "Postgres@"
	dbname   = "postgres"
)

func main() {

	psqlInfo := fmt.Sprintf("host=%s user=%s "+
		"password=%s dbname=%s sslmode=require",
		host, user, password, dbname)

	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		panic(err)
	}
	err = db.Ping()
	if err != nil {
		panic(err)
	}

	var filme_server *FilmeServiceServer = NewFilmeServiceServer()
	filme_server.conn = db

	if err := filme_server.Run(); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
