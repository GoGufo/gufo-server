// Copyright 2020 Alexey Yanchenko <mail@yanchenko.me>
//
// This file is part of the Gufo library.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//    http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
//
//
// This file content curent app version and System DB API VERSION
// DB API Version need to compare with plugins DB Vesrions
// If DB version is same it means that plagin use right System DB structure
// System DB Structure descibes in functions/dbstructure.go

package version

// VERSION is current Gifo version
const VERSION = "1.4.9.#CI_PID-#CI_HASH (#CI_DATE)"

// VERSIONDB is current System DB API VERSION
const VERSIONDB = "1.0"
