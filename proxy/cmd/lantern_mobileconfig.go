package main

import (
	_ "github.com/lib/pq"

	"flag"
	"text/template"
	"io/ioutil"
	"encoding/base64"
	"bytes"
)

type Data struct {
	CertificateBase64  string
}


func main() {
	//parse flags
	caCertPath := flag.String("ca-path", "/srv/lantern/certs/ca.crt", "Path to CA file to embed in mobileconfig")
	mobileConfigPath := flag.String("mobileconfig-path", "/srv/lantern/lantern.mobileconfig", "Path to mobileconfig file to generate")
	flag.Parse()

	templateContent, err := ioutil.ReadFile("/defaults/lantern.mobileconfig.tmpl")
	if err != nil {
		panic(err)
	}

	t, err := template.New("mobileconfig").Parse(string(templateContent))
	if err != nil {
		panic(err)
	}

	caCertificateContent, err := ioutil.ReadFile(*caCertPath)
	if err != nil {
		panic(err)
	}

	tmplData := Data{
		CertificateBase64:  base64.StdEncoding.EncodeToString(caCertificateContent),
	}

	var tpl bytes.Buffer
	err = t.Execute(&tpl, tmplData)
	if err != nil {
		panic(err)
	}

	err = ioutil.WriteFile(*mobileConfigPath, tpl.Bytes(), 0644)
	if err != nil {
		panic(err)
	}
}