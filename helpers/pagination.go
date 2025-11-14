package helpers

import (
	"backend-city/structs"
	"net/http"
	"reflect"

	"strconv"

	"github.com/gin-gonic/gin"
)

// Struktur untuk link pagination (seperti Laravel)
type PaginationLink struct {
	URL    string `json:"url"`
	Label  string `json:"label"`
	Active bool   `json:"active"`
}

// Konversi string ke integer dengan fallback ke 1 jika gagal atau < 1
func StringToInt(str string) int {
	i, err := strconv.Atoi(str)
	if err != nil || i < 1 {
		return 1
	}
	return i
}

// Hitung jumlah halaman berdasarkan total data dan batas per halaman
func TotalPage(total int64, perPage int) int {
	if perPage == 0 {
		return 1
	}
	pages := int(total) / perPage
	if int(total)%perPage != 0 {
		pages++
	}
	return pages
}

// Buat link pagination seperti Laravel
func BuildPaginationLinks(currentPage, lastPage int, baseURL, search string) []PaginationLink {
	links := []PaginationLink{}

	// Link ke halaman sebelumnya
	links = append(links, PaginationLink{
		URL:    PageURL(baseURL, currentPage-1, lastPage, search),
		Label:  "&laquo; Previous",
		Active: false,
	})

	// Link ke semua halaman
	for i := 1; i <= lastPage; i++ {
		links = append(links, PaginationLink{
			URL:    baseURL + "?page=" + strconv.Itoa(i) + QueryString(search),
			Label:  strconv.Itoa(i),
			Active: i == currentPage,
		})
	}

	// Link ke halaman berikutnya
	links = append(links, PaginationLink{
		URL:    PageURL(baseURL, currentPage+1, lastPage, search),
		Label:  "Next &raquo;",
		Active: false,
	})

	return links
}

// Buat URL untuk halaman tertentu (prev/next)
func PageURL(baseURL string, page, lastPage int, search string) string {
	if page < 1 || page > lastPage {
		return ""
	}
	return baseURL + "?page=" + strconv.Itoa(page) + QueryString(search)
}

// Tambahkan query string jika ada parameter pencarian
func QueryString(search string) string {
	if search == "" {
		return ""
	}
	return "&search=" + search
}

// Ambil parameter pencarian dan pagination dari query string
func GetPaginationParams(c *gin.Context) (search string, page, limit, offset int) {
	search = c.Query("search")
	page = StringToInt(c.DefaultQuery("page", "1"))
	limit = StringToInt(c.DefaultQuery("limit", "5"))
	offset = (page - 1) * limit
	return
}

// Bangun URL dasar untuk pagination (mempertimbangkan HTTP/HTTPS dan proxy)
func BuildBaseURL(c *gin.Context) string {
	// Cek X-Forwarded-Proto dari reverse proxy jika ada
	scheme := c.Request.Header.Get("X-Forwarded-Proto")
	if scheme == "" {
		if c.Request.TLS != nil {
			scheme = "https"
		} else {
			scheme = "http"
		}
	}
	return scheme + "://" + c.Request.Host + c.Request.URL.Path
}

// Kirim response JSON pagination lengkap seperti Laravel
func PaginateResponse(c *gin.Context, data interface{}, total int64, page, limit int, baseURL, search, message string) {
	lastPage := TotalPage(total, limit)
	from := (page-1)*limit + 1
	to := from + reflect.ValueOf(data).Len() - 1

	links := BuildPaginationLinks(page, lastPage, baseURL, search)

	c.JSON(http.StatusOK, structs.SuccessResponse{
		Success: true,
		Message: message,
		Data: gin.H{
			"current_page":   page,
			"data":           data,
			"first_page_url": baseURL + "?page=1" + QueryString(search),
			"from":           from,
			"last_page":      lastPage,
			"last_page_url":  baseURL + "?page=" + strconv.Itoa(lastPage) + QueryString(search),
			"links":          links,
			"next_page_url":  PageURL(baseURL, page+1, lastPage, search),
			"path":           baseURL,
			"per_page":       limit,
			"prev_page_url":  PageURL(baseURL, page-1, lastPage, search),
			"to":             to,
			"total":          total,
		},
	})
}
