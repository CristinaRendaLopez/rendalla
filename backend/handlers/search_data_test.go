package handlers_test

import "github.com/CristinaRendaLopez/rendalla-backend/models"

// --- SONGS ---

var SongLoveOfMyLife = models.Song{
	ID:     "1",
	Title:  "Love of My Life",
	Author: "Queen",
}

var SongSomebodyToLove = models.Song{
	ID:     "2",
	Title:  "Somebody to Love",
	Author: "Queen",
}

var SongSevenSeasOfRhye = models.Song{
	ID:        "3",
	Title:     "Seven Seas of Rhye",
	Author:    "Queen",
	CreatedAt: "1974-01-01T00:00:00Z",
}

var SongRadioGaGa = models.Song{
	ID:        "4",
	Title:     "Radio Ga Ga",
	Author:    "Queen",
	CreatedAt: "1984-01-01T00:00:00Z",
}

var SongTheShowMustGoOn = models.Song{
	ID:     "5",
	Title:  "The Show Must Go On",
	Author: "Queen",
}

var SongIWantItAll = models.Song{
	ID:     "6",
	Title:  "I Want It All",
	Author: "Queen",
}

var SongBohemianRhapsody = models.Song{
	ID:     "2",
	Title:  "Bohemian Rhapsody",
	Author: "Queen",
}

var SongOneVision = models.Song{
	ID:     "4",
	Title:  "One Vision",
	Author: "Queen",
}

// --- DOCUMENTS ---

var DocSheetMusicGuitar = models.Document{
	ID:              "1",
	SongID:          "s1",
	TitleNormalized: "queen",
	Type:            "sheet_music",
	Instrument:      []string{"Guitar"},
}

var DocSheetMusicPiano = models.Document{
	ID:              "2",
	SongID:          "s2",
	TitleNormalized: "bohemian rhapsody",
	Type:            "sheet_music",
	Instrument:      []string{"Piano"},
}

var DocTablatureGuitar = models.Document{
	ID:              "3",
	SongID:          "s3",
	TitleNormalized: "we will rock you",
	Type:            "tablature",
	Instrument:      []string{"Guitar"},
}

var DocViolinLoveOfMyLife = models.Document{
	ID:              "3",
	SongID:          "s3",
	TitleNormalized: "love of my life",
	Type:            "sheet_music",
	Instrument:      []string{"Violin"},
}

var DocViolinSomebodyToLove = models.Document{
	ID:              "4",
	SongID:          "s4",
	TitleNormalized: "somebody to love",
	Type:            "sheet_music",
	Instrument:      []string{"Violin"},
}

var DocSortedA = models.Document{
	ID:              "7",
	SongID:          "s7",
	TitleNormalized: "a kind of magic",
	Type:            "sheet_music",
}

var DocSortedB = models.Document{
	ID:              "8",
	SongID:          "s8",
	TitleNormalized: "bohemian rhapsody",
	Type:            "sheet_music",
}

var DocUnderPressure = models.Document{
	ID:              "5",
	SongID:          "s5",
	TitleNormalized: "under pressure",
	CreatedAt:       "1985-01-01T00:00:00Z",
}

var DocInnuendo = models.Document{
	ID:              "6",
	SongID:          "s6",
	TitleNormalized: "innuendo",
	CreatedAt:       "1991-01-01T00:00:00Z",
}
