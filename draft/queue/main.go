package main

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
	"sync"

	"github.com/gin-gonic/gin"
)

// type Note struct {
// 	Fields []struct{}
// 	Front
// }

// type Card struct {
// 	NoteID int64
// 	Factor int64
// }

type CardQueue struct{}
type Deck struct {
	NewCards     CardQueue
	LearnCards   CardQueue
	ReviewCards  CardQueue
	RelearnCards CardQueue
}

const DefaultStartingEase = 250
const DefaultIntervalModifier = 100

func GetNewInterval(currentInterval int64, ease int64, intervalModifier int64) int64 {
	return currentInterval * ease * intervalModifier / 100 / 100
}

func main() {
	if err := run(); err != nil {
		fmt.Println(err)
	}
}

func run() (err error) {
	r := gin.New()
	r.Use(gin.Logger())
	r.Use(gin.Recovery())

	r.Static("/static", "./static")
	r.LoadHTMLGlob("./templates/*")
	Router(r)

	log.Println("server started")
	r.Run()
	return
}

func Router(r *gin.Engine) {
	r.GET("/", index)
	r.GET("/review/back/:id", showBack)
	r.GET("/review/again/:id", answerAgain)
	r.GET("/review/pass/:id", answerPass)
}

func index(c *gin.Context) {
	queuedCards := cardManager.GetQueuedCards(1)
	if len(queuedCards) == 0 {
		congrats(c)
		return
	}

	c.HTML(http.StatusOK, "index.html", gin.H{
		"Card": queuedCards[0],
		"Back": false,
	})
}

func GetCardID(c *gin.Context) (id int64, err error) {
	idParam := c.Param("id")

	id, err = strconv.ParseInt(idParam, 10, 64)
	if err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	return
}

func showBack(c *gin.Context) {
	id, err := GetCardID(c)
	if err != nil {
		return
	}

	_ = id

	card, err := cardManager.GetCard(id)
	if err != nil {
		index(c)
		return
	}

	c.HTML(http.StatusOK, "index.html", gin.H{
		"Card": card,
		"Back": true,
		// "CardID":   card.ID,
		// "Back":     true,
		// "Question": card.Question,
		// "Answer":   card.Answer,
	})
}

func congrats(c *gin.Context) {
	c.HTML(
		http.StatusOK,
		"congrats.html",
		gin.H{},
	)
}

func answerAgain(c *gin.Context) {
	id, err := GetCardID(c)
	if err != nil {
		return
	}

	_ = id

	err = cardManager.Answer(id, 0)
	if err != nil {
		fmt.Println("cardManager.Answer:", err)
	}

	index(c)
}

func answerPass(c *gin.Context) {
}

var cardManager *CardManager = &CardManager{
	cards: map[int64]Card{
		1: Card{
			ID:       1,
			State:    CardNew,
			Question: "How big is sun?",
			Answer:   "Pretty big",
			Answered: false,
		},
		2: Card{
			ID:       2,
			State:    CardNew,
			Question: "How big is earth?",
			Answer:   "Smaller than sun",
			Answered: false,
		},
		3: Card{
			ID:       3,
			State:    CardNew,
			Question: "How big is moon?",
			Answer:   "Smaller than earth",
			Answered: false,
		},
	},
}

type CardState int64

const (
	CardNew CardState = iota
	CardLearning
	CardReview
	CardRelearn
)

type Card struct {
	ID       int64
	State    CardState
	Question string
	Answer   string
	Answered bool
}

type CardManager struct {
	cards map[int64]Card
	mtx   sync.RWMutex
}

func (c *CardManager) GetCard(id int64) (card Card, err error) {
	c.mtx.RLock()
	card, ok := c.cards[id]
	c.mtx.RUnlock()
	if !ok {
		err = fmt.Errorf("card(%d) not found", id)
		return
	}

	return
}

func (c *CardManager) GetQueuedCards(limit int64) (result []Card) {
	for _, c := range c.cards {
		if !c.Answered {
			result = append(result, c)
		}
	}
	fmt.Println("queue size:", len(result))
	return
}

func (c *CardManager) Answer(id int64, answer int64) (err error) {
	c.mtx.Lock()
	defer c.mtx.Unlock()
	card, ok := c.cards[id]
	if !ok {
		return fmt.Errorf("card(%d) not found", id)
	}

	if c.cards[id].Answered {
		return fmt.Errorf("card(%d) already answered", id)
	}

	card.Answered = true
	c.cards[id] = card
	return
}
