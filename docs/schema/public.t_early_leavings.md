# public.t_early_leavings

## Description

## Columns

| Name | Type | Default | Nullable | Children | Parents | Comment |
| ---- | ---- | ------- | -------- | -------- | ------- | ------- |
| t_early_leavings_pkey | bigint | nextval('t_early_leavings_t_early_leavings_pkey_seq'::regclass) | false |  |  |  |
| early_leaving_id | uuid | uuid_generate_v4() | false |  |  |  |
| attendance_id | uuid |  | false |  | [public.t_attendances](public.t_attendances.md) |  |
| leave_time | timestamp with time zone |  | false |  |  |  |

## Constraints

| Name | Type | Definition |
| ---- | ---- | ---------- |
| fk_t_early_leavings_attendance_id | FOREIGN KEY | FOREIGN KEY (attendance_id) REFERENCES t_attendances(attendance_id) ON UPDATE RESTRICT ON DELETE RESTRICT |
| t_early_leavings_pkey | PRIMARY KEY | PRIMARY KEY (t_early_leavings_pkey) |

## Indexes

| Name | Definition |
| ---- | ---------- |
| t_early_leavings_pkey | CREATE UNIQUE INDEX t_early_leavings_pkey ON public.t_early_leavings USING btree (t_early_leavings_pkey) |
| idx_t_early_leavings_id | CREATE UNIQUE INDEX idx_t_early_leavings_id ON public.t_early_leavings USING btree (early_leaving_id) |

## Relations

![er](public.t_early_leavings.svg)

---

> Generated by [tbls](https://github.com/k1LoW/tbls)
