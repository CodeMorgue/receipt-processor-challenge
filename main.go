package main

import (
	"math"
	"net/http"
	"regexp"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// Define a map of receipts that stores a UUID as a key, and a receipt object as the value
var receiptMap map[string]Receipt

// Define an items struct to be used in the Receipt struct
type ItemsList struct {
	ShortDescription string `json:"shortDescription"`
	Price            string `json:"price"`
}

// Define a Receipt struct to bind the given json object to
type Receipt struct {
	Retailer     string      `json:"retailer"`
	PurchaseDate string      `json:"purchaseDate"`
	PurchaseTime string      `json:"purchaseTime"`
	Total        string      `json:"total"`
	Items        []ItemsList `json:"items"`
}

// When a json object is posted, bind its values to a new Receipt and generate a UUID. Add the UUID to the receiptMap as a key and the Receipt as a value
func createId(c *gin.Context) {
	var newReceipt Receipt
	if err := c.BindJSON(&newReceipt); err != nil {
		return
	}
	id_string := (uuid.New()).String()
	receiptMap[id_string] = newReceipt
	c.IndentedJSON(http.StatusCreated, gin.H{"id": id_string})
}

// When a GET request is made with the appropriate path with an id, find the associated Receipt and calculate its point total
// If receipt does not exist, return http.StatusNotFound code
func processReceipt(c *gin.Context) {
	id := c.Param("id")
	for key, value := range receiptMap {
		if key == id {
			c.IndentedJSON(http.StatusOK, gin.H{"points": calculatePoints(value)})
			return
		}
	}
	c.IndentedJSON(http.StatusNotFound, gin.H{"message": "receipt not found"})
}

// Takes a Receipt as an argument and calculates its point total using function calls
func calculatePoints(r Receipt) int {
	pointTotal := 0

	pointTotal += calculateAlphanumeric(r.Retailer)
	pointTotal += calculateRoundTotal(r.Total)
	pointTotal += calculateMultipleTotal(r.Total)
	pointTotal += calculateItemCount(r.Items)
	pointTotal += calculateDescLength(r.Items)
	pointTotal += calculatePurchaseDay(r.PurchaseDate)
	pointTotal += calculatePurchaseTime(r.PurchaseTime)

	return pointTotal
}

// Calculates the number of alphanumeric values in the retailer name of a Receipt
func calculateAlphanumeric(s string) int {
	var isAlpha = regexp.MustCompile(`^[a-zA-Z0-9]*$`).MatchString
	if isAlpha(s) == true {
		return len(s)
	} else {
		var result strings.Builder
		for i := 0; i < len(s); i++ {
			b := s[i]
			if ('a' <= b && b <= 'z') || ('A' <= b && b <= 'Z') || ('0' <= b && b <= '9') {
				result.WriteByte(b)
			}
		}
		return len(result.String())
	}
}

// Caclulates if the Receipt total is a rounded value with no decimals and returns the appropriate point value
func calculateRoundTotal(s string) int {
	val, err := strconv.ParseFloat(s, 32)
	if err == nil {
		if math.Mod(val, 1.0) == 0 {
			return 50
		} else {
			return 0
		}
	} else {
		return 0
	}
}

// Calculates if the total is a multiple of 0.25 and returns the appropriate point value
func calculateMultipleTotal(s string) int {
	val, err := strconv.ParseFloat(s, 32)
	if err == nil {
		if math.Mod(val, 0.25) == 0 {
			return 25
		} else {
			return 0
		}
	} else {
		return 0
	}
}

// Calculates the number of pairs of items on a Receipt and returns the appropriate point value
func calculateItemCount(l []ItemsList) int {
	length := len(l)
	return ((length - (length % 2)) / 2) * 5
}

func calculateDescLength(l []ItemsList) int {
	total := 0
	for i := 0; i < len(l); i++ {
		b := l[i]
		if len(strings.TrimSpace(b.ShortDescription))%3 == 0 {
			val, err := strconv.ParseFloat(b.Price, 32)
			if err == nil {
				total += int(math.Ceil(val * 0.2))
			}
		}
	}
	return total
}

// Calculates if the purchase day on a Receipt is odd and returns the appropriate point value
func calculatePurchaseDay(s string) int {
	val, err := strconv.Atoi(string(s[9]))
	if err != nil {
		return 0
	}
	if val%2 != 0 {
		return 6
	}
	return 0
}

// Calculates if the purchase time on a Receipt is after 2:00 pm and before 4:00 pm and returns the appropriate point value
func calculatePurchaseTime(s string) int {
	val, err := strconv.Atoi(string(s[0:2]))
	val2, err2 := strconv.Atoi(string(s[3]))

	if err != nil || err2 != nil {
		return 0
	}
	if (val < 14) || (val >= 16) {
		return 0
	} else if val == 14 && val2 != 0 {
		return 10
	} else {
		return 10
	}
}

// Initializes the receiptMap, gin router, and http request paths, then runs the program listening on localhost:8080
func main() {
	receiptMap = make(map[string]Receipt)
	router := gin.Default()
	router.POST("/receipts/process", createId)
	router.GET("/receipts/:id/points", processReceipt)
	router.Run("localhost:8080")
}
