package internal
import "context"
type Article struct { ID int64; Slug, Title, Body string }
func LookupTenantByHost(ctx context.Context, q Querier, host string) (int64, error) {
	var id int64
	err := q.QueryRow(ctx, `select t.id from tenants t join domains d on d.tenant_id=t.id where d.host=$1 limit 1`, host).Scan(&id)
	if err != nil { return 1, err }
	return id, nil
}
func LoadArticleBySlug(ctx context.Context, q Querier, slug string) (*Article, error) {
	var id int64; var title, body string
	err := q.QueryRow(ctx, `select e.id, e.data->>'title', e.data->>'body' from api_published_articles e where e.slug=$1 limit 1`, slug).Scan(&id, &title, &body)
	if err != nil { return nil, err }
	return &Article{ID:id, Slug:slug, Title:title, Body:body}, nil
}
