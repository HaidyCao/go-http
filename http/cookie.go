package http

import (
	"errors"
	"net/http"
	"net/url"
	"time"
)

type GoCookie struct {
	Name  string
	Value string

	Path       string // optional
	Domain     string // optional
	Expires    string // optional
	RawExpires string // for reading cookies only

	// MaxAge=0 means no 'Max-Age' attribute specified.
	// MaxAge<0 means delete cookie now, equivalently 'Max-Age: 0'
	// MaxAge>0 means Max-Age attribute present and given in seconds
	MaxAge   int
	Secure   bool
	HttpOnly bool
	Raw      string
	Unparsed []string // Raw text of unparsed attribute-value pairs
}

func (cookie *GoCookie) AddUnparsed(unparsed string) {
	cookie.Unparsed = append(cookie.Unparsed, unparsed)
}

func (cookie *GoCookie) GenHttpCookie() (*http.Cookie, error) {
	expires, err := time.Parse("2006-01-02 15:04:05", cookie.Expires)
	if err != nil {
		return nil, err
	}

	return &http.Cookie{
		Name:       cookie.Name,
		Value:      cookie.Value,
		Path:       cookie.Path,
		Domain:     cookie.Domain,
		Expires:    expires,
		RawExpires: cookie.RawExpires,
		MaxAge:     cookie.MaxAge,
		Secure:     cookie.Secure,
		HttpOnly:   cookie.HttpOnly,
		Raw:        cookie.Raw,
		Unparsed:   cookie.Unparsed,
	}, nil
}

type GoCookies struct {
	cookies []*GoCookie
}

func (cookies *GoCookies) AppendCookie(cookie *GoCookie) {
	cookies.cookies = append(cookies.cookies, cookie)
}

// cookie jar

type GoCookieJar struct {
	Jar http.CookieJar
}

func (jar *GoCookieJar) SetCookies(urlStr string, goCookies *GoCookies) error {
	url, err := url.Parse(urlStr)
	if err != nil {
		return errors.New("url is illegal")
	}
	cookies := make([]*http.Cookie, len(goCookies.cookies))
	for i := 0; i < len(cookies); i++ {
		c, err := goCookies.cookies[i].GenHttpCookie()
		if err != nil {
			return err
		}
		cookies[i] = c
	}

	jar.Jar.SetCookies(url, cookies)
	return nil
}

func (jar *GoCookieJar) Cookies(urlStr string) (*GoCookies, error) {
	url, err := url.Parse(urlStr)
	if err != nil {
		return nil, errors.New("url is illegal")
	}

	cookies := jar.Jar.Cookies(url)
	goCookies := &GoCookies{}
	for i := 0; i < len(cookies); i++ {
		cookie := cookies[i]
		goCookies.AppendCookie(&GoCookie{
			Name:       cookie.Name,
			Value:      cookie.Value,
			Path:       cookie.Path,
			Domain:     cookie.Domain,
			Expires:    cookie.Expires.Format("2006-01-02 15:04:05"),
			RawExpires: cookie.RawExpires,
			MaxAge:     cookie.MaxAge,
			Secure:     cookie.Secure,
			HttpOnly:   cookie.HttpOnly,
			Raw:        cookie.Raw,
			Unparsed:   cookie.Unparsed,
		})
	}
	return goCookies, nil
}
