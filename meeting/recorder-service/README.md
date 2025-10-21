flow


1. request o doloczenie do rozmowy jako nagrywarka
```go

type RecordMeetingRequest struct {
	RoomID string `json:"roomId"`
	FinishTime time.Time `json:"finishTime"`
}

```

2. zapis do bazy ze bot o jakims id jest w pokoju roomid -> redis ttl to jest dlugosc rozmowy
ttl jest liczony na podstwie now i finish time + d - opoznienie

- zapis w bazie pozwoli sledzic czy mozna dodac ponownie bota do rozmo

3. tym samym czasie jest generowane zdarzenie to do gateway'a z requestem dodania bota

```go
type RecordRequestedEvent struct {
	RoomID string `json:"roomId"`
	Bot
}
```
oglnie mozna zrobic typ bota w domenie gdzie tego bota sobie przygotujemy

```go
type Bot struct {
	ID string `json:"botId"`
	Name string `json:"botName"` // moze byc generowane z fakera
	Tag string `json:"botTag"`
}
```
exchange name = meeting


4. gateway obsluguje dodanie i generuje powiadomienie -> nie potrzeba dto mozna dac imie bota
title:
Rozpoczeto nagrywanie spotkania
body:
Bot Alan dolaczyl do spotaknia i zajmie sie nagrywaniem rozmowy


CZYLI
gateway wtedy dodaje bota do pokoju i wysyla event z powidomieneim
emituje informacje o room userach

5. userzy ogarniaja offery itp

6. po ogarnienciu tych eventow bot moze nasluchiwac rozmo
- z racji ze webtrc dziala tak ze user wysyla swoje audio do kazdego to mamy nagrania per user - nie problem nawet lepiej

approach:
zapis rozmowy kazdego usera osobno
a potem sklejenie w jedna rozmowe ffmpegiem

7 po meetingu i zmergowaniu nagrania mozna generowac event dla
- powiadomienia ze user ma dostep do nagrania mozna se je odsluchac (w pewnym sensie notatka)
- evnet dla note-service ze moze przetwarzac

Zapis w bazie tylko path'y
loklanie rozmowy **.ogg**

do impl bedzie jakis get do rozmow czy czegos takiego
no mozna miec jakies metadane dot. rozmowy czyli date startu i zakonczenia, tytul itp reszta juz bedzie nalezala do ai kto bral udzial itp.
metadane w sumie moga byc w jakims postgresie

nagrania

/recordings
  /{meetingid}
    merged.ogg
    {user1uid}.ogg
    ...
