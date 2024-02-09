/*
 * Copyright 2024 CloudWeGo Authors
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package template

import (
	"bytes"
	"fmt"
	"text/template"
)

var baseTemplate = `// Code generated by cwgo ({{.Version}}). DO NOT EDIT.

package {{.PackageName}}

import (
{{.GetImports}}
)` + "\n"

type BaseRender struct {
	Version     string            // cwgo version
	PackageName string            // package name in target generation go file
	Imports     map[string]string // key:import path value:import name
}

func (bt *BaseRender) RenderObj(buffer *bytes.Buffer) error {
	if err := templateRender(buffer, "baseTemplate", baseTemplate, bt); err != nil {
		return err
	}
	return nil
}

func (bt *BaseRender) GetImports() string {
	result := ""

	for key, value := range bt.Imports {
		if value != "" {
			result += fmt.Sprintf("\tvalue "+`"%s"`+"\n", key)
		} else {
			result += fmt.Sprintf("\t"+`"%s"`+"\n", key)
		}
	}

	return result
}

func templateRender(buffer *bytes.Buffer, templateName, parseText string, data any) error {
	tmpl, err := template.New(templateName).Parse(parseText)
	if err != nil {
		return err
	}

	if err = tmpl.Execute(buffer, data); err != nil {
		return err
	}

	return nil
}
