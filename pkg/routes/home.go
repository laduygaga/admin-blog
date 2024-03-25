package routes

import (
	"fmt"
	"strings"
	"time"

	"github.com/mikestefanello/pagoda/ent"
	"github.com/mikestefanello/pagoda/pkg/controller"
	"github.com/mikestefanello/pagoda/templates"

	"github.com/labstack/echo/v4"
)

type (
	home struct {
		controller.Controller
	}
)

func (c *home) Get(ctx echo.Context) error {
	page := controller.NewPage(ctx)
	page.Layout = templates.LayoutMain
	page.Name = templates.PageHome
	page.Metatags.Description = "Welcome to the homepage."
	page.Metatags.Keywords = []string{"Go", "MVC", "Web", "Software"}
	page.Pager = controller.NewPager(ctx, 4)
	page.Data = fetchPosts(c.Controller, ctx, &page.Pager)

	return c.RenderPage(ctx, page)
}

func fetchPosts(c controller.Controller, ctx echo.Context, pager *controller.Pager) []ent.Post {
	from := ctx.QueryParam("from")
	to := ctx.QueryParam("to")
	daterange := ctx.QueryParam("daterange")
	var total int
	var p []*ent.Post
	if from == "" && to == "" && daterange == "" {
		total, p = getPosts(c, ctx, pager)
	}
	if from != "" && to != "" {
		fromDate, err := time.Parse("2006-01-02", from)
		if err != nil {
			return nil
		}
		toDate, err := time.Parse("2006-01-02", to)
		if err != nil {
			return nil
		}
		total, p = getPostsByDate(c, ctx, pager, fromDate, toDate)
	}
	if daterange != "" {
		from := strings.Split(daterange, "-")[0]
		to := strings.Split(daterange, "-")[1]
		fromDate, err := time.Parse("01/02/2006", strings.TrimSpace(from))
		if err != nil {
			return nil
		}
		toDate, err := time.Parse("01/02/2006", strings.TrimSpace(to))
		if err != nil {
			return nil
		}
		total, p = getPostsByDate(c, ctx, pager, fromDate, toDate)
	}

	pager.SetItems(total)
	if pager.Page < 1 {
		pager.Page = 1
	}
	posts := make([]ent.Post, len(p))
	for k, v := range p {
		if len(v.Title) > 30 {
			v.Title = v.Title[:30] + "..."
		}
		if len(v.Body) > 80 {
			v.Body = v.Body[:80] + "..."
		}
		posts[k] = ent.Post{
			ID:    v.ID,
			Title: fmt.Sprintf("%s", v.Title),
			Body:  fmt.Sprintf("%s", v.Body),
		}
	}
	return posts
}

func routeToPostPage(c controller.Controller, ctx echo.Context, p int) error {
	page := controller.NewPage(ctx)
	page.Layout = templates.LayoutMain
	page.Name = templates.PageHome
	page.Metatags.Description = "Posts Page"
	page.Pager = controller.NewPager(ctx, 4)
	page.Pager.Page = p
	page.Data = fetchPosts(c, ctx, &page.Pager)

	return c.RenderPage(ctx, page)
}
