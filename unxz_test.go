package unxz

import (
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"testing"
)

func currentDir() string {
	dir, err := os.Getwd()
	if err != nil {
		return ""
	}
	return dir
}

func testNullContent(t *testing.T) {
	aFile, err := os.Open(path.Join(currentDir(), "./test/a.txt"))
	if err != nil {
		t.Fatal(err)
	}
	defer aFile.Close()

	aFileContent, err := ioutil.ReadAll(aFile)
	if err != nil {
		t.Fatal(err)
	}

	aOutFile, err := os.Open(path.Join(currentDir(), "./test/out/a.txt"))
	if err != nil {
		t.Fatal(err)
	}
	defer aOutFile.Close()

	aOutFileContent, err := ioutil.ReadAll(aOutFile)
	if err != nil {
		t.Fatal(err)
	}

	if string(aFileContent) != string(aOutFileContent) {
		t.Fatal("Unxz file content error.")
	}

	if string(aFileContent) != "" {
		t.Fatal("Unxz file content error.")
	}
}

func testContentEqual(t *testing.T) {
	bFile, err := os.Open(path.Join(currentDir(), "./test/b.txt"))
	if err != nil {
		t.Fatal(err)
	}
	defer bFile.Close()

	bFileContent, err := ioutil.ReadAll(bFile)
	if err != nil {
		t.Fatal(err)
	}

	bOutFile, err := os.Open(path.Join(currentDir(), "./test/out/b.txt"))
	if err != nil {
		t.Fatal(err)
	}
	defer bOutFile.Close()

	bOutFileContent, err := ioutil.ReadAll(bOutFile)
	if err != nil {
		t.Fatal(err)
	}

	if string(bFileContent) != string(bOutFileContent) {
		t.Log(string(bFileContent))
		t.Log(string(bOutFileContent))
		t.Fatal("Unxz file content error.")
	}
}

func testLink(t *testing.T) {
	cOutFile, err := os.OpenFile(path.Join(currentDir(), "./test/out/c.txt"), os.O_RDWR|os.O_CREATE|os.O_APPEND, 0644)
	if err != nil {
		t.Fatal(err)
	}
	defer cOutFile.Close()

	cOutFileContent, err := ioutil.ReadAll(cOutFile)
	if err != nil {
		t.Fatal(err)
	}

	cOutLinkFile, err := os.Open(path.Join(currentDir(), "./test/out/dir/c.txt"))
	if err != nil {
		t.Fatal(err)
	}
	defer cOutLinkFile.Close()

	cOutLinkFileContent, err := ioutil.ReadAll(cOutLinkFile)
	if err != nil {
		t.Fatal(err)
	}

	if string(cOutFileContent) != string(cOutLinkFileContent) {
		t.Fatal("Unxz file content error or link file content error.")
	}

	t.Log("Old c.txt content:", string(cOutFileContent))
	t.Log("Change c.txt content.")
	_, err = cOutFile.WriteString("123")
	if err != nil {
		t.Fatal(err)
	}

	cOutFileNew, err := os.Open(path.Join(currentDir(), "./test/out/c.txt"))
	if err != nil {
		t.Fatal(err)
	}
	defer cOutFileNew.Close()

	cOutFileContentNew, err := ioutil.ReadAll(cOutFileNew)
	if err != nil {
		t.Fatal(err)
	}
	t.Log("New c.txt content:", string(cOutFileContentNew))

	cOutLinkFileNew, err := os.Open(path.Join(currentDir(), "./test/out/dir/c.txt"))
	if err != nil {
		t.Fatal(err)
	}
	defer cOutLinkFileNew.Close()

	cOutLinkFileContentNew, err := ioutil.ReadAll(cOutLinkFileNew)
	if err != nil {
		t.Fatal(err)
	}

	if string(cOutFileContentNew) != string(cOutLinkFileContentNew) {
		t.Fatal("Unxz file link error.")
	}
}

func testSymlink(t *testing.T) {
	targetURL, err := filepath.EvalSymlinks(path.Join(currentDir(), "./test/out/d.txt"))
	if err != nil {
		t.Fatal(err)
	}

	if targetURL != path.Join(currentDir(), "./test/out/dir/d.txt") {
		t.Fatal("Unxz file Symlink error.")
	}

	dOutSymlinkFile, err := os.Open(path.Join(currentDir(), "./test/out/d.txt"))
	if err != nil {
		t.Fatal(err)
	}
	defer dOutSymlinkFile.Close()

	dOutSymlinkFileContent, err := ioutil.ReadAll(dOutSymlinkFile)
	if err != nil {
		t.Fatal(err)
	}

	dOutFile, err := os.Open(path.Join(currentDir(), "./test/out/dir/d.txt"))
	if err != nil {
		t.Fatal(err)
	}
	defer dOutFile.Close()

	dOutFileContent, err := ioutil.ReadAll(dOutFile)
	if err != nil {
		t.Fatal(err)
	}

	if string(dOutSymlinkFileContent) != string(dOutFileContent) {
		t.Fatal("The content of symlink file is not equal to that of target file")
	}
}

func testLintInDir(t *testing.T) {
	eOutFile, err := os.OpenFile(path.Join(currentDir(), "./test/out/dir/dir/e.txt"), os.O_RDWR|os.O_CREATE|os.O_APPEND, 0644)
	if err != nil {
		t.Fatal(err)
	}
	defer eOutFile.Close()

	eOutFileContent, err := ioutil.ReadAll(eOutFile)
	if err != nil {
		t.Fatal(err)
	}

	eOutLinkFile, err := os.Open(path.Join(currentDir(), "./test/out/dir/e.txt"))
	if err != nil {
		t.Fatal(err)
	}
	defer eOutLinkFile.Close()

	eOutLinkFileContent, err := ioutil.ReadAll(eOutLinkFile)
	if err != nil {
		t.Fatal(err)
	}

	if string(eOutFileContent) != string(eOutLinkFileContent) {
		t.Fatal("Unxz file content error or link file content error.")
	}

	t.Log("Old dir/dir/e.txt content:", string(eOutFileContent))
	t.Log("Change dir/dir/e.txt content.")
	_, err = eOutFile.WriteString("abcde")
	if err != nil {
		t.Fatal(err)
	}

	eOutFileNew, err := os.Open(path.Join(currentDir(), "./test/out/dir/dir/e.txt"))
	if err != nil {
		t.Fatal(err)
	}
	defer eOutFileNew.Close()

	eOutFileContentNew, err := ioutil.ReadAll(eOutFileNew)
	if err != nil {
		t.Fatal(err)
	}
	t.Log("New dir/dir/e.txt content:", string(eOutFileContentNew))

	eOutLinkFileNew, err := os.Open(path.Join(currentDir(), "./test/out/dir/e.txt"))
	if err != nil {
		t.Fatal(err)
	}
	defer eOutLinkFileNew.Close()

	eOutLinkFileContentNew, err := ioutil.ReadAll(eOutLinkFileNew)
	if err != nil {
		t.Fatal(err)
	}

	if string(eOutFileContentNew) != string(eOutLinkFileContentNew) {
		t.Fatal("Unxz file link error.")
	}
}

func testSymlinkInDir(t *testing.T) {
	inDirTargetURL, err := filepath.EvalSymlinks(path.Join(currentDir(), "./test/out/dir/dir/f.txt"))
	if err != nil {
		t.Fatal(err)
	}
	t.Log("F outSymlinkFile link:", inDirTargetURL)

	if inDirTargetURL != path.Join(currentDir(), "./test/out/dir/f.txt") {
		t.Fatal("Unxz file Symlink error.")
	}

	fOutSymlinkFile, err := os.Open(path.Join(currentDir(), "./test/out/dir/dir/f.txt"))
	if err != nil {
		t.Fatal(err)
	}
	defer fOutSymlinkFile.Close()

	fOutSymlinkFileContent, err := ioutil.ReadAll(fOutSymlinkFile)
	if err != nil {
		t.Fatal(err)
	}

	fOutFile, err := os.Open(path.Join(currentDir(), "./test/out/dir/f.txt"))
	if err != nil {
		t.Fatal(err)
	}
	defer fOutFile.Close()

	fOutFileContent, err := ioutil.ReadAll(fOutFile)
	if err != nil {
		t.Fatal(err)
	}

	if string(fOutSymlinkFileContent) != string(fOutFileContent) {
		t.Fatal("The content of symlink file is not equal to that of target file")
	}
}

func TestUnxz(t *testing.T) {
	filePath := filepath.FromSlash(path.Join(currentDir(), "./test/t.tar.xz"))
	outDir := filepath.FromSlash(path.Join(currentDir(), "./test/out") + "/")

	os.RemoveAll(outDir)

	unxz := New(filePath, outDir)
	err := unxz.Extract()
	if err != nil {
		t.Fatal(err)
	}

	t.Run("Content is null", testNullContent)
	t.Run("Content is equal", testContentEqual)
	t.Run("Link", testLink)
	t.Run("Symlink", testSymlink)
	t.Run("Link in dir", testLintInDir)
	t.Run("Symlink in dir", testSymlinkInDir)
}
