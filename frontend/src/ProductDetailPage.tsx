import React, { useState, useEffect } from 'react';
import { useParams, Link } from 'react-router-dom';
import { useCart } from './CartContext';
import './App.css';

interface Product {
  id: number;
  name: string;
  price: number;
  description: string;
  imageUrl: string;
  category: string;
}

function ProductDetailPage() {
  const { id } = useParams<{ id: string }>();
  const { addToCart, getTotalQuantity } = useCart();
  const [product, setProduct] = useState<Product | null>(null);
  const [loading, setLoading] = useState<boolean>(true);
  const [error, setError] = useState<string>('');

  useEffect(() => {
    if (id) {
      fetchProduct(id);
    }
  }, [id]);

  const fetchProduct = async (productId: string) => {
    try {
      setLoading(true);
      setError('');
      const response = await fetch(`/api/products/${productId}`);
      
      if (!response.ok) {
        if (response.status === 404) {
          throw new Error('商品が見つかりませんでした。');
        }
        throw new Error(`HTTP error! status: ${response.status}`);
      }
      
      const data = await response.json();
      setProduct(data);
    } catch (err) {
      setError(err instanceof Error ? err.message : 'エラーが発生しました');
    } finally {
      setLoading(false);
    }
  };

  const handleAddToCart = () => {
    if (product) {
      addToCart({
        id: product.id,
        name: product.name,
        price: product.price,
        imageUrl: product.imageUrl,
      });
      alert(`${product.name} をカートに追加しました！`);
    }
  };

  if (loading) {
    return (
      <div className="App">
        <header className="App-header">
          <h1>おみせやさん♪</h1>
          <nav className="header-nav">
            <Link to="/" className="nav-link">🏠 ホーム</Link>
            <Link to="/cart" className="nav-link">🛒 カート ({getTotalQuantity()})</Link>
            <Link to="/register" className="nav-link">👤 アカウント登録</Link>
          </nav>
        </header>
        <main>
          <p className="loading">商品を読み込み中...</p>
        </main>
      </div>
    );
  }

  if (error) {
    return (
      <div className="App">
        <header className="App-header">
          <h1>おみせやさん♪</h1>
          <nav className="header-nav">
            <Link to="/" className="nav-link">🏠 ホーム</Link>
            <Link to="/cart" className="nav-link">🛒 カート ({getTotalQuantity()})</Link>
            <Link to="/register" className="nav-link">👤 アカウント登録</Link>
          </nav>
        </header>
        <main>
          <div className="error-container">
            <p className="error">エラー: {error}</p>
            <Link to="/" className="nav-link">商品一覧に戻る</Link>
          </div>
        </main>
      </div>
    );
  }

  if (!product) {
    return (
      <div className="App">
        <header className="App-header">
          <h1>おみせやさん♪</h1>
          <nav className="header-nav">
            <Link to="/" className="nav-link">🏠 ホーム</Link>
            <Link to="/cart" className="nav-link">🛒 カート ({getTotalQuantity()})</Link>
            <Link to="/register" className="nav-link">👤 アカウント登録</Link>
          </nav>
        </header>
        <main>
          <p>商品が見つかりませんでした。</p>
          <Link to="/" className="nav-link">商品一覧に戻る</Link>
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
          <Link to="/cart" className="nav-link">🛒 カートを見る</Link>
          <Link to="/register" className="nav-link">👤 アカウント登録</Link>
        </nav>
      </header>
      
      <main>
        <div className="product-detail-container">
          <div className="product-detail">
            <div className="product-image-container">
              <img
                src={product.imageUrl}
                alt={product.name}
                className="product-detail-image"
              />
            </div>
            
            <div className="product-detail-info">
              <span className="product-detail-category">{product.category}</span>
              <h2 className="product-detail-name">{product.name}</h2>
              <p className="product-detail-price">¥{product.price.toLocaleString()}</p>
              <p className="product-detail-description">{product.description}</p>
              
              <div className="product-actions">
                <button 
                  onClick={handleAddToCart}
                  className="add-to-cart-button"
                >
                  🛒 カートに入れる
                </button>
              </div>
            </div>
          </div>
          
          <div className="navigation-links">
            <Link to="/" className="nav-link">← 商品一覧に戻る</Link>
          </div>
        </div>
      </main>
    </div>
  );
}

export default ProductDetailPage;