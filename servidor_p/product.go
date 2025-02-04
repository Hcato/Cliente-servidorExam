package servidorp

import (
	"fmt"
	"net/http"
	"net/url"
	"time"

	"github.com/gin-gonic/gin"
)

type Product struct {
	ID      int64  `json:"id"`
	Name    string `json:"name"`
	Amount  string `json:"amount"`
	CodeBar string `json:"codebar"`
}

type Cambio struct {
	Accion string  `json:"accion"`
	User   Product `json:"user"`
}

var (
	bd      []Product
	cambios []Cambio
)

func sendToReplicationServer(product Product, accion string) {
	url := fmt.Sprintf("http://localhost:8800/replication?product_id=%d&name=%s&amount=%s&codebar=%s&accion=%s",
		product.ID,
		url.QueryEscape(product.Name),
		url.QueryEscape(product.Amount),
		url.QueryEscape(product.CodeBar),
		url.QueryEscape(accion),
	)

	client := &http.Client{Timeout: 10 * time.Second}

	resp, err := client.Get(url)
	if err != nil {
		fmt.Println("Error al enviar al servidor de replicación:", err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		fmt.Println("Error replicando el producto:", resp.Status)
	}
}

func sendProductToReplication(c *gin.Context) {
	if len(cambios) > 0 {
		lastChange := cambios[len(cambios)-1]
		sendToReplicationServer(lastChange.User, lastChange.Accion)
		c.JSON(http.StatusOK, gin.H{"mensaje": "Producto enviado a replicación", "product": lastChange.User})
	} else {
		c.JSON(http.StatusOK, gin.H{"mensaje": "No hay cambios nuevos"})
	}
}

func createProduct(c *gin.Context) {
	var newUser Product
	if err := c.ShouldBindJSON(&newUser); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	newUser.ID = int64(len(bd) + 1)
	bd = append(bd, newUser)

	cambios = append(cambios, Cambio{Accion: "create", User: newUser})
	sendToReplicationServer(newUser, "create")

	c.JSON(http.StatusCreated, newUser)
}

func getProducts(c *gin.Context) {
	c.JSON(http.StatusOK, bd)
}

func getCambios(c *gin.Context) {
	if len(cambios) == 0 {
		c.JSON(http.StatusOK, gin.H{"mensaje": "No hay cambios nuevos"})
		return
	}

	response := cambios
	cambios = []Cambio{}

	c.JSON(http.StatusOK, response)
}
func updateProduct(c *gin.Context) {
	var updatedProduct Product
	if err := c.ShouldBindJSON(&updatedProduct); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	for i, p := range bd {
		if p.ID == updatedProduct.ID {
			bd[i] = updatedProduct
			cambios = append(cambios, Cambio{Accion: "update", User: updatedProduct})
			sendToReplicationServer(updatedProduct, "update")
			c.JSON(http.StatusOK, updatedProduct)
			return
		}
	}
	c.JSON(http.StatusNotFound, gin.H{"error": "Producto no encontrado"})
}
func deleteProduct(c *gin.Context) {
	var productToDelete Product
	if err := c.ShouldBindJSON(&productToDelete); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	for i, p := range bd {
		if p.ID == productToDelete.ID {
			bd = append(bd[:i], bd[i+1:]...)
			cambios = append(cambios, Cambio{Accion: "delete", User: productToDelete})
			sendToReplicationServer(productToDelete, "delete")
			c.JSON(http.StatusOK, gin.H{"mensaje": "Producto eliminado"})
			return
		}
	}
	c.JSON(http.StatusNotFound, gin.H{"error": "Producto no encontrado"})
}
