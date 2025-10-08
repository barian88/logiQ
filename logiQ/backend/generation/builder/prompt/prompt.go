package prompt

import (
	"backend/generation/config"
	"errors"
	"fmt"
	"math/rand/v2"
	"sort"
	"strings"

	"backend/models"
)

type Prompt struct {
	Text string
}

type Builder struct{}

func NewBuilder() Builder { return Builder{} }

// Params bundles the inputs needed to render a prompt from a template.
type Params struct {
	Intent config.IntentSpec
	QType  models.QuestionType
	Data   map[string]string
}

// BuildPrompt 根据意图、问题类型和数据生成问题文本。
//   - params: 包含意图名称、意图规格、问题类型以及模板替换所需的数据。
//   - rng:    随机数生成器，用于随机选择模板。
//
// 返回生成的Prompt，包含文本和使用的模板ID。
func (Builder) BuildPrompt(params Params, rng *rand.Rand) (Prompt, error) {
	if rng == nil {
		return Prompt{}, errors.New("generation: rng must not be nil")
	}

	templates := params.Intent.Templates.TemplatesFor(params.QType)
	// 确保有可用模板，应该不会发生，因为在samplePlan已经剔除了模板为空的intent
	if len(templates) == 0 {
		return Prompt{}, errors.New("generation: intent has no templates for requested question type")
	}

	// 随机抽取一个模板
	idx := 0
	if len(templates) > 1 {
		idx = rng.IntN(len(templates))
	}
	template := templates[idx]

	// 按键名排序，确保替换顺序一致（便于测试）
	keys := make([]string, 0, len(params.Data))
	for key := range params.Data {
		keys = append(keys, key)
	}
	sort.Strings(keys)

	text := template
	for _, key := range keys {
		placeholder := fmt.Sprintf("{%s}", key)
		text = strings.ReplaceAll(text, placeholder, params.Data[key])
	}

	return Prompt{Text: text}, nil
}
