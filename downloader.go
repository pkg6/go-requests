package requests

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"runtime"
	"strconv"
	"sync"
)

var (
	ErrDownloaderFileIncomplete = errors.New("incomplete file")
	ErrDownloaderPartLength     = errors.New("the length of the segmented download is incorrect")
)

// filePart 文件分片
type downloaderPart struct {
	// 文件分片的序号
	Index int
	// 开始byte
	From int
	// 结束byte
	To int
	// http下载得到的文件分片内容
	Data []byte
}

type Downloader struct {
	req *Client
	// 下载源链接
	URL string
	// 下载完成文件名
	FileName string
	//协程数量
	CoroutineNumber int

	userAgent string
	ctx       context.Context
	// 待下载文件总大小
	fileSize int
	// 已完成文件切片
	doneFilePart []downloaderPart
}

func (c *Client) Download(ctx context.Context, uri, fileName string) error {
	return NewDownloader(c, ctx, uri, fileName).Download()
}

func NewDownloader(req *Client, ctx context.Context, url, fileName string, coroutineNumbers ...int) *Downloader {
	coroutineNumber := runtime.NumCPU()
	if len(coroutineNumbers) > 0 {
		coroutineNumber = coroutineNumbers[0]
	}
	return &Downloader{
		req:             req,
		userAgent:       RandomUserAgent(),
		ctx:             ctx,
		fileSize:        0,
		URL:             url,
		FileName:        fileName,
		CoroutineNumber: coroutineNumber,
		doneFilePart:    make([]downloaderPart, coroutineNumber),
	}
}

func (d *Downloader) Download() error {
	r, err := d.newReq(http.MethodHead)
	if err != nil {
		return err
	}
	resp, err := d.req.callRequest(r)
	if err != nil {
		return err
	}
	if resp.IsError() {
		return fmt.Errorf("incorrect response status code %v", resp.StatusCode)
	}
	if d.FileName == "" {
		d.FileName = filepath.Base(resp.Request.URL.Path)
	}
	// https://developer.mozilla.org/en-US/docs/Web/HTTP/Headers/Content-Length
	d.fileSize, err = strconv.Atoi(resp.Header.Get(HttpHeaderContentLength))
	if err != nil {
		return err
	}
	// 检查是否支持断点续传
	// https://developer.mozilla.org/en-US/docs/Web/HTTP/Headers/Accept-Ranges
	if resp.Header.Get(HttpHeaderAcceptRanges) != "bytes" {
		return d.simpleDownload()
	}
	return d.bytesDownload()
}

func (d *Downloader) newReq(method string) (req *http.Request, err error) {
	req, err = http.NewRequest(method, d.URL, nil)
	if err != nil {
		return
	}
	req.Header.Set(HttpHeaderUserAgent, d.userAgent)
	return
}

func (d *Downloader) simpleDownload() error {
	r, err := d.newReq(http.MethodGet)
	if err != nil {
		return err
	}
	resp, err := d.req.callRequest(r)
	if err != nil {
		return err
	}
	if resp.IsError() {
		return fmt.Errorf("incorrect response status code %v", resp.StatusCode)
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	file, err := d.openFile()
	if err != nil {
		return err
	}
	totalSize := len(body)
	_, _ = file.Write(body)
	if totalSize != d.fileSize {
		return ErrDownloaderFileIncomplete
	}
	return nil
}

func (d *Downloader) bytesDownload() error {
	jobs := make([]downloaderPart, d.CoroutineNumber)
	eachSize := d.fileSize / d.CoroutineNumber
	for i := range jobs {
		jobs[i].Index = i
		if i == 0 {
			jobs[i].From = 0
		} else {
			jobs[i].From = jobs[i-1].To + 1
		}
		if i < d.CoroutineNumber-1 {
			jobs[i].To = jobs[i].From + eachSize
		} else {
			// 最后一个filePart
			jobs[i].To = d.fileSize - 1
		}
	}
	var wg sync.WaitGroup
	for _, j := range jobs {
		wg.Add(1)
		go func(job downloaderPart) {
			defer wg.Done()
			if err := d.downloadPart(job); err != nil {
				d.req.Logger.Errorf("【Downloader】Blocked download failed %v job %v", err, job)
			}
		}(j)
	}
	wg.Wait()
	return d.mergeFileParts()
}

func (d *Downloader) openFile() (*os.File, error) {
	return os.Create(d.FileName)
}

func (d *Downloader) mergeFileParts() error {
	fullFileName, err := d.openFile()
	if err != nil {
		return err
	}
	defer fullFileName.Close()
	totalSize := 0
	for _, s := range d.doneFilePart {
		_, _ = fullFileName.Write(s.Data)
		totalSize += len(s.Data)
	}
	if totalSize != d.fileSize {
		return ErrDownloaderFileIncomplete
	}
	return nil
}

func (d *Downloader) downloadPart(c downloaderPart) error {
	r, err := d.newReq(http.MethodGet)
	if err != nil {
		return err
	}
	if d.ctx != nil {
		r.WithContext(d.ctx)
	}
	r.Header.Set("Range", fmt.Sprintf("bytes=%v-%v", c.From, c.To))
	resp, err := d.req.callRequest(r)
	if err != nil {
		return err
	}
	if resp.IsError() {
		return fmt.Errorf("incorrect response status code %v", resp.StatusCode)
	}
	defer resp.Body.Close()
	bs, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	if len(bs) != (c.To - c.From + 1) {
		return ErrDownloaderPartLength
	}
	c.Data = bs
	d.doneFilePart[c.Index] = c
	return nil
}
