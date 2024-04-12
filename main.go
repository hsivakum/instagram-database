package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"sync"

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
