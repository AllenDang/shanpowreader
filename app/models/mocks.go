// Copyright 2013 The Shanpow Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// mocks包含所有不需要存储到数据的实体定义

package models

type AjaxResult struct {
  Result   bool
  ErrorMsg string
  Data     interface{}
}
