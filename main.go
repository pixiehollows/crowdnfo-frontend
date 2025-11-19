package main

import (
	"context"
	"net/http"
	"strconv"

	"github.com/pixiehollows/crowdnfo-frontend/internal/crowdnfo"
	"github.com/pixiehollows/crowdnfo-frontend/internal/templates"

	"github.com/a-h/templ"
	"github.com/labstack/echo/v4"
)

//go:generate go tool oapi-codegen -config internal/crowdnfo/cfg.yaml internal/crowdnfo/crowdnfo.v1.json

func main() {
	e := echo.New()
	crowdnfoClient, err := crowdnfo.NewClientWithResponses("https://crowdnfo.net/api")
	if err != nil {
		e.Logger.Error(err)
	}

	releases, err := crowdnfoClient.GetReleasesWithResponse(context.TODO(), &crowdnfo.GetReleasesParams{})
	if err != nil {
		e.Logger.Error(err)
	}

	releasegroups, err := crowdnfoClient.GetReleasegroupsWithResponse(context.TODO())
	if err != nil {
		e.Logger.Error(err)
	}

	e.Static("/", "public")
	e.GET("/", func(c echo.Context) error {
		index := templates.Index()
		return render(c, index)
	})
	e.GET("/releases", func(c echo.Context) error {
		releases := templates.Releases(*releases.JSON200)
		return render(c, releases)
	})
	e.GET("releases/:id", func(c echo.Context) error {
		id, err := strconv.Atoi(c.Param("id"))
		if err != nil {
			return err //TODO
		}
		release, err := crowdnfoClient.GetReleasesIdWithResponse(context.TODO(), int32(id))
		if err != nil {
			return err //TODO
		}
		releasesId := templates.ReleasesId(*release.JSON200)
		return render(c, releasesId)
	})
	e.GET("/releasegroups", func(c echo.Context) error {
		releasegroups := templates.ReleaseGroups(*releasegroups.JSON200)
		return render(c, releasegroups)
	})
	e.Start(":8080")
}

func render(c echo.Context, component templ.Component) error {
	buf := templ.GetBuffer()
	defer templ.ReleaseBuffer(buf)

	if err := component.Render(c.Request().Context(), buf); err != nil {
		return err
	}
	return c.HTML(http.StatusOK, buf.String())
}
