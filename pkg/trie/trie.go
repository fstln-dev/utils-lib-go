package trie

import (
	"errors"
	"regexp"
	"strings"
)

var Methods = []string{
	"GET", "POST", "PUT", "DELETE", "PATCH",
}

type PathTrie interface {
	Add(path string, method string, data interface{}) error
	Find(path string, method string) interface{}
}

type trie struct {
	regexp string
	trees  map[string]path_tree
}

type path_tree struct {
	data   interface{}           // 节点内容
	pare   string                // 节点 key
	childs map[string]*path_tree // 前缀树数据
}

func Init() PathTrie {

	trees := make(map[string]path_tree)
	for _, v := range Methods {
		trees[v] = path_tree{
			childs: make(map[string]*path_tree),
		}
	}

	return &trie{
		regexp: `^\{[a-zA-Z0-9_-]+\}$`,
		trees:  trees,
	}
}

func (t *trie) Add(path string, method string, data interface{}) error {

	paths := t.split_path(path)
	if len(paths) < 1 {
		return errors.New("path is not nil")
	}

	// 过滤 http 方法
	tree, ok := t.trees[method]
	if !ok {
		return errors.New("methods must GET, POST, PUT, DELETE or PATCH")
	}

	handle := &tree
	for _, v := range paths {

		key := v
		if is_wild, err := t.regrex(v); err != nil {
			return err
		} else if is_wild {
			key = "*"
		}

		if _, ok := handle.childs[key]; !ok {
			handle.childs[key] = &path_tree{
				pare:   v,
				childs: make(map[string]*path_tree),
			}
		}
		// 下一层级
		handle = handle.childs[key]
	}

	// 最后一个路径上标记为叶节点
	handle.data = data
	return nil
}

func (t *trie) Find(path string, method string) interface{} {

	paths := t.split_path(path)
	if len(paths) < 1 {
		return nil
	}

	// 过滤 http 方法
	tree, ok := t.trees[method]
	if !ok {
		// return nil
		return errors.New("methods must GET, POST, PUT, DELETE or PATCH")
	}

	handle := &tree
	for _, v := range paths {
		if next, ok := handle.childs[v]; ok {
			// 下一层级
			handle = next
			continue
		} else if next, ok := handle.childs["*"]; ok {
			handle = next
			continue
		}
		return nil
	}

	// 取最后一栏的数据
	return handle.data
}

func (t *trie) split_path(path string) []string {
	if len(path) < 1 {
		return make([]string, 0)
	}
	if path[0] == '/' {
		path = path[1:]
	}
	if x := len(path); path[x-1] == '/' {
		path = path[:x-1]
	}
	return strings.Split(path, "/")
}

func (t *trie) regrex(key string) (bool, error) {
	return regexp.Match(t.regexp, []byte(key))
}
