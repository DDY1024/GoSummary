package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/facebookgo/inject"
)

// 依赖注入主要解决的问题：简化初始化代码的管理，初始化复杂的依赖关系
// 1. 依赖反射实现的运行时依赖注入
// 2. 使用 code-gen 实现的依赖注入

type RootNode struct {
	Root *HomePlanetRenderApp `inject:""`
}

func (root *RootNode) Do(id uint64) string {
	return root.Root.Render(id)
}

type HomePlanetRenderApp struct {
	NameAPI   *NameAPI   `inject:""`
	PlanetAPI *PlanetAPI `inject:""`
}

func (a *HomePlanetRenderApp) Render(id uint64) string {
	return fmt.Sprintf("%s is from the planet %s.", a.NameAPI.Name(id), a.PlanetAPI.Planet(id))
}

type NameAPI struct {
	HTTPTransport http.RoundTripper `inject:""`
}

func (n *NameAPI) Name(id uint64) string {
	return "Spock"
}

type PlanetAPI struct {
	HTTPTransport http.RoundTripper `inject:""`
}

func (p *PlanetAPI) Planet(id uint64) string {
	return "Vulcan"
}

func main() {
	var g inject.Graph
	var root RootNode

	// root 为需要被注入的根对象
	// http.DefaultTransport 实现了 http.RoundTripper 接口，作为 NameAPI、PlanetAPI 的成员，同样需要被注入
	// 自底向上构造 root 对象，解决依赖对象的构建工作
	err := g.Provide(&inject.Object{Value: &root}, &inject.Object{Value: http.DefaultTransport})
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	if err := g.Populate(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	fmt.Println(root.Do(42))
}
