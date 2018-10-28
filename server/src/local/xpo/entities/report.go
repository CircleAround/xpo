package entities

import (
	"time"

	"google.golang.org/appengine/datastore"
)

var Languages = []string {"1c", "abnf", "accesslog", "actionscript", "ada", "apache", "applescript", "cpp", "arduino", "armasm", "xml", "asciidoc", "aspectj", "autohotkey", "autoit", "avrasm", "awk", "axapta", "bash", "basic", "bnf", "brainfuck", "cal", "capnproto", "ceylon", "clean", "clojure", "clojure-repl", "cmake", "coffeescript", "coq", "cos", "crmsh", "crystal", "cs", "csp", "css", "d", "markdown", "dart", "delphi", "diff", "django", "dns", "dockerfile", "dos", "dsconfig", "dts", "dust", "ebnf", "elixir", "elm", "ruby", "erb", "erlang-repl", "erlang", "excel", "fix", "flix", "fortran", "fsharp", "gams", "gauss", "gcode", "gherkin", "glsl", "go", "golo", "gradle", "groovy", "haml", "handlebars", "haskell", "haxe", "hsp", "htmlbars", "http", "hy", "inform7", "ini", "irpf90", "java", "javascript", "jboss-cli", "json", "julia", "julia-repl", "kotlin", "lasso", "ldif", "leaf", "less", "lisp", "livecodeserver", "livescript", "llvm", "lsl", "lua", "makefile", "mathematica", "matlab", "maxima", "mel", "mercury", "mipsasm", "mizar", "perl", "mojolicious", "monkey", "moonscript", "n1ql", "nginx", "nimrod", "nix", "nsis", "objectivec", "ocaml", "openscad", "oxygene", "parser3", "pf", "php", "pony", "powershell", "processing", "profile", "prolog", "protobuf", "puppet", "purebasic", "python", "q", "qml", "r", "rib", "roboconf", "routeros", "rsl", "ruleslanguage", "rust", "scala", "scheme", "scilab", "scss", "shell", "smali", "smalltalk", "sml", "sqf", "sql", "stan", "stata", "step21", "stylus", "subunit", "swift", "taggerscript", "yaml", "tap", "tcl", "tex", "thrift", "tp", "twig", "typescript", "vala", "vbnet", "vbscript", "vbscript-html", "verilog", "vhdl", "vim", "x86asm", "xl", "xquery", "zephir" }

// Report struct
type Report struct {
	ID             int64          `json:"id" datastore:"-" goon:"id"`
	AuthorKey      *datastore.Key `json:"-" datastore:"-" goon:"parent" validate:"required"`
	AuthorID       string         `json:"authorId" validate:"required"`
	Author         string         `json:"author" validate:"required"`
	AuthorNickname string         `json:"authorNickname" validate:"required"`
	Content        string         `json:"content" validate:"required,max=20000" datastore:"Content,noindex"`
	ContentType    string         `json:"contentType" validate:"required"`
	Languages      []string       `json:"languages" validate:"oneof=1c abnf accesslog actionscript ada apache applescript cpp arduino armasm xml asciidoc aspectj autohotkey autoit avrasm awk axapta bash basic bnf brainfuck cal capnproto ceylon clean clojure clojure-repl cmake coffeescript coq cos crmsh crystal cs csp css d markdown dart delphi diff django dns dockerfile dos dsconfig dts dust ebnf elixir elm ruby erb erlang-repl erlang excel fix flix fortran fsharp gams gauss gcode gherkin glsl go golo gradle groovy haml handlebars haskell haxe hsp htmlbars http hy inform7 ini irpf90 java javascript jboss-cli json julia julia-repl kotlin lasso ldif leaf less lisp livecodeserver livescript llvm lsl lua makefile mathematica matlab maxima mel mercury mipsasm mizar perl mojolicious monkey moonscript n1ql nginx nimrod nix nsis objectivec ocaml openscad oxygene parser3 pf php pony powershell processing profile prolog protobuf puppet purebasic python q qml r rib roboconf routeros rsl ruleslanguage rust scala scheme scilab scss shell smali smalltalk sml sqf sql stan stata step21 stylus subunit swift taggerscript yaml tap tcl tex thrift tp twig typescript vala vbnet vbscript vbscript-html verilog vhdl vim x86asm xl xquery zephir"`
	ReportedAt     time.Time      `json:"reportedAt"`
	CreatedAt      time.Time      `json:"createdAt"`
	UpdatedAt      time.Time      `json:"updatedAt"`
}
