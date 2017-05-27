type InMemoryCookieJar struct{
	storage map[string][]http.Cookie
}

// buggy... but works

func (jar InMemoryCookieJar) SetCookies(u *url.URL, cookies []*http.Cookie) {
	for _, ck := range cookies {
		path := ck.Domain
		if ck.Path != "/" {
			path = ck.Domain + ck.Path
		}
		if ck.Domain[0] == '.' {
			path = path[1:]		// FIXME: .hi.baidu.com won't match hi.baidu.com
		}
		if _, found := jar.storage[path]; found {
			setted := false
			for i, v := range jar.storage[path] {
				if v.Name == ck.Name {
					jar.storage[path][i] = *ck
					setted = true
					break
				}
			}
			if ! setted {
				jar.storage[path] = append(jar.storage[path], *ck)
			}
		} else {
			jar.storage[path] = []http.Cookie{*ck}
		}
	}
}

func (jar InMemoryCookieJar) Cookies(u *url.URL) []*http.Cookie {
	cks := []*http.Cookie{}
	log.Println("get cookies", u)
	for pattern, cookies := range jar.storage {
		if strings.Contains(u.String(), pattern) {
			for i := range cookies {
				cks = append(cks, &cookies[i])
				log.Println("add cookie", cookies[i].Name)
			}
		}
	}
	return cks
}

func newCookieJar() InMemoryCookieJar {
	storage := make(map[string][]http.Cookie)
	return InMemoryCookieJar{storage}
}
