package server

import (
	"log/slog"
	"net/http"
	"strings"

	"github.com/danielgtaylor/huma/v2"
	"github.com/ibbycs/private-markets-api/internal/config"
	"github.com/ibbycs/private-markets-api/internal/handler"
	"github.com/ibbycs/private-markets-api/internal/humaecho"
	"github.com/ibbycs/private-markets-api/internal/repository"
	"github.com/ibbycs/private-markets-api/internal/service"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/labstack/echo/v5"
	"github.com/labstack/echo/v5/middleware"
)

func NewServer(cfg *config.Config, logger *slog.Logger, db *pgxpool.Pool) (*echo.Echo, error) {
	e := echo.New()

	e.Use(middleware.RequestLogger())
	e.Use(middleware.RequestID())
	e.Use(middleware.Recover())

	e.GET("/", func(c *echo.Context) error {
		return c.JSON(http.StatusOK, map[string]string{
			"name": "private-markets-api",
			"env":  cfg.Env,
		})
	})

	config := huma.DefaultConfig("Private Markets API", "1.0.0")
	config.DocsPath = ""

	api := humaecho.New(e, config)

	e.GET("/docs", func(c *echo.Context) error {
		csp := []string{
			"default-src 'none'",
			"base-uri 'none'",
			"connect-src 'self'",
			"form-action 'none'",
			"frame-ancestors 'none'",
			"sandbox allow-same-origin allow-scripts",
			"script-src 'unsafe-eval' https://unpkg.com/@scalar/api-reference@1.44.20/dist/browser/standalone.js",
			"style-src 'unsafe-inline'",
		}
		c.Response().Header().Set("Content-Security-Policy", strings.Join(csp, "; "))
		return c.HTML(http.StatusOK, `
			<!doctype html>
			<html lang="en">
				<head>
					<meta charset="utf-8">
					<meta name="referrer" content="no-referrer">
					<meta name="viewport" content="width=device-width, initial-scale=1">
					<title>API Reference</title>
				</head>
				<body>
					<script id="api-reference" data-url="/openapi.json"></script>
					<script src="https://unpkg.com/@scalar/api-reference@1.44.20/dist/browser/standalone.js" crossorigin integrity="sha384-tMz7GAo6dMy55x9tLFtH+sHtogji6Scmb+feBR31TAHmvSPRUTboK9H3M5NFaP4R"></script>
				</body>
			</html>`)
	})

	queries := repository.New(db)

	fundSvc := service.NewFundService(queries)
	handler.NewFundHandler(api, fundSvc)

	investorSvc := service.NewInvestorService(queries)
	handler.NewInvestorHandler(api, investorSvc)

	investmentSvc := service.NewInvestmentService(queries, fundSvc, investorSvc)
	handler.NewInvestmentHandler(api, investmentSvc)

	return e, nil
}
