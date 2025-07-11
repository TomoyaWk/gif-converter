package main

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"time"

	ffmpeg "github.com/u2takey/ffmpeg-go"
)

// App struct
type App struct {
	ctx context.Context
}

// NewApp creates a new App application struct
func NewApp() *App {
	return &App{}
}

// startup is called when the app starts. The context is saved
// so we can call the runtime methods
func (a *App) startup(ctx context.Context) {
	a.ctx = ctx
}

// GifConvertOptions GIF変換のオプション設定
type GifConvertOptions struct {
	FPS     int    // フレームレート（デフォルト: 10）
	Width   int    // 幅（デフォルト:-1で元の比率維持）
	Height  int    // 高さ（デフォルト: -1で元の比率維持）
	Quality string // 品質フラグ（デフォルト: "lanczos"）
}

// NewGifConvertOptions デフォルトオプションを返す
func NewGifConvertOptions() *GifConvertOptions {
	return &GifConvertOptions{
		FPS:     10,
		Width:   -1, // 幅は比率維持
		Height:  -1, // 高さは比率維持
		Quality: "lanczos",
	}
}

// ConvertVideoToGif 動画をGIFに変換する
func (a *App) ConvertVideoToGif(videoPath string, options *GifConvertOptions) (string, error) {
	ext := filepath.Ext(videoPath)
	tempOutputPath := strings.TrimSuffix(videoPath, ext) + ".gif"

	// ffmpegの実行ファイルパスを取得
	ffmpegPath := getFFmpegPath()

	// ffmpeg-goを使用してGIF変換
	err := ffmpeg.Input(videoPath).
		Output(tempOutputPath, ffmpeg.KwArgs{
			"vf": fmt.Sprintf("fps=%d,scale=%d:%d", options.FPS, options.Width, options.Height),
			"f":  "gif",
		}).
		OverWriteOutput().
		SetFfmpegPath(ffmpegPath).
		Run()

	if err != nil {
		return "", fmt.Errorf("GIF変換エラー (ffmpeg path: %s): %v", ffmpegPath, err)
	}

	// outputフォルダに移動
	outputPath, err := a.moveToOutputFolder(tempOutputPath)
	if err != nil {
		return "", fmt.Errorf("outputフォルダへの移動エラー: %v", err)
	}

	return outputPath, nil
}

// SaveUploadedFile アップロードされたファイルを一時保存し、パスを返す
func (a *App) SaveUploadedFile(fileData []byte, fileName string) (string, error) {
	tempDir := filepath.Join(os.TempDir(), "gif-converter")
	if err := os.MkdirAll(tempDir, 0755); err != nil {
		return "", fmt.Errorf("一時ディレクトリ作成エラー: %v", err)
	}

	timestamp := time.Now().Format("20060102_150405")
	tempFilePath := filepath.Join(tempDir, fmt.Sprintf("%s_%s", timestamp, fileName))

	if err := os.WriteFile(tempFilePath, fileData, 0644); err != nil {
		return "", fmt.Errorf("ファイル保存エラー: %v", err)
	}

	return tempFilePath, nil
}

// GetFFmpegInfo ffmpegの情報を取得する（デバッグ用）
func (a *App) GetFFmpegInfo() string {
	ffmpegPath := getFFmpegPath()
	info := fmt.Sprintf("FFmpeg path: %s\n", ffmpegPath)

	// ファイルの存在確認
	if _, err := os.Stat(ffmpegPath); err == nil {
		info += "Status: File exists\n"

		// バージョン確認
		if cmd := exec.Command(ffmpegPath, "-version"); cmd != nil {
			if output, err := cmd.Output(); err == nil {
				lines := strings.Split(string(output), "\n")
				if len(lines) > 0 {
					info += fmt.Sprintf("Version: %s\n", lines[0])
				}
			}
		}
	} else {
		info += fmt.Sprintf("Status: File not found (%v)\n", err)
	}

	return info
}

// ConvertUploadedVideoToGif アップロードされたファイルデータを直接GIFに変換する
func (a *App) ConvertUploadedVideoToGif(fileData []byte, fileName string) (string, error) {
	// 一時ファイルに保存
	tempFilePath, err := a.SaveUploadedFile(fileData, fileName)
	if err != nil {
		return "", fmt.Errorf("ファイル保存エラー: %v", err)
	}

	// 一時ファイルを削除する処理をdeferで登録
	defer func() {
		if err := os.Remove(tempFilePath); err != nil {
			fmt.Printf("一時ファイル削除エラー: %v\n", err)
		}
	}()

	// デフォルトオプションでGIF変換
	options := NewGifConvertOptions()
	return a.ConvertVideoToGif(tempFilePath, options)
}

// GetOutputDir 出力ディレクトリのパスを取得する（フロントエンド用）
func (a *App) GetOutputDir() string {
	return getOutputDir()
}

// OpenOutputDir 出力ディレクトリをファイルマネージャーで開く
func (a *App) OpenOutputDir() error {
	outputDir := getOutputDir()

	// ディレクトリが存在しない場合は作成
	if err := os.MkdirAll(outputDir, 0755); err != nil {
		return fmt.Errorf("outputディレクトリ作成エラー: %v", err)
	}

	var cmd *exec.Cmd
	switch runtime.GOOS {
	case "darwin":
		cmd = exec.Command("open", outputDir)
	case "windows":
		cmd = exec.Command("explorer", outputDir)
	case "linux":
		cmd = exec.Command("xdg-open", outputDir)
	default:
		return fmt.Errorf("サポートされていないOS: %s", runtime.GOOS)
	}

	return cmd.Run()
}

// getOutputDir 出力ディレクトリのパスを取得する
func getOutputDir() string {
	// 1. ユーザーのホームディレクトリを取得
	homeDir, err := os.UserHomeDir()
	if err != nil {
		// フォールバック: 現在のディレクトリのoutputフォルダ
		return "output"
	}

	// 2. アプリケーション専用のフォルダを作成
	// macOSの場合: ~/Documents/GIF-Converter
	// Windowsの場合: %USERPROFILE%\Documents\GIF-Converter
	// Linuxの場合: ~/Documents/GIF-Converter
	var outputDir string
	switch runtime.GOOS {
	case "darwin":
		outputDir = filepath.Join(homeDir, "Documents", "GIF-Converter")
	case "windows":
		outputDir = filepath.Join(homeDir, "Documents", "GIF-Converter")
	case "linux":
		outputDir = filepath.Join(homeDir, "Documents", "GIF-Converter")
	default:
		outputDir = filepath.Join(homeDir, "Documents", "GIF-Converter")
	}

	return outputDir
}

// moveToOutputFolder ファイルをoutputフォルダに移動
func (a *App) moveToOutputFolder(filePath string) (string, error) {
	outputDir := getOutputDir()
	if err := os.MkdirAll(outputDir, 0755); err != nil {
		return "", fmt.Errorf("outputディレクトリ作成エラー: %v", err)
	}

	fileName := filepath.Base(filePath)
	timestamp := time.Now().Format("20060102_150405")
	ext := filepath.Ext(fileName)
	nameWithoutExt := strings.TrimSuffix(fileName, ext)
	finalFileName := fmt.Sprintf("%s_%s%s", nameWithoutExt, timestamp, ext)
	finalPath := filepath.Join(outputDir, finalFileName)

	if err := os.Rename(filePath, finalPath); err != nil {
		// リネームが失敗した場合はコピー＆削除
		if err := a.copyFile(filePath, finalPath); err != nil {
			return "", fmt.Errorf("ファイルコピーエラー: %v", err)
		}
		os.Remove(filePath)
	}

	absPath, err := filepath.Abs(finalPath)
	if err != nil {
		return finalPath, nil
	}

	return absPath, nil
}

// copyFile ファイルをコピー
func (a *App) copyFile(src, dst string) error {
	sourceFile, err := os.Open(src)
	if err != nil {
		return err
	}
	defer sourceFile.Close()

	destFile, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer destFile.Close()

	_, err = destFile.ReadFrom(sourceFile)
	return err
}

// getFFmpegPath ffmpegの実行ファイルパスを取得する
func getFFmpegPath() string {
	// 1. PATH環境変数から検索
	if path, err := exec.LookPath("ffmpeg"); err == nil {
		return path
	}

	// 2. 一般的なインストール場所をチェック
	commonPaths := []string{
		"/opt/homebrew/bin/ffmpeg", // macOS Homebrew (Apple Silicon)
		"/usr/local/bin/ffmpeg",    // macOS Homebrew (Intel)
		"/usr/bin/ffmpeg",          // Linux
		"./bin/ffmpeg",             // アプリケーションと同じディレクトリのbinフォルダ
	}

	// macOSの場合
	if runtime.GOOS == "darwin" {
		commonPaths = append([]string{
			"/opt/homebrew/bin/ffmpeg",
			"/usr/local/bin/ffmpeg",
		}, commonPaths...)
	}

	for _, path := range commonPaths {
		if _, err := os.Stat(path); err == nil {
			return path
		}
	}

	// デフォルトで"ffmpeg"を返す（PATH環境変数に依存）
	return "ffmpeg"
}
