# GIF Converter

動画ファイルをGIFアニメーションに変換するデスクトップアプリケーション

## 📋 概要

このアプリケーションは、Wails フレームワークを使用して構築されたクロスプラットフォーム対応の動画→GIF変換ツールです。直感的なドラッグ&ドロップUIを提供し、FFmpegを活用して高品質な変換を実現します。

## 🛠️ 技術スタック

- **フロントエンド**: React + TypeScript + Vite
- **バックエンド**: Go
- **フレームワーク**: Wails v2
- **動画処理**: FFmpeg (ffmpeg-go ライブラリ)
- **UI**: カスタムCSS + モダンデザイン

## ✨ 機能

- 📁 ドラッグ&ドロップでの動画ファイルアップロード
- 🎬 多様な動画フォーマットをサポート (.mp4, .mov, .avi, .mkv, .webm など)
- 🔄 FFmpegを使用した高品質な変換処理
- 📂 出力ファイルの自動整理 (~/Documents/GIF-Converter/)
- 🖥️ クロスプラットフォーム対応 (macOS, Windows, Linux)
- 📊 変換進捗とエラー表示
- 🛠️ デバッグ機能 (FFmpeg情報確認)

## 🚀 セットアップ

### 前提条件

1. **Go**: 1.21以上
2. **Node.js**: 18以上
3. **Wails CLI**: v2.10.1以上
4. **FFmpeg**: システムにインストール済み

### FFmpegのインストール

```bash
# macOS (Homebrew)
brew install ffmpeg

# Windows (Chocolatey)
choco install ffmpeg

# Linux (Ubuntu/Debian)
sudo apt update && sudo apt install ffmpeg
```

### プロジェクトのセットアップ

```bash
# リポジトリをクローン
git clone <repository-url>
cd gif-converter

# 依存関係のインストール
go mod tidy

# フロントエンドの依存関係をインストール
cd frontend && npm install && cd ..
```

## 🎯 使用方法

### 開発モード

```bash
# 開発サーバーを起動
wails dev
```

- ホットリロード対応
- デバッグツール利用可能
- ブラウザ版: http://localhost:34115

### 本番ビルド

```bash
# 本番用パッケージをビルド
wails build
```

ビルドされたアプリケーションは `build/bin/` に生成されます。

## 📁 プロジェクト構成

```
gif-converter/
├── app.go                 # メインのGoアプリケーション
├── main.go               # エントリーポイント
├── go.mod                # Go依存関係
├── wails.json            # Wails設定
├── frontend/             # Reactフロントエンド
│   ├── src/
│   │   ├── App.tsx       # メインコンポーネント
│   │   ├── App.css       # スタイル
│   │   └── main.tsx      # エントリーポイント
│   ├── wailsjs/          # Wails自動生成バインディング
│   └── package.json      # npm依存関係
├── build/                # ビルド出力
└── output/               # 開発時の出力フォルダ
```

## 🔧 主要な機能詳細

### 動画変換処理

- **変換ライブラリ**: github.com/u2takey/ffmpeg-go
- **デフォルト設定**: 
  - フレームレート: 10fps
  - サイズ: 元動画の縦横比維持
  - 品質: lanczos フィルター
- **対応フォーマット**: FFmpegでサポートされる全ての動画フォーマット

### ファイル管理

- **一時ファイル**: システムの一時ディレクトリに保存
- **出力先**: `~/Documents/GIF-Converter/`
- **ファイル名**: `元ファイル名_YYYYMMDD_HHMMSS.gif`
- **自動削除**: 変換後に一時ファイルは自動削除

### エラーハンドリング

- FFmpegパスの自動検出
- 詳細なエラーメッセージ表示
- 変換失敗時の適切な処理

## 🎨 UI機能

### ファイルアップロード

- ドラッグ&ドロップエリア
- ファイル選択ダイアログ
- ファイル情報表示 (ファイル名、サイズ)
- 動画ファイルの形式チェック

### 操作ボタン

- **GIFに変換**: メインの変換処理
- **FFmpeg情報を確認**: デバッグ情報表示
- **出力フォルダを確認**: 出力先パスを表示
- **出力フォルダを開く**: ファイルマネージャーで開く

## 🛠️ 開発情報

### 主要な依存関係

**Go**:
- `github.com/wailsapp/wails/v2`
- `github.com/u2takey/ffmpeg-go`

**Frontend**:
- `react` + `react-dom`
- `typescript`
- `vite`

### API関数

- `ConvertUploadedVideoToGif(fileData, fileName)`: メイン変換処理
- `GetFFmpegInfo()`: FFmpeg情報取得
- `GetOutputDir()`: 出力ディレクトリパス取得
- `OpenOutputDir()`: 出力ディレクトリを開く

## 🤖 開発支援

このプロジェクトは GitHub Copilot を活用して開発されました。主に以下の領域でAIアシスタントを使用：

- Go言語でのffmpeg-goライブラリ統合
- React + TypeScriptでの型安全なUI開発
- エラーハンドリングとデバッグ機能の実装
- クロスプラットフォーム対応のファイルパス処理
- Wails フレームワークの活用とバインディング生成

## 📝 開発履歴

- 動画→GIF変換の基本機能実装
- ドラッグ&ドロップUI対応
- FFmpegバイナリパスの自動検出
- 出力ファイル管理の改善
- デバッグ機能の追加
- クロスプラットフォーム対応

