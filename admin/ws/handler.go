package ws

// import (
// 	"net/http"
// 	"strconv"

// 	"github.com/gin-gonic/gin"
// 	"github.com/gorilla/websocket"
// )

// var upgrader = websocket.Upgrader{
// 	CheckOrigin: func(r *http.Request) bool { return true },
// }

// func WebSocketHandler(c *gin.Context) {
// 	userIDParam := c.Query("user_id")
// 	userID, err := strconv.Atoi(userIDParam)
// 	if err != nil {
// 		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid user_id"})
// 		return
// 	}

// 	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
// 	if err != nil {
// 		return
// 	}

// 	WSManager.AddClient(uint(userID), conn)

// 	for {
// 		_, _, err := conn.ReadMessage()
// 		if err != nil {
// 			break
// 		}
// 	}

// 	WSManager.RemoveClient(uint(userID))
// 	conn.Close()
// }
