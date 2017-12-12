package abstract_factory

type Item interface {
    toString() string
}

type Link interface {
    Item
}

type Tray interface {
    Item
    AddToTray(item Item)
}

type baseTray struct {
    tray []Item
}

func (self *baseTray) AddToTray(item Item) {
    self.tray = append(self.tray, item)
}

type Page interface {
    AddToContent(item Item)
    Output() string
}

type basePage struct {
    content []Item
}

func (self *basePage) AddToContent(item Item) {
    self.content = append(self.content, item)
}

type Factory interface {
    CreateLink(caption, url string) Link
    CreateTray(caption string) Tray
    CreatePage(title, author string) Page
}

type mdLink struct {
    caption, url string
}

func (self *mdLink) toString() string {
    return "[" + self.caption + "](" + self.url + ")"
}

type mdTray struct {
    baseTray
    caption string
}

func (self *mdTray) toString() string {
    tray := "- " + self.caption + "\n"
    for _, Item := range self.tray {
        tray += Item.toString() + "\n"
    }
    return tray
}

type mdPage struct {
    basePage
    title, author string
}

func (self *mdPage) Output() string {
    content := "title: " + self.title + "\n"
    content += "author: " + self.author + "\n"
    for _, Item := range self.content {
        content += Item.toString() + "\n"
    }
    return content
}

type MdFactory struct {
}

func (self *MdFactory) CreateLink(caption, url string) link {
    return &mdLink{caption, url}
}
func (self *MdFactory) CreateTray(caption string) tray {
    return &mdTray{caption: caption}
}
func (self *MdFactory) CreatePage(title, author string) page {
    return &mdPage{title: title, author: author}
}
