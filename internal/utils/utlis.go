package utils

import (
	"html/template"
	"io"
	"net"
	"os"
	"path/filepath"
	"strconv"
	"time"

	"github.com/labstack/echo/v4"
)

// TemplateRenderer implements echo.Renderer
type TemplateRenderer struct {
	templates *template.Template
}

// NewTemplateRenderer creates a new TemplateRenderer
func NewTemplateRenderer() (*TemplateRenderer, error) {
	rootDir, err := os.Getwd()
	if err != nil {
		return nil, err
	}

	templateDir := filepath.Join(rootDir, "templates")
	t := template.New("")

	err = filepath.Walk(templateDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() && filepath.Ext(path) == ".html" {
			_, err = t.ParseFiles(path)
			if err != nil {
				return err
			}
		}
		return nil
	})

	if err != nil {
		return nil, err
	}

	return &TemplateRenderer{templates: t}, nil
}

// Render implements echo.Renderer
func (tr *TemplateRenderer) Render(w io.Writer, name string, data interface{}, c echo.Context) error {
	return tr.templates.ExecuteTemplate(w, name, data)
}

// used to auto detect the active local IP address - not used yet
func GetLocalIP() net.IP {

	conn, err := net.Dial("udp", "8.8.8.8:80")
	if err != nil {
		return nil
		//log.Fatal(err)
	}
	defer conn.Close()

	localAddress := conn.LocalAddr().(*net.UDPAddr)

	return localAddress.IP
}

func IsValidExamCode(code string) bool {
	//validate the examid
	//2026S1ITCS5.100
	yr := code[:3]         //year component
	sem := code[4:5]       //semester component
	coursecode := code[6:] //	coursecode component

	//e.g 2026
	y, err := strconv.Atoi(yr)
	if err != nil || len(yr) != 4 || y < time.Now().Year() {
		return false
	}

	//e.g S1 or S2 or S3
	if sem != "S1" && sem != "S2" && sem != "S3" {
		return false
	}

	//e.g ITCS5.100
	if len(coursecode) < 9 {
		return false
	}
	return true
}

// TODO: implement a more robust course code validation function
func IsValidCourseCode(code string) bool {
	return len(code) == 9
}

//basic sets for validating status fields
/*
	stat = utils.statusSet{}
	stat.add("active")
	stat.add("closed")
	if stat.has(status) {
		// valid status
	} else {
		// invalid status
	}

*/

type StatusSet map[string]struct{}

func (s StatusSet) Add(value string) {
	s[value] = struct{}{}
}

func (s StatusSet) Remove(value string) {
	delete(s, value)
}

func (s StatusSet) Has(value string) bool {
	_, ok := s[value]
	return ok
}
