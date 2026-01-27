
---

## 1) Ruzta Language Spec v0.1 (Draft)

| Rule / Feature           | Decision                                                                        |
| ------------------------ | --------------------------------------------------------------------------------|
| File Source units        | `.rz` primary source, `.rc` codegen binary LLVM IR. File are classes by default |
| Comments                 | `//`, `#` line, `/* ... */` block                                               |
| Identifers               | ASCII/`_` + Unicode letters/digits/`_`                                          |
| Terminator               | `;` optional, newline can end stmt(golang-like)                                 |
| Scope                    | `{ ... }` defines scope always (no indentation semantics)                       |
| Variables                | `var name = expr;` (mutable) and `const name = expr;` (immutable)               |
| Gradual typing           | `var x Int = 3;`, `var x = 3;`(infered), `var x := 3; var y`(variant typed)     |
| Definition               | `fn name(params) Type { ... }` return type optional                             |
| Default args             | `fn f(x = 3) {}` allowed                                                        |
| Traits: Interface + mixin| `trait T { fn m(self, ...) void { ... } }` use `uses ...` to include the trait  |
| Member Receiver          | `self`  but members are available in method scope and can not be shadowed.      |
| Members Allowed          | Traits & classes may contain: variables, functions, signals, annotations, inner
 traits, inner classes. functions in traits can be bodyless in traits                                        |
| Module                   | `mod Name { ... }` hold for classes,traits, other modulues                      |
| Import system            | `import "../path_to_file"/mod.class_name.inner_child as alias;`                 |
| Annotations              | `@feature`, `@rpc`, `@export`, `@onready` - applyed members                     |
| Calls                    | `f(a, b)` and method calls `obj.m(x)`                                           |
| If/while/for             | `if (...) {}`; `while (...) {}`; `for x in iter {}`                             |
| Match                    | `match (expr) { pattern: checks [when guard_condition] ... _ { stmt_or_block }}` Brace-based cases.`_` for default, optional guards . No fallthrough; first match wins. Example: 
`match (x) { 1 when t == true { print("one") } 2,3 { print("two or 3") } _: { print("other") } }`            |
| Signals                  |  `signal signal_name(param)` and `signal_name.connect(handler)`                 |
| Normal constructor       | `Type.new()` -> invokes `init` (if defined)                                     |
| Builder constructor      | `Type { new(...); prop = expr; ChildType { ... } }` creates object, assigns 
props, calls, and nested adds children via `add_child`. `init` if defined with param that dont have default 
must call `new(x,y)` it top of body                                                                          |
| Class Inheritance        | `class Clildclass extends ParentClass { ... }` (optional single inheritance)    |
| Alias system             | `type Type/class as alias;`                                                     |
| CADRe/RAII               | Constructor Acquires, Destructor Releases `func init(...)` then `func deinit(...)' |


---

## 2) Ruzta Constants

- PI = 3.14159265358979
- TAU = 6.28318530717959
- INF = inf
- NAN = nan

---

## 3) Ruzta DataType Spec

1. bool - Default: false
2. byte / i8 - Default: 0
3. int / i32 - Default: 0
4. long / i64 - Default: 0
4. i128- Default: 0
5. float / f32 - Default: 0.0
6. double / f64 - Default: 0.0

7. string  - Default: ""
8. variant (hold any thing) - Default: null

9. array (can be typed [...]) - Default: []
10. dictionary/ dict  (can be typed [...,...]) - Default: {}

11. class - hold class as a meta type - Default: null
12. trait - hold trait as a meta type - Default: null

13. Signal (first class) - Default: null

14. <Object> (user defined classes) - Default: null

15. Enum 


---

## 5) Ruzta Integer Spec

support digit separators: `_` or `,`

1. base 10 - `45`
2. binary / base 2 - `0b1010`/ `0B1010`
3. Octal / base 8 - `0o755`/ `0O755`
4. Hex / base 16 - `0xFF`/ `0Xff`
5. Scientific notation - `1e-10` / `1E-10`

---

## 6) Ruzta Annotations Spec

can have a block `{...}` for scope

1. @feature(names: String, ...) - make the line, block, func, class only availabe in feature listed
2. @private - Access modifier not be overritten, object can not call but child can  
3. @abstract - Marks a class or a method as abstract.An abstract class is a class that cannot be instantiated directly.

2. @onready - class must have  fn named _ready()
4. @export
    - i) @export_category(name: String)
    - ii) @export_color_no_alpha()
    - iii) @export_custom(hint: PropertyHint, hint_string: String, usage: BitField[PropertyUsageFlags] = 6) 
    - iv) @export_dir() 
    - v) @export_enum(names: String, ...) vararg 
    - vi) @export_exp_easing(hints: String = "", ...) vararg
    - vii) @export_file(filter: String = "", ...) vararg
    - viii) @export_file_path(filter: String = "", ...) vararg
    - ix) @export_flags(names: String, ...) vararg 
    - x) @export_flags_2d_navigation() 
    - ...