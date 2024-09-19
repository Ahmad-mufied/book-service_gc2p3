package service

import (
	"book-service_gc2p3/entity"
	"book-service_gc2p3/pb"
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type BookService struct {
	db *mongo.Client
}

func NewBookService(db *mongo.Client) *BookService {
	return &BookService{db: db}
}

func (s *BookService) Create(ctx context.Context, req *pb.CreateBookRequest) (*pb.CreateBookResponse, error) {
	// Get data from request
	title := req.GetTitle()
	author := req.GetAuthor()
	publishData := req.GetPublishDate()
	status := req.GetStatus()

	// Insert data to database
	coll := s.db.Database("book_db").Collection("books")
	book := entity.Book{
		Title:       title,
		Author:      author,
		PublishData: publishData,
		Status:      status,
	}
	result, err := coll.InsertOne(ctx, book)
	if err != nil {
		return nil, err
	}

	book.Id = result.InsertedID.(primitive.ObjectID)

	// Return response
	res := &pb.CreateBookResponse{
		Success: true,
		Id:      book.Id.Hex(),
	}

	return res, nil
}

func (s *BookService) GetBookById(ctx context.Context, req *pb.GetBookByIdRequest) (*pb.GetBookByIdResponse, error) {
	// Get data from request
	bookId := req.GetId()

	// Find book by id
	bookIdObj, _ := primitive.ObjectIDFromHex(bookId)
	coll := s.db.Database("book_db").Collection("books")
	var book entity.Book

	err := coll.FindOne(ctx, bson.M{"_id": bookIdObj}).Decode(&book)
	//log.Printf("Book ID: %s\n", bookId)
	if err != nil {
		res := &pb.GetBookByIdResponse{
			Success: false,
		}
		return res, nil
	}

	// Return response
	res := &pb.GetBookByIdResponse{
		Title:       book.Title,
		Author:      book.Author,
		PublishDate: book.PublishData,
		Status:      book.Status,
		Success:     true,
	}

	return res, nil
}

func (s *BookService) GetBookByTitle(ctx context.Context, req *pb.GetBookByTitleRequest) (*pb.GetBookByTitleResponse, error) {
	// Get data from request
	title := req.GetTitle()

	// Find book by title
	coll := s.db.Database("book_db").Collection("books")
	var book entity.Book
	err := coll.FindOne(ctx, bson.M{"title": title}).Decode(&book)
	if err != nil {
		return nil, err
	}

	// Return response
	res := &pb.GetBookByTitleResponse{
		Id:          book.Id.Hex(),
		Author:      book.Author,
		PublishDate: book.PublishData,
		Status:      book.Status,
	}

	return res, nil
}

func (s *BookService) Update(ctx context.Context, req *pb.UpdateBookRequest) (*pb.UpdateBookResponse, error) {
	// Get data from request
	bookId := req.GetId()
	title := req.GetTitle()
	author := req.GetAuthor()
	publishData := req.GetPublishDate()
	status := req.GetStatus()

	// Check book by id

	// Update book by id
	coll := s.db.Database("book_db").Collection("books")
	bookIdObj, _ := primitive.ObjectIDFromHex(bookId)
	_, err := coll.UpdateOne(ctx, bson.M{"_id": bookIdObj}, bson.M{"$set": bson.M{
		"title":        title,
		"author":       author,
		"publish_data": publishData,
		"status":       status,
	}})
	if err != nil {
		return nil, err
	}

	// Return response
	res := &pb.UpdateBookResponse{
		Success: true,
	}

	return res, nil
}

func (s *BookService) Delete(ctx context.Context, req *pb.DeleteBookRequest) (*pb.DeleteBookResponse, error) {
	// Get data from request
	bookId := req.GetId()

	// Delete book by id
	bookIdObj, _ := primitive.ObjectIDFromHex(bookId)
	coll := s.db.Database("book_db").Collection("books")
	_, err := coll.DeleteOne(ctx, bson.M{"_id": bookIdObj})
	if err != nil {
		return nil, err
	}

	// Return response
	res := &pb.DeleteBookResponse{
		Success: true,
	}

	return res, nil
}

func (s *BookService) Check(ctx context.Context, req *pb.CheckBookRequest) (*pb.CheckBookResponse, error) {
	// Get data from request
	bookID := req.GetId()
	bookIdObj, _ := primitive.ObjectIDFromHex(bookID)

	// Check book by id
	coll := s.db.Database("book_db").Collection("books")

	// check if book exist by id
	count, err := coll.CountDocuments(ctx, bson.M{"_id": bookIdObj})
	if err != nil {
		return nil, err
	}

	if count == 0 {
		res := &pb.CheckBookResponse{
			Exist: false,
		}
		return res, nil
	}

	// Return response
	res := &pb.CheckBookResponse{
		Exist: true,
	}

	return res, nil
}
