package main

import (
	"fmt"
	"image/color"
	"log"
	"math/rand"

	"github.com/golang/freetype/truetype"
	"golang.org/x/image/font"

	"github.com/hajimehoshi/ebiten"
	"github.com/hajimehoshi/ebiten/ebitenutil"
	"github.com/hajimehoshi/ebiten/examples/resources/fonts"
	"github.com/hajimehoshi/ebiten/text"
)

const (
	screenWidth, screenHeight = 1280, 960
)

var (
	err           error
	background    *ebiten.Image
	spaceShip     *ebiten.Image
	aliveChicken  *ebiten.Image
	deadChicken   *ebiten.Image
	bullet        *ebiten.Image
	playerOne     item
	bullets       []item
	chickens      []item
	cooked        []item
	didYouLose    bool
	score         int64
	scoreColor    = color.NRGBA{200, 0, 200, 0x80}
	gameOverColor = color.NRGBA{200, 0, 0, 0x80}
)

const (
	arcadeFontBaseSize = 8
)

var (
	arcadeFonts map[int]font.Face
)

type item struct {
	image      *ebiten.Image
	xPos, yPos float64
	speed      float64
}

func init() {
	background, _, err = ebitenutil.NewImageFromFile("assets/space.png", ebiten.FilterDefault)
	if err != nil {
		log.Fatal(err)
	}

	spaceShip, _, err = ebitenutil.NewImageFromFile("assets/spaceship.png", ebiten.FilterDefault)
	if err != nil {
		log.Fatal(err)
	}

	aliveChicken, _, err = ebitenutil.NewImageFromFile("assets/rooster.png", ebiten.FilterDefault)
	if err != nil {
		log.Fatal(err)
	}

	deadChicken, _, err = ebitenutil.NewImageFromFile("assets/chicken.png", ebiten.FilterDefault)
	if err != nil {
		log.Fatal(err)
	}

	bullet, _, err = ebitenutil.NewImageFromFile("assets/bulletM.png", ebiten.FilterDefault)
	if err != nil {
		log.Fatal(err)
	}

	playerOne = item{spaceShip, screenWidth / 2.0, screenHeight / 2.0, 10}
}

func removeItem(s []item, index int) []item {
	return append(s[:index], s[index+1:]...)
}

func addNewChicken() {
	newChicken := item{aliveChicken, float64(rand.Intn(1200)), 0, 4}
	collisionWithOtherChicken := true
	for collisionWithOtherChicken == true {
		collisionWithOtherChicken = false
		newChicken.xPos = float64(rand.Intn(1200))
		for i := 0; i < len(chickens); i++ {
			if newChicken.xPos >= chickens[i].xPos && newChicken.xPos <= chickens[i].xPos+64 {
				collisionWithOtherChicken = true
			}
			if newChicken.xPos+64 >= chickens[i].xPos && newChicken.xPos+64 <= chickens[i].xPos+64 {
				collisionWithOtherChicken = true
			}
		}
	}
	chickens = append(chickens, newChicken)
}

func didCollisionHappen(x1, y1, siz1 float64, x2, y2, siz2 float64) bool {
	if x1+siz1 < x2 || x2+siz2 < x1 || y1+siz1 < y2 || y2+siz2 < y1 {
		return false
	}
	return true
}

func captureAndUpdate() {
	if ebiten.IsKeyPressed(ebiten.KeyUp) {
		playerOne.yPos -= playerOne.speed
		if playerOne.yPos < 0 {
			playerOne.yPos = 0
		}
	}
	if ebiten.IsKeyPressed(ebiten.KeyDown) {
		playerOne.yPos += playerOne.speed
		if playerOne.yPos > screenHeight-64 {
			playerOne.yPos = screenHeight - 64
		}
	}
	if ebiten.IsKeyPressed(ebiten.KeyLeft) {
		playerOne.xPos -= playerOne.speed
		if playerOne.xPos < 0 {
			playerOne.xPos = 0
		}
	}
	if ebiten.IsKeyPressed(ebiten.KeyRight) {
		playerOne.xPos += playerOne.speed
		if playerOne.xPos > screenWidth-64 {
			playerOne.xPos = screenWidth - 64
		}
	}

	if ebiten.IsKeyPressed(ebiten.KeySpace) {
		if len(bullets) == 0 || bullets[len(bullets)-1].yPos+64 < playerOne.yPos {
			newBullet := item{bullet, playerOne.xPos + 16, playerOne.yPos - 32, 6}
			bullets = append(bullets, newBullet)
		}
	}

	if len(chickens) < 6 {
		addNewChicken()
	}

	for i := 0; i < len(bullets); i++ {
		bullets[i].yPos -= bullets[i].speed
		if bullets[i].yPos < -0 {
			bullets = removeItem(bullets, i)
		}
	}

	for i := 0; i < len(chickens); i++ {
		chickens[i].yPos += chickens[i].speed
		if didCollisionHappen(playerOne.xPos, playerOne.yPos, float64(64),
			chickens[i].xPos, chickens[i].yPos, float64(64)) {
			didYouLose = true
		}
		for j := 0; j < len(bullets); j++ {
			if didCollisionHappen(bullets[j].xPos, bullets[j].yPos, float64(24),
				chickens[i].xPos, chickens[i].yPos, float64(64)) {
				chickens[i].image = deadChicken
				chickens[i].speed = 6
				cooked = append(cooked, chickens[i])
				chickens = removeItem(chickens, i)
				i = -1
				break
			}
		}
		if i >= 0 && i < len(chickens) && chickens[i].yPos+64 >= screenHeight {
			// chickens = removeItem(chickens, i)
			didYouLose = true
		}
	}

	for i := 0; i < len(cooked); i++ {
		cooked[i].yPos += cooked[i].speed
		if didCollisionHappen(playerOne.xPos, playerOne.yPos, float64(64),
			cooked[i].xPos, cooked[i].yPos, 32) {
			score++
			cooked = removeItem(cooked, i)
			i--
		} else if cooked[i].yPos > screenHeight {
			cooked = removeItem(cooked, i)
			i--
		}
	}
}

func captureRestart() {
	if ebiten.IsKeyPressed(ebiten.KeyEnter) {
		bullets = bullets[:0]
		chickens = chickens[:0]
		cooked = cooked[:0]
		playerOne = item{spaceShip, screenWidth / 2.0, screenHeight / 2.0, 10}
		score = 0
		didYouLose = false
	}
}

func getArcadeFonts(scale int) font.Face {
	if arcadeFonts == nil {
		tt, err := truetype.Parse(fonts.ArcadeN_ttf)
		if err != nil {
			log.Fatal(err)
		}

		arcadeFonts = map[int]font.Face{}
		for i := 1; i <= 10; i++ {
			const dpi = 72
			arcadeFonts[i] = truetype.NewFace(tt, &truetype.Options{
				Size:    float64(arcadeFontBaseSize * i),
				DPI:     dpi,
				Hinting: font.HintingFull,
			})
		}
	}
	return arcadeFonts[scale]
}

func update(screen *ebiten.Image) error {
	if didYouLose == false {
		captureAndUpdate()
	} else {
		captureRestart()
	}

	if ebiten.IsDrawingSkipped() {
		return nil
	}
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(0, 0)
	screen.DrawImage(background, op)

	playerOp := &ebiten.DrawImageOptions{}
	playerOp.GeoM.Translate(playerOne.xPos, playerOne.yPos)
	screen.DrawImage(playerOne.image, playerOp)

	for i := 0; i < len(bullets); i++ {
		blletOp := &ebiten.DrawImageOptions{}
		blletOp.GeoM.Translate(bullets[i].xPos, bullets[i].yPos)
		screen.DrawImage(bullets[i].image, blletOp)
	}

	for i := 0; i < len(chickens); i++ {
		chickensOp := &ebiten.DrawImageOptions{}
		chickensOp.GeoM.Translate(chickens[i].xPos, chickens[i].yPos)
		screen.DrawImage(chickens[i].image, chickensOp)
	}

	for i := 0; i < len(cooked); i++ {
		cookedOp := &ebiten.DrawImageOptions{}
		cookedOp.GeoM.Translate(cooked[i].xPos, cooked[i].yPos)
		screen.DrawImage(cooked[i].image, cookedOp)
	}

	text.Draw(screen, fmt.Sprintf("Your Score: %d", score), getArcadeFonts(2), 20, 40, scoreColor)

	if didYouLose == true {
		text.Draw(screen, fmt.Sprintf("Game Over"), getArcadeFonts(8), 360, 500, gameOverColor)
		text.Draw(screen, fmt.Sprintf("Your Score: %d", score), getArcadeFonts(3), 500, 550, scoreColor)
		text.Draw(screen, fmt.Sprintf("Press Enter to Restart"), getArcadeFonts(2), 480, 600, scoreColor)
	}

	return nil
}

func main() {
	if err := ebiten.Run(update, screenWidth, screenHeight, 1, "Space Chickens"); err != nil {
		log.Fatal(err)
	}
}
