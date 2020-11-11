package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"os"
	"github.com/jmoiron/sqlx"
	"google.golang.org/grpc"
)

type server struct{}

func (*server) CreateProductGroup(ctx context.Context, req *prdgpb.CreateProductReq) (*prdgpb.CreateProductRes, error) {
	db := connectDB()
	defer db.Close()
	prd := req.GetProduct()

	prdInsert := `INSERT INTO product_groups (id, name, created_at, updated_at) VALUES ($1, $2, $3, $4)`
	_, err := db.Exec(prdInsert, req.Product.GetId(), req.Product.GetName(), req.Product.GetCreatedAt(), req.Product.GetUpdatedAt())
	if err != nil {
		return nil, err
	}
	return &prdgpb.CreateProductRes{Product: prd}, nil
}

func (*server) ReadProductGroup(ctx context.Context, req *prdgpb.ReadProductReq) (*prdgpb.ReadProductRes, error) {
	db := connectDB()

	defer db.Close()
	pgrp := model.ProductGroup{}
	id := req.GetId()
	err := db.QueryRowx("SELECT * FROM product_groups where id = $1", id).StructScan(&pgrp)
	if err != nil {
		return nil, err
	}

	return &prdgpb.ReadProductRes{
		Product: &prdgpb.ProductGroup{
			Id:        pgrp.Id,
			Name:      pgrp.Name.String,
			CreatedAt: pgrp.CreatedAt,
			UpdatedAt: pgrp.UpdatedAt,
		},
	}, nil
}

func (*server) UpdateProductGroup(ctx context.Context, req *prdgpb.UpdateProductReq) (*prdgpb.UpdateProductRes, error) {
	db := connectDB()

	prod := req.GetProduct()
	id := req.GetProduct().GetId()
	prd := model.ProductGroup{}
	err := db.QueryRowx("SELECT * FROM product_groups where id = $1", id).StructScan(&prd)
	if err != nil {
		return nil, err
	}

	defer db.Close()

	update := `UPDATE product_groups SET name=$1, created_at=$2, updated_at=$3 WHERE id=$4;`
	_, err = db.Exec(update, req.Product.GetName(), req.Product.GetCreatedAt(), req.Product.GetUpdatedAt(), req.Product.GetId())
	if err != nil {
		panic(err)
	}

	return &prdgpb.UpdateProductRes{Product: prod}, nil

}

func (*server) DeleteProductGroup(ctx context.Context, req *prdgpb.DeleteProductReq) (*prdgpb.DeleteProductRes, error) {
	db := connectDB()
	defer db.Close()

	status := false

	id := req.GetId()
	res, err := db.Exec("DELETE FROM product_groups where id = $1", id)
	if err != nil {
		return &prdgpb.DeleteProductRes{Del: status}, err
	}

	count, err := res.RowsAffected()
	if err != nil {
		panic(err)
	}

	if count > 0 {
		status = true
	}

	fmt.Println(count)

	return &prdgpb.DeleteProductRes{Del: status}, nil
}

func (*server) ListProductGroup(ctx context.Context, req *prdgpb.ListProductReq) (*prdgpb.ListProductRes, error) {
	db := connectDB()
	defer db.Close()

	rows, err := db.Queryx("SELECT * FROM product_groups ORDER BY id ASC LIMIT $1 OFFSET $2;", req.Limit, req.Offset)
	if err != nil {
		return nil, err
	}

	list := []*prdgpb.ProductGroup{}
	for rows.Next() {
		prd := new(model.ProductGroup)
		err = rows.StructScan(&prd)
		if err != nil {
			return nil, err
		}
		list = append(list, dataProduct(prd))
	}

	return &prdgpb.ListProductRes{
		Products: list,
	}, nil
}

func dataProduct(data *model.ProductGroup) *prdgpb.ProductGroup {
	return &prdgpb.ProductGroup{
		Id:        data.Id,
		Name:      data.Name.String,
		CreatedAt: data.CreatedAt,
		UpdatedAt: data.UpdatedAt,
	}
}

func connectDB() *sqlx.DB {
	db, err := database.GetSqlxConnection()
	if err != nil {
		panic(err)
	}
	return db
}

func main() {
	addr := fmt.Sprintf(os.Getenv("VENDOR_ADDR"))
	lis, err := net.Listen("tcp", addr)
	log.Printf("Starting server on: %v", addr)

	if err != nil {
		log.Fatalf("Cannot connect to: %v", err)
	}

	srv := grpc.NewServer()
	prdgpb.RegisterProductServiceServer(srv, &server{})
	if err := srv.Serve(lis); err != nil {
		log.Fatalf("Cannot listen to: %v", err)
	}
}
