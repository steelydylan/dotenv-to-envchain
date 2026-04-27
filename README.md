# env-to-envchain

`.env` ファイルの内容を [envchain](https://github.com/sorah/envchain) に namespace 付きで保存する CLI ツールです。

## 前提条件

- Go 1.21+
- [envchain](https://github.com/sorah/envchain) がインストール済みであること

```bash
brew install envchain
```

## インストール

```bash
go install github.com/steelydylan/env-to-envchain@latest
```

## 使い方

```bash
env-to-envchain <namespace> <envfile>
```

### 例

```bash
# .env の内容を "myapp" namespace に保存
env-to-envchain myapp .env

# 保存した変数を確認
envchain myapp env

# 保存した変数を使ってコマンドを実行
envchain myapp rails server
```

### 対応する .env フォーマット

```env
# コメント行（無視されます）
DATABASE_URL=postgres://localhost:5432/mydb
REDIS_URL="redis://localhost:6379"
API_KEY='sk-test-12345'
export SECRET_TOKEN=abc123
```

- `#` で始まる行はスキップ
- `export` プレフィックスに対応
- ダブルクォート・シングルクォートの値に対応
- 空行はスキップ
