import React, { useState, useEffect } from 'react';
import { Link } from 'react-router-dom';
import './HistoryPage.css';

interface HistoryEntry {
  id: number;
  timestamp: string;
  description: string;
  effortHours: number;
  claudePrompt: string;
}

const HistoryPage: React.FC = () => {
  const [historyData, setHistoryData] = useState<HistoryEntry[]>([]);
  const [loading, setLoading] = useState<boolean>(true);
  const [error, setError] = useState<string>('');

  useEffect(() => {
    fetchHistory();
  }, []);

  const fetchHistory = async () => {
    try {
      setLoading(true);
      setError('');
      const response = await fetch('/api/history');
      if (!response.ok) {
        throw new Error(`HTTP error! status: ${response.status}`);
      }
      const data = await response.json();
      setHistoryData(data);
    } catch (err) {
      setError(err instanceof Error ? err.message : 'An error occurred');
    } finally {
      setLoading(false);
    }
  };

  const formatPrompt = (prompt: string) => {
    if (!prompt) return '';
    
    // Split by lines and format as list if it contains bullet points or numbered items
    const lines = prompt.split('\n').filter(line => line.trim());
    
    // Check if it looks like a list (starts with -, *, numbers, etc.)
    const isListLike = lines.some(line => 
      /^\s*[-*•]\s/.test(line) || /^\s*\d+\.\s/.test(line)
    );

    if (isListLike) {
      return (
        <ul className="prompt-list">
          {lines.map((line, index) => {
            const cleanLine = line.replace(/^\s*[-*•]\s*/, '').replace(/^\s*\d+\.\s*/, '');
            return <li key={index}>{cleanLine}</li>;
          })}
        </ul>
      );
    }

    // Otherwise, just preserve line breaks
    return (
      <div className="prompt-text">
        {lines.map((line, index) => (
          <React.Fragment key={index}>
            {line}
            {index < lines.length - 1 && <br />}
          </React.Fragment>
        ))}
      </div>
    );
  };

  return (
    <div className="history-page">
      <header className="history-header">
        <h1>更新履歴</h1>
        <p className="database-note">📊 データベースから履歴を取得中</p>
        <Link to="/" className="back-link">← ホームに戻る</Link>
      </header>

      <main className="history-content">
        {loading && <p className="loading">データベースから読み込み中...</p>}
        {error && (
          <div className="error-container">
            <p className="error">エラー: {error}</p>
            <p className="error-hint">データベース接続に問題がある可能性があります。</p>
          </div>
        )}
        
        {!loading && !error && (
          <>
            <div className="history-summary">
              <p>総件数: {historyData.length}件</p>
              <button onClick={fetchHistory} className="refresh-button">
                更新
              </button>
            </div>
            
            {historyData.length === 0 ? (
              <div className="no-data-container">
                <p className="no-data">履歴データがありません。</p>
                <p className="no-data-hint">ホームページから履歴を追加してください。</p>
              </div>
            ) : (
              <div className="table-container">
                <table className="history-table">
                  <thead>
                    <tr>
                      <th>時刻</th>
                      <th>変更内容</th>
                      <th>工数（時間）</th>
                      <th>Claudeへのプロンプト</th>
                    </tr>
                  </thead>
                  <tbody>
                    {historyData.map((entry) => (
                      <tr key={entry.id}>
                        <td className="timestamp-cell">{entry.timestamp}</td>
                        <td className="description-cell">{entry.description}</td>
                        <td className="effort-cell">{entry.effortHours}</td>
                        <td className="prompt-cell">
                          {formatPrompt(entry.claudePrompt)}
                        </td>
                      </tr>
                    ))}
                  </tbody>
                </table>
              </div>
            )}
          </>
        )}
      </main>
    </div>
  );
};

export default HistoryPage;