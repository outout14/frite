package main

import (
	"bufio"
	"flag"
	"fmt"
	"net/http"
	"os"
	"strings"

	log "github.com/sirupsen/logrus"
)

type link struct {
	short string
	to    string
}

type App struct {
	links      []link
	httpIP     string
	httpPort   int
	httpSubdir string
	debug      bool
}

func fixuri(uri string) string {
	if !strings.HasSuffix(uri, "/") {
		uri = uri + "/"
	}
	if !strings.HasPrefix(uri, "/") {
		uri = "/" + uri
	}
	return uri
}

func (a *App) getlinks(path string) error {
	f, err := os.Open(path)
	if err != nil {
		return err
	}
	defer f.Close()

	var links []link

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		data := strings.Fields(scanner.Text())
		if len(data) != 2 {
			log.WithFields(log.Fields{
				"short": data,
			}).Error("CORE> Failed to load short link (invalid syntax):")
		} else {
			l := link{short: data[0], to: data[1]}
			links = append(links, l)
		}
	}

	if err := scanner.Err(); err != nil {
		return err
	}

	a.links = links
	return nil
}

func (a App) httpGetLink(w http.ResponseWriter, r *http.Request) {
	uri := strings.Replace(r.RequestURI, a.httpSubdir, "", 1)
	log.WithFields(log.Fields{
		"uri": uri,
	}).Debugf("HTTP> got request")
	for _, l := range a.links {
		if l.short == uri {
			log.Debug("CORE> short found")
			http.Redirect(w, r, l.to, http.StatusFound)
			return
		}
	}
	log.Debug("CORE> short not found")
	w.WriteHeader(404)
	fmt.Fprintf(w, "404\n")
}

func main() {
	linksPath := flag.String("links", "links.txt", "Path to the links file")
	webPort := flag.Int("http-port", 8080, "HTTP Listen port")
	webIP := flag.String("http-host", "127.0.0.1", "HTTP Listen IP")
	webDir := flag.String("http-dir", "/", "If proxied in subfolder")
	debugState := flag.Bool("debug", false, "Enable debug logs")

	flag.Parse()

	app := App{
		debug:      *debugState,
		httpIP:     *webIP,
		httpPort:   *webPort,
		httpSubdir: fixuri(*webDir),
	}
	log.WithFields(log.Fields{
		"debug": app.debug,
		"http":  fmt.Sprintf("%s:%v%s", app.httpIP, app.httpPort, app.httpSubdir),
	}).Info("CORE> Starting frite-web application")

	if app.debug {
		log.SetLevel(log.DebugLevel)
	}

	app.getlinks(*linksPath)

	log.Infof("CORE> Loaded %v urls", len(app.links))

	log.Infof("HTTP> Listening on http://%s:%v", app.httpIP, app.httpPort)

	http.HandleFunc("/", app.httpGetLink)
	err := http.ListenAndServe(fmt.Sprintf("%s:%v", app.httpIP, app.httpPort), nil)
	if err != nil {
		log.WithFields(log.Fields{
			"err": err,
		}).Error("HTTP:")
	}
}
