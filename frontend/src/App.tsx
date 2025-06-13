import React, { useState, useEffect } from 'react';
import { Link } from 'react-router-dom';
import './App.css';

function App() {
  const [message, setMessage] = useState<string>('');
  const [loading, setLoading] = useState<boolean>(true);
  const [error, setError] = useState<string>('');

  // History form state
  const [historyForm, setHistoryForm] = useState({
    description: '',
    effortHours: 0,
    claudePrompt: ''
  });
  const [historySubmitting, setHistorySubmitting] = useState<boolean>(false);
  const [historySuccess, setHistorySuccess] = useState<boolean>(false);
  const [historyError, setHistoryError] = useState<string>('');

  useEffect(() => {
    fetchMessage();
  }, []);

  const fetchMessage = async () => {
    try {
      setLoading(true);
      setError('');
      const response = await fetch('/api/hello');
      if (!response.ok) {
        throw new Error(`HTTP error! status: ${response.status}`);
      }
      const data = await response.text();
      setMessage(data);
    } catch (err) {
      setError(err instanceof Error ? err.message : 'An error occurred');
    } finally {
      setLoading(false);
    }
  };

  const handleHistorySubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    
    if (!historyForm.description.trim()) {
      setHistoryError('変更内容は必須項目です。');
      return;
    }

    try {
      setHistorySubmitting(true);
      setHistoryError('');
      setHistorySuccess(false);

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

      // Clear form and show success
      setHistoryForm({
        description: '',
        effortHours: 0,
        claudePrompt: ''
      });
      setHistorySuccess(true);
      
      // Hide success message after 3 seconds
      setTimeout(() => setHistorySuccess(false), 3000);

    } catch (err) {
      setHistoryError(err instanceof Error ? err.message : 'エラーが発生しました');
    } finally {
      setHistorySubmitting(false);
    }
  };

  const handleHistoryInputChange = (e: React.ChangeEvent<HTMLInputElement | HTMLTextAreaElement>) => {
    const { name, value } = e.target;
    setHistoryForm(prev => ({
      ...prev,
      [name]: name === 'effortHours' ? parseFloat(value) || 0 : value
    }));
  };

  return (
    <div className="App">
      <header className="App-header">
        <h1>Full-stack Demo App</h1>
        <p>Go Backend + React Frontend</p>
        <Link to="/history" className="history-link">
          更新履歴を見る
        </Link>
      </header>
      
      <main>
        <div className="message-container">
          <h2>Message from Go Backend:</h2>
          {loading && <p className="loading">Loading...</p>}
          {error && <p className="error">Error: {error}</p>}
          {!loading && !error && <p>{message}</p>}
        </div>
        
        <button onClick={fetchMessage} disabled={loading}>
          Refresh Message
        </button>

        {/* History Form Section */}
        <div className="history-form-container">
          <h2>更新履歴の追加</h2>
          <form onSubmit={handleHistorySubmit} className="history-form">
            <div className="form-group">
              <label htmlFor="description">変更内容 *</label>
              <textarea
                id="description"
                name="description"
                value={historyForm.description}
                onChange={handleHistoryInputChange}
                placeholder="実装した機能や修正内容を入力してください"
                rows={3}
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
                onChange={handleHistoryInputChange}
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
                onChange={handleHistoryInputChange}
                placeholder="この作業で使用したClaudeへのプロンプトを入力してください&#10;- 複数行での入力可能&#10;- 箇条書きも対応"
                rows={6}
              />
            </div>

            {historyError && <p className="form-error">{historyError}</p>}
            {historySuccess && <p className="form-success">履歴が正常に追加されました！</p>}

            <button 
              type="submit" 
              disabled={historySubmitting}
              className="submit-button"
            >
              {historySubmitting ? '送信中...' : '履歴を追加'}
            </button>
          </form>
        </div>
      </main>
    </div>
  );
}

export default App;