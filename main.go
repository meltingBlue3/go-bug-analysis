package main

import (
	"context"
	"flag"
	"fmt"
	"io/fs"
	"log"
	"net/http"
	"os"
	"os/exec"
	"os/signal"
	"runtime"
	"time"

	"go-bug-analysis/internal/server"
	"go-bug-analysis/web"
)

func main() {
	port := flag.Int("port", 18088, "HTTP 服务端口")
	flag.Parse()

	// 从嵌入的文件系统中剥离 "static" 前缀
	staticFS, err := fs.Sub(web.StaticFiles, "static")
	if err != nil {
		log.Fatalf("无法加载静态资源: %v", err)
	}

	state := &server.AppState{}
	handler := server.New(staticFS, state)

	addr := fmt.Sprintf(":%d", *port)
	srv := &http.Server{
		Addr:    addr,
		Handler: handler,
	}

	// 在 goroutine 中启动 HTTP 服务器
	go func() {
		log.Printf("服务已启动: http://localhost:%d\n", *port)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("服务启动失败: %v", err)
		}
	}()

	// 延迟 300ms 后自动打开浏览器
	go func() {
		time.Sleep(300 * time.Millisecond)
		openBrowser(fmt.Sprintf("http://localhost:%d", *port))
	}()

	// 监听中断信号，优雅退出
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	<-quit

	log.Println("正在关闭服务...")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("服务关闭失败: %v", err)
	}
	log.Println("服务已退出")
}

// openBrowser 根据操作系统打开默认浏览器
func openBrowser(url string) {
	var cmd *exec.Cmd

	switch runtime.GOOS {
	case "windows":
		cmd = exec.Command("cmd", "/c", "start", url)
	case "darwin":
		cmd = exec.Command("open", url)
	default:
		cmd = exec.Command("xdg-open", url)
	}

	// 忽略错误 — 浏览器打不开不应阻止服务运行
	_ = cmd.Start()
}
