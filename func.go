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
// date: 2023-04-18

package log

import (
	"context"
	"github.com/go-wares/log/base"
	"github.com/go-wares/log/config"
	"github.com/go-wares/log/managers"
)

// +---------------------------------------------------------------------------+
// | Logger methods                                                            |
// +---------------------------------------------------------------------------+

func Debug(text string) {
	if config.Config.DebugOn() {
		managers.Manager.Log(nil, nil, base.Debug, text)
	}
}

func Info(text string) {
	if config.Config.InfoOn() {
		managers.Manager.Log(nil, nil, base.Info, text)
	}
}

func Warn(text string) {
	if config.Config.WarnOn() {
		managers.Manager.Log(nil, nil, base.Warn, text)
	}
}

func Error(text string) {
	if config.Config.ErrorOn() {
		managers.Manager.Log(nil, nil, base.Error, text)
	}
}

func Fatal(text string) {
	if config.Config.FatalOn() {
		managers.Manager.Log(nil, nil, base.Fatal, text)
	}
}

// +---------------------------------------------------------------------------+
// | Logger methods with formatter                                             |
// +---------------------------------------------------------------------------+

func Debugf(format string, args ...interface{}) {
	if config.Config.DebugOn() {
		managers.Manager.Log(nil, nil, base.Debug, format, args...)
	}
}

func Infof(format string, args ...interface{}) {
	if config.Config.InfoOn() {
		managers.Manager.Log(nil, nil, base.Info, format, args...)
	}
}

func Warnf(format string, args ...interface{}) {
	if config.Config.WarnOn() {
		managers.Manager.Log(nil, nil, base.Warn, format, args...)
	}
}

func Errorf(format string, args ...interface{}) {
	if config.Config.ErrorOn() {
		managers.Manager.Log(nil, nil, base.Error, format, args...)
	}
}

func Fatalf(format string, args ...interface{}) {
	if config.Config.FatalOn() {
		managers.Manager.Log(nil, nil, base.Fatal, format, args...)
	}
}

// +---------------------------------------------------------------------------+
// | Logger methods with formatter and context                                 |
// +---------------------------------------------------------------------------+

func Debugfc(ctx context.Context, format string, args ...interface{}) {
	if config.Config.DebugOn() {
		managers.Manager.Log(ctx, nil, base.Debug, format, args...)
	}
}

func Infofc(ctx context.Context, format string, args ...interface{}) {
	if config.Config.InfoOn() {
		managers.Manager.Log(ctx, nil, base.Info, format, args...)
	}
}

func Warnfc(ctx context.Context, format string, args ...interface{}) {
	if config.Config.WarnOn() {
		managers.Manager.Log(ctx, nil, base.Warn, format, args...)
	}
}

func Errorfc(ctx context.Context, format string, args ...interface{}) {
	if config.Config.ErrorOn() {
		managers.Manager.Log(ctx, nil, base.Error, format, args...)
	}
}

func Fatalfc(ctx context.Context, format string, args ...interface{}) {
	if config.Config.FatalOn() {
		managers.Manager.Log(ctx, nil, base.Fatal, format, args...)
	}
}
