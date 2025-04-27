package server

import (
	"net/http"
	"time"

	"github.com/ZivojinMilutinovic/golang-coding-test/store_api"
	"github.com/gin-gonic/gin"
)

type API struct {
	Store *store_api.Store
}

func StartServer() {
	api := &API{
		Store: store_api.NewStore(),
	}
	r := gin.Default()
	api.registerRoutes(r)
	r.Run(":8080")
}

func (api *API) registerRoutes(r *gin.Engine) {
	r.GET("/get/:key", api.handleGet)
	r.POST("/set/:key", api.handleSet)
	r.POST("/update/:key", api.handleUpdate)
	r.DELETE("/remove/:key", api.handleRemove)
	r.POST("/push/:key", api.handlePush)
	r.POST("/pop/:key", api.handlePop)
}

func (api *API) handleGet(c *gin.Context) {
	key := c.Param("key")
	val := api.Store.Get(key)
	if val == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "key not found"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"value": val})
}

func (api *API) handleSet(c *gin.Context) {
	key := c.Param("key")
	var body struct {
		Value interface{} `json:"value"`
		TTL   int         `json:"ttl"` // in seconds
	}
	if err := c.BindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}
	api.Store.Set(key, body.Value, time.Duration(body.TTL)*time.Second)
}

func (api *API) handleUpdate(c *gin.Context) {
	key := c.Param("key")
	var body struct {
		Value interface{} `json:"value"`
	}
	if err := c.BindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}
	ok := api.Store.Update(key, body.Value)
	if !ok {
		c.JSON(http.StatusNotFound, gin.H{"error": "key not found"})
		return
	}
	c.Status(http.StatusOK)
}

func (api *API) handleRemove(c *gin.Context) {
	key := c.Param("key")
	api.Store.Remove(key)
	c.Status(http.StatusOK)
}

func (api *API) handlePush(c *gin.Context) {
	key := c.Param("key")
	var body struct {
		Value string `json:"value"`
	}
	if err := c.BindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}
	api.Store.Push(key, body.Value)
	c.Status(http.StatusOK)
}

func (api *API) handlePop(c *gin.Context) {
	key := c.Param("key")
	val := api.Store.Pop(key)
	if val == "" {
		c.JSON(http.StatusNotFound, gin.H{"error": "list empty or not found"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"value": val})
}
