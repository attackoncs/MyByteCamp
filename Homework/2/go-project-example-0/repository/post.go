package repository

import (
	"bufio"
	"encoding/json"
	"errors"
	"os"
	"sync"
)

//回帖列表，通过parent_id关联到话题
type Post struct {
	Id         int64  `json:"id"`
	ParentId   int64  `json:"parent_id"`
	Content    string `json:"content"`
	CreateTime int64  `json:"create_time"`
}
type PostDao struct {
}

var (
	postDao  *PostDao
	postOnce sync.Once
)

func NewPostDaoInstance() *PostDao {
	postOnce.Do(
		func() {
			postDao = &PostDao{}
		})
	return postDao
}
func (*PostDao) QueryPostsByParentId(parentId int64) []*Post {
	return postIndexMap[parentId]
}

func (d *PostDao) AddPost(post *Post) error {
	lock := sync.Mutex{}
	lock.Lock()
	posts, ok := postIndexMap[post.ParentId]
	if !ok {
		return errors.New("parentId not exist")
	}

	postIndexMap[post.ParentId] = append(posts, post)

	err := d.InsertPost2Data("./data/", post)
	if err != nil {
		return err
	}
	lock.Unlock()
	return nil
}

func (d *PostDao) InsertPost2Data(path string, post *Post) error {
	file, err := os.OpenFile(path+"post", os.O_WRONLY|os.O_APPEND, 777)
	defer file.Close()
	if err != nil {
		return err
	}
	writer := bufio.NewWriter(file)
	data, err := json.Marshal(*post)
	if err != nil {
		return err
	}
	writer.Write(data)
	writer.WriteString("\n")
	writer.Flush()
	return nil
}
