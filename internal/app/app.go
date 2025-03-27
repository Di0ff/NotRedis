package app

import (
	"NotRedis/internal/compute/parser"
	"NotRedis/internal/storage"
	"bufio"
	"fmt"
	"go.uber.org/zap"
	"os"
	"strings"
	"sync"
)

type App struct {
	parser  *parser.Parser
	storage *storage.Storage
	logger  *zap.Logger
	ch      chan string
	wg      sync.WaitGroup
}

func NewApp(parser *parser.Parser, storage *storage.Storage, logger *zap.Logger) *App {
	return &App{
		parser:  parser,
		storage: storage,
		logger:  logger,
		ch:      make(chan string, 52),
	}
}

func (app *App) handle(request string) error {
	cli, err := app.parser.Parse(request)
	if err != nil {
		return err
	}

	switch cli.Type {
	case "SET":
		err := app.storage.Set(cli.Key, cli.Value)
		if err != nil {
			return err
		}
		fmt.Println("Апрувд")
	case "GET":
		value, err := app.storage.Get(cli.Key)
		if err != nil {
			return err
		}
		fmt.Println(value)
	case "DEL":
		err := app.storage.Delete(cli.Key)
		if err != nil {
			return err
		}
		fmt.Println("Апрувд")
	default:
		return fmt.Errorf("unknown command: %v", cli.Type)
	}

	return nil
}

func (app *App) async() {
	defer app.wg.Done()

	for request := range app.ch {
		if err := app.handle(request); err != nil {
			fmt.Println("Одна ошибка — и ты ошибся!", err)
		}
	}
}

func (app *App) Run() {
	defer close(app.ch)

	app.wg.Add(1)

	go app.async()

	scanner := bufio.NewScanner(os.Stdin)
	fmt.Println("NotRedis готов к работе)")

	for scanner.Scan() {
		app.ch <- strings.TrimSpace(scanner.Text())
	}

	if err := scanner.Err(); err != nil {
		app.logger.Error("Scan fail", zap.Error(err))
	}

	app.wg.Wait()
}
