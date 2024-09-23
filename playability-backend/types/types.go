package types

type Game struct {
	ID                    int    `json:"id"`
	Name                  string `json:"name"`
	Cover                 int    `json:"cover"`
	Summary               string `json:"summary"`
	Platforms             []int  `json:"platforms"`
	InvolvedCompanies     []int  `json:"involved_companies"`
	ExternalGames         []int  `json:"external_games"`
	CoverArt              string `json:"cover_art"`
	ClosedCaptions        string `json:"closed_captions"`
	ColorBlind            string `json:"color_blind"`
	FullControllerSupport string `json:"full_controller_support"`
	ControllerRemapping   string `json:"controller_remapping"`
	SteamAvailability     bool   `json:"steam_availability"`
}

// CoverArt represents the structure of cover art data from the IGDB API
type CoverArt struct {
	ID      int    `json:"id"`
	ImageID string `json:"image_id"`
}

type ExternalGame struct {
	ID  int    `json:"id"`
	UID string `json:"uid"`
}

// PCGamingWikiResponse represents the structure of responses from PCGamingWiki API
type PCGamingWikiResponse struct {
	CargoQuery []struct {
		Title struct {
			ClosedCaptions        string `json:"Closed captions,omitempty"`
			ColorBlind            string `json:"Color blind,omitempty"`
			FullControllerSupport string `json:"Full controller support,omitempty"`
			ControllerRemapping   string `json:"Controller remapping,omitempty"`
		} `json:"title"`
	} `json:"cargoquery"`
}
type UserLogin struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type UserRegister struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type UserRow struct {
	ID           int    `json:"id"`
	Username     string `json:"username"`
	Email        string `json:"email"`
	Hash         string `json:"hash"`
	NumOfReports int    `json:"num_of_reports"`
}
type ReportRegister struct {
	GameID                int    `json:"game_id"`
	UserID                int    `json:"user_id"`
	ClosedCaptions        string `json:"closed_captions"`
	ColorBlind            string `json:"color_blind"`
	FullControllerSupport string `json:"full_controller_support"`
	ControllerRemapping   string `json:"controller_remapping"`
	Score                 int    `json:"score"`
	Report                string `json:"report"`
}
type ReportRow struct {
	ID                    int    `json:"id"`
	GameID                int    `json:"game_id"`
	UserID                int    `json:"user_id"`
	ClosedCaptions        string `json:"closed_captions"`
	ColorBlind            string `json:"color_blind"`
	FullControllerSupport string `json:"full_controller_support"`
	ControllerRemapping   string `json:"controller_remapping"`
	Score                 int    `json:"score"`
	Report                string `json:"report"`
}
type ReportCards struct {
	ID     int    `json:"id"`
	GameID int    `json:"game_id"`
	UserID int    `json:"user_id"`
	Score  int    `json:"score"`
	Report string `json:"report"`
}
type FeatureReport struct {
	ID                    int    `json:"id"`
	GameID                int    `json:"game_id"`
	ClosedCaptions        string `json:"closed_captions"`
	ColorBlind            string `json:"color_blind"`
	FullControllerSupport string `json:"full_controller_support"`
	ControllerRemapping   string `json:"controller_remapping"`
}
type FeatureStat struct {
	FeatureName        string  `json:"feature_name"`
	Consensus          string  `json:"consensus"`
	SecondaryConsensus string  `json:"secondary_consensus"`
	TruePercentage     float64 `json:"true_percentage"`
	LimitedPercentage  float64 `json:"limited_percentage"`
	FalsePercentage    float64 `json:"false_percentage"`
}
