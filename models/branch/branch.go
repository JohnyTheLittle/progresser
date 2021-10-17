package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type Branch struct {
	UpdatedAt    primitive.Timestamp `bson:"updated_at" json:"updated_at"`
	ID           string              `bson:"_id" json:"id"`
	Belongs      string              `bson:"user" json:"user"`
	Name         string              `bson:"name_of_branch" json:"name_of_branch"`
	Books        []Book              `bson:"books" json:"books"`
	Projects     []Project           `bson:"projects" json:"projects"`
	Improvement  []Improvement       `bson:"impovments" json:"impovments"`
	Hours        int                 `bson:"hours" json:"hours"`
	IsPrivate    bool                `bson:"is_private" json:"is_private"`
	VideoCourses []VideoCourse       `bson:"video_courses" json:"video_courses"`
	Queue        []QueueElement      `bson:"queue" json:"queue"`
}

type VideoCourse struct {
	Name string `bson:"video_course_name" json:"video_course_name"`
	Link string `bson:"video_courses_link" json:"video_courses_link"`
}
type Improvement struct {
	Name string `bson:"name_of_improvment" json:"name_of_improvment"`
}
type Project struct {
	Name string             `bson:"name_of_project"`
	Date primitive.DateTime `bson:"date_of_accomplished"`
}

type Book struct {
	Name         string `bson:"name_of_book" json:"name_of_book"`
	Accomplished bool   `bson:"is_accomplished" json:"is_accomplished"`
	TotalPages   int    `bson:"total_pages" json:"total_pages"`
	Pages        int    `bson:"pages" json:"pages"`
	Hours        int    `bson:"hours" json:"hours"`
}

type QueueElement struct {
	ElementName string `json:"element_name" bson:"element_name"`
}
