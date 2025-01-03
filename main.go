package main

import (
	"fmt"
	"html/template"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/gomarkdown/markdown"
	"github.com/gomarkdown/markdown/parser"
)

// 定义模板数据
type TemplateData struct {
	Title    string
	SiteName string
	Articles []Article // 文章列表
}

// 文章结构体
type Article struct {
	Title   string        // 文章标题
	Link    string        // 文章链接（文件名）
	Content template.HTML // 文章内容（HTML 格式）
}

func main() {
	// 确保 dist 目录存在
	if err := os.MkdirAll("dist", os.ModePerm); err != nil {
		panic(err)
	}

	// 加载模板文件
	tmpl, err := template.ParseFiles("template/index.html")
	if err != nil {
		panic(err)
	}

	// 获取文章列表
	articles, err := getArticles("posts")
	if err != nil {
		panic(err)
	}

	// 准备模板数据
	data := TemplateData{
		Title:    "首页",
		SiteName: "AcoFork的Cloudflare全栈博客",
		Articles: articles,
	}

	// 生成 HTML 文件
	outputPath := filepath.Join("dist", "index.html")
	outputFile, err := os.Create(outputPath)
	if err != nil {
		panic(err)
	}
	defer outputFile.Close()

	// 执行模板
	if err := tmpl.Execute(outputFile, data); err != nil {
		panic(err)
	}

	fmt.Println("构建完成！")
}

// 获取文章列表
func getArticles(dir string) ([]Article, error) {
	files, err := ioutil.ReadDir(dir)
	if err != nil {
		return nil, err
	}

	var articles []Article
	for _, file := range files {
		if !file.IsDir() && strings.HasSuffix(file.Name(), ".md") {
			// 读取 Markdown 文件内容
			content, err := ioutil.ReadFile(filepath.Join(dir, file.Name()))
			if err != nil {
				return nil, err
			}

			// 将 Markdown 转换为 HTML
			htmlContent := mdToHTML(content)

			// 添加到文章列表
			articles = append(articles, Article{
				Title:   strings.TrimSuffix(file.Name(), ".md"), // 使用文件名作为标题
				Link:    strings.TrimSuffix(file.Name(), ".md"),
				Content: template.HTML(htmlContent),
			})
		}
	}
	return articles, nil
}

// 将 Markdown 转换为 HTML
func mdToHTML(md []byte) string {
	extensions := parser.CommonExtensions | parser.AutoHeadingIDs
	parser := parser.NewWithExtensions(extensions)
	return string(markdown.ToHTML(md, parser, nil))
}