package core

import (
	"archive/zip"
	"bytes"
	"context"
	"io"
	"io/ioutil"
	"strings"
	"sync"
	"unicode/utf8"

	"gitlab.badanamu.com.cn/calmisland/common-log/log"
	"gitlab.badanamu.com.cn/calmisland/kidsloop-file-processing-service/entity"
)

const (
	OOXMLStartTagCreator = `<dc:creator>`
	OOXMLEndTagCreator   = `</dc:creator>`

	OOXMLStartTagLastModifiedBy = `<cp:lastModifiedBy>`
	OOXMLEndTagLastModifiedBy   = `</cp:lastModifiedBy>`

	OOXMLStartTagDescription = `<dc:description>`
	OOXMLEndTagDescription   = `</dc:description>`

	OOXMLStartTagKeywords = `<cp:keywords>`
	OOXMLEndTagKeywords   = `</cp:keywords>`

	OOXMLStartTagCategory = `<cp:category>`
	OOXMLEndTagCategory   = `</cp:category>`

	OOXMLStartTagManager = `<Manager>`
	OOXMLEndTagManager   = `</Manager>`

	OOXMLStartTagCompany = `<Company>`
	OOXMLEndTagCompany   = `</Company>`
)

type RemoveOOXMLMetaDataHandler struct {
}

func (h *RemoveOOXMLMetaDataHandler) Do(ctx context.Context, f *entity.HandleFileParams) error {
	dst, err := f.CreateOutputFile(ctx)
	if err != nil {
		log.Error(ctx, "Can't create output file",
			log.Err(err),
			log.Any("params", f))
		return err
	}

	r, err := zip.OpenReader(f.LocalPath)
	if err != nil {
		log.Error(ctx, "Can't open zip file",
			log.Err(err),
			log.String("path", f.LocalPath))
		return err
	}
	defer r.Close()

	buf := new(bytes.Buffer)

	// Create a new zip archive.
	zipWriter := zip.NewWriter(buf)
	defer zipWriter.Close()

	for _, f := range r.File {
		// Write same name file to output zip file
		zipFile, err := zipWriter.Create(f.Name)
		if err != nil {
			log.Error(ctx, "Can't create zip file",
				log.Err(err),
				log.String("name", f.Name))
			return err
		}

		rc, err := f.Open()
		if err != nil {
			log.Error(ctx, "Can't open zip file",
				log.Err(err),
				log.Any("file", f))
			return err
		}
		defer rc.Close()

		if f.Name == "docProps/core.xml" {
			data, err := ioutil.ReadAll(rc)
			if err != nil {
				log.Error(ctx, "Can't read zip file",
					log.Err(err),
					log.Any("file", f))
				return err
			}

			xmlContent := string(data)

			xmlContent, err = h.removeOOXMLCoreMetadata(ctx, xmlContent)
			if err != nil {
				log.Error(ctx, "Can't remove core metadata",
					log.Err(err),
					log.Any("xmlContent", xmlContent))
				return err
			}

			_, err = zipFile.Write([]byte(xmlContent))
			if err != nil {
				log.Error(ctx, "Can't write data to dst zip file",
					log.Err(err),
					log.Any("xmlContent", xmlContent))
				return err
			}
		} else if f.Name == "docProps/app.xml" {
			data, err := ioutil.ReadAll(rc)
			if err != nil {
				log.Error(ctx, "Can't read zip file",
					log.Err(err),
					log.Any("file", f))
				return err
			}

			xmlContent := string(data)

			xmlContent, err = h.removeOOXMLAppMetadata(ctx, xmlContent)
			if err != nil {
				log.Error(ctx, "Can't remove core metadata",
					log.Err(err),
					log.Any("xmlContent", xmlContent))
				return err
			}

			_, err = zipFile.Write([]byte(xmlContent))
			if err != nil {
				log.Error(ctx, "Can't write data to dst zip file",
					log.Err(err),
					log.Any("xmlContent", xmlContent))
				return err
			}
		} else {
			_, err = io.Copy(zipFile, rc)
			if err != nil {
				log.Error(ctx, "Can't copy zip file",
					log.Err(err),
					log.Any("zipFile", zipFile))
				return err
			}
		}
	}

	_, err = buf.WriteTo(dst)
	if err != nil {
		if err != nil {
			log.Error(ctx, "Can't write to dst file",
				log.Err(err),
				log.Any("dstFile", dst))
			return err
		}
	}

	return nil
}

func (h *RemoveOOXMLMetaDataHandler) removeOOXMLCoreMetadata(ctx context.Context, xmlContent string) (string, error) {
	xmlContent, err := h.removeOOXMLTag(ctx, xmlContent, OOXMLStartTagCreator, OOXMLEndTagCreator)
	if err != nil {
		log.Error(ctx, "Can't remove creator metadata",
			log.Err(err),
			log.String("xmlContent", xmlContent),
			log.String("startTag", OOXMLStartTagCreator),
			log.String("endTag", OOXMLEndTagCreator))
		return xmlContent, err
	}

	xmlContent, err = h.removeOOXMLTag(ctx, xmlContent, OOXMLStartTagKeywords, OOXMLEndTagKeywords)
	if err != nil {
		log.Error(ctx, "Can't remove keywords metadata",
			log.Err(err),
			log.String("xmlContent", xmlContent),
			log.String("startTag", OOXMLStartTagKeywords),
			log.String("endTag", OOXMLEndTagKeywords))
		return xmlContent, err
	}

	xmlContent, err = h.removeOOXMLTag(ctx, xmlContent, OOXMLStartTagDescription, OOXMLEndTagDescription)
	if err != nil {
		log.Error(ctx, "Can't remove description metadata",
			log.Err(err),
			log.String("xmlContent", xmlContent),
			log.String("startTag", OOXMLStartTagDescription),
			log.String("endTag", OOXMLEndTagDescription))
		return xmlContent, err
	}

	xmlContent, err = h.removeOOXMLTag(ctx, xmlContent, OOXMLStartTagLastModifiedBy, OOXMLEndTagLastModifiedBy)
	if err != nil {
		log.Error(ctx, "Can't remove lastModifiedBy metadata",
			log.Err(err),
			log.String("xmlContent", xmlContent),
			log.String("startTag", OOXMLStartTagLastModifiedBy),
			log.String("endTag", OOXMLEndTagLastModifiedBy))
		return xmlContent, err
	}

	xmlContent, err = h.removeOOXMLTag(ctx, xmlContent, OOXMLStartTagCategory, OOXMLEndTagCategory)
	if err != nil {
		log.Error(ctx, "Can't remove category metadata",
			log.Err(err),
			log.String("xmlContent", xmlContent),
			log.String("startTag", OOXMLStartTagCategory),
			log.String("endTag", OOXMLEndTagCategory))
		return xmlContent, err
	}

	return xmlContent, nil
}

func (h *RemoveOOXMLMetaDataHandler) removeOOXMLAppMetadata(ctx context.Context, xmlContent string) (string, error) {
	xmlContent, err := h.removeOOXMLTag(ctx, xmlContent, OOXMLStartTagCompany, OOXMLEndTagCompany)
	if err != nil {
		log.Error(ctx, "Can't remove company metadata",
			log.Err(err),
			log.String("xmlContent", xmlContent),
			log.String("startTag", OOXMLStartTagCompany),
			log.String("endTag", OOXMLEndTagCompany))
		return xmlContent, err
	}

	xmlContent, err = h.removeOOXMLTag(ctx, xmlContent, OOXMLStartTagManager, OOXMLEndTagManager)
	if err != nil {
		log.Error(ctx, "Can't remove manager metadata",
			log.Err(err),
			log.String("xmlContent", xmlContent),
			log.String("startTag", OOXMLStartTagManager),
			log.String("endTag", OOXMLEndTagManager))
		return xmlContent, err
	}

	return xmlContent, nil
}

func (h *RemoveOOXMLMetaDataHandler) removeOOXMLTag(ctx context.Context, xmlContent string, startTag, endTag string) (string, error) {
	startIndex := strings.Index(xmlContent, startTag)
	if startIndex == -1 {
		log.Debug(ctx, "startTag not exist",
			log.String("startTag", startTag),
			log.String("xmlContent", xmlContent))
		return xmlContent, nil
	}
	endIndex := strings.Index(xmlContent, endTag)
	if startIndex == -1 {
		log.Debug(ctx, "endTag not exist",
			log.String("endTag", endTag),
			log.String("xmlContent", xmlContent))
		return xmlContent, nil
	}

	newContent := xmlContent[:startIndex] + xmlContent[endIndex+utf8.RuneCountInString(endTag):]

	return newContent, nil
}

var (
	_RemoveOOXMLMetaDataHandler     IFileHandler
	_RemoveOOXMLMetaDataHandlerOnce sync.Once
)

func GetRemoveOOXMLMetaDataHandler() IFileHandler {
	_RemoveOOXMLMetaDataHandlerOnce.Do(func() {
		_RemoveOOXMLMetaDataHandler = new(RemoveOOXMLMetaDataHandler)
	})
	return _RemoveOOXMLMetaDataHandler
}
