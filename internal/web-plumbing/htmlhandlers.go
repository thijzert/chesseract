package plumbing

import (
	"bytes"
	"fmt"
	"html/template"
	"io"
	"net/http"
	"path"
	"strings"

	"github.com/thijzert/chesseract/chesseract"
	"github.com/thijzert/chesseract/internal/assets"
	weberrors "github.com/thijzert/chesseract/internal/web-plumbing/errors"
	"github.com/thijzert/chesseract/web"
)

type htmlHandler struct {
	Server       *Server
	TemplateName string
	Handler      web.Handler
}

// HTMLFunc creates a HTTP handler that outputs HTML
func (s *Server) HTMLFunc(handler web.Handler, templateName string) http.Handler {
	return htmlHandler{
		Server:       s,
		TemplateName: templateName,
		Handler:      handler,
	}
}

var csp string

func init() {
	csp = ""
	csp += "default-src 'self' blob: data: ; "
	csp += "script-src 'self' blob: ; "
	csp += "style-src 'self' data: 'unsafe-inline'; "
	csp += "img-src 'self' blob: data: ; "
	csp += "connect-src 'self' blob: data: ; "
	csp += "frame-src 'none' ; "
	csp += "frame-ancestors 'none'; "
	csp += "form-action 'self'; "
}

func (h htmlHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	req, err := h.Handler.DecodeRequest(r)
	if err != nil {
		h.Error(w, r, err)
		return
	}

	tpl, err := h.Server.getTemplate(h.TemplateName)
	if err != nil {
		h.Error(w, r, err)
		return
	}

	provider, sessionID := h.Server.getProvider(r)
	if _, ok := h.Handler.(sessionlessHandler); !ok {
		if sessionID.IsEmpty() {
			h.Error(w, r, errNoSession)
			return
		}
	}

	resp, err := h.Handler.HandleRequest(provider, req)
	if err != nil {
		h.Error(w, r, err)
		return
	}

	w.Header()["Content-Type"] = []string{"text/html; charset=UTF-8"}

	w.Header()["Content-Security-Policy"] = []string{csp}
	w.Header()["X-Frame-Options"] = []string{"deny"}
	w.Header()["X-XSS-Protection"] = []string{"1; mode=block"}
	w.Header()["Referrer-Policy"] = []string{"strict-origin-when-cross-origin"}
	w.Header()["X-Content-Type-Options"] = []string{"nosniff"}

	tpData := struct {
		Version       string
		AppRoot       string
		AssetLocation string
		PageCSS       string
		Request       web.Request
		Response      web.Response
	}{
		Version:       chesseract.PackageVersion,
		AppRoot:       h.appRoot(r),
		AssetLocation: h.appRoot(r) + "assets",
		Request:       req,
		Response:      resp,
	}

	cssTemplate := h.TemplateName
	if len(cssTemplate) > 5 && cssTemplate[0:5] == "full/" {
		cssTemplate = cssTemplate[5:]
	}

	if _, err := assets.GetAsset(path.Join("web", "dist", "css", "pages", cssTemplate+".css")); err == nil {
		tpData.PageCSS = cssTemplate
	}

	var b bytes.Buffer
	err = tpl.ExecuteTemplate(&b, "basePage", tpData)
	if err != nil {
		h.Error(w, r, err)
		return
	}
	io.Copy(w, &b)
}

func (s *Server) getTemplate(name string) (*template.Template, error) {
	if s.parsedTemplates == nil {
		s.parsedTemplates = make(map[string]*template.Template)
	}

	if assets.AssetsEmbedded() {
		if tp, ok := s.parsedTemplates[name]; ok {
			return tp, nil
		}
	}

	var tp *template.Template

	b, err := assets.GetAsset(path.Join("web", "templates", name+".html"))
	if err != nil {
		return nil, err
	}

	funcs := template.FuncMap{
		"add":     templateAddi,
		"mul":     templateMuli,
		"addf":    templateAddf,
		"mulf":    templateMulf,
		"urlfrag": template.URLQueryEscaper,
	}

	if name == "full/basePage" {
		tp, err = template.New("basePage").Funcs(funcs).Parse(string(b))
		if err != nil {
			return nil, err
		}
	} else if len(name) > 5 && name[:5] == "full/" {
		basePage, err := s.getTemplate("full/basePage")
		if err != nil {
			return nil, err
		}

		tp, err = basePage.Clone()
		if err != nil {
			return nil, err
		}

		_, err = tp.Parse(string(b))
		if err != nil {
			return nil, err
		}
	} else {
		tp, err = template.New("basePage").Funcs(funcs).Parse(string(b))
		if err != nil {
			return nil, err
		}
	}

	s.parsedTemplates[name] = tp
	return tp, nil
}

func templateMuli(a, b int) int {
	return a * b
}
func templateAddi(a, b int) int {
	return a + b
}
func templateMulf(a, b float64) float64 {
	return a * b
}
func templateAddf(a, b float64) float64 {
	return a + b
}

// appRoot finds the relative path to the application root
func (htmlHandler) appRoot(r *http.Request) string {
	// Find the relative path for the application root by counting the number of slashes in the relative URL
	c := strings.Count(r.URL.Path, "/") - 1
	if c == 0 {
		return "./"
	}
	return strings.Repeat("../", c)
}

func (h htmlHandler) Error(w http.ResponseWriter, r *http.Request, err error) {
	st, _ := weberrors.HTTPStatusCode(err)

	if redir, ok := err.(weberrors.Redirector); ok {
		if st == 0 {
			st = 302
		}
		h.Server.redirect(w, r, h.appRoot(r)+redir.RedirectLocation(), st)
		return
	}

	if st == 0 {
		st = 500
	}

	headline := weberrors.Headline(err)
	message := weberrors.Message(err)

	w.WriteHeader(st)
	fmt.Fprintf(w, "Error: %s - %s", headline, message)
}
