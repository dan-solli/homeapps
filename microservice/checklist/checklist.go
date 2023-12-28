package main

type CheckList struct {
	items []CheckListItem
	name  string
}

func (c *CheckList) New(name string) CheckList {
	c.items = []CheckListItem{}
	c.name = name
	return *c
}

func (c *CheckList) add(item CheckListItem) {
	c.items = append(c.items, item)
}

func (c *CheckList) AddTodo(todo string) string {
	ci := NewCheckListItem().
		WithTitle(todo)
	c.add(*ci)

	return ci.title
}

func (c *CheckList) CountTodo() int {
	return len(c.items)
}

type BlockType uint8

const (
	BLOCKED BlockType = iota
	BLOCKING
)

type Blocks struct {
	internal_id string
	kind        BlockType
}

func (c *CheckListItem) WithBlock(kind BlockType, internal_id string) *CheckListItem {
	c.blocks = append(c.blocks, Blocks{internal_id: internal_id, kind: kind})
	return c
}

type CheckListItem struct {
	title       string
	internal_id string // Used for database
	parent_id   string // Used for database

	external_id string // Used for API/client

	geolocation *Geolocation
	storage     *Storage
	priority    *Priority
	raci        *Raci

	deadline   []Deadlines
	recurrance []Recurrence
	blocks     []Blocks
}

func NewCheckListItem() *CheckListItem {
	return &CheckListItem{}
}

func (c *CheckListItem) GetTitle() string {
	return c.title
}

func (c *CheckListItem) WithTitle(title string) *CheckListItem {
	c.title = title
	return c
}

func (c *CheckListItem) WithGeolocation(geolocation *Geolocation) *CheckListItem {
	c.geolocation = geolocation
	return c
}

func (c *CheckListItem) WithStorage(storage *Storage) *CheckListItem {
	c.storage = storage
	return c
}

func (c *CheckListItem) Store() (string, error) {
	return c.internal_id, nil
}

func (c *CheckListItem) Fetch(int_id string) (string, error) {
	return int_id, nil
}

type Geolocation struct {
	external_id string
}

type Storage struct {
	external_id string
}

type DateTime struct {
	year  uint16
	month uint8
	day   uint8
	hour  uint8
	min   uint8
	sec   uint8
}

type DeadlineType uint8

const (
	SOFT DeadlineType = iota
	HARD
)

type Deadlines struct {
	deadline DateTime
	kind     DeadlineType
}

type Recurrence struct {
	fancyPattern string
}

func (c *CheckListItem) WithDeadline(deadline DateTime, kind DeadlineType) *CheckListItem {
	c.deadline = append(c.deadline, Deadlines{deadline: deadline, kind: kind})
	return c
}

func (c *CheckListItem) WithRecurrence(fancyPattern string) *CheckListItem {
	c.recurrance = append(c.recurrance, Recurrence{fancyPattern: fancyPattern})
	return c
}

type Priority uint8

const (
	LOW Priority = iota
	MEDIUM
	HIGH
)

func (c *CheckListItem) WithPriority(priority Priority) *CheckListItem {
	c.priority = &priority
	return c
}

type Party string // external_id against party microservice

type Raci struct {
	creator     Party
	responsible Party
	accountable Party
	consulted   []Party
	informed    []Party
}

func (c *CheckListItem) WithRaci(raci *Raci) *CheckListItem {
	c.raci = raci
	return c
}
