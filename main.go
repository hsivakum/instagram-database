package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/google/uuid"
	_ "github.com/lib/pq"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"data-loader/models"
)

func init() {
	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname='%s' sslmode=disable", "localhost", "5432", "SYS", "instaadmin", "")
	var err error
	db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal(err)
	}
	rawDB, err = db.DB()
	if err != nil {
		log.Fatal(err)
	}
}

type Data struct {
	Account             string `json:"account,omitempty"`
	Biography           string `json:"biography,omitempty"`
	BusinessAddressJson struct {
		CityName      string  `json:"city_name,omitempty"`
		CityId        int64   `json:"city_id,omitempty"`
		Latitude      float64 `json:"latitude,omitempty"`
		Longitude     float64 `json:"longitude,omitempty"`
		StreetAddress string  `json:"street_address,omitempty"`
		ZipCode       string  `json:"zip_code,omitempty"`
	} `json:"business_address_json,omitempty"`
	BusinessCategoryName string `json:"business_category_name,omitempty"`
	BusinessEmail        string `json:"business_email,omitempty"`
	ExternalUrl          string `json:"external_url,omitempty"`
	Fbid                 string `json:"fbid,omitempty"`
	Followers            int64  `json:"followers,omitempty"`
	Following            int64  `json:"following,omitempty"`
	Highlights           []struct {
		Id    string `json:"id,omitempty"`
		Title string `json:"title,omitempty"`
		Image string `json:"image,omitempty"`
		Owner string `json:"owner,omitempty"`
	} `json:"highlights,omitempty"`
	Id                    string `json:"id,omitempty"`
	IsBusinessAccount     bool   `json:"is_business_account,omitempty"`
	IsProfessionalAccount bool   `json:"is_professional_account,omitempty"`
	IsVerified            bool   `json:"is_verified,omitempty"`
	Posts                 []struct {
		Caption  string `json:"caption,omitempty"`
		Likes    int64  `json:"likes,omitempty"`
		Datetime int    `json:"datetime,omitempty"`
		ImageUrl string `json:"image_url,omitempty"`
		Id       string `json:"id,omitempty"`
		Location *struct {
			Id            string `json:"id,omitempty"`
			HasPublicPage bool   `json:"has_public_page,omitempty"`
			Name          string `json:"name,omitempty"`
			Slug          string `json:"slug,omitempty"`
		} `json:"location,omitempty"`
		Url            string `json:"url,omitempty"`
		Comments       int64  `json:"comments,omitempty"`
		VideoViewCount int64  `json:"video_view_count,omitempty,omitempty"`
		VideoUrl       string `json:"video_url,omitempty,omitempty"`
	} `json:"posts,omitempty"`
	PostsCount       int64    `json:"posts_count,omitempty"`
	ProfileImageLink string   `json:"profile_image_link,omitempty"`
	ProfileName      string   `json:"profile_name,omitempty"`
	HighlightsCount  int64    `json:"highlights_count,omitempty"`
	CountryCode      string   `json:"country_code,omitempty"`
	Region           string   `json:"region,omitempty"`
	AvgEngagement    float64  `json:"avg_engagement,omitempty"`
	PostHashtags     []string `json:"post_hashtags,omitempty"`
}

var (
	db    *gorm.DB
	rawDB *sql.DB
)

func main() {

	file, err := os.ReadFile("instagram_profiles_Github Hashtag_dataset.json")
	if err != nil {
		log.Fatal(err)
	}

	var collection []Data
	err = json.Unmarshal(file, &collection)
	if err != nil {
		log.Fatal(err)
	}

	locations := []*models.Location{}
	mapLocations := map[string]*models.Location{}
	for _, data := range collection {
		for _, post := range data.Posts {
			if post.Location != nil {
				if _, ok := mapLocations[post.Location.Name]; !ok {
					elems := &models.Location{
						HasPublicPage: post.Location.HasPublicPage,
						Name:          post.Location.Name,
						Slug:          post.Location.Slug,
					}
					locations = append(locations, elems)
					mapLocations[post.Location.Name] = elems
				}
			}
		}
	}

	users := []*models.User{}
	businesss := []*models.Business{}
	posts := []*models.Post{}
	tags := []*models.HashTag{}

	postCaption := map[*int64]string{}
	hashTags := map[string]*int64{}
	allTags := []string{}
	highlights := []*models.Highlight{}
	highlightSet := map[string]*models.Highlight{}
	for _, data := range collection {
		userID := uuid.NewString()
		user := &models.User{
			ID:               userID,
			Username:         data.Account,
			FollowingCount:   data.Following,
			FollowersCount:   data.Followers,
			Bio:              data.Biography,
			PostsCount:       data.PostsCount,
			HighlightsCount:  data.HighlightsCount,
			Name:             data.ProfileName,
			ProfileImageLink: data.ProfileImageLink,
			IsBusiness:       data.IsBusinessAccount,
			IsVerified:       data.IsVerified,
			Country:          data.CountryCode,
			Region:           data.Region,
		}
		users = append(users, user)
		if data.IsBusinessAccount {
			atoi, _ := strconv.Atoi(data.BusinessAddressJson.ZipCode)
			business := &models.Business{
				ID:            uuid.NewString(),
				CityName:      data.BusinessAddressJson.CityName,
				Latitude:      data.BusinessAddressJson.Latitude,
				Longitude:     data.BusinessAddressJson.Longitude,
				StreetAddress: data.BusinessAddressJson.StreetAddress,
				ZipCode:       atoi,
				UserID:        userID,
			}
			businesss = append(businesss, business)
		}

		allTags = append(allTags, data.PostHashtags...)

		for _, postData := range data.Posts {
			elems := &models.Post{
				UserID:          userID,
				Caption:         postData.Caption,
				LikesCount:      postData.Likes,
				CommentsCount:   postData.Comments,
				VideoViewCount:  postData.VideoViewCount,
				PrimaryImageURL: postData.ImageUrl,
				PrimaryVideoURL: postData.VideoUrl,
				URL:             postData.Url,
			}
			if postData.Location != nil {
				elems.LocationID = &mapLocations[postData.Location.Name].ID
				elems.IsSponsored = postData.Location.Name == "Sponsered"
			}
			posts = append(posts, elems)
		}

		for _, highlight := range data.Highlights {
			highlightSet[fmt.Sprintf("%s_%s", userID, highlight.Title)] = &models.Highlight{
				UserID: userID,
				Title:  highlight.Title,
				Image:  highlight.Image,
			}
		}
	}

	for _, v := range highlightSet {
		highlights = append(highlights, v)
	}

	tagSet := map[string]bool{}
	for _, tag := range allTags {
		tagSet[tag] = true
	}

	allTags = []string{}
	for tag, _ := range tagSet {
		elems := &models.HashTag{
			Name: tag,
		}
		tags = append(tags, elems)
		hashTags[tag] = &elems.ID
		allTags = append(allTags, tag)
	}

	createUser(users)
	createBusiness(businesss)
	createLocations(locations)
	createPosts(posts)

	for _, post := range posts {
		postCaption[post.ID] = post.Caption
	}
	createHashTags(tags)

	postTags := []*models.PostTag{}
	createPostTagsConcurrently(postTags, allTags, postCaption, hashTags)

	createHighlights(highlights)

	createFollowers()

	createComments()

	createPostLikes()

	createCommentLikes()

	createPostImages()

	createStories()

	createStoryTags()

	createStoryViews()

	createHighlightStories()
}

func createHighlightStories() {
	type storyData struct {
		UserID  string          `json:"user_id"`
		Stories json.RawMessage `json:"stories"`
	}
	var storiesJson []storyData
	err := db.Table("stories as s").Select("json_agg(json_build_object('id', s.id, 'user_id', s.user_id , 'media_url', s.media_url, 'created_at', s.created_at, 'updated_at', s.updated_at, 'deleted_at', s.deleted_at )) AS stories", "s.user_id").Group("s.user_id").Scan(&storiesJson).Error
	if err != nil {
		log.Fatal(err)
	}

	storyByUser := map[string][]models.Story{}
	for _, story := range storiesJson {
		var userStories []models.Story
		err := json.Unmarshal(story.Stories, &userStories)
		if err != nil {
			log.Fatal(err)
		}
		storyByUser[story.UserID] = userStories
	}

	type highlightsData struct {
		UserID     string          `json:"user_id"`
		Highlights json.RawMessage `json:"highlights"`
	}
	var hData []highlightsData
	err = db.Table("highlights as h").Select("json_agg(h.id) as highlights", "h.user_id").Group("h.user_id").Scan(&hData).Error
	if err != nil {
		log.Fatal(err)
	}

	userHighlights := map[string][]int64{}
	for _, data := range hData {
		var ids []int64
		err := json.Unmarshal(data.Highlights, &ids)
		if err != nil {
			log.Fatal(err)
		}
		userHighlights[data.UserID] = ids
	}

	allHighlightStories := []models.HighlightsStory{}

	for userID, highLights := range userHighlights {
		stories := storyByUser[userID]

		// Calculate the maximum number of stories that can be added to highlights
		maxStories := len(stories)
		if maxStories < len(highLights) {
			maxStories = len(highLights)
		}

		// Randomly generate a number of stories to add to highlights
		numStories := rand.Intn(maxStories + 1)

		// Add the selected stories to highlights
		for i := 0; i < numStories && i < len(stories) && i < len(highLights); i++ {
			allHighlightStories = append(allHighlightStories, models.HighlightsStory{
				HighlightID: highLights[i],
				StoryID:     stories[i].ID,
			})
		}
	}

	err = db.CreateInBatches(allHighlightStories, 10000).Error
	if err != nil {
		log.Fatal(err)
	}

	log.Println("Highlights stories created!")
}

func createStoryViews() {
	var stories []models.Story
	err := db.Model(&models.Story{}).Scan(&stories).Error
	if err != nil {
		log.Fatal(err)
	}

	numbers := getRandomNumbers(int64(len(stories)), 300)

	type Followers struct {
		UserID    string          `json:"user_id"`
		Followers json.RawMessage `json:"followers"`
	}
	var userFollowers []Followers
	err = db.Table("users as u").Select("u.id AS user_id", "json_agg(follower_id) AS followers").Joins("LEFT JOIN followers f ON u.id = f.following_id").Group("u.id").Scan(&userFollowers).Error
	if err != nil {
		log.Fatal(err)
	}

	userXFollowers := map[string][]string{}
	for _, user := range userFollowers {
		var followers []string
		err := json.Unmarshal(user.Followers, &followers)
		if err != nil {
			log.Fatal(err)
		}
		userXFollowers[user.UserID] = followers
	}

	var storyViews []models.StoryView
	for i := 0; i < len(numbers); i++ {
		followers := userXFollowers[stories[i].UserID]
		for j, viewerID := range followers[:numbers[i]] {
			storyViews = append(storyViews, models.StoryView{
				StoryID:  stories[i].ID,
				ViewerID: viewerID,
				IsLiked: func() bool {
					return j%2 == 0
				}(),
			})
		}
	}

	err = db.CreateInBatches(storyViews, 10000).Error
	if err != nil {
		log.Fatal(err)
	}

	log.Println("successfully created story views!")
}

func createStoryTags() {
	var tags []models.HashTag
	err := db.Model(&models.HashTag{}).Scan(&tags).Error
	if err != nil {
		log.Fatal(err)
	}

	var stories []models.Story
	err = db.Model(&models.Story{}).Scan(&stories).Error
	if err != nil {
		log.Fatal(err)
	}

	numbers := getRandomNumbers(int64(len(stories)), 3)

	allStoryTags := []models.StoryTag{}
	storyTags := map[string]models.StoryTag{}
	for i := 0; i < len(numbers); i++ {
		number := numbers[i]
		randomNumbers := getRandomNumbers(int64(number), number)
		for _, tagIndex := range randomNumbers {
			storyTags[fmt.Sprintf("%s_%d", stories[i].ID, tags[tagIndex].ID)] = models.StoryTag{
				StoryID: stories[i].ID,
				TagID:   tags[tagIndex].ID,
			}
		}
	}

	for _, storyTag := range storyTags {
		allStoryTags = append(allStoryTags, storyTag)
	}

	err = db.CreateInBatches(allStoryTags, 10000).Error
	if err != nil {
		log.Fatal(err)
	}

	log.Println("Story tags created successfully!")
}

func createStories() {
	fileBytes, err := os.Open("stories.json")
	if err != nil {
		log.Fatal(err)
	}

	var storiesData []*models.Story
	err = json.NewDecoder(fileBytes).Decode(&storiesData)
	if err != nil {
		log.Fatal(err)
	}

	var users []*models.User
	err = db.Model(&models.User{}).Scan(&users).Error
	if err != nil {
		log.Fatal(err)
	}

	usersCount := int64(len(users))
	numbers := getRandomNumbers(usersCount, 100)

	allStories := []*models.Story{}
	start := int64(0)
	for i := int64(0); i < usersCount; i++ {
		storyCount := int64(numbers[i])
		if users[i].HighlightsCount > storyCount {
			storyCount = users[i].HighlightsCount
		}
		stories := storiesData[start : start+storyCount]
		for j, story := range stories {
			story.ID = uuid.NewString()
			story.UserID = users[i].ID
			stories[j] = story
		}

		allStories = append(allStories, stories...)
		start = start + storyCount
	}

	err = db.CreateInBatches(allStories, 10000).Error
	if err != nil {
		log.Fatal(err)
	}

	log.Println("Stories created successfully")
}

func createPostImages() {
	fileBytes, err := os.Open("post_images.json")
	if err != nil {
		log.Fatal(err)
	}

	var postImagesData []*models.PostImage
	err = json.NewDecoder(fileBytes).Decode(&postImagesData)
	if err != nil {
		log.Fatal(err)
	}

	var postCount int64
	err = db.Model(&models.Post{}).Select("count(*) as post_count").Scan(&postCount).Error
	if err != nil {
		log.Fatal(err)
	}

	var posts []*models.Post
	err = db.Model(&models.Post{}).Scan(&posts).Error
	if err != nil {
		log.Fatal(err)
	}

	var allPostImages []*models.PostImage
	start := 0
	postImagesCount := getRandomNumbers(postCount, 10)
	for i := int64(0); i < postCount; i++ {
		singlePostImageCount := postImagesCount[i]
		images := postImagesData[start : start+singlePostImageCount]
		for i, image := range images {
			image.PostOrder = i + 1
			image.PostID = *posts[i].ID
			images[i] = image
		}
		allPostImages = append(allPostImages, images...)
		start = start + singlePostImageCount
	}

	err = db.CreateInBatches(allPostImages, 9300).Error
	if err != nil {
		log.Fatal(err)
	}

	log.Println("Post images created")
}

func createCommentLikes() {
	type CommentSchema struct {
		ID              int64  `json:"id"`
		PostID          int64  `json:"post_id"`
		UserID          string `json:"user_id"`
		ParentCommentID int64  `json:"parent_comment_id"`
		PostAuthorID    string `json:"post_author_id"`
	}
	var comments []*CommentSchema
	err := db.Table("comments c").Select("c.id", "c.post_id", "c.user_id", "c.parent_comment_id", "p.user_id as post_author_id").Joins("inner join posts p on p.id = c.post_id").Scan(&comments).Error
	if err != nil {
		log.Fatal(err)
	}

	type Followers struct {
		UserID    string          `json:"user_id"`
		Followers json.RawMessage `json:"followers"`
	}
	var userFollowers []Followers
	err = db.Table("users as u").Select("u.id AS user_id", "json_agg(follower_id) AS followers").Joins("LEFT JOIN followers f ON u.id = f.following_id").Group("u.id").Scan(&userFollowers).Error
	if err != nil {
		log.Fatal(err)
	}

	userXFollowers := map[string][]string{}
	for _, user := range userFollowers {
		var followers []string
		err := json.Unmarshal(user.Followers, &followers)
		if err != nil {
			log.Fatal(err)
		}
		userXFollowers[user.UserID] = followers
	}

	randomNumbers := getRandomNumbers(int64(len(comments)), 200)

	commentLikes := []*models.CommentLike{}
	for i, comment := range comments {
		noOfLikes := randomNumbers[i]
		for j := 0; j < noOfLikes; j++ {
			commentLikes = append(commentLikes, &models.CommentLike{
				CommentID: comment.ID,
				LikedBy:   userXFollowers[comment.PostAuthorID][j],
			})
		}
	}

	err = db.CreateInBatches(commentLikes, 10000).Error
	if err != nil {
		log.Fatal(err)
	}

	log.Println("Comment likes generated")
}

func getRandomNumbers(size int64, maxValue int) []int {
	rand.Seed(time.Now().UnixNano())

	// Create an array to store the random numbers
	randomNumbers := make([]int, size)

	// Generate random numbers and store them in the array
	for i := int64(0); i < size; i++ {
		randomNumbers[i] = rand.Intn(maxValue) + 1 // Generates a random number between 1 and 200
	}

	return randomNumbers
}

func createPostLikes() {
	var posts []*models.Post
	err := db.Model(&models.Post{}).Select("id", "likes_count", "user_id").Where("likes_count > 0").Scan(&posts).Error
	if err != nil {
		log.Fatal(err)
	}

	var likes []*models.PostLikes
	for _, post := range posts {
		var followingUsers []models.Follower
		err := db.Model(&models.Follower{}).Select("follower_id").Where("following_id = ?", post.UserID).Scan(&followingUsers).Error
		if err != nil {
			log.Fatal(err)
		}

		selectedUsers := []models.Follower{}
		if int64(len(followingUsers)) > post.LikesCount {
			selectedUsers = followingUsers[:post.LikesCount]
		}

		for _, user := range selectedUsers {
			likes = append(likes, &models.PostLikes{
				PostID: *post.ID,
				UserID: user.FollowerID,
			})
		}
	}

	err = db.CreateInBatches(likes, 10000).Error
	if err != nil {
		log.Fatal(err)
	}

	_, err = rawDB.Exec(`UPDATE posts AS p
	SET likes_count = (
		SELECT COUNT(*)
		FROM post_likes AS pl
		WHERE pl.post_id = p.id
	)`)
	if err != nil {
		log.Fatal(err)
	}

	log.Println("Post likes successfully created")
}

func createFollowers() {
	// Query for user IDs and their follower and following counts
	rows, err := rawDB.Query("SELECT id, following_count, followers_count FROM users")
	if err != nil {
		log.Fatal(err)
	}

	defer rows.Close()

	type User struct {
		UserID         string
		FollowingCount int
		FollowersCount int
	}

	users := []User{}
	followers := []*models.Follower{}

	// Iterate through each user and randomly generate follower relationships
	for rows.Next() {
		var user User
		if err := rows.Scan(&user.UserID, &user.FollowingCount, &user.FollowersCount); err != nil {
			log.Fatal(err)
		}
		users = append(users, user)
	}

	if err := rows.Err(); err != nil {
		log.Fatal(err)
	}

	followersMap := map[string]bool{}

	for _, user := range users {
		// Generate follower relationships based on following and followers count
		followingIDs := generateRandomUserIDs(user.UserID, user.FollowingCount, rawDB, false)
		for _, id := range followingIDs {
			key := fmt.Sprintf("%s_%s", user.UserID, id)
			if _, ok := followersMap[key]; !ok {
				followersMap[key] = true
				follower := &models.Follower{
					FollowerID:  user.UserID,
					FollowingID: id,
				}
				followers = append(followers, follower)
			}
		}

		followersIDs := generateRandomUserIDs(user.UserID, user.FollowersCount, rawDB, false)
		for _, id := range followersIDs {
			key := fmt.Sprintf("%s_%s", id, user.UserID)
			if _, ok := followersMap[key]; !ok {
				followersMap[key] = true
				follower := &models.Follower{
					FollowerID:  id,
					FollowingID: user.UserID,
				}
				followers = append(followers, follower)
			}
		}
	}

	tx := db.CreateInBatches(followers, 10000)
	if tx.Error != nil {
		log.Fatal(tx.Error)
	}

	tx = db.Raw(`UPDATE users AS u
SET
    following_count = (
        SELECT COUNT(*)
        FROM followers AS f
        WHERE f.follower_id = u.id
    ),
    followers_count = (
        SELECT COUNT(*)
        FROM followers AS f
        WHERE f.following_id = u.id
    )`)

	if tx.Error != nil {
		log.Fatal(tx.Error)
	}

	fmt.Println("Follower relationships have been generated successfully!")
}

func createComments() {
	var commentsStore []*models.Comment
	file, err := os.Open("comments.json")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	err = json.NewDecoder(file).Decode(&commentsStore)
	if err != nil {
		log.Fatal(err)
	}

	var totalComments int
	err = db.Model(&models.Post{}).Select("sum(comments_count) as total_comments").Scan(&totalComments).Error
	if err != nil {
		log.Fatal(err)
	}

	requiredComments := commentsStore[:totalComments]

	var posts []models.Post
	tx := db.Model(&models.Post{}).Select("comments_count", "id", "user_id").Where("comments_count > 0").Scan(&posts)
	if tx.Error != nil {
		log.Fatal(tx.Error)
	}

	finalComments := []*models.Comment{}

	postComments := map[int64][]*models.Comment{}
	start := int64(0)
	counter := int64(1)
	for _, post := range posts {
		var followingUsers []models.Follower
		err := db.Model(&models.Follower{}).Select("follower_id").Where("following_id = ?", post.UserID).Scan(&followingUsers).Error
		if err != nil {
			log.Fatal(err)
		}
		selectedComments := requiredComments[start : start+post.CommentsCount]
		for i, comment := range selectedComments {
			rand.Seed(time.Now().UnixNano())
			index := rand.Intn(len(followingUsers))
			comment.UserID = followingUsers[index].FollowerID
			comment.PostID = *post.ID
			selectedComments[i] = comment
			comment.ID = counter
			counter = counter + 1
		}

		postComments[*post.ID] = selectedComments
		finalComments = append(finalComments, selectedComments...)
		start = start + post.CommentsCount
	}

	for _, v := range postComments {
		fillParentCommentID(v)
	}

	tx = db.CreateInBatches(finalComments, 8190)
	if tx.Error != nil {
		log.Fatal(tx.Error)
	}

	log.Println("Comments successfully loaded")
}

// Function to generate random user IDs based on following or followers
func generateRandomUserIDs(excludeID string, count int, db *sql.DB, includeSelf bool) []string {
	var userIDs []string
	var query string
	// Exclude the current user from the random selection
	query = fmt.Sprintf("SELECT id FROM users WHERE id != $1 ORDER BY random() LIMIT %d", count)
	// Query for random user IDs based on the condition
	rows, err := db.Query(query, excludeID)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	for rows.Next() {
		var userID string
		if err := rows.Scan(&userID); err != nil {
			log.Fatal(err)
		}
		userIDs = append(userIDs, userID)
	}
	if err := rows.Err(); err != nil {
		log.Fatal(err)
	}

	return userIDs
}

func createUser(users []*models.User) {
	tx := db.CreateInBatches(users, 10000)
	if tx.Error != nil {
		log.Fatal(tx.Error)
	}
}

func createBusiness(businesses []*models.Business) {
	tx := db.CreateInBatches(businesses, 10000)
	if tx.Error != nil {
		log.Fatal(tx.Error)
	}
}

func createLocations(locations []*models.Location) {
	tx := db.CreateInBatches(locations, 10000)
	if tx.Error != nil {
		log.Fatal(tx.Error)
	}
}

func createPosts(posts []*models.Post) {
	tx := db.CreateInBatches(posts, 4000)
	if tx.Error != nil {
		log.Fatal(tx.Error)
	}
}

func createHashTags(hashTags []*models.HashTag) {
	tx := db.CreateInBatches(hashTags, 10000)
	if tx.Error != nil {
		log.Fatal(tx.Error)
	}
}

func createHighlights(highlights []*models.Highlight) {
	tx := db.CreateInBatches(highlights, 10000)
	if tx.Error != nil {
		log.Fatal(tx.Error)
	}
}

func createPostTags(postTags []*models.PostTag) {
	tx := db.CreateInBatches(postTags, 10000)
	if tx.Error != nil {
		log.Fatal(tx.Error)
	}
}

func createPostTagsConcurrently(postTags []*models.PostTag, allTags []string, postCaption map[*int64]string, hashTags map[string]*int64) {
	resultCh := make(chan *models.PostTag)

	var wg sync.WaitGroup
	wg.Add(len(allTags))

	// Start worker goroutines
	for _, tag := range allTags {
		go func(tag string) {
			defer wg.Done()
			worker(tag, postCaption, hashTags, resultCh)
		}(tag)
	}

	// Close result channel when all workers finish
	go func() {
		wg.Wait()
		close(resultCh)
	}()

	// Collect results from worker goroutines
	for postTag := range resultCh {
		postTags = append(postTags, postTag)
	}

	createPostTags(postTags)
}

// Worker function to process tags concurrently
func worker(tag string, postCaption map[*int64]string, hashTags map[string]*int64, resultCh chan<- *models.PostTag) {
	for postID, caption := range postCaption {
		if strings.Contains(caption, tag) {
			resultCh <- &models.PostTag{
				PostID: *postID,
				TagID:  *hashTags[tag],
			}
		}
	}
}

// Function to fill in parent_comment_id for comments based on a binary tree structure
func fillParentCommentID(comments []*models.Comment) {
	// Shuffle comment IDs
	rand.Shuffle(len(comments), func(i, j int) {
		comments[i], comments[j] = comments[j], comments[i]
	})

	// Build parent_comment_id based on shuffled IDs
	for i := 1; i < len(comments); i++ {
		comments[i].ParentCommentID = comments[(i-1)/2].ID
	}
}
