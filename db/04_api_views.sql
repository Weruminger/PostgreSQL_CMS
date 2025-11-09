create or replace view api_published_articles as
select e.id, e.slug, e.published_at, e.data, d.host as domain, e.tenant_id
from entries e
join content_types t on t.id = e.type_id
join domains d on d.tenant_id = e.tenant_id and d.primary_domain
where t.slug='article' and e.status='published';
grant select on api_published_articles to web_anon;
