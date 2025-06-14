import React, { useEffect } from 'react';
import { Link, useNavigate } from 'react-router-dom';
import { useCart } from './CartContext';
import './App.css';

interface User {
  id: number;
  name: string;
  address: string;
  phoneNumber: string;
  email: string;
}

interface MyAccountPageProps {
  loggedInUser: User | null;
  onLogout: () => void;
}

function MyAccountPage({ loggedInUser, onLogout }: MyAccountPageProps) {
  const navigate = useNavigate();
  const { getTotalQuantity } = useCart();

  // Redirect to login if not authenticated
  useEffect(() => {
    if (!loggedInUser) {
      navigate('/login');
    }
  }, [loggedInUser, navigate]);

  const handleLogout = () => {
    if (window.confirm('ログアウトしますか？')) {
      onLogout();
      navigate('/');
    }
  };

  // Show loading or redirect message if user is not logged in
  if (!loggedInUser) {
    return (
      <div className="App">
        <header className="App-header">
          <h1>おみせやさん♪</h1>
          <nav className="header-nav">
            <Link to="/" className="nav-link">🏠 ホーム</Link>
            <Link to="/cart" className="nav-link">🛒 カート ({getTotalQuantity()})</Link>
            <Link to="/login" className="nav-link">🔑 ログイン</Link>
            <Link to="/register" className="nav-link">👤 アカウント登録</Link>
          </nav>
        </header>
        
        <main>
          <div className="account-container">
            <div className="account-redirect">
              <h2>ログインが必要です</h2>
              <p>マイアカウントページにアクセスするにはログインしてください。</p>
              <Link to="/login" className="nav-link">ログインページへ</Link>
            </div>
          </div>
        </main>
      </div>
    );
  }

  return (
    <div className="App">
      <header className="App-header">
        <h1>おみせやさん♪</h1>
        <nav className="header-nav">
          <Link to="/" className="nav-link">🏠 ホーム</Link>
          <Link to="/cart" className="nav-link">🛒 カート ({getTotalQuantity()})</Link>
          <Link to="/order-history" className="nav-link">📋 注文履歴</Link>
          <Link to="/history" className="nav-link">📝 更新履歴</Link>
        </nav>
      </header>
      
      <main>
        <div className="account-container">
          <h2>マイアカウント</h2>
          <p className="welcome-message">こんにちは、{loggedInUser.name}さん！</p>
          
          <div className="account-content">
            <div className="account-info-section">
              <h3>アカウント情報</h3>
              <div className="account-info">
                <div className="info-row">
                  <span className="info-label">ユーザーID:</span>
                  <span className="info-value">#{loggedInUser.id}</span>
                </div>
                
                <div className="info-row">
                  <span className="info-label">お名前:</span>
                  <span className="info-value">{loggedInUser.name}</span>
                </div>
                
                <div className="info-row">
                  <span className="info-label">メールアドレス:</span>
                  <span className="info-value">{loggedInUser.email}</span>
                </div>
                
                <div className="info-row">
                  <span className="info-label">住所:</span>
                  <span className="info-value">
                    {loggedInUser.address || '未設定'}
                  </span>
                </div>
                
                <div className="info-row">
                  <span className="info-label">電話番号:</span>
                  <span className="info-value">
                    {loggedInUser.phoneNumber || '未設定'}
                  </span>
                </div>
              </div>
            </div>
            
            <div className="account-actions-section">
              <h3>アカウント操作</h3>
              <div className="account-actions">
                <Link to="/order-history" className="action-btn primary">
                  📋 注文履歴を見る
                </Link>
                
                <Link to="/cart" className="action-btn secondary">
                  🛒 カートを見る
                </Link>
                
                <Link to="/" className="action-btn secondary">
                  🛍️ お買い物を続ける
                </Link>
                
                <button onClick={handleLogout} className="action-btn logout">
                  🚪 ログアウト
                </button>
              </div>
            </div>
            
            <div className="account-features-section">
              <h3>ご利用ガイド</h3>
              <div className="features-list">
                <div className="feature-item">
                  <h4>🛒 簡単お買い物</h4>
                  <p>商品をカートに追加して、簡単に注文できます。</p>
                </div>
                
                <div className="feature-item">
                  <h4>📋 注文履歴</h4>
                  <p>過去の注文を確認して、再注文も可能です。</p>
                </div>
                
                <div className="feature-item">
                  <h4>🔒 安全・安心</h4>
                  <p>お客様の情報は安全に保護されています。</p>
                </div>
              </div>
            </div>
          </div>
          
          <div className="navigation-links">
            <Link to="/" className="nav-link">← ホームに戻る</Link>
          </div>
        </div>
      </main>
    </div>
  );
}

export default MyAccountPage;