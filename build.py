#!/usr/bin/env python3
"""
Python build script for LLVM 21.x + Go project

Usage:
    python build.py [build|run|all]

- LLVM source: c_pkg/llvm-project
- Build dir: c_pkg/llvm-project/build
- Go project: current repo
- Works on Windows/Linux/macOS
"""

from numpy import true_divide
import os
import subprocess
from pathlib import Path
import platform
import sys
from shutil import which

# ===== CONFIG =====
LLVM_SRC_DIR = Path("c_pkg/llvm-project")       # may already exist
BUILD_DIR = LLVM_SRC_DIR / "build"             # LLVM build output
GO_PROJECT_DIR = Path(".")                     # current repo
LLVM_VERSION = "release/21.x"                 # branch/tag
LLVM_TARGETS = "X86"                           # minimal target

# ===== UTILS =====
def run(cmd, cwd=None, env=None):
    print(f"> {' '.join(cmd)}")
    subprocess.check_call(cmd, cwd=cwd, env=env)

def detect_os():
    os_name = platform.system()
    if os_name == "Darwin":
        return "mac"
    elif os_name == "Windows":
        return "windows"
    else:
        return "linux"

def llvm_config_path():
    exe = "llvm-config.exe" if detect_os() == "windows" else "llvm-config"
    return str(BUILD_DIR / "bin" / exe)

# ===== Clone / Update LLVM =====
def ensure_llvm():
    if not LLVM_SRC_DIR.exists():
        print(f"Cloning LLVM {LLVM_VERSION} to {LLVM_SRC_DIR} ...")
        run([
            "git", "clone",
            "--branch", LLVM_VERSION,
            "--config core.autocrlf=false",# recomended in LLVM docs
            "--depth", "1",
            "https://github.com/llvm/llvm-project.git",
            str(LLVM_SRC_DIR)
        ])
    else:
        print(f"Existing LLVM repo in {LLVM_SRC_DIR}, fetching latest ...")
        run(["git", "fetch"], cwd=LLVM_SRC_DIR)
        run(["git", "checkout", LLVM_VERSION], cwd=LLVM_SRC_DIR)
        run(["git", "pull"], cwd=LLVM_SRC_DIR)

# =====Build LLVM =====
def build_llvm():
    llvm_bin = BUILD_DIR / "bin" / ("llvm-config.exe" if detect_os() == "windows" else "llvm-config")
    if llvm_bin.exists():
        print("LLVM already built, skipping rebuild.")
        return

    BUILD_DIR.mkdir(parents=True, exist_ok=True)

    cmake_args = [
        "cmake", "-S", "llvm", "-B", "build",
        "-G", "Ninja",
        "-DCMAKE_BUILD_TYPE=Release",
        "-DLLVM_ENABLE_PROJECTS='clang;lld;mlir;clang-tools-extra'",
        "-DLLVM_ENABLE_RUNTIMES=compiler-rt",
        f"-DLLVM_TARGETS_TO_BUILD={LLVM_TARGETS}",
        "-DLLVM_AR=lib.exe"
    ]

    using_mscv = True
    # using_mscv = False
    if using_mscv:
        # MSCV
        cmake_args += [
            "-DCMAKE_C_COMPILER=cl",
            "-DCMAKE_CXX_COMPILER=cl",
            # "-DCMAKE_ASM_COMPILER=cl",
            # "-DCMAKE_LINKER=link"
        ]
    else:
        #CLANG
        cmake_args += [ 
            "-DCMAKE_C_COMPILER=clang",
            "-DCMAKE_CXX_COMPILER=clang++",
            "-DCMAKE_ASM_COMPILER=clang",
            "-DCMAKE_LINKER=lld"
        ]


    
    build_cmd = ["ninja", f"-j{os.cpu_count() - 1}"]
    # build_cmd = ["cmake", "--build", "build", "--config", "Release", "--target", "check-all", f"-j{os.cpu_count() - 1}"]

    env = os.environ.copy()

    print("💡 Running CMake...")
    run(cmake_args, cwd=LLVM_SRC_DIR, env=env)
    print("💡 Building LLVM...")
    run(build_cmd, cwd=LLVM_SRC_DIR, env=env)

# =====Install go-llvm bindings ===== experiment
def install_go_llvm():
    env = os.environ.copy()
    env["LLVM_CONFIG"] = llvm_config_path()
    run(["go", "get", "-u", "tinygo.org/x/go-llvm"], env=env)

# =====Build Go project =====
def build_go_project():
    env = os.environ.copy()
    env["LLVM_CONFIG"] = llvm_config_path()
    run(["go", "mod", "tidy"], cwd=GO_PROJECT_DIR, env=env)
    run(["go", "build", "--tags=llvm21", "main.go"], cwd=GO_PROJECT_DIR, env=env)

# =====Run Go project =====
def run_go_project():
    env = os.environ.copy()
    env["LLVM_CONFIG"] = llvm_config_path()
    run(["go", "run", "--tags=llvm21", "main.go"], cwd=GO_PROJECT_DIR, env=env)

# ===== MAIN =====
def main():
    # Parse action
    action = sys.argv[1].lower() if len(sys.argv) > 1 else "all"
    if action not in ("build", "run", "all", "clean"):
        print("Usage: python build.py [build|run|all]clean")
        sys.exit(1)
    
    if action in ("clean"):
        env = os.environ.copy()
        run(["make", "clean"], cwd=LLVM_SRC_DIR, env=env)
        quit(0)

    ensure_llvm()
    build_llvm()
    print(f"LLVM_CONFIG set to: {llvm_config_path()}")
    install_go_llvm()

    if action in ("build", "all"):
        build_go_project()

    if action in ("run", "all"):
        run_go_project()

    print("✅ Done!")

if __name__ == "__main__":
    main()
