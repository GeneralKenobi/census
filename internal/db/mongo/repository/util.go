package repository

import (
	"fmt"
	"github.com/GeneralKenobi/census/internal/db"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// idAsObjectIdOptional is a more lenient idAsObjectId - it returns primitive.NilObjectID if id is an empty string of if it's not a valid
// primitive.ObjectID. This function can be used when id is known to be valid or when treating empty/invalid id as primitive.NilObjectID is
// acceptable.
func idAsObjectIdOptional(id string) primitive.ObjectID {
	if id == "" {
		return primitive.NilObjectID
	}

	objectId, err := idAsObjectId(id)
	if err != nil {
		return primitive.NilObjectID
	}
	return objectId
}

// idAsObjectId tries to convert the given ID to primitive.ObjectID. It returns db.ErrNoRows on failure because non-ObjectID IDs can't match
// any documents because every document has an ID of type ObjectID.
func idAsObjectId(id string) (primitive.ObjectID, error) {
	objectId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return primitive.NilObjectID, fmt.Errorf("non-ObjectID IDs like %s are not supported: %w", id, db.ErrNoRows)
	}
	return objectId, nil
}

// filterById returns a bson filter document that matches only documents with the given id.
func filterById(id primitive.ObjectID) bson.M {
	return bson.M{"_id": id}
}

// updateBySet returns a bson document that applies update operation based on updateSpec using the $set operator.
func updateBySet(updateSpec any) bson.D {
	return bson.D{{"$set", updateSpec}}
}
