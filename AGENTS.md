# Repository Guidelines

## プロジェクト構成とモジュール整理
プロトタイプは `main.go` が Ebiten のシーン遷移を束ねる Go モジュールです（`main.go` orchestrates Ebiten scenes）。ロジックは `core/`、表示まわりは `component/` と `scene/` に収容（logic lives in `core/` while UI shells stay in `component/`/`scene/`）。入出力と描画ユーティリティは `frontend/` と `widget/`、イベント制御は `sequence/`、シナリオデータは `data/`、アセットとデバッグ補助は `image/`、`font/`、`debug/` に配置されています（assets land in `image/` & `font/`, debugging helpers in `debug/`）。

## ビルド・テスト・開発コマンド
- `go run main.go` : 現在のバトルプロトタイプを起動し、手動プレイテストを行います。
- `go build ./...` : 依存を含めたビルドを実行し、実行可能ファイルをワーキングディレクトリに生成します。
- `go test ./...` : すべての Go テストを走らせ、イベントアダプタや UI ヘルパーの回帰を早期検知します。
- `go test -run Sequence ./sequence` : シーケンス周りに変更がある場合に対象テストだけを素早く回します。

## コーディングスタイルと命名規約
Go の慣習に従い、コミット前に `go fmt ./...` を必ず実行します（always `go fmt` before committing）。パッケージ名はスネークケース、公開シンボルはパスカルケース、内部用はキャメルケースで統一し、ファイルごとに責務を絞ります（keep packages snake_case, exported identifiers PascalCase, internals camelCase）。描画コードは Ebiten の draw 呼び出しを論理単位で束ね、UI 更新を追いやすくしてください（group draw calls so the frame diff stays predictable）。

## テスト指針
ユニットテストは Go の `testing` と YAML フィクスチャで構成します（tests pair `testing` with YAML fixtures）。仕様追加時は実装ファイル横の `_test.go`（例: `sequence/sequence_test.go`）にケースを積み、状態遷移と表示メッセージの両面を検証します（mirror production filenames to keep intent obvious）。

## コミットとプルリクエストの指針
コミットメッセージは「window close on submit」「battle initializeを切り出してみる」のように短い命令形で72文字以内を目安に（imperative commits under ~72 chars）。必要なら日本語で補足し、PR ではゲームへの影響、実行したテスト、UI変更多発時のスクリーンショットやGIF、関連Issueや変更した YAML を明記します（document gameplay impact, test commands, visual evidence, and YAML touchpoints）。

## データとデバッグのヒント
`data/` の YAML は起動時に読み込まれるため、スキーマを変えるときは `invoke/export.go` のヘルパーで整合性を確認してください（run the export helper when schemas move）。

## バトルシーンUI方針
`scene/battle.go` は非UIロジックを保持し、UI処理は `scene/battle_ui.go` の純粋関数群と `battleUI` 構造体で扱います（keep battle UI inside the functional helpers introduced here）。レイアウト計算は `computeBattleUILayout` のような純粋関数で行い、状態を変更する副作用は呼び出し側で明示的に処理してください（計算と副作用の分離を徹底）。`BattleScene.Update/Draw` は `battleUI.Update/Draw` を呼び出すだけに留め、UI オブジェクトは欠損しない前提で扱い、nil ガードを追加せずに異常はクラッシュで顕在化させます。UI 以外の責務（例: `BattleEventSequencer` 等）は `BattleScene` 側に保持し、UI 構造体へ押し込まないでください。新しいUIを追加する場合も同じ抽象を拡張し、関数型の流れ（計算→適用）、フレームごとの追加割り当てゼロ、そして副作用を持つ関数を極力許容しない方針を維持してください（avoid stateful helpers, keep hot paths allocation-free）。

## 実装方針

副作用のない純粋関数をベースに実装することを前提としてください。UpdateやDrawから呼び出される想定の処理の中ではメモリアロケーションの無いように注意を払ってください。
