import React, { useState } from 'react';
import { Link, useNavigate } from 'react-router-dom';
import { useCart } from './CartContext';
import './App.css';

interface UpdateHistoryForm {
  description: string;
  effortHours: number;
  claudePrompt: string;
}

function UpdateHistoryRegisterPage() {
  const navigate = useNavigate();
  const { getTotalQuantity } = useCart();
  
  // 更新履歴フォームの状態管理
  const [historyForm, setHistoryForm] = useState<UpdateHistoryForm>({
    description: '',
    effortHours: 0,
    claudePrompt: ''
  });
  const [historySubmitting, setHistorySubmitting] = useState<boolean>(false);
  const [historySuccess, setHistorySuccess] = useState<boolean>(false);
  const [historyError, setHistoryError] = useState<string>('');

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

      // バックエンドの /api/history エンドポイントに送信
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

      // 登録成功後、フォームクリアして成功メッセージ表示
      setHistoryForm({
        description: '',
        effortHours: 0,
        claudePrompt: ''
      });
      setHistorySuccess(true);
      
      // 2秒後に /update-history へリダイレクト
      setTimeout(() => {
        navigate('/update-history');
      }, 2000);

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
        <h1>おみせやさん♪</h1>
        <p>更新履歴の登録</p>
        <nav className="header-nav">
          <Link to="/" className="nav-link">🏠 ホーム</Link>
          <Link to="/cart" className="nav-link">🛒 カート ({getTotalQuantity()})</Link>
          <Link to="/update-history" className="nav-link">📋 更新履歴</Link>
          <Link to="/order-history" className="nav-link">🛍️ 注文履歴</Link>
        </nav>
      </header>
      
      <main>
        <div className="history-form-container">
          <h2>更新履歴の登録</h2>
          <p>開発の進捗や変更内容を記録してください。</p>
          
          {/* 成功メッセージ */}
          {historySuccess && (
            <div className="form-success">
              ✅ 更新履歴が正常に登録されました！更新履歴ページにリダイレクトします...
            </div>
          )}
          
          <form onSubmit={handleHistorySubmit} className="history-form">
            {/* 変更内容入力フィールド */}
            <div className="form-group">
              <label htmlFor="description">変更内容 *</label>
              <textarea
                id="description"
                name="description"
                value={historyForm.description}
                onChange={handleHistoryInputChange}
                placeholder="実装した機能や修正内容を入力してください"
                rows={4}
                required
              />
            </div>

            {/* 工数入力フィールド */}
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

            {/* Claudeプロンプト入力フィールド */}
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

            {/* エラーメッセージ */}
            {historyError && <p className="form-error">{historyError}</p>}

            {/* 送信ボタンとナビゲーション */}
            <div style={{ display: 'flex', gap: '15px', marginTop: '20px' }}>
              <button 
                type="submit" 
                disabled={historySubmitting}
                className="submit-button"
              >
                {historySubmitting ? '登録中...' : '📝 履歴を登録'}
              </button>
              
              <Link to="/update-history" className="action-btn secondary">
                ← 更新履歴一覧に戻る
              </Link>
            </div>
          </form>
        </div>
      </main>
    </div>
  );
}

export default UpdateHistoryRegisterPage;