package main

import (
	"flag";
	"http";
	"io";
	"log";
	"strings";
	"template";
)

var addr = flag.String("addr", ":1718", "http service address") // Q=17, R=18
var fmap = template.FormatterMap{
	"html": template.HTMLFormatter,
	"url+html": UrlHtmlFormatter,
}
var template_dir = "templates";

func main() {
	flag.Parse();
  apptemplate := LoadTemplate("application.html");
/*  addtemplate := LoadTemplate("add.html");*/
	
	http.Handle("/", http.HandlerFunc(func(c *http.Conn, req *http.Request) {
	  QR(c, req, apptemplate)
	}));
	http.Handle("/css/", http.FileServer("public/css", "/css/"));
	http.Handle("/js/", http.FileServer("public/js", "/js/"));
	
	err := http.ListenAndServe(*addr, nil);
	if err != nil {
		log.Exit("ListenAndServe:", err);
	}
}

func LoadTemplate(path string) *template.Template {
  log.Stdout(template_dir + "/" + path);
  data, err := io.ReadFile(template_dir + "/" + path);
  if err != nil {
		log.Exit("ReadFile:", err);
  }
  return template.MustParse(string(data), fmap);
}

func QR(c *http.Conn, req *http.Request, templ *template.Template) {
	templ.Execute(req.FormValue("s"), c);
}

func UrlHtmlFormatter(w io.Writer, v interface{}, fmt string) {
	template.HTMLEscape(w, strings.Bytes(http.URLEscape(v.(string))));
}


const templateStr = `
<html>
<head>
<meta http-equiv="Content-Type" content="text/html; charset=utf-8" /> 
<meta name="viewport" content="width=480" /> 
<title>QR Link Generator</title>
<link rel="stylesheet" href="/css/screen.css" type="text/css" media="screen, projection" />

<script src="http://ajax.googleapis.com/ajax/libs/jquery/1.3.2/jquery.min.js"></script>
<script src="/js/script.js"></script>
</head>
<body>
{.section @}
<img src="http://chart.apis.google.com/chart?chs=300x300&cht=qr&choe=UTF-8&chl={@|url+html}"
/>
<br>
{@|html}
<br>
<br>
{.end}
<form action="/" name=f method="GET"><input maxLength=1024 size=70
name=s value="" title="Text to QR Encode"><input type=submit
value="Show QR" name=qr>
</form>
</body>
</html>
`
