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

// Specifies the JSON body request for storing source code
type SourceCodeOnlyRequest struct {
	// Represents the User's ID from frontend
	UsersID bson.ObjectID `json:"users_id" binding:"required"`
	// Represents the User's source code
	Code string `json:"source_code" binding:"required"`
	// User's project name
	Project_Name string `json:"project_name" binding:"required"`
}

// Specifies the JSON body request.
type RulesRequest struct {
	// Represents the pairs of Type and Regex
	Pairs []services.TypeRegex `json:"pairs" binding:"required"`
	// Represents the User's ID from frontend
	UsersID bson.ObjectID `json:"users_id" binding:"required"`
	// User's project name
	Project_Name string `json:"project_name" binding:"required"`
}

// Specifies the JSON body request for the Users ID.
type IDRequest struct {
	// Represents the User's ID from frontend
	UsersID bson.ObjectID `json:"users_id" binding:"required"`
	// User's project name
	Project_Name string `json:"project_name" binding:"required"`
}

// @Summary Store the user's source code
// @Description Takes the user's source code and stores it with the user's ID. If a source code already exists, it overwrites that code and removes any other created fields (tokens, dfa, nfa, rules)
// @Tags Lexing
// @Accept json
// @Produce json
// @Param request body SourceCodeOnlyRequest true "Read Source Code From User"
// @Success 200 {object} map[string]string "Source code successfully stored"
// @Failure 400 {object} map[string]string "Invalid input"
// @Failure 500 {object} map[string]string "Internal server error"
// @Router /lexing/code [post]
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

	filters := bson.M{
		"users_id":     req.UsersID,
		"project_name": req.Project_Name,
	}
	var userexisting bson.M

	err := collection.FindOne(ctx, filters).Decode(&userexisting)

	if err == mongo.ErrNoDocuments {
		_, err = collection.InsertOne(ctx, bson.M{
			"code":         req.Code,
			"users_id":     req.UsersID,
			"project_name": req.Project_Name,
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

// @Summary Create rules from stored source code
// @Description Searches the database for the user's source code. If found, the pairs defined by the user are used to create the rules. The rules are either created or ,if already existing, updated. If the source code is not found, returns an error
// @Tags Lexing
// @Accept json
// @Produce json
// @Param request body RulesRequest true "Read Pairs and From User to Create Rules"
// @Success 200 {object} map[string]string "Rules successfully created and stored"
// @Failure 400 {object} map[string]string "Invalid input/Regex rules creation failed"
// @Failure 404 {object} map[string]string "Source code not found"
// @Failure 500 {object} map[string]string "Internal server error"
// @Router /lexing/rules [post]
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

	filters := bson.M{
		"users_id":     req.UsersID,
		"project_name": req.Project_Name,
	}
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

// @Summary Lexes the user's stored rules
// @Description Searches the database for the user's rules. If found, the source code and rules are used in the lexer to create the tokens and/or unidentified tokens. The tokens are either created or ,if already existing, updated. If the source code or rules are not found, returns an error
// @Tags Lexing
// @Accept json
// @Produce json
// @Param request body IDRequest true "Create Tokens from Stored Code and Rules"
// @Success 200 {object} map[string]string "Tokens successfully created and stored"
// @Failure 400 {object} map[string]string "Invalid input/Lexing failed"
// @Failure 404 {object} map[string]string "Source code and/or rules not found"
// @Failure 500 {object} map[string]string "Internal server error"
// @Router /lexing/lexer [post]
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

	err := collection.FindOne(ctx, bson.M{"users_id": req.UsersID, "project_name": req.Project_Name}).Decode(&res)
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

// Private struct for the DFA Request
type readDFARequest struct {
	// Represents the states of the automata
	States []string `json:"states"`
	// Represents the transitions of the automata
	Transitions []services.Transition `json:"transitions"`
	// Represents the start state of the automata
	Start string `json:"start_state"`
	// Represents the accepting states of the automata
	Accepting []services.AcceptingState `json:"accepting_states"`
	// Represents the User's ID from frontend
	UsersID bson.ObjectID `json:"users_id" binding:"required"`
	// User's project name
	Project_Name string `json:"project_name" binding:"required"`
}

// @Summary Reads DFA from user
// @Description Takes the user's DFA and stores it with the user's ID. If a DFA already exists, it overwrites that DFA and removes any other created fields (tokens, rules)
// @Tags Lexing
// @Accept json
// @Produce json
// @Param request body readDFARequest true "Read DFA from User"
// @Success 200 {object} map[string]string "DFA successfully and stored"
// @Failure 400 {object} map[string]string "Invalid input"
// @Failure 500 {object} map[string]string "Internal server error"
// @Router /lexing/dfa [post]
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

	filters := bson.M{
		"users_id":     req.UsersID,
		"project_name": req.Project_Name,
	}

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

// @Summary Creates tokens from a stored DFA and source code
// @Description Searches the database for the user's DFA. If found, the DFA and code are used to create the tokens and/or unidentified tokens. The tokens are either created or ,if already existing, updated. If the DFA and/or source code is not found, returns an error
// @Tags Lexing
// @Accept json
// @Produce json
// @Param request body IDRequest true "Create Tokens from Stored DFA"
// @Success 200 {object} map[string]string "Tokens successfully created and stored"
// @Failure 400 {object} map[string]string "Invalid input/Tokenization failed"
// @Failure 404 {object} map[string]string "DFA/source code not found"
// @Failure 500 {object} map[string]string "Internal server error"
// @Router /lexing/dfaToTokens [post]
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

	err := collection.FindOne(ctx, bson.M{"users_id": req.UsersID, "project_name": req.Project_Name}).Decode(&res)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Source code not found"})
		return
	}

	tokens, unidentified, error_caught := services.CreateTokensFromDFA(res.Code, res.DFA)
	if error_caught != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Tokenization from DFA failed", "details": error_caught.Error()})
		return
	}

	filters := bson.M{
		"users_id":     req.UsersID,
		"project_name": req.Project_Name,
	}
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

// @Summary Converts stored DFA to Regular Expressions (Rules)
// @Description Searches the database for the user's DFA. If found, the DFA is used to create the Regular Expressions (Rules). The Rules are either created or ,if already existing, updated. If the DFA is not found, returns an error
// @Tags Lexing
// @Accept json
// @Produce json
// @Param request body IDRequest true "Create Regex from Stored DFA"
// @Success 200 {object} map[string]string "Rules successfully created and stored"
// @Failure 400 {object} map[string]string "Invalid input/Conversion failed"
// @Failure 404 {object} map[string]string "DFA not found"
// @Failure 500 {object} map[string]string "Internal server error"
// @Router /lexing/dfaToRegex [post]
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

	err := collection.FindOne(ctx, bson.M{"users_id": req.UsersID, "project_name": req.Project_Name}).Decode(&res)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "DFA not found. Please create one"})
		return
	}

	rules, error_caught := services.ConvertDFAToRegex(res.DFA)
	if error_caught != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Conversion from DFA to Regex failed", "details": error_caught.Error()})
		return
	}

	filters := bson.M{
		"users_id":     req.UsersID,
		"project_name": req.Project_Name,
	}
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

// @Summary Converts stored Rules to an NFA
// @Description Searches the database for the user's Rules. If found, the Rules are used to create the NFA. The NFA is either created or ,if already existing, updated. If the Rules are not found, returns an error
// @Tags Lexing
// @Accept json
// @Produce json
// @Param request body IDRequest true "Create NFA from Stored Rules"
// @Success 200 {object} map[string]string "NFA successfully created and stored"
// @Failure 400 {object} map[string]string "Invalid input/Conversion failed"
// @Failure 404 {object} map[string]string "Rules not found"
// @Failure 500 {object} map[string]string "Internal server error"
// @Router /lexing/regexToNFA [post]
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

	err := collection.FindOne(ctx, bson.M{"users_id": req.UsersID, "project_name": req.Project_Name}).Decode(&res)
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

	filters := bson.M{"users_id": req.UsersID, "project_name": req.Project_Name}
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

// @Summary Converts stored Rules to an DFA
// @Description Searches the database for the user's Rules. If found, the Rules are used to create the DFA. The DFA is either created or ,if already existing, updated. If the Rules are not found, returns an error
// @Tags Lexing
// @Accept json
// @Produce json
// @Param request body IDRequest true "Create DFA from Stored Rules"
// @Success 200 {object} map[string]string "DFA successfully created and stored"
// @Failure 400 {object} map[string]string "Invalid input/Conversion failed"
// @Failure 404 {object} map[string]string "Rules not found"
// @Failure 500 {object} map[string]string "Internal server error"
// @Router /lexing/regexToDFA [post]
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

	err := collection.FindOne(ctx, bson.M{"users_id": req.UsersID, "project_name": req.Project_Name}).Decode(&res)
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

	filters := bson.M{"users_id": req.UsersID, "project_name": req.Project_Name}
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

// @Summary Converts stored NFA to an DFA
// @Description Searches the database for the user's NFA. If found, the NFA is used to create the DFA. The DFA is either created or ,if already existing, updated. If the NFA is not found, returns an error
// @Tags Lexing
// @Accept json
// @Produce json
// @Param request body IDRequest true "Create DFA from NFA"
// @Success 200 {object} map[string]string "DFA successfully created and stored"
// @Failure 400 {object} map[string]string "Invalid input/Conversion failed"
// @Failure 404 {object} map[string]string "Rules not found"
// @Failure 500 {object} map[string]string "Internal server error"
// @Router /lexing/nfaToDFA [post]
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

	err := collection.FindOne(ctx, bson.M{"users_id": req.UsersID, "project_name": req.Project_Name}).Decode(&res)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "NFA not found. Please create one"})
		return
	}

	dfa, error_caught := services.ConvertNFAToDFA(res.NFA)
	if error_caught != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Conversion from NFA to DFA failed", "details": error_caught.Error()})
		return
	}

	filters := bson.M{"users_id": req.UsersID, "project_name": req.Project_Name}
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
