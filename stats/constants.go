package stats

const typescript = "Typescript"
const javascript = "Javascript"
const gql = "GraphQL"
const html = "HTML"
const css = "CSS"
const less = "LESS"
const sass = "SASS"
const golang = "GO"
const ruby = "Ruby"
const scala = "Scala"
const perl = "Perl"
const python = "Python"
const csharp = "C#"
const swift = "Swift"
const c = "C"
const cpp = "C++"
const java = "Java"
const groovy = "Groovy"
const yml = "YAML"
const xml = "XML"
const assets = "Assets"
const php = "PHP"
const jSON = "JSON"
const sql = "SQL"
const md = "Markdown"
const text = "Text"
const kotlin = "Kotlin"
const shell = "Shell"
const r = "R"
const clojure = "Clojure"
const lua = "Lua"
const tests = "Tests"
const other = "Other"
const powershell = "Powershell"
const cassandra = "Cassandra"
const executable = "Executable"

var extensionMap = map[string]string{
	"ts":         typescript,
	"tsx":        typescript,
	"js":         javascript,
	"jsx":        javascript,
	"ejs":        javascript,
	"jst":        javascript,
	"pug":        javascript,
	"gql":        gql,
	"html":       html,
	"xhtml":      html,
	"mustache":   html,
	"handlebars": html,
	"cshtml":     html,
	"css":        css,
	"scss":       css,
	"less":       less,
	"overrides":  less,
	"variables":  less,
	"sass":       sass,
	"go":         golang,
	"rb":         ruby,
	"sc":         scala,
	"scala":      scala,
	"pl":         perl,
	"pm":         perl,
	"t":          perl,
	"py":         python,
	"j2":         python,
	"pyc":        python,
	"cs":         csharp,
	"swift":      swift,
	"h":          c,
	"c":          c,
	"m":          c,
	"so":         c,
	"cpp":        cpp,
	"java":       java,
	"class":      java,
	"jar":        java,
	"yml":        yml,
	"yaml":       yml,
	"map":        "Sourcemaps",
	"xml":        xml,
	"XML":        xml,
	"png":        assets,
	"svg":        assets,
	"ttf":        assets,
	"woff":       assets,
	"woff2":      assets,
	"jpg":        assets,
	"jpeg":       assets,
	"gif":        assets,
	"php":        php,
	"phpt":       php,
	"json":       jSON,
	"sql":        sql,
	"md":         md,
	"snap":       tests,
	"sh":         shell,
	"bash":       shell,
	"txt":        text,
	"kt":         kotlin,
	"R":          r,
	"Rd":         r,
	"groovy":     groovy,
	"ps1":        powershell,
	"lua":        lua,
	"clj":        clojure,
	"cql":        cassandra,
	"exe":        executable,
}
