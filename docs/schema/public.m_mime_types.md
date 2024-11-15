# public.m_mime_types

## Description

## Columns

| Name | Type | Default | Nullable | Children | Parents | Comment |
| ---- | ---- | ------- | -------- | -------- | ------- | ------- |
| m_mime_types_pkey | bigint | nextval('m_mime_types_m_mime_types_pkey_seq'::regclass) | false |  |  |  |
| mime_type_id | uuid | uuid_generate_v4() | false | [public.t_attachable_items](public.t_attachable_items.md) |  |  |
| name | varchar(255) |  | false |  |  |  |
| kind | varchar(255) |  | false |  |  |  |
| key | varchar(255) |  | false |  |  |  |

## Constraints

| Name | Type | Definition |
| ---- | ---- | ---------- |
| m_mime_types_pkey | PRIMARY KEY | PRIMARY KEY (m_mime_types_pkey) |

## Indexes

| Name | Definition |
| ---- | ---------- |
| m_mime_types_pkey | CREATE UNIQUE INDEX m_mime_types_pkey ON public.m_mime_types USING btree (m_mime_types_pkey) |
| idx_m_mime_types_id | CREATE UNIQUE INDEX idx_m_mime_types_id ON public.m_mime_types USING btree (mime_type_id) |
| idx_m_mime_types_key | CREATE UNIQUE INDEX idx_m_mime_types_key ON public.m_mime_types USING btree (key) |
| idx_m_mime_types_kind | CREATE INDEX idx_m_mime_types_kind ON public.m_mime_types USING btree (kind) |

## Relations

![er](public.m_mime_types.svg)

---

> Generated by [tbls](https://github.com/k1LoW/tbls)
