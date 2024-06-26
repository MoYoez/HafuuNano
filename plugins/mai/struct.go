package mai

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/FloatTech/gg"
	"github.com/FloatTech/imgfactory"
	nano "github.com/fumiama/NanoBot"
	"github.com/moyoez/HafuuNano/utils"
	"golang.org/x/image/font"
	"golang.org/x/text/width"
	"image"
	"image/color"
	"io"
	"net/http"
	"os"
	"strconv"
	"strings"
	"sync"
	"unicode/utf8"
)

type DivingFishB50UserName struct {
	Username string `json:"username"`
	B50      bool   `json:"b50"`
}

type DivingFishDevFullDataRecords struct {
	AdditionalRating int    `json:"additional_rating"`
	Nickname         string `json:"nickname"`
	Plate            string `json:"plate"`
	Rating           int    `json:"rating"`
	Records          []struct {
		Achievements float64 `json:"achievements"`
		Ds           float64 `json:"ds"`
		DxScore      int     `json:"dxScore"`
		Fc           string  `json:"fc"`
		Fs           string  `json:"fs"`
		Level        string  `json:"level"`
		LevelIndex   int     `json:"level_index"`
		LevelLabel   string  `json:"level_label"`
		Ra           int     `json:"ra"`
		Rate         string  `json:"rate"`
		SongId       int     `json:"song_id"`
		Title        string  `json:"title"`
		Type         string  `json:"type"`
	} `json:"records"`
	Username string `json:"username"`
}

type WebPingStauts struct {
	Details struct {
		MaimaiDXCN struct {
			Uptime float64 `json:"uptime"`
		} `json:"maimai DX CN"`
		MaimaiDXCNDXNet struct {
			Uptime float64 `json:"uptime"`
		} `json:"maimai DX CN DXNet"`
		MaimaiDXCNMain struct {
			Uptime float64 `json:"uptime"`
		} `json:"maimai DX CN Main"`
		MaimaiDXCNNetLogin struct {
			Uptime float64 `json:"uptime"`
		} `json:"maimai DX CN NetLogin"`
		MaimaiDXCNTitle struct {
			Uptime float64 `json:"uptime"`
		} `json:"maimai DX CN Title"`
		MaimaiDXCNUpdate struct {
			Uptime float64 `json:"uptime"`
		} `json:"maimai DX CN Update"`
	} `json:"details"`
	Status bool `json:"status"`
}

type RealConvertPlay struct {
	ReturnValue []struct {
		SkippedCount  int `json:"skippedCount"`
		RetriedCount  int `json:"retriedCount"`
		RetryCountSum int `json:"retryCountSum"`
		TotalCount    int `json:"totalCount"`
		FailedCount   int `json:"failedCount"`
	} `json:"returnValue"`
}

type ZlibErrorStatus struct {
	Full struct {
		Field1 int `json:"10"`
		Field2 int `json:"30"`
		Field3 int `json:"60"`
	} `json:"full"`
	FullError struct {
		Field1 int `json:"10"`
		Field2 int `json:"30"`
		Field3 int `json:"60"`
	} `json:"full_Error"`
	ZlibError struct {
		Field1 int `json:"10"`
		Field2 int `json:"30"`
		Field3 int `json:"60"`
	} `json:"zlib_Error"`
}

type player struct {
	AdditionalRating int `json:"additional_rating"`
	Charts           struct {
		Dx []struct {
			Achievements float64 `json:"achievements"`
			Ds           float64 `json:"ds"`
			DxScore      int     `json:"dxScore"`
			Fc           string  `json:"fc"`
			Fs           string  `json:"fs"`
			Level        string  `json:"level"`
			LevelIndex   int     `json:"level_index"`
			LevelLabel   string  `json:"level_label"`
			Ra           int     `json:"ra"`
			Rate         string  `json:"rate"`
			SongId       int     `json:"song_id"`
			Title        string  `json:"title"`
			Type         string  `json:"type"`
		} `json:"dx"`
		Sd []struct {
			Achievements float64 `json:"achievements"`
			Ds           float64 `json:"ds"`
			DxScore      int     `json:"dxScore"`
			Fc           string  `json:"fc"`
			Fs           string  `json:"fs"`
			Level        string  `json:"level"`
			LevelIndex   int     `json:"level_index"`
			LevelLabel   string  `json:"level_label"`
			Ra           int     `json:"ra"`
			Rate         string  `json:"rate"`
			SongId       int     `json:"song_id"`
			Title        string  `json:"title"`
			Type         string  `json:"type"`
		} `json:"sd"`
	} `json:"charts"`
	Nickname string      `json:"nickname"`
	Plate    string      `json:"plate"`
	Rating   int         `json:"rating"`
	UserData interface{} `json:"user_data"`
	UserId   interface{} `json:"user_id"`
	Username string      `json:"username"`
}

type playerData struct {
	Achievements float64 `json:"achievements"`
	Ds           float64 `json:"ds"`
	DxScore      int     `json:"dxScore"`
	Fc           string  `json:"fc"`
	Fs           string  `json:"fs"`
	Level        string  `json:"level"`
	LevelIndex   int     `json:"level_index"`
	LevelLabel   string  `json:"level_label"`
	Ra           int     `json:"ra"`
	Rate         string  `json:"rate"`
	SongId       int     `json:"song_id"`
	Title        string  `json:"title"`
	Type         string  `json:"type"`
}

var (
	loadMaiPic        = Root + "pic/"
	defaultCoverLink  = Root + "default_cover.png"
	typeImageDX       = loadMaiPic + "chart_type_dx.png"
	typeImageSD       = loadMaiPic + "chart_type_sd.png"
	titleFontPath     = maifont + "NotoSansSC-Bold.otf"
	UniFontPath       = maifont + "Montserrat-Bold.ttf"
	nameFont          = maifont + "NotoSansSC-Regular.otf"
	maifont           = Root + "font/"
	b50bgOriginal     = loadMaiPic + "b50_bg.png"
	b50bg             = loadMaiPic + "b50_bg.png"
	b50Custom         = loadMaiPic + "b50_bg_custom.png"
	Root              = utils.ReturnLucyMainDataIndex("maidx") + "resources/maimai/"
	userPlate         = utils.ResLoader("mai") + "user/"
	titleFont         font.Face
	scoreFont         font.Face
	rankFont          font.Face
	levelFont         font.Face
	ratingFont        font.Face
	nameTypeFont      font.Face
	diffColor         []color.RGBA
	ratingBgFilenames = []string{
		"rating_white.png",
		"rating_blue.png",
		"rating_green.png",
		"rating_yellow.png",
		"rating_red.png",
		"rating_purple.png",
		"rating_copper.png",
		"rating_silver.png",
		"rating_gold.png",
		"rating_rainbow.png",
	}
)

func init() {
	if _, err := os.Stat(userPlate); os.IsNotExist(err) {
		err := os.MkdirAll(userPlate, 0777)
		if err != nil {
			return
		}
	}
	nameTypeFont = utils.LoadFontFace(nameFont, 36)
	titleFont = utils.LoadFontFace(titleFontPath, 20)
	scoreFont = utils.LoadFontFace(UniFontPath, 32)
	rankFont = utils.LoadFontFace(UniFontPath, 24)
	levelFont = utils.LoadFontFace(UniFontPath, 20)
	ratingFont = utils.LoadFontFace(UniFontPath, 24)
	diffColor = []color.RGBA{
		{69, 193, 36, 255},
		{255, 186, 1, 255},
		{255, 90, 102, 255},
		{134, 49, 200, 255},
		{207, 144, 240, 255},
	}
}

func QueryMaiBotDataFromUserName(username string) (playerdata []byte, err error) {
	// packed json and sent.
	jsonStruct := DivingFishB50UserName{Username: username, B50: true}
	jsonStructData, err := json.Marshal(jsonStruct)
	if err != nil {
		return nil, err
	}
	req, err := http.NewRequest("POST", "https://www.diving-fish.com/api/maimaidxprober/query/player", bytes.NewBuffer(jsonStructData))
	req.Header.Set("Content-Type", "application/json")
	if err != nil {
		return
	}
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return
	}
	defer resp.Body.Close()
	if resp.StatusCode == 400 {
		return nil, errors.New("- 未找到用户或者用户数据丢失\n\n - 请检查您是否在 水鱼查分器 上 上传过成绩")
	}
	if resp.StatusCode == 403 {
		return nil, errors.New("- 该用户设置禁止查分\n\n - 请检查您是否在 水鱼查分器/ 上 是否关闭了允许他人查分功能")
	}
	playerDataByte, err := io.ReadAll(resp.Body)
	return playerDataByte, err
}

// FullPageRender  Render Full Page
func FullPageRender(data player, ctx *nano.Ctx) (raw image.Image) {
	// muilt-threading.
	getUserID := ctx.UserID()
	var avatarHandler sync.WaitGroup
	avatarHandler.Add(1)
	var getAvatarFormat *gg.Context
	// avatar handler.
	go func() {
		// avatar Round Style
		defer avatarHandler.Done()
		getAvatar := utils.GetUserAvatar(ctx)
		if getAvatar != nil {
			avatarFormat := imgfactory.Size(getAvatar, 180, 180)
			getAvatarFormat = gg.NewContext(180, 180)
			getAvatarFormat.DrawRoundedRectangle(0, 0, 178, 178, 20)
			getAvatarFormat.Clip()
			getAvatarFormat.DrawImage(avatarFormat.Image(), 0, 0)
			getAvatarFormat.Fill()
		}
	}()
	userPlatedCustom := utils.QueryUserbaseMaiData(ctx.UserID()).Background
	// render Header.
	b50Render := gg.NewContext(2090, 1660)
	rawPlateData, errs := gg.LoadImage(userPlate + strconv.Itoa(int(getUserID)) + ".png")
	if errs == nil {
		b50bg = b50Custom
		b50Render.DrawImage(rawPlateData, 595, 30)
		b50Render.Fill()
	} else {
		if userPlatedCustom != "" {
			b50bg = b50Custom
			images, _ := GetDefaultPlate(userPlatedCustom)
			b50Render.DrawImage(images, 595, 30)
			b50Render.Fill()
		} else {
			// show nil
			b50bg = b50bgOriginal
		}
	}
	getContent, _ := gg.LoadImage(b50bg)
	b50Render.DrawImage(getContent, 0, 0)
	b50Render.Fill()
	// render user info
	avatarHandler.Wait()
	if getAvatarFormat != nil {
		b50Render.DrawImage(getAvatarFormat.Image(), 610, 50)
		b50Render.Fill()
	}
	// render Userinfo
	b50Render.SetFontFace(nameTypeFont)
	b50Render.SetColor(color.Black)
	b50Render.DrawStringAnchored(width.Widen.String(data.Nickname), 825, 160, 0, 0)
	b50Render.Fill()
	b50Render.SetFontFace(titleFont)
	setPlateLocalStatus := utils.QueryUserbaseMaiData(ctx.UserID()).Plate
	if setPlateLocalStatus != "" {
		data.Plate = setPlateLocalStatus
	}
	b50Render.DrawStringAnchored(strings.Join(strings.Split(data.Plate, ""), " "), 1050, 207, 0.5, 0.5)
	b50Render.Fill()
	getRating := getRatingBg(data.Rating)
	getRatingBG, err := gg.LoadImage(loadMaiPic + getRating)
	if err != nil {
		return
	}
	b50Render.DrawImage(getRatingBG, 800, 40)
	b50Render.Fill()
	// render Rank
	imgs, err := GetRankPicRaw(data.AdditionalRating)
	if err != nil {
		return
	}
	b50Render.DrawImage(imgs, 1080, 50)
	b50Render.Fill()
	// draw number
	b50Render.SetFontFace(scoreFont)
	b50Render.SetRGBA255(236, 219, 113, 255)
	b50Render.DrawStringAnchored(strconv.Itoa(data.Rating), 1056, 60, 1, 1)
	b50Render.Fill()
	// Render Card Type
	getSDLength := len(data.Charts.Sd)
	getDXLength := len(data.Charts.Dx)
	getDXinitX := 45
	getDXinitY := 1225
	getInitX := 45
	getInitY := 285
	var i int
	for i = 0; i < getSDLength; i++ {
		b50Render.DrawImage(RenderCard(data.Charts.Sd[i], i+1, false), getInitX, getInitY)
		getInitX += 400
		if getInitX == 2045 {
			getInitX = 45
			getInitY += 125
		}
	}

	for dx := 0; dx < getDXLength; dx++ {
		b50Render.DrawImage(RenderCard(data.Charts.Dx[dx], dx+1, false), getDXinitX, getDXinitY)
		getDXinitX += 400
		if getDXinitX == 2045 {
			getDXinitX = 45
			getDXinitY += 125
		}
	}
	return b50Render.Image()
}

// RenderCard Main Lucy Render Page , if isSimpleRender == true, then render count will not show here.
func RenderCard(data playerData, num int, isSimpleRender bool) image.Image {
	getType := data.Type
	var CardBackGround string
	var multiTypeRender sync.WaitGroup
	var CoverDownloader sync.WaitGroup
	CoverDownloader.Add(1)
	multiTypeRender.Add(1)
	// choose Type.
	if getType == "SD" {
		CardBackGround = typeImageSD
	} else {
		CardBackGround = typeImageDX
	}
	charCount := 0.0
	setBreaker := false
	var truncated string
	var charFloatNum float64
	getSongName := data.Title
	var getSongId string
	switch {
	case data.SongId < 1000:
		getSongId = fmt.Sprintf("%05d", data.SongId)
	case data.SongId < 10000:
		getSongId = fmt.Sprintf("1%d", data.SongId)
	default:
		getSongId = strconv.Itoa(data.SongId)
	}
	var Image image.Image
	go func() {
		defer CoverDownloader.Done()
		Image, _ = GetCover(getSongId)
	}()
	// set rune count
	go func() {
		defer multiTypeRender.Done()
		for _, runeValue := range getSongName {
			charWidth := utf8.RuneLen(runeValue)
			if charWidth == 3 {
				charFloatNum = 1.5
			} else {
				charFloatNum = float64(charWidth)
			}
			if charCount+charFloatNum > 19 {
				setBreaker = true
				break
			}
			truncated += string(runeValue)
			charCount += charFloatNum
		}
		if setBreaker {
			getSongName = truncated + ".."
		} else {
			getSongName = truncated
		}
	}()
	loadSongType, _ := gg.LoadImage(CardBackGround)
	// draw pic
	drawBackGround := gg.NewContextForImage(GetChartType(data.LevelLabel))
	// draw song pic
	CoverDownloader.Wait()
	drawBackGround.DrawImage(Image, 25, 25)
	// draw name
	drawBackGround.SetColor(color.White)
	drawBackGround.SetFontFace(titleFont)
	multiTypeRender.Wait()
	drawBackGround.DrawStringAnchored(getSongName, 130, 32.5, 0, 0.5)
	drawBackGround.Fill()
	// draw acc
	drawBackGround.SetFontFace(scoreFont)
	drawBackGround.DrawStringAnchored(strconv.FormatFloat(data.Achievements, 'f', 4, 64)+"%", 129, 62.5, 0, 0.5)
	// draw rate
	drawBackGround.DrawImage(GetRateStatusAndRenderToImage(data.Rate), 305, 45)
	drawBackGround.Fill()
	drawBackGround.SetFontFace(rankFont)
	drawBackGround.SetColor(diffColor[data.LevelIndex])
	if !isSimpleRender {
		drawBackGround.DrawString("#"+strconv.Itoa(num), 130, 111)
	}
	drawBackGround.FillPreserve()
	// draw rest of card.
	drawBackGround.SetFontFace(levelFont)
	drawBackGround.DrawString(strconv.FormatFloat(data.Ds, 'f', 1, 64), 195, 111)
	drawBackGround.FillPreserve()
	drawBackGround.SetFontFace(ratingFont)
	drawBackGround.DrawString("▶", 235, 111)
	drawBackGround.FillPreserve()
	drawBackGround.SetFontFace(ratingFont)
	drawBackGround.DrawString(strconv.Itoa(data.Ra), 250, 111)
	drawBackGround.FillPreserve()
	if data.Fc != "" {
		drawBackGround.DrawImage(LoadComboImage(data.Fc), 290, 84)
	}
	if data.Fs != "" {
		drawBackGround.DrawImage(LoadSyncImage(data.Fs), 325, 84)
	}
	drawBackGround.DrawImage(loadSongType, 68, 88)
	return drawBackGround.Image()
}

func GetRankPicRaw(id int) (image.Image, error) {
	var idStr string
	if id < 10 {
		idStr = "0" + strconv.FormatInt(int64(id), 10)
	} else {
		idStr = strconv.FormatInt(int64(id), 10)
	}
	if id == 22 {
		idStr = "21"
	}
	data := Root + "rank/UI_CMN_DaniPlate_" + idStr + ".png"
	imgRaw, err := gg.LoadImage(data)
	if err != nil {
		return nil, err
	}
	return imgRaw, nil
}

func GetDefaultPlate(id string) (image.Image, error) {
	data := Root + "plate/plate_" + id + ".png"
	imgRaw, err := gg.LoadImage(data)
	if err != nil {
		return nil, err
	}
	return imgRaw, nil
}

// GetCover Careful The nil data
func GetCover(id string) (image.Image, error) {
	fileName := id + ".png"
	filePath := Root + "cover/" + fileName
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		// Auto download cover from diving fish's site
		downloadURL := "https://www.diving-fish.com/covers/" + fileName
		cover, err := utils.DownloadImage(downloadURL)
		if err != nil {
			return utils.LoadPictureWithResize(defaultCoverLink, 90, 90), nil
		}
		utils.SaveImage(cover, filePath)
	}
	imageFile, err := os.Open(filePath)
	if err != nil {
		return utils.LoadPictureWithResize(defaultCoverLink, 90, 90), nil
	}
	defer func(imageFile *os.File) {
		err := imageFile.Close()
		if err != nil {
			return
		}
	}(imageFile)
	img, _, err := image.Decode(imageFile)
	if err != nil {
		return utils.LoadPictureWithResize(defaultCoverLink, 90, 90), nil
	}
	return utils.Resize(img, 90, 90), nil
}

// GetRateStatusAndRenderToImage Get Rate
func GetRateStatusAndRenderToImage(rank string) image.Image {
	// Load rank images
	return utils.LoadPictureWithResize(loadMaiPic+"rate_"+rank+".png", 80, 40)
}

// GetChartType Get Chart Type
func GetChartType(chart string) image.Image {
	data, _ := gg.LoadImage(loadMaiPic + "chart_" + NoHeadLineCase(chart) + ".png")
	return data
}

// LoadComboImage Load combo images
func LoadComboImage(imageName string) image.Image {
	link := loadMaiPic + "combo_" + imageName + ".png"
	return utils.LoadPictureWithResize(link, 60, 40)
}

// LoadSyncImage Load sync images
func LoadSyncImage(imageName string) image.Image {
	link := loadMaiPic + "sync_" + imageName + ".png"
	return utils.LoadPictureWithResize(link, 60, 40)
}

// NoHeadLineCase No HeadLine.
func NoHeadLineCase(word string) string {
	text := strings.ToLower(word)
	textNewer := strings.ReplaceAll(text, ":", "")
	return textNewer
}

func getRatingBg(rating int) string {
	index := 0
	switch {
	case rating >= 15000:
		index++
		fallthrough
	case rating >= 14000:
		index++
		fallthrough
	case rating >= 13000:
		index++
		fallthrough
	case rating >= 12000:
		index++
		fallthrough
	case rating >= 10000:
		index++
		fallthrough
	case rating >= 7000:
		index++
		fallthrough
	case rating >= 4000:
		index++
		fallthrough
	case rating >= 2000:
		index++
		fallthrough
	case rating >= 1000:
		index++
	}
	return ratingBgFilenames[index]
}

/*
// MixedRegionWriter Some Mixed Magic, looking for your region information.
func MixedRegionWriter(regionID int, playCount int, createdDate string) string {
	getCountryID := returnCountryID(regionID)
	return fmt.Sprintf(" - 在 regionID 为 %d (%s) 的省/直辖市 游玩过 %d 次, 第一次游玩时间于 %s", regionID+1, getCountryID, playCount, createdDate)
}

// ReportToEndPoint Report Some Error To Wahlap Server.
func ReportToEndPoint(getReport int, getReportType string) string {
	url := "https://maihook.lemonkoi.one/api/zlib?report=" + strconv.Itoa(getReport) + "&reportType=" + getReportType
	method := "GET"
	client := &http.Client{}
	req, err := http.NewRequest(method, url, nil)
	if err != nil {
		fmt.Println(err)
		return ""
	}
	req.Header.Add("authkey", authKey)
	req.Header.Add("Host", "maihook.lemonkoi.one")
	req.Header.Add("Connection", "keep-alive")
	res, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		return ""
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		fmt.Println(err)
		return ""
	}
	return string(body)
}

// ReturnZlibError Return Zlib ERROR
func ReturnZlibError() ZlibErrorStatus {
	url := "https://maihook.lemonkoi.one/api/zlib"
	method := "GET"
	client := &http.Client{}
	req, err := http.NewRequest(method, url, nil)

	if err != nil {
		fmt.Println(err)
		return ZlibErrorStatus{}
	}
	req.Header.Add("Host", "maihook.lemonkoi.one")
	req.Header.Add("Connection", "keep-alive")
	res, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		return ZlibErrorStatus{}
	}
	defer res.Body.Close()
	body, err := io.ReadAll(res.Body)
	if err != nil {
		fmt.Println(err)
		return ZlibErrorStatus{}
	}
	var returnData ZlibErrorStatus
	json.Unmarshal(body, &returnData)
	return returnData
}

func ConvertZlib(value, total int) string {
	if total == 0 {
		return "0.000%"
	}
	percentage := float64(value) / float64(total) * 100
	return fmt.Sprintf("%.3f%%", percentage)
}

func ConvertRealPlayWords(retry RealConvertPlay) string {
	var pickedWords string
	var count = 0
	var header = " - 错误率数据收集自机台的真实网络通信，可以反映舞萌 DX 的网络状况。\n"

	for _, word := range retry.ReturnValue {
		var timeCount int
		var UserReturnLogs string
		switch {
		case count == 0:
			timeCount = 10
		case count == 1:
			timeCount = 30
		case count == 2:
			timeCount = 60
		}

		if word.TotalCount < 20 {
			UserReturnLogs = "没有收集到足够的数据进行分析~"
		} else {
			totalSuccess := word.TotalCount - word.FailedCount
			skippedRate := float64(word.SkippedCount) / float64(totalSuccess) * 100
			otherErrorRate := float64(word.RetryCountSum) / float64(totalSuccess+word.RetryCountSum) * 100
			overallErrorRate := (float64(word.SkippedCount+word.RetryCountSum) / float64(totalSuccess+word.RetryCountSum)) * 100
			skippedRate = math.Round(skippedRate*100) / 100
			otherErrorRate = math.Round(otherErrorRate*100) / 100
			overallErrorRate = math.Round(overallErrorRate*100) / 100
			UserReturnLogs = fmt.Sprintf("共 %d 个成功的请求中，有 %d 次未压缩（%.2f%%），有 %d 个请求共 %d 次其他错误（%.2f%%），整体错误率为 %.2f%%。", totalSuccess, word.SkippedCount, skippedRate, word.RetriedCount, word.RetryCountSum, otherErrorRate, overallErrorRate)
		}
		pickedWords = pickedWords + fmt.Sprintf("\n - 在 %d 分钟内%s", timeCount, UserReturnLogs)
		count = count + 1

	}
	return header + pickedWords + "\n"
}

*/
