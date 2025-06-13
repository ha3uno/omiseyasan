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
  // State to track which prompts are expanded (by entry ID)
  const [expandedPrompts, setExpandedPrompts] = useState<Set<number>>(new Set());

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

  // Toggle expanded state for a specific entry
  const togglePromptExpansion = (entryId: number) => {
    setExpandedPrompts(prev => {
      const newSet = new Set(prev);
      if (newSet.has(entryId)) {
        newSet.delete(entryId);
      } else {
        newSet.add(entryId);
      }
      return newSet;
    });
  };

  // Check if prompt should be truncated (more than 3 lines or 150 characters)
  const shouldTruncatePrompt = (prompt: string) => {
    if (!prompt) return false;
    const lines = prompt.split('\n');
    return lines.length > 3 || prompt.length > 150;
  };

  // Get truncated version of the prompt
  const getTruncatedPrompt = (prompt: string) => {
    if (!prompt) return '';
    const lines = prompt.split('\n');
    if (lines.length > 3) {
      return lines.slice(0, 3).join('\n') + '...';
    }
    if (prompt.length > 150) {
      return prompt.substring(0, 150) + '...';
    }
    return prompt;
  };

  // Format prompt with collapsible functionality
  const formatPrompt = (entry: HistoryEntry) => {
    const { id, claudePrompt } = entry;
    if (!claudePrompt) return '';

    const isExpanded = expandedPrompts.has(id);
    const shouldTruncate = shouldTruncatePrompt(claudePrompt);
    const displayText = isExpanded || !shouldTruncate ? claudePrompt : getTruncatedPrompt(claudePrompt);

    return (
      <div className="prompt-container">
        <div className={`prompt-text ${isExpanded ? 'expanded' : 'collapsed'}`}>
          {displayText}
        </div>
        {shouldTruncate && (
          <button
            className="prompt-toggle-button"
            onClick={() => togglePromptExpansion(id)}
            type="button"
          >
            {isExpanded ? '閉じる' : 'もっと見る'}
          </button>
        )}
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
                          {formatPrompt(entry)}
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