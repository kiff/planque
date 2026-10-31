// Package store keeps planque's records: a small key-value table with
// conditional writes, which is all planque needs from DynamoDB. Memory is
// the same contract in one process, for tests and a self-hosted trial.
package store

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"errors"
	"slices"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/kiff/planque/internal/awsapi"
)

// Item is one record: a JSON document under (PK, SK), with a version for
// optimistic writes, an optional lease (unix ms) and an expiry (unix s).
type Item struct {
	PK, SK string
	Doc    []byte
	Ver    int64
	Lease  int64
	TTL    int64
}

// Cond is what must hold for a write to happen.
type Cond struct {
	// NotExists: no item has this key.
	NotExists bool
	// Ver: the stored item has this version.
	Ver *int64
	// LeaseFreeAt: the stored item's lease is absent or before this time
	// (unix ms).
	LeaseFreeAt int64
}

// ErrCondition: the write's condition did not hold.
var ErrCondition = errors.New("store: condition failed")

// KV is the table.
type KV interface {
	Get(ctx context.Context, pk, sk string) (Item, bool, error)
	Put(ctx context.Context, it Item, c Cond) error
	// Delete removes the item and returns it as it was, if it existed.
	Delete(ctx context.Context, pk, sk string) (Item, bool, error)
	// Query lists the items under pk whose SK starts with prefix.
	Query(ctx context.Context, pk, prefix string) ([]Item, error)
	// QueryTo lists the items under pk whose SK is at most maxSK, in SK
	// order: reading a range never costs more than the range.
	QueryTo(ctx context.Context, pk, maxSK string) ([]Item, error)
}

// Memory is KV in one process.
type Memory struct {
	mu    sync.Mutex
	items map[[2]string]Item
	Now   func() time.Time
}

func (m *Memory) now() time.Time {
	if m.Now != nil {
		return m.Now()
	}
	return time.Now()
}

func (m *Memory) live(it Item) bool { return it.TTL == 0 || it.TTL > m.now().Unix() }

func (m *Memory) Get(_ context.Context, pk, sk string) (Item, bool, error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	it, ok := m.items[[2]string{pk, sk}]
	if !ok || !m.live(it) {
		return Item{}, false, nil
	}
	it.Doc = slices.Clone(it.Doc)
	return it, true, nil
}

func (m *Memory) Put(_ context.Context, it Item, c Cond) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	if m.items == nil {
		m.items = map[[2]string]Item{}
	}
	old, ok := m.items[[2]string{it.PK, it.SK}]
	ok = ok && m.live(old)
	if c.NotExists && ok {
		return ErrCondition
	}
	if c.Ver != nil && (!ok || old.Ver != *c.Ver) {
		return ErrCondition
	}
	if c.LeaseFreeAt != 0 && ok && old.Lease != 0 && old.Lease >= c.LeaseFreeAt {
		return ErrCondition
	}
	it.Doc = slices.Clone(it.Doc)
	m.items[[2]string{it.PK, it.SK}] = it
	return nil
}

func (m *Memory) Delete(_ context.Context, pk, sk string) (Item, bool, error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	it, ok := m.items[[2]string{pk, sk}]
	delete(m.items, [2]string{pk, sk})
	return it, ok && m.live(it), nil
}

func (m *Memory) Query(_ context.Context, pk, prefix string) ([]Item, error) {
	return m.list(pk, func(sk string) bool { return strings.HasPrefix(sk, prefix) }), nil
}

func (m *Memory) QueryTo(_ context.Context, pk, maxSK string) ([]Item, error) {
	return m.list(pk, func(sk string) bool { return sk <= maxSK }), nil
}

func (m *Memory) list(pk string, keep func(sk string) bool) []Item {
	m.mu.Lock()
	defer m.mu.Unlock()
	var out []Item
	for k, it := range m.items {
		if k[0] == pk && keep(k[1]) && m.live(it) {
			out = append(out, it)
		}
	}
	slices.SortFunc(out, func(a, b Item) int { return strings.Compare(a.SK, b.SK) })
	return out
}

// Dynamo is KV on one DynamoDB table with string keys pk and sk.
type Dynamo struct {
	Table string
	API   *awsapi.Client
	Now   func() time.Time
}

type av map[string]any

func (d *Dynamo) key(pk, sk string) av {
	return av{"pk": av{"S": pk}, "sk": av{"S": sk}}
}

func encode(it Item) av {
	m := av{"pk": av{"S": it.PK}, "sk": av{"S": it.SK}, "doc": av{"B": it.Doc},
		"ver": av{"N": strconv.FormatInt(it.Ver, 10)}}
	if it.Lease != 0 {
		m["lease"] = av{"N": strconv.FormatInt(it.Lease, 10)}
	}
	if it.TTL != 0 {
		m["ttl"] = av{"N": strconv.FormatInt(it.TTL, 10)}
	}
	return m
}

type attr struct {
	S *string `json:"S"`
	N *string `json:"N"`
	B *string `json:"B"`
}

func decode(m map[string]attr) Item {
	var it Item
	if v := m["pk"].S; v != nil {
		it.PK = *v
	}
	if v := m["sk"].S; v != nil {
		it.SK = *v
	}
	if v := m["doc"].B; v != nil {
		it.Doc, _ = base64.StdEncoding.DecodeString(*v)
	}
	num := func(a attr) int64 {
		if a.N == nil {
			return 0
		}
		n, _ := strconv.ParseInt(*a.N, 10, 64)
		return n
	}
	it.Ver, it.Lease, it.TTL = num(m["ver"]), num(m["lease"]), num(m["ttl"])
	return it
}

func (d *Dynamo) live(it Item) bool {
	now := time.Now
	if d.Now != nil {
		now = d.Now
	}
	// DynamoDB deletes expired items lazily; an expired one is absent.
	return it.TTL == 0 || it.TTL > now().Unix()
}

// readTries is how many times a read is sent when the connection drops.
const readTries = 3

// read sends a read again when the connection dropped before DynamoDB answered
// (an EOF, a reset): reading twice changes nothing. An answer from DynamoDB,
// or the request's context ending, is final. Writes are never resent here: a
// conditional write whose answer was lost may already have landed.
func (d *Dynamo) read(ctx context.Context, action string, in, out any) error {
	var err error
	for try := 0; try < readTries; try++ {
		if err = d.API.Do(ctx, action, in, out); err == nil || ctx.Err() != nil {
			return err
		}
		var answered *awsapi.Error
		if errors.As(err, &answered) {
			return err
		}
	}
	return err
}

func (d *Dynamo) Get(ctx context.Context, pk, sk string) (Item, bool, error) {
	var out struct {
		Item map[string]attr `json:"Item"`
	}
	err := d.read(ctx, "GetItem", av{"TableName": d.Table, "Key": d.key(pk, sk), "ConsistentRead": true}, &out)
	if err != nil || out.Item == nil {
		return Item{}, false, err
	}
	it := decode(out.Item)
	return it, d.live(it), nil
}

func (d *Dynamo) Put(ctx context.Context, it Item, c Cond) error {
	in := av{"TableName": d.Table, "Item": encode(it)}
	var conds []string
	vals := av{}
	if c.NotExists {
		// An expired item counts as absent, as in Memory.
		conds = append(conds, "(attribute_not_exists(pk) OR (attribute_exists(#ttl) AND #ttl < :nows))")
		vals[":nows"] = av{"N": strconv.FormatInt(time.Now().Unix(), 10)}
		in["ExpressionAttributeNames"] = av{"#ttl": "ttl"}
	}
	if c.Ver != nil {
		conds = append(conds, "ver = :ver")
		vals[":ver"] = av{"N": strconv.FormatInt(*c.Ver, 10)}
	}
	if c.LeaseFreeAt != 0 {
		conds = append(conds, "(attribute_not_exists(lease) OR lease < :lease)")
		vals[":lease"] = av{"N": strconv.FormatInt(c.LeaseFreeAt, 10)}
	}
	if len(conds) > 0 {
		in["ConditionExpression"] = strings.Join(conds, " AND ")
		in["ExpressionAttributeValues"] = vals
	}
	err := d.API.Do(ctx, "PutItem", in, nil)
	if awsapi.Is(err, "ConditionalCheckFailedException") {
		return ErrCondition
	}
	return err
}

func (d *Dynamo) Delete(ctx context.Context, pk, sk string) (Item, bool, error) {
	var out struct {
		Attributes map[string]attr `json:"Attributes"`
	}
	err := d.API.Do(ctx, "DeleteItem", av{"TableName": d.Table, "Key": d.key(pk, sk), "ReturnValues": "ALL_OLD"}, &out)
	if err != nil || out.Attributes == nil {
		return Item{}, false, err
	}
	it := decode(out.Attributes)
	return it, d.live(it), nil
}

func (d *Dynamo) Query(ctx context.Context, pk, prefix string) ([]Item, error) {
	if prefix == "" { // DynamoDB takes no empty key value, even to begin with
		return d.query(ctx, "pk = :pk", av{":pk": av{"S": pk}})
	}
	return d.query(ctx, "pk = :pk AND begins_with(sk, :p)", av{":pk": av{"S": pk}, ":p": av{"S": prefix}})
}

func (d *Dynamo) QueryTo(ctx context.Context, pk, maxSK string) ([]Item, error) {
	return d.query(ctx, "pk = :pk AND sk <= :max", av{":pk": av{"S": pk}, ":max": av{"S": maxSK}})
}

// query reads every page of a key condition, in SK order.
func (d *Dynamo) query(ctx context.Context, cond string, values av) ([]Item, error) {
	var items []Item
	var start json.RawMessage
	for {
		in := av{"TableName": d.Table, "ConsistentRead": true,
			"KeyConditionExpression": cond, "ExpressionAttributeValues": values}
		if start != nil {
			in["ExclusiveStartKey"] = start
		}
		var out struct {
			Items            []map[string]attr `json:"Items"`
			LastEvaluatedKey json.RawMessage   `json:"LastEvaluatedKey"`
		}
		if err := d.read(ctx, "Query", in, &out); err != nil {
			return nil, err
		}
		for _, m := range out.Items {
			if it := decode(m); d.live(it) {
				items = append(items, it)
			}
		}
		if len(out.LastEvaluatedKey) == 0 || string(out.LastEvaluatedKey) == "null" {
			return items, nil
		}
		start = out.LastEvaluatedKey
	}
}
