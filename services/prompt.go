package services

import (
	"fmt"
	"strings"
)

// PromptOptions 用于存储生成 prompt 的可选参数。
// 每个字段代表 MidJourney 中可配置的选项。
type PromptOptions struct {
	// Prompt 核心提示词，描述图像内容的自然语言文本。
	// 例子: "A futuristic city with flying cars"
	Prompt string

	// Style 风格选择，指定图像生成的风格。
	// 可选值：
	// - "original"：原版风格。
	// - "4a"：第四代风格 a 版本。
	// - "4b"：第四代风格 b 版本。
	// - "4c"：第四代风格 c 版本。
	// - "niji"：适用于二次元风格。
	Style string

	// Version 版本选择，指定 MidJourney 使用的模型版本。
	// 可选值：
	// - "v1"：第 1 版模型。
	// - "v2"：第 2 版模型。
	// - "v3"：第 3 版模型。
	// - "v4"：第 4 版模型。
	// - "v5"：第 5 版模型。
	// - "v6"：第 6 版模型（如有发布）。
	Version string

	// Niji 是否使用 Niji 模式，专为二次元风格设计的绘画模式。
	// true：启用 Niji 模式，适合绘制动画风格的图像。
	// false：禁用 Niji 模式，使用常规风格生成图像。
	Niji bool

	// Aspect 宽高比，定义生成图像的宽高比。
	// 可选值：
	// - "1:1"：正方形图片。
	// - "16:9"：宽屏图片（常用于电影场景）。
	// - "9:16"：竖屏图片（适用于移动设备的展示）。
	// - "4:3"：常见的相机宽高比。
	// - "3:4"：相对较长的竖屏比例。
	Aspect string

	// Quality 图片质量，定义生成图像的质量。质量越高，生成时间越长。
	// 默认值为 1。
	// 可选值：
	// - 0.25：低质量，适合快速生成粗略图像。
	// - 0.5：中等质量，速度较快。
	// - 1：标准质量，适合大部分场景。
	// - 2：高质量，生成时间较长但图像更细致。
	Quality float64

	// Speed 绘画速度，控制生成图像的速度模式。
	// 可选值：
	// - "fast"：快速模式，减少生成时间，但可能降低图像质量。
	// - "relax"：慢速模式，提高生成质量，但会增加生成时间。
	Speed string

	// Seed 随机种子，用于控制图像生成的随机性。相同的种子值会生成相似的图像。
	// 0 表示不指定种子，生成全新随机图像。
	// 例子：Seed = 123456 会基于该种子值生成图像，便于复现相似风格。
	Seed int

	// Model 自定义绘画模型，允许选择不同的模型进行绘制。
	// 可选值：
	// - "openjourney"：开源模型，生成具有实验性风格的图像。
	// - "default"：默认模型，使用 MidJourney 标准生成模式。
	// - 其他模型根据 MidJourney 支持的模型进行扩展。
	Model string

	// JobID 唯一任务 ID，用于标识每个生成请求的唯一性。
	// 用于在生成和回调时区分不同的任务。通常会通过系统生成的唯一值赋予。
	//JobID string
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

//// SetJobID 设置唯一的任务 ID (job-id)
//func (pg *PromptGenerator) SetJobID(jobID string) {
//	pg.options.JobID = jobID
//}

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

	// 最后加入 job-id
	//if pg.options.JobID != "" {
	//	promptParts = append(promptParts, fmt.Sprintf("--job-id %s", pg.options.JobID))
	//}

	// 返回最终拼接的 prompt
	return strings.Join(promptParts, " ")
}

// ParseJobIDFromContent 解析 content 字段中的 job-id
//func ParseJobIDFromContent(content string) string {
//	// 使用正则表达式解析 job-id
//	re := regexp.MustCompile(`--job-id (\d+)`)
//	matches := re.FindStringSubmatch(content)
//
//	if len(matches) > 1 {
//		return matches[1] // 返回 job-id
//	}
//	return ""
//}

func promptTest() {
	// 生成唯一的 nonce 这个应该本地业务时候就生成在外部传入
	//nonce := "145491650871820289"

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
	//generator.SetJobID(nonce)         // 设置 job-id

	// 生成最终的 prompt
	finalPrompt := generator.Generate()
	fmt.Println("Generated Prompt:", finalPrompt)

	//// 模拟解析 content 中的 job-id
	//parsedJobID := ParseJobIDFromContent(finalPrompt)
	//fmt.Println("Parsed JobID:", parsedJobID)
}
