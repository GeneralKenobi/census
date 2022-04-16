package repository

import (
	"fmt"
	"github.com/GeneralKenobi/census/internal/db"
	"github.com/GeneralKenobi/census/pkg/util"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

// decodeSingleResult decodes the result to TBson and maps it to T using mappingFunc. It handles the error in result.
func decodeSingleResult[TBson, T any](queryName string, result *mongo.SingleResult, mappingFunc func(TBson) T) (T, error) {
	err := result.Err()
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return util.ZeroValue[T](), fmt.Errorf("%s: %w", queryName, db.ErrNoRows)
		}
		return util.ZeroValue[T](), fmt.Errorf("%s: error running query: %w", queryName, err)
	}

	var bsonEntity TBson
	err = result.Decode(&bsonEntity)
	if err != nil {
		return util.ZeroValue[T](), fmt.Errorf("%s: error processing query result: %w", queryName, err)
	}

	return mappingFunc(bsonEntity), nil
}

// processInsertOneResult is a helper for working with a result of an insert one operation. It returns a standardised error if queryErr is
// non-nil. It also tries to convert the inserted document ID to primitive.ObjectID and returns an error on failure.
func processInsertOneResult(queryName string, result *mongo.InsertOneResult, queryErr error) (primitive.ObjectID, error) {
	if queryErr != nil {
		return primitive.NilObjectID, fmt.Errorf("%s: error running query: %w", queryName, queryErr)
	}

	objectId, ok := result.InsertedID.(primitive.ObjectID)
	if !ok {
		return primitive.NilObjectID, fmt.Errorf("%s: inserted document ID (%v) is not an ObjectID", queryName, result.InsertedID)
	}

	return objectId, nil
}

// processUpdateOneResult is a helper for working with a result of an update operation that should've matched a single document. It returns
// a standardised error if queryErr is non-nil. It also returns an error if the number of matched documents is not 1.
func processUpdateOneResult(queryName string, result *mongo.UpdateResult, queryErr error) error {
	if queryErr != nil {
		return fmt.Errorf("%s: error running query: %w", queryName, queryErr)
	}
	if result.MatchedCount == 0 {
		return fmt.Errorf("%s: %w", queryName, db.ErrNoRows)
	}
	if result.MatchedCount > 1 {
		return fmt.Errorf("%s: %w", queryName, db.ErrTooManyRows)
	}
	return nil
}

// processDeleteOneResult is a helper for working with a result of a delete operation that should've matched a single document. It returns
// a standardised error if queryErr is non-nil. It also returns an error if the number of deleted documents is not 1.
func processDeleteOneResult(queryName string, result *mongo.DeleteResult, queryErr error) error {
	if queryErr != nil {
		return fmt.Errorf("%s: error running query: %w", queryName, queryErr)
	}
	if result.DeletedCount == 0 {
		return fmt.Errorf("%s: %w", queryName, db.ErrNoRows)
	}
	if result.DeletedCount > 1 {
		return fmt.Errorf("%s: %w", queryName, db.ErrTooManyRows)
	}
	return nil
}
