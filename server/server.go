package main

import (
	"context"
	"log"
	"net"

	pb "github.com/tcc-grpc/v2/filmepb"

	"github.com/jackc/pgx/v4"
	"google.golang.org/grpc"
)

const (
	port = ":50052"
)

func NewFilmeServiceServer() *FilmeServiceServer {
	return &FilmeServiceServer{}
}

type FilmeServiceServer struct {
	conn *pgx.Conn
	pb.UnimplementedMoviesServiceServer
}

func (server *FilmeServiceServer) Run() error {
	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	s := grpc.NewServer()
	pb.RegisterMoviesServiceServer(s, server)
	log.Printf("server listening at %v", lis.Addr())
	return s.Serve(lis)
}

func (server *FilmeServiceServer) ListarFilmesGrpc(ctx context.Context, in *pb.ReadFilmesRequest) (*pb.ReadFilmesResponse, error) {

	var filmes_list *pb.ReadFilmesResponse = &pb.ReadFilmesResponse{}

	rows, err := server.conn.Query(context.Background(), `select distinct m.movieid, m.title, m.genres, r.userid, r.rating, r.timestamp, t.tag from movies m 
	inner join ratings r on r.movieid = m.movieid 
	inner join tags t on t.movieid = m.movieid`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		filme := pb.Filme{}
		err = rows.Scan(&filme.Movieid, &filme.Title, &filme.Genres, &filme.Userid, &filme.Rating, &filme.Timestamp, &filme.Tag)
		if err != nil {
			return nil, err
		}
		filmes_list.Filmes = append(filmes_list.Filmes, &filme)
	}

	return filmes_list, nil

}
func main() {

	database_url := "postgres://postgres:postgrespw@localhost:55001"
	conn, err := pgx.Connect(context.Background(), database_url)
	if err != nil {
		log.Fatalf("Unable to establish connection: %v", err)
	}
	defer conn.Close(context.Background())
	var filme_server *FilmeServiceServer = NewFilmeServiceServer()
	filme_server.conn = conn

	if err := filme_server.Run(); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
