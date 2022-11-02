package main

import (
	"context"
	"sort"

	"github.com/examples/bookstore/rpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	emptypb "google.golang.org/protobuf/types/known/emptypb"
)

type shelf struct {
	id    int64
	theme string
	books map[int64]*book
}

type book struct {
	id     int64
	author string
	title  string
}

type bookstoreServer struct {
	storage map[int64]*shelf

	rpc.UnimplementedBookstoreServer
}

func NewBookstoreServer() *bookstoreServer {
	return &bookstoreServer{
		storage: make(map[int64]*shelf),
	}
}

func (bs *bookstoreServer) Reset(ctx context.Context, req *emptypb.Empty) (*emptypb.Empty, error) {
	bs.storage = make(map[int64]*shelf)
	return &emptypb.Empty{}, nil
}

// Returns a list of all shelves in the bookstore.
func (bs *bookstoreServer) ListShelves(ctx context.Context, req *emptypb.Empty) (*rpc.ListShelvesResponse, error) {
	shelves := make([]*rpc.Shelf, 0)
	for _, v := range bs.storage {
		shelves = append(shelves, &rpc.Shelf{Id: v.id, Theme: v.theme})
	}
	sort.Slice(shelves, func(i, j int) bool {
		return shelves[i].Id < shelves[j].Id
	})
	return &rpc.ListShelvesResponse{Shelves: shelves}, nil
}

// Creates a new shelf in the bookstore.
func (bs *bookstoreServer) CreateShelf(ctx context.Context, req *rpc.CreateShelfRequest) (*rpc.Shelf, error) {
	if bs.storage[req.Shelf.Id] != nil {
		return nil, status.Errorf(codes.AlreadyExists, "a shelf with id %d already exists", req.Shelf.Id)
	}
	bs.storage[req.Shelf.Id] = &shelf{
		id:    req.Shelf.Id,
		theme: req.Shelf.Theme,
		books: make(map[int64]*book)}
	return req.Shelf, nil
}

// Returns a specific bookstore shelf.
func (bs *bookstoreServer) GetShelf(ctx context.Context, req *rpc.GetShelfRequest) (*rpc.Shelf, error) {
	s := bs.storage[req.Shelf]
	if s == nil {
		return nil, status.Errorf(codes.NotFound, "a shelf with id %d does not exist", req.Shelf)
	}
	return &rpc.Shelf{Id: s.id, Theme: s.theme}, nil
}

// Deletes a shelf, including all books that are stored on the shelf.
func (bs *bookstoreServer) DeleteShelf(ctx context.Context, req *rpc.DeleteShelfRequest) (*emptypb.Empty, error) {
	if bs.storage[req.Shelf] == nil {
		return nil, status.Errorf(codes.NotFound, "a shelf with id %d does not exist", req.Shelf)
	}
	delete(bs.storage, req.Shelf)
	return &emptypb.Empty{}, nil
}

// Returns a list of books on a shelf.
func (bs *bookstoreServer) ListBooks(ctx context.Context, req *rpc.ListBooksRequest) (*rpc.ListBooksResponse, error) {
	s := bs.storage[req.Shelf]
	if s == nil {
		return nil, status.Errorf(codes.NotFound, "a shelf with id %d does not exist", req.Shelf)
	}
	books := make([]*rpc.Book, 0)
	for _, v := range s.books {
		books = append(books, &rpc.Book{Id: v.id, Author: v.author, Title: v.title})
	}
	sort.Slice(books, func(i, j int) bool {
		return books[i].Id < books[j].Id
	})
	return &rpc.ListBooksResponse{Books: books}, nil
}

// Creates a new book.
func (bs *bookstoreServer) CreateBook(ctx context.Context, req *rpc.CreateBookRequest) (*rpc.Book, error) {
	s := bs.storage[req.Shelf]
	if s == nil {
		return nil, status.Errorf(codes.NotFound, "a shelf with id %d does not exist", req.Shelf)
	}
	if s.books[req.Book.Id] != nil {
		return nil, status.Errorf(codes.AlreadyExists, "a book with id %d already exists", req.Book.Id)
	}
	s.books[req.Book.Id] = &book{
		id:     req.Book.Id,
		author: req.Book.Author,
		title:  req.Book.Title,
	}
	return req.Book, nil
}

// Returns a specific book.
func (bs *bookstoreServer) GetBook(ctx context.Context, req *rpc.GetBookRequest) (*rpc.Book, error) {
	s := bs.storage[req.Shelf]
	if s == nil {
		return nil, status.Errorf(codes.NotFound, "a shelf with id %d does not exist", req.Shelf)
	}
	b := s.books[req.Book]
	if b == nil {
		return nil, status.Errorf(codes.NotFound, "a book with id %d does not exist", req.Book)
	}
	return &rpc.Book{Id: b.id, Author: b.author, Title: b.title}, nil
}

// Deletes a book from a shelf.
func (bs *bookstoreServer) DeleteBook(ctx context.Context, req *rpc.DeleteBookRequest) (*emptypb.Empty, error) {
	s := bs.storage[req.Shelf]
	if s == nil {
		return nil, status.Errorf(codes.NotFound, "a shelf with id %d does not exist", req.Shelf)
	}
	if s.books[req.Book] == nil {
		return nil, status.Errorf(codes.NotFound, "a book with id %d does not exist", req.Book)
	}
	delete(s.books, req.Book)
	return &emptypb.Empty{}, nil
}
