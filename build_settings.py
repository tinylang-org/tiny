# Copyright (c) 2022 Salimgereyev Adi
# All rights reserved.

# Permission is hereby granted, free of charge, to any person obtaining a copy
# of this software and associated documentation files (the "Software"), to deal
# in the Software without restriction, including without limitation the rights
# to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
# copies of the Software, and to permit persons to whom the Software is
# furnished to do so, subject to the following conditions:

# The above copyright notice and this permission notice shall be included in all
# copies or substantial portions of the Software.

# THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
# IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
# FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
# AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
# LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
# OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
# SOFTWARE.

PLATFORMS = [
    # windows
    ('windows', 'amd64', '.exe'),
    ('windows', '386', '.exe'),

    # linux
    ('linux', '386', ''),
    ('linux', 'amd64', ''),

    # linux arm
    ('linux', 'arm', ''),
    ('linux', 'arm64', ''),
]

PACKAGE_PREFIX = "github.com/tinylang-org/tiny"

# packages to test
TEST_PACKAGES = [
    "pkg/lexer",
    "pkg/ast",
    "pkg/parser",
]

VERSION = "alpha_0.1.0"
