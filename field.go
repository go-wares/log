// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
//
// author: wsfuyibing <websearch@163.com>
// date: 2023-05-13

package log

import (
	"context"
	"github.com/go-wares/log/base"
	"github.com/go-wares/log/config"
	"github.com/go-wares/log/managers"
)

type (
	// Field
	// 自定义字段.
	//
	//   log.Field{"key": "value"}.Info("info")
	Field map[string]interface{}
)

// +---------------------------------------------------------------------------+
// | Logger methods                                                            |
// +---------------------------------------------------------------------------+

func (o Field) Debug(text string) {
	if config.Config.DebugOn() {
		managers.Manager.Log(nil, o, base.Debug, text)
	}
}

func (o Field) Info(text string) {
	if config.Config.InfoOn() {
		managers.Manager.Log(nil, o, base.Info, text)
	}
}

func (o Field) Warn(text string) {
	if config.Config.WarnOn() {
		managers.Manager.Log(nil, o, base.Warn, text)
	}
}

func (o Field) Error(text string) {
	if config.Config.ErrorOn() {
		managers.Manager.Log(nil, o, base.Error, text)
	}
}

func (o Field) Fatal(text string) {
	if config.Config.FatalOn() {
		managers.Manager.Log(nil, o, base.Fatal, text)
	}
}

// +---------------------------------------------------------------------------+
// | Logger methods with formatter                                             |
// +---------------------------------------------------------------------------+

func (o Field) Debugf(format string, args ...interface{}) {
	if config.Config.DebugOn() {
		managers.Manager.Log(nil, o, base.Debug, format, args...)
	}
}

func (o Field) Infof(format string, args ...interface{}) {
	if config.Config.InfoOn() {
		managers.Manager.Log(nil, o, base.Info, format, args...)
	}
}

func (o Field) Warnf(format string, args ...interface{}) {
	if config.Config.WarnOn() {
		managers.Manager.Log(nil, o, base.Warn, format, args...)
	}
}

func (o Field) Errorf(format string, args ...interface{}) {
	if config.Config.ErrorOn() {
		managers.Manager.Log(nil, o, base.Error, format, args...)
	}
}

func (o Field) Fatalf(format string, args ...interface{}) {
	if config.Config.FatalOn() {
		managers.Manager.Log(nil, o, base.Fatal, format, args...)
	}
}

// +---------------------------------------------------------------------------+
// | Logger methods with formatter and context                                 |
// +---------------------------------------------------------------------------+

func (o Field) Debugfc(ctx context.Context, format string, args ...interface{}) {
	if config.Config.DebugOn() {
		managers.Manager.Log(ctx, o, base.Debug, format, args...)
	}
}

func (o Field) Infofc(ctx context.Context, format string, args ...interface{}) {
	if config.Config.InfoOn() {
		managers.Manager.Log(ctx, o, base.Info, format, args...)
	}
}

func (o Field) Warnfc(ctx context.Context, format string, args ...interface{}) {
	if config.Config.WarnOn() {
		managers.Manager.Log(ctx, o, base.Warn, format, args...)
	}
}

func (o Field) Errorfc(ctx context.Context, format string, args ...interface{}) {
	if config.Config.ErrorOn() {
		managers.Manager.Log(ctx, o, base.Error, format, args...)
	}
}

func (o Field) Fatalfc(ctx context.Context, format string, args ...interface{}) {
	if config.Config.FatalOn() {
		managers.Manager.Log(ctx, o, base.Fatal, format, args...)
	}
}
