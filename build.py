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

from threading import Thread
import os
import sys
import time
import zipfile
import build_settings
from colorama import *
from rich.console import Console

version = build_settings.VERSION
package_prefix = build_settings.PACKAGE_PREFIX
packages_to_test = build_settings.TEST_PACKAGES
platforms = build_settings.PLATFORMS

console = Console()

started_threads_counter = 0  # amount of started build threads
finished_threads_counter = 0  # amount of finished build threads

error_occured = False
problematic_command = ""

print("⚙️  build started")

threads = []


def compile(package, platform):
    """ Compile go package `package`, for platform
    (os=platform[0], architecture=platform[1], binary_file_extension=platform[2])
    """
    global started_threads_counter, finished_threads_counter, error_occured, problematic_command

    started_threads_counter += 1
    now = time.time()
    os.putenv("GOOS", platform[0])
    os.putenv("GOARCH", platform[1])

    command = f"go build -ldflags \"-s -w\" -o .build/{platform[0]}.{platform[1]}/{package.split('/')[-1]}{platform[2]} {package}"
    os.system(command)

    finished_threads_counter += 1

    print(
        f"{Fore.GREEN}{Style.BRIGHT}[build]{Style.RESET_ALL}{Fore.GREEN}: finish {package.split('/')[-1]} [OS]: {platform[0]} [ARCH]: {platform[1]} in {time.time() - now} s{Fore.RESET}")


try:
    os.mkdir(".build")
except:
    pass

for directory in os.listdir("cmd"):
    """ Every `main` package, that will be compiled into binary is located in `cmd` folder, that package is compiling and binary file
    is located at .build/%OS%.%ARCHITECTURE% folder.
    """
    package = f"{package_prefix}/cmd/{directory}"

    for platform in platforms:
        try:
            os.mkdir(f".build/{platform[0]}.{platform[1]}")
        except:
            pass

        # start compile thread
        thread = Thread(target=compile, args=(package, platform))
        thread.start()
        print(
            f"{Fore.BLUE}{Style.BRIGHT}[build]{Style.RESET_ALL}{Fore.BLUE}: start building {directory} [OS]: {platform[0]} [ARCH]: {platform[1]}")

# while all compile and test threads are not finished, print out status in console
while started_threads_counter != finished_threads_counter:
    with console.status(f"[bold yellow]waiting for build threads") as st:
        while started_threads_counter != finished_threads_counter:
            pass

print("======== DONE BUILD AND TEST THREADS ========")
