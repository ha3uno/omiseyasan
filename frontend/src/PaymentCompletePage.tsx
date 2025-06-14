import React from 'react';
import { Link } from 'react-router-dom';
import './App.css';

function PaymentCompletePage() {
  return (
    <div className="App">
      <header className="App-header">
        <h1>おみせやさん♪</h1>
        <nav className="header-nav">
          <Link to="/" className="nav-link">🏠 ホーム</Link>
          <Link to="/cart" className="nav-link">🛒 カート</Link>
          <Link to="/register" className="nav-link">👤 アカウント登録</Link>
        </nav>
      </header>
      
      <main>
        <div className="payment-complete-container">
          <div className="payment-complete-content">
            <div className="success-icon">
              🎉
            </div>
            
            <h2>ご注文ありがとうございます！</h2>
            
            <div className="completion-message">
              <p className="main-message">
                ご注文が正常に完了しました。
              </p>
              
              <div className="order-details">
                <h3>📦 配送について</h3>
                <ul>
                  <li>商品は3〜5営業日以内に発送予定です</li>
                  <li>配送状況は別途メールでお知らせいたします</li>
                  <li>お急ぎの場合は、お電話でお問い合わせください</li>
                </ul>
              </div>
              
              <div className="payment-info">
                <h3>💳 決済について</h3>
                <p className="dummy-notice">
                  ※ これはデモアプリケーションです。実際の決済は行われていません。
                </p>
                <ul>
                  <li>決済処理が完了しました（ダミー）</li>
                  <li>領収書は登録いただいたメールアドレスに送付されます</li>
                  <li>ご不明点がございましたらお気軽にお問い合わせください</li>
                </ul>
              </div>
              
              <div className="contact-info">
                <h3>📞 お問い合わせ</h3>
                <p>カスタマーサポート: 0120-XXX-XXX</p>
                <p>営業時間: 平日 9:00〜18:00</p>
                <p>メール: support@omiseyasan.com</p>
              </div>
            </div>
            
            <div className="action-buttons">
              <Link to="/" className="continue-shopping-btn">
                🛍️ お買い物を続ける
              </Link>
              
              <button 
                onClick={() => window.print()} 
                className="print-btn"
              >
                🖨️ この画面を印刷
              </button>
            </div>
            
            <div className="thank-you-message">
              <p>
                またのご利用を心よりお待ちしております。<br />
                今後ともおみせやさん♪をよろしくお願いいたします。
              </p>
            </div>
          </div>
        </div>
      </main>
    </div>
  );
}

export default PaymentCompletePage;