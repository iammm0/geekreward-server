package utils

import (
	"errors"
	"github.com/google/uuid"
	"reflect"

	"github.com/go-playground/locales/zh"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	zhTranslations "github.com/go-playground/validator/v10/translations/zh"
)

// Validator 是一个封装了 go-playground/validator 的结构体，支持错误消息的本地化
type Validator struct {
	validate *validator.Validate
	trans    ut.Translator
}

// NewValidator 初始化一个新的 Validator 实例，并注册中文翻译
func NewValidator() (*Validator, error) {
	// 创建一个中文的 locale
	zhLocale := zh.New()

	// 创建一个 UniversalTranslator 实例
	uni := ut.New(zhLocale, zhLocale)

	// 获取中文的 Translator
	trans, found := uni.GetTranslator("zh")
	if !found {
		return nil, errors.New("translator not found")
	}

	// 创建一个新的验证器实例
	v := validator.New()

	// 注册自定义标签名解析器，使验证错误信息中显示的字段名更具可读性
	v.RegisterTagNameFunc(func(fld reflect.StructField) string {
		// 优先使用 json 标签
		name := fld.Tag.Get("json")
		if name == "-" || name == "" {
			// 如果没有 json 标签，则使用字段名
			name = fld.Name
		}
		return name
	})

	// 注册中文翻译
	if err := zhTranslations.RegisterDefaultTranslations(v, trans); err != nil {
		return nil, err
	}

	// 注册 UUID 验证
	err := v.RegisterValidation("uuid", func(fl validator.FieldLevel) bool {
		_, err := uuid.Parse(fl.Field().String())
		return err == nil
	})
	if err != nil {
		return nil, err
	}

	return &Validator{
		validate: v,
		trans:    trans,
	}, nil
}

// ValidateStruct 执行对结构体的验证
func (v *Validator) ValidateStruct(obj interface{}) error {
	return v.validate.Struct(obj)
}

// Engine 返回底层的验证器实例
func (v *Validator) Engine() interface{} {
	return v.validate
}

// RegisterCustomValidation 注册自定义验证器
func (v *Validator) RegisterCustomValidation(tag string, fn validator.Func) error {
	return v.validate.RegisterValidation(tag, fn)
}

// TranslateValidationErrors 将验证错误信息转换为更具可读性的格式（中文）
func (v *Validator) TranslateValidationErrors(err error) map[string]string {
	// 类型断言，确保错误是 validator.ValidationErrors 类型
	var validationErrors validator.ValidationErrors
	ok := errors.As(err, &validationErrors)
	if !ok {
		// 如果不是 ValidationErrors，返回一个通用错误
		return map[string]string{
			"error": "验证失败",
		}
	}

	// 创建一个 map 来存储字段名和对应的错误消息
	errorsMap := make(map[string]string)

	// 遍历所有的验证错误
	for _, fieldError := range validationErrors {
		// 获取翻译后的错误消息
		errorsMap[fieldError.Field()] = fieldError.Translate(v.trans)
	}

	return errorsMap
}
