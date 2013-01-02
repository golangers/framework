package urlmanage

import (
	"golanger.com/framework/log"
	"net/http"
	"regexp"
	"strconv"
	"strings"
	"sync"
)

func parse(r string) (string, string, []flag) {
	var expr, replace string
	var flags []flag
	r = regexp.MustCompile(`[[:blank:]]+`).ReplaceAllString(r, "`")
	rs := strings.Split(r, "`")
	lrs := len(rs)
	if lrs >= 2 {
		expr = rs[0]
		replace = rs[1]
		if lrs >= 3 {
			flagss := strings.Split(strings.Trim(rs[2], `[]`), ",")
			for _, f := range flagss {
				if f == "NC" {
					expr = `(?i)` + expr
					continue
				}

				var fl flag
				fs := strings.Split(f, "=")
				fl = flag{
					name: fs[0],
				}
				if len(fs) > 1 {
					fl.param = fs[1]
				}

				flags = append(flags, fl)
			}
		}
	}

	return expr, replace, flags
}

/*
NC - Nocase:URL地址匹配对大小写敏感
S - Skip:忽略之后的规则
R - Redirect:发出一个HTTP重定向
N - Next:再次从第一个规则开始处理，但是使用当前重写后的URL地址
L - Last:停止处理接下来的规则
QSA - Qsappend:在新的URL地址后附加查询字符串部分，而不是替代
*/
type flag struct {
	name  string
	param string
}

type rule struct {
	regexp  *regexp.Regexp
	replace string
	flags   []flag
}

type UrlManage struct {
	manage bool
	rules  []rule
	mutex  sync.RWMutex
}

func New() *UrlManage {
	return &UrlManage{}
}

func (u *UrlManage) Manage() bool {
	return u.manage
}

func (u *UrlManage) Start() {
	u.mutex.Lock()
	u.manage = true
	u.mutex.Unlock()
}

func (u *UrlManage) Stop() {
	u.mutex.Lock()
	u.manage = false
	u.mutex.Unlock()
}

func (u *UrlManage) addRule(expr, replace string, flags ...flag) {
	if expr == "" || replace == "" {
		log.Warn("UrlManage.addUrl: expr and reolace is empty")
		return
	}

	r, err := regexp.Compile(expr)
	if err != nil {
		log.Warn("UrlManage.addUrl: regexp compile failed - ", err)
		return
	}

	rl := rule{
		regexp:  r,
		replace: replace,
		flags:   flags,
	}

	u.mutex.Lock()
	u.rules = append(u.rules, rl)
	u.mutex.Unlock()
}

func (u *UrlManage) doRule(rw http.ResponseWriter, req *http.Request) string {
	in := req.URL.Path
	out := in
	u.mutex.RLock()
	lrs := len(u.rules)
	u.mutex.RUnlock()
RESTART:
	for i := 0; i < lrs; i++ {
		u.mutex.RLock()
		ur := u.rules[i]
		u.mutex.RUnlock()
		if !ur.regexp.MatchString(in) {
			continue
		}

		var skip bool
		var restart bool
		var last bool
		var redirect bool
		var redirectCode int
		var queryStringAppend bool

		if len(ur.flags) > 0 {
			for _, f := range ur.flags {
				switch f.name {
				case "R":
					redirect = true
					redirectCode, _ = strconv.Atoi(f.param)
					if redirectCode == 0 {
						redirectCode = http.StatusFound
					}
				case "S":
					skip = true
					skipNum, _ := strconv.Atoi(f.param)
					//循环后会自动加1，所以这里减1
					skipNum = skipNum - 1
					if skipNum > 0 {
						i = i + skipNum
					}
				case "N":
					restart = true
				case "L":
					last = true
				case "QSA":
					queryStringAppend = true
				}
			}
		}

		out = ur.regexp.ReplaceAllString(in, ur.replace)

		if queryStringAppend {
			if strings.Contains(out, "?") {
				out += "&"
			} else {
				out += "?"
			}

			out += req.URL.RawQuery
		}

		if redirect {
			http.Redirect(rw, req, out, redirectCode)
			return `redirect`
		}

		if skip {
			continue
		}

		if restart {
			in = out
			goto RESTART
		}

		if last {
			break
		}
	}

	return out
}

func (u *UrlManage) clearRule() {
	u.mutex.Lock()
	defer u.mutex.Unlock()
	u.rules = make([]rule, 0)
}

func (u *UrlManage) loadRule(rules string) {
	for _, r := range strings.Split(rules, "\n") {
		u.AddRule(r)
	}
}

func (u *UrlManage) AddRule(r string) {
	expr, replace, flags := parse(r)
	u.addRule(expr, replace, flags...)
}

func (u *UrlManage) ReWrite(rw http.ResponseWriter, req *http.Request) string {
	out := req.URL.String()
	u.mutex.RLock()
	manage := u.manage
	u.mutex.RUnlock()
	if manage {
		out = u.doRule(rw, req)
	}

	return out
}

func (u *UrlManage) LoadRule(rules string, reload bool) {
	if reload {
		u.clearRule()
	}

	u.loadRule(strings.TrimSpace(rules))
}
