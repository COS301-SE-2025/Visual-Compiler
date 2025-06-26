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

type SourceCodeOnlyRequest struct {
	// Represents the User's ID from frontend
	UsersID bson.ObjectID `json:"users_id" binding:"required"`
	Code    string        `json:"source_code" binding:"required"`
}

// Specifies the JSON body request.
type RulesRequest struct {
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

func StoreSourceCode(c *gin.Context) {
	var req SourceCodeOnlyRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Input is invalid", "details": err.Error()})
		return
	}

	mongo_cli := db.ConnectClient()
	collection := mongo_cli.Database("visual-compiler").Collection("lexing")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	filters := bson.M{"users_id": req.UsersID}
	var userexisting bson.M

	err := collection.FindOne(ctx, filters).Decode(&userexisting)

	if err == mongo.ErrNoDocuments {
		_, err = collection.InsertOne(ctx, bson.M{
			"code":     req.Code,
			"users_id": req.UsersID,
		})
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Database Insertion error"})
			return
		}
	} else if err == nil {
		update_existing := bson.D{
			bson.E{Key: "$unset", Value: bson.M{
				"tokens":              "",
				"tokens_unidentified": "",
				"nfa":                 "",
				"dfa":                 "",
				"rules":               "",
			}},
			bson.E{Key: "$set", Value: bson.M{
				"code": req.Code,
			}},
		}
		_, err = collection.UpdateOne(ctx, filters, update_existing)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Database Update error"})
			return
		}
	} else {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database lookup error"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Code is ready for further processing"})
}

// Locally store a user's source code and regex expressions.
// Gets the source code from a JSON request.
// Formats the response as a JSON Body
//
// Returns:
//   - A JSON response body.
//   - A 200 OK response if successful
//   - A 500 Internal Server Error if any errors are caught for parsing errors
func CreateRulesFromCode(c *gin.Context) {
	var req RulesRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Input is invalid", "details": err.Error()})
		return
	}

	json_as_bytes, err := json.Marshal(req.Pairs)
	if err != nil {
		panic(err)
	}

	pairs := json_as_bytes

	rules, err := services.ReadRegexRules(pairs)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Regex rule creation failed"})
	}

	mongo_cli := db.ConnectClient()
	collection := mongo_cli.Database("visual-compiler").Collection("lexing")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var res struct {
		Code string `bson:"code"`
	}

	err = collection.FindOne(ctx, bson.M{"users_id": req.UsersID}).Decode(&res)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Source code not found. Please enter a source code"})
		return
	}

	filters := bson.M{"users_id": req.UsersID}
	update_users_rules := bson.M{
		"$set": bson.M{
			"rules": rules,
		},
	}

	_, err = collection.UpdateOne(ctx, filters, update_users_rules)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to insert rules"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Rules successfully created."})
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

	mongo_cli := db.ConnectClient()
	collection := mongo_cli.Database("visual-compiler").Collection("lexing")

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

	tokens, unidentified, error_caught := services.CreateTokens(res.Code, res.Rules)
	if error_caught != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Tokenization failed", "details": error_caught.Error()})
		return
	}

	filters := bson.M{"users_id": req.UsersID}
	update_users_lexing := bson.M{"$set": bson.M{
		"tokens":              tokens,
		"tokens_unidentified": unidentified,
	}}

	_, err = collection.UpdateOne(ctx, filters, update_users_lexing)
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

	json_as_bytes, err := json.Marshal(req)
	if err != nil {
		panic(err)
	}

	dfa := services.Automata{}

	err = json.Unmarshal(json_as_bytes, &dfa)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error extracting details for DFA creation"})
		return
	}

	mongo_cli := db.ConnectClient()
	collection := mongo_cli.Database("visual-compiler").Collection("lexing")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	filters := bson.M{"users_id": req.UsersID}

	var user_existing bson.M
	err = collection.FindOne(ctx, filters).Decode(&user_existing)

	if err == mongo.ErrNoDocuments {
		_, err = collection.InsertOne(ctx, bson.M{
			"dfa": dfa,
		})
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Database Insertion error"})
			return
		}
	} else if err == nil {
		update_existing := bson.D{
			bson.E{Key: "$unset", Value: bson.M{
				"tokens":              "",
				"tokens_unidentified": "",
				"rules":               "",
			}},
			bson.E{Key: "$set", Value: bson.M{
				"dfa": dfa,
			}},
		}
		_, err = collection.UpdateOne(ctx, filters, update_existing)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Database Update error"})
			return
		}
		_, err = collection.UpdateOne(ctx, filters, update_existing)
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

	mongo_cli := db.ConnectClient()
	collection := mongo_cli.Database("visual-compiler").Collection("lexing")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var res struct {
		Code string            `bson:"code"`
		DFA  services.Automata `bson:"dfa"`
	}

	err := collection.FindOne(ctx, bson.M{"users_id": req.UsersID}).Decode(&res)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Source code not found"})
		return
	}

	tokens, unidentified, error_caught := services.CreateTokensFromDFA(res.Code, res.DFA)
	if error_caught != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Tokenization from DFA failed", "details": error_caught.Error()})
		return
	}

	filters := bson.M{"users_id": req.UsersID}
	update_users_lexing := bson.D{
		bson.E{
			Key: "$set", Value: bson.M{
				"tokens":              tokens,
				"tokens_unidentified": unidentified,
			},
		}}

	_, err = collection.UpdateOne(ctx, filters, update_users_lexing)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update tokens"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":             "Successfully tokenised your code",
		"tokens":              tokens,
		"tokens_unidentified": unidentified,
	})
}

func ConvertDFAToRG(c *gin.Context) {
	var req IDRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Input is invalid", "details": err.Error()})
		return
	}

	mongo_cli := db.ConnectClient()
	collection := mongo_cli.Database("visual-compiler").Collection("lexing")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var res struct {
		DFA services.Automata `bson:"dfa"`
	}

	err := collection.FindOne(ctx, bson.M{"users_id": req.UsersID}).Decode(&res)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "DFA not found. Please create one"})
		return
	}

	rules, error_caught := services.ConvertDFAToRegex(res.DFA)
	if error_caught != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Conversion from DFA to Regex failed", "details": error_caught.Error()})
		return
	}

	filters := bson.M{"users_id": req.UsersID}
	update_users_lexing := bson.M{"$set": bson.M{
		"rules": rules,
	}}

	_, err = collection.UpdateOne(ctx, filters, update_users_lexing)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to insert rules"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Successfully converted DFA to Regex",
		"rules":   rules,
	})
}

func ConvertRGToNFA(c *gin.Context) {
	var req IDRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Input is invalid", "details": err.Error()})
		return
	}

	mongo_cli := db.ConnectClient()
	collection := mongo_cli.Database("visual-compiler").Collection("lexing")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var res struct {
		Rules []services.TypeRegex `bson:"rules"`
	}

	err := collection.FindOne(ctx, bson.M{"users_id": req.UsersID}).Decode(&res)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Regex rules not found. Please create one"})
		return
	}

	regexes_from_rules := make(map[string]string)
	for _, rule := range res.Rules {
		regexes_from_rules[rule.Type] = rule.Regex
	}

	nfa, error_caught := services.ConvertRegexToNFA(regexes_from_rules)
	if error_caught != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Conversion from Regex to NFA failed", "details": error_caught.Error()})
		return
	}

	filters := bson.M{"users_id": req.UsersID}
	update_users_lexing := bson.M{"$set": bson.M{
		"nfa": nfa,
	}}

	_, err = collection.UpdateOne(ctx, filters, update_users_lexing)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to insert NFA"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Successfully converted Regex to NFA",
		"nfa":     nfa,
	})
}

func ConvertRGToDFA(c *gin.Context) {
	var req IDRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Input is invalid", "details": err.Error()})
		return
	}

	mongo_cli := db.ConnectClient()
	collection := mongo_cli.Database("visual-compiler").Collection("lexing")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var res struct {
		Rules []services.TypeRegex `bson:"rules"`
	}

	err := collection.FindOne(ctx, bson.M{"users_id": req.UsersID}).Decode(&res)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Regex rules not found. Please create one"})
		return
	}

	regexes_from_rules := make(map[string]string)
	for _, rule := range res.Rules {
		regexes_from_rules[rule.Type] = rule.Regex
	}

	dfa, error_caught := services.ConvertRegexToDFA(regexes_from_rules)
	if error_caught != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Conversion from Regex to DFA failed", "details": error_caught.Error()})
		return
	}

	filters := bson.M{"users_id": req.UsersID}
	update_users_lexing := bson.M{"$set": bson.M{
		"dfa": dfa,
	}}

	_, err = collection.UpdateOne(ctx, filters, update_users_lexing)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to insert DFA"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Successfully converted Regex to DFA",
		"dfa":     dfa,
	})
}

func ConvertNFAToDFA(c *gin.Context) {
	var req IDRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Input is invalid", "details": err.Error()})
		return
	}

	mongo_cli := db.ConnectClient()
	collection := mongo_cli.Database("visual-compiler").Collection("lexing")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var res struct {
		NFA services.Automata `bson:"nfa"`
	}

	err := collection.FindOne(ctx, bson.M{"users_id": req.UsersID}).Decode(&res)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "NFA not found. Please create one"})
		return
	}

	dfa, error_caught := services.ConvertNFAToDFA(res.NFA)
	if error_caught != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Conversion from NFA to DFA failed", "details": error_caught.Error()})
		return
	}

	filters := bson.M{"users_id": req.UsersID}
	update_users_lexing := bson.M{"$set": bson.M{
		"dfa": dfa,
	}}

	_, err = collection.UpdateOne(ctx, filters, update_users_lexing)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to insert DFA"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Successfully converted NFA to DFA",
		"dfa":     dfa,
	})
}
