package storage

import (
	"context"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/utopiosphe/titan-storage-sdk/memfile"
)

var (
	apiKey     = os.Getenv("API_KEY")
	locatorURL = os.Getenv("LOCATOR_URL")
	token      = os.Getenv("TOKEN")
)

func TestCalculateCarCID(t *testing.T) {
	f, err := os.Open("./example/main.go")
	if err != nil {
		t.Fatal(err)
	}

	cid, err := CalculateCid(f)
	if err != nil {
		t.Fatal(err)
	}

	t.Log("cid ", cid.String())
}

func TestCreateCarWithFile(t *testing.T) {

	input := "xx.zip"
	output := "./xx.car"

	root, err := createCar(input, output)
	if err != nil {
		t.Fatal(err)
	}

	t.Log("root ", root.String())

}

func TestCreateCarWithStream(t *testing.T) {
	f, err := os.Open("./example/example.exe")
	if err != nil {
		t.Fatal(err)
	}
	defer f.Close()

	mFile := memfile.New([]byte{})
	root, err := createCarStream(f, mFile)
	if err != nil {
		t.Fatal(err)
	}

	t.Log("root ", root.String())
}

func TestUpload(t *testing.T) {
	storage, err := Initialize(&Config{TitanURL: locatorURL, APIKey: apiKey})
	if err != nil {
		t.Fatal("NewStorage error ", err)
	}

	progress := func(doneSize int64, totalSize int64) {
		t.Logf("upload %d of %d", doneSize, totalSize)
	}

	filePath := "./"
	visitFile := func(fp string, fi os.DirEntry, err error) error {
		// Check for and handle errors
		if err != nil {
			fmt.Println(err) // Can be used to handle errors (e.g., permission denied)
			return nil
		}
		if fi.IsDir() {
			return nil
		} else {
			// This is a file, you can perform file-specific operations here
			if strings.HasSuffix(fp, ".go") {
				path, err := filepath.Abs(fp)
				if err != nil {
					t.Fatal(err)
				}
				_, err = storage.UploadFilesWithPath(context.Background(), path, progress, true)
				if err != nil {
					t.Log("upload file failed ", err.Error())
					return nil
				}

				t.Logf("totalSize %s success", fp)
			}

		}
		return nil
	}

	err = filepath.WalkDir(filePath, visitFile)
	if err != nil {
		t.Fatal("WalkDir ", err)
	}
}

func TestUploadStream(t *testing.T) {
	storage, err := Initialize(&Config{TitanURL: locatorURL, APIKey: apiKey})
	if err != nil {
		t.Fatal("NewStorage error ", err)
	}

	progress := func(doneSize int64, totalSize int64) {
		t.Logf("upload %d of %d", doneSize, totalSize)
	}

	filePath := "./storage_test.go"
	f, err := os.Open(filePath)
	if err != nil {
		t.Fatal(err)
	}

	cid, err := storage.UploadStream(context.Background(), f, f.Name(), progress)
	if err != nil {
		t.Fatal("upload file failed ", err.Error())
	}

	t.Logf("totalSize %s success, cid %s", filePath, cid.String())

}

func TestGetFile(t *testing.T) {
	s, err := Initialize(&Config{TitanURL: locatorURL, APIKey: apiKey})
	if err != nil {
		t.Fatal("NewStorage error ", err)
	}

	// storageObject := s.(*storage)
	// t.Log("candidate node ", storageObject.candidateID)

	// progress := func(doneSize int64, totalSize int64) {
	// 	t.Logf("upload %d of %d", doneSize, totalSize)
	// }

	// filePath := "./storage_test.go"
	// f, err := os.Open(filePath)
	// if err != nil {
	// 	t.Fatal(err)
	// }

	// cid, err := s.UploadStream(context.Background(), f, f.Name(), progress)
	// if err != nil {
	// 	t.Fatal("upload file failed ", err.Error())
	// }

	// res, err := s.GetURL(context.Background(), cid.String())
	// if err != nil {
	// 	t.Fatal("get url ", err)
	// }

	// t.Log("url:", res.URLs)

	reader, fn, err := s.GetFileWithCid(context.Background(), "bafybeibhwmjstnkv3pvlswwcgs2as24bvlbr5tvba6mthtgf77mw4uown4")
	if err != nil {
		t.Fatal("get url ", err)
	}
	defer reader.Close()

	data, err := io.ReadAll(reader)

	if err != nil {
		t.Fatal("get url ", err)
	}

	newFilePath := fmt.Sprintf("./example/%s", fn)
	newFile, err := os.Create(newFilePath)
	if err != nil {
		t.Fatal("Create file", err)
	}
	defer newFile.Close()

	newFile.Write(data)

	t.Logf("write file %s %d", fn, len(data))
}

func TestUploadFileWithURL(t *testing.T) {
	s, err := Initialize(&Config{TitanURL: locatorURL, APIKey: apiKey})
	if err != nil {
		t.Fatal("NewStorage error ", err)
	}

	storageObject := s.(*storage)
	t.Log("candidate ", storageObject.candidateID)

	url := "https://files.oaiusercontent.com/file-JRBvuyBZDh7g6Garcgd9HQLl?se=2023-11-07T10%3A43%3A54Z&sp=r&sv=2021-08-06&sr=b&rscc=max-age%3D31536000%2C%20immutable&rscd=attachment%3B%20filename%3D27c485be-ad30-4b24-80ca-4d0e9ccfef8d.webp&sig=nzo0wiDe/a3oVAT5JPHLTP%2B7WdMouvVrmCyYkHSVUBE%3D"
	name, err := getFileNameFromURL(url)
	if err != nil {
		t.Fatal(err)
	}

	t.Log("name:", name)
	cid, newURL, err := s.UploadFileWithURL(context.Background(), url, nil)
	if err != nil {
		t.Fatal(err)
	}

	t.Logf("cid %s, newURL %s", cid, newURL)

}

func TestListAsset(t *testing.T) {
	s, err := Initialize(&Config{TitanURL: locatorURL, APIKey: apiKey})
	if err != nil {
		t.Fatal("NewStorage error ", err)
	}

	rsp, err := s.ListUserAssets(context.Background(), 0, 20, 1)
	if err != nil {
		t.Fatal("ListUserAssets ", err)
	}

	t.Logf("total assets %d, len:%d", rsp.Total, len(rsp.AssetOverviews))
	for _, asset := range rsp.AssetOverviews {
		t.Logf("cid:%s name:%s size:%d", asset.AssetRecord.CID, asset.UserAssetDetail.AssetName, asset.AssetRecord.TotalSize)
	}
}

func TestGetFileWithCid(t *testing.T) {
	s, err := Initialize(&Config{TitanURL: locatorURL, Token: token})
	if err != nil {
		t.Fatal("NewStorage error ", err)
	}
	r, fn, err := s.GetFileWithCid(context.Background(), "bafybeic7d7hf67yfgvv2wtdwth7uka2suqzr3cg3a52pt3bg3hhorsmp3m")
	if err != nil {
		t.Fatal("GetFileWithCid ", err)
	}
	defer r.Close()
	data, err := io.ReadAll(r)
	if err != nil {
		t.Fatal("GetFileWithCid ", err)
	}
	t.Logf("file:%s, size:%d", fn, len(data))
}

func TestUploadCarFetchBlock(t *testing.T) {
	carFilePath := "./example.car"
	// locatorURL := "https://storage-test-api.titannet.io"
	// token := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3MzM0ODQ0NjUsImlkIjoiMHg5ZGFjNjNhYjE5OTQzMjdmZjAyMjY4MDUyZmQwYThkODlmZjgxMzNjIiwib3JpZ19pYXQiOjE3MzMzOTgwNjUsInJvbGUiOjB9.XJsveS9BWHzHZG8q4tJd0rBU8m290YIAlf3tyG7VCFo"
	s, err := Initialize(&Config{TitanURL: locatorURL, Token: token, APIKey: apiKey})
	if err != nil {
		t.Fatal("NewStorage error ", err)
	}

	cid, err := s.UploadFilesWithPath(context.Background(), carFilePath, nil, false)
	if err != nil {
		t.Fatalf("UploadFilesWithPath %s", err)
	}

	list, err := s.ListAllBlocks(context.Background(), cid.String())
	if err != nil {
		t.Fatal("ListAllBlocks ", err)
	}

	var lastBlock string
	for _, block := range list {
		t.Log("block ", block)
		lastBlock = block
	}

	if lastBlock == "" {
		t.Log("no block")
	}

	reader, err := s.FetchBlockFromRoot(context.Background(), cid.String(), lastBlock)
	if err != nil {
		t.Fatal("FetchBlockFromRoot ", err)
	}

	content, err := io.ReadAll(reader)
	if err != nil {
		t.Fatal("FetchBlockFromRoot ", err)
	}

	t.Logf("rootCid: %s, subcid: %s", cid.String(), lastBlock)
	t.Logf("content size %d", len(content))
	t.Log("content ", string(content))
}
