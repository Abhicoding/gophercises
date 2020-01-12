package main

import (
	"fmt"
	"os"

	"golang.org/x/net/html"
)

type Link struct {
	Href string
	Text string
}

var aTags []Link

func main() {

	args := os.Args
	if len(args) < 2 {
		fmt.Println("Please provide HTML file to parse")
		return
	}

	file, err := os.Open(args[1])
	if err != nil {
		fmt.Println(err)
		return
	}
	if file == nil {
		fmt.Println("Couldn't open the file provided. Exiting ...")
		return
	}

	z, err := html.Parse(file)
	if err != nil {
		fmt.Printf("Error parsing %s", file.Name())
		return
	}
	arrayOfNode := []*html.Node{z}
	breadthFirst(arrayOfNode)
	printLinks(aTags)
}

func printLinks(links []Link) {
	for _, anchor := range links {
		fmt.Printf("href=\"%s\" text=\"%s\"\n", anchor.Href, anchor.Text)
	}
}

func isAnchorTag(n *html.Node) *html.Node {
	if n.DataAtom.String() == "a" {
		return n
	}
	return nil
}

func getAnchor(a *html.Node) Link {
	var link Link
	for _, attr := range a.Attr {
		if attr.Key == "href" {
			link.Href = attr.Val
		}
	}
	link.Text = a.FirstChild.Data
	return link
}

func breadthFirst(n []*html.Node) {
	if len(n) == 0 {
		return
	}
	childNodes := []*html.Node{}
	for _, node := range n {
		if isAnchorTag(node) != nil {
			aTags = append(aTags, getAnchor(node))
		}
		if node.FirstChild != nil {
			f := node.FirstChild
			childNodes = append(childNodes, f)
			s := f.NextSibling
			for {
				if s != nil {
					childNodes = append(childNodes, s)
					s = s.NextSibling
					continue
				}
				break
			}
		}
	}
	breadthFirst(childNodes)
}
