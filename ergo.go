package main

import (
	"bytes";
	"flag";
	"http";
	"io";
	"log";
	"strings";
	"template";
	pathutil "path";
)

var addr = flag.String("addr", ":1718", "http service address") // Q=17, R=18

var fmap = template.FormatterMap{
	"html": template.HTMLFormatter,
	"url+html": UrlHtmlFormatter,
}
var tmplroot = "templates";

func readTemplate(name string) *template.Template {
	path := pathutil.Join(tmplroot, name);
	data, err := io.ReadFile(path);
	if err != nil {
		log.Exitf("ReadFile %s: %v", path, err)
	}
	t, err := template.Parse(string(data), fmap);
	if err != nil {
		log.Exitf("%s: %v", name, err)
	}
	return t;
}

var (
	  appHTML,
		addHTML *template.Template;
)

func readTemplates() {
	// have to delay until after flags processing,
	// so that main has chdir'ed to goroot.
	appHTML = readTemplate("application.html");
	addHTML = readTemplate("add.html");
}

/*func readTemplate(path string) *template.Template {
  log.Stdout(template_dir + "/" + path);
  data, err := io.ReadFile(template_dir + "/" + path);
  if err != nil {
    log.Exit("ReadFile:", err);
  }
  return template.MustParse(string(data), fmap);
}*/

func servePage(c *http.Conn, title, query string, content string) {
	type Data struct {
		Title		string;
		//Timestamp	uint64;	// int64 to be compatible with os.Dir.Mtime_ns
		Query		string;
		Content		string;
	}

	//_, ts := fsTree.get();
	d := Data{
		Title: title,
		//Timestamp: uint64(ts) * 1e9,	// timestamp in ns
		Query: query,
		Content: content,
	};

	if err := appHTML.Execute(&d, c); err != nil {
		log.Stderrf("appHTML.Execute: %s", err)
	}
}

func main() {
	flag.Parse();
  readTemplates();
	
	http.Handle("/", http.HandlerFunc(func(c *http.Conn, req *http.Request) {
	  
	}));
	http.Handle("/add", http.HandlerFunc(func(c *http.Conn, req *http.Request) {
	  // Process the add template
  	var buf bytes.Buffer;
/*    if err := parseerrorHTML.Execute(errors, &buf); err != nil {*/
	  err := addHTML.Execute("x", &buf);
	  if err != nil {
  		log.Stderrf("addHTML.Execute: %s", err)
	  }
  	//templ.Execute(req.FormValue("s"), c);
  	servePage(c, "Add", "", string(buf.Bytes()));
	}));
	http.Handle("/css/", http.FileServer("public/css", "/css/"));
	http.Handle("/js/", http.FileServer("public/js", "/js/"));
	
	err := http.ListenAndServe(*addr, nil);
	if err != nil {
		log.Exit("ListenAndServe:", err);
	}
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
