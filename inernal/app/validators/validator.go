package validators

import (
	"errors"
	"fmt"
	"github.com/go-playground/validator/v10"
	"reflect"
)

type Validator struct {
	validate *validator.Validate
}

/*
基础验证：直接调用 ValidateStruct 对结构体进行常规验证。
扩展验证：通过 RegisterCustomValidation 注册自定义验证规则，满足复杂的业务需求。
错误处理：通过 TranslateValidationErrors 将错误信息转换为可读的格式，提高可维护性。

在小型项目或初期阶段，使用一个全局通用的验证器对象是最简单且快速的方式。
可以直接在 main 方法中初始化一个全局验证器，并在需要时通过上下文或全局变量调用。
这种方式减少了依赖注入的复杂性，并且适用于项目的早期阶段或实验性质的项目。

随着项目规模的扩大，模块或服务的复杂性增加，不同模块可能需要不同的验证逻辑。
此时，采用依赖注入模式，将验证器与具体模块或服务绑定是更合理的选择。
每个模块可以根据自身的需求注入一个专用的验证器实例，以应对更加复杂的业务场景和验证逻辑。
*/

func NewValidator() *Validator {
	v := &Validator{
		validate: validator.New(),
	}

	// 注册自定义标签名解析器，使验证错误信息中显示的字段名更具可读性
	v.validate.RegisterTagNameFunc(func(fld reflect.StructField) string {
		name := fld.Tag.Get("json")
		if name == "" {
			name = fld.Name
		}
		return name
	})

	return v
}

// ValidateStruct 执行对结构体的验证
func (v *Validator) ValidateStruct(s interface{}) error {
	return v.validate.Struct(s)
}

// RegisterCustomValidation 注册自定义验证器
func (v *Validator) RegisterCustomValidation(tag string, fn validator.Func) error {
	return v.validate.RegisterValidation(tag, fn)
}

// TranslateValidationErrors 将验证错误信息转换为更具可读性的格式
func (v *Validator) TranslateValidationErrors(err error) map[string]string {
	var validationErrors validator.ValidationErrors
	errors.As(err, &validationErrors)
	errors := make(map[string]string)
	for _, fieldError := range validationErrors {
		errors[fieldError.Field()] = fmt.Sprintf("Validation failed on the '%s' tag", fieldError.Tag())
	}
	return errors
}
