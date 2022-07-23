package email

import "testing"
import testify "github.com/stretchr/testify/assert"

func TestValidateData(t *testing.T) {
	assert := testify.New(t)
	d := Data{SenderEmail: ""}
	r := validateData(d)
	assert.Contains(r, "address")

	d = Data{SenderEmail: "e"}
	r = validateData(d)
	assert.Contains(r, "address")

	d = Data{SenderEmail: "lorem@l.net", SenderName: " "}
	r = validateData(d)
	assert.NotContains(r, "address")
	assert.Contains(r, "name")

	d = Data{SenderEmail: "lorem@l.t", SenderName: " Lara Croft", Subject: ""}
	r = validateData(d)
	assert.NotContains(r, "address")
	assert.NotContains(r, "name")
	assert.Contains(r, "subject")

	d = Data{SenderEmail: "antioch@k.com", SenderName: "A", Subject: "随机的中文单词"}
	r = validateData(d)
	assert.NotContains(r, "address")
	assert.NotContains(r, "name")
	assert.NotContains(r, "subject")
	assert.Contains(r, "message")

	d = Data{SenderEmail: "k@r", SenderName: "أدامو أحمد", Subject: "كلمات عربية عشوائية", Message: " "}
	r = validateData(d)
	assert.NotContains(r, "address")
	assert.NotContains(r, "name")
	assert.NotContains(r, "subject")
	assert.Contains(r, "message")

	d = Data{SenderEmail: "ä@l.net", SenderName: "Lara Croft  ", Subject: " K 词 ,", Message: " Hallo, das sind zufällige deutsche Wörter", CaptchaResponse: " "}
	r = validateData(d)
	assert.NotContains(r, "address")
	assert.NotContains(r, "name")
	assert.NotContains(r, "subject")
	assert.NotContains(r, "message")
	assert.Contains(r, "humanness")

	d = Data{SenderEmail: "ä@l.net", SenderName: "Lara Croft  ", Subject: " K 词 ,", Message: " Hallo, das sind zufällige deutsche Wörter", CaptchaResponse: "A"}
	r = validateData(d)
	assert.True(len(r) == 0)
}
