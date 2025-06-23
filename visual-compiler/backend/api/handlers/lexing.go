package handlers

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	"github.com/COS301-SE-2025/Visual-Compiler/backend/core/db"
	"github.com/COS301-SE-2025/Visual-Compiler/backend/core/services"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

// Specifies the JSON body request.
type SourceCodeRequest struct {
	// Represents the source code the user enters
	Code string `json:"source_code" binding:"required"`
	// Represents the pairs of Type and Regex
	Pairs []services.TypeRegex `json:"pairs" binding:"required"`
	// Represents the User's ID from frontend
	UsersID bson.ObjectID `json:"users_id" binding:"required"`
}

// Specifies the JSON body request for the Users ID.
type IDRequest struct {
	// Represents the User's ID from frontend
	UsersID bson.ObjectID `json:"users_id" binding:"required"`
}

// Locally store a user's source code and regex expressions.
// Gets the source code from a JSON request.
// Formats the response as a JSON Body
//
// Returns:
//   - A JSON response body.
//   - A 200 OK response if successful
//   - A 500 Internal Server Error if any errors are caught for parsing errors
func StoreSourceCode(c *gin.Context) {
	var req SourceCodeRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Input is invalid", "details": err.Error()})
		return
	}

	jsonAsBytes, err := json.Marshal(req.Pairs)
	if err != nil {
		panic(err)
	}

	pairs := jsonAsBytes

	// services.SourceCode(req.Code)
	rules, err := services.ReadRegexRules(pairs)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Regex rule creation failed"})
	}

	mongoCli := db.ConnectClient()
	collection := mongoCli.Database("visual-compiler").Collection("lexing")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	filters := bson.M{"users_id": req.UsersID}
	var userexisting bson.M

	err = collection.FindOne(ctx, filters).Decode(&userexisting)

	if err == mongo.ErrNoDocuments {
		_, err = collection.InsertOne(ctx, bson.M{
			"code":     req.Code,
			"rules":    rules,
			"users_id": req.UsersID,
		})
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Database Insertion error"})
			return
		}
	} else if err == nil {
		updateexisting := bson.D{
			bson.E{Key: "$unset", Value: bson.M{
				"tokens":              "",
				"tokens_unidentified": "",
			}},
			bson.E{Key: "$set", Value: bson.M{
				"code":  req.Code,
				"rules": rules,
			}},
		}
		_, err = collection.UpdateOne(ctx, filters, updateexisting)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Database Update error"})
			return
		}
	} else {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database lookup error"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Code is ready for lexing"})
}

// Lexes the user's source code that is locally stored.
// Stores the tokens, unidentified input, user's source code and the user's id in the database.
// Gets the source code from the function GetSourceCode.
// Gets the users id from a global variable `UsersID`.
// Formats the response as a JSON Body
//
// Returns:
//   - A JSON response body.
//   - A 200 OK response if successful
//   - A 500 Internal Server Error if any errors are caught for parsing and lexing errors
func Lexing(c *gin.Context) {
	var req IDRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Input is invalid", "details": err.Error()})
		return
	}

	mongoCli := db.ConnectClient()
	collection := mongoCli.Database("visual-compiler").Collection("lexing")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var res struct {
		Code  string               `bson:"code"`
		Rules []services.TypeRegex `bson:"rules"`
	}

	err := collection.FindOne(ctx, bson.M{"users_id": req.UsersID}).Decode(&res)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Source code not found"})
		return
	}

	tokens, unidentified, errorcaught := services.CreateTokens(res.Code, res.Rules)
	if errorcaught != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Tokenization failed", "details": errorcaught.Error()})
		return
	}

	filters := bson.M{"users_id": req.UsersID}
	updateuserslexing := bson.M{"$set": bson.M{
		"tokens":              tokens,
		"tokens_unidentified": unidentified,
	}}

	_, err = collection.UpdateOne(ctx, filters, updateuserslexing)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to insert tokens"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"users_id":            req.UsersID,
		"message":             "Successfully tokenised your code",
		"tokens":              tokens,
		"tokens_unidentified": unidentified,
	})
}

type readDFARequest struct {
	States      []string                  `json:"states"`
	Transitions []services.Transition     `json:"transitions"`
	Start       string                    `json:"start_state"`
	Accepting   []services.AcceptingState `json:"accepting_states"`
	// Represents the User's ID from frontend
	UsersID bson.ObjectID `json:"users_id" binding:"required"`
}

func ReadDFAFromUser(c *gin.Context) {
	var req readDFARequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Input is invalid", "details": err.Error()})
		return
	}

	jsonAsBytes, err := json.Marshal(req)
	if err != nil {
		panic(err)
	}

	dfa := services.Automata{}

	err = json.Unmarshal(jsonAsBytes, &dfa)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error extracting details for DFA creation"})
		return
	}

	mongoCli := db.ConnectClient()
	collection := mongoCli.Database("visual-compiler").Collection("lexing")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	filters := bson.M{"users_id": req.UsersID}

	var userexisting bson.M
	err = collection.FindOne(ctx, filters).Decode(&userexisting)

	if err == mongo.ErrNoDocuments {
		_, err = collection.InsertOne(ctx, bson.M{
			"dfa":      dfa,
			"users_id": req.UsersID,
		})
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Database Insertion error"})
			return
		}
	} else if err == nil {
		updateexisting := bson.D{
			bson.E{Key: "$unset", Value: bson.M{
				"tokens":              "",
				"tokens_unidentified": "",
				"rules":               "",
			}},
			bson.E{Key: "$set", Value: bson.M{
				"dfa": dfa,
			}},
		}
		_, err = collection.UpdateOne(ctx, filters, updateexisting)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Database Update error"})
			return
		}
		_, err = collection.UpdateOne(ctx, filters, updateexisting)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Database Update error"})
			return
		}
	} else {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database lookup error"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "DFA successfuly created. Ready to create tokens"})
}

func TokensFromDFA(c *gin.Context) {
	var req IDRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Input is invalid", "details": err.Error()})
		return
	}

	mongoCli := db.ConnectClient()
	collection := mongoCli.Database("visual-compiler").Collection("lexing")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var res struct {
		Code string              `bson:"code"`
		DFA  []services.Automata `bson:"dfa"`
	}

	err := collection.FindOne(ctx, bson.M{"users_id": req.UsersID}).Decode(&res)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Source code not found"})
		return
	}

	// tokens, unidentified, errorcaught := services.CreateTokensFromDFA(res.Code, res.DFA)
	// if errorcaught != nil {
	// 	c.JSON(http.StatusInternalServerError, gin.H{"error": "Tokenization failed", "details": errorcaught.Error()})
	// 	return
	// }

	// filters := bson.M{"users_id": req.UsersID}
	// updateuserslexing := bson.M{"$set": bson.M{
	// 	"tokens":              tokens,
	// 	"tokens_unidentified": unidentified,
	// }}

	// _, err = collection.UpdateOne(ctx, filters, updateuserslexing)
	// if err != nil {
	// 	c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update tokens"})
	// 	return
	// }

	// c.JSON(http.StatusOK, gin.H{
	// 	"users_id":            req.UsersID,
	// 	"message":             "Successfully tokenised your code",
	// 	"tokens":              tokens,
	// 	"tokens_unidentified": unidentified,
	// })
}

func ConvertDFAToRG(c *gin.Context) {
	var req IDRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Input is invalid", "details": err.Error()})
		return
	}

	mongoCli := db.ConnectClient()
	collection := mongoCli.Database("visual-compiler").Collection("lexing")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var res struct {
		DFA []services.Automata `bson:"dfa"`
	}

	err := collection.FindOne(ctx, bson.M{"users_id": req.UsersID}).Decode(&res)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "DFA not found. Please create one"})
		return
	}

	// rules, errorCaught := services.ConvertDFAToRegex(res.DFA)
	// if errorcaught != nil {
	// 	c.JSON(http.StatusInternalServerError, gin.H{"error": "Tokenization failed", "details": errorcaught.Error()})
	// 	return
	// }

	// filters := bson.M{"users_id": req.UsersID}
	// updateuserslexing := bson.M{"$set": bson.M{
	// 	"rules" : rules,
	// }}

	// _, err = collection.UpdateOne(ctx, filters, updateuserslexing)
	// if err != nil {
	// 	c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to insert rules"})
	// 	return
	// }

	// c.JSON(http.StatusOK, gin.H{
	// 	"users_id":            req.UsersID,
	// 	"message":             "Successfully converted DFA to Regex",
	// 	"rules":              rules,
	// })
}
