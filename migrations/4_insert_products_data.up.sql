-- Insert products data (only if not exists)
-- This file contains all product data with categories

-- 食品・食材
INSERT INTO products (id, name, price, description, image_url, category) 
SELECT 1, 'フレッシュアボカド', 280, '栄養豊富でクリーミーな食感のアボカド', 'https://cdn.pixabay.com/photo/2016/03/05/19/02/avocado-1238250_1280.jpg', '食品・食材'
WHERE NOT EXISTS (SELECT 1 FROM products WHERE id = 1);

INSERT INTO products (id, name, price, description, image_url, category) 
SELECT 2, 'オーガニックバナナ', 198, '自然栽培された甘いバナナ', 'https://cdn.pixabay.com/photo/2017/06/28/10/14/bananas-2450718_1280.jpg', '食品・食材'
WHERE NOT EXISTS (SELECT 1 FROM products WHERE id = 2);

INSERT INTO products (id, name, price, description, image_url, category) 
SELECT 3, '新鮮トマト', 320, 'ジューシーで甘みのある完熟トマト', 'https://cdn.pixabay.com/photo/2016/08/10/16/11/tomatoes-1583288_1280.jpg', '食品・食材'
WHERE NOT EXISTS (SELECT 1 FROM products WHERE id = 3);

INSERT INTO products (id, name, price, description, image_url, category) 
SELECT 4, '有機人参', 250, '甘くて栄養価の高い有機人参', 'https://cdn.pixabay.com/photo/2016/08/09/10/30/carrots-1579717_1280.jpg', '食品・食材'
WHERE NOT EXISTS (SELECT 1 FROM products WHERE id = 4);

INSERT INTO products (id, name, price, description, image_url, category) 
SELECT 5, 'レモン', 180, '爽やかな香りの国産レモン', 'https://cdn.pixabay.com/photo/2016/09/15/19/24/salad-1672505_1280.jpg', '食品・食材'
WHERE NOT EXISTS (SELECT 1 FROM products WHERE id = 5);

-- パン・ベーカリー
INSERT INTO products (id, name, price, description, image_url, category) 
SELECT 6, 'クロワッサン', 280, 'バターたっぷりのサクサククロワッサン', 'https://cdn.pixabay.com/photo/2014/04/22/02/56/croissant-329136_1280.jpg', 'パン・ベーカリー'
WHERE NOT EXISTS (SELECT 1 FROM products WHERE id = 6);

INSERT INTO products (id, name, price, description, image_url, category) 
SELECT 7, '食パン', 320, 'しっとりふわふわの食パン', 'https://cdn.pixabay.com/photo/2017/06/23/23/49/bread-2434370_1280.jpg', 'パン・ベーカリー'
WHERE NOT EXISTS (SELECT 1 FROM products WHERE id = 7);

INSERT INTO products (id, name, price, description, image_url, category) 
SELECT 8, 'ベーグル', 250, 'もちもち食感のプレーンベーグル', 'https://cdn.pixabay.com/photo/2017/05/07/08/56/blank-2292428_1280.jpg', 'パン・ベーカリー'
WHERE NOT EXISTS (SELECT 1 FROM products WHERE id = 8);

-- 飲み物
INSERT INTO products (id, name, price, description, image_url, category) 
SELECT 9, 'コーヒー豆', 980, '香り豊かなアラビカコーヒー豆', 'https://cdn.pixabay.com/photo/2015/07/02/20/57/coffee-beans-829734_1280.jpg', '飲み物'
WHERE NOT EXISTS (SELECT 1 FROM products WHERE id = 9);

INSERT INTO products (id, name, price, description, image_url, category) 
SELECT 10, '緑茶', 680, '上質な日本茶の緑茶', 'https://cdn.pixabay.com/photo/2015/10/12/15/46/tea-984367_1280.jpg', '飲み物'
WHERE NOT EXISTS (SELECT 1 FROM products WHERE id = 10);

INSERT INTO products (id, name, price, description, image_url, category) 
SELECT 11, 'フルーツジュース', 380, '100%フレッシュオレンジジュース', 'https://cdn.pixabay.com/photo/2017/05/12/08/29/juice-2306330_1280.jpg', '飲み物'
WHERE NOT EXISTS (SELECT 1 FROM products WHERE id = 11);

-- お菓子・スイーツ
INSERT INTO products (id, name, price, description, image_url, category) 
SELECT 12, 'チョコレートケーキ', 1280, '濃厚なチョコレートケーキ', 'https://cdn.pixabay.com/photo/2017/01/11/11/33/cake-1971552_1280.jpg', 'お菓子・スイーツ'
WHERE NOT EXISTS (SELECT 1 FROM products WHERE id = 12);

INSERT INTO products (id, name, price, description, image_url, category) 
SELECT 13, 'マカロン', 480, 'カラフルなフレンチマカロン', 'https://cdn.pixabay.com/photo/2017/07/05/15/41/macarons-2474188_1280.jpg', 'お菓子・スイーツ'
WHERE NOT EXISTS (SELECT 1 FROM products WHERE id = 13);

INSERT INTO products (id, name, price, description, image_url, category) 
SELECT 14, 'ドーナツ', 220, 'ふわふわの手作りドーナツ', 'https://cdn.pixabay.com/photo/2017/11/22/00/18/donuts-2969490_1280.jpg', 'お菓子・スイーツ'
WHERE NOT EXISTS (SELECT 1 FROM products WHERE id = 14);

INSERT INTO products (id, name, price, description, image_url, category) 
SELECT 15, 'アイスクリーム', 380, 'なめらかなバニラアイスクリーム', 'https://cdn.pixabay.com/photo/2017/06/03/23/30/ice-cream-2371033_1280.jpg', 'お菓子・スイーツ'
WHERE NOT EXISTS (SELECT 1 FROM products WHERE id = 15);

-- 雑貨・インテリア
INSERT INTO products (id, name, price, description, image_url, category) 
SELECT 16, 'コーヒーカップ', 1200, 'エレガントな白いコーヒーカップ', 'https://cdn.pixabay.com/photo/2016/11/29/12/51/cafe-1869656_1280.jpg', '雑貨・インテリア'
WHERE NOT EXISTS (SELECT 1 FROM products WHERE id = 16);

INSERT INTO products (id, name, price, description, image_url, category) 
SELECT 17, '観葉植物', 2800, '室内を彩る美しい観葉植物', 'https://cdn.pixabay.com/photo/2016/11/21/16/21/aloe-1846759_1280.jpg', '雑貨・インテリア'
WHERE NOT EXISTS (SELECT 1 FROM products WHERE id = 17);

INSERT INTO products (id, name, price, description, image_url, category) 
SELECT 18, 'キャンドル', 680, 'リラックス効果のあるアロマキャンドル', 'https://cdn.pixabay.com/photo/2017/02/15/12/12/candles-2069021_1280.jpg', '雑貨・インテリア'
WHERE NOT EXISTS (SELECT 1 FROM products WHERE id = 18);

INSERT INTO products (id, name, price, description, image_url, category) 
SELECT 19, 'ノート', 580, '上質な紙を使った手帳ノート', 'https://cdn.pixabay.com/photo/2016/03/26/22/13/books-1281581_1280.jpg', '雑貨・インテリア'
WHERE NOT EXISTS (SELECT 1 FROM products WHERE id = 19);

INSERT INTO products (id, name, price, description, image_url, category) 
SELECT 20, '時計', 4800, 'シンプルでスタイリッシュな壁掛け時計', 'https://cdn.pixabay.com/photo/2017/06/08/17/32/clock-2383642_1280.jpg', '雑貨・インテリア'
WHERE NOT EXISTS (SELECT 1 FROM products WHERE id = 20);

-- 文房具
INSERT INTO products (id, name, price, description, image_url, category) 
SELECT 21, '色鉛筆セット', 1280, '48色の豊富なカラー色鉛筆', 'https://cdn.pixabay.com/photo/2015/11/07/11/22/pencils-1030092_1280.jpg', '文房具'
WHERE NOT EXISTS (SELECT 1 FROM products WHERE id = 21);

INSERT INTO products (id, name, price, description, image_url, category) 
SELECT 22, '万年筆', 3800, '書き心地なめらかな万年筆', 'https://cdn.pixabay.com/photo/2016/03/26/22/22/fountain-pen-1281576_1280.jpg', '文房具'
WHERE NOT EXISTS (SELECT 1 FROM products WHERE id = 22);

INSERT INTO products (id, name, price, description, image_url, category) 
SELECT 23, 'スケッチブック', 880, 'アーティスト向け高品質スケッチブック', 'https://cdn.pixabay.com/photo/2017/07/31/11/21/people-2557396_1280.jpg', '文房具'
WHERE NOT EXISTS (SELECT 1 FROM products WHERE id = 23);

-- キッチン用品
INSERT INTO products (id, name, price, description, image_url, category) 
SELECT 24, 'エプロン', 1980, 'かわいいデザインのキッチンエプロン', 'https://cdn.pixabay.com/photo/2017/06/01/18/46/cooking-2363358_1280.jpg', 'キッチン用品'
WHERE NOT EXISTS (SELECT 1 FROM products WHERE id = 24);

INSERT INTO products (id, name, price, description, image_url, category) 
SELECT 25, 'まな板', 2480, '抗菌仕様の木製まな板', 'https://cdn.pixabay.com/photo/2016/11/29/04/19/cooking-1867761_1280.jpg', 'キッチン用品'
WHERE NOT EXISTS (SELECT 1 FROM products WHERE id = 25);

INSERT INTO products (id, name, price, description, image_url, category) 
SELECT 26, '包丁', 5800, 'プロ仕様の切れ味抜群包丁', 'https://cdn.pixabay.com/photo/2016/12/26/17/28/spaghetti-1932466_1280.jpg', 'キッチン用品'
WHERE NOT EXISTS (SELECT 1 FROM products WHERE id = 26);

-- 衣類・ファッション
INSERT INTO products (id, name, price, description, image_url, category) 
SELECT 27, 'Tシャツ', 1980, 'コットン100%の快適Tシャツ', 'https://cdn.pixabay.com/photo/2016/12/06/09/31/blank-1886008_1280.jpg', '衣類・ファッション'
WHERE NOT EXISTS (SELECT 1 FROM products WHERE id = 27);

INSERT INTO products (id, name, price, description, image_url, category) 
SELECT 28, 'スニーカー', 7800, '歩きやすいカジュアルスニーカー', 'https://cdn.pixabay.com/photo/2016/06/03/17/35/shoes-1433925_1280.jpg', '衣類・ファッション'
WHERE NOT EXISTS (SELECT 1 FROM products WHERE id = 28);

INSERT INTO products (id, name, price, description, image_url, category) 
SELECT 29, '帽子', 2800, '紫外線対策にぴったりの帽子', 'https://cdn.pixabay.com/photo/2016/04/09/17/17/fashion-1319032_1280.jpg', '衣類・ファッション'
WHERE NOT EXISTS (SELECT 1 FROM products WHERE id = 29);

INSERT INTO products (id, name, price, description, image_url, category) 
SELECT 30, 'バッグ', 4800, '機能的でおしゃれなトートバッグ', 'https://cdn.pixabay.com/photo/2017/08/01/11/48/woman-2563491_1280.jpg', '衣類・ファッション'
WHERE NOT EXISTS (SELECT 1 FROM products WHERE id = 30);

-- Update sequence for auto-increment (set to a safe high value)
SELECT setval('products_id_seq', 100, false);