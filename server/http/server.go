package http

import (
	"NewCsTeamServer/config"
	"NewCsTeamServer/utils"
	"context"
	"github.com/gin-gonic/gin"
	"github.com/kataras/golog"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func HttpServer(port string, ssl bool) {
	app := gin.New()
	app.Use(gin.Recovery())
	app.Use(utils.Cors()) //解决跨域问题
	app.NoRoute(func(c *gin.Context) {
		c.String(404, "404")
	})
	app.GET(config.ProfileConfig.GetUrl, HandleBeacon)         //处理Beacon请求
	app.POST(config.ProfileConfig.PostUrl, HandleBeaconResult) //处理Beacon请求
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
	golog.Warn(`HttpServer is shutting down ...`)

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
