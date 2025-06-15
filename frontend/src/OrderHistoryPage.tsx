import React, { useState, useEffect } from 'react';
import { Link } from 'react-router-dom';
import { useCart } from './CartContext';
import './App.css';

interface OrderItem {
  productId: number;
  name: string;
  price: number;
  quantity: number;
  subtotal: number;
}

interface ShippingInfo {
  name: string;
  address: string;
  phoneNumber: string;
}

interface Order {
  id: number;
  timestamp: string;
  items: OrderItem[];
  totalAmount: number;
  shippingInfo: ShippingInfo;
}

function OrderHistoryPage() {
  const { getTotalQuantity } = useCart();
  const [orders, setOrders] = useState<Order[]>([]);
  const [loading, setLoading] = useState<boolean>(true);
  const [error, setError] = useState<string>('');
  const [expandedOrders, setExpandedOrders] = useState<Set<number>>(new Set());

  useEffect(() => {
    fetchOrders();
  }, []);

  const fetchOrders = async () => {
    try {
      setLoading(true);
      setError('');
      const response = await fetch('/api/orders');
      
      if (!response.ok) {
        throw new Error(`HTTP error! status: ${response.status}`);
      }
      
      const data = await response.json();
      setOrders(data || []);
    } catch (err) {
      setError(err instanceof Error ? err.message : 'エラーが発生しました');
    } finally {
      setLoading(false);
    }
  };

  const toggleOrderExpansion = (orderId: number) => {
    const newExpanded = new Set(expandedOrders);
    if (newExpanded.has(orderId)) {
      newExpanded.delete(orderId);
    } else {
      newExpanded.add(orderId);
    }
    setExpandedOrders(newExpanded);
  };

  const formatDate = (timestamp: string) => {
    try {
      // Parse the timestamp (format: "2006-01-02 15:04:05")
      const date = new Date(timestamp.replace(' ', 'T'));
      return date.toLocaleDateString('ja-JP', {
        year: 'numeric',
        month: 'long',
        day: 'numeric',
        hour: '2-digit',
        minute: '2-digit'
      });
    } catch {
      return timestamp;
    }
  };

  return (
    <div className="App">
      <header className="App-header">
        <h1>おみせやさん♪</h1>
        <nav className="header-nav">
          <Link to="/" className="nav-link">🏠 ホーム</Link>
          <Link to="/cart" className="nav-link">🛒 カート ({getTotalQuantity()})</Link>
          <Link to="/update-history" className="nav-link">📝 更新履歴</Link>
        </nav>
      </header>
      
      <main>
        <div className="order-history-container">
          {/* ページヘッダーに新規登録ボタンを追加 */}
          <div style={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center', marginBottom: '30px' }}>
            <h2>注文履歴・更新履歴</h2>
            {/* 新規更新履歴登録ボタンを追加 */}
            <Link to="/update-history/register" className="action-btn primary">
              ➕ 新規更新履歴を登録
            </Link>
          </div>
          
          {loading && <p className="loading">注文履歴を読み込み中...</p>}
          {error && (
            <div className="error-container">
              <p className="error">エラー: {error}</p>
              <button onClick={fetchOrders} className="retry-btn">再試行</button>
            </div>
          )}
          
          {!loading && !error && orders.length === 0 && (
            <div className="empty-history">
              <h3>注文履歴がありません</h3>
              <p>まだ注文されていません。お気に入りの商品を見つけてお買い物をお楽しみください。</p>
              <div style={{ display: 'flex', gap: '15px', justifyContent: 'center', marginTop: '20px' }}>
                <Link to="/" className="nav-link">商品を見る</Link>
                {/* 空の履歴でも更新履歴登録へのリンクを表示 */}
                <Link to="/update-history/register" className="action-btn primary">
                  更新履歴を登録
                </Link>
              </div>
            </div>
          )}
          
          {!loading && !error && orders.length > 0 && (
            <div className="orders-list">
              {orders.map((order) => (
                <div key={order.id} className="order-card">
                  <div className="order-header" onClick={() => toggleOrderExpansion(order.id)}>
                    <div className="order-header-info">
                      <h3 className="order-id">注文番号: #{order.id}</h3>
                      <p className="order-date">{formatDate(order.timestamp)}</p>
                      <p className="order-total">合計: ¥{order.totalAmount.toLocaleString()}</p>
                    </div>
                    <div className="order-summary">
                      <p className="order-items-count">
                        {order.items.length}種類 / {order.items.reduce((total, item) => total + item.quantity, 0)}個
                      </p>
                      <button className="expand-btn">
                        {expandedOrders.has(order.id) ? '△ 詳細を閉じる' : '▽ 詳細を見る'}
                      </button>
                    </div>
                  </div>
                  
                  {expandedOrders.has(order.id) && (
                    <div className="order-details">
                      <div className="order-items-section">
                        <h4>注文内容</h4>
                        <div className="order-items-list">
                          {order.items.map((item, index) => (
                            <div key={index} className="order-item-row">
                              <div className="order-item-info">
                                <span className="order-item-name">{item.name}</span>
                                <span className="order-item-price">¥{item.price.toLocaleString()}</span>
                              </div>
                              <div className="order-item-quantity">×{item.quantity}</div>
                              <div className="order-item-subtotal">¥{item.subtotal.toLocaleString()}</div>
                            </div>
                          ))}
                        </div>
                      </div>
                      
                      <div className="shipping-info-section">
                        <h4>配送先情報</h4>
                        <div className="shipping-details">
                          <p><strong>お名前:</strong> {order.shippingInfo.name}</p>
                          <p><strong>住所:</strong> {order.shippingInfo.address}</p>
                          {order.shippingInfo.phoneNumber && (
                            <p><strong>電話番号:</strong> {order.shippingInfo.phoneNumber}</p>
                          )}
                        </div>
                      </div>
                      
                      <div className="order-actions">
                        <button 
                          onClick={() => window.print()} 
                          className="print-order-btn"
                        >
                          🖨️ この注文を印刷
                        </button>
                      </div>
                    </div>
                  )}
                </div>
              ))}
            </div>
          )}
          
          <div className="navigation-links">
            <Link to="/" className="nav-link">← ホームに戻る</Link>
            <Link to="/cart" className="nav-link">カートを見る →</Link>
            {/* 新規更新履歴登録への追加リンク */}
            <Link to="/update-history/register" className="nav-link">新規更新履歴を登録 →</Link>
          </div>
        </div>
      </main>
    </div>
  );
}

export default OrderHistoryPage;