package frame

import (
	"fmt"
	"html/template"
)

func GetFrameButton(index int, content string) template.HTML {
	button := fmt.Sprintf(`<meta property="fc:frame:button:%d" content="%s" />`, index, content)
	return template.HTML(button)
}

func GetFramePostButton(url string) template.HTML {
	post := fmt.Sprintf(`<meta property="fc:frame:post_url" content="%s" />`, url)
	return template.HTML(post)
}
