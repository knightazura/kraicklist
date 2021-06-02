package vendors

import (
	"encoding/json"
	"reflect"
	"strconv"

	"github.com/go-ego/riot"
	"github.com/go-ego/riot/types"
	"github.com/knightazura/contracts"
	"github.com/knightazura/domain"
	"github.com/knightazura/utils"
)

var (
	se      = riot.Engine{}
	options = types.RankOpts{
		ScoringCriteria: AdvertisementScoringCriteria{},
		OutputOffset:    0,
		MaxOutputs:      100,
	}
)

const MaxTokenProximity = 150

type Riot struct {
	Logger *utils.Logger
	Client *riot.Engine
}

func InitRiot() contracts.SearchEngine {
	// Init. & configure Riot
	se.Init(types.EngineOpts{
		Using:   1,
		GseMode: true,
		IndexerOpts: &types.IndexerOpts{
			IndexType: types.LocsIndex,
		},
		DefRankOpts: &options,
	})

	return &Riot{
		Logger: utils.InitLogger(),
		Client: &se,
	}
}

func (r *Riot) Search(indexName string, query string) domain.SearchedDocument {
	var ads []domain.Advertisement
	var hits []interface{}

	records := r.Client.Search(types.SearchReq{Text: query})
	br, _ := json.Marshal(records.Docs)
	_ = json.Unmarshal(br, &ads)

	for _, adv := range ads {
		var ad domain.Advertisement
		_ = json.Unmarshal([]byte(adv.Content), &ad)

		ads = append(ads, ad)
		hits = append(hits, ad)
	}
	defer r.Client.Close()

	r.Logger.LogAccess("Riot: Search %s in %s documents. Found %d", query, indexName, len(ads))

	return domain.SearchedDocument{
		Hits:      hits,
		TotalHits: int64(len(ads)),
		Query:     query,
	}
}

func (r *Riot) Add(doc *domain.GeneralDocument, indexName string) {
	ad := doc.Data.(*domain.Advertisement)
	fields := AdvertisementScoringFields{
		Title:   ad.Title,
		Content: ad.Content,
	}

	fj, _ := json.Marshal(fields)
	r.Client.Index(strconv.FormatInt(ad.ID, 10), types.DocData{
		Content: string(fj),
		Fields:  fields,
	})
	defer r.Client.Close()

	r.Client.Flush()

	r.Logger.LogAccess("Riot: Success add new %s, %s document", ad.Title, indexName)

	return
}

func (r *Riot) BulkInsert(docs *domain.GeneralDocuments, indexName string) {
	for _, doc := range *docs {
		ad := doc.Data.(domain.Advertisement)

		fields := AdvertisementScoringFields{
			Title:   ad.Title,
			Content: ad.Content,
		}

		fj, _ := json.Marshal(fields)
		r.Client.Index(strconv.FormatInt(ad.ID, 10), types.DocData{
			Content: string(fj),
			Fields:  fields,
		})
	}
	defer r.Client.Close()

	r.Client.Flush()

	r.Logger.LogAccess("Riot: Bulk insert %s documents", indexName)

	return
}

func (r *Riot) DeleteDocument(docID string, indexName string) {
	r.Client.RemoveDoc(docID)
	defer r.Client.Close()

	r.Client.Flush()
}

func (r *Riot) DeleteIndex(indexName string) {
	r.Logger.LogAccess("Riot: There's no remove delete index in Riot: %s index", indexName)
}

func (r *Riot) TotalDocuments(indexName string) int64 {
	return int64(0)
}

type AdvertisementScoringFields struct {
	Title   string
	Content string
}

type AdvertisementScoringCriteria struct{}

func (crit AdvertisementScoringCriteria) Score(doc types.IndexedDoc, fields interface{}) []float32 {
	if doc.TokenProximity > MaxTokenProximity {
		return []float32{}
	}
	if reflect.TypeOf(fields) != reflect.TypeOf(AdvertisementScoringFields{}) {
		return []float32{}
	}

	output := make([]float32, 3)
	asf := fields.(AdvertisementScoringFields)
	title, _ := strconv.ParseFloat(asf.Title, 32)
	content, _ := strconv.ParseFloat(asf.Content, 32)
	output[0] = float32(title)
	output[1] = float32(content)
	output[2] = float32(int(doc.BM25))
	return output
}
