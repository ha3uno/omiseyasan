import React, { useState } from 'react';
import { Link, useNavigate } from 'react-router-dom';
import { useCart } from './CartContext';
import './App.css';

interface ShippingInfo {
  name: string;
  address: string;
  phoneNumber: string;
}

interface OrderItem {
  productId: number;
  name: string;
  price: number;
  quantity: number;
  subtotal: number;
}

function OrderConfirmationPage() {
  const navigate = useNavigate();
  const { cartItems, getTotalPrice, clearCart } = useCart();
  
  const [shippingInfo, setShippingInfo] = useState<ShippingInfo>({
    name: '',
    address: '',
    phoneNumber: ''
  });
  
  const [submitting, setSubmitting] = useState<boolean>(false);
  const [error, setError] = useState<string>('');

  // If cart is empty, redirect to home
  if (cartItems.length === 0) {
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
          <div className="order-container">
            <div className="empty-order">
              <h2>カートが空です</h2>
              <p>注文を進めるにはカートに商品を追加してください。</p>
              <Link to="/" className="nav-link">商品を見る</Link>
            </div>
          </div>
        </main>
      </div>
    );
  }

  const handleInputChange = (e: React.ChangeEvent<HTMLInputElement | HTMLTextAreaElement>) => {
    const { name, value } = e.target;
    setShippingInfo(prev => ({
      ...prev,
      [name]: value
    }));
  };

  const handleSubmitOrder = async (e: React.FormEvent) => {
    e.preventDefault();
    
    // Basic validation
    if (!shippingInfo.name.trim() || !shippingInfo.address.trim()) {
      setError('お名前と住所は必須項目です。');
      return;
    }

    try {
      setSubmitting(true);
      setError('');

      // Convert cart items to order items format
      const orderItems: OrderItem[] = cartItems.map(item => ({
        productId: item.id,
        name: item.name,
        price: item.price,
        quantity: item.quantity,
        subtotal: item.price * item.quantity
      }));

      const orderData = {
        items: orderItems,
        totalAmount: getTotalPrice(),
        shippingInfo: shippingInfo
      };

      const response = await fetch('/api/orders', {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
        },
        body: JSON.stringify(orderData),
      });

      if (!response.ok) {
        const errorText = await response.text();
        throw new Error(errorText || `HTTP error! status: ${response.status}`);
      }

      const createdOrder = await response.json();
      console.log('Order created:', createdOrder);
      
      // Clear cart after successful order
      clearCart();
      
      // Navigate to payment complete page
      navigate('/payment-complete');

    } catch (err) {
      setError(err instanceof Error ? err.message : '注文処理中にエラーが発生しました');
    } finally {
      setSubmitting(false);
    }
  };

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
        <div className="order-container">
          <h2>注文確認</h2>
          
          {/* Order Summary */}
          <div className="order-summary">
            <h3>ご注文内容</h3>
            <div className="order-items-list">
              {cartItems.map((item) => (
                <div key={item.id} className="order-item">
                  <img 
                    src={item.imageUrl} 
                    alt={item.name} 
                    className="order-item-image"
                  />
                  <div className="order-item-details">
                    <h4 className="order-item-name">{item.name}</h4>
                    <p className="order-item-price">¥{item.price.toLocaleString()} × {item.quantity}</p>
                    <p className="order-item-subtotal">小計: ¥{(item.price * item.quantity).toLocaleString()}</p>
                  </div>
                </div>
              ))}
            </div>
            
            <div className="order-total">
              <h3>合計金額: ¥{getTotalPrice().toLocaleString()}</h3>
            </div>
          </div>

          {/* Shipping Information Form */}
          <div className="shipping-form-section">
            <h3>配送先情報</h3>
            <form onSubmit={handleSubmitOrder} className="shipping-form">
              <div className="form-group">
                <label htmlFor="name">お名前 *</label>
                <input
                  type="text"
                  id="name"
                  name="name"
                  value={shippingInfo.name}
                  onChange={handleInputChange}
                  placeholder="山田太郎"
                  required
                />
              </div>

              <div className="form-group">
                <label htmlFor="address">住所 *</label>
                <textarea
                  id="address"
                  name="address"
                  value={shippingInfo.address}
                  onChange={handleInputChange}
                  placeholder="〒123-4567 東京都渋谷区..."
                  rows={3}
                  required
                />
              </div>

              <div className="form-group">
                <label htmlFor="phoneNumber">電話番号</label>
                <input
                  type="tel"
                  id="phoneNumber"
                  name="phoneNumber"
                  value={shippingInfo.phoneNumber}
                  onChange={handleInputChange}
                  placeholder="090-1234-5678"
                />
              </div>

              {error && <p className="form-error">{error}</p>}

              <div className="order-actions">
                <Link to="/cart" className="back-to-cart-btn">
                  ← カートに戻る
                </Link>
                
                <button 
                  type="submit" 
                  disabled={submitting}
                  className="confirm-order-btn"
                >
                  {submitting ? '処理中...' : '🚀 注文を確定する（ダミー決済）'}
                </button>
              </div>
            </form>
          </div>
        </div>
      </main>
    </div>
  );
}

export default OrderConfirmationPage;