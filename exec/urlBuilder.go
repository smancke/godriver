package exec

import (
        "text/template"
        "fmt"
        "bytes"
)

type UrlStruct struct {
        SearchBaseUrl string
        Term string
        Limit string
        Page string
        Facet1 string
        Facet2 string
        Facet3 string
}

func BuildUrl(urlStruct UrlStruct) string {

        var buffer bytes.Buffer

        t := template.New("urlTemplate")

        const url =
                "{{.SearchBaseUrl}}"+
                        "{{if .Term}}"+ "?term=" +"{{.Term}}" +"{{- end}}"+
                        "{{if .Facet1}}" + "&facets[0]" + "{{.Facet1}}" + "{{- end}}"+
                        "{{if .Facet2}}" + "&facets[1]" + "{{.Facet2}}" + "{{- end}}"+
                        "{{if .Facet3}}" + "&facets[2]" + "{{.Facet3}}" + "{{- end}}"+
                        "{{if .Limit}}" + "&limit=" + "{{.Limit}}" + "{{- end}}"+
                        "{{if .Page}}" + "&page=" + "{{.Page}}" + "{{- end}}"

        t,_ = t.Parse(url)

        err := t.Execute(&buffer, urlStruct)

        if err != nil {
                fmt.Println(err.Error())
        }
        return buffer.String()
}