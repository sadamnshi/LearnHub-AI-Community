package repositories

import (
	"gin_demo/databases"
	"gin_demo/models"
)

// PostRepository 帖子数据访问接口
// 这个接口定义了对帖子数据的所有操作方法（增删改查）
// 作用：为服务层提供统一的数据访问入口，隐藏具体的数据库操作细节
// 好处：可以方便地更换数据库实现（比如从 PostgreSQL 换成 MySQL），而不影响上层业务逻辑
type PostRepository interface {
	// List 分页查询帖子列表，可按分类筛选
	// 参数说明：
	//   - page: 页码，从 1 开始（第 1 页、第 2 页...）
	//   - pageSize: 每页显示多少条数据（比如 20 条）
	//   - categoryID: 分类 ID，如果为 0 表示不按分类筛选（查询所有分类的帖子）
	// 返回值说明：
	//   - []models.Post: 查询到的帖子列表
	//   - int64: 总共有多少条帖子（用于计算总页数）
	//   - error: 如果出错会返回错误信息，否则为 nil
	List(page, pageSize int, categoryID uint) ([]models.Post, int64, error)

	// FindByID 根据帖子 ID 查询单个帖子的完整详情
	// 参数说明：
	//   - id: 帖子的 ID
	// 返回值说明：
	//   - *models.Post: 指向帖子结构体的指针（如果找不到会返回 nil）
	//   - error: 如果出错会返回错误信息
	FindByID(id uint) (*models.Post, error)
}

// postRepository 帖子数据访问的具体实现
// 注意：结构体名字首字母小写 (postRepository)，这表示私有结构体，只能在本包内使用
// 目前这个结构体没有字段，因为所有数据库操作都通过全局的 databases.DB 连接进行
type postRepository struct{}

// NewPostRepository 工厂函数，用来创建 PostRepository 实例
// 设计模式：工厂模式 - 统一创建对象的方式
// 返回值：返回 PostRepository 接口类型（不是具体的 postRepository 实现），这样调用者只需关心接口，不需关心具体实现
func NewPostRepository() PostRepository {
	return &postRepository{}
}

// List 分页查询帖子列表 - 这是 PostRepository 接口的具体实现
// 作用：从数据库查询帖子，支持分页和按分类筛选
// 工作流程：
//  1. 构建查询条件（只查询已发布的帖子）
//  2. 如果指定了分类，添加分类筛选条件
//  3. 统计符合条件的总帖子数
//  4. 执行分页查询（跳过前几条，取指定数量）
//  5. 返回帖子列表和总数
func (r *postRepository) List(page, pageSize int, categoryID uint) ([]models.Post, int64, error) {
	// 声明两个变量来存储查询结果
	var posts []models.Post // 用来存储查询到的帖子列表
	var total int64         // 用来存储符合条件的总帖子数（int64 是为了处理大数据量）

	// ========== 第一步：构建基础查询 ==========
	// 解释：
	//   - databases.DB: 全局数据库连接（在 databases/database.go 中初始化）
	//   - Model(&models.Post{}): 指定要查询的表（Post 对应数据库的 posts 表）
	//   - Where("status = ?", "published"): 添加条件 - 只查询状态为 "published"（已发布）的帖子
	//   - 这里使用了参数化查询 (?) 来防止 SQL 注入攻击，第二个参数 "published" 会自动填入 ?
	query := databases.DB.Model(&models.Post{}).Where("status = ?", "published")

	// ========== 第二步：按分类筛选（可选） ==========
	// 解释：
	//   - 如果 categoryID > 0，说明用户要按某个分类筛选（比如只看 "技术分享" 分类）
	//   - 如果 categoryID == 0，说明不按分类筛选（查询所有分类的帖子）
	// 按分类筛选（categoryID=0 时不过滤）
	if categoryID > 0 {
		// 在已有的查询条件基础上，再添加一个条件
		// 现在的条件变成了：status = 'published' AND category_id = ?
		query = query.Where("category_id = ?", categoryID)
	}

	// ========== 第三步：统计总数 ==========
	// 解释：
	//   - Count(&total): 统计符合条件的记录数
	//   - .Error: 获取执行过程中的错误信息（如果有的话）
	//   - 如果查询出错，立即返回错误，不继续执行后续代码
	// 先统计总数
	if err := query.Count(&total).Error; err != nil {
		// 如果出错，返回空列表、0 条总数、和错误信息
		return nil, 0, err
	}

	// ========== 第四步：分页计算 ==========
	// 解释：假设 page=2, pageSize=20
	//   - page=2 表示第 2 页
	//   - offset = (2-1) * 20 = 20
	//   - 意思是：跳过前 20 条记录，然后取 20 条
	//   - 这样就得到了第 21-40 条记录（第 2 页的内容）
	offset := (page - 1) * pageSize

	// ========== 第五步：执行分页查询 ==========
	err := query.
		// Preload 的作用：预加载关联的数据，避免 N+1 查询问题
		// 说白了，就是在查询帖子的同时，也把关联的作者和分类信息一并查出来

		// Preload("Author") - 自动加载每个帖子的作者信息
		// 不加这行的话，帖子对象中的 Author 字段会是空的，需要再单独查询
		// 加上这行后，GORM 会帮我们自动做一个 JOIN 查询，一次性把所有数据都查出来
		Preload("Author"). // 预加载作者信息，避免 N+1 查询

		// Preload("Category") - 自动加载每个帖子的分类信息
		// 原理同上
		Preload("Category"). // 预加载分类信息

		// Order("created_at DESC") - 按创建时间降序排列
		// DESC = Descending = 降序
		// created_at DESC: 最新的帖子排在前面
		// 如果用 ASC（升序），最旧的帖子就排在前面
		Order("created_at DESC").

		// Limit(pageSize) - 限制返回的记录数
		// pageSize = 20，意思是最多返回 20 条记录
		Limit(pageSize).

		// Offset(offset) - 跳过前 offset 条记录
		// 比如 offset=20，就是跳过前 20 条，从第 21 条开始
		// 配合 Limit 使用就实现了分页功能
		Offset(offset).

		// Find(&posts) - 执行查询，把结果存到 posts 变量中
		// & 表示取地址（传指针），因为我们要修改 posts 这个变量，所以要传指针
		// .Error - 获取查询是否出错
		Find(&posts).Error

	// 如果查询出错，返回错误信息
	if err != nil {
		return nil, 0, err
	}

	// ========== 第六步：返回结果 ==========
	// 返回三个值：
	//   1. posts - 查询到的帖子列表
	//   2. total - 符合条件的总帖子数
	//   3. nil - 没有错误，所以错误值为 nil
	return posts, total, nil
}

// FindByID 根据 ID 查询单个帖子的详细信息 - 这是 PostRepository 接口的另一个实现
// 作用：当用户点击某个帖子时，查询这个帖子的完整内容
// 参数：id - 要查询的帖子 ID
// 返回值：帖子对象指针 (*models.Post) 和错误信息
func (r *postRepository) FindByID(id uint) (*models.Post, error) {
	// 声明一个 Post 结构体变量来存储查询结果
	var post models.Post

	// ========== 执行查询 ==========
	err := databases.DB.
		// Preload("Author") - 预加载帖子的作者信息
		// 这样我们不仅能得到帖子内容，还能得到完整的作者信息（username、avatar 等）
		Preload("Author").

		// Preload("Category") - 预加载帖子所属的分类信息
		// 这样我们能直接看到这个帖子属于哪个分类
		Preload("Category").

		// First(&post, id) - 查询匹配条件的第一条记录
		// 参数说明：
		//   - &post: 查询结果存到这个变量里（需要传指针）
		//   - id: 查询条件 - WHERE id = ?
		// First 函数会自动根据传入的第二个参数生成 WHERE 条件
		// 所以这里相当于：SELECT * FROM posts WHERE id = ? LIMIT 1
		First(&post, id).

		// .Error - 获取查询过程中是否出错
		Error

	// ========== 错误处理 ==========
	// 如果查询出错（比如网络问题、数据库异常、ID 不存在等），返回错误
	if err != nil {
		// 返回 nil 表示没查到数据，返回错误信息
		return nil, err
	}

	// ========== 返回结果 ==========
	// &post: 返回 post 变量的地址（指针）
	// nil: 没有错误
	// 为什么要返回指针而不是值？这是 Go 的常见做法：
	//   - 如果是小对象，返回值更高效
	//   - 如果是大对象，返回指针更高效（避免复制整个对象）
	//   - 如果可能不存在，通常返回指针（可以判断是否为 nil）
	return &post, nil
}
