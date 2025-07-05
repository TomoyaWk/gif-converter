import {useState} from 'react';
import './App.css';
import {ConvertUploadedVideoToGif, GetFFmpegInfo, GetOutputDir, OpenOutputDir} from "../wailsjs/go/main/App";

function App() {
    const [resultText, setResultText] = useState("動画ファイルを選択してGIFに変換してください 👇");
    const [selectedFile, setSelectedFile] = useState<File | null>(null);
    const [isConverting, setIsConverting] = useState(false);
    const [outputPath, setOutputPath] = useState<string>('');
    const [isDragOver, setIsDragOver] = useState(false);
    const [ffmpegInfo, setFFmpegInfo] = useState<string>('');
    const [outputDir, setOutputDir] = useState<string>('');

    const handleFileSelect = (file: File) => {
        if (file) {
            setSelectedFile(file);
            setResultText(`選択されたファイル: ${file.name} (${(file.size / 1024 / 1024).toFixed(2)} MB)`);
            setOutputPath('');
        }
    };

    const handleInputChange = (e: React.ChangeEvent<HTMLInputElement>) => {
        const file = e.target.files?.[0];
        if (file) {
            handleFileSelect(file);
        }
    };

    const handleDragOver = (e: React.DragEvent) => {
        e.preventDefault();
        setIsDragOver(true);
    };

    const handleDragLeave = (e: React.DragEvent) => {
        e.preventDefault();
        setIsDragOver(false);
    };

    const handleDrop = (e: React.DragEvent) => {
        e.preventDefault();
        setIsDragOver(false);
        
        const files = e.dataTransfer.files;
        if (files.length > 0) {
            const file = files[0];
            if (file.type.startsWith('video/')) {
                handleFileSelect(file);
            } else {
                setResultText("動画ファイルを選択してください");
            }
        }
    };

    const convertToGif = async () => {
        if (!selectedFile) {
            setResultText("ファイルを選択してください");
            return;
        }

        setIsConverting(true);
        setResultText("変換中です。しばらくお待ちください...");

        try {
            // ファイルをArrayBufferとして読み込み
            const arrayBuffer = await selectedFile.arrayBuffer();
            const uint8Array = new Uint8Array(arrayBuffer);
            
            const result = await ConvertUploadedVideoToGif(Array.from(uint8Array), selectedFile.name);
            setOutputPath(result);
            setResultText(`変換完了！\n出力ファイル: ${result}`);
        } catch (error) {
            console.error('変換エラー:', error);
            setResultText(`変換エラー: ${error}`);
        } finally {
            setIsConverting(false);
        }
    };

    const checkFFmpegInfo = async () => {
        try {
            const info = await GetFFmpegInfo();
            setFFmpegInfo(info);
        } catch (error) {
            setFFmpegInfo(`FFmpeg情報取得エラー: ${error}`);
        }
    };

    const getOutputDirInfo = async () => {
        try {
            const dir = await GetOutputDir();
            setOutputDir(dir);
        } catch (error) {
            setOutputDir(`出力フォルダ情報取得エラー: ${error}`);
        }
    };

    const openOutputFolder = async () => {
        try {
            await OpenOutputDir();
        } catch (error) {
            setResultText(`フォルダを開けませんでした: ${error}`);
        }
    };

    return (
        <div id="App">
            <div id="result" className="result" style={{whiteSpace: 'pre-line'}}>{resultText}</div>
            
            <div 
                className={`upload-area ${isDragOver ? 'drag-over' : ''} ${selectedFile ? 'has-file' : ''}`}
                onDragOver={handleDragOver}
                onDragLeave={handleDragLeave}
                onDrop={handleDrop}
            >
                <div className="upload-content">
                    <div className="upload-icon">📁</div>
                    <div className="upload-text">
                        {selectedFile ? (
                            <>
                                <p><strong>{selectedFile.name}</strong></p>
                                <p>{(selectedFile.size / 1024 / 1024).toFixed(2)} MB</p>
                            </>
                        ) : (
                            <>
                                <p>動画ファイルをドラッグ&ドロップ</p>
                                <p>または</p>
                            </>
                        )}
                    </div>
                    <input 
                        type="file" 
                        accept="video/*,.mp4,.mov,.avi,.mkv,.webm" 
                        onChange={handleInputChange}
                        className="file-input"
                        disabled={isConverting}
                        id="fileInput"
                    />
                    <label htmlFor="fileInput" className="file-label">
                        {selectedFile ? 'ファイルを変更' : 'ファイルを選択'}
                    </label>
                </div>
            </div>

            <div className="button-area">
                <button 
                    className="btn" 
                    onClick={convertToGif}
                    disabled={!selectedFile || isConverting}
                >
                    {isConverting ? '変換中...' : 'GIFに変換'}
                </button>
                
                <button 
                    className="btn" 
                    onClick={checkFFmpegInfo}
                    style={{marginLeft: '10px', backgroundColor: '#6c757d'}}
                >
                    FFmpeg情報を確認
                </button>
                
                <button 
                    className="btn" 
                    onClick={getOutputDirInfo}
                    style={{marginLeft: '10px', backgroundColor: '#17a2b8'}}
                >
                    出力フォルダを確認
                </button>
                
                <button 
                    className="btn" 
                    onClick={openOutputFolder}
                    style={{marginLeft: '10px', backgroundColor: '#28a745'}}
                >
                    📁 出力フォルダを開く
                </button>
            </div>

            {outputPath && (
                <div className="output-section">
                    <p>変換完了: {outputPath}</p>
                    <button 
                        className="btn" 
                        onClick={openOutputFolder}
                        style={{marginTop: '10px', backgroundColor: '#28a745'}}
                    >
                        📁 出力フォルダを開く
                    </button>
                </div>
            )}

            {outputDir && (
                <div className="output-info-section" style={{marginTop: '20px', padding: '10px', backgroundColor: '#e9ecef', borderRadius: '5px',color: '#333'}}>
                    <h4>出力フォルダ:</h4>
                    <p style={{fontSize: '14px', wordBreak: 'break-all'}}>{outputDir}</p>
                </div>
            )}

            {ffmpegInfo && (
                <div className="debug-section" style={{marginTop: '20px', padding: '10px', backgroundColor: '#f8f9fa', borderRadius: '5px', color: '#333'}}>
                    <h4>FFmpeg情報:</h4>
                    <pre style={{fontSize: '12px', whiteSpace: 'pre-wrap'}}>{ffmpegInfo}</pre>
                </div>
            )}
        </div>
    )
}

export default App

