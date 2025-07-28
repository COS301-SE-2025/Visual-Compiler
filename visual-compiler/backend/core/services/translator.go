package services

import (
	"encoding/json"
	"fmt"
	"regexp"
	"strings"
)

// Struct for the translation rules
type TranslationRule struct {
	Sequence    []string
	Translation []string
}