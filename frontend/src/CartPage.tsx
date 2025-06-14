import React from 'react';
import { Link } from 'react-router-dom';
import { useCart } from './CartContext';
import './App.css';

function CartPage() {
  const { 
    cartItems, 
    removeFromCart, 
    updateQuantity, 
    clearCart, 
    getTotalQuantity, 
    getTotalPrice 
  } = useCart();

  const handleQuantityChange = (productId: number, newQuantity: number) => {
    if (newQuantity < 1) {
      removeFromCart(productId);
    } else {
      updateQuantity(productId, newQuantity);
    }
  };

  const handleClearCart = () => {
    if (window.confirm('カート内の全ての商品を削除しますか？')) {
      clearCart();
    }
  };

  const handleProceedToCheckout = () => {
    alert('購入手続き機能は準備中です！');
  };

  return (
    <div className="App">
      <header className="App-header">
        <h1>おみせやさん♪</h1>
        <nav className="header-nav">
          <Link to="/" className="nav-link">🏠 ホーム</Link>
          <Link to="/register" className="nav-link">👤 アカウント登録</Link>
          <Link to="/history" className="nav-link">📝 更新履歴</Link>
        </nav>
      </header>
      
      <main>
        <div className="cart-container">
          <h2>ショッピングカート</h2>
          
          {cartItems.length === 0 ? (
            <div className="empty-cart">
              <p>カートは空です</p>
              <Link to="/" className="nav-link">商品を見る</Link>
            </div>
          ) : (
            <>
              <div className="cart-summary">
                <p>商品数: {getTotalQuantity()}点</p>
                <p className="total-price">合計: ¥{getTotalPrice().toLocaleString()}</p>
              </div>

              <div className="cart-items">
                <div className="cart-table">
                  <div className="cart-table-header">
                    <div className="cart-header-item">商品</div>
                    <div className="cart-header-price">単価</div>
                    <div className="cart-header-quantity">数量</div>
                    <div className="cart-header-subtotal">小計</div>
                    <div className="cart-header-actions">操作</div>
                  </div>
                  
                  {cartItems.map((item) => (
                    <div key={item.id} className="cart-table-row">
                      <div className="cart-item-info">
                        <img 
                          src={item.imageUrl} 
                          alt={item.name} 
                          className="cart-item-image"
                        />
                        <div className="cart-item-details">
                          <h3 className="cart-item-name">{item.name}</h3>
                          <Link 
                            to={`/products/${item.id}`} 
                            className="cart-item-link"
                          >
                            商品詳細を見る
                          </Link>
                        </div>
                      </div>
                      
                      <div className="cart-item-price">
                        ¥{item.price.toLocaleString()}
                      </div>
                      
                      <div className="cart-item-quantity">
                        <div className="quantity-controls">
                          <button 
                            onClick={() => handleQuantityChange(item.id, item.quantity - 1)}
                            className="quantity-btn"
                            disabled={item.quantity <= 1}
                          >
                            -
                          </button>
                          <span className="quantity-display">{item.quantity}</span>
                          <button 
                            onClick={() => handleQuantityChange(item.id, item.quantity + 1)}
                            className="quantity-btn"
                          >
                            +
                          </button>
                        </div>
                      </div>
                      
                      <div className="cart-item-subtotal">
                        ¥{(item.price * item.quantity).toLocaleString()}
                      </div>
                      
                      <div className="cart-item-actions">
                        <button 
                          onClick={() => removeFromCart(item.id)}
                          className="remove-btn"
                        >
                          削除
                        </button>
                      </div>
                    </div>
                  ))}
                </div>
              </div>

              <div className="cart-actions">
                <button 
                  onClick={handleClearCart}
                  className="clear-cart-btn"
                >
                  カートを空にする
                </button>
                
                <button 
                  onClick={handleProceedToCheckout}
                  className="checkout-btn"
                >
                  🚀 購入手続きへ進む
                </button>
              </div>
            </>
          )}
          
          <div className="navigation-links">
            <Link to="/" className="nav-link">← 買い物を続ける</Link>
          </div>
        </div>
      </main>
    </div>
  );
}

export default CartPage;