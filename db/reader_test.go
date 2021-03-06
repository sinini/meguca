package db

import (
	"github.com/bakape/meguca/config"
	"github.com/bakape/meguca/types"
	r "github.com/dancannon/gorethink"
	. "gopkg.in/check.v1"
)

var (
	genericImage = &types.Image{
		ImageCommon: types.ImageCommon{
			SHA1: "foo",
		},
	}

	sampleThreads = []types.DatabaseThread{
		{
			ID:       1,
			Board:    "a",
			ImageCtr: 1,
			PostCtr:  2,
			Posts: map[string]types.Post{
				"1": {
					Board: "a",
					ID:    1,
					Image: genericImage,
				},
				"2": {
					Board: "a",
					ID:    2,
					Image: genericImage,
				},
				"3": {
					Board: "a",
					ID:    3,
				},
			},
			Log: [][]byte{
				{1, 3, 4},
				{1, 3, 2},
				{1},
			},
		},
		{
			ID:    4,
			Board: "a",
			Posts: map[string]types.Post{
				"4": {
					Board: "a",
					ID:    4,
					Image: genericImage,
				},
			},
		},
		{
			ID:    5,
			Board: "c",
			Posts: map[string]types.Post{
				"5": {
					Board: "c",
					ID:    5,
					Image: genericImage,
				},
			},
		},
	}

	boardStandard = types.Board{
		Ctr: 7,
		Threads: []types.Thread{
			{
				Post: types.Post{
					ID:    4,
					Board: "a",
					Image: genericImage,
				},
				Posts: nil,
			},
			{
				ImageCtr: 1,
				PostCtr:  2,
				LogCtr:   3,
				Post: types.Post{
					ID:    1,
					Board: "a",
					Image: genericImage,
				},
				Posts: nil,
			},
		},
	}
)

func (*DBSuite) TestGetPost(c *C) {
	conf := config.ServerConfigs{}
	conf.Boards.Enabled = []string{"a"}
	config.Set(conf)
	std := types.Post{
		ID:    2,
		Board: "a",
	}
	threads := []types.DatabaseThread{
		{
			ID:    1,
			Board: "a",
			Posts: map[string]types.Post{
				"2": std,
			},
		},
		{
			ID:    4,
			Board: "q",
			Posts: map[string]types.Post{
				"5": {
					Board: "q",
					ID:    5,
				},
			},
		},
	}
	c.Assert(Write(r.Table("threads").Insert(threads)), IsNil)

	empties := [...]struct {
		id, op int64
	}{
		{2, 76}, // Thread does not exist
		{8, 1},  // Post does not exist
	}

	for _, args := range empties {
		post, err := GetPost(args.id, args.op)
		c.Assert(err, Equals, r.ErrEmptyResult)
		c.Assert(post, DeepEquals, types.Post{})
	}

	// Valid read
	post, err := GetPost(2, 1)
	c.Assert(err, IsNil)
	c.Assert(post, DeepEquals, std)
}

func (*DBSuite) TestGetBoard(c *C) {
	setEnabledBoards("a")
	c.Assert(Write(r.Table("threads").Insert(sampleThreads)), IsNil)

	boardCounters := map[string]interface{}{
		"id": "histCounts",
		"a":  7,
	}
	c.Assert(Write(r.Table("main").Insert(boardCounters)), IsNil)

	board, err := GetBoard("a")
	c.Assert(err, IsNil)
	c.Assert(board, DeepEquals, &boardStandard)
}

func setEnabledBoards(boards ...string) {
	conf := config.ServerConfigs{}
	conf.Boards.Enabled = boards
	config.Set(conf)
}

func (*DBSuite) TestGetEmptyBoard(c *C) {
	setEnabledBoards("a")
	boardCounters := Document{"histCounts"}
	c.Assert(Write(r.Table("main").Insert(boardCounters)), IsNil)

	board, err := GetBoard("a")
	c.Assert(err, IsNil)
	c.Assert(board, DeepEquals, new(types.Board))
}

func (*DBSuite) TestGetAllBoard(c *C) {
	setEnabledBoards("a")
	c.Assert(Write(r.Table("threads").Insert(sampleThreads)), IsNil)
	info := infoDocument{
		Document: Document{"info"},
		PostCtr:  7,
	}
	c.Assert(Write(r.Table("main").Insert(info)), IsNil)

	std := boardStandard
	std.Threads = []types.Thread{
		boardStandard.Threads[0],
		{
			Post: types.Post{
				ID:    5,
				Board: "c",
				Image: genericImage,
			},
			Posts: nil,
		},
		boardStandard.Threads[1],
	}

	board, err := GetAllBoard()
	c.Assert(err, IsNil)
	c.Assert(board, DeepEquals, &std)
}

func (*DBSuite) TestGetEmptyAllBoard(c *C) {
	setEnabledBoards("a")
	info := infoDocument{
		Document: Document{"info"},
	}
	c.Assert(Write(r.Table("main").Insert(info)), IsNil)
	board, err := GetAllBoard()
	c.Assert(err, IsNil)
	c.Assert(board, DeepEquals, new(types.Board))
}

func (*DBSuite) TestReaderGetThread(c *C) {
	conf := config.ServerConfigs{}
	conf.Boards.Enabled = []string{"a"}
	config.Set(conf)
	c.Assert(Write(r.Table("threads").Insert(sampleThreads)), IsNil)

	// No replies ;_;
	std := &types.Thread{
		Post: types.Post{
			Board: "a",
			ID:    4,
			Image: genericImage,
		},
		Posts: map[string]types.Post{},
	}
	thread, err := GetThread(4, 0)
	c.Assert(err, IsNil)
	c.Assert(thread, DeepEquals, std)

	// With replies
	std = &types.Thread{
		ImageCtr: 1,
		PostCtr:  2,
		LogCtr:   3,
		Post: types.Post{
			Board: "a",
			ID:    1,
			Image: genericImage,
		},
		Posts: map[string]types.Post{
			"2": {
				Board: "a",
				ID:    2,
				Image: genericImage,
			},
			"3": {
				Board: "a",
				ID:    3,
			},
		},
	}
	thread, err = GetThread(1, 0)
	c.Assert(err, IsNil)
	c.Assert(thread, DeepEquals, std)

	// Last 1 post
	delete(std.Posts, "2")
	thread, err = GetThread(1, 1)
	c.Assert(err, IsNil)
	c.Assert(thread, DeepEquals, std)
}
