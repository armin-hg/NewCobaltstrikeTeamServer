package manager

import (
	"NewCsTeamServer/server/manager/admin"
	"NewCsTeamServer/utils"
	"context"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"github.com/kataras/golog"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func RunTeamServer(port string, ssl bool) { //TODO 后续开发成websocket，目前初版使用http
	app := gin.New()
	app.Use(gin.Recovery())
	app.Use(utils.Cors()) //解决跨域问题
	//	app.LoadHTMLGlob("templates/*") //测试用
	app.Static("/web", "./html")
	app.NoRoute(func(c *gin.Context) {
		c.String(404, "404")
	})

	app.GET("/ws", func(c *gin.Context) { //TODO 鉴权
		upgrader := websocket.Upgrader{
			CheckOrigin: func(r *http.Request) bool {
				return true
			},
			HandshakeTimeout: 10 * time.Second,
			ReadBufferSize:   1024,
			WriteBufferSize:  1024,
		}
		conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
		if err != nil {
			log.Printf("管理端 WebSocket 升级失败: %v", err)
			return
		}
		clientIP := c.ClientIP()
		//adminid := admin.UserInfo(c)
		log.Printf("管理端  %s 已连接", clientIP)
		admin2 := admin.GetConnectionManager().AddAdmin(conn)

		admin.HandleAdminMessages(admin2, clientIP) //TODO 还需要传入用户权限
	}) //WebSocket
	srv := &http.Server{
		Addr:    "0.0.0.0:" + port, // HTTPS server will listen on this port
		Handler: app,
		ConnContext: func(ctx context.Context, c net.Conn) context.Context {
			return context.WithValue(ctx, `Conn`, c)
		},
	}
	if ssl {
		go func() {

			if err := srv.ListenAndServeTLS("", ""); err != nil && err != http.ErrServerClosed {
				golog.Fatal(`Failed to bind address: `, err)
			}
		}()
	}
	httpSrv := &http.Server{
		Addr:    "0.0.0.0:" + port, // HTTP server will listen on a different port
		Handler: app,
	}
	go func() {
		if err := httpSrv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			golog.Fatal(`Http Failed to bind address: `, err)
		}
	}()

	quit := make(chan os.Signal, 3)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	golog.Warn(`TeamServer is shutting down ...`)

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		golog.Fatal(`Https Server shutdown: `, err)
	}
	if err := httpSrv.Shutdown(ctx); err != nil {
		golog.Fatal(`Http Server shutdown: `, err)
	}
	<-ctx.Done()
}
