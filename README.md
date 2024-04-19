# recs-seem-mono-container

## Getting started

It runs on Linux environment.

```shell
git clone https://github.com/micro-service-lab/recs-seem-mono-container.git
cd recs-seem-mono-container
aqua policy allow "${PWD}/aqua-policy.yaml"
aqua i -l
```

## mageコマンド

- コンテナ操作
    - up
	``` sh
	mage app:up
	```
    - down
	``` sh
	mage app:down
	```
    - ps
	``` sh
	mage app:ps
	```
    - log
	``` sh
	mage app:log {{サービス名}}
	```
    - clean
	``` sh
	mage app:clean
	```
    - bash
	``` sh
	mage app:bash {{サービス名}}
	```
    - migrate
	``` sh
	mage app:migrate
	```
    - rollback
	``` sh
	mage app:rollback
	```
    - seed
	``` sh
	mage app:seed
	```
	- dev
	``` sh
	mage app:dev
	```
	- kill
	``` sh
	mage app:kill
	```
	- serve
	``` sh
	mage app:serve
	```

- マイグレーションファイル生成
	``` sh
	mage generate:migration {{ファイル名}}
	```
    - 使用例

	``` sh
	mage generate:migration create_m_mime_types_table
	```
