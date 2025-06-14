import React, { useState, useEffect } from 'react';
import { Link } from 'react-router-dom';
import { useCart } from './CartContext';
import './App.css';

interface Product {
  id: number;
  name: string;
  price: number;
  description: string;
  imageUrl: string;
}

interface User {
  id: number;
  name: string;
  address: string;
  phoneNumber: string;
  email: string;
}

const USER_STORAGE_KEY = 'omiseyasan-user';

interface AppProps {
  loggedInUser: User | null;
  onLogin: (user: User) => void;
  onLogout: () => void;
}

function App({ loggedInUser, onLogin, onLogout }: AppProps) {
  const { getTotalQuantity } = useCart();
  const [message, setMessage] = useState<string>('');
  const [loading, setLoading] = useState<boolean>(true);
  const [error, setError] = useState<string>('');
  
  // Products state
  const [products, setProducts] = useState<Product[]>([]);
  const [productsLoading, setProductsLoading] = useState<boolean>(true);
  const [productsError, setProductsError] = useState<string>('');

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
    fetchProducts();
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

  const fetchProducts = async () => {
    try {
      setProductsLoading(true);
      setProductsError('');
      const response = await fetch('/api/products');
      if (!response.ok) {
        throw new Error(`HTTP error! status: ${response.status}`);
      }
      const data = await response.json();
      setProducts(data);
    } catch (err) {
      setProductsError(err instanceof Error ? err.message : 'An error occurred');
    } finally {
      setProductsLoading(false);
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
        <h1>おみせやさん♪</h1>
        <p>かわいい商品がいっぱいのオンラインショップ</p>
        <nav className="header-nav">
          <Link to="/cart" className="nav-link">
            🛒 カート ({getTotalQuantity()})
          </Link>
          <Link to="/order-history" className="nav-link">
            📋 注文履歴
          </Link>
          {loggedInUser ? (
            <>
              <Link to="/my-account" className="nav-link">
                👤 {loggedInUser.name}さん
              </Link>
              <button onClick={onLogout} className="nav-link logout-btn">
                🚪 ログアウト
              </button>
            </>
          ) : (
            <>
              <Link to="/login" className="nav-link">
                🔑 ログイン
              </Link>
              <Link to="/register" className="nav-link">
                👤 アカウント登録
              </Link>
            </>
          )}
          <Link to="/history" className="nav-link">
            📝 更新履歴
          </Link>
        </nav>
      </header>
      
      <main>
        {/* Products Section */}
        <section className="products-section">
          <h2>おすすめ商品</h2>
          {productsLoading && <p className="loading">商品を読み込み中...</p>}
          {productsError && <p className="error">エラー: {productsError}</p>}
          
          {!productsLoading && !productsError && (
            <div className="products-grid">
              {products.map((product) => (
                <Link
                  key={product.id}
                  to={`/products/${product.id}`}
                  className="product-card"
                >
                  <img
                    src={product.imageUrl}
                    alt={product.name}
                    className="product-image"
                  />
                  <div className="product-info">
                    <h3 className="product-name">{product.name}</h3>
                    <p className="product-price">¥{product.price.toLocaleString()}</p>
                    <p className="product-description">{product.description}</p>
                  </div>
                </Link>
              ))}
            </div>
          )}
        </section>

        {/* Debug Section - Keep the original hello message for debugging */}
        <section className="debug-section">
          <details>
            <summary>デバッグ情報</summary>
            <div className="message-container">
              <h3>Backend 接続テスト:</h3>
              {loading && <p className="loading">Loading...</p>}
              {error && <p className="error">Error: {error}</p>}
              {!loading && !error && <p>{message}</p>}
              <button onClick={fetchMessage} disabled={loading}>
                接続テスト
              </button>
            </div>
          </details>
        </section>

        {/* History Form Section */}
        <section className="admin-section">
          <details>
            <summary>管理者機能</summary>
            <div className="history-form-container">
              <h3>更新履歴の追加</h3>
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
          </details>
        </section>
      </main>
      
      <footer className="App-footer">
        <div className="footer-content">
          <div className="footer-section">
            <h4>おみせやさん♪</h4>
            <p>かわいい商品がいっぱいのオンラインショップ</p>
          </div>
          
          <div className="footer-section">
            <h4>サービス</h4>
            <ul>
              <li><Link to="/">商品一覧</Link></li>
              <li><Link to="/cart">ショッピングカート</Link></li>
              <li><Link to="/order-history">注文履歴</Link></li>
            </ul>
          </div>
          
          <div className="footer-section">
            <h4>アカウント</h4>
            <ul>
              <li><Link to="/register">アカウント登録</Link></li>
              <li><Link to="/history">更新履歴</Link></li>
            </ul>
          </div>
          
          <div className="footer-section">
            <h4>お問い合わせ</h4>
            <p>📞 0120-XXX-XXX</p>
            <p>📧 support@omiseyasan.com</p>
            <p>営業時間: 平日 9:00〜18:00</p>
          </div>
        </div>
        
        <div className="footer-bottom">
          <p>&copy; 2024 おみせやさん♪ All rights reserved.</p>
          <p className="demo-notice">※ これはデモアプリケーションです</p>
        </div>
      </footer>
    </div>
  );
}

export default App;