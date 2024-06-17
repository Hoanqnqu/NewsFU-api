package service

import (
	"github.com/google/uuid"
	inport "news-api/application/port/in"
	outport "news-api/application/port/out"
)

type CommentService struct {
	commentPort outport.CommentPort
}

func NewCommentService(commentPort outport.CommentPort) *CommentService {
	return &CommentService{commentPort: commentPort}

}

func (c *CommentService) InsertComment(newsID, userID uuid.UUID, comment string) error {
	return c.commentPort.InsertComment(newsID, userID, comment)
}

func (c *CommentService) GetCommentsByNewsID(newsID uuid.UUID) ([]*inport.Comment, error) {
	comments, err := c.commentPort.GetCommentsByNews(newsID)
	if err != nil {
		return nil, err
	}
	return func() []*inport.Comment {
		rs := make([]*inport.Comment, len(comments))
		for i, v := range comments {
			rs[i] = &inport.Comment{
				Text:        v.Text,
				PublishedAt: v.PublishedAt,
				UserName:    v.UserName,
				UserAvatar:  v.UserAvatar,
			}
		}
		return rs
	}(), nil
}
