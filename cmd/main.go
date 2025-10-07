package main

import (
	"context"
	_ "embed"
	"flag"
	"fmt"
	"os"
	"os/signal"
	"path/filepath"
	"spahttp/app"
	"spahttp/checkdbg"
	"spahttp/config"
	"spahttp/domain/models/application"
	"spahttp/reductor"
	"spahttp/repo"
	"spahttp/spaserver"
	"spahttp/zaplog"
	"strings"
	"syscall"

	"github.com/containers/winquit/pkg/winquit"
	"github.com/mechiko/dbscan"
	"golang.org/x/sync/errgroup"
)

const modError = "main"

// var version = "0.0.0"
var fileExe string
var dir string

// если home true то папка создается в home каталоге
var home = flag.Bool("home", false, "")

// set pwd to path exe for deamon
func init() {
	fileExe = os.Args[0]
	var err error
	dir, err = filepath.Abs(filepath.Dir(fileExe))
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to get absolute path: %v\n", err)
		os.Exit(1)
	}
	if err := os.Chdir(dir); err != nil {
		fmt.Fprintf(os.Stderr, "Failed to change directory: %v\n", err)
		os.Exit(1)
	}
}

func main() {
	flag.Parse()

	// ctx, cancel := context.WithCancel(context.Background())
	// defer cancel()

	done := make(chan bool, 1)

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()
	// ctx := context.Background()

	group, groupCtx := errgroup.WithContext(ctx)

	cfg, err := config.New(*home)
	if err != nil {
		errMessageExit(nil, "ошибка конфигурации", err)
	}

	debug := strings.ToLower(config.Mode) == "development"
	var logsOutConfig = map[string]zaplog.LogConfig{
		"logger": {
			ErrorOutputPaths: []string{"stdout", filepath.Join(cfg.LogPath(), config.Name)},
			Debug:            debug,
			Console:          true,
			Name:             filepath.Join(cfg.LogPath(), config.Name),
		},
		"echo": {
			ErrorOutputPaths: []string{filepath.Join(cfg.LogPath(), "echo")},
			Debug:            debug,
			Console:          false,
			Name:             filepath.Join(cfg.LogPath(), "echo"),
		},
		"reductor": {
			ErrorOutputPaths: []string{filepath.Join(cfg.LogPath(), "reductor")},
			Debug:            debug,
			Console:          true,
			Name:             filepath.Join(cfg.LogPath(), "reductor"),
		},
	}

	zl, err := zaplog.New(logsOutConfig)
	if err != nil {
		errMessageExit(nil, "ошибка создания логера", err)
	}

	lg, err := zl.GetLogger("logger")
	if err != nil {
		errMessageExit(nil, "ошибка получения логера", err)
	}
	loger := lg.Sugar()
	loger.Debug("zaplog started")
	loger.Infof("mode = %s", config.Mode)
	if cfg.Warning() != "" {
		loger.Infof("pkg:config warning %s", cfg.Warning())
	}

	group.Go(func() error {
		go func() {
			defer stop()
			<-done
			loger.Info("получен сигнал winquit.NotifyOnQuit завершения работы")
		}()
		return nil
	})

	// используем инлайн функцию для захвата loger
	errProcessExit := func(title string, err error) {
		errMessageExit(loger, title, err)
	}
	// создаем приложение с опциями из конфига и логером основным
	app := app.New(cfg, loger, dir)
	// инициализируем пути необходимые приложению
	app.CreatePath()
	// создаем редуктор для хранения моделей приложения
	reductorLogger, err := zl.GetLogger("reductor")
	if err != nil {
		errProcessExit("Ошибка получения логера для редуктора", err)
	}

	if err := reductor.New(reductorLogger.Sugar()); err != nil {
		errProcessExit("Ошибка создания редуктора", err)
	}

	loger.Info("start repo")
	// инициализируем REPO
	// TODO изменить получение путей из конфига
	listDbs := make(dbscan.ListDbInfoForScan)
	listDbs[dbscan.Other] = &dbscan.DbInfo{
		File:   app.Options().Db.File,
		Driver: app.Options().Db.Driver,
		Path:   app.DbPath(),
	}

	err = repo.New(listDbs, app.DbPath())
	if err != nil {
		errProcessExit("Ошибки запуска репозитория", err)
	}
	repoStart, err := repo.GetRepository()
	if err != nil {
		errProcessExit("Ошибки получения репозитория", err)
	}

	appModel, err := application.New(app)
	if err != nil {
		errProcessExit("Ошибка создания модели для редуктора", err)
	}
	if err := reductor.Instance().SetModel(appModel, false); err != nil {
		errProcessExit("Ошибка редуктора", err)
	}
	group.Go(func() error {
		go func() {
			<-groupCtx.Done()
			err := repoStart.Shutdown()
			loger.Infof("repo shutdown %v", err)
		}()
		return repoStart.Run(groupCtx)
	})
	// тесты
	checker, err := checkdbg.NewChecks(loger, repoStart)
	if err != nil {
		stop()
		// Wait for cleanup to complete
		group.Wait()
		errProcessExit("Check failed", err)
	}
	err = checker.Run()
	if err != nil {
		stop()
		// Wait for cleanup to complete
		group.Wait()
		errProcessExit("Check failed", err)
	}

	loger.Info("start up webapp")

	port := cfg.Configuration().HostPort
	if port == "" || port == "auto" {
		errProcessExit("Ошибка port http server", fmt.Errorf("port %v wrong", port))
	}
	loger.Infof("http port %s", port)

	// тут инициализируются так же модели для всех видов
	spaServerLogger, err := zl.GetLogger("echo")
	if err != nil {
		errProcessExit("Ошибка получения логера для http server", err)
	}
	httpServer, err := spaserver.New(app, spaServerLogger, port, true)
	if err != nil {
		errProcessExit("Ошибка создания http server", err)
	}
	loger.Infof("отладка шаблонов %v", httpServer.TemplateIsDebug())
	loger.Infof("путь шаблонов %s", httpServer.RootPathTemplates())
	// запускаем сервер эхо через него SSE работает для флэш сообщений
	// httpServer.Start()
	group.Go(func() error {
		go func() {
			// предположим, что httpServer (как и http.ListenAndServe, кстати) не умеет останавливаться по отмене
			// контекста, тогда придётся добавить обработку отмены вручную.
			// ошибка у какого то другого члена группы или он завершился принудительно
			<-groupCtx.Done()
			app.Logger().Debugf("%s получен сигнал завершения контекста группы в HTTP", modError)
			if err := httpServer.Shutdown(); err != nil {
				app.Logger().Debugf("%s stopped http server with error: %v", modError, err)
			}
		}()
		httpServer.Start()
		// по ошибке сервера возвращаем в группу код ошибки
		return <-httpServer.Notify()
	})

	// только в винде откроет брауезер на индекс сайта
	openUrl(app)

	// Simulate SIGTERM when a quit occurs
	winquit.NotifyOnQuit(done)

	// ожидание завершения всех в группе
	if err := group.Wait(); err != nil {
		fmt.Printf("game over! error %s\n", err.Error())
	} else {
		fmt.Println("game over!")
	}
	// завершаем все логи
	zl.Shutdown()
}
