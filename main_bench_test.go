package helper

import (
	"encoding/json"
	"strconv"
	"strings"
	"testing"
)

func BenchmarkNumToAlphabet(t *testing.B) {
	for i := 0; i < t.N; i++ {
		_ = NumToAlphabet(int64(i))
	}
}

func BenchmarkJsonImplement(t *testing.B) {
	for i := 0; i < t.N; i++ {
		tt := &struct {
			A *UnstableInt
		}{}
		err := json.Unmarshal([]byte(`{"A":"`+strconv.Itoa(i)+`"}`), tt)
		if err != nil {
			t.Error(err)
		}
		_, err = json.Marshal(tt)
		if err != nil {
			t.Error(err)
		}
		// t.Log(string(rsp))
	}
}

// 构造一段有代表性的 HTML 文本，包含嵌套标签和文本
const testHTML = `
<html>
<head><title>Test Page</title></head>
<body>
<h1>Hello, World!</h1>
<p>This is a <strong>benchmark</strong> test for <em>HTML text extraction</em>.</p>
<ul>
<li>Item 1</li>
<li>Item 2</li>
<li>Item 3</li>
</ul>
<div><span>Nested <span>content</span> here</span></div>
</body>
</html>
`

// Benchmark for HtmlExtractText
func BenchmarkExtractHtmlText(b *testing.B) {
	// 预构建 Reader，避免在循环内部分配
	htmlStr := testHTML
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		r := strings.NewReader(htmlStr)
		_, err := ExtractHtmlText(r)
		if err != nil {
			b.Fatal(err)
		}
	}
}

// Benchmark for RegexExtractText
func BenchmarkExtractHtmlTextWithRegex(b *testing.B) {
	htmlStr := testHTML
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = ExtractHtmlTextWithRegex(htmlStr)
	}
}
