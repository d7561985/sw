package newapp

import (
	"encoding/json"
	"fmt"
	"github.com/astaxie/beego/swagger"
	"github.com/beego/bee/logger"
	"github.com/d7561985/sw/generators/buffalo"
	"go/parser"
	"go/token"
	"gopkg.in/yaml.v2"
	"log"
	"os"
	"path"
	"strings"
	"unicode"
)

var GeneratorList = []string{"buffalo"}
var rootapi swagger.Swagger

type Generator struct {
	Type string
	gen  Gen
}

// New ...
func New(t string) (*Generator, error) {
	res := &Generator{Type: t}
	switch t {
	case "buffalo":
		res.gen = &buffalo.Buffalo{}
	default:
		log.Fatalf("not supported framework %s", t)
	}
	return res, res.Validate()
}

// Validate ...
func (g *Generator) Validate() error {
	ok := false
	for _, v := range GeneratorList {
		if v == g.Type {
			ok = true
		}
	}
	if !ok {
		return fmt.Errorf("format %q not supported", g.Type)
	}
	return nil
}

// Run ...
func (g *Generator) Run(dirpath string) error {
	return GenerateDocs(dirpath, g.gen)
}

// GenerateDocs generates documentations from given path.
func GenerateDocs(curpath string, gen Gen) error {

	fset := token.NewFileSet()
	f, err := parser.ParseFile(fset, gen.RootPath(curpath), nil, parser.ParseComments)
	if err != nil {
		return err
	}

	fmt.Println("===>>", gen.RootPath(curpath))

	rootapi.Infos = swagger.Information{}
	rootapi.SwaggerVersion = "2.0"

	// Analyse API comments
	if f.Comments != nil {
		for _, c := range f.Comments {
			for _, s := range strings.Split(c.Text(), "\n") {
				if strings.HasPrefix(s, "@APIVersion") {
					rootapi.Infos.Version = strings.TrimSpace(s[len("@APIVersion"):])
				} else if strings.HasPrefix(s, "@Title") {
					rootapi.Infos.Title = strings.TrimSpace(s[len("@Title"):])
				} else if strings.HasPrefix(s, "@Description") {
					rootapi.Infos.Description = strings.TrimSpace(s[len("@Description"):])
				} else if strings.HasPrefix(s, "@TermsOfServiceUrl") {
					rootapi.Infos.TermsOfService = strings.TrimSpace(s[len("@TermsOfServiceUrl"):])
				} else if strings.HasPrefix(s, "@Contact") {
					rootapi.Infos.Contact.EMail = strings.TrimSpace(s[len("@Contact"):])
				} else if strings.HasPrefix(s, "@Name") {
					rootapi.Infos.Contact.Name = strings.TrimSpace(s[len("@Name"):])
				} else if strings.HasPrefix(s, "@URL") {
					rootapi.Infos.Contact.URL = strings.TrimSpace(s[len("@URL"):])
				} else if strings.HasPrefix(s, "@LicenseUrl") {
					if rootapi.Infos.License == nil {
						rootapi.Infos.License = &swagger.License{URL: strings.TrimSpace(s[len("@LicenseUrl"):])}
					} else {
						rootapi.Infos.License.URL = strings.TrimSpace(s[len("@LicenseUrl"):])
					}
				} else if strings.HasPrefix(s, "@License") {
					if rootapi.Infos.License == nil {
						rootapi.Infos.License = &swagger.License{Name: strings.TrimSpace(s[len("@License"):])}
					} else {
						rootapi.Infos.License.Name = strings.TrimSpace(s[len("@License"):])
					}
				} else if strings.HasPrefix(s, "@Schemes") {
					rootapi.Schemes = strings.Split(strings.TrimSpace(s[len("@Schemes"):]), ",")
				} else if strings.HasPrefix(s, "@Host") {
					rootapi.Host = strings.TrimSpace(s[len("@Host"):])
				} else if strings.HasPrefix(s, "@SecurityDefinition") {
					if len(rootapi.SecurityDefinitions) == 0 {
						rootapi.SecurityDefinitions = make(map[string]swagger.Security)
					}
					var out swagger.Security
					p := getparams(strings.TrimSpace(s[len("@SecurityDefinition"):]))
					if len(p) < 2 {
						beeLogger.Log.Fatalf("Not enough params for security: %d\n", len(p))
					}
					out.Type = p[1]
					switch out.Type {
					case "oauth2":
						if len(p) < 6 {
							beeLogger.Log.Fatalf("Not enough params for oauth2: %d\n", len(p))
						}
						if !(p[3] == "implicit" || p[3] == "password" || p[3] == "application" || p[3] == "accessCode") {
							beeLogger.Log.Fatalf("Unknown flow type: %s. Possible values are `implicit`, `password`, `application` or `accessCode`.\n", p[1])
						}
						out.AuthorizationURL = p[2]
						out.Flow = p[3]
						if len(p)%2 != 0 {
							out.Description = strings.Trim(p[len(p)-1], `" `)
						}
						out.Scopes = make(map[string]string)
						for i := 4; i < len(p)-1; i += 2 {
							out.Scopes[p[i]] = strings.Trim(p[i+1], `" `)
						}
					case "apiKey":
						if len(p) < 4 {
							beeLogger.Log.Fatalf("Not enough params for apiKey: %d\n", len(p))
						}
						if !(p[3] == "header" || p[3] == "query") {
							beeLogger.Log.Fatalf("Unknown in type: %s. Possible values are `query` or `header`.\n", p[4])
						}
						out.Name = p[2]
						out.In = p[3]
						if len(p) > 4 {
							out.Description = strings.Trim(p[4], `" `)
						}
					case "basic":
						if len(p) > 2 {
							out.Description = strings.Trim(p[2], `" `)
						}
					default:
						beeLogger.Log.Fatalf("Unknown security type: %s. Possible values are `oauth2`, `apiKey` or `basic`.\n", p[1])
					}
					rootapi.SecurityDefinitions[p[0]] = out
				} else if strings.HasPrefix(s, "@Security") {
					if len(rootapi.Security) == 0 {
						rootapi.Security = make([]map[string][]string, 0)
					}
					rootapi.Security = append(rootapi.Security, getSecurity(s))
				}
			}
		}
	}

	os.Mkdir(path.Join(curpath, "swagger"), 0755)
	fd, err := os.Create(path.Join(curpath, "swagger", "swagger.json"))
	if err != nil {
		panic(err)
	}
	fdyml, err := os.Create(path.Join(curpath, "swagger", "swagger.yml"))
	if err != nil {
		panic(err)
	}
	defer fdyml.Close()
	defer fd.Close()
	dt, err := json.MarshalIndent(rootapi, "", "    ")
	dtyml, erryml := yaml.Marshal(rootapi)
	if err != nil || erryml != nil {
		panic(err)
	}
	_, err = fd.Write(dt)
	_, erryml = fdyml.Write(dtyml)
	if err != nil || erryml != nil {
		panic(err)
	}
	return nil
}

func getSecurity(t string) (security map[string][]string) {
	security = make(map[string][]string)
	p := getparams(strings.TrimSpace(t[len("@Security"):]))
	if len(p) == 0 {
		beeLogger.Log.Fatalf("No params for security specified\n")
	}
	security[p[0]] = make([]string, 0)
	for i := 1; i < len(p); i++ {
		security[p[0]] = append(security[p[0]], p[i])
	}
	return
}

// analisys params return []string
// @Param	query		form	 string	true		"The email for login"
// [query form string true "The email for login"]
func getparams(str string) []string {
	var s []rune
	var j int
	var start bool
	var r []string
	var quoted int8
	for _, c := range str {
		if unicode.IsSpace(c) && quoted == 0 {
			if !start {
				continue
			} else {
				start = false
				j++
				r = append(r, string(s))
				s = make([]rune, 0)
				continue
			}
		}

		start = true
		if c == '"' {
			quoted ^= 1
			continue
		}
		s = append(s, c)
	}
	if len(s) > 0 {
		r = append(r, string(s))
	}
	return r
}
