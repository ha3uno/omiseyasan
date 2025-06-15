import React, { useState, useEffect } from 'react';
import { Link } from 'react-router-dom';
import { useCart } from './CartContext';
import './HistoryPage.css';
import './App.css';

interface HistoryEntry {
  id: number;
  timestamp: string;
  description: string;
  effortHours: number;
  claudePrompt: string;
}

interface UpdateHistoryForm {
  description: string;
  effortHours: number;
  claudePrompt: string;
}

const HistoryPage: React.FC = () => {
  const { getTotalQuantity } = useCart();
  const [historyData, setHistoryData] = useState<HistoryEntry[]>([]);
  const [loading, setLoading] = useState<boolean>(true);
  const [error, setError] = useState<string>('');
  // State to track which prompts are expanded (by entry ID)
  const [expandedPrompts, setExpandedPrompts] = useState<Set<number>>(new Set());
  
  // 新規履歴登録用の状態
  const [showRegistrationForm, setShowRegistrationForm] = useState<boolean>(false);
  const [historyForm, setHistoryForm] = useState<UpdateHistoryForm>({
    description: '',
    effortHours: 0,
    claudePrompt: ''
  });
  const [submitting, setSubmitting] = useState<boolean>(false);
  const [submitSuccess, setSubmitSuccess] = useState<boolean>(false);
  const [submitError, setSubmitError] = useState<string>('');

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

  // 新規履歴登録フォーム送信処理
  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    
    if (!historyForm.description.trim()) {
      setSubmitError('変更内容は必須項目です。');
      return;
    }

    try {
      setSubmitting(true);
      setSubmitError('');
      setSubmitSuccess(false);

      const response = await fetch('/api/history', {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
        },
        body: JSON.stringify(historyForm),
      });

      if (!response.ok) {
        throw new Error(`HTTP error! status: ${response.status}`);
      }

      // 登録成功後、フォームクリア
      setHistoryForm({
        description: '',
        effortHours: 0,
        claudePrompt: ''
      });
      setSubmitSuccess(true);
      setShowRegistrationForm(false);
      
      // 履歴を再取得
      fetchHistory();
      
      // 成功メッセージを3秒後に非表示
      setTimeout(() => {
        setSubmitSuccess(false);
      }, 3000);

    } catch (err) {
      setSubmitError(err instanceof Error ? err.message : 'エラーが発生しました');
    } finally {
      setSubmitting(false);
    }
  };

  const handleInputChange = (e: React.ChangeEvent<HTMLInputElement | HTMLTextAreaElement>) => {
    const { name, value } = e.target;
    setHistoryForm(prev => ({
      ...prev,
      [name]: name === 'effortHours' ? parseFloat(value) || 0 : value
    }));
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
    <div className="App">
      <header className="App-header">
        <h1>おみせやさん♪</h1>
        <p>更新履歴の管理</p>
        <nav className="header-nav">
          <Link to="/" className="nav-link">🏠 ホーム</Link>
          <Link to="/cart" className="nav-link">🛒 カート ({getTotalQuantity()})</Link>
        </nav>
      </header>

      <main>
        <div className="history-container">
          <div style={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center', marginBottom: '30px' }}>
            <h2>更新履歴</h2>
            <button 
              onClick={() => setShowRegistrationForm(!showRegistrationForm)}
              className="action-btn primary"
            >
              {showRegistrationForm ? '📝 登録フォームを閉じる' : '➕ 新規履歴を登録'}
            </button>
          </div>
          
          {/* 成功メッセージ */}
          {submitSuccess && (
            <div className="form-success">
              ✅ 更新履歴が正常に登録されました！
            </div>
          )}
          
          {/* 新規登録フォーム */}
          {showRegistrationForm && (
            <div className="registration-form-container">
              <h3>新規更新履歴の登録</h3>
              <form onSubmit={handleSubmit} className="history-form">
                <div className="form-group">
                  <label htmlFor="description">変更内容 *</label>
                  <textarea
                    id="description"
                    name="description"
                    value={historyForm.description}
                    onChange={handleInputChange}
                    placeholder="実装した機能や修正内容を入力してください"
                    rows={4}
                    required
                  />
                </div>

                <div className="form-group">
                  <label htmlFor="effortHours">工数（時間）</label>
                  <input
                    type="number"
                    id="effortHours"
                    name="effortHours"
                    value={historyForm.effortHours}
                    onChange={handleInputChange}
                    min="0"
                    step="0.5"
                    placeholder="0.5"
                  />
                </div>

                <div className="form-group">
                  <label htmlFor="claudePrompt">Claudeへのプロンプト</label>
                  <textarea
                    id="claudePrompt"
                    name="claudePrompt"
                    value={historyForm.claudePrompt}
                    onChange={handleInputChange}
                    placeholder="この作業で使用したClaudeへのプロンプトを入力してください&#10;- 複数行での入力可能&#10;- 箇条書きも対応"
                    rows={6}
                  />
                </div>

                {submitError && <p className="form-error">{submitError}</p>}

                <div style={{ display: 'flex', gap: '15px', marginTop: '20px' }}>
                  <button 
                    type="submit" 
                    disabled={submitting}
                    className="submit-button"
                  >
                    {submitting ? '登録中...' : '📝 履歴を登録'}
                  </button>
                  
                  <button 
                    type="button"
                    onClick={() => setShowRegistrationForm(false)}
                    className="cancel-button"
                  >
                    キャンセル
                  </button>
                </div>
              </form>
            </div>
          )}
          
          {/* 履歴一覧セクション */}
          <div className="history-list-section">
            <h3>登録済み履歴一覧</h3>
            
            {loading && <p className="loading">データベースから読み込み中...</p>}
            {error && (
              <div className="error-container">
                <p className="error">エラー: {error}</p>
                <button onClick={fetchHistory} className="retry-btn">再試行</button>
              </div>
            )}
            
            {!loading && !error && (
              <>
                <div className="history-summary">
                  <p>総件数: {historyData.length}件</p>
                  <button onClick={fetchHistory} className="refresh-button">
                    🔄 更新
                  </button>
                </div>
                
                {historyData.length === 0 ? (
                  <div className="no-data-container">
                    <p className="no-data">履歴データがありません。</p>
                    <p className="no-data-hint">上の「新規履歴を登録」ボタンから履歴を追加してください。</p>
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
          </div>
          
          <div className="navigation-links">
            <Link to="/" className="nav-link">← ホームに戻る</Link>
          </div>
        </div>
      </main>
    </div>
  );
};

export default HistoryPage;