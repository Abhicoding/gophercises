package main

import (
	"encoding/xml"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
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
	args := os.Args
	url := args[len(args)-1]
	resp, err := http.Get(url)
	if err != nil {
		fmt.Println(err)
		return
	}
	z, err := html.Parse(resp.Body)
	body, _ := ioutil.ReadAll(resp.Body)
	fmt.Println(string(body))

	if err != nil {
		fmt.Println(err)
		return
	}
	html, _ := FindHTML(z)
	anchorNodes := Crawler(html)
	printAnchorNodes(anchorNodes)
	//printChild(anchorTags)
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

func printChild(childTags []*html.Node) {
	for _, tag := range childTags {
		fmt.Println(tag.DataAtom.String())
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
	printChild(anchorNodes)
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
