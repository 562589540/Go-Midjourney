package services

import (
	"fmt"
	"strings"
)

// PromptOptions 用于存储生成 prompt 的可选参数
type PromptOptions struct {
	Prompt  string  // 核心 prompt
	Style   string  // 风格选择（如 original、4a）
	Version string  // 版本选择（如 v5）
	Niji    bool    // 是否使用 Niji 模式
	Aspect  string  // 宽高比（如 16:9）
	Quality float64 // 图片质量（默认 1，范围 0.25 - 2）
	Speed   string  // 绘画速度（fast 或 relax）
	Seed    int     // 随机种子
	Model   string  // 自定义绘画模型（如 openjourney、其他）
}

// PromptGenerator 生成最终的 MJ prompt
type PromptGenerator struct {
	options PromptOptions
}

// NewPromptGenerator 初始化 PromptGenerator
func NewPromptGenerator(prompt string) *PromptGenerator {
	// 默认使用给定的 prompt 初始化
	return &PromptGenerator{
		options: PromptOptions{
			Prompt:  prompt,
			Quality: 1, // 默认质量为 1
		},
	}
}

// SetStyle 设置风格（可选）
func (pg *PromptGenerator) SetStyle(style string) {
	pg.options.Style = style
}

// SetVersion 设置 MJ 版本（如 v5）
func (pg *PromptGenerator) SetVersion(version string) {
	pg.options.Version = version
}

// SetNiji 设置是否使用 Niji 模式
func (pg *PromptGenerator) SetNiji(useNiji bool) {
	pg.options.Niji = useNiji
}

// SetAspectRatio 设置宽高比
func (pg *PromptGenerator) SetAspectRatio(aspect string) {
	pg.options.Aspect = aspect
}

// SetQuality 设置图像质量（范围 0.25 到 2）
func (pg *PromptGenerator) SetQuality(quality float64) {
	if quality >= 0.25 && quality <= 2 {
		pg.options.Quality = quality
	}
}

// SetSpeed 设置绘画速度（fast 或 relax）
func (pg *PromptGenerator) SetSpeed(speed string) {
	if speed == "fast" || speed == "relax" {
		pg.options.Speed = speed
	}
}

// SetSeed 设置随机种子（可选）
func (pg *PromptGenerator) SetSeed(seed int) {
	pg.options.Seed = seed
}

// SetModel 设置自定义绘画模型
func (pg *PromptGenerator) SetModel(model string) {
	pg.options.Model = model
}

// Generate 生成最终的 MJ prompt
func (pg *PromptGenerator) Generate() string {
	var promptParts []string
	promptParts = append(promptParts, pg.options.Prompt)

	// 加入风格
	if pg.options.Style != "" {
		promptParts = append(promptParts, fmt.Sprintf("--style %s", pg.options.Style))
	}

	// 加入版本
	if pg.options.Version != "" {
		promptParts = append(promptParts, fmt.Sprintf("--v %s", pg.options.Version))
	}

	// 是否使用 Niji 模式
	if pg.options.Niji {
		promptParts = append(promptParts, "--niji")
	}

	// 加入宽高比
	if pg.options.Aspect != "" {
		promptParts = append(promptParts, fmt.Sprintf("--ar %s", pg.options.Aspect))
	}

	// 加入图片质量
	if pg.options.Quality != 1 {
		promptParts = append(promptParts, fmt.Sprintf("--q %.2f", pg.options.Quality))
	}

	// 加入绘画速度
	if pg.options.Speed != "" {
		promptParts = append(promptParts, fmt.Sprintf("--%s", pg.options.Speed))
	}

	// 加入随机种子
	if pg.options.Seed != 0 {
		promptParts = append(promptParts, fmt.Sprintf("--seed %d", pg.options.Seed))
	}

	// 加入自定义模型
	if pg.options.Model != "" {
		promptParts = append(promptParts, fmt.Sprintf("--model %s", pg.options.Model))
	}

	// 返回最终拼接的 prompt
	return strings.Join(promptParts, " ")
}

func promptTest() {

	//// 生成唯一的 nonce 这个应该本地业务时候就生成在外部传入
	//nonce, _ := GetNextIDGenerator().Generate()
	//
	//// 在 prompt 中附加 nonce，确保它与结果消息相关联
	//fullPrompt := fmt.Sprintf("[%d] %s", nonce, prompt)

	//储存在字典
	//saveRequest(fmt.Sprintf("%d", nonce), prompt)

	// 创建新的 prompt 生成器
	generator := NewPromptGenerator("A futuristic city with flying cars")

	// 设置一些可选参数
	generator.SetStyle("4a")
	generator.SetVersion("5")
	generator.SetNiji(true)
	generator.SetAspectRatio("16:9")
	generator.SetQuality(1.5)
	generator.SetSpeed("fast")
	generator.SetSeed(123456)
	generator.SetModel("openjourney") // 自定义模型

	// 生成最终的 prompt
	finalPrompt := generator.Generate()
	fmt.Println("Generated Prompt:", finalPrompt)
}
