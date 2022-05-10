package service

import (
	"errors"
	"github.com/Moonlight-Zhao/go-project-example/repository"
	"time"
)

func PublishPost(parentId int64, content string) (int64, error) {
	return NewPostInfoFlow(parentId, content).Do()
}

func NewPostInfoFlow(parentId int64, content string) *PublishPostFlow {
	return &PublishPostFlow{
		parentId: parentId,
		content:  content,
	}
}

type PublishPostFlow struct {
	id       int64
	parentId int64
	content  string
}

func (f *PublishPostFlow) Do() (int64, error) {
	if err := f.checkParam(); err != nil {
		return 0, err
	}
	if err := f.publishPost(); err != nil {
		return 0, err
	}
	return f.id, nil
}

func (f *PublishPostFlow) checkParam() error {
	if f.parentId <= 0 {
		return errors.New("topic id must be larger than 0")
	}
	return nil
}

func (f *PublishPostFlow) publishPost() error {
	post := &repository.Post{
		Id:         f.id,
		ParentId:   f.parentId,
		Content:    f.content,
		CreateTime: time.Now().Unix(),
	}
	repository.LastId++
	post.Id = repository.LastId
	if err := repository.NewPostDaoInstance().AddPost(post); err != nil {
		return err
	}
	f.id = post.Id
	return nil
}
