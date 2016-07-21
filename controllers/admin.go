package controllers

import (
	"net/http"
	"runtime"
	"strconv"
	"time"

	"github.com/Cloakaac/cloak/models"
	"github.com/julienschmidt/httprouter"
	"github.com/raggaer/pigo"
)

type AdminController struct {
	*pigo.Controller
}

// Dashboard shows the admin dashboard
func (base *AdminController) Dashboard(w http.ResponseWriter, req *http.Request, ps httprouter.Params) {
	adminInfo := models.GetAdminInformation()
	m := &runtime.MemStats{}
	runtime.ReadMemStats(m)
	onlineRecords, err := models.GetOnlineRecords(10)
	if err != nil {
		base.Error = "Error while trying to get online records"
		return
	}
	base.Data["Records"] = onlineRecords
	base.Data["Memstats"] = m
	base.Data["Goversion"] = runtime.Version()
	base.Data["Numcpu"] = runtime.NumCPU()
	base.Data["Numroutine"] = runtime.NumGoroutine()
	base.Data["Numcgo"] = runtime.NumCgoCall()
	base.Data["AccountTotal"] = adminInfo.Accounts
	base.Data["PlayerTotal"] = adminInfo.Players
	base.Data["MaleTotal"] = adminInfo.Males
	base.Data["FemaleTotal"] = adminInfo.Females
	base.Data["SorcererTotal"] = adminInfo.Sorcerers
	base.Data["DruidTotal"] = adminInfo.Druids
	base.Data["PaladinTotal"] = adminInfo.Paladins
	base.Data["KnightTotal"] = adminInfo.Knights
	base.Template = "admin.html"
}

// Server shows the TFS server manager
func (base *AdminController) Server(w http.ResponseWriter, req *http.Request, ps httprouter.Params) {
	base.Template = "admin_server.html"
}

// ArticleList shows the news admin manager
func (base *AdminController) ArticleList(w http.ResponseWriter, req *http.Request, _ httprouter.Params) {
	articles, err := models.GetArticles(100)
	if err != nil {
		base.Error = err.Error()
		return
	}
	base.Data["Success"] = base.Session.GetFlashes("success")
	base.Data["News"] = articles
	base.Template = "admin_news.html"
}

// ArticleEdit shows the form to edit an article
func (base *AdminController) ArticleEdit(w http.ResponseWriter, req *http.Request, ps httprouter.Params) {
	articleID, err := strconv.ParseInt(ps.ByName("id"), 10, 64)
	if err != nil {
		base.Error = err.Error()
		return
	}
	article := models.GetArticle(articleID)
	if article.ID == -1 {
		base.Redirect = "/admin/news"
		return
	}
	base.Data["Article"] = article
	base.Template = "admin_news_edit.html"

}

// ArticleEditProcess process the article edit form
func (base *AdminController) ArticleEditProcess(w http.ResponseWriter, req *http.Request, ps httprouter.Params) {
	articleID, err := strconv.ParseInt(ps.ByName("id"), 10, 64)
	if err != nil {
		base.Error = err.Error()
		return
	}
	article := models.GetArticle(articleID)
	if article.ID == -1 {
		base.Redirect = "/admin/news"
		return
	}
	article.Text = req.FormValue("text")
	article.Title = req.FormValue("title")
	if err := article.Update(); err != nil {
		base.Error = err.Error()
		return
	}
	base.Session.AddFlash("Article edited successfully", "success")
	base.Redirect = "/admin/news"
}

// ArticleCreate shows the form to create a new article
func (base *AdminController) ArticleCreate(w http.ResponseWriter, req *http.Request, _ httprouter.Params) {
	base.Template = "admin_news_create.html"
}

// ArticleCreateProcess process the article create form
func (base *AdminController) ArticleCreateProcess(w http.ResponseWriter, req *http.Request, _ httprouter.Params) {
	article := models.NewArticle()
	article.Text = req.FormValue("text")
	article.Title = req.FormValue("title")
	article.Created = time.Now().Unix()
	if err := article.Insert(); err != nil {
		base.Error = err.Error()
		return
	}
	base.Session.AddFlash("Article created successfully", "success")
	base.Redirect = "/admin/news"
}

// ArticleDelete deletes the given article
func (base *AdminController) ArticleDelete(w http.ResponseWriter, req *http.Request, ps httprouter.Params) {
	articleID, err := strconv.ParseInt(ps.ByName("id"), 10, 64)
	if err != nil {
		base.Error = err.Error()
		return
	}
	article := models.GetArticle(articleID)
	if article.ID == -1 {
		base.Redirect = "/admin/news"
		return
	}
	err = article.Delete()
	if err != nil {
		base.Error = err.Error()
		return
	}
	base.Session.AddFlash("Article deleted successfully", "success")
	base.Redirect = "/admin/news"
}

// ShopCategories shows the donation shop categories
func (base *AdminController) ShopCategories(w http.ResponseWriter, req *http.Request, _ httprouter.Params) {
	categories, err := models.GetCategories()
	if err != nil {
		base.Error = err.Error()
		return
	}
	base.Data["Success"] = base.Session.GetFlashes("success")
	base.Data["Errors"] = base.Session.GetFlashes("errors")
	base.Data["Categories"] = categories
	base.Template = "admin_shop_categories.html"
}

// CreateCategory shows the form to create a new shop category
func (base *AdminController) CreateCategory(w http.ResponseWriter, req *http.Request, _ httprouter.Params) {
	base.Template = "admin_shop_categories_create.html"
}

// CreateCategoryProcess creates a new shop category
func (base *AdminController) CreateCategoryProcess(w http.ResponseWriter, req *http.Request, _ httprouter.Params) {
	category := models.GetCategory(req.FormValue("name"))
	if category.ID != -1 {
		base.Session.AddFlash("Category name already in use", "errors")
		base.Redirect = "/admin/shop/categories"
		return
	}
	category.Name = req.FormValue("name")
	category.Description = req.FormValue("desc")
	err := category.Insert()
	if err != nil {
		base.Error = err.Error()
		return
	}
	base.Session.AddFlash("Category created successfully", "success")
	base.Redirect = "/admin/shop/categories"
}
