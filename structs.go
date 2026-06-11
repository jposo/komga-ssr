package main

import "time"

type Library struct {
	ID                string    `json:"id"`
	Name              string    `json:"name"`
	Root              string    `json:"root"`
	ImportComicInfo   bool      `json:"importComicInfo"`
	ImportEpubBook    bool      `json:"importEpubBook"`
	ImportEpubManga   bool      `json:"importEpubManga"`
	ImportMangaInfo   bool      `json:"importMangaInfo"`
	ImportCbx         bool      `json:"importCbx"`
	ScanForceModified bool      `json:"scanForceModified"`
	EmptyTrash        bool      `json:"emptyTrash"`
	Created           time.Time `json:"created"`
	LastModified      time.Time `json:"lastModified"`
}

type SeriesPage struct {
	Content          []Series `json:"content"`
	Empty            bool     `json:"empty"`
	First            bool     `json:"first"`
	Last             bool     `json:"last"`
	Number           int      `json:"number"`
	NumberOfElements int      `json:"numberOfElements"`
	Size             int      `json:"size"`
	TotalElements    int      `json:"totalElements"`
	TotalPages       int      `json:"totalPages"`
}

type Series struct {
	ID           string         `json:"id"`
	LibraryID    string         `json:"libraryId"`
	Name         string         `json:"name"`
	Url          string         `json:"url"`
	Created      time.Time      `json:"created"`
	LastModified time.Time      `json:"lastModified"`
	BooksCount   int            `json:"booksCount"`
	Metadata     SeriesMetadata `json:"metadata"`
}

type SeriesDetails struct {
	ID                   string         `json:"id"`
	LibraryID            string         `json:"libraryId"`
	Name                 string         `json:"name"`
	Url                  string         `json:"url"`
	Created              time.Time      `json:"created"`
	LastModified         time.Time      `json:"lastModified"`
	FileLastModified     time.Time      `json:"fileLastModified"`
	BooksCount           int            `json:"booksCount"`
	BooksReadCount       int            `json:"booksReadCount"`
	BooksUnreadCount     int            `json:"booksUnreadCount"`
	BooksInProgressCount int            `json:"booksInProgressCount"`
	Metadata             SeriesMetadata `json:"metadata"`
	Deleted              bool           `json:"deleted"`
	Oneshot              bool           `json:"oneshot"`
}

type SeriesMetadata struct {
	Title                string    `json:"title"`
	TitleLock            bool      `json:"titleLock"`
	TitleSort            string    `json:"titleSort"`
	TitleSortLock        bool      `json:"titleSortLock"`
	Summary              string    `json:"summary"`
	SummaryLock          bool      `json:"summaryLock"`
	Status               string    `json:"status"`
	StatusLock           bool      `json:"statusLock"`
	ReadingDirection     string    `json:"readingDirection,omitempty"`
	ReadingDirectionLock bool      `json:"readingDirectionLock"`
	Publisher            string    `json:"publisher"`
	PublisherLock        bool      `json:"publisherLock"`
	AgeRating            *int      `json:"ageRating"`
	AgeRatingLock        bool      `json:"ageRatingLock"`
	Language             string    `json:"language"`
	LanguageLock         bool      `json:"languageLock"`
	Genres               []string  `json:"genres"`
	GenresLock           bool      `json:"genresLock"`
	Tags                 []string  `json:"tags"`
	TagsLock             bool      `json:"tagsLock"`
	TotalBookCount       *int      `json:"totalBookCount"`
	TotalBookCountLock   bool      `json:"totalBookCountLock"`
	Created              time.Time `json:"created"`
	LastModified         time.Time `json:"lastModified"`
}

type BookPage struct {
	Content          []Book `json:"content"`
	Empty            bool   `json:"empty"`
	First            bool   `json:"first"`
	Last             bool   `json:"last"`
	Number           int    `json:"number"`
	NumberOfElements int    `json:"numberOfElements"`
	Size             int    `json:"size"`
	TotalElements    int    `json:"totalElements"`
	TotalPages       int    `json:"totalPages"`
}

type Book struct {
	ID           string        `json:"id"`
	SeriesID     string        `json:"seriesId"`
	SeriesTitle  string        `json:"seriesTitle"`
	LibraryID    string        `json:"libraryId"`
	Name         string        `json:"name"`
	Url          string        `json:"url"`
	Number       int           `json:"number"`
	Created      time.Time     `json:"created"`
	LastModified time.Time     `json:"lastModified"`
	SizeBytes    int64         `json:"sizeBytes"`
	Size         string        `json:"size"`
	Media        BookMedia     `json:"media"`
	Metadata     BookMetadata  `json:"metadata"`
	ReadProgress *ReadProgress `json:"readProgress,omitempty"`
	Deleted      bool          `json:"deleted"`
	FileHash     string        `json:"fileHash"`
	Oneshot      bool          `json:"oneshot"`
}

type BookMedia struct {
	Status               string `json:"status"`
	MediaType            string `json:"mediaType"`
	PagesCount           int    `json:"pagesCount"`
	Comment              string `json:"comment"`
	EpubDivinaCompatible bool   `json:"epubDivinaCompatible"`
	EpubIsKepub          bool   `json:"epubIsKepub"`
	MediaProfile         string `json:"mediaProfile"`
}

type BookMetadata struct {
	Title           string    `json:"title"`
	TitleLock       bool      `json:"titleLock"`
	Summary         string    `json:"summary"`
	SummaryLock     bool      `json:"summaryLock"`
	Number          string    `json:"number"`
	NumberLock      bool      `json:"numberLock"`
	NumberSort      float64   `json:"numberSort"`
	NumberSortLock  bool      `json:"numberSortLock"`
	ReleaseDate     *string   `json:"releaseDate"`
	ReleaseDateLock bool      `json:"releaseDateLock"`
	Authors         []Author  `json:"authors"`
	AuthorsLock     bool      `json:"authorsLock"`
	Tags            []string  `json:"tags"`
	TagsLock        bool      `json:"tagsLock"`
	Isbn            string    `json:"isbn"`
	IsbnLock        bool      `json:"isbnLock"`
	Links           []string  `json:"links"`
	LinksLock       bool      `json:"linksLock"`
	Created         time.Time `json:"created"`
	LastModified    time.Time `json:"lastModified"`
}

type Author struct {
	Name string `json:"name"`
	Role string `json:"role"`
}

type ReadProgress struct {
	Page         int64     `json:"page"`
	Completed    bool      `json:"completed"`
	ReadDate     time.Time `json:"readDate,omitempty"`
	Created      time.Time `json:"created,omitempty"`
	LastModified time.Time `json:"lastModified,omitempty"`
	DeviceID     string    `json:"deviceId,omitempty"`
	DeviceName   string    `json:"deviceName,omitempty"`
}

type ReadProgressUpdateRequest struct {
	Page      int64 `json:"page"`
	Completed bool  `json:"completed"`
}

type BookPageItem struct {
	Number    int    `json:"number"`
	FileName  string `json:"fileName"`
	MediaType string `json:"mediaType"`
	Width     int    `json:"width,omitempty"`
	Height    int    `json:"height,omitempty"`
}
