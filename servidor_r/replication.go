package servidorr

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type Product struct {
	ID      int64  `json:"id"`
	Name    string `json:"name"`
	Amount  string `json:"amount"`
	CodeBar string `json:"codebar"`
}

var bdReplication []Product

func getReplicatedProducts(c *gin.Context) {
	productIDStr := c.DefaultQuery("product_id", "")
	name := c.DefaultQuery("name", "")
	amount := c.DefaultQuery("amount", "")
	codebar := c.DefaultQuery("codebar", "")
	accion := c.DefaultQuery("accion", "")

	if productIDStr != "" && name != "" && amount != "" && codebar != "" && accion != "" {
		productID, err := strconv.ParseInt(productIDStr, 10, 64)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid product_id"})
			return
		}

		var newBdReplication []Product
		for _, p := range bdReplication {
			if p.ID != productID {
				newBdReplication = append(newBdReplication, p)
			}
		}
		bdReplication = newBdReplication

		if accion == "delete" {
			fmt.Println("Product deleted from replication:", productID)
			c.JSON(http.StatusOK, gin.H{"message": "Product deleted from replication", "id": productID})
			return
		}

		newProduct := Product{ID: productID, Name: name, Amount: amount, CodeBar: codebar}
		bdReplication = append(bdReplication, newProduct)
		fmt.Println("Product replicated:", newProduct)
	}

	c.JSON(http.StatusOK, bdReplication)
}
