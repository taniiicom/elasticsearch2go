# elasticsearch2go

[English](https://github.com/taniiicom/elasticsearch2go/blob/main/README.md)
| [日本語](https://github.com/taniiicom/elasticsearch2go/blob/main/README.ja.md)

`elasticsearch2go`は, Elasticsearch のマッピング定義を基に Go の構造体を自動生成するツールです. Elasticsearch の JSON スキーマから Go の構造体を生成し, コードの自動化と一貫性の維持を支援します.

## インストール方法 (Installation)

このパッケージをインストールするには, 以下のコマンドを使用してください.

```bash
go get github.com/taniiicom/elasticsearch2go
```

## 使用方法 (Usage)

このパッケージを使用して, Elasticsearch のマッピング定義から Go の構造体を生成する方法を以下に示します.

### コマンドラインでの使用例

```bash
elasticsearch2go --in=mapping.json --out=model.go --package=mypackage --struct=MyStruct
```

### コードでの使用例

```go
package main

import (
    "log"
    "github.com/taniiicom/elasticsearch2go"
)

func main() {
    err := elasticsearch2go.GenerateStructs(
        "mapping.json",
        "model.go",
        "mypackage",
        "MyStruct",
        "MyWrapperStruct",
        "custom_mapping.json",
        "custom_field_exceptions.json",
        "custom_type_exceptions.json",
        "skip_fields.json",
        "field_comments.json",
        "custom_template.tmpl",
    )
    if err != nil {
        log.Fatalf("Failed to generate structs: %v", err)
    }
}
```

## コマンドラインオプション (Command-line Options)

このパッケージには, 以下のコマンドラインオプションがあります.

- `--in`: **(必須)** 入力の JSON スキーマファイルのパスを指定します.
- `--out`: **(必須)** 出力先の Go ファイルのパスを指定します.
- `--package`: **(必須)** 出力される Go ファイルのパッケージ名を指定します.
- `--struct`: **(必須)** 生成される構造体の名前を指定します.
- `--init`: 初期化用のラッパー構造体の名前を指定します（オプション）.
- `--type-mapping`: Elasticsearch の型を Go の型にマッピングする JSON ファイルのパスを指定します（オプション）.
- `--exception-field`: field 名の例外を定義する JSON ファイルのパスを指定します（オプション）.
- `--exception-type`: 型の例外を定義する JSON ファイルのパスを指定します（オプション）.
- `--skip-field`: 生成からスキップする field を定義する JSON ファイルのパスを指定します（オプション）.
- `--field-comment`: field にコメントを追加するための JSON ファイルのパスを指定します（オプション）.
- `--tmpl`: カスタム Go テンプレートファイルのパスを指定します（オプション）.

## カスタマイズ方法 (Customization)

このパッケージは, 様々なカスタマイズが可能です. 以下に, いくつかのカスタマイズ方法を示します.

### カスタム型マッピングの使用

Elasticsearch の型を特定の Go 型にマッピングするためには, `custom_mapping.json`ファイルを使用します. 例:

```json
{
  "text": "*string",
  "integer": "int"
}
```

これにより, `text`型が`*string`に, `integer`型が`int`にマッピングされます.

### field 名の例外設定

field 名に特定の変換を適用するためには, `custom_field_exceptions.json`ファイルを使用します. 例:

```json
{
  "my_field": "MyCustomField"
}
```

これにより, `my_field`が`MyCustomField`という field 名に変換されます.

### スキップ field の設定

特定の field を生成からスキップするためには, `skip_fields.json`ファイルを使用します. 例:

```json
{
  "unnecessary_field": true
}
```

これにより, `unnecessary_field`が構造体から除外されます.

## 機能 (Features)

### ネストされた構造体の生成に対応

Elasticsearch で`"type": "nested"`として定義された field は, Go の構造体生成時にサブ構造体として扱われます. このサブ構造体は, メインの構造体の下に, 同じ出力ファイル内で生成されます. これにより, ネストされた構造をそのまま Go のコードに反映することができます.

### 柔軟なカスタマイズ対応

elasticsearch -> go の, { field 名, type } の例外的な対応を JSON ファイルで指定できます. 例外的な対応は, field 名と型の両方に適用できます.

また, 特定の property をスキップするための JSON ファイルを指定できます. これにより, 不要な field を生成から除外できます.

### field のコメントの追加

elasticsearch の field にコメントを追加するための JSON ファイルを指定できます. この機能を使用すると, 生成された Go の構造体にコメントを追加できます.

### カスタムテンプレートの使用

生成するファイルのフォーマットをカスタマイズするために, カスタムテンプレートファイルを指定できます. この機能を使用すると, 出力される Go のコードのフォーマットを制御できます.

### [todo] 構造体の field 順序の指定

現時点では, 生成される構造体の field 順は, アルファベット順です. 今後のバージョンで, field 順序を, オリジナルの Elasticsearch の定義を維持できるようになる予定です.

## 貢献方法 (Contributing)

このパッケージに対する貢献は歓迎されます. 貢献するには, 以下の手順に従ってください.

1. このリポジトリを Fork します.
2. 新しいブランチを作成します (`git checkout -b feature/your-feature-name`).
3. 変更をコミットします (`git commit -m 'Add some feature'`).
4. ブランチにプッシュします (`git push origin feature/your-feature-name`).
5. GitHub でプルリクエストを作成します.

## ライセンス (License)

このプロジェクトは MIT ライセンスの下で公開されています. 詳細は `LICENSE` ファイルを参照してください.

## メンテナンスとサポート (Maintenance and Support)

このパッケージは現在, @taniiicom (Taniii.com) によってメンテナンスされています. 問題が発生した場合は, GitHub の Issues を使用して報告してください.
