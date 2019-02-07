package consts

const Debug = true
const ServerVersion = "1.0.1"
const AssetsVersion = "1"
const AssetsPath = "assets"
const DirIndexFile = "index.html"

// Bootstrap resources
const AssetsBootstrapCss = AssetsPath + "/bootstrap.css"
const AssetsBootstrapJs = AssetsPath + "/bootstrap.js"
const AssetsJqueryJs = AssetsPath + "/jquery.js"
const AssetsPopperJs = AssetsPath + "/popper.js"

// System resources
const AssetsCpScriptsJs = AssetsPath + "/cp/scripts.js"
const AssetsCpStylesCss = AssetsPath + "/cp/styles.css"
const AssetsSysBgPng = AssetsPath + "/sys/bg.png"
const AssetsSysFaveIco = AssetsPath + "/sys/fave.ico"
const AssetsSysLogoPng = AssetsPath + "/sys/logo.png"
const AssetsSysLogoSvg = AssetsPath + "/sys/logo.svg"
const AssetsSysStylesCss = AssetsPath + "/sys/styles.css"

// Template data
type TmplSystem struct {
	PathIcoFav       string
	PathSvgLogo      string
	PathCssStyles    string
	PathCssCpStyles  string
	PathCssBootstrap string
	PathJsJquery     string
	PathJsPopper     string
	PathJsBootstrap  string
	PathJsCpScripts  string
}

type TmplError struct {
	ErrorMessage string
}

type TmplData struct {
	System TmplSystem
	Data   interface{}
}
