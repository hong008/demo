package pipeline

import "demo/7.go_webcrawler/errors"

func genError(errMsg string) error {
	return errors.NewCrawlerError(errors.ERROR_TYPE_PIPELINE, errMsg)
}

func genParameterError(errMsg string) error {
	return errors.NewCrawlerErrorBy(errors.ERROR_TYPE_PIPELINE, errors.NewIllegalParameterError(errMsg))
}
