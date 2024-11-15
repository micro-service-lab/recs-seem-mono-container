# public.m_chat_room_action_types

## Description

## Columns

| Name | Type | Default | Nullable | Children | Parents | Comment |
| ---- | ---- | ------- | -------- | -------- | ------- | ------- |
| m_chat_room_action_types_pkey | bigint | nextval('m_chat_room_action_types_m_chat_room_action_types_pkey_seq'::regclass) | false |  |  |  |
| chat_room_action_type_id | uuid | uuid_generate_v4() | false | [public.t_chat_room_actions](public.t_chat_room_actions.md) |  |  |
| name | varchar(255) |  | false |  |  |  |
| key | varchar(255) |  | false |  |  |  |

## Constraints

| Name | Type | Definition |
| ---- | ---- | ---------- |
| m_chat_room_action_types_pkey | PRIMARY KEY | PRIMARY KEY (m_chat_room_action_types_pkey) |

## Indexes

| Name | Definition |
| ---- | ---------- |
| m_chat_room_action_types_pkey | CREATE UNIQUE INDEX m_chat_room_action_types_pkey ON public.m_chat_room_action_types USING btree (m_chat_room_action_types_pkey) |
| idx_m_chat_room_action_types_id | CREATE UNIQUE INDEX idx_m_chat_room_action_types_id ON public.m_chat_room_action_types USING btree (chat_room_action_type_id) |
| idx_m_chat_room_action_types_key | CREATE UNIQUE INDEX idx_m_chat_room_action_types_key ON public.m_chat_room_action_types USING btree (key) |

## Relations

![er](public.m_chat_room_action_types.svg)

---

> Generated by [tbls](https://github.com/k1LoW/tbls)
