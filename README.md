# recs-seem-mono-container

## Getting started

It runs on Linux environment.

```shell
git clone https://github.com/micro-service-lab/recs-seem-mono-container.git
cd recs-seem-mono-container
aqua policy allow "${PWD}/aqua-policy.yaml"
aqua i -l
```

- devで`too many open files`が出た場合、ホスト上で以下を実行。

``` sh
echo "fs.inotify.max_user_watches = 524288" >> /etc/sysctl.conf
echo "fs.inotify.max_user_instances = 256" >> /etc/sysctl.conf
sysctl -p
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
