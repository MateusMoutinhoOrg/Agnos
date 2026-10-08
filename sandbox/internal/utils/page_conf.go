package utils

import (
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/api"
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/stagedfs"
)

// A page is a file of assets/front/ and nothing else: no route, no
// declaration. The front route front-init writes serves the whole tree, so
// a page exists the moment its html does, and a bundler's dist/ dropped there
// is as much the front as one add-page wrote. The helpers below are the one
// place a page name is spelled as a file, so add-page, remove-page and the
// interview all read it the same way.

// FrontDir is the root of the asset tree the front layer serves.
// Everything under it is the project's own content: no build writes there, and
// front-purge leaves it whole.
const FrontDir = "assets/front"

// FrontRouteName is the route front-init writes to serve FrontDir.
const FrontRouteName = "front"

// FrontIndexPage is the page front-init scaffolds, the one "/" answers.
const FrontIndexPage = "index"

// FrontNotFoundPage is the page front-init scaffolds for a path that names
// no file: the front route answers it with a 404. It is the
// OpinionatedAgnosFront lib's NotFound spelled as a page name.
const FrontNotFoundPage = "404"

// pageExtension is what a page name is written with on disk. The front
// route tries it after the bare path, so /about answers about.html.
const pageExtension = ".html"

// PageName normalizes what add-page and remove-page are handed into the name
// a page is known by: slashes kept, surrounding slashes and a trailing .html
// dropped, so "blog/post", "/blog/post" and "blog/post.html" are one page.
func PageName(sandbox *api.Sandbox, name string) string {
	trimmed := sandbox.Deps.StringsDeps.Trim(sandbox.Deps.StringsDeps.TrimSpace(name), "/")
	return sandbox.Deps.StringsDeps.TrimSuffix(trimmed, pageExtension)
}

// ValidatePageName refuses a name that is not a plain path under FrontDir:
// every segment one or more of letters, digits, '-', '_' and '.', and none of
// them "." or "..". It is the same shape the front route accepts from a
// request, so every page add-page writes is one that route can answer.
func ValidatePageName(sandbox *api.Sandbox, name string) error {
	page := PageName(sandbox, name)
	if page == "" {
		return sandbox.Deps.StdDeps.Errorf("a page needs a name")
	}

	for _, segment := range sandbox.Deps.StringsDeps.Split(page, "/") {
		if segment == "" || segment == "." || segment == ".." {
			return sandbox.Deps.StdDeps.Errorf("invalid page name %q: every segment of the path names a directory or the file", name)
		}
		for _, letter := range segment {
			valid := (letter >= 'a' && letter <= 'z') ||
				(letter >= 'A' && letter <= 'Z') ||
				(letter >= '0' && letter <= '9') ||
				letter == '-' || letter == '_' || letter == '.'
			if !valid {
				return sandbox.Deps.StdDeps.Errorf("invalid page name %q: only letters, digits, dashes, underscores, dots and slashes are allowed (it becomes %s)",
					name, PageAsset(sandbox, name))
			}
		}
	}
	return nil
}

// PageAsset is the project-relative path of a page's html file.
func PageAsset(sandbox *api.Sandbox, name string) string {
	return FrontDir + "/" + PageName(sandbox, name) + pageExtension
}

// PageUrl is the path the front route answers a page on: "/" for the index
// page, "/<name>" for every other one.
func PageUrl(sandbox *api.Sandbox, name string) string {
	page := PageName(sandbox, name)
	if page == FrontIndexPage {
		return "/"
	}
	return "/" + sandbox.Deps.StringsDeps.TrimSuffix(page, "/"+FrontIndexPage)
}

// ListPages is every html file under FrontDir, named the way add-page and
// remove-page spell a page, sorted.
func ListPages(sandbox *api.Sandbox, io *stagedfs.StagedFS) []string {
	if !io.IsDir(FrontDir) {
		return []string{}
	}

	pages := []string{}
	for _, file := range io.ListFilesRecursively(FrontDir) {
		relative := sandbox.Deps.StringsDeps.TrimPrefix(file, FrontDir+"/")
		if !sandbox.Deps.StringsDeps.HasSuffix(relative, pageExtension) {
			continue
		}
		pages = append(pages, PageName(sandbox, relative))
	}
	sandbox.Deps.SortDeps.Strings(pages)
	return pages
}

// WritePage renders the scaffold of one page, plain html with no template
// syntax left in it, into PageAsset. io.WriteFile refuses an existing file, so
// a caller that means to keep one checks first.
func WritePage(sandbox *api.Sandbox, io *stagedfs.StagedFS, name string, title string) error {
	page := PageAsset(sandbox, name)
	sandbox.Deps.StdDeps.Logf("creating %s, answered on %s \n", page, PageUrl(sandbox, name))

	content, err := sandbox.Deps.EmbedDeps.RenderTemplate("templates/front_page.html", map[string]interface{}{
		"Title": title,
		"File":  PageName(sandbox, name) + pageExtension,
	})
	if err != nil {
		return err
	}
	return io.CreateFile(page, content)
}

// WriteNotFoundPage renders the scaffold of FrontNotFoundPage, a formatted
// page with no template syntax left in it, into its PageAsset. Like WritePage,
// io.WriteFile refuses an existing file, so a caller that means to keep one
// checks first.
func WriteNotFoundPage(sandbox *api.Sandbox, io *stagedfs.StagedFS) error {
	page := PageAsset(sandbox, FrontNotFoundPage)
	sandbox.Deps.StdDeps.Logf("creating %s, answered on every path that names no file \n", page)

	content, err := sandbox.Deps.EmbedDeps.RenderTemplate("templates/front_404.html", map[string]interface{}{})
	if err != nil {
		return err
	}
	return io.CreateFile(page, content)
}
