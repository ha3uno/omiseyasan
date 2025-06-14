package data

// Product represents a product in the e-commerce site
type Product struct {
	ID          int     `json:"id"`
	Name        string  `json:"name"`
	Price       float64 `json:"price"`
	Description string  `json:"description"`
	ImageURL    string  `json:"imageUrl"`
}

// GetProducts returns a comprehensive list of products with high-quality images
func GetProducts() []Product {
	return []Product{
		// 食品・食材
		{ID: 1, Name: "フレッシュアボカド", Price: 280, Description: "栄養豊富でクリーミーな食感のアボカド", ImageURL: "https://cdn.pixabay.com/photo/2016/03/05/19/02/avocado-1238250_1280.jpg"},
		{ID: 2, Name: "オーガニックバナナ", Price: 198, Description: "自然栽培された甘いバナナ", ImageURL: "https://cdn.pixabay.com/photo/2017/06/28/10/14/bananas-2450718_1280.jpg"},
		{ID: 3, Name: "新鮮トマト", Price: 320, Description: "ジューシーで甘みのある完熟トマト", ImageURL: "https://cdn.pixabay.com/photo/2016/08/10/16/11/tomatoes-1583288_1280.jpg"},
		{ID: 4, Name: "有機人参", Price: 250, Description: "甘くて栄養価の高い有機人参", ImageURL: "https://cdn.pixabay.com/photo/2016/08/09/10/30/carrots-1579717_1280.jpg"},
		{ID: 5, Name: "レモン", Price: 180, Description: "爽やかな香りの国産レモン", ImageURL: "https://cdn.pixabay.com/photo/2016/09/15/19/24/salad-1672505_1280.jpg"},
		
		// パン・ベーカリー
		{ID: 6, Name: "クロワッサン", Price: 280, Description: "バターたっぷりのサクサククロワッサン", ImageURL: "https://cdn.pixabay.com/photo/2014/04/22/02/56/croissant-329136_1280.jpg"},
		{ID: 7, Name: "食パン", Price: 320, Description: "しっとりふわふわの食パン", ImageURL: "https://cdn.pixabay.com/photo/2017/06/23/23/49/bread-2434370_1280.jpg"},
		{ID: 8, Name: "ベーグル", Price: 250, Description: "もちもち食感のプレーンベーグル", ImageURL: "https://cdn.pixabay.com/photo/2017/05/07/08/56/blank-2292428_1280.jpg"},
		
		// 飲み物
		{ID: 9, Name: "コーヒー豆", Price: 980, Description: "香り豊かなアラビカコーヒー豆", ImageURL: "https://cdn.pixabay.com/photo/2015/07/02/20/57/coffee-beans-829734_1280.jpg"},
		{ID: 10, Name: "緑茶", Price: 680, Description: "上質な日本茶の緑茶", ImageURL: "https://cdn.pixabay.com/photo/2015/10/12/15/46/tea-984367_1280.jpg"},
		{ID: 11, Name: "フルーツジュース", Price: 380, Description: "100%フレッシュオレンジジュース", ImageURL: "https://cdn.pixabay.com/photo/2017/05/12/08/29/juice-2306330_1280.jpg"},
		
		// お菓子・スイーツ
		{ID: 12, Name: "チョコレートケーキ", Price: 1280, Description: "濃厚なチョコレートケーキ", ImageURL: "https://cdn.pixabay.com/photo/2017/01/11/11/33/cake-1971552_1280.jpg"},
		{ID: 13, Name: "マカロン", Price: 480, Description: "カラフルなフレンチマカロン", ImageURL: "https://cdn.pixabay.com/photo/2017/07/05/15/41/macarons-2474188_1280.jpg"},
		{ID: 14, Name: "ドーナツ", Price: 220, Description: "ふわふわの手作りドーナツ", ImageURL: "https://cdn.pixabay.com/photo/2017/11/22/00/18/donuts-2969490_1280.jpg"},
		{ID: 15, Name: "アイスクリーム", Price: 380, Description: "なめらかなバニラアイスクリーム", ImageURL: "https://cdn.pixabay.com/photo/2017/06/03/23/30/ice-cream-2371033_1280.jpg"},
		
		// 雑貨・インテリア
		{ID: 16, Name: "コーヒーカップ", Price: 1200, Description: "エレガントな白いコーヒーカップ", ImageURL: "https://cdn.pixabay.com/photo/2016/11/29/12/51/cafe-1869656_1280.jpg"},
		{ID: 17, Name: "観葉植物", Price: 2800, Description: "室内を彩る美しい観葉植物", ImageURL: "https://cdn.pixabay.com/photo/2016/11/21/16/21/aloe-1846759_1280.jpg"},
		{ID: 18, Name: "キャンドル", Price: 680, Description: "リラックス効果のあるアロマキャンドル", ImageURL: "https://cdn.pixabay.com/photo/2017/02/15/12/12/candles-2069021_1280.jpg"},
		{ID: 19, Name: "ノート", Price: 580, Description: "上質な紙を使った手帳ノート", ImageURL: "https://cdn.pixabay.com/photo/2016/03/26/22/13/books-1281581_1280.jpg"},
		{ID: 20, Name: "時計", Price: 4800, Description: "シンプルでスタイリッシュな壁掛け時計", ImageURL: "https://cdn.pixabay.com/photo/2017/06/08/17/32/clock-2383642_1280.jpg"},
		
		// 文房具
		{ID: 21, Name: "色鉛筆セット", Price: 1280, Description: "48色の豊富なカラー色鉛筆", ImageURL: "https://cdn.pixabay.com/photo/2015/11/07/11/22/pencils-1030092_1280.jpg"},
		{ID: 22, Name: "万年筆", Price: 3800, Description: "書き心地なめらかな万年筆", ImageURL: "https://cdn.pixabay.com/photo/2016/03/26/22/22/fountain-pen-1281576_1280.jpg"},
		{ID: 23, Name: "スケッチブック", Price: 880, Description: "アーティスト向け高品質スケッチブック", ImageURL: "https://cdn.pixabay.com/photo/2017/07/31/11/21/people-2557396_1280.jpg"},
		
		// キッチン用品
		{ID: 24, Name: "エプロン", Price: 1980, Description: "かわいいデザインのキッチンエプロン", ImageURL: "https://cdn.pixabay.com/photo/2017/06/01/18/46/cooking-2363358_1280.jpg"},
		{ID: 25, Name: "まな板", Price: 2480, Description: "抗菌仕様の木製まな板", ImageURL: "https://cdn.pixabay.com/photo/2016/11/29/04/19/cooking-1867761_1280.jpg"},
		{ID: 26, Name: "包丁", Price: 5800, Description: "プロ仕様の切れ味抜群包丁", ImageURL: "https://cdn.pixabay.com/photo/2016/12/26/17/28/spaghetti-1932466_1280.jpg"},
		
		// 衣類・ファッション
		{ID: 27, Name: "Tシャツ", Price: 1980, Description: "コットン100%の快適Tシャツ", ImageURL: "https://cdn.pixabay.com/photo/2016/12/06/09/31/blank-1886008_1280.jpg"},
		{ID: 28, Name: "スニーカー", Price: 7800, Description: "歩きやすいカジュアルスニーカー", ImageURL: "https://cdn.pixabay.com/photo/2016/06/03/17/35/shoes-1433925_1280.jpg"},
		{ID: 29, Name: "帽子", Price: 2800, Description: "紫外線対策にぴったりの帽子", ImageURL: "https://cdn.pixabay.com/photo/2016/04/09/17/17/fashion-1319032_1280.jpg"},
		{ID: 30, Name: "バッグ", Price: 4800, Description: "機能的でおしゃれなトートバッグ", ImageURL: "https://cdn.pixabay.com/photo/2017/08/01/11/48/woman-2563491_1280.jpg"},
		
		// 本・書籍
		{ID: 31, Name: "料理本", Price: 1680, Description: "初心者向けの簡単料理レシピ本", ImageURL: "https://cdn.pixabay.com/photo/2016/11/29/01/34/man-1867175_1280.jpg"},
		{ID: 32, Name: "小説", Price: 1280, Description: "ベストセラー小説", ImageURL: "https://cdn.pixabay.com/photo/2015/11/19/21/10/glasses-1052010_1280.jpg"},
		{ID: 33, Name: "写真集", Price: 2980, Description: "美しい風景写真集", ImageURL: "https://cdn.pixabay.com/photo/2015/09/05/21/51/reading-925589_1280.jpg"},
		
		// 花・植物
		{ID: 34, Name: "バラの花束", Price: 3800, Description: "上品な赤いバラの花束", ImageURL: "https://cdn.pixabay.com/photo/2017/02/15/13/40/tulips-2068692_1280.jpg"},
		{ID: 35, Name: "多肉植物", Price: 980, Description: "手入れ簡単な可愛い多肉植物", ImageURL: "https://cdn.pixabay.com/photo/2016/11/29/05/07/plants-1867564_1280.jpg"},
		{ID: 36, Name: "ひまわり", Price: 580, Description: "元気いっぱいのひまわり", ImageURL: "https://cdn.pixabay.com/photo/2015/06/19/20/13/sunflower-815666_1280.jpg"},
		
		// 日用品
		{ID: 37, Name: "タオル", Price: 1280, Description: "ふわふわの高品質タオル", ImageURL: "https://cdn.pixabay.com/photo/2017/03/25/23/32/bathroom-2174743_1280.jpg"},
		{ID: 38, Name: "石鹸", Price: 480, Description: "天然成分のハンドメイド石鹸", ImageURL: "https://cdn.pixabay.com/photo/2016/04/22/16/15/soap-1345226_1280.jpg"},
		{ID: 39, Name: "歯ブラシ", Price: 380, Description: "やわらかい毛先の歯ブラシ", ImageURL: "https://cdn.pixabay.com/photo/2017/03/30/15/47/toothbrush-2188368_1280.jpg"},
		{ID: 40, Name: "シャンプー", Price: 1580, Description: "髪に優しいオーガニックシャンプー", ImageURL: "https://cdn.pixabay.com/photo/2016/04/22/04/57/spa-1344997_1280.jpg"},
		
		// おもちゃ・ホビー
		{ID: 41, Name: "積み木", Price: 2480, Description: "知育に最適な木製積み木", ImageURL: "https://cdn.pixabay.com/photo/2017/08/25/22/31/cubes-2682124_1280.jpg"},
		{ID: 42, Name: "パズル", Price: 1880, Description: "1000ピースの美しい風景パズル", ImageURL: "https://cdn.pixabay.com/photo/2017/10/24/20/33/puzzle-2883634_1280.jpg"},
		{ID: 43, Name: "ぬいぐるみ", Price: 2280, Description: "ふわふわで可愛いテディベア", ImageURL: "https://cdn.pixabay.com/photo/2017/03/27/14/46/teddy-2179048_1280.jpg"},
		
		// 電子機器・アクセサリー
		{ID: 44, Name: "イヤホン", Price: 4800, Description: "高音質ワイヤレスイヤホン", ImageURL: "https://cdn.pixabay.com/photo/2017/05/26/19/36/air-pods-2344840_1280.jpg"},
		{ID: 45, Name: "スマホケース", Price: 1980, Description: "衝撃に強いスマートフォンケース", ImageURL: "https://cdn.pixabay.com/photo/2017/01/13/01/22/mobile-phone-1976082_1280.jpg"},
		{ID: 46, Name: "充電器", Price: 2880, Description: "急速充電対応モバイルバッテリー", ImageURL: "https://cdn.pixabay.com/photo/2016/11/29/07/15/smartphone-1868423_1280.jpg"},
		
		// スポーツ・アウトドア
		{ID: 47, Name: "ヨガマット", Price: 3800, Description: "滑りにくい高品質ヨガマット", ImageURL: "https://cdn.pixabay.com/photo/2017/08/07/14/02/people-2604149_1280.jpg"},
		{ID: 48, Name: "水筒", Price: 2480, Description: "保温保冷機能付きステンレス水筒", ImageURL: "https://cdn.pixabay.com/photo/2017/06/02/18/24/fruit-2367029_1280.jpg"},
		{ID: 49, Name: "ランニングシューズ", Price: 8800, Description: "軽量で快適なランニングシューズ", ImageURL: "https://cdn.pixabay.com/photo/2016/11/19/18/06/feet-1840619_1280.jpg"},
		{ID: 50, Name: "テント", Price: 15800, Description: "2人用コンパクトテント", ImageURL: "https://cdn.pixabay.com/photo/2016/11/21/17/44/adventure-1846482_1280.jpg"},
		
		// 追加の食品
		{ID: 51, Name: "パスタ", Price: 480, Description: "本格イタリアン スパゲッティ", ImageURL: "https://cdn.pixabay.com/photo/2017/02/15/15/17/pasta-2069547_1280.jpg"},
		{ID: 52, Name: "オリーブオイル", Price: 1280, Description: "エクストラバージンオリーブオイル", ImageURL: "https://cdn.pixabay.com/photo/2016/04/15/08/04/olive-oil-1330418_1280.jpg"},
		{ID: 53, Name: "ハチミツ", Price: 980, Description: "天然100%の純粋ハチミツ", ImageURL: "https://cdn.pixabay.com/photo/2015/04/08/13/13/honey-713153_1280.jpg"},
		{ID: 54, Name: "チーズ", Price: 1580, Description: "濃厚な味わいのナチュラルチーズ", ImageURL: "https://cdn.pixabay.com/photo/2016/12/10/21/26/food-1898194_1280.jpg"},
		{ID: 55, Name: "ワイン", Price: 2980, Description: "芳醇な香りの赤ワイン", ImageURL: "https://cdn.pixabay.com/photo/2016/10/22/20/34/wines-1761613_1280.jpg"},
		
		// 美容・コスメ
		{ID: 56, Name: "化粧水", Price: 2480, Description: "保湿効果の高い化粧水", ImageURL: "https://cdn.pixabay.com/photo/2017/04/03/15/52/lotion-2199447_1280.jpg"},
		{ID: 57, Name: "リップクリーム", Price: 680, Description: "潤いを与えるリップクリーム", ImageURL: "https://cdn.pixabay.com/photo/2017/08/07/05/84/lipstick-2603584_1280.jpg"},
		{ID: 58, Name: "香水", Price: 4800, Description: "上品で優雅な香りの香水", ImageURL: "https://cdn.pixabay.com/photo/2016/03/27/07/32/perfume-1282291_1280.jpg"},
		
		// インテリア・家具
		{ID: 59, Name: "クッション", Price: 1980, Description: "ふわふわで快適なクッション", ImageURL: "https://cdn.pixabay.com/photo/2016/11/18/17/20/couch-1835923_1280.jpg"},
		{ID: 60, Name: "額縁", Price: 1280, Description: "写真や絵を美しく飾る額縁", ImageURL: "https://cdn.pixabay.com/photo/2017/03/27/14/50/picture-frame-2179221_1280.jpg"},
		{ID: 61, Name: "ランプ", Price: 3800, Description: "温かい光のテーブルランプ", ImageURL: "https://cdn.pixabay.com/photo/2016/11/21/12/51/light-1845539_1280.jpg"},
		{ID: 62, Name: "カーテン", Price: 4800, Description: "遮光性の高いエレガントなカーテン", ImageURL: "https://cdn.pixabay.com/photo/2016/11/18/15/44/room-1835482_1280.jpg"},
		
		// 楽器・音楽
		{ID: 63, Name: "ギター", Price: 28000, Description: "初心者にも優しいアコースティックギター", ImageURL: "https://cdn.pixabay.com/photo/2016/11/23/15/48/audience-1853662_1280.jpg"},
		{ID: 64, Name: "ピアノ楽譜", Price: 1280, Description: "クラシック名曲集の楽譜", ImageURL: "https://cdn.pixabay.com/photo/2015/05/07/11/02/guitar-756326_1280.jpg"},
		
		// 工具・DIY
		{ID: 65, Name: "ドライバーセット", Price: 1980, Description: "多機能精密ドライバーセット", ImageURL: "https://cdn.pixabay.com/photo/2017/09/23/21/21/tools-2779591_1280.jpg"},
		{ID: 66, Name: "ハンマー", Price: 2480, Description: "使いやすい軽量ハンマー", ImageURL: "https://cdn.pixabay.com/photo/2016/04/04/14/12/hammer-1307834_1280.jpg"},
		
		// ペット用品
		{ID: 67, Name: "ペットフード", Price: 1580, Description: "栄養バランスの良いプレミアムフード", ImageURL: "https://cdn.pixabay.com/photo/2016/02/19/15/46/dog-1210559_1280.jpg"},
		{ID: 68, Name: "猫のおもちゃ", Price: 680, Description: "猫ちゃんが喜ぶ可愛いおもちゃ", ImageURL: "https://cdn.pixabay.com/photo/2017/02/20/18/03/cat-2083492_1280.jpg"},
		
		// 季節商品
		{ID: 69, Name: "扇子", Price: 1280, Description: "涼しい風を運ぶ美しい扇子", ImageURL: "https://cdn.pixabay.com/photo/2016/06/02/13/58/fan-1432666_1280.jpg"},
		{ID: 70, Name: "マフラー", Price: 2800, Description: "暖かくて肌触りの良いマフラー", ImageURL: "https://cdn.pixabay.com/photo/2016/11/29/01/15/scarf-1867063_1280.jpg"},
		
		// 健康・ウェルネス
		{ID: 71, Name: "サプリメント", Price: 3800, Description: "ビタミン豊富な健康サプリメント", ImageURL: "https://cdn.pixabay.com/photo/2017/02/12/10/29/pills-2059687_1280.jpg"},
		{ID: 72, Name: "マッサージオイル", Price: 1980, Description: "リラックス効果のあるマッサージオイル", ImageURL: "https://cdn.pixabay.com/photo/2017/08/07/05/84/essential-oils-2603467_1280.jpg"},
		
		// ガーデニング
		{ID: 73, Name: "種子セット", Price: 880, Description: "家庭菜園用の野菜種子セット", ImageURL: "https://cdn.pixabay.com/photo/2016/06/02/13/20/seeds-1432416_1280.jpg"},
		{ID: 74, Name: "じょうろ", Price: 1580, Description: "ガーデニング用のおしゃれなじょうろ", ImageURL: "https://cdn.pixabay.com/photo/2016/03/26/22/22/watering-can-1281581_1280.jpg"},
		
		// アート・クラフト
		{ID: 75, Name: "絵の具セット", Price: 2480, Description: "アクリル絵の具24色セット", ImageURL: "https://cdn.pixabay.com/photo/2017/05/15/17/43/paint-2314569_1280.jpg"},
		{ID: 76, Name: "キャンバス", Price: 1280, Description: "油絵・アクリル画用キャンバス", ImageURL: "https://cdn.pixabay.com/photo/2016/03/26/22/14/art-1281582_1280.jpg"},
		
		// 追加の雑貨
		{ID: 77, Name: "貯金箱", Price: 1680, Description: "可愛いブタさんの貯金箱", ImageURL: "https://cdn.pixabay.com/photo/2017/01/14/10/56/piggy-bank-1979430_1280.jpg"},
		{ID: 78, Name: "カレンダー", Price: 980, Description: "美しい風景写真のカレンダー", ImageURL: "https://cdn.pixabay.com/photo/2017/01/31/13/14/agenda-2023741_1280.jpg"},
		{ID: 79, Name: "定規セット", Price: 580, Description: "学習・仕事に便利な定規セット", ImageURL: "https://cdn.pixabay.com/photo/2016/03/26/22/17/ruler-1281588_1280.jpg"},
		{ID: 80, Name: "マグネット", Price: 380, Description: "冷蔵庫用可愛いマグネット", ImageURL: "https://cdn.pixabay.com/photo/2017/01/13/01/22/magnets-1976077_1280.jpg"},
		
		// 食器・キッチン雑貨
		{ID: 81, Name: "お皿セット", Price: 3800, Description: "上品な白いお皿6枚セット", ImageURL: "https://cdn.pixabay.com/photo/2016/11/29/12/45/dishes-1869649_1280.jpg"},
		{ID: 82, Name: "コップ", Price: 880, Description: "透明で美しいガラスのコップ", ImageURL: "https://cdn.pixabay.com/photo/2017/02/15/12/12/glass-2069067_1280.jpg"},
		{ID: 83, Name: "急須", Price: 2480, Description: "伝統的な日本の急須", ImageURL: "https://cdn.pixabay.com/photo/2015/10/12/15/46/tea-984367_1280.jpg"},
		
		// 最終追加商品
		{ID: 84, Name: "ミラー", Price: 2800, Description: "おしゃれな手鏡", ImageURL: "https://cdn.pixabay.com/photo/2016/11/29/09/10/mirror-1868559_1280.jpg"},
		{ID: 85, Name: "靴下", Price: 680, Description: "履き心地の良いコットン靴下", ImageURL: "https://cdn.pixabay.com/photo/2017/08/06/06/49/socks-2589618_1280.jpg"},
		{ID: 86, Name: "手袋", Price: 1280, Description: "暖かい冬用の手袋", ImageURL: "https://cdn.pixabay.com/photo/2016/11/29/01/34/gloves-1867176_1280.jpg"},
		{ID: 87, Name: "傘", Price: 2480, Description: "雨の日も安心の丈夫な傘", ImageURL: "https://cdn.pixabay.com/photo/2016/04/12/21/17/umbrella-1325558_1280.jpg"},
		{ID: 88, Name: "懐中電灯", Price: 1980, Description: "停電時に便利なLED懐中電灯", ImageURL: "https://cdn.pixabay.com/photo/2017/01/13/01/22/flashlight-1976086_1280.jpg"},
		{ID: 89, Name: "電卓", Price: 1480, Description: "仕事・学習に便利な電卓", ImageURL: "https://cdn.pixabay.com/photo/2017/01/13/01/22/calculator-1976081_1280.jpg"},
		{ID: 90, Name: "USBメモリ", Price: 1680, Description: "大容量32GBのUSBメモリ", ImageURL: "https://cdn.pixabay.com/photo/2017/01/13/01/22/usb-stick-1976084_1280.jpg"},
		{ID: 91, Name: "マウス", Price: 2800, Description: "使いやすいワイヤレスマウス", ImageURL: "https://cdn.pixabay.com/photo/2017/05/24/21/33/workplace-2340933_1280.jpg"},
		{ID: 92, Name: "キーボード", Price: 4800, Description: "打ちやすいメカニカルキーボード", ImageURL: "https://cdn.pixabay.com/photo/2016/11/29/06/18/home-office-1867761_1280.jpg"},
		{ID: 93, Name: "ヘッドホン", Price: 6800, Description: "高音質オーバーイヤーヘッドホン", ImageURL: "https://cdn.pixabay.com/photo/2017/05/24/21/33/headphones-2340932_1280.jpg"},
		{ID: 94, Name: "スピーカー", Price: 8800, Description: "Bluetooth対応ポータブルスピーカー", ImageURL: "https://cdn.pixabay.com/photo/2017/08/07/05/84/speaker-2603591_1280.jpg"},
		{ID: 95, Name: "デスクライト", Price: 3800, Description: "目に優しいLEDデスクライト", ImageURL: "https://cdn.pixabay.com/photo/2016/11/21/12/51/light-1845538_1280.jpg"},
		{ID: 96, Name: "ペン立て", Price: 980, Description: "机上整理に便利なペン立て", ImageURL: "https://cdn.pixabay.com/photo/2017/01/13/01/22/pencil-holder-1976088_1280.jpg"},
		{ID: 97, Name: "ゴミ箱", Price: 1880, Description: "おしゃれなインテリアゴミ箱", ImageURL: "https://cdn.pixabay.com/photo/2017/01/13/01/22/trash-can-1976089_1280.jpg"},
		{ID: 98, Name: "洗剤", Price: 680, Description: "環境に優しい天然洗剤", ImageURL: "https://cdn.pixabay.com/photo/2017/03/30/15/47/cleaning-2188370_1280.jpg"},
		{ID: 99, Name: "ティッシュ", Price: 480, Description: "やわらかいボックスティッシュ", ImageURL: "https://cdn.pixabay.com/photo/2017/03/30/15/47/tissues-2188367_1280.jpg"},
		{ID: 100, Name: "マスク", Price: 580, Description: "快適な不織布マスク50枚入り", ImageURL: "https://cdn.pixabay.com/photo/2020/04/29/07/54/coronavirus-5107715_1280.jpg"},
	}
}