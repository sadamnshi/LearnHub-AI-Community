package services

import (
	"errors"
	"time"

	"gin_demo/models"
	"gin_demo/repositories"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

// UserService 业务接口定义：
// Register 负责注册（创建新用户并哈希密码）
// Login 校验用户名密码合法性（不返回 token，保持职责单一）
// Authenticate 登录并返回 JWT（组合 Login + 签发 Token）
// GenerateToken 可以被 Authenticate 或刷新流程调用。
// 在企业级中还会有：UpdateProfile / ChangePassword / ResetPassword / LockUser / List / Disable 等。
type UserService interface {
	Register(user *models.User) error
	Login(username, password string) (*models.User, error)
	Authenticate(username, password string) (string, *models.User, error)
	GenerateToken(u *models.User) (string, error)
	GetProfile(userID uint) (*models.User, error)
	UpdatePassword(userID uint, oldPassword, newPassword string) error
}

// userService 实现 UserService，聚合仓储。
type userService struct {
	repo repositories.UserRepository
	// 可以扩展：logger, cache, config 等
}

// NewUserService 构造函数，用于依赖注入。
func NewUserService(repo repositories.UserRepository) UserService {
	return &userService{repo: repo}
}

// hashPassword 使用 bcrypt 生成哈希。成本（cost）默认采用 bcrypt.DefaultCost。
// 企业级注意：
// 1. 成本值可以通过配置文件调整，示例使用默认值便于演示。
// 2. 密码哈希不能重复使用旧值，不能使用 MD5/SHA1 等不可逆、弱算法。
// 3. 不在服务层返回哈希后的密码给调用方。
func hashPassword(plain string) (string, error) {
	//bcrypt哈希算法函数接收一个字节数组，所以需要将传入的参数强转为字节数组，返回的也是一个字节数组，所以在return是强转回字符串
	h, err := bcrypt.GenerateFromPassword([]byte(plain), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(h), nil
}

// verifyPassword 校验明文与哈希是否匹配。
// 第一个参数是根据传入的用户名从数据库中查找的哈希值，第二个参数是明文密码
func verifyPassword(hashed, plain string) bool {
	return bcrypt.CompareHashAndPassword([]byte(hashed), []byte(plain)) == nil
}

// Register 注册新用户：
// 1. 检查用户名是否已存在。
// 2. bcrypt 哈希密码。
// 3. 持久化保存。
// 错误策略：
// - 已存在返回业务错误，不暴露底层 DB 错误细节（可包装）。
// - 哈希失败返回内部错误。
func (s *userService) Register(user *models.User) error {

	if existing, _ := s.repo.FindByUsername(user.Username); existing != nil {
		return errors.New("username already exists")
	}
	//生成哈希
	h, err := hashPassword(user.Password)
	//判断哈希是否出错
	if err != nil {
		return errors.New("internal password hash error")
	}
	//将生成的哈希赋值给用户密码
	user.Password = h
	//写入数据库
	return s.repo.Create(user)
}

// Login 只做身份校验，不签发 Token，便于被组合复用。
func (s *userService) Login(username, password string) (*models.User, error) {
	u, err := s.repo.FindByUsername(username)
	if err != nil {
		return nil, errors.New("user not found")
	}
	if !verifyPassword(u.Password, password) {
		return nil, errors.New("invalid password")
	}
	return u, nil
}

// jwtSecret 演示用静态秘钥。生产环境：
// 1. 使用环境变量或配置文件加载；
// 2. 采用长度足够的随机字符串；
// 3. 支持密钥轮换（Key Rotation）。
var jwtSecret = []byte("CHANGE_ME_TO_ENV_SECRET")

// GenerateToken 根据用户信息签发 JWT。
// Claims 设计：
//   - sub: 主题，使用用户 ID
//   - name: 用户名
//   - exp: 过期时间（短 Token，一般 15m ~ 2h）
//   - iat: 签发时间
//   - nbf: 生效时间
//   - iss: 签发者（示例写死，可配置）
//
// 企业级还可添加：
//   - jti: 唯一 ID 便于做 Token 黑名单撤销
//   - roles / permissions: RBAC 权限控制
//   - trace 信息：便于链路追踪
//
// 生成jwt
func (s *userService) GenerateToken(u *models.User) (string, error) {

	claims := jwt.MapClaims{
		"sub":  u.ID,
		"name": u.Username,
		"exp":  time.Now().Add(2 * time.Hour).Unix(), // 过期时间
		"iat":  time.Now().Unix(),
		"nbf":  time.Now().Unix(),
		"iss":  "gin_demo_app",
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(jwtSecret)
}

// Authenticate = Login + GenerateToken
func (s *userService) Authenticate(username, password string) (string, *models.User, error) {
	u, err := s.Login(username, password)
	if err != nil {
		return "", nil, err
	}
	t, err := s.GenerateToken(u)
	if err != nil {
		return "", nil, errors.New("token generation failed")
	}
	return t, u, nil
}

// GetProfile 根据用户 ID 从数据库查询完整用户信息（不含密码）
func (s *userService) GetProfile(userID uint) (*models.User, error) {
	u, err := s.repo.FindByID(userID)
	if err != nil {
		return nil, errors.New("user not found")
	}
	u.Password = "" // 绝不返回密码哈希
	return u, nil
}

// UpdatePassword 修改密码业务逻辑：
// 1. 查出当前用户（需要旧的哈希密码用于校验）
// 2. 验证旧密码是否正确
// 3. 检查新密码与旧密码不能相同
// 4. 对新密码哈希后写入数据库
func (s *userService) UpdatePassword(userID uint, oldPassword, newPassword string) error {
	// 第一步：从数据库查出当前用户，拿到旧的哈希密码
	u, err := s.repo.FindByID(userID)
	if err != nil {
		return errors.New("user not found")
	}

	// 第二步：用 bcrypt 校验旧密码是否与数据库中的哈希匹配
	// 如果旧密码输入错误，直接拒绝，防止别人在拿到登录态后乱改密码
	if !verifyPassword(u.Password, oldPassword) {
		return errors.New("旧密码错误")
	}

	// 第三步：新密码不能与旧密码相同
	if verifyPassword(u.Password, newPassword) {
		return errors.New("新密码不能与旧密码相同")
	}

	// 第四步：对新密码进行哈希，绝不明文存库
	newHashed, err := hashPassword(newPassword)
	if err != nil {
		return errors.New("密码加密失败")
	}

	// 第五步：调用 Repository 将新哈希写入数据库
	return s.repo.UpdatePassword(userID, newHashed)
}
