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
  category: string;
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
  
  // Categories and filtering state
  const [categories, setCategories] = useState<string[]>([]);
  const [selectedCategory, setSelectedCategory] = useState<string>('all');
  const [searchQuery, setSearchQuery] = useState<string>('');
  const [filteredProducts, setFilteredProducts] = useState<Product[]>([]);


  useEffect(() => {
    fetchMessage();
    fetchProducts();
    fetchCategories();
  }, []);
  
  // Filter products when products, category, or search changes
  useEffect(() => {
    applyFilters();
  }, [products, selectedCategory, searchQuery]);

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
  
  const fetchCategories = async () => {
    try {
      const response = await fetch('/api/categories');
      if (!response.ok) {
        throw new Error(`HTTP error! status: ${response.status}`);
      }
      const data = await response.json();
      setCategories(data);
    } catch (err) {
      console.error('Failed to fetch categories:', err);
    }
  };
  
  const applyFilters = () => {
    let filtered = products;
    
    // Apply category filter
    if (selectedCategory !== 'all') {
      filtered = filtered.filter(product => product.category === selectedCategory);
    }
    
    // Apply search filter
    if (searchQuery.trim()) {
      const query = searchQuery.toLowerCase();
      filtered = filtered.filter(product => 
        product.name.toLowerCase().includes(query) ||
        product.category.toLowerCase().includes(query) ||
        product.description.toLowerCase().includes(query)
      );
    }
    
    setFilteredProducts(filtered);
  };
  
  const handleCategoryChange = (category: string) => {
    setSelectedCategory(category);
  };
  
  const handleSearchChange = (query: string) => {
    setSearchQuery(query);
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
          <Link to="/update-history" className="nav-link">
            📝 更新履歴
          </Link>
        </nav>
      </header>
      
      <main>
        <div className="main-content">
          {/* Sidebar */}
          <aside className="sidebar">
            <div className="search-section">
              <h3>商品検索</h3>
              <input
                type="text"
                placeholder="商品名やカテゴリで検索..."
                value={searchQuery}
                onChange={(e) => handleSearchChange(e.target.value)}
                className="search-input"
              />
            </div>
            
            <div className="categories-section">
              <h3>カテゴリ</h3>
              <div className="category-list">
                <button
                  className={`category-item ${selectedCategory === 'all' ? 'active' : ''}`}
                  onClick={() => handleCategoryChange('all')}
                >
                  すべて ({products.length})
                </button>
                {categories.map((category) => {
                  const count = products.filter(p => p.category === category).length;
                  return (
                    <button
                      key={category}
                      className={`category-item ${selectedCategory === category ? 'active' : ''}`}
                      onClick={() => handleCategoryChange(category)}
                    >
                      {category} ({count})
                    </button>
                  );
                })}
              </div>
            </div>
          </aside>
          
          {/* Products Section */}
          <section className="products-section">
            <div className="products-header">
              <h2>
                {selectedCategory === 'all' ? 'すべての商品' : selectedCategory}
                {searchQuery && ` - 「${searchQuery}」の検索結果`}
              </h2>
              <p className="products-count">{filteredProducts.length}件の商品</p>
            </div>
            
            {productsLoading && <p className="loading">商品を読み込み中...</p>}
            {productsError && <p className="error">エラー: {productsError}</p>}
            
            {!productsLoading && !productsError && (
              <>
                {filteredProducts.length === 0 ? (
                  <div className="no-products">
                    <p>条件に一致する商品が見つかりませんでした。</p>
                    <button onClick={() => { setSelectedCategory('all'); setSearchQuery(''); }}>すべての商品を表示</button>
                  </div>
                ) : (
                  <div className="products-grid">
                    {filteredProducts.map((product) => (
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
                          <span className="product-category">{product.category}</span>
                          <h3 className="product-name">{product.name}</h3>
                          <p className="product-price">¥{product.price.toLocaleString()}</p>
                          <p className="product-description">{product.description}</p>
                        </div>
                      </Link>
                    ))}
                  </div>
                )}
              </>
            )}
          </section>
        </div>

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
            <h4>サポート</h4>
            <ul>
              <li><Link to="/register">アカウント登録</Link></li>
              <li><Link to="/update-history">更新履歴</Link></li>
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