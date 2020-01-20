package main

import (
	"encoding/xml"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"os"

	"golang.org/x/net/html"
)

type UrlSet struct {
	XMLName xml.Name `xml:"urlset"`
	Url     []*Url   `xml:"url"`
}
type Url struct {
	Loc string `xml:"loc"`
}

func main() {
	var siteMap UrlSet
	args := os.Args
	URL := args[len(args)-1]
	u, _ := url.Parse(URL)
	resp, err := http.Get(u.String())
	if err != nil {
		fmt.Println(err)
		return
	}
	z, err := html.Parse(resp.Body)

	if err != nil {
		fmt.Println(err)
		return
	}
	html, _ := FindHTML(z)
	anchorNodes := Crawler(html)
	filteredTags := FilterSiteLinks(anchorNodes, u)
	siteMap = getSiteMap(filteredTags)

	printXML(siteMap)
	defer resp.Body.Close()
}

func FindHTML(doc *html.Node) (*html.Node, error) {
	if doc.DataAtom.String() == "html" {
		return doc, nil
	}
	if doc == nil {
		return nil, errors.New("No HTML tag found")
	}
	for elem := doc.FirstChild; elem != nil; elem = elem.NextSibling {
		newElem, err := FindHTML(elem)
		if err != nil {
			continue
		}
		return newElem, nil
	}
	return nil, errors.New("No HTML tag found")
}

func GetChildNodes(n *html.Node) []*html.Node {
	var childNodes []*html.Node
	for childNode := n.FirstChild; childNode != nil; childNode = childNode.NextSibling {
		childNodes = append(childNodes, childNode)
	}
	if len(childNodes) == 0 {
		return nil
	}
	return childNodes
}

func FilterAnchorNodes(nodes []*html.Node) []*html.Node {
	var anchorNodes []*html.Node
	for _, node := range nodes {
		if node.DataAtom.String() == "a" {
			anchorNodes = append(anchorNodes, node)
		}
	}
	if len(anchorNodes) == 0 {
		return nil
	}
	return anchorNodes
}

func printXML(siteMap UrlSet) {
	enc := xml.NewEncoder(os.Stdout)
	enc.Indent("  ", "    ")
	if err := enc.Encode(siteMap); err != nil {
		fmt.Printf("error: %v\n", err)
	}
}

func printAnchorNodes(childTags []*html.Node) {
	for _, tag := range childTags {
		for _, attr := range tag.Attr {
			if attr.Key == "href" {
				fmt.Println(tag.DataAtom.String(), attr.Val)
				break
			}
		}
	}
}

func Crawler(n *html.Node) []*html.Node {
	var childNodes, anchorNodes []*html.Node
	childNodes = GetChildNodes(n)
	if childNodes == nil {
		return nil
	}
	childNodes = Traverse([]*html.Node{n}, childNodes)
	anchorNodes = FilterAnchorNodes(childNodes)
	return anchorNodes
}

func Traverse(Nodes, childNodes []*html.Node) []*html.Node {
	var nextChildNodes []*html.Node
	if childNodes != nil {
		for _, node := range childNodes {
			nextChildNodes = append(nextChildNodes, GetChildNodes(node)...)
		}
	}
	if len(nextChildNodes) == 0 {
		return append(Nodes, childNodes...)
	}
	return Traverse(append(Nodes, childNodes...), nextChildNodes)
}

func FilterSiteLinks(list []*html.Node, URL *url.URL) map[string]bool {
	temp := make(map[string]bool)
	for _, tag := range list {
		for _, attr := range tag.Attr {
			if attr.Key == "href" {
				u, _ := url.Parse(attr.Val)
				if u.Host == "" && u.Fragment != "" {
					continue
				}
				if u.Scheme == "" {
					u.Scheme = URL.Scheme
				}
				if u.Scheme != URL.Scheme {
					continue
				}
				if u.Host == "" {
					u.Host = URL.Hostname()
				}
				if u.Hostname() == URL.Hostname() {
					if _, ok := temp[u.String()]; !ok {
						temp[u.String()] = true
						continue
					}
					continue
				}

				continue
			}
			continue
		}
	}
	return temp
}

func getSiteMap(links map[string]bool) UrlSet {
	var siteMap UrlSet
	for k, _ := range links {
		uu := Url{
			Loc: k,
		}
		siteMap.Url = append(siteMap.Url, &uu)
	}
	return siteMap
}
