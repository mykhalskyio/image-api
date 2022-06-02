package service

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"sync"

	"github.com/mykhalskyio/image-api/internal/config"
	"github.com/mykhalskyio/image-api/internal/entity"
	"github.com/mykhalskyio/image-api/pkg/resize"
	"github.com/rabbitmq/amqp091-go"
)

// image storage interface
type ImageStorage interface {
	Insert(img *entity.Image) error
	Get(id int) (*entity.Image, error)
	Delete(id int) error
}

// image struct
type ImageService struct {
	storage ImageStorage
	cfg     *config.Config
}

// new image struct
func NewImageService(storage ImageStorage, cfg *config.Config) *ImageService {
	return &ImageService{storage: storage, cfg: config.GetConfig()}
}

var once sync.Once

// upload image to RabbitMQ
func (img *ImageService) Upload(imageOriginal *entity.ImageUpload) error {
	conn, err := amqp091.Dial(fmt.Sprintf("amqp://%s:%s@%s:%s/",
		img.cfg.RabbitMQ.User,
		img.cfg.RabbitMQ.Pass,
		img.cfg.RabbitMQ.Host,
		img.cfg.RabbitMQ.Port))
	if err != nil {
		return err
	}
	defer conn.Close()

	ch, err := conn.Channel()
	if err != nil {
		return err
	}
	defer ch.Close()

	_, err = ch.QueueDeclare("image", true, false, false, false, nil)
	if err != nil {
		return err
	}

	imageOriginalJson, err := json.Marshal(imageOriginal)
	if err != nil {
		return err
	}

	err = ch.Publish("", "image", false, false, amqp091.Publishing{
		DeliveryMode: amqp091.Persistent,
		ContentType:  "application/json",
		Body:         imageOriginalJson,
	})
	if err != nil {
		return err
	}

	// once run goroutine
	var ok bool
	once.Do(func() {
		ok = true
	})
	if ok {
		go img.UploadToDB()
	}

	return nil
}

// save image from RabbitMQ and upload to DB
func (img *ImageService) UploadToDB() {
	conn, err := amqp091.Dial(fmt.Sprintf("amqp://%s:%s@%s:%s/",
		img.cfg.RabbitMQ.User,
		img.cfg.RabbitMQ.Pass,
		img.cfg.RabbitMQ.Host,
		img.cfg.RabbitMQ.Port))
	if err != nil {
		return
	}
	defer conn.Close()

	ch, err := conn.Channel()
	if err != nil {
		return
	}
	defer ch.Close()

	msgs, err := ch.Consume("image", "", true, false, false, false, nil)
	if err != nil {
		return
	}

	var (
		path          string
		path75        string
		path50        string
		path25        string
		image75       []byte
		image50       []byte
		image25       []byte
		imageStruct   entity.Image
		imageOriginal entity.ImageUpload
	)
	for msg := range msgs {
		json.Unmarshal(msg.Body, &imageOriginal)

		path = fmt.Sprintf("100-%s", imageOriginal.ImageName)
		ioutil.WriteFile(path, imageOriginal.Image, 0644)

		image75, _ = resize.ImgResize(path, 75)
		path75 = fmt.Sprintf("75-%s", imageOriginal.ImageName)
		ioutil.WriteFile(path75, image75, 0644)

		image50, _ = resize.ImgResize(path, 50)
		path50 = fmt.Sprintf("50-%s", imageOriginal.ImageName)
		ioutil.WriteFile(path50, image50, 0644)

		image25, _ = resize.ImgResize(path, 25)
		path25 = fmt.Sprintf("25-%s", imageOriginal.ImageName)
		ioutil.WriteFile(path25, image25, 0644)

		imageStruct = entity.Image{
			ImagePathQualityOriginal: path,
			ImagePathQuality75:       path75,
			ImagePathQuality50:       path50,
			ImagePathQuality25:       path25}
		img.storage.Insert(&imageStruct)
	}
}

// download image
func (img *ImageService) Download(id int, quality int) (string, error) {
	image, err := img.storage.Get(id)
	if err != nil {
		return "", err
	}
	var imagePath string
	switch quality {
	case 100:
		imagePath = image.ImagePathQualityOriginal
	case 75:
		imagePath = image.ImagePathQuality75
	case 50:
		imagePath = image.ImagePathQuality50
	case 25:
		imagePath = image.ImagePathQuality25
	}

	return imagePath, nil
}

// delete image
func (img *ImageService) Delete(id int) error {
	image, err := img.storage.Get(id)
	if err != nil {
		return err
	}

	err = img.storage.Delete(id)
	if err != nil {
		return err
	}

	err = os.Remove(image.ImagePathQualityOriginal)
	if err != nil {
		return err
	}
	err = os.Remove(image.ImagePathQuality75)
	if err != nil {
		return err
	}
	err = os.Remove(image.ImagePathQuality50)
	if err != nil {
		return err
	}
	err = os.Remove(image.ImagePathQuality25)
	if err != nil {
		return err
	}

	return nil
}
