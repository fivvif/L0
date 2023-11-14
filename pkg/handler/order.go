package handler

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func (h *Handler) getByUID(c *gin.Context) {
	uid := c.Param("uid")
	order, err := h.service.Cache.GetCache(uid)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, err.Error())
		return
	}
	c.HTML(http.StatusOK, "index.html", gin.H{
		"OrderUID":          order.OrderUID,
		"TrackNumber":       order.TrackNumber,
		"Entry":             order.Entry,
		"Locale":            order.Locale,
		"InternalSignature": order.InternalSignature,
		"CustomerID":        order.CustomerID,
		"DeliveryService":   order.DeliveryService,
		"Shardkey":          order.Shardkey,
		"SMID":              order.SMID,
		"DateCreated":       order.DateCreated.Format("2006-01-02 15:04:05"),
		"OofShard":          order.OofShard,
		"Delivery":          order.Delivery,
		"Payment":           order.Payment,
		"Items":             order.Items,
	})

}
