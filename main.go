package main

import (
    "fmt"
    "tinygo.org/x/go-llvm"
)

func main() {
    llvm.InitializeNativeTarget()
    llvm.InitializeNativeAsmPrinter()
    llvm.InitializeNativeAsmParser()

    mod := llvm.NewModule("my_module")
    builder := llvm.NewBuilder()

    fnType := llvm.FunctionType(llvm.Int32Type(), nil, false)
    fn := llvm.AddFunction(mod, "foo", fnType)

    entry := llvm.AddBasicBlock(fn, "entry")
    builder.SetInsertPointAtEnd(entry)

    builder.CreateRet(llvm.ConstInt(llvm.Int32Type(), 42, false))

    fmt.Println(mod.String())
    fmt.Println("Hello, Ruzta World!")

}

